package builder

import (
	"fmt"
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// NodeGroupResourceSet stores the resource information of the nodegroup
type NodeGroupResourceSet struct {
	rs                   *resourceSet
	clusterSpec          *api.ClusterConfig
	spec                 *api.NodeGroup
	supportsManagedNodes bool
	forceAddCNIPolicy    bool
	provider             api.ClusterProvider
	clusterStackName     string
	instanceProfileARN   *gfnt.Value
	securityGroups       []*gfnt.Value
	vpc                  *gfnt.Value
	userData             *gfnt.Value
}

// NewNodeGroupResourceSet returns a resource set for a nodegroup embedded in a cluster config
func NewNodeGroupResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig, clusterStackName string, ng *api.NodeGroup,
	supportsManagedNodes, forceAddCNIPolicy bool) *NodeGroupResourceSet {
	return &NodeGroupResourceSet{
		rs:                   newResourceSet(),
		clusterStackName:     clusterStackName,
		supportsManagedNodes: supportsManagedNodes,
		forceAddCNIPolicy:    forceAddCNIPolicy,
		clusterSpec:          spec,
		spec:                 ng,
		provider:             provider,
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
	launchTemplateData := newLaunchTemplateData(n)

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

		if n.spec.OSVolumeName != nil && *n.spec.OSVolumeName != "" {
			launchTemplateData.BlockDeviceMappings = append(launchTemplateData.BlockDeviceMappings, gfnec2.LaunchTemplate_BlockDeviceMapping{
				DeviceName: gfnt.NewString(*n.spec.OSVolumeName),
				Ebs: &gfnec2.LaunchTemplate_Ebs{
					Encrypted: gfnt.NewBoolean(*n.spec.VolumeEncrypted),
					KmsKeyId:  kmsKeyID,
				},
			})
		}
	}

	n.newResource("NodeGroupLaunchTemplate", &gfnec2.LaunchTemplate{
		LaunchTemplateName: launchTemplateName,
		LaunchTemplateData: launchTemplateData,
	})

	vpcZoneIdentifier, err := AssignSubnets(n.spec.NodeGroupBase, n.clusterStackName, n.clusterSpec)
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

	asg := nodeGroupResource(launchTemplateName, vpcZoneIdentifier, tags, n.spec)
	n.newResource("NodeGroup", asg)

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
func AssignSubnets(spec *api.NodeGroupBase, clusterStackName string, clusterSpec *api.ClusterConfig) (*gfnt.Value, error) {
	// currently goformation type system doesn't allow specifying `VPCZoneIdentifier: { "Fn::ImportValue": ... }`,
	// and tags don't have `PropagateAtLaunch` field, so we have a custom method here until this gets resolved

	if numNodeGroupsAZs, numNodeGroupsSubnets := len(spec.AvailabilityZones), len(spec.Subnets); numNodeGroupsAZs > 0 || numNodeGroupsSubnets > 0 {
		subnets := clusterSpec.VPC.Subnets.Public
		if spec.PrivateNetworking {
			subnets = clusterSpec.VPC.Subnets.Private
		}
		subnetIDs, err := vpc.SelectNodeGroupSubnets(spec.AvailabilityZones, spec.Subnets, subnets)
		return gfnt.NewStringSlice(subnetIDs...), err
	}

	var subnets *gfnt.Value
	if spec.PrivateNetworking {
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

func newLaunchTemplateData(n *NodeGroupResourceSet) *gfnec2.LaunchTemplate_LaunchTemplateData {

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

	if n.spec.Placement != nil {
		launchTemplateData.Placement = &gfnec2.LaunchTemplate_Placement{
			GroupName: gfnt.NewString(n.spec.Placement.GroupName),
		}
	}

	return launchTemplateData
}

func makeMetadataOptions(ng *api.NodeGroupBase) *gfnec2.LaunchTemplate_MetadataOptions {
	imdsv2TokensRequired := "optional"
	if api.IsEnabled(ng.DisableIMDSv1) || api.IsEnabled(ng.DisablePodIMDS) {
		imdsv2TokensRequired = "required"
	}
	hopLimit := 2
	if api.IsEnabled(ng.DisablePodIMDS) {
		hopLimit = 1
	}
	return &gfnec2.LaunchTemplate_MetadataOptions{
		HttpPutResponseHopLimit: gfnt.NewInteger(hopLimit),
		HttpTokens:              gfnt.NewString(imdsv2TokensRequired),
	}
}

func nodeGroupResource(launchTemplateName *gfnt.Value, vpcZoneIdentifier interface{}, tags []map[string]interface{}, ng *api.NodeGroup) *awsCloudFormationResource {
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
		ngProps["MixedInstancesPolicy"] = *mixedInstancesPolicy(launchTemplateName, ng)
	} else {
		ngProps["LaunchTemplate"] = map[string]interface{}{
			"LaunchTemplateName": launchTemplateName,
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
	}
}

func mixedInstancesPolicy(launchTemplateName *gfnt.Value, ng *api.NodeGroup) *map[string]interface{} {
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

	return &policy
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
