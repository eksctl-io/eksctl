package builder

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs                   *resourceSet
	spec                 *api.ClusterConfig
	ec2API               ec2iface.EC2API
	region               string
	supportsManagedNodes bool
	vpcResourceSet       *VPCResourceSet
	securityGroups       []*gfnt.Value
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(ec2API ec2iface.EC2API, region string, spec *api.ClusterConfig, supportsManagedNodes bool, existingStack *gjson.Result) *ClusterResourceSet {
	if existingStack != nil {
		unsetExistingResources(existingStack, spec)
	}
	rs := newResourceSet()
	return &ClusterResourceSet{
		rs:                   rs,
		spec:                 spec,
		ec2API:               ec2API,
		region:               region,
		supportsManagedNodes: supportsManagedNodes,
		vpcResourceSet:       NewVPCResourceSet(rs, spec, ec2API),
	}
}

// AddAllResources adds all the information about the cluster to the resource set
func (c *ClusterResourceSet) AddAllResources() error {
	if err := c.spec.HasSufficientSubnets(); err != nil {
		return err
	}

	vpcResource, err := c.vpcResourceSet.AddResources()
	if err != nil {
		return errors.Wrap(err, "error adding VPC resources")
	}

	c.vpcResourceSet.AddOutputs()
	clusterSG := c.addResourcesForSecurityGroups(vpcResource)

	if privateCluster := c.spec.PrivateCluster; privateCluster.Enabled && !privateCluster.SkipEndpointCreation {
		vpcEndpointResourceSet := NewVPCEndpointResourceSet(c.ec2API, c.region, c.rs, c.spec, vpcResource.VPC, vpcResource.SubnetDetails.Private, clusterSG.ClusterSharedNode)

		if err := vpcEndpointResourceSet.AddResources(); err != nil {
			return errors.Wrap(err, "error adding resources for VPC endpoints")
		}
	}

	c.addResourcesForIAM()
	c.addResourcesForControlPlane(vpcResource.SubnetDetails)

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

	c.Template().Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings

	c.rs.template.Description = fmt.Sprintf(
		"%s (dedicated VPC: %v, dedicated IAM: %v) %s",
		clusterTemplateDescription,
		c.spec.VPC.ID == "",
		c.rs.withIAM,
		templateDescriptionSuffix,
	)

	return nil
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
func (c *ClusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
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

func (c *ClusterResourceSet) addResourcesForControlPlane(subnetDetails *subnetDetails) {
	clusterVPC := &gfneks.Cluster_ResourcesVpcConfig{
		SecurityGroupIds: gfnt.NewSlice(c.securityGroups...),
	}

	clusterVPC.SubnetIds = gfnt.NewSlice(append(subnetDetails.PublicSubnetRefs(), subnetDetails.PrivateSubnetRefs()...)...)

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
		Name:               gfnt.NewString(c.spec.Metadata.Name),
		RoleArn:            serviceRoleARN,
		Version:            gfnt.NewString(c.spec.Metadata.Version),
		ResourcesVpcConfig: clusterVPC,
		EncryptionConfig:   encryptionConfigs,
	}
	if c.spec.KubernetesNetworkConfig != nil && c.spec.KubernetesNetworkConfig.ServiceIPv4CIDR != "" {
		cluster.KubernetesNetworkConfig = &gfneks.Cluster_KubernetesNetworkConfig{
			ServiceIpv4Cidr: gfnt.NewString(c.spec.KubernetesNetworkConfig.ServiceIPv4CIDR),
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

	if c.supportsManagedNodes {
		// This exports the cluster security group ID that EKS creates by default. To enable communication between both
		// managed and unmanaged nodegroups, they must share a security group.
		// EKS attaches this to Managed Nodegroups by default, but we need to add this for unmanaged nodegroups.
		// This exported value is imported in the CloudFormation resource for unmanaged nodegroups
		c.rs.defineOutputFromAtt(outputs.ClusterDefaultSecurityGroup, "ControlPlane", "ClusterSecurityGroupId",
			true, func(s string) error {
				return nil
			})
	}
}

func (c *ClusterResourceSet) addResourcesForFargate() {
	_ = addResourcesForFargate(c.rs, c.spec)
}
