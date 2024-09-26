package builder

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfncfn "github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs             *resourceSet
	spec           *api.ClusterConfig
	ec2API         awsapi.EC2
	region         string
	vpcResourceSet VPCResourceSet
	securityGroups *gfnt.Value
}

// NewClusterResourceSet returns a resource set for the new cluster.
func NewClusterResourceSet(ec2API awsapi.EC2, region string, spec *api.ClusterConfig, existingStack *gjson.Result, extendForOutposts bool) *ClusterResourceSet {
	var usesExistingVPC bool
	if existingStack != nil {
		unsetExistingResources(existingStack, spec)
		usesExistingVPC = !existingStack.Get(cfnVPCResource).Exists()
	} else {
		usesExistingVPC = spec.VPC.ID != ""
	}

	var (
		vpcResourceSet VPCResourceSet
		rs             = newResourceSet()
	)

	switch {
	case usesExistingVPC:
		vpcResourceSet = NewExistingVPCResourceSet(rs, spec, ec2API)
	case spec.IPv6Enabled():
		vpcResourceSet = NewIPv6VPCResourceSet(rs, spec, ec2API)
	default:
		vpcResourceSet = NewIPv4VPCResourceSet(rs, spec, ec2API, extendForOutposts)
	}

	return &ClusterResourceSet{
		rs:             rs,
		spec:           spec,
		ec2API:         ec2API,
		region:         region,
		vpcResourceSet: vpcResourceSet,
	}
}

// AddAllResources adds all the information about the cluster to the resource set
func (c *ClusterResourceSet) AddAllResources(ctx context.Context) error {
	if err := c.spec.HasSufficientSubnets(); err != nil {
		return err
	}

	vpcID, subnetDetails, err := c.vpcResourceSet.CreateTemplate(ctx)
	if err != nil {
		return errors.Wrap(err, "error adding VPC resources")
	}

	clusterSG := c.addResourcesForSecurityGroups(vpcID)

	if privateCluster := c.spec.PrivateCluster; privateCluster.Enabled && !privateCluster.SkipEndpointCreation {
		vpcEndpointResourceSet := NewVPCEndpointResourceSet(c.ec2API, c.region, c.rs, c.spec, vpcID, subnetDetails.Private, clusterSG.ClusterSharedNode)

		if err := vpcEndpointResourceSet.AddResources(ctx); err != nil {
			return errors.Wrap(err, "error adding resources for VPC endpoints")
		}
	}

	c.addResourcesForIAM()
	c.addResourcesForControlPlane(subnetDetails)

	if len(c.spec.FargateProfiles) > 0 {
		c.addResourcesForFargate()
	}

	c.rs.defineOutput(outputs.ClusterStackName, gfnt.RefStackName, false, func(v string) error {
		if c.spec.Status == nil {
			c.spec.Status = &api.ClusterStatus{}
		}
		c.spec.Status.StackName = v
		return nil
	})

	c.Template().Mappings[servicePrincipalPartitionMapName] = api.Partitions.ServicePrincipalPartitionMappings()

	c.rs.template.Description = fmt.Sprintf(
		"%s (dedicated VPC: %v, dedicated IAM: %v) %s",
		clusterTemplateDescription,
		c.spec.VPC.ID == "",
		c.rs.withIAM,
		templateDescriptionSuffix,
	)

	return nil
}

