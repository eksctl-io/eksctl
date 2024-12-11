package builder

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"k8s.io/utils/strings/slices"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	gfn "goformation/v4/cloudformation"
	gfncfn "goformation/v4/cloudformation/cloudformation"
	gfnec2 "goformation/v4/cloudformation/ec2"
	gfneks "goformation/v4/cloudformation/eks"
	gfnt "goformation/v4/cloudformation/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/az"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// MaximumTagNumber for ASGs as described here https://docs.aws.amazon.com/autoscaling/ec2/userguide/autoscaling-tagging.html
const MaximumTagNumber = 50
const MaximumCreatedTagNumberPerCall = 25

const (
	nodeTemplatePrefix = "k8s.io/cluster-autoscaler/node-template/"
	labelsPrefix       = nodeTemplatePrefix + "label/"
	taintsPrefix       = nodeTemplatePrefix + "taint/"
)

// NodeGroupOptions represents options passed to a NodeGroupResourceSet.
type NodeGroupOptions struct {
	ClusterConfig      *api.ClusterConfig
	NodeGroup          *api.NodeGroup
	Bootstrapper       nodebootstrap.Bootstrapper
	ForceAddCNIPolicy  bool
	VPCImporter        vpc.Importer
	SkipEgressRules    bool
	DisableAccessEntry bool
	// DisableAccessEntryResource disables creation of an access entry resource but still attaches the UsesAccessEntry tag.
	DisableAccessEntryResource bool
}

// NodeGroupResourceSet stores the resource information of the nodegroup
type NodeGroupResourceSet struct {
	rs      *resourceSet
	iamAPI  awsapi.IAM
	ec2API  awsapi.EC2
	options NodeGroupOptions

	instanceProfileARN *gfnt.Value
	securityGroups     []*gfnt.Value
	vpc                *gfnt.Value
}

// NewNodeGroupResourceSet returns a resource set for a nodegroup embedded in a cluster config.
func NewNodeGroupResourceSet(ec2API awsapi.EC2, iamAPI awsapi.IAM, options NodeGroupOptions) *NodeGroupResourceSet {
	return &NodeGroupResourceSet{
		rs:      newResourceSet(),
		ec2API:  ec2API,
		iamAPI:  iamAPI,
		options: options,
	}
}

// AddAllResources adds all the information about the nodegroup to the resource set
func (n *NodeGroupResourceSet) AddAllResources(ctx context.Context) error {

	if n.options.ClusterConfig.IPv6Enabled() {
		return errors.New("unmanaged nodegroups are not supported with IPv6 clusters")
	}

	ng := n.options.NodeGroup
	n.rs.template.Description = fmt.Sprintf(
		"%s (AMI family: %s, SSH access: %v, private networking: %v) %s",
		nodeGroupTemplateDescription,
		ng.AMIFamily, api.IsEnabled(ng.SSH.Allow), ng.PrivateNetworking,
		templateDescriptionSuffix)

	n.Template().Mappings[servicePrincipalPartitionMapName] = api.Partitions.ServicePrincipalPartitionMappings()

	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeaturePrivateNetworking, ng.PrivateNetworking, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureSharedSecurityGroup, ng.SecurityGroups.WithShared, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureLocalSecurityGroup, ng.SecurityGroups.WithLocal, false)

	n.vpc = n.options.VPCImporter.VPC()

	if ng.Tags == nil {
		ng.Tags = map[string]string{}
	}

	for k, v := range n.options.ClusterConfig.Metadata.Tags {
		if _, exists := ng.Tags[k]; !exists {
			ng.Tags[k] = v
		}
	}

	// Ensure MinSize is set, as it is required by the ASG cfn resource
	// TODO this validation and default setting should happen way earlier than this
	if ng.MinSize == nil {
		if ng.DesiredCapacity == nil {
			defaultNodeCount := api.DefaultNodeCount
			ng.MinSize = &defaultNodeCount
		} else {
			ng.MinSize = ng.DesiredCapacity
		}
		logger.Info("--nodes-min=%d was set automatically for nodegroup %s", *ng.MinSize, ng.Name)
	} else if ng.DesiredCapacity != nil && *ng.DesiredCapacity < *ng.MinSize {
		return fmt.Errorf("--nodes value (%d) cannot be lower than --nodes-min value (%d)", *ng.DesiredCapacity, *ng.MinSize)
	}

	// Ensure MaxSize is set, as it is required by the ASG cfn resource
	if ng.MaxSize == nil {
		if ng.DesiredCapacity == nil {
			ng.MaxSize = ng.MinSize
		} else {
			ng.MaxSize = ng.DesiredCapacity
		}
		logger.Info("--nodes-max=%d was set automatically for nodegroup %s", *ng.MaxSize, ng.Name)
	} else if ng.DesiredCapacity != nil && *ng.DesiredCapacity > *ng.MaxSize {
		return fmt.Errorf("--nodes value (%d) cannot be greater than --nodes-max value (%d)", *ng.DesiredCapacity, *ng.MaxSize)
	} else if *ng.MaxSize < *ng.MinSize {
		return fmt.Errorf("--nodes-min value (%d) cannot be greater than --nodes-max value (%d)", *ng.MinSize, *ng.MaxSize)
	}

	if err := n.addResourcesForIAM(ctx); err != nil {
		return err
	}
	n.addResourcesForSecurityGroups()
	if !n.options.DisableAccessEntry {
		n.addAccessEntry()
	}

	return n.addResourcesForNodeGroup(ctx)
}

