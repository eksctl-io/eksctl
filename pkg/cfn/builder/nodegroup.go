package builder

import (
	"fmt"

	"github.com/spotinst/spotinst-sdk-go/spotinst/credentials"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

// NodeGroupResourceSet stores the resource information of the node group
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
	launchTemplateData := n.newLaunchTemplateData()

	if api.IsEnabled(n.spec.SSH.Allow) && api.IsSetAndNonEmptyString(n.spec.SSH.PublicKeyName) {
		launchTemplateData.KeyName = gfn.NewString(*n.spec.SSH.PublicKeyName)
	}

	if volumeSize := n.spec.VolumeSize; volumeSize != nil && *volumeSize > 0 {
		launchTemplateData.BlockDeviceMappings = []gfn.AWSEC2LaunchTemplate_BlockDeviceMapping{{
			DeviceName: gfn.NewString(*n.spec.VolumeName),
			Ebs: &gfn.AWSEC2LaunchTemplate_Ebs{
				VolumeSize: gfn.NewInteger(*volumeSize),
				VolumeType: gfn.NewString(*n.spec.VolumeType),
				Encrypted:  gfn.NewBoolean(*n.spec.VolumeEncrypted),
			},
		}}
		if api.IsSetAndNonEmptyString(n.spec.VolumeKmsKeyID) {
			launchTemplateData.BlockDeviceMappings[0].Ebs.KmsKeyId = gfn.NewString(*n.spec.VolumeKmsKeyID)
		}
	}

	launchTemplate := &gfn.AWSEC2LaunchTemplate{
		LaunchTemplateName: launchTemplateName,
		LaunchTemplateData: launchTemplateData,
	}
	n.newResource("NodeGroupLaunchTemplate", launchTemplate)

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

	g, err := n.newNodeGroupResource(launchTemplate, &vpcZoneIdentifier, tags)
	if err != nil {
		return fmt.Errorf("failed to build node group resource: %v", err)
	}
	n.newResource("NodeGroup", g)

	return nil
}

// GetAllOutputs collects all outputs of the node group
func (n *NodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack)
}

func (n *NodeGroupResourceSet) newLaunchTemplateData() *gfn.AWSEC2LaunchTemplate_LaunchTemplateData {
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
	}

	return launchTemplateData
}

func (n *NodeGroupResourceSet) newNodeGroupResource(launchTemplate *gfn.AWSEC2LaunchTemplate,
	vpcZoneIdentifier *interface{}, tags []map[string]interface{}) (*awsCloudFormationResource, error) {

	if n.spec.Spotinst != nil {
		logger.Debug("creating nodegroup using spotinst ocean")
		return n.newNodeGroupSpotinstResource(launchTemplate, vpcZoneIdentifier, tags)
	} else {
		logger.Debug("creating nodegroup using aws auto scaling group")
		return n.newNodeGroupAutoScalingGroupResource(launchTemplate, vpcZoneIdentifier, tags)
	}
}

func (n *NodeGroupResourceSet) newNodeGroupAutoScalingGroupResource(launchTemplate *gfn.AWSEC2LaunchTemplate,
	vpcZoneIdentifier *interface{}, tags []map[string]interface{}) (*awsCloudFormationResource, error) {

	ng := n.spec
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
		ngProps["MixedInstancesPolicy"] = *n.newMixedInstancesPolicy(launchTemplate.LaunchTemplateName, ng)
	} else {
		ngProps["LaunchTemplate"] = map[string]interface{}{
			"LaunchTemplateName": launchTemplate.LaunchTemplateName,
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
	}, nil
}

