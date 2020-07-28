package builder

import (
	"encoding/json"
	"fmt"
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/spot"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfnv4 "github.com/weaveworks/goformation/v4/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

// NodeGroupResourceSet stores the resource information of the nodegroup
type NodeGroupResourceSet struct {
	rs                   *resourceSet
	clusterSpec          *api.ClusterConfig
	spec                 *api.NodeGroup
	supportsManagedNodes bool
	provider             api.ClusterProvider
	clusterStackName     string
	instanceProfileARN   *gfnt.Value
	securityGroups       []*gfnt.Value
	vpc                  *gfnt.Value
	userData             *gfnt.Value
	sharedTags           []*cfn.Tag
}

// NewNodeGroupResourceSet returns a resource set for a nodegroup embedded in a cluster config
func NewNodeGroupResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig,
	clusterStackName string, sharedTags []*cfn.Tag, ng *api.NodeGroup,
	supportsManagedNodes bool) *NodeGroupResourceSet {
	return &NodeGroupResourceSet{
		rs:                   newResourceSet(),
		clusterStackName:     clusterStackName,
		supportsManagedNodes: supportsManagedNodes,
		clusterSpec:          spec,
		spec:                 ng,
		provider:             provider,
		sharedTags:           sharedTags,
	}
}

// AddAllResources adds all the information about the nodegroup to the resource set
func (n *NodeGroupResourceSet) AddAllResources() error {
	n.rs.template.Description = fmt.Sprintf(
		"%s (AMI family: %s, SSH access: %v, private networking: %v) %s",
		nodeGroupTemplateDescription,
		n.spec.AMIFamily, api.IsEnabled(n.spec.SSH.Allow), n.spec.PrivateNetworking,
		templateDescriptionSuffix)

	n.Template().Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings

	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeaturePrivateNetworking, n.spec.PrivateNetworking, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureSharedSecurityGroup, n.spec.SecurityGroups.WithShared, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureLocalSecurityGroup, n.spec.SecurityGroups.WithLocal, false)

	n.vpc = makeImportValue(n.clusterStackName, outputs.ClusterVPC)

	userData, err := nodebootstrap.NewUserData(n.clusterSpec, n.spec)
	if err != nil {
		return err
	}
	n.userData = gfnt.NewString(userData)

	// Ensure MinSize is set, as it is required by the ASG cfn resource
	if n.spec.MinSize == nil {
		if n.spec.DesiredCapacity == nil {
			defaultNodeCount := api.DefaultNodeCount
			n.spec.MinSize = &defaultNodeCount
		} else {
			n.spec.MinSize = n.spec.DesiredCapacity
		}
		logger.Info("--nodes-min=%d was set automatically for nodegroup %s", *n.spec.MinSize, n.spec.Name)
	} else if n.spec.DesiredCapacity != nil && *n.spec.DesiredCapacity < *n.spec.MinSize {
		return fmt.Errorf("cannot use --nodes-min=%d and --nodes=%d at the same time", *n.spec.MinSize, *n.spec.DesiredCapacity)
	}

	// Ensure MaxSize is set, as it is required by the ASG cfn resource
	if n.spec.SpotOcean == nil {
		if n.spec.MaxSize == nil {
			if n.spec.DesiredCapacity == nil {
				n.spec.MaxSize = n.spec.MinSize
			} else {
				n.spec.MaxSize = n.spec.DesiredCapacity
			}
			logger.Info("--nodes-max=%d was set automatically for nodegroup %s", *n.spec.MaxSize, n.spec.Name)
		} else if n.spec.DesiredCapacity != nil && *n.spec.DesiredCapacity > *n.spec.MaxSize {
			return fmt.Errorf("cannot use --nodes-max=%d and --nodes=%d at the same time", *n.spec.MaxSize, *n.spec.DesiredCapacity)
		} else if *n.spec.MaxSize < *n.spec.MinSize {
			return fmt.Errorf("cannot use --nodes-min=%d and --nodes-max=%d at the same time", *n.spec.MinSize, *n.spec.MaxSize)
		}
	}

	if err := n.addResourcesForIAM(); err != nil {
		return err
	}
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

func (n *NodeGroupResourceSet) newResource(name string, resource gfn.Resource) *gfnt.Value {
	return n.rs.newResource(name, resource)
}