// A PartialEgressRule represents a partial security group egress rule.
type PartialEgressRule struct {
	FromPort   int
	ToPort     int
	IPProtocol string
}

var controlPlaneEgressInterCluster = PartialEgressRule{
	FromPort:   1025,
	ToPort:     65535,
	IPProtocol: "tcp",
}

var controlPlaneEgressInterClusterAPI = PartialEgressRule{
	FromPort:   443,
	ToPort:     443,
	IPProtocol: "tcp",
}

// ControlPlaneNodeGroupEgressRules is a slice of egress rules attached to the control plane security group.
var ControlPlaneNodeGroupEgressRules = []PartialEgressRule{
	controlPlaneEgressInterCluster,
	controlPlaneEgressInterClusterAPI,
}

// ControlPlaneEgressRuleDescriptionPrefix is the prefix applied to the description for control plane security group egress rules.
var ControlPlaneEgressRuleDescriptionPrefix = "Allow control plane to communicate with "

func (n *NodeGroupResourceSet) addAccessEntry() {
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupUsesAccessEntry, true, false)
	if n.options.DisableAccessEntryResource {
		return
	}

	n.newResource("AccessEntry", &gfneks.AccessEntry{
		PrincipalArn: gfnt.MakeFnGetAttString(cfnIAMInstanceRoleName, "Arn"),
		ClusterName:  gfnt.NewString(n.options.ClusterConfig.Metadata.Name),
		Type:         gfnt.NewString(string(api.GetAccessEntryType(n.options.NodeGroup))),
	})
}

func (n *NodeGroupResourceSet) addResourcesForSecurityGroups() {
	ng := n.options.NodeGroup
	for _, id := range ng.SecurityGroups.AttachIDs {
		n.securityGroups = append(n.securityGroups, gfnt.NewString(id))
	}

	if api.IsEnabled(ng.SecurityGroups.WithShared) {
		n.securityGroups = append(n.securityGroups, n.options.VPCImporter.SharedNodeSecurityGroup())
	}

	if api.IsDisabled(ng.SecurityGroups.WithLocal) {
		return
	}

	desc := "worker nodes in group " + ng.Name
	vpcID := n.options.VPCImporter.VPC()
	refControlPlaneSG := n.options.VPCImporter.ControlPlaneSecurityGroup()

	refNodeGroupLocalSG := n.newResource("SG", &gfnec2.SecurityGroup{
		VpcId:            vpcID,
		GroupDescription: gfnt.NewString("Communication between the control plane and " + desc),
		Tags: []gfncfn.Tag{{
			Key:   gfnt.NewString("kubernetes.io/cluster/" + n.options.ClusterConfig.Metadata.Name),
			Value: gfnt.NewString("owned"),
		}},
		SecurityGroupIngress: makeNodeIngressRules(ng.NodeGroupBase, refControlPlaneSG, n.options.ClusterConfig.VPC.CIDR.String(), desc),
	})

	n.securityGroups = append(n.securityGroups, refNodeGroupLocalSG)

	if api.IsEnabled(ng.EFAEnabled) {
		efaSG := n.rs.addEFASecurityGroup(vpcID, n.options.ClusterConfig.Metadata.Name, desc)
		n.securityGroups = append(n.securityGroups, efaSG)
	}

	if !n.options.SkipEgressRules {
		n.newResource("EgressInterCluster", &gfnec2.SecurityGroupEgress{
			GroupId:                    refControlPlaneSG,
			DestinationSecurityGroupId: refNodeGroupLocalSG,
			Description:                gfnt.NewString(ControlPlaneEgressRuleDescriptionPrefix + desc + " (kubelet and workload TCP ports)"),
			IpProtocol:                 gfnt.NewString(controlPlaneEgressInterCluster.IPProtocol),
			FromPort:                   gfnt.NewInteger(controlPlaneEgressInterCluster.FromPort),
			ToPort:                     gfnt.NewInteger(controlPlaneEgressInterCluster.ToPort),
		})
		n.newResource("EgressInterClusterAPI", &gfnec2.SecurityGroupEgress{
			GroupId:                    refControlPlaneSG,
			DestinationSecurityGroupId: refNodeGroupLocalSG,
			Description:                gfnt.NewString(ControlPlaneEgressRuleDescriptionPrefix + desc + " (workloads using HTTPS port, commonly used with extension API servers)"),
			IpProtocol:                 gfnt.NewString(controlPlaneEgressInterClusterAPI.IPProtocol),
			FromPort:                   gfnt.NewInteger(controlPlaneEgressInterClusterAPI.FromPort),
			ToPort:                     gfnt.NewInteger(controlPlaneEgressInterClusterAPI.ToPort),
		})
	}
	n.newResource("IngressInterClusterCP", &gfnec2.SecurityGroupIngress{
		GroupId:               refControlPlaneSG,
		SourceSecurityGroupId: refNodeGroupLocalSG,
		Description:           gfnt.NewString("Allow control plane to receive API requests from " + desc),
		IpProtocol:            sgProtoTCP,
		FromPort:              sgPortHTTPS,
		ToPort:                sgPortHTTPS,
	})
}

