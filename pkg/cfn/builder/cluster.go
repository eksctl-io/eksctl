package builder

import (
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	"github.com/weaveworks/eksctl/pkg/eks/api"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs             *resourceSet
	spec           *api.ClusterConfig
	vpc            *gfn.Value
	subnets        map[api.SubnetTopology][]*gfn.Value
	securityGroups []*gfn.Value
	outputs        *ClusterStackOutputs
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(spec *api.ClusterConfig) *ClusterResourceSet {
	return &ClusterResourceSet{
		rs:      newResourceSet(),
		spec:    spec,
		outputs: &ClusterStackOutputs{},
	}
}

// AddAllResources adds all the information about the cluster to the resource set
func (c *ClusterResourceSet) AddAllResources() error {

	templateDescriptionFeatures := clusterTemplateDescriptionDefaultFeatures

	if err := c.spec.HasSufficientSubnets(); err != nil {
		return err
	}

	if c.spec.VPC.ID != "" {
		c.importResourcesForVPC()
		templateDescriptionFeatures = " (with shared VPC and dedicated IAM role) "
	} else {
		c.addResourcesForVPC()
	}
	c.addOutputsForVPC()

	c.addResourcesForSecurityGroups()
	c.addResourcesForIAM()
	c.addResourcesForControlPlane("1.10")

	c.rs.newOutput(cfnOutputClusterStackName, gfn.RefStackName, false)

	c.rs.template.Description = clusterTemplateDescription
	c.rs.template.Description += templateDescriptionFeatures
	c.rs.template.Description += templateDescriptionSuffix

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

func (c *ClusterResourceSet) newResource(name string, resource interface{}) *gfn.Value {
	return c.rs.newResource(name, resource)
}

func (c *ClusterResourceSet) addResourcesForControlPlane(version string) {
	clusterVPC := &gfn.AWSEKSCluster_ResourcesVpcConfig{
		SecurityGroupIds: c.securityGroups,
	}
	for topology := range c.spec.VPC.Subnets {
		clusterVPC.SubnetIds = append(clusterVPC.SubnetIds, c.subnets[topology]...)
	}

	c.newResource("ControlPlane", &gfn.AWSEKSCluster{
		Name:               gfn.NewString(c.spec.ClusterName),
		RoleArn:            gfn.MakeFnGetAttString("ServiceRole.Arn"),
		Version:            gfn.NewString(version),
		ResourcesVpcConfig: clusterVPC,
	})

	c.rs.newOutputFromAtt(cfnOutputClusterCertificateAuthorityData, "ControlPlane.CertificateAuthorityData", false)
	c.rs.newOutputFromAtt(cfnOutputClusterEndpoint, "ControlPlane.Endpoint", true)
	c.rs.newOutputFromAtt(cfnOutputClusterARN, "ControlPlane.Arn", true)
}

// GetAllOutputs collects all outputs of the cluster
func (c *ClusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	if err := c.rs.GetAllOutputs(stack, c.outputs); err != nil {
		return err
	}

	c.spec.VPC.ID = c.outputs.VPC
	c.spec.VPC.SecurityGroup = c.outputs.SecurityGroup

	// TODO: shouldn't assume the order - https://github.com/weaveworks/eksctl/issues/293
	for i, subnet := range c.outputs.SubnetsPrivate {
		c.spec.ImportSubnet(api.SubnetTopologyPrivate, c.spec.AvailabilityZones[i], subnet)
	}

	for i, subnet := range c.outputs.SubnetsPublic {
		c.spec.ImportSubnet(api.SubnetTopologyPublic, c.spec.AvailabilityZones[i], subnet)
	}

	c.spec.ClusterStackName = c.outputs.ClusterStackName
	c.spec.Endpoint = c.outputs.Endpoint
	c.spec.CertificateAuthorityData = c.outputs.CertificateAuthorityData
	c.spec.ARN = c.outputs.ARN

	return nil
}