func (n *NodeGroupResourceSet) addResourcesForNodeGroup() error {
	launchTemplateName := gfnt.MakeFnSubString(fmt.Sprintf("${%s}", gfnt.StackName))
	launchTemplateData := n.newLaunchTemplateData()

	if n.spec.SSH != nil && api.IsSetAndNonEmptyString(n.spec.SSH.PublicKeyName) {
		launchTemplateData.KeyName = gfnt.NewString(*n.spec.SSH.PublicKeyName)
	}

	if volumeSize := n.spec.VolumeSize; volumeSize != nil && *volumeSize > 0 {
		var (
			kmsKeyID   *gfnt.Value
			volumeIOPS *gfnt.Value
		)
		if api.IsSetAndNonEmptyString(n.spec.VolumeKmsKeyID) {
			kmsKeyID = gfnt.NewString(*n.spec.VolumeKmsKeyID)
		}

		if *n.spec.VolumeType == api.NodeVolumeTypeIO1 {
			volumeIOPS = gfnt.NewInteger(*n.spec.VolumeIOPS)
		}

		launchTemplateData.BlockDeviceMappings = []gfnec2.LaunchTemplate_BlockDeviceMapping{{
			DeviceName: gfnt.NewString(*n.spec.VolumeName),
			Ebs: &gfnec2.LaunchTemplate_Ebs{
				VolumeSize: gfnt.NewInteger(*volumeSize),
				VolumeType: gfnt.NewString(*n.spec.VolumeType),
				Encrypted:  gfnt.NewBoolean(*n.spec.VolumeEncrypted),
				KmsKeyId:   kmsKeyID,
				Iops:       volumeIOPS,
			},
		}}
	}

	launchTemplate := &gfnec2.LaunchTemplate{
		LaunchTemplateName: launchTemplateName,
		LaunchTemplateData: launchTemplateData,
	}

	// Do not create a Launch Template resource for Spot-managed nodegroups.
	if n.spec.SpotOcean == nil {
		n.newResource("NodeGroupLaunchTemplate", launchTemplate)
	}

	vpcZoneIdentifier, err := AssignSubnets(n.spec.AvailabilityZones, n.clusterStackName, n.clusterSpec, n.spec.PrivateNetworking)
	if err != nil {
		return err
	}

	tags := []map[string]interface{}{
		{
			"Key":               "Name",
			"Value":             generateNodeName(n.spec.NodeGroupBase, n.clusterSpec.Metadata),
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
		return fmt.Errorf("failed to build nodegroup resource: %v", err)
	}
	n.newResource("NodeGroup", g)

	return nil
}

// generateNodeName formulates the name based on the configuration in input
func generateNodeName(ng *api.NodeGroupBase, meta *api.ClusterMeta) string {
	var nameParts []string
	if ng.InstancePrefix != "" {
		nameParts = append(nameParts, ng.InstancePrefix, "-")
	}
	// this overrides the default naming convention
	if ng.InstanceName != "" {
		nameParts = append(nameParts, ng.InstanceName)
	} else {
		nameParts = append(nameParts, fmt.Sprintf("%s-%s-Node", meta.Name, ng.Name))
	}
	return strings.Join(nameParts, "")
}

// AssignSubnets subnets based on the specified availability zones
func AssignSubnets(availabilityZones []string, clusterStackName string, clusterSpec *api.ClusterConfig, privateNetworking bool) (*gfnt.Value, error) {
	// currently goformation type system doesn't allow specifying `VPCZoneIdentifier: { "Fn::ImportValue": ... }`,
	// and tags don't have `PropagateAtLaunch` field, so we have a custom method here until this gets resolved

	if numNodeGroupsAZs := len(availabilityZones); numNodeGroupsAZs > 0 {
		subnets := clusterSpec.VPC.Subnets.Private
		if !privateNetworking {
			subnets = clusterSpec.VPC.Subnets.Public
		}
		makeErrorDesc := func() string {
			return fmt.Sprintf("(subnets=%#v AZs=%#v)", subnets, availabilityZones)
		}
		if len(subnets) < numNodeGroupsAZs {
			return nil, fmt.Errorf("VPC doesn't have enough subnets for nodegroup AZs %s", makeErrorDesc())
		}
		subnetIDs := make([]*gfnt.Value, numNodeGroupsAZs)
		for i, az := range availabilityZones {
			subnet, ok := subnets[az]
			if !ok {
				return nil, fmt.Errorf("VPC doesn't have subnets in %s %s", az, makeErrorDesc())
			}

			subnetIDs[i] = gfnt.NewString(subnet.ID)
		}
		return gfnt.NewSlice(subnetIDs...), nil
	}

	var subnets *gfnt.Value
	if privateNetworking {
		subnets = makeImportValue(clusterStackName, outputs.ClusterSubnetsPrivate)
	} else {
		subnets = makeImportValue(clusterStackName, outputs.ClusterSubnetsPublic)
	}

	return gfnt.MakeFnSplit(",", subnets), nil
}

// GetAllOutputs collects all outputs of the nodegroup
func (n *NodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack)
}