func makeNodeIngressRules(ng *api.NodeGroupBase, controlPlaneSG *gfnt.Value, vpcCIDR, description string) []gfnec2.SecurityGroup_Ingress {
	ingressRules := []gfnec2.SecurityGroup_Ingress{
		{
			SourceSecurityGroupId: controlPlaneSG,
			Description:           gfnt.NewString(fmt.Sprintf("[IngressInterCluster] Allow %s to communicate with control plane (kubelet and workload TCP ports)", description)),
			IpProtocol:            gfnt.NewString(controlPlaneEgressInterCluster.IPProtocol),
			FromPort:              gfnt.NewInteger(controlPlaneEgressInterCluster.FromPort),
			ToPort:                gfnt.NewInteger(controlPlaneEgressInterCluster.ToPort),
		},
		{
			SourceSecurityGroupId: controlPlaneSG,
			Description:           gfnt.NewString(fmt.Sprintf("[IngressInterClusterAPI] Allow %s to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)", description)),
			IpProtocol:            gfnt.NewString(controlPlaneEgressInterClusterAPI.IPProtocol),
			FromPort:              gfnt.NewInteger(controlPlaneEgressInterClusterAPI.FromPort),
			ToPort:                gfnt.NewInteger(controlPlaneEgressInterClusterAPI.ToPort),
		},
	}

	return append(ingressRules, makeSSHIngressRules(ng, vpcCIDR, description)...)
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

func (n *NodeGroupResourceSet) addResourcesForNodeGroup(ctx context.Context) error {
	launchTemplateName := gfnt.MakeFnSubString(fmt.Sprintf("${%s}", gfnt.StackName))
	launchTemplateData, err := newLaunchTemplateData(ctx, n)
	if err != nil {
		return fmt.Errorf("could not add resources for nodegroup: %w", err)
	}

	ng := n.options.NodeGroup
	if ng.SSH != nil && api.IsSetAndNonEmptyString(ng.SSH.PublicKeyName) {
		launchTemplateData.KeyName = gfnt.NewString(*ng.SSH.PublicKeyName)
	}

	launchTemplateData.BlockDeviceMappings = makeBlockDeviceMappings(ng.NodeGroupBase)

	n.newResource("NodeGroupLaunchTemplate", &gfnec2.LaunchTemplate{
		LaunchTemplateName: launchTemplateName,
		LaunchTemplateData: launchTemplateData,
	})

	vpcZoneIdentifier, err := AssignSubnets(ctx, ng, n.options.ClusterConfig, n.ec2API)
	if err != nil {
		return err
	}

	tags := []map[string]string{
		{
			"Key":               "Name",
			"Value":             generateNodeName(ng.NodeGroupBase, n.options.ClusterConfig.Metadata),
			"PropagateAtLaunch": "true",
		},
		{
			"Key":               "kubernetes.io/cluster/" + n.options.ClusterConfig.Metadata.Name,
			"Value":             "owned",
			"PropagateAtLaunch": "true",
		},
	}
	if api.IsEnabled(ng.IAM.WithAddonPolicies.AutoScaler) {
		tags = append(tags,
			map[string]string{
				"Key":               "k8s.io/cluster-autoscaler/enabled",
				"Value":             "true",
				"PropagateAtLaunch": "true",
			},
			map[string]string{
				"Key":               "k8s.io/cluster-autoscaler/" + n.options.ClusterConfig.Metadata.Name,
				"Value":             "owned",
				"PropagateAtLaunch": "true",
			},
		)
	}

	if api.IsEnabled(ng.PropagateASGTags) {
		var clusterTags []map[string]string
		GenerateClusterAutoscalerTags(ng, func(key, value string) {
			clusterTags = append(clusterTags, map[string]string{
				"Key":               key,
				"Value":             value,
				"PropagateAtLaunch": "true",
			})
		})
		tags = append(tags, clusterTags...)
		if len(tags) > MaximumTagNumber {
			return fmt.Errorf("number of tags is exceeding the configured amount %d, was: %d. Due to desiredCapacity==0 we added an extra %d number of tags to ensure the nodegroup is scaled correctly", MaximumTagNumber, len(tags), len(clusterTags))
		}
	}

	asg := nodeGroupResource(launchTemplateName, vpcZoneIdentifier, tags, ng)
	n.newResource("NodeGroup", asg)

	return nil
}

// GenerateClusterAutoscalerTags generates Cluster Autoscaler tags for labels and taints.
func GenerateClusterAutoscalerTags(np api.NodePool, addTag func(key, value string)) {
	// labels
	for k, v := range np.BaseNodeGroup().Labels {
		addTag(labelsPrefix+k, v)
	}

	var taints []api.NodeGroupTaint
	switch ng := np.(type) {
	case *api.NodeGroup:
		taints = ng.Taints
	case *api.ManagedNodeGroup:
		taints = ng.Taints
	}

	// taints
	for _, taint := range taints {
		addTag(taintsPrefix+taint.Key, taint.Value)
	}
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

// AssignSubnets assigns subnets based on the availability zones, local zones and subnet IDs in the specified nodegroup.
func AssignSubnets(ctx context.Context, np api.NodePool, clusterConfig *api.ClusterConfig, ec2API awsapi.EC2) (*gfnt.Value, error) {
	ng := np.BaseNodeGroup()
	if !shouldImportSubnetsFromVPC(np, clusterConfig) {
		subnetIDs, err := vpc.SelectNodeGroupSubnets(ctx, np, clusterConfig, ec2API)
		if err != nil {
			return nil, err
		}
		return gfnt.NewStringSlice(subnetIDs...), nil
	}

	supportedZones, err := az.FilterBasedOnAvailability(ctx, clusterConfig.AvailabilityZones, []api.NodePool{np}, ec2API)
	if err != nil {
		return nil, err
	}

	subnetMapping := clusterConfig.VPC.Subnets.Public
	if ng.PrivateNetworking {
		subnetMapping = clusterConfig.VPC.Subnets.Private
	}

	subnetIDs := []string{}
	// only assign a subnet if the AZ to which it belongs supports all required instance types
	for key, subnetSpec := range subnetMapping {
		az := subnetSpec.AZ
		if az == "" {
			az = key
		}
		if !slices.Contains(supportedZones, az) {
			continue
		}
		subnetIDs = append(subnetIDs, subnetSpec.ID)
	}
	return gfnt.NewStringSlice(subnetIDs...), nil
}

func shouldImportSubnetsFromVPC(np api.NodePool, cfg *api.ClusterConfig) bool {
	ng := np.BaseNodeGroup()
	if nodeGroup, ok := np.(*api.NodeGroup); (!ok || len(nodeGroup.LocalZones) == 0) && len(ng.AvailabilityZones) == 0 && len(ng.Subnets) == 0 && (ng.OutpostARN == "" || cfg.IsControlPlaneOnOutposts()) {
		return true
	}
	return false
}

// GetAllOutputs collects all outputs of the nodegroup
func (n *NodeGroupResourceSet) GetAllOutputs(stack types.Stack) error {
	return n.rs.GetAllOutputs(stack)
}

func newLaunchTemplateData(ctx context.Context, n *NodeGroupResourceSet) (*gfnec2.LaunchTemplate_LaunchTemplateData, error) {
	userData, err := n.options.Bootstrapper.UserData()
	if err != nil {
		return nil, err
	}

	ng := n.options.NodeGroup
	launchTemplateData := &gfnec2.LaunchTemplate_LaunchTemplateData{
		IamInstanceProfile: &gfnec2.LaunchTemplate_IamInstanceProfile{
			Arn: n.instanceProfileARN,
		},
		ImageId:           gfnt.NewString(ng.AMI),
		UserData:          gfnt.NewString(userData),
		MetadataOptions:   makeMetadataOptions(ng.NodeGroupBase),
		TagSpecifications: makeTags(ng.NodeGroupBase, n.options.ClusterConfig.Metadata),
	}

	if ng.CapacityReservation != nil {
		valueOrNil := func(value *string) *gfnt.Value {
			if value != nil {
				return gfnt.NewString(*value)
			}
			return nil
		}
		launchTemplateData.CapacityReservationSpecification = &gfnec2.LaunchTemplate_CapacityReservationSpecification{}
		launchTemplateData.CapacityReservationSpecification.CapacityReservationPreference = valueOrNil(ng.CapacityReservation.CapacityReservationPreference)
		if ng.CapacityReservation.CapacityReservationTarget != nil {
			launchTemplateData.CapacityReservationSpecification.CapacityReservationTarget = &gfnec2.LaunchTemplate_CapacityReservationTarget{
				CapacityReservationId:               valueOrNil(ng.CapacityReservation.CapacityReservationTarget.CapacityReservationID),
				CapacityReservationResourceGroupArn: valueOrNil(ng.CapacityReservation.CapacityReservationTarget.CapacityReservationResourceGroupARN),
			}
		}
	}

	if err := buildNetworkInterfaces(ctx, launchTemplateData, ng.InstanceTypeList(), api.IsEnabled(ng.EFAEnabled), n.securityGroups, n.ec2API); err != nil {
		return nil, fmt.Errorf("couldn't build network interfaces for launch template data: %w", err)
	}

	if api.IsEnabled(ng.EFAEnabled) && ng.Placement == nil {
		groupName := n.newResource("NodeGroupPlacementGroup", &gfnec2.PlacementGroup{
			Strategy: gfnt.NewString("cluster"),
		})
		launchTemplateData.Placement = &gfnec2.LaunchTemplate_Placement{
			GroupName: groupName,
		}
	}

	if !api.HasMixedInstances(ng) {
		launchTemplateData.InstanceType = gfnt.NewString(ng.InstanceType)
	} else {
		launchTemplateData.InstanceType = gfnt.NewString(ng.InstancesDistribution.InstanceTypes[0])
	}
	if ng.EBSOptimized != nil {
		launchTemplateData.EbsOptimized = gfnt.NewBoolean(*ng.EBSOptimized)
	}

	if ng.CPUCredits != nil {
		launchTemplateData.CreditSpecification = &gfnec2.LaunchTemplate_CreditSpecification{
			CpuCredits: gfnt.NewString(strings.ToLower(*ng.CPUCredits)),
		}
	}

	if ng.Placement != nil {
		launchTemplateData.Placement = &gfnec2.LaunchTemplate_Placement{
			GroupName: gfnt.NewString(ng.Placement.GroupName),
		}
	}

	if ng.EnableDetailedMonitoring != nil {
		launchTemplateData.Monitoring = &gfnec2.LaunchTemplate_Monitoring{
			Enabled: gfnt.NewBoolean(*ng.EnableDetailedMonitoring),
		}
	}

	return launchTemplateData, nil
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

func nodeGroupResource(launchTemplateName *gfnt.Value, vpcZoneIdentifier interface{}, tags []map[string]string, ng *api.NodeGroup) *awsCloudFormationResource {
	ngProps := map[string]interface{}{
		"VPCZoneIdentifier": vpcZoneIdentifier,
		"Tags":              tags,
	}

	if ng.InstancesDistribution != nil && ng.InstancesDistribution.CapacityRebalance {
		ngProps["CapacityRebalance"] = ng.InstancesDistribution.CapacityRebalance
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

	if ng.MaxInstanceLifetime != nil {
		ngProps["MaxInstanceLifetime"] = *ng.MaxInstanceLifetime
	}

	rollingUpdate := map[string]interface{}{}
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