func (c *ClusterResourceSet) addResourcesForSecurityGroups(vpcID *gfnt.Value) *clusterSecurityGroup {
	var refControlPlaneSG, refClusterSharedNodeSG *gfnt.Value

	if sg := c.spec.VPC.SecurityGroup; sg != "" {
		refControlPlaneSG = gfnt.NewString(sg)
		c.securityGroups = gfnt.NewStringSlice(sg)
	} else if securityGroupIDs := c.spec.VPC.ControlPlaneSecurityGroupIDs; len(securityGroupIDs) > 0 {
		refControlPlaneSG = gfnt.NewString(securityGroupIDs[0])
		c.securityGroups = gfnt.NewStringSlice(securityGroupIDs...)
	} else {
		refControlPlaneSG = c.newResource(cfnControlPlaneSGResource, &gfnec2.SecurityGroup{
			GroupDescription: gfnt.NewString("Communication between the control plane and worker nodegroups"),
			VpcId:            vpcID,
		})

		if len(c.spec.VPC.ExtraCIDRs) > 0 {
			for i, cidr := range c.spec.VPC.ExtraCIDRs {
				c.newResource(fmt.Sprintf("IngressControlPlaneExtraCIDR%d", i), &gfnec2.SecurityGroupIngress{
					GroupId:     refControlPlaneSG,
					CidrIp:      gfnt.NewString(cidr),
					Description: gfnt.NewString(fmt.Sprintf("Allow Extra CIDR %d (%s) to communicate to controlplane", i, cidr)),
					IpProtocol:  gfnt.NewString("tcp"),
					FromPort:    sgPortHTTPS,
					ToPort:      sgPortHTTPS,
				})
			}
		}

		if len(c.spec.VPC.ExtraIPv6CIDRs) > 0 {
			for i, cidr := range c.spec.VPC.ExtraIPv6CIDRs {
				c.newResource(fmt.Sprintf("IngressControlPlaneExtraIPv6CIDR%d", i), &gfnec2.SecurityGroupIngress{
					GroupId:     refControlPlaneSG,
					CidrIpv6:    gfnt.NewString(cidr),
					Description: gfnt.NewString(fmt.Sprintf("Allow Extra IPv6 CIDR %d (%s) to communicate to controlplane", i, cidr)),
					IpProtocol:  gfnt.NewString("tcp"),
					FromPort:    sgPortHTTPS,
					ToPort:      sgPortHTTPS,
				})
			}
		}
		c.securityGroups = gfnt.NewSlice(refControlPlaneSG)
	}

	if c.spec.VPC.SharedNodeSecurityGroup == "" {
		refClusterSharedNodeSG = c.newResource(cfnSharedNodeSGResource, &gfnec2.SecurityGroup{
			GroupDescription: gfnt.NewString("Communication between all nodes in the cluster"),
			VpcId:            vpcID,
		})
		c.newResource("IngressInterNodeGroupSG", &gfnec2.SecurityGroupIngress{
			GroupId:               refClusterSharedNodeSG,
			SourceSecurityGroupId: refClusterSharedNodeSG,
			Description:           gfnt.NewString("Allow nodes to communicate with each other (all ports)"),
			IpProtocol:            gfnt.NewString("-1"),
			FromPort:              sgPortZero,
			ToPort:                sgMaxNodePort,
		})
	} else {
		refClusterSharedNodeSG = gfnt.NewString(c.spec.VPC.SharedNodeSecurityGroup)
	}

	if api.IsEnabled(c.spec.VPC.ManageSharedNodeSecurityGroupRules) {
		// To enable communication between both managed and unmanaged nodegroups, this allows ingress traffic from
		// the default cluster security group ID that EKS creates by default
		// EKS attaches this to Managed Nodegroups by default, but we need to handle this for unmanaged nodegroups
		c.newResource(cfnIngressClusterToNodeSGResource, &gfnec2.SecurityGroupIngress{
			GroupId:               refClusterSharedNodeSG,
			SourceSecurityGroupId: gfnt.MakeFnGetAttString("ControlPlane", outputs.ClusterDefaultSecurityGroup),
			Description:           gfnt.NewString("Allow managed and unmanaged nodes to communicate with each other (all ports)"),
			IpProtocol:            gfnt.NewString("-1"),
			FromPort:              sgPortZero,
			ToPort:                sgMaxNodePort,
		})
		if c.spec.IsControlPlaneOnOutposts() && c.spec.IsFullyPrivate() {
			if subnets := c.spec.VPC.Subnets; subnets != nil && subnets.Private != nil {
				for az, subnet := range subnets.Private {
					c.newResource(fmt.Sprintf("IngressPrivateSubnet%s", formatAZ(az)), &gfnec2.SecurityGroupIngress{
						GroupId:     refClusterSharedNodeSG,
						CidrIp:      gfnt.NewString(subnet.CIDR.String()),
						Description: gfnt.NewString("Allow private subnets to communicate with VPC endpoints"),
						IpProtocol:  gfnt.NewString("tcp"),
						FromPort:    sgPortHTTPS,
						ToPort:      sgPortHTTPS,
					})
				}
			}

		}
		c.newResource("IngressNodeToDefaultClusterSG", &gfnec2.SecurityGroupIngress{
			GroupId:               gfnt.MakeFnGetAttString("ControlPlane", outputs.ClusterDefaultSecurityGroup),
			SourceSecurityGroupId: refClusterSharedNodeSG,
			Description:           gfnt.NewString("Allow unmanaged nodes to communicate with control plane (all ports)"),
			IpProtocol:            gfnt.NewString("-1"),
			FromPort:              sgPortZero,
			ToPort:                sgMaxNodePort,
		})
	}

	if c.spec.VPC == nil {
		c.spec.VPC = &api.ClusterVPC{}
	}
	c.rs.defineOutput(outputs.ClusterSecurityGroup, refControlPlaneSG, true, func(v string) error {
		c.spec.VPC.SecurityGroup = v
		return nil
	})
	c.rs.defineOutput(outputs.ClusterSharedNodeSecurityGroup, refClusterSharedNodeSG, true, func(v string) error {
		c.spec.VPC.SharedNodeSecurityGroup = v
		return nil
	})

	return &clusterSecurityGroup{
		ControlPlane:      refControlPlaneSG,
		ClusterSharedNode: refClusterSharedNodeSG,
	}
}

