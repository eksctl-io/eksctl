package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/utils/strings/slices"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/spotinst/spotinst-sdk-go/spotinst"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfncfn "github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/az"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/spot"
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

// NodeGroupResourceSet stores the resource information of the nodegroup
type NodeGroupResourceSet struct {
	rs                *resourceSet
	clusterSpec       *api.ClusterConfig
	spec              *api.NodeGroup
	forceAddCNIPolicy bool
	ec2API            awsapi.EC2

	iamAPI             awsapi.IAM
	instanceProfileARN *gfnt.Value
	securityGroups     []*gfnt.Value
	vpc                *gfnt.Value
	vpcImporter        vpc.Importer
	bootstrapper       nodebootstrap.Bootstrapper
	sharedTags         []types.Tag
}

// NewNodeGroupResourceSet returns a resource set for a nodegroup embedded in a cluster config
func NewNodeGroupResourceSet(ec2API awsapi.EC2, iamAPI awsapi.IAM, spec *api.ClusterConfig, ng *api.NodeGroup, bootstrapper nodebootstrap.Bootstrapper, sharedTags []types.Tag, forceAddCNIPolicy bool, vpcImporter vpc.Importer) *NodeGroupResourceSet {
	return &NodeGroupResourceSet{
		rs:                newResourceSet(),
		forceAddCNIPolicy: forceAddCNIPolicy,
		clusterSpec:       spec,
		spec:              ng,
		ec2API:            ec2API,
		iamAPI:            iamAPI,
		vpcImporter:       vpcImporter,
		bootstrapper:      bootstrapper,
		sharedTags:        sharedTags,
	}
}

// AddAllResources adds all the information about the nodegroup to the resource set
func (n *NodeGroupResourceSet) AddAllResources(ctx context.Context) error {

	if n.clusterSpec.IPv6Enabled() {
		return errors.New("unmanaged nodegroups are not supported with IPv6 clusters")
	}

	n.rs.template.Description = fmt.Sprintf(
		"%s (AMI family: %s, SSH access: %v, private networking: %v) %s",
		nodeGroupTemplateDescription,
		n.spec.AMIFamily, api.IsEnabled(n.spec.SSH.Allow), n.spec.PrivateNetworking,
		templateDescriptionSuffix)

	n.Template().Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings

	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeaturePrivateNetworking, n.spec.PrivateNetworking, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureSharedSecurityGroup, n.spec.SecurityGroups.WithShared, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureLocalSecurityGroup, n.spec.SecurityGroups.WithLocal, false)

	n.vpc = n.vpcImporter.VPC()

	if n.spec.Tags == nil {
		n.spec.Tags = map[string]string{}
	}

	for k, v := range n.clusterSpec.Metadata.Tags {
		if _, exists := n.spec.Tags[k]; !exists {
			n.spec.Tags[k] = v
		}
	}

	// Ensure MinSize is set, as it is required by the ASG cfn resource
	// TODO this validation and default setting should happen way earlier than this
	if n.spec.MinSize == nil {
		if n.spec.DesiredCapacity == nil {
			defaultNodeCount := api.DefaultNodeCount
			n.spec.MinSize = &defaultNodeCount
		} else {
			n.spec.MinSize = n.spec.DesiredCapacity
		}
		logger.Info("--nodes-min=%d was set automatically for nodegroup %s", *n.spec.MinSize, n.spec.Name)
	} else if n.spec.DesiredCapacity != nil && *n.spec.DesiredCapacity < *n.spec.MinSize {
		return fmt.Errorf("--nodes value (%d) cannot be lower than --nodes-min value (%d)", *n.spec.DesiredCapacity, *n.spec.MinSize)
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
		return fmt.Errorf("--nodes value (%d) cannot be greater than --nodes-max value (%d)", *n.spec.DesiredCapacity, *n.spec.MaxSize)
	} else if *n.spec.MaxSize < *n.spec.MinSize {
		return fmt.Errorf("--nodes-min value (%d) cannot be greater than --nodes-max value (%d)", *n.spec.MinSize, *n.spec.MaxSize)
	}

	// Avoid creating IAM resources for the Ocean Cluster resource set as it
	// will only be used as a template for Ocean Virtual Node Groups.
	if n.spec.Name != api.SpotOceanClusterNodeGroupName {
		if err := n.addResourcesForIAM(ctx); err != nil {
			return err
		}
	}
	n.addResourcesForSecurityGroups()

	return n.addResourcesForNodeGroup(ctx)
}

