package builder

import (
	"fmt"

	"github.com/kris-nova/logger"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

// NodeGroupResourceSet stores the resource information of the node group
type NodeGroupResourceSet struct {
	rs               *resourceSet
	clusterSpec      *api.ClusterConfig
	spec             *api.NodeGroup
	provider         api.ClusterProvider
	clusterStackName string
	nodeGroupName    string
	instanceProfile  *gfn.Value
	securityGroups   []*gfn.Value
	vpc              *gfn.Value
	userData         *gfn.Value
}

// NewNodeGroupResourceSet returns a resource set for a node group embedded in a cluster config
func NewNodeGroupResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig, clusterStackName string, ng *api.NodeGroup) *NodeGroupResourceSet {
	return &NodeGroupResourceSet{
		rs:               newResourceSet(),
		clusterStackName: clusterStackName,
		nodeGroupName:    ng.Name,
		clusterSpec:      spec,
		spec:             ng,
		provider:         provider,
	}
}

// AddAllResources adds all the information about the node group to the resource set
func (n *NodeGroupResourceSet) AddAllResources() error {
	n.rs.template.Description = fmt.Sprintf(
		"%s (AMI family: %s, SSH access: %v, private networking: %v) %s",
		nodeGroupTemplateDescription,
		n.spec.AMIFamily, n.spec.SSH.Allow, n.spec.PrivateNetworking,
		templateDescriptionSuffix)

	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeaturePrivateNetworking, n.spec.PrivateNetworking, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureSharedSecurityGroup, n.spec.SecurityGroups.WithShared, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureLocalSecurityGroup, n.spec.SecurityGroups.WithLocal, false)

	n.vpc = makeImportValue(n.clusterStackName, outputs.ClusterVPC)

	userData, err := nodebootstrap.NewUserData(n.clusterSpec, n.spec)
	if err != nil {
		return err
	}
	n.userData = gfn.NewString(userData)

	// Ensure MinSize is set, as it is required by the ASG cfn resource
	if n.spec.MinSize == nil {
		if n.spec.DesiredCapacity == nil {
			defaultNodeCount := api.DefaultNodeCount
			n.spec.MinSize = &defaultNodeCount
		} else {
			n.spec.MinSize = n.spec.DesiredCapacity
		}
		logger.Info("--nodes-min=%d was set automatically for nodegroup %s", *n.spec.MinSize, n.nodeGroupName)
	} else if n.spec.DesiredCapacity != nil && *n.spec.DesiredCapacity < *n.spec.MinSize {
		return fmt.Errorf("cannot use --nodes-min=%d and --nodes=%d at the same time", *n.spec.MinSize, *n.spec.DesiredCapacity)
	}

	// Ensure MaxSize is set, as it is required by the ASG cfn resource
	if n.spec.MaxSize == nil {
		if n.spec.DesiredCapacity == nil {
			n.spec.MaxSize = n.spec.MinSize
		} else {
			n.spec.MaxSize = n.spec.DesiredCapacity
		}
		logger.Info("--nodes-max=%d was set automatically for nodegroup %s", *n.spec.MaxSize, n.nodeGroupName)
	} else if n.spec.DesiredCapacity != nil && *n.spec.DesiredCapacity > *n.spec.MaxSize {
		return fmt.Errorf("cannot use --nodes-max=%d and --nodes=%d at the same time", *n.spec.MaxSize, *n.spec.DesiredCapacity)
	} else if *n.spec.MaxSize < *n.spec.MinSize {
		return fmt.Errorf("cannot use --nodes-min=%d and --nodes-max=%d at the same time", *n.spec.MinSize, *n.spec.MaxSize)
	}

	n.addResourcesForIAM()
	n.addResourcesForSecurityGroups()

	return n.addResourcesForNodeGroup()
}

// RenderJSON returns the rendered JSON
func (n *NodeGroupResourceSet) RenderJSON() ([]byte, error) {
	return n.rs.renderJSON()
}

// Template returns the CloudFormation template
func (n *NodeGroupResourceSet) Template() gfn.Template {
	return *n.rs.template
}

func (n *NodeGroupResourceSet) newResource(name string, resource interface{}) *gfn.Value {
	return n.rs.newResource(name, resource)
}

