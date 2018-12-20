package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs             *resourceSet
	spec           *api.ClusterConfig
	provider       api.ClusterProvider
	vpc            *gfn.Value
	subnets        map[api.SubnetTopology][]*gfn.Value
	securityGroups []*gfn.Value
	outputs        *ClusterStackOutputs
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig) *ClusterResourceSet {
	return &ClusterResourceSet{
		rs:       newResourceSet(),
		spec:     spec,
		provider: provider,
		outputs:  &ClusterStackOutputs{},
	}
}

// AddAllResources adds all the information about the cluster to the resource set
func (c *ClusterResourceSet) AddAllResources() error {
	dedicatedVPC := c.spec.VPC.ID == ""

	c.rs.template.Description = fmt.Sprintf(
		"%s (dedicated VPC: %v, dedicated IAM: %v) %s",
		clusterTemplateDescription,
		dedicatedVPC, true,
		templateDescriptionSuffix)

	if err := c.spec.HasSufficientSubnets(); err != nil {
		return err
	}

	if dedicatedVPC {
		c.addResourcesForVPC()
	} else {
		c.importResourcesForVPC()
	}
	c.addOutputsForVPC()

	c.addResourcesForSecurityGroups()
	c.addResourcesForIAM()
	c.addResourcesForControlPlane()

	c.rs.newOutput(cfnOutputClusterStackName, gfn.RefStackName, false)

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

func (c *ClusterResourceSet) addResourcesForControlPlane() {
	clusterVPC := &gfn.AWSEKSCluster_ResourcesVpcConfig{
		SecurityGroupIds: c.securityGroups,
	}
	for topology := range c.spec.VPC.Subnets {
		clusterVPC.SubnetIds = append(clusterVPC.SubnetIds, c.subnets[topology]...)
	}

	c.newResource("ControlPlane", &gfn.AWSEKSCluster{
		Name:               gfn.NewString(c.spec.Metadata.Name),
		RoleArn:            gfn.MakeFnGetAttString("ServiceRole.Arn"),
		Version:            gfn.NewString(c.spec.Metadata.Version),
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

	if err := vpc.UseSubnets(c.provider, c.spec, api.SubnetTopologyPrivate, c.outputs.SubnetsPrivate); err != nil {
		return err
	}

	if err := vpc.UseSubnets(c.provider, c.spec, api.SubnetTopologyPublic, c.outputs.SubnetsPublic); err != nil {
		return err
	}

	c.spec.ClusterStackName = c.outputs.ClusterStackName
	c.spec.Endpoint = c.outputs.Endpoint
	c.spec.CertificateAuthorityData = c.outputs.CertificateAuthorityData
	c.spec.ARN = c.outputs.ARN

	return nil
}