func (n *NodeGroupResourceSet) newLaunchTemplateData() *gfnec2.LaunchTemplate_LaunchTemplateData {
	launchTemplateData := &gfnec2.LaunchTemplate_LaunchTemplateData{
		IamInstanceProfile: &gfnec2.LaunchTemplate_IamInstanceProfile{
			Arn: n.instanceProfileARN,
		},
		ImageId:  gfnt.NewString(n.spec.AMI),
		UserData: n.userData,
		NetworkInterfaces: []gfnec2.LaunchTemplate_NetworkInterface{{
			// Explicitly un-setting this so that it doesn't get defaulted to true
			AssociatePublicIpAddress: nil,
			DeviceIndex:              gfnt.NewInteger(0),
			Groups:                   gfnt.NewSlice(n.securityGroups...),
		}},
		MetadataOptions: makeMetadataOptions(n.spec.NodeGroupBase),
	}

	if !api.HasMixedInstances(n.spec) {
		launchTemplateData.InstanceType = gfnt.NewString(n.spec.InstanceType)
	} else {
		launchTemplateData.InstanceType = gfnt.NewString(n.spec.InstancesDistribution.InstanceTypes[0])
	}
	if n.spec.EBSOptimized != nil {
		launchTemplateData.EbsOptimized = gfnt.NewBoolean(*n.spec.EBSOptimized)
	}

	if n.spec.CPUCredits != nil {
		launchTemplateData.CreditSpecification = &gfnec2.LaunchTemplate_CreditSpecification{
			CpuCredits: gfnt.NewString(strings.ToLower(*n.spec.CPUCredits)),
		}
	}

	return launchTemplateData
}

func makeMetadataOptions(ng *api.NodeGroupBase) *gfnec2.LaunchTemplate_MetadataOptions {
	imdsv2TokensRequired := "optional"
	if api.IsEnabled(ng.DisableIMDSv1) {
		imdsv2TokensRequired = "required"
	}
	return &gfnec2.LaunchTemplate_MetadataOptions{
		HttpPutResponseHopLimit: gfnt.NewInteger(2),
		HttpTokens:              gfnt.NewString(imdsv2TokensRequired),
	}
}

func (n *NodeGroupResourceSet) newNodeGroupResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, tags []map[string]interface{}) (*awsCloudFormationResource, error) {

	if n.spec.SpotOcean != nil {
		return n.newNodeGroupSpotOceanResource(launchTemplate, vpcZoneIdentifier, tags)
	}

	return n.newNodeGroupAutoScalingGroupResource(launchTemplate, vpcZoneIdentifier, tags)
}