func (n *NodeGroupResourceSet) addResourcesForNodeGroup() error {
	lc := &gfn.AWSAutoScalingLaunchConfiguration{
		IamInstanceProfile: n.instanceProfile,
		SecurityGroups:     n.securityGroups,
		ImageId:            gfn.NewString(n.spec.AMI),
		InstanceType:       gfn.NewString(n.spec.InstanceType),
		UserData:           n.userData,
	}
	if n.spec.SSH.Allow && n.spec.SSH.PublicKeyName != nil {
		lc.KeyName = gfn.NewString(*n.spec.SSH.PublicKeyName)
	}
	if n.spec.PrivateNetworking {
		lc.AssociatePublicIpAddress = gfn.False()
	} else {
		lc.AssociatePublicIpAddress = gfn.True()
	}
	if n.spec.VolumeSize > 0 {
		lc.BlockDeviceMappings = []gfn.AWSAutoScalingLaunchConfiguration_BlockDeviceMapping{
			{
				DeviceName: gfn.NewString("/dev/xvda"),
				Ebs: &gfn.AWSAutoScalingLaunchConfiguration_BlockDevice{
					VolumeSize: gfn.NewInteger(n.spec.VolumeSize),
					VolumeType: gfn.NewString(n.spec.VolumeType),
				},
			},
		}
	}
	refLC := n.newResource("NodeLaunchConfig", lc)
	// currently goformation type system doesn't allow specifying `VPCZoneIdentifier: { "Fn::ImportValue": ... }`,
	// and tags don't have `PropagateAtLaunch` field, so we have a custom method here until this gets resolved
	var vpcZoneIdentifier interface{}
	if numNodeGroupsAZs := len(n.spec.AvailabilityZones); numNodeGroupsAZs > 0 {
		subnets := n.clusterSpec.VPC.Subnets.Private
		if !n.spec.PrivateNetworking {
			subnets = n.clusterSpec.VPC.Subnets.Public
		}
		errorDesc := fmt.Sprintf("(subnets=%#v AZs=%#v)", subnets, n.spec.AvailabilityZones)
		if len(subnets) < numNodeGroupsAZs {
			return fmt.Errorf("VPC doesn't have enough subnets for nodegroup AZs %s", errorDesc)
		}
		vpcZoneIdentifier = make([]interface{}, numNodeGroupsAZs)
		for i, az := range n.spec.AvailabilityZones {
			subnet, ok := subnets[az]
			if !ok {
				return fmt.Errorf("VPC doesn't have subnets in %s %s", az, errorDesc)
			}
			vpcZoneIdentifier.([]interface{})[i] = subnet.ID
		}
	} else {
		subnets := makeImportValue(n.clusterStackName, outputs.ClusterSubnetsPrivate)
		if !n.spec.PrivateNetworking {
			subnets = makeImportValue(n.clusterStackName, outputs.ClusterSubnetsPublic)
		}
		vpcZoneIdentifier = map[string][]interface{}{
			gfn.FnSplit: {",", subnets},
		}
	}
	tags := []map[string]interface{}{
		{
			"Key":               "Name",
			"Value":             fmt.Sprintf("%s-%s-Node", n.clusterSpec.Metadata.Name, n.nodeGroupName),
			"PropagateAtLaunch": "true",
		},
		{
			"Key":               "kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name,
			"Value":             "owned",
			"PropagateAtLaunch": "true",
		},
	}
	if v := n.spec.IAM.WithAddonPolicies.AutoScaler; v != nil && *v {
		tags = append(tags,
			map[string]interface{}{
				"Key":               "k8s.io/cluster-autoscaler/enabled",
				"Value":             "true",
				"PropagateAtLaunch": "true",
			},
			map[string]interface{}{
				"Key":               "k8s.io/cluster-autoscaler/" + n.clusterSpec.Metadata.Name,
				"Value":             "owned",
				"PropagateAtLaunch": "true",
			},
		)
	}
	ngProps := map[string]interface{}{
		"LaunchConfigurationName": refLC,
		"VPCZoneIdentifier":       vpcZoneIdentifier,
		"Tags":                    tags,
	}
	if n.spec.DesiredCapacity != nil {
		ngProps["DesiredCapacity"] = fmt.Sprintf("%d", *n.spec.DesiredCapacity)
	}
	if n.spec.MinSize != nil {
		ngProps["MinSize"] = fmt.Sprintf("%d", *n.spec.MinSize)
	}
	if n.spec.MaxSize != nil {
		ngProps["MaxSize"] = fmt.Sprintf("%d", *n.spec.MaxSize)
	}
	n.newResource("NodeGroup", &awsCloudFormationResource{
		Type:       "AWS::AutoScaling::AutoScalingGroup",
		Properties: ngProps,
		UpdatePolicy: map[string]map[string]string{
			"AutoScalingRollingUpdate": {
				"MinInstancesInService": "1",
				"MaxBatchSize":          "1",
			},
		},
	})

	return nil
}

// GetAllOutputs collects all outputs of the node group
func (n *NodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack)
}