// RenderJSON returns the rendered JSON
func (c *ClusterResourceSet) RenderJSON() ([]byte, error) {
	return c.rs.renderJSON()
}

// Template returns the CloudFormation template
func (c *ClusterResourceSet) Template() gfn.Template {
	return *c.rs.template
}

// GetAllOutputs collects all outputs of the cluster
func (c *ClusterResourceSet) GetAllOutputs(stack types.Stack) error {
	return c.rs.GetAllOutputs(stack)
}

// HasManagedNodesSG reports whether the stack has the security group required for communication between
// managed and unmanaged nodegroups
func HasManagedNodesSG(stackResources *gjson.Result) bool {
	return stackResources.Get(cfnIngressClusterToNodeSGResource).Exists()
}

// unsetExistingResources unsets fields for CloudFormation resources that were created by eksctl (and not user-supplied)
// in order to trigger execution of code that relies on these fields
func unsetExistingResources(existingStack *gjson.Result, clusterConfig *api.ClusterConfig) {
	controlPlaneSG := existingStack.Get(cfnControlPlaneSGResource)
	if controlPlaneSG.Exists() {
		clusterConfig.VPC.SecurityGroup = ""
	}
	sharedNodeSG := existingStack.Get(cfnSharedNodeSGResource)
	if sharedNodeSG.Exists() {
		clusterConfig.VPC.SharedNodeSecurityGroup = ""
	}
}

func (c *ClusterResourceSet) newResource(name string, resource gfn.Resource) *gfnt.Value {
	return c.rs.newResource(name, resource)
}