func (n *NodeGroupResourceSet) newNodeGroupAutoScalingGroupResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, tags []map[string]interface{}) (*awsCloudFormationResource, error) {

	ng := n.spec
	ngProps := map[string]interface{}{
		"VPCZoneIdentifier": vpcZoneIdentifier,
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
	if len(ng.ASGMetricsCollection) > 0 {
		ngProps["MetricsCollection"] = metricsCollectionResource(ng.ASGMetricsCollection)
	}
	if len(ng.ClassicLoadBalancerNames) > 0 {
		ngProps["LoadBalancerNames"] = ng.ClassicLoadBalancerNames
	}
	if len(ng.TargetGroupARNs) > 0 {
		ngProps["TargetGroupARNs"] = ng.TargetGroupARNs
	}
	if api.HasMixedInstances(ng) {
		ngProps["MixedInstancesPolicy"] = n.newMixedInstancesPolicy(launchTemplate.LaunchTemplateName)
	} else {
		ngProps["LaunchTemplate"] = map[string]interface{}{
			"LaunchTemplateName": launchTemplate.LaunchTemplateName,
			"Version":            gfnt.MakeFnGetAttString("NodeGroupLaunchTemplate", "LatestVersionNumber"),
		}
	}

	rollingUpdate := map[string]interface{}{
		"MinInstancesInService": "0",
		"MaxBatchSize":          "1",
	}
	if len(ng.ASGSuspendProcesses) > 0 {
		rollingUpdate["SuspendProcesses"] = ng.ASGSuspendProcesses
	}

	return &awsCloudFormationResource{
		Type:       "AWS::AutoScaling::AutoScalingGroup",
		Properties: ngProps,
		UpdatePolicy: map[string]map[string]interface{}{
			"AutoScalingRollingUpdate": rollingUpdate,
		},
	}, nil
}

func (n *NodeGroupResourceSet) newMixedInstancesPolicy(launchTemplateName *gfnt.Value) map[string]interface{} {
	ng := n.spec
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
				"Version":            gfnt.MakeFnGetAttString("NodeGroupLaunchTemplate", "LatestVersionNumber"),
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

	if ng.InstancesDistribution.SpotAllocationStrategy != nil {
		instancesDistribution["SpotAllocationStrategy"] = *ng.InstancesDistribution.SpotAllocationStrategy
	}

	policy["InstancesDistribution"] = instancesDistribution

	return policy
}

func metricsCollectionResource(asgMetricsCollection []api.MetricsCollection) []map[string]interface{} {
	var metricsCollections []map[string]interface{}
	for _, m := range asgMetricsCollection {
		newCollection := make(map[string]interface{})

		if len(m.Metrics) > 0 {
			newCollection["Metrics"] = m.Metrics
		}
		newCollection["Granularity"] = m.Granularity

		metricsCollections = append(metricsCollections, newCollection)
	}
	return metricsCollections
}

// newNodeGroupSpotOceanResource returns a Spot Ocean resource.
func (n *NodeGroupResourceSet) newNodeGroupSpotOceanResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, tags []map[string]interface{}) (*awsCloudFormationResource, error) {

	var res *spot.NodeGroupResource
	var out awsCloudFormationResource
	var err error

	// Resource.
	{
		if n.spec.Name == api.SpotOceanNodeGroupName {
			logger.Debug("ocean: creating cluster for nodegroup %q", n.spec.Name)
			res, err = n.newNodeGroupSpotOceanClusterResource(
				launchTemplate, vpcZoneIdentifier, tags)
		} else {
			logger.Debug("ocean: creating launchspec for nodegroup %q", n.spec.Name)
			res, err = n.newNodeGroupSpotOceanLaunchSpecResource(
				launchTemplate, vpcZoneIdentifier, tags)
		}
		if err != nil {
			return nil, err
		}
	}

	// Credentials.
	{
		var profile *string
		if n.spec.SpotOcean.Metadata != nil {
			profile = n.spec.SpotOcean.Metadata.Profile
		}

		token, account, err := spot.LoadCredentials(profile)
		if err != nil {
			return nil, err
		}

		if token != "" {
			res.Token = n.rs.newParameter(spot.CredentialsTokenParameterKey, &gfnv4.Parameter{
				Type:    "String",
				Default: token,
			})
		}

		if account != "" {
			res.Account = n.rs.newParameter(spot.CredentialsAccountParameterKey, &gfnv4.Parameter{
				Type:    "String",
				Default: account,
			})
		}
	}

	// Feature Flags.
	{
		ff := spot.LoadFeatureFlags()
		if ff != "" {
			res.FeatureFlags = n.rs.newParameter(spot.FeatureFlagsParameterKey, &gfnv4.Parameter{
				Type:    "String",
				Default: ff,
			})
		}
	}

	// Service Token.
	{
		svc, err := spot.LoadServiceToken()
		if err != nil {
			return nil, err
		}

		if svc != "" {
			res.ServiceToken = gfnt.MakeFnSubString(svc)
		}
	}

	// Convert.
	{
		b, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(b, &out); err != nil {
			return nil, err
		}
	}

	return &out, nil
}