func (n *NodeGroupResourceSet) addResourcesForSecurityGroups() {
	for _, id := range n.spec.SecurityGroups.AttachIDs {
		n.securityGroups = append(n.securityGroups, gfnt.NewString(id))
	}

	if api.IsEnabled(n.spec.SecurityGroups.WithShared) {
		n.securityGroups = append(n.securityGroups, n.vpcImporter.SharedNodeSecurityGroup())
	}

	if api.IsDisabled(n.spec.SecurityGroups.WithLocal) {
		return
	}

	desc := "worker nodes in group " + n.spec.Name
	vpcID := n.vpcImporter.VPC()
	refControlPlaneSG := n.vpcImporter.ControlPlaneSecurityGroup()

	refNodeGroupLocalSG := n.newResource("SG", &gfnec2.SecurityGroup{
		VpcId:            vpcID,
		GroupDescription: gfnt.NewString("Communication between the control plane and " + desc),
		Tags: []gfncfn.Tag{{
			Key:   gfnt.NewString("kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name),
			Value: gfnt.NewString("owned"),
		}},
		SecurityGroupIngress: makeNodeIngressRules(n.spec.NodeGroupBase, refControlPlaneSG, n.clusterSpec.VPC.CIDR.String(), desc),
	})

	n.securityGroups = append(n.securityGroups, refNodeGroupLocalSG)

	if api.IsEnabled(n.spec.EFAEnabled) {
		efaSG := n.rs.addEFASecurityGroup(vpcID, n.clusterSpec.Metadata.Name, desc)
		n.securityGroups = append(n.securityGroups, efaSG)
	}

	n.newResource("EgressInterCluster", &gfnec2.SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                gfnt.NewString("Allow control plane to communicate with " + desc + " (kubelet and workload TCP ports)"),
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgMinNodePort,
		ToPort:                     sgMaxNodePort,
	})
	n.newResource("EgressInterClusterAPI", &gfnec2.SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                gfnt.NewString("Allow control plane to communicate with " + desc + " (workloads using HTTPS port, commonly used with extension API servers)"),
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgPortHTTPS,
		ToPort:                     sgPortHTTPS,
	})
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
			IpProtocol:            sgProtoTCP,
			FromPort:              sgMinNodePort,
			ToPort:                sgMaxNodePort,
		},
		{
			SourceSecurityGroupId: controlPlaneSG,
			Description:           gfnt.NewString(fmt.Sprintf("[IngressInterClusterAPI] Allow %s to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)", description)),
			IpProtocol:            sgProtoTCP,
			FromPort:              sgPortHTTPS,
			ToPort:                sgPortHTTPS,
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
		return errors.Wrap(err, "could not add resources for nodegroup")
	}

	if n.spec.SSH != nil && api.IsSetAndNonEmptyString(n.spec.SSH.PublicKeyName) {
		launchTemplateData.KeyName = gfnt.NewString(*n.spec.SSH.PublicKeyName)
	}

	launchTemplateData.BlockDeviceMappings = makeBlockDeviceMappings(n.spec.NodeGroupBase)

	launchTemplate := &gfnec2.LaunchTemplate{
		LaunchTemplateName: launchTemplateName,
		LaunchTemplateData: launchTemplateData,
	}

	// Do not create a Launch Template resource for Spot-managed nodegroups.
	if n.spec.SpotOcean == nil {
		n.newResource("NodeGroupLaunchTemplate", launchTemplate)
	}

	vpcZoneIdentifier, err := AssignSubnets(ctx, n.spec, n.clusterSpec, n.ec2API)
	if err != nil {
		return err
	}

	tags := []map[string]string{
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
			map[string]string{
				"Key":               "k8s.io/cluster-autoscaler/enabled",
				"Value":             "true",
				"PropagateAtLaunch": "true",
			},
			map[string]string{
				"Key":               "k8s.io/cluster-autoscaler/" + n.clusterSpec.Metadata.Name,
				"Value":             "owned",
				"PropagateAtLaunch": "true",
			},
		)
	}

	if api.IsEnabled(n.spec.PropagateASGTags) {
		var clusterTags []map[string]string
		GenerateClusterAutoscalerTags(n.spec, func(key, value string) {
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

	g, err := n.newNodeGroupResource(launchTemplate, &vpcZoneIdentifier, tags)

	if g == nil {
		return fmt.Errorf("failed to build nodegroup resource: %v", err)
	}
	n.newResource("NodeGroup", g)

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
	userData, err := n.bootstrapper.UserData()
	if err != nil {
		return nil, err
	}

	launchTemplateData := &gfnec2.LaunchTemplate_LaunchTemplateData{
		IamInstanceProfile: &gfnec2.LaunchTemplate_IamInstanceProfile{
			Arn: n.instanceProfileARN,
		},
		ImageId:           gfnt.NewString(n.spec.AMI),
		UserData:          gfnt.NewString(userData),
		MetadataOptions:   makeMetadataOptions(n.spec.NodeGroupBase),
		TagSpecifications: makeTags(n.spec.NodeGroupBase, n.clusterSpec.Metadata),
	}

	if n.spec.CapacityReservation != nil {
		valueOrNil := func(value *string) *gfnt.Value {
			if value != nil {
				return gfnt.NewString(*value)
			}
			return nil
		}
		launchTemplateData.CapacityReservationSpecification = &gfnec2.LaunchTemplate_CapacityReservationSpecification{}
		launchTemplateData.CapacityReservationSpecification.CapacityReservationPreference = valueOrNil(n.spec.CapacityReservation.CapacityReservationPreference)
		if n.spec.CapacityReservation.CapacityReservationTarget != nil {
			launchTemplateData.CapacityReservationSpecification.CapacityReservationTarget = &gfnec2.LaunchTemplate_CapacityReservationTarget{
				CapacityReservationId:               valueOrNil(n.spec.CapacityReservation.CapacityReservationTarget.CapacityReservationID),
				CapacityReservationResourceGroupArn: valueOrNil(n.spec.CapacityReservation.CapacityReservationTarget.CapacityReservationResourceGroupARN),
			}
		}
	}

	if err := buildNetworkInterfaces(ctx, launchTemplateData, n.spec.InstanceTypeList(), api.IsEnabled(n.spec.EFAEnabled), n.securityGroups, n.ec2API); err != nil {
		return nil, errors.Wrap(err, "couldn't build network interfaces for launch template data")
	}

	if api.IsEnabled(n.spec.EFAEnabled) && n.spec.Placement == nil {
		groupName := n.newResource("NodeGroupPlacementGroup", &gfnec2.PlacementGroup{
			Strategy: gfnt.NewString("cluster"),
		})
		launchTemplateData.Placement = &gfnec2.LaunchTemplate_Placement{
			GroupName: groupName,
		}
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

	if n.spec.EnableDetailedMonitoring != nil {
		launchTemplateData.Monitoring = &gfnec2.LaunchTemplate_Monitoring{
			Enabled: gfnt.NewBoolean(*n.spec.EnableDetailedMonitoring),
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

func (n *NodeGroupResourceSet) newNodeGroupResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, tags []map[string]string) (*awsCloudFormationResource, error) {

	if n.spec.SpotOcean != nil {
		return n.newNodeGroupSpotOceanResource(launchTemplate, vpcZoneIdentifier, tags)
	}

	return nodeGroupResource(launchTemplate.LaunchTemplateName, vpcZoneIdentifier, tags, n.spec), nil
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

// newNodeGroupSpotOceanResource returns a Spot Ocean resource.
func (n *NodeGroupResourceSet) newNodeGroupSpotOceanResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, tags []map[string]string) (*awsCloudFormationResource, error) {

	var res *spot.ResourceNodeGroup
	var out awsCloudFormationResource
	var err error

	// Resource.
	{
		if n.spec.Name == api.SpotOceanClusterNodeGroupName {
			logger.Debug("ocean: building nodegroup %q as cluster", n.spec.Name)
			res, err = n.newNodeGroupSpotOceanClusterResource(
				launchTemplate, vpcZoneIdentifier, tags)
		} else {
			logger.Debug("ocean: building nodegroup %q as virtual node group", n.spec.Name)
			n.populateNodeGroupSpotOceanVirtualNodeGroupResourcesWithClusterConfig()
			res, err = n.newNodeGroupSpotOceanVirtualNodeGroupResource(
				launchTemplate, vpcZoneIdentifier, tags)
		}
		if err != nil {
			return nil, err
		}
	}

	// Service Token.
	{
		res.ServiceToken = gfnt.MakeFnSubString(spot.LoadServiceToken())
	}

	// Credentials.
	{
		token, account, err := spot.LoadCredentials()
		if err != nil {
			return nil, err
		}
		if token != "" {
			res.Token = n.rs.newParameter(spot.CredentialsTokenParameterKey, gfn.Parameter{
				Type:    "String",
				Default: token,
			})
		}
		if account != "" {
			res.Account = n.rs.newParameter(spot.CredentialsAccountParameterKey, gfn.Parameter{
				Type:    "String",
				Default: account,
			})
		}
	}

	// Feature Flags.
	{
		if ff := spot.LoadFeatureFlags(); ff != "" {
			res.FeatureFlags = n.rs.newParameter(spot.FeatureFlagsParameterKey, gfn.Parameter{
				Type:    "String",
				Default: ff,
			})
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

	return &out, err
}

// newNodeGroupSpotOceanClusterResource returns a Spot Ocean Cluster resource.
func (n *NodeGroupResourceSet) newNodeGroupSpotOceanClusterResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, resourceTags []map[string]string) (*spot.ResourceNodeGroup, error) {

	template := launchTemplate.LaunchTemplateData
	cluster := &spot.Cluster{
		Name:      spotinst.String(n.clusterSpec.Metadata.Name),
		ClusterID: spotinst.String(n.clusterSpec.Metadata.Name),
		Region:    gfnt.MakeRef("AWS::Region"),
		Compute: &spot.Compute{
			LaunchSpecification: &spot.VirtualNodeGroup{
				ImageID:           template.ImageId,
				UserData:          template.UserData,
				KeyPair:           template.KeyName,
				EBSOptimized:      template.EbsOptimized,
				UseAsTemplateOnly: spotinst.Bool(true),
			},
			SubnetIDs: vpcZoneIdentifier,
		},
	}

	// Networking.
	{
		if ifaces := template.NetworkInterfaces; len(ifaces) > 0 {
			cluster.Compute.LaunchSpecification.AssociatePublicIPAddress = ifaces[0].AssociatePublicIpAddress
		}
	}

	// Security Groups.
	{
		if len(n.securityGroups) > 0 {
			cluster.Compute.LaunchSpecification.SecurityGroupIDs = gfnt.NewSlice(n.securityGroups...)
		}
	}

	// Load Balancers.
	{
		var lbs []*spot.LoadBalancer

		// ELBs.
		if len(n.spec.ClassicLoadBalancerNames) > 0 {
			for _, name := range n.spec.ClassicLoadBalancerNames {
				lbs = append(lbs, &spot.LoadBalancer{
					Type: spotinst.String("CLASSIC"),
					Name: spotinst.String(name),
				})
			}
		}

		// ALBs.
		if len(n.spec.TargetGroupARNs) > 0 {
			for _, arn := range n.spec.TargetGroupARNs {
				lbs = append(lbs, &spot.LoadBalancer{
					Type: spotinst.String("TARGET_GROUP"),
					ARN:  spotinst.String(arn),
				})
			}
		}

		if len(lbs) > 0 {
			cluster.Compute.LaunchSpecification.LoadBalancers = lbs
		}
	}

	// Tags.
	{
		tagMap := make(map[string]string)

		// Nodegroup tags.
		if len(n.spec.Tags) > 0 {
			for key, value := range n.spec.Tags {
				tagMap[key] = value
			}
		}

		// Resource tags (Name, kubernetes.io/*, k8s.io/*, etc.).
		if len(resourceTags) > 0 {
			for _, tag := range resourceTags {
				tagMap[tag["Key"]] = tag["Value"]
			}
		}

		// Shared tags (metadata.tags + eksctl's tags).
		if len(n.sharedTags) > 0 {
			for _, tag := range n.sharedTags {
				tagMap[spotinst.StringValue(tag.Key)] = spotinst.StringValue(tag.Value)
			}
		}

		if len(tagMap) > 0 {
			tags := make([]*spot.Tag, 0, len(tagMap))
			for k, v := range tagMap {
				tags = append(tags, &spot.Tag{
					Key:   spotinst.String(k),
					Value: spotinst.String(v),
				})
			}
			cluster.Compute.LaunchSpecification.Tags = tags
		}
	}

	if spotOcean := n.clusterSpec.SpotOcean; spotOcean != nil {
		// Strategy.
		{
			if strategy := spotOcean.Strategy; strategy != nil {

				cluster.Strategy = &spot.Strategy{
					SpotPercentage:           strategy.SpotPercentage,
					UtilizeReservedInstances: strategy.UtilizeReservedInstances,
					UtilizeCommitments:       strategy.UtilizeCommitments,
					FallbackToOnDemand:       strategy.FallbackToOnDemand,
				}
				if strategy.ClusterOrientation != nil {
					cluster.Strategy.ClusterOrientation = &spot.ClusterOrientation{
						AvailabilityVsCost: strategy.ClusterOrientation.AvailabilityVsCost,
					}
				}

			}
		}

		// Instance Types.
		{
			if compute := spotOcean.Compute; compute != nil && compute.InstanceTypes != nil {
				cluster.Compute.InstanceTypes = &spot.InstanceTypes{
					Whitelist: compute.InstanceTypes.Whitelist,
					Blacklist: compute.InstanceTypes.Blacklist,
				}
			}
		}

		// Scheduling.
		{
			if scheduling := spotOcean.Scheduling; scheduling != nil {
				if hours := scheduling.ShutdownHours; hours != nil {
					cluster.Scheduling = &spot.Scheduling{
						ShutdownHours: &spot.ShutdownHours{
							IsEnabled:   hours.IsEnabled,
							TimeWindows: hours.TimeWindows,
						},
					}
				}
				if tasks := scheduling.Tasks; len(tasks) > 0 {
					if cluster.Scheduling == nil {
						cluster.Scheduling = new(spot.Scheduling)
					}

					cluster.Scheduling.Tasks = make([]*spot.Task, 0)
					for _, task := range tasks {
						if *task.Type != api.SpotOceanTaskTypeManualHeadroomUpdate {
							clusterTask := &spot.Task{
								IsEnabled:      task.IsEnabled,
								Type:           task.Type,
								CronExpression: task.CronExpression,
							}
							cluster.Scheduling.Tasks = append(cluster.Scheduling.Tasks, clusterTask)
						}
					}
				}
				if cluster.Scheduling != nil && cluster.Scheduling.Tasks != nil && len(cluster.Scheduling.Tasks) == 0 {
					cluster.Scheduling.Tasks = nil
				}
			}
		}

		// Auto Scaler.
		{
			if autoScaler := spotOcean.AutoScaler; autoScaler != nil {
				cluster.AutoScaler = &spot.AutoScaler{
					IsEnabled:    autoScaler.Enabled,
					IsAutoConfig: autoScaler.AutoConfig,
					Cooldown:     autoScaler.Cooldown,
				}
				if h := autoScaler.Headroom; h != nil {
					cluster.AutoScaler.Headroom = &spot.Headroom{
						CPUPerUnit:    h.CPUPerUnit,
						GPUPerUnit:    h.GPUPerUnit,
						MemoryPerUnit: h.MemoryPerUnit,
						NumOfUnits:    h.NumOfUnits,
					}
				}
				if l := autoScaler.ResourceLimits; l != nil {
					cluster.AutoScaler.ResourceLimits = &spot.ResourceLimits{
						MaxVCPU:      l.MaxVCPU,
						MaxMemoryGiB: l.MaxMemoryGiB,
					}
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

	return &spot.ResourceNodeGroup{Cluster: cluster}, nil
}

// newNodeGroupSpotOceanVirtualNodeGroupResource returns a Spot Ocean Virtual Node Group resource.
func (n *NodeGroupResourceSet) newNodeGroupSpotOceanVirtualNodeGroupResource(launchTemplate *gfnec2.LaunchTemplate,
	vpcZoneIdentifier interface{}, resourceTags []map[string]string) (*spot.ResourceNodeGroup, error) {

	// Import the Ocean Cluster identifier.
	oceanClusterStackName := fmt.Sprintf("eksctl-%s-nodegroup-ocean", n.clusterSpec.Metadata.Name)
	oceanClusterID := gfnt.MakeFnImportValueString(fmt.Sprintf("%s::%s",
		oceanClusterStackName,
		outputs.NodeGroupSpotOceanClusterID))

	template := launchTemplate.LaunchTemplateData
	spec := &spot.VirtualNodeGroup{
		Name:      spotinst.String(n.spec.Name),
		OceanID:   oceanClusterID,
		ImageID:   template.ImageId,
		UserData:  template.UserData,
		SubnetIDs: vpcZoneIdentifier,
	}

	// Strategy.
	{
		if strategy := n.spec.SpotOcean.Strategy; strategy != nil {
			spec.Strategy = &spot.Strategy{
				SpotPercentage: strategy.SpotPercentage,
			}
		}
	}

	// Block Device Mappings.
	{
		if devs := template.BlockDeviceMappings; len(devs) > 0 {
			spec.BlockDeviceMappings = make([]*spot.BlockDevice, len(devs))
			for i, d := range devs {
				dev := &spot.BlockDevice{
					DeviceName: d.DeviceName,
				}
				if d.Ebs != nil {
					dev.EBS = &spot.BlockDeviceEBS{
						VolumeSize: d.Ebs.VolumeSize,
						VolumeType: d.Ebs.VolumeType,
						Encrypted:  d.Ebs.Encrypted,
						KMSKeyID:   d.Ebs.KmsKeyId,
						IOPS:       d.Ebs.Iops,
						Throughput: d.Ebs.Throughput,
					}
				}
				spec.BlockDeviceMappings[i] = dev
			}
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
		if ifaces := template.NetworkInterfaces; len(ifaces) > 0 {
			spec.AssociatePublicIPAddress = ifaces[0].AssociatePublicIpAddress
		}
	}

	// Security Groups.
	{
		if len(n.securityGroups) > 0 {
			spec.SecurityGroupIDs = gfnt.NewSlice(n.securityGroups...)
		}
	}

	// Tags.
	{
		tagMap := make(map[string]string)

		// Nodegroup tags.
		if len(n.spec.Tags) > 0 {
			for k, v := range n.spec.Tags {
				tagMap[k] = v
			}
		}

		// Resource tags (Name, kubernetes.io/*, k8s.io/*, etc.).
		if len(resourceTags) > 0 {
			for _, tag := range resourceTags {
				tagMap[tag["Key"]] = tag["Value"]
			}
		}

		// Shared tags (metadata.tags + eksctl's tags).
		if len(n.sharedTags) > 0 {
			for _, tag := range n.sharedTags {
				tagMap[spotinst.StringValue(tag.Key)] = spotinst.StringValue(tag.Value)
			}
		}

		if len(tagMap) > 0 {
			tags := make([]*spot.Tag, 0, len(tagMap))
			for k, v := range tagMap {
				tags = append(tags, &spot.Tag{
					Key:   spotinst.String(k),
					Value: spotinst.String(v),
				})
			}
			spec.Tags = tags
		}
	}

	// Instance Types.
	{
		if compute := n.spec.SpotOcean.Compute; compute != nil && compute.InstanceTypes != nil {
			spec.InstanceTypes = compute.InstanceTypes
		}
	}

	// Instance Metadata Options.
	{
		if compute := n.spec.SpotOcean.Compute; compute != nil && compute.InstanceMetadataOptions != nil {
			spec.InstanceMetadataOptions = &spot.InstanceMetadataOptions{
				HttpPutResponseHopLimit: compute.InstanceMetadataOptions.HttpPutResponseHopLimit,
				HttpTokens:              compute.InstanceMetadataOptions.HttpTokens,
			}
		}
	}

	// Scheduling.
	{
		if scheduling := n.spec.SpotOcean.Scheduling; scheduling != nil {
			if hours := scheduling.ShutdownHours; hours != nil {
				spec.Scheduling = &spot.Scheduling{
					ShutdownHours: &spot.ShutdownHours{
						IsEnabled:   hours.IsEnabled,
						TimeWindows: hours.TimeWindows,
					},
				}
			}
			if tasks := scheduling.Tasks; len(tasks) > 0 {
				if spec.Scheduling == nil {
					spec.Scheduling = new(spot.Scheduling)
				}

				spec.Scheduling.Tasks = make([]*spot.Task, len(tasks))
				for i, task := range tasks {
					var headrooms []*spot.Headroom

					if config := task.Config; config != nil && config.Headrooms != nil {
						headrooms = make([]*spot.Headroom, len(config.Headrooms))

						for j, SpotOceanHeadroom := range config.Headrooms {
							headrooms[j] = &spot.Headroom{
								CPUPerUnit:    SpotOceanHeadroom.CPUPerUnit,
								GPUPerUnit:    SpotOceanHeadroom.GPUPerUnit,
								MemoryPerUnit: SpotOceanHeadroom.MemoryPerUnit,
								NumOfUnits:    SpotOceanHeadroom.NumOfUnits,
							}
						}
					}

					spec.Scheduling.Tasks[i] = &spot.Task{
						IsEnabled:      task.IsEnabled,
						Type:           task.Type,
						CronExpression: task.CronExpression,
						Config: &spot.TaskConfig{
							Headrooms: headrooms,
						},
					}
				}
			}
		}
	}

	// Labels.
	{
		if len(n.spec.Labels) > 0 {
			labels := make([]*spot.Label, 0, len(n.spec.Labels))

			for key, value := range n.spec.Labels {
				labels = append(labels, &spot.Label{
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
			taints := make([]*spot.Taint, len(n.spec.Taints))

			for i, t := range n.spec.Taints {
				taints[i] = &spot.Taint{
					Key:    spotinst.String(t.Key),
					Value:  spotinst.String(t.Value),
					Effect: spotinst.String(string(t.Effect)),
				}
			}

			spec.Taints = taints
		}
	}

	// Auto Scaler.
	{
		if autoScaler := n.spec.SpotOcean.AutoScaler; autoScaler != nil {
			if len(autoScaler.Headrooms) > 0 {
				headrooms := make([]*spot.Headroom, len(autoScaler.Headrooms))

				for i, h := range autoScaler.Headrooms {
					headrooms[i] = &spot.Headroom{
						CPUPerUnit:    h.CPUPerUnit,
						GPUPerUnit:    h.GPUPerUnit,
						MemoryPerUnit: h.MemoryPerUnit,
						NumOfUnits:    h.NumOfUnits,
					}
				}

				spec.AutoScaler = &spot.AutoScaler{
					Headrooms: headrooms,
				}
			}

			if autoScaler.ResourceLimits != nil {
				spec.ResourceLimits = &spot.ResourceLimits{
					MinInstanceCount: autoScaler.ResourceLimits.MinInstanceCount,
					MaxInstanceCount: autoScaler.ResourceLimits.MaxInstanceCount,
				}
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

	// Initial nodes.
	{
		if len(n.spec.Taints) == 0 {
			if n.spec.MinSize == nil && n.spec.DesiredCapacity != nil {
				n.spec.MinSize = n.spec.DesiredCapacity
			}
			if spotinst.IntValue(n.spec.MinSize) == 0 {
				initialNodes := api.DefaultNodeCount
				n.spec.MinSize = &initialNodes
			}
		}
	}

	return &spot.ResourceNodeGroup{
		VirtualNodeGroup: spec,
		Resource: spot.Resource{
			Parameters: spot.ResourceParameters{
				OnCreate: map[string]interface{}{
					"initialNodes": spotinst.IntValue(n.spec.MinSize),
				},
				OnDelete: map[string]interface{}{
					"deleteNodes": true,
					"forceDelete": true,
				},
			},
		},
	}, nil
}

func (n *NodeGroupResourceSet) populateNodeGroupSpotOceanVirtualNodeGroupResourcesWithClusterConfig() {
	clusterSpec := n.clusterSpec.SpotOcean
	launchSpec := n.spec.SpotOcean

	if clusterSpec != nil {

		// Instance Metadata Options.
		if compute := clusterSpec.Compute; compute != nil && compute.InstanceMetadataOptions != nil &&
			(launchSpec.Compute == nil || launchSpec.Compute.InstanceMetadataOptions == nil) {
			if launchSpec.Compute == nil {
				launchSpec.Compute = new(api.SpotOceanVirtualNodeGroupCompute)
			}

			launchSpec.Compute.InstanceMetadataOptions = &api.InstanceMetadataOptions{
				HttpPutResponseHopLimit: compute.InstanceMetadataOptions.HttpPutResponseHopLimit,
				HttpTokens:              compute.InstanceMetadataOptions.HttpTokens,
			}
		}

		// Scheduling.
		if scheduling := clusterSpec.Scheduling; scheduling != nil && launchSpec.Scheduling == nil {
			launchSpec.Scheduling = new(api.SpotOceanClusterScheduling)

			if tasks := scheduling.Tasks; len(tasks) > 0 {
				launchSpec.Scheduling.Tasks = make([]*api.SpotOceanTask, len(tasks))
				for i, task := range tasks {
					var headrooms []*api.SpotOceanHeadroom

					launchSpec.Scheduling.Tasks[i] = &api.SpotOceanTask{
						IsEnabled:      task.IsEnabled,
						Type:           task.Type,
						CronExpression: task.CronExpression,
					}

					if config := task.Config; config != nil && config.Headrooms != nil {
						headrooms = make([]*api.SpotOceanHeadroom, len(config.Headrooms))

						for j, SpotOceanHeadroom := range config.Headrooms {
							headrooms[j] = &api.SpotOceanHeadroom{
								CPUPerUnit:    SpotOceanHeadroom.CPUPerUnit,
								GPUPerUnit:    SpotOceanHeadroom.GPUPerUnit,
								MemoryPerUnit: SpotOceanHeadroom.MemoryPerUnit,
								NumOfUnits:    SpotOceanHeadroom.NumOfUnits,
							}
						}

						launchSpec.Scheduling.Tasks[i].Config = &api.SpotOceanTaskConfig{
							Headrooms: headrooms,
						}
					}
				}
			}
		}
	}
}