func (c *ClusterResourceSet) addResourcesForControlPlane(subnetDetails *SubnetDetails) {
	clusterVPC := &gfneks.Cluster_ResourcesVpcConfig{
		EndpointPublicAccess:  gfnt.NewBoolean(*c.spec.VPC.ClusterEndpoints.PublicAccess),
		EndpointPrivateAccess: gfnt.NewBoolean(*c.spec.VPC.ClusterEndpoints.PrivateAccess),
		SecurityGroupIds:      c.securityGroups,
		PublicAccessCidrs:     gfnt.NewStringSlice(c.spec.VPC.PublicAccessCIDRs...),
	}
	if subnetIDs := c.spec.VPC.ControlPlaneSubnetIDs; len(subnetIDs) > 0 {
		clusterVPC.SubnetIds = gfnt.NewStringSlice(subnetIDs...)
	} else {
		clusterVPC.SubnetIds = gfnt.NewSlice(subnetDetails.ControlPlaneSubnetRefs()...)
	}

	serviceRoleARN := gfnt.MakeFnGetAttString("ServiceRole", "Arn")
	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRoleARN) {
		serviceRoleARN = gfnt.NewString(*c.spec.IAM.ServiceRoleARN)
	}

	var encryptionConfigs []gfneks.Cluster_EncryptionConfig
	if c.spec.SecretsEncryption != nil && c.spec.SecretsEncryption.KeyARN != "" {
		encryptionConfigs = []gfneks.Cluster_EncryptionConfig{
			{
				Resources: gfnt.NewSlice(gfnt.NewString("secrets")),
				Provider: &gfneks.Cluster_Provider{
					KeyArn: gfnt.NewString(c.spec.SecretsEncryption.KeyARN),
				},
			},
		}
	}

	cluster := gfneks.Cluster{
		EncryptionConfig:           encryptionConfigs,
		Logging:                    makeClusterLogging(c.spec),
		Name:                       gfnt.NewString(c.spec.Metadata.Name),
		ResourcesVpcConfig:         clusterVPC,
		RoleArn:                    serviceRoleARN,
		BootstrapSelfManagedAddons: gfnt.NewBoolean(false),
		AccessConfig: &gfneks.Cluster_AccessConfig{
			AuthenticationMode:                      gfnt.NewString(string(c.spec.AccessConfig.AuthenticationMode)),
			BootstrapClusterCreatorAdminPermissions: gfnt.NewBoolean(!api.IsDisabled(c.spec.AccessConfig.BootstrapClusterCreatorAdminPermissions)),
		},
		Tags:    makeCFNTags(c.spec),
		Version: gfnt.NewString(c.spec.Metadata.Version),
	}

	if c.spec.IsControlPlaneOnOutposts() {
		cluster.OutpostConfig = &gfneks.Cluster_OutpostConfig{
			OutpostArns:              gfnt.NewStringSlice(c.spec.Outpost.ControlPlaneOutpostARN),
			ControlPlaneInstanceType: gfnt.NewString(c.spec.Outpost.ControlPlaneInstanceType),
		}
		if c.spec.Outpost.HasPlacementGroup() {
			cluster.OutpostConfig.ControlPlanePlacement = &gfneks.Cluster_ControlPlanePlacement{
				GroupName: gfnt.NewString(c.spec.Outpost.ControlPlanePlacement.GroupName),
			}
		}
	}

	kubernetesNetworkConfig := &gfneks.Cluster_KubernetesNetworkConfig{}
	if knc := c.spec.KubernetesNetworkConfig; knc != nil {
		if knc.ServiceIPv4CIDR != "" {
			kubernetesNetworkConfig.ServiceIpv4Cidr = gfnt.NewString(knc.ServiceIPv4CIDR)
		}

		ipFamily := knc.IPFamily
		if ipFamily == "" {
			ipFamily = api.IPV4Family
		}
		kubernetesNetworkConfig.IpFamily = gfnt.NewString(strings.ToLower(ipFamily))
	}
	cluster.KubernetesNetworkConfig = kubernetesNetworkConfig
	if c.spec.ZonalShiftConfig != nil && api.IsEnabled(c.spec.ZonalShiftConfig.Enabled) {
		cluster.ZonalShiftConfig = &gfneks.Cluster_ZonalShift{
			Enabled: gfnt.NewBoolean(true),
		}
	}

	c.newResource("ControlPlane", &cluster)

	if c.spec.Status == nil {
		c.spec.Status = &api.ClusterStatus{}
	}

	c.rs.defineOutputFromAtt(outputs.ClusterCertificateAuthorityData, "ControlPlane", "CertificateAuthorityData", false, func(v string) error {
		caData, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return errors.Wrap(err, "decoding certificate authority data")
		}
		c.spec.Status.CertificateAuthorityData = caData
		return nil
	})
	c.rs.defineOutputFromAtt(outputs.ClusterEndpoint, "ControlPlane", "Endpoint", true, func(v string) error {
		c.spec.Status.Endpoint = v
		return nil
	})
	c.rs.defineOutputFromAtt(outputs.ClusterARN, "ControlPlane", "Arn", true, func(v string) error {
		c.spec.Status.ARN = v
		return nil
	})

	// This exports the cluster security group ID that EKS creates by default. To enable communication between both
	// managed and unmanaged nodegroups, they must share a security group.
	// EKS attaches this to Managed Nodegroups by default, but we need to add this for unmanaged nodegroups.
	// This exported value is imported in the CloudFormation resource for unmanaged nodegroups
	c.rs.defineOutputFromAtt(outputs.ClusterDefaultSecurityGroup, "ControlPlane", "ClusterSecurityGroupId",
		true, func(s string) error {
			return nil
		})
}

func makeCFNTags(clusterConfig *api.ClusterConfig) []gfncfn.Tag {
	var tags []gfncfn.Tag
	for k, v := range clusterConfig.Metadata.Tags {
		tags = append(tags, gfncfn.Tag{
			Key:   gfnt.NewString(k),
			Value: gfnt.NewString(v),
		})
	}
	return tags
}

func (c *ClusterResourceSet) addResourcesForFargate() {
	_ = addResourcesForFargate(c.rs, c.spec)
}