func (n *NodeGroupResourceSet) newNodeGroupSpotinstResource(launchTemplate *gfn.AWSEC2LaunchTemplate,
	vpcZoneIdentifier *interface{}, tags []map[string]interface{}) (*awsCloudFormationResource, error) {

	var resource *awsCloudFormationResource
	var err error

	// Resource.
	{
		if oceanID := n.spec.Spotinst.Ocean.ID; oceanID == nil {
			logger.Debug("creating a new spotinst ocean cluster")
			resource, err = n.newNodeGroupSpotinstOceanResource(
				launchTemplate, vpcZoneIdentifier, tags)
		} else {
			logger.Debug("joining to an existing spotinst ocean cluster %q", *oceanID)
			resource, err = n.newNodeGroupSpotinstOceanLaunchSpecResource(
				launchTemplate, vpcZoneIdentifier, tags, *oceanID)
		}

		if err != nil {
			return nil, err
		}
	}

	// Credentials.
	{
		resource.Properties["ServiceToken"] = gfn.MakeFnSubString("arn:aws:lambda:${AWS::Region}:178579023202:function:spotinst-cloudformation")

		provider := credentials.NewChainCredentials(
			new(credentials.EnvProvider),
			new(credentials.FileProvider))

		creds, err := provider.Get()
		if err != nil {
			return nil, err
		}
		if creds.Token != "" {
			resource.Properties["accessToken"] = creds.Token
		}
		if creds.Account != "" {
			resource.Properties["accountId"] = creds.Account
		}
	}

	return resource, nil
}

func (n *NodeGroupResourceSet) newNodeGroupSpotinstOceanResource(launchTemplate *gfn.AWSEC2LaunchTemplate,
	vpcZoneIdentifier *interface{}, tags []map[string]interface{}) (*awsCloudFormationResource, error) {

	ng := n.spec

	ngProps := make(map[string]interface{})
	ocProps := make(map[string]interface{})

	ocProps["name"] = n.clusterSpec.Metadata.Name
	ocProps["controllerClusterId"] = n.clusterSpec.Metadata.Name
	ocProps["region"] = gfn.MakeRef("AWS::Region")
	ngProps["ocean"] = ocProps

	// Strategy.
	{
		strategyProps := make(map[string]interface{})
		strategyProps["spotPercentage"] = 100
		strategyProps["fallbackToOd"] = true
		strategyProps["utilizeReservedInstances"] = true

		if ng.Spotinst.Strategy != nil {
			if v := ng.Spotinst.Strategy.SpotPercentage; v != nil {
				strategyProps["spotPercentage"] = *v
			}

			if v := ng.Spotinst.Strategy.FallbackToOnDemand; v != nil {
				strategyProps["fallbackToOd"] = *v
			}

			if v := ng.Spotinst.Strategy.UtilizeReservedInstances; v != nil {
				strategyProps["utilizeReservedInstances"] = *v
			}
		}

		if len(strategyProps) > 0 {
			ocProps["strategy"] = strategyProps
		}
	}

	// Capacity.
	{
		capacityProps := make(map[string]interface{})

		if ng.DesiredCapacity != nil {
			capacityProps["target"] = fmt.Sprintf("%d", *ng.DesiredCapacity)
		}

		if ng.MinSize != nil {
			capacityProps["minimum"] = fmt.Sprintf("%d", *ng.MinSize)
		}

		if ng.MaxSize != nil {
			capacityProps["maximum"] = fmt.Sprintf("%d", *ng.MaxSize)
		}

		if len(capacityProps) > 0 {
			ocProps["capacity"] = capacityProps
		}
	}

	// Compute.
	{
		computeProps := make(map[string]interface{})

		// Subnet IDs.
		{
			if vpcZoneIdentifier != nil {
				computeProps["subnetIds"] = *vpcZoneIdentifier
			}
		}

		// Launch Specification.
		{
			launchSpecProps := make(map[string]interface{})
			launchTemplateData := launchTemplate.LaunchTemplateData

			if launchTemplateData.ImageId != nil {
				launchSpecProps["imageId"] = launchTemplateData.ImageId
			}

			if launchTemplateData.UserData != nil {
				launchSpecProps["userData"] = launchTemplateData.UserData
			}

			if launchTemplateData.KeyName != nil {
				launchSpecProps["keyPair"] = launchTemplateData.KeyName
			}

			if launchTemplateData.IamInstanceProfile != nil {
				launchSpecProps["iamInstanceProfile"] = map[string]interface{}{
					"arn": launchTemplateData.IamInstanceProfile.Arn,
				}
			}

			if launchTemplateData.NetworkInterfaces != nil {
				if launchTemplateData.NetworkInterfaces[0].AssociatePublicIpAddress != nil {
					launchSpecProps["associatePublicIpAddress"] = launchTemplateData.NetworkInterfaces[0].AssociatePublicIpAddress
				}

				if len(launchTemplateData.NetworkInterfaces[0].Groups) > 0 {
					launchSpecProps["securityGroupIds"] = launchTemplateData.NetworkInterfaces[0].Groups
				}
			}

			if len(ng.TargetGroupARNs) > 0 {
				targetGroups := make([]interface{}, len(ng.TargetGroupARNs))

				for i, arn := range ng.TargetGroupARNs {
					targetGroups[i] = map[string]interface{}{
						"type": "TARGET_GROUP",
						"arn":  arn,
					}
				}

				launchSpecProps["loadBalancers"] = targetGroups
			}

			if len(tags) > 0 {
				tagsKV := make([]map[string]interface{}, len(tags))

				for i, tag := range tags {
					tagsKV[i] = map[string]interface{}{
						"tagKey":   tag["Key"],
						"tagValue": tag["Value"],
					}
				}

				launchSpecProps["tags"] = tagsKV
			}

			if len(launchSpecProps) > 0 {
				computeProps["launchSpecification"] = launchSpecProps
			}
		}

		if len(computeProps) > 0 {
			ocProps["compute"] = computeProps
		}
	}

	// Auto Scaler.
	{
		autoScalerProps := make(map[string]interface{})
		autoScalerProps["isEnabled"] = true
		autoScalerProps["isAutoConfig"] = true

		if ng.Spotinst.AutoScaler != nil {
			if v := ng.Spotinst.AutoScaler.Enabled; v != nil {
				autoScalerProps["isEnabled"] = *v
			}

			if v := ng.Spotinst.AutoScaler.AutoConfig; v != nil {
				autoScalerProps["isAutoConfig"] = *v
			}

			if v := ng.Spotinst.AutoScaler.Cooldown; v != nil {
				autoScalerProps["cooldown"] = *v
			}

			if ng.Spotinst.AutoScaler.Headroom != nil {
				headroomProps := make(map[string]interface{})

				if v := ng.Spotinst.AutoScaler.Headroom.CPUPerUnit; v != nil {
					headroomProps["cpuPerUnit"] = *v
				}

				if v := ng.Spotinst.AutoScaler.Headroom.GPUPerUnit; v != nil {
					headroomProps["gpuPerUnit"] = *v
				}

				if v := ng.Spotinst.AutoScaler.Headroom.MemPerUnit; v != nil {
					headroomProps["memoryPerUnit"] = *v
				}

				if v := ng.Spotinst.AutoScaler.Headroom.NumOfUnits; v != nil {
					headroomProps["numOfUnits"] = *v
				}

				if len(headroomProps) > 0 {
					autoScalerProps["headroom"] = headroomProps
				}
			}
		}

		if len(autoScalerProps) > 0 {
			ocProps["autoScaler"] = autoScalerProps
		}
	}

	return &awsCloudFormationResource{
		Type:       "Custom::ocean",
		Properties: ngProps,
	}, nil
}