// newNodeGroupSpotOceanClusterResource returns a Spot Ocean cluster resource.
func (n *NodeGroupResourceSet) newNodeGroupSpotOceanClusterResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, tags []map[string]interface{}) (*spot.NodeGroupResource, error) {

	template := launchTemplate.LaunchTemplateData
	cluster := &spot.NodeGroupCluster{
		Name:      spotinst.String(n.clusterSpec.Metadata.Name),
		ClusterID: spotinst.String(n.clusterSpec.Metadata.Name),
		Region:    gfnt.MakeRef("AWS::Region"),
		Compute: &spot.NodeGroupCompute{
			LaunchSpecification: &spot.NodeGroupLaunchSpec{
				ImageID:      template.ImageId,
				UserData:     template.UserData,
				KeyPair:      template.KeyName,
				EBSOptimized: n.spec.EBSOptimized,
			},
			SubnetIDs: vpcZoneIdentifier,
		},
		Capacity: &spot.NodeGroupCapacity{
			Target:  n.spec.DesiredCapacity,
			Minimum: n.spec.MinSize,
			Maximum: n.spec.MaxSize,
		},
	}

	// Storage.
	{
		if n.spec.VolumeSize != nil && spotinst.IntValue(n.spec.VolumeSize) > 0 {
			cluster.Compute.LaunchSpecification.VolumeSize = n.spec.VolumeSize
		}
	}

	// Strategy.
	{
		if strategy := n.spec.SpotOcean.Strategy; strategy != nil {
			cluster.Strategy = &spot.NodeGroupStrategy{
				SpotPercentage:           strategy.SpotPercentage,
				UtilizeReservedInstances: strategy.UtilizeReservedInstances,
				FallbackToOnDemand:       strategy.FallbackToOnDemand,
			}
		}
	}

	// IAM.
	{
		if template.IamInstanceProfile != nil {
			cluster.Compute.LaunchSpecification.IAMInstanceProfile = map[string]*gfnt.Value{
				"arn": template.IamInstanceProfile.Arn,
			}
		}
	}

	// Networking.
	{
		if len(template.NetworkInterfaces) > 0 {
			if template.NetworkInterfaces[0].AssociatePublicIpAddress != nil {
				cluster.Compute.LaunchSpecification.AssociatePublicIPAddress = template.NetworkInterfaces[0].AssociatePublicIpAddress
			}

			if template.NetworkInterfaces[0].Groups != nil {
				cluster.Compute.LaunchSpecification.SecurityGroupIDs = template.NetworkInterfaces[0].Groups
			}
		}
	}

	// Load Balancers.
	{
		var lbs []*spot.NodeGroupLoadBalancer

		// ELBs.
		if len(n.spec.ClassicLoadBalancerNames) > 0 {
			for _, name := range n.spec.ClassicLoadBalancerNames {
				lbs = append(lbs, &spot.NodeGroupLoadBalancer{
					Type: spotinst.String("CLASSIC"),
					Name: spotinst.String(name),
				})
			}
		}

		// ALBs.
		if len(n.spec.TargetGroupARNs) > 0 {
			for _, arn := range n.spec.TargetGroupARNs {
				lbs = append(lbs, &spot.NodeGroupLoadBalancer{
					Type: spotinst.String("TARGET_GROUP"),
					Arn:  spotinst.String(arn),
				})
			}
		}

		if len(lbs) > 0 {
			cluster.Compute.LaunchSpecification.LoadBalancers = lbs
		}
	}

	// Tags.
	{
		var tagsKV []*spot.NodeGroupTag

		// Nodegroup tags.
		if len(n.spec.Tags) > 0 {
			for key, value := range n.spec.Tags {
				tagsKV = append(tagsKV, &spot.NodeGroupTag{
					Key:   spotinst.String(key),
					Value: spotinst.String(value),
				})
			}
		}

		// Resource tags (Name, kubernetes.io/*, k8s.io/*, etc.).
		if len(tags) > 0 {
			for _, tag := range tags {
				tagsKV = append(tagsKV, &spot.NodeGroupTag{
					Key:   tag["Key"],
					Value: tag["Value"],
				})
			}
		}

		// Shared tags (metadata.tags + eksctl's tags).
		if len(n.sharedTags) > 0 {
			for _, tag := range n.sharedTags {
				tagsKV = append(tagsKV, &spot.NodeGroupTag{
					Key:   spotinst.StringValue(tag.Key),
					Value: spotinst.StringValue(tag.Value),
				})
			}
		}

		if len(tagsKV) > 0 {
			cluster.Compute.LaunchSpecification.Tags = tagsKV
		}
	}

	// Instance Types.
	{
		if compute := n.spec.SpotOcean.Compute; compute != nil && compute.InstanceTypes != nil {
			cluster.Compute.InstanceTypes = &spot.NodeGroupInstanceTypes{
				Whitelist: compute.InstanceTypes.Whitelist,
				Blacklist: compute.InstanceTypes.Blacklist,
			}
		}
	}

	// Scheduling.
	{
		if scheduling := n.spec.SpotOcean.Scheduling; scheduling != nil {
			if hours := scheduling.ShutdownHours; hours != nil {
				cluster.Scheduling = &spot.NodeGroupScheduling{
					ShutdownHours: &spot.NodeGroupSchedulingShutdownHours{
						IsEnabled:   hours.IsEnabled,
						TimeWindows: hours.TimeWindows,
					},
				}
			}

			if tasks := scheduling.Tasks; len(tasks) > 0 {
				if cluster.Scheduling == nil {
					cluster.Scheduling = new(spot.NodeGroupScheduling)
				}

				cluster.Scheduling.Tasks = make([]*spot.NodeGroupSchedulingTask, len(tasks))
				for i, task := range tasks {
					cluster.Scheduling.Tasks[i] = &spot.NodeGroupSchedulingTask{
						IsEnabled:      task.IsEnabled,
						Type:           task.Type,
						CronExpression: task.CronExpression,
					}
				}
			}
		}
	}

	// Auto Scaler.
	{
		if autoScaler := n.spec.SpotOcean.AutoScaler; autoScaler != nil {
			cluster.AutoScaler = &spot.NodeGroupAutoScaler{
				IsEnabled:    autoScaler.Enabled,
				IsAutoConfig: autoScaler.AutoConfig,
				Cooldown:     autoScaler.Cooldown,
			}

			if headrooms := autoScaler.Headrooms; len(headrooms) > 0 {
				cluster.AutoScaler.Headroom = &spot.NodeGroupAutoScalerHeadroom{
					CPUPerUnit:    headrooms[0].CPUPerUnit,
					GPUPerUnit:    headrooms[0].GPUPerUnit,
					MemoryPerUnit: headrooms[0].MemoryPerUnit,
					NumOfUnits:    headrooms[0].NumOfUnits,
				}
			}
		}
	}

	// Outputs.
	{
		n.rs.defineOutputWithoutCollector(
			outputs.NodeGroupSpotOceanClusterID,
			gfnt.MakeRef("NodeGroup"),
			true)
	}

	return &spot.NodeGroupResource{
		OceanCluster: cluster,
		OceanSummary: &spot.NodeGroupSummary{
			ImageID:  cluster.Compute.LaunchSpecification.ImageID,
			Capacity: cluster.Capacity,
		},
	}, nil
}

