package builder

import (
	"fmt"

	"github.com/kris-nova/logger"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

// NodeGroupResourceSet stores the resource information of the nodegroup
type NodeGroupResourceSet struct {
	rs                 *resourceSet
	clusterSpec        *api.ClusterConfig
	spec               *api.NodeGroup
	provider           api.ClusterProvider
	clusterStackName   string
	nodeGroupName      string
	instanceProfileARN *gfn.Value
	securityGroups     []*gfn.Value
	vpc                *gfn.Value
	userData           *gfn.Value
}

// NewNodeGroupResourceSet returns a resource set for a nodegroup embedded in a cluster config
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

// AddAllResources adds all the information about the nodegroup to the resource set
func (n *NodeGroupResourceSet) AddAllResources() error {
	n.rs.template.Description = fmt.Sprintf(
		"%s (AMI family: %s, SSH access: %v, private networking: %v) %s",
		nodeGroupTemplateDescription,
		n.spec.AMIFamily, api.IsEnabled(n.spec.SSH.Allow), n.spec.PrivateNetworking,
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
	launchTemplateName := gfn.MakeFnSubString(fmt.Sprintf("${%s}", gfn.StackName))
	launchTemplateData := newLaunchTemplateData(n)

	if api.IsSetAndNonEmptyString(n.spec.SSH.PublicKeyName) {
		launchTemplateData.KeyName = gfn.NewString(*n.spec.SSH.PublicKeyName)
	}

	if volumeSize := n.spec.VolumeSize; volumeSize != nil && *volumeSize > 0 {
		var (
			kmsKeyID   *gfn.Value
			volumeIOPS *gfn.Value
		)
		if api.IsSetAndNonEmptyString(n.spec.VolumeKmsKeyID) {
			kmsKeyID = gfn.NewString(*n.spec.VolumeKmsKeyID)
		}

		if *n.spec.VolumeType == api.NodeVolumeTypeIO1 {
			volumeIOPS = gfn.NewInteger(*n.spec.VolumeIOPS)
		}

		launchTemplateData.BlockDeviceMappings = []gfn.AWSEC2LaunchTemplate_BlockDeviceMapping{{
			DeviceName: gfn.NewString(*n.spec.VolumeName),
			Ebs: &gfn.AWSEC2LaunchTemplate_Ebs{
				VolumeSize: gfn.NewInteger(*volumeSize),
				VolumeType: gfn.NewString(*n.spec.VolumeType),
				Encrypted:  gfn.NewBoolean(*n.spec.VolumeEncrypted),
				KmsKeyId:   kmsKeyID,
				Iops:       volumeIOPS,
			},
		}}
	}

	n.newResource("NodeGroupLaunchTemplate", &gfn.AWSEC2LaunchTemplate{
		LaunchTemplateName: launchTemplateName,
		LaunchTemplateData: launchTemplateData,
	})

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
	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.AutoScaler) {
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

	asg := nodeGroupResource(launchTemplateName, &vpcZoneIdentifier, tags, n.spec)
	n.newResource("NodeGroup", asg)

	return nil
}

// GetAllOutputs collects all outputs of the nodegroup
func (n *NodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack)
}

func newLaunchTemplateData(n *NodeGroupResourceSet) *gfn.AWSEC2LaunchTemplate_LaunchTemplateData {
	launchTemplateData := &gfn.AWSEC2LaunchTemplate_LaunchTemplateData{
		IamInstanceProfile: &gfn.AWSEC2LaunchTemplate_IamInstanceProfile{
			Arn: n.instanceProfileARN,
		},
		ImageId:  gfn.NewString(n.spec.AMI),
		UserData: n.userData,
		NetworkInterfaces: []gfn.AWSEC2LaunchTemplate_NetworkInterface{{
			AssociatePublicIpAddress: gfn.NewBoolean(!n.spec.PrivateNetworking),
			DeviceIndex:              gfn.NewInteger(0),
			Groups:                   n.securityGroups,
		}},
	}
	if !api.HasMixedInstances(n.spec) {
		launchTemplateData.InstanceType = gfn.NewString(n.spec.InstanceType)
	} else {
		launchTemplateData.InstanceType = gfn.NewString(n.spec.InstancesDistribution.InstanceTypes[0])
	}
	if n.spec.EBSOptimized != nil {
		launchTemplateData.EbsOptimized = gfn.NewBoolean(*n.spec.EBSOptimized)
	}

	return launchTemplateData
}

func nodeGroupResource(launchTemplateName *gfn.Value, vpcZoneIdentifier *interface{}, tags []map[string]interface{}, ng *api.NodeGroup) *awsCloudFormationResource {
	ngProps := map[string]interface{}{
		"VPCZoneIdentifier": *vpcZoneIdentifier,
		"Tags":              tags,
	}
	if ng.DesiredCapacity != nil {
		ngProps["DesiredCapacity"] = fmt.Sprintf("%d", *ng.DesiredCapacity)
	}
	if ng.MinSize != nil {
		ngProps["MinSize"] = fmt.Sprintf("%d", *ng.MinSize)
	}
	if ng.MaxSize != nil {
		ngProps["MaxSize"] = fmt.Sprintf("%d", *ng.MaxSize)
	}
	if len(ng.TargetGroupARNs) > 0 {
		ngProps["TargetGroupARNs"] = ng.TargetGroupARNs
	}
	if api.HasMixedInstances(ng) {
		ngProps["MixedInstancesPolicy"] = *mixedInstancesPolicy(launchTemplateName, ng)
	} else {
		ngProps["LaunchTemplate"] = map[string]interface{}{
			"LaunchTemplateName": launchTemplateName,
			"Version":            gfn.MakeFnGetAttString("NodeGroupLaunchTemplate.LatestVersionNumber"),
		}
	}

	return &awsCloudFormationResource{
		Type:       "AWS::AutoScaling::AutoScalingGroup",
		Properties: ngProps,
		UpdatePolicy: map[string]map[string]string{
			"AutoScalingRollingUpdate": {
				"MinInstancesInService": "0",
				"MaxBatchSize":          "1",
			},
		},
	}
}

func mixedInstancesPolicy(launchTemplateName *gfn.Value, ng *api.NodeGroup) *map[string]interface{} {
	instanceTypes := ng.InstancesDistribution.InstanceTypes
	overrides := make([]map[string]string, len(instanceTypes))

	for i, instanceType := range instanceTypes {
		overrides[i] = map[string]string{
			"InstanceType": instanceType,
		}
	}
	policy := map[string]interface{}{
		"LaunchTemplate": map[string]interface{}{
			"LaunchTemplateSpecification": map[string]interface{}{
				"LaunchTemplateName": launchTemplateName,
				"Version":            gfn.MakeFnGetAttString("NodeGroupLaunchTemplate.LatestVersionNumber"),
			},

			"Overrides": overrides,
		},
	}

	instancesDistribution := map[string]string{}

	// Only set the price if it was specified so otherwise AWS picks "on-demand price" as the default
	if ng.InstancesDistribution.MaxPrice != nil {
		instancesDistribution["SpotMaxPrice"] = fmt.Sprintf("%f", *ng.InstancesDistribution.MaxPrice)
	}
	if ng.InstancesDistribution.OnDemandBaseCapacity != nil {
		instancesDistribution["OnDemandBaseCapacity"] = fmt.Sprintf("%d", *ng.InstancesDistribution.OnDemandBaseCapacity)
	}
	if ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity != nil {
		instancesDistribution["OnDemandPercentageAboveBaseCapacity"] = fmt.Sprintf("%d", *ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity)
	}
	if ng.InstancesDistribution.SpotInstancePools != nil {
		instancesDistribution["SpotInstancePools"] = fmt.Sprintf("%d", *ng.InstancesDistribution.SpotInstancePools)
	}

	policy["InstancesDistribution"] = instancesDistribution

	return &policy
}