func (n *NodeGroupResourceSet) newNodeGroupSpotinstOceanLaunchSpecResource(launchTemplate *gfn.AWSEC2LaunchTemplate,
	vpcZoneIdentifier *interface{}, tags []map[string]interface{}, oceanID string) (*awsCloudFormationResource, error) {

	ngProps := make(map[string]interface{})
	lsProps := make(map[string]interface{})

	lsProps["oceanId"] = oceanID
	ngProps["oceanLaunchSpec"] = lsProps

	launchTemplateData := launchTemplate.LaunchTemplateData

	if launchTemplateData.ImageId != nil {
		lsProps["imageId"] = launchTemplateData.ImageId
	}

	if launchTemplateData.UserData != nil {
		lsProps["userData"] = launchTemplateData.UserData
	}

	if launchTemplateData.IamInstanceProfile != nil {
		lsProps["iamInstanceProfile"] = map[string]interface{}{
			"arn": launchTemplateData.IamInstanceProfile.Arn,
		}
	}

	if launchTemplateData.NetworkInterfaces != nil {
		if len(launchTemplateData.NetworkInterfaces[0].Groups) > 0 {
			lsProps["securityGroupIds"] = launchTemplateData.NetworkInterfaces[0].Groups
		}
	}

	return &awsCloudFormationResource{
		Type:       "Custom::oceanLaunchSpec",
		Properties: ngProps,
	}, nil
}

func (n *NodeGroupResourceSet) newMixedInstancesPolicy(launchTemplateName *gfn.Value, ng *api.NodeGroup) *map[string]interface{} {
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