// newNodeGroupSpotOceanLaunchSpecResource returns a Spot Ocean launchspec resource.
func (n *NodeGroupResourceSet) newNodeGroupSpotOceanLaunchSpecResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, tags []map[string]interface{}) (*spot.NodeGroupResource, error) {

	// Import the Ocean cluster identifier.
	oceanClusterStackName := fmt.Sprintf("eksctl-%s-nodegroup-ocean", n.clusterSpec.Metadata.Name)
	oceanClusterID := makeImportValue(oceanClusterStackName, outputs.NodeGroupSpotOceanClusterID)

	template := launchTemplate.LaunchTemplateData
	spec := &spot.NodeGroupLaunchSpec{
		Name:      spotinst.String(n.spec.Name),
		OceanID:   oceanClusterID,
		ImageID:   template.ImageId,
		UserData:  template.UserData,
		SubnetIDs: vpcZoneIdentifier,
	}

	// Storage.
	{
		if n.spec.VolumeSize != nil && spotinst.IntValue(n.spec.VolumeSize) > 0 {
			spec.VolumeSize = n.spec.VolumeSize
		}
	}

	// IAM.
	{
		if template.IamInstanceProfile != nil {
			spec.IAMInstanceProfile = map[string]*gfnt.Value{
				"arn": template.IamInstanceProfile.Arn,
			}
		}
	}

	// Networking.
	{
		if len(template.NetworkInterfaces) > 0 {
			if template.NetworkInterfaces[0].Groups != nil {
				spec.SecurityGroupIDs = template.NetworkInterfaces[0].Groups
			}
		}
	}

	// Tags.
	{
		var tagsKV []*spot.NodeGroupTag

		// Nodegroup tags.
		if len(n.spec.Tags) > 0 {
			for key, value := range n.spec.Tags {
				tagsKV = append(tagsKV, &spot.NodeGroupTag{
					Key:   spotinst.String(key),
					Value: spotinst.String(value),
				})
			}
		}

		// Resource tags (Name, kubernetes.io/*, k8s.io/*, etc.).
		if len(tags) > 0 {
			for _, tag := range tags {
				tagsKV = append(tagsKV, &spot.NodeGroupTag{
					Key:   tag["Key"],
					Value: tag["Value"],
				})
			}
		}

		// Shared tags (metadata.tags + eksctl's tags).
		if len(n.sharedTags) > 0 {
			for _, tag := range n.sharedTags {
				tagsKV = append(tagsKV, &spot.NodeGroupTag{
					Key:   spotinst.StringValue(tag.Key),
					Value: spotinst.StringValue(tag.Value),
				})
			}
		}

		if len(tagsKV) > 0 {
			spec.Tags = tagsKV
		}
	}

	// Instance Types.
	{
		if compute := n.spec.SpotOcean.Compute; compute != nil && compute.InstanceTypes != nil {
			spec.InstanceTypes = compute.InstanceTypes.Whitelist
		}
	}

	// Labels.
	{
		if len(n.spec.Labels) > 0 {
			labels := make([]*spot.NodeGroupLabel, 0, len(n.spec.Labels))

			for key, value := range n.spec.Labels {
				labels = append(labels, &spot.NodeGroupLabel{
					Key:   spotinst.String(key),
					Value: spotinst.String(value),
				})
			}

			spec.Labels = labels
		}
	}

	// Taints.
	{
		if len(n.spec.Taints) > 0 {
			taints := make([]*spot.NodeGroupTaint, 0, len(n.spec.Taints))

			for key, valueEffect := range n.spec.Taints {
				taint := &spot.NodeGroupTaint{
					Key: spotinst.String(key),
				}
				parts := strings.Split(valueEffect, ":")
				if len(parts) >= 1 {
					taint.Value = spotinst.String(parts[0])
				}
				if len(parts) > 1 {
					taint.Effect = spotinst.String(parts[1])
				}
				taints = append(taints, taint)
			}

			spec.Taints = taints
		}
	}

	// Auto Scaler.
	{
		if autoScaler := n.spec.SpotOcean.AutoScaler; autoScaler != nil && len(autoScaler.Headrooms) > 0 {
			headrooms := make([]*spot.NodeGroupAutoScalerHeadroom, len(autoScaler.Headrooms))

			for i, headroom := range autoScaler.Headrooms {
				headrooms[i] = &spot.NodeGroupAutoScalerHeadroom{
					CPUPerUnit:    headroom.CPUPerUnit,
					GPUPerUnit:    headroom.GPUPerUnit,
					MemoryPerUnit: headroom.MemoryPerUnit,
					NumOfUnits:    headroom.NumOfUnits,
				}
			}

			spec.AutoScaler = &spot.NodeGroupAutoScaler{
				Headrooms: headrooms,
			}
		}
	}

	// Outputs.
	{
		n.rs.defineOutputWithoutCollector(
			outputs.NodeGroupSpotOceanLaunchSpecID,
			gfnt.MakeRef("NodeGroup"),
			true)
	}

	return &spot.NodeGroupResource{
		OceanLaunchSpec: spec,
		OceanSummary: &spot.NodeGroupSummary{
			ImageID: spec.ImageID,
			Capacity: &spot.NodeGroupCapacity{
				Target:  n.spec.DesiredCapacity,
				Minimum: n.spec.MinSize,
				Maximum: n.spec.MaxSize,
			},
		},
	}, nil
}
