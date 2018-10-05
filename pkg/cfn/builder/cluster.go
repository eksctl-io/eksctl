package builder

import (
	"net"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	"github.com/weaveworks/eksctl/pkg/eks/api"
)

const (
	cfnOutputClusterCertificateAuthorityData = "CertificateAuthorityData"
	cfnOutputClusterEndpoint                 = "Endpoint"
	cfnOutputClusterARN                      = "ARN"
	cfnOutputClusterStackName                = "ClusterStackName"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs             *resourceSet
	spec           *api.ClusterConfig
	subnets        []*gfn.Value
	securityGroups []*gfn.Value
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(spec *api.ClusterConfig) *ClusterResourceSet {
	return &ClusterResourceSet{
		rs:   newResourceSet(),
		spec: spec,
	}
}

// AddAllResources adds all the information about the cluster to the resource set
func (c *ClusterResourceSet) AddAllResources() error {
	c.rs.template.Description = clusterTemplateDescription
	c.rs.template.Description += clusterTemplateDescriptionDefaultFeatures
	c.rs.template.Description += templateDescriptionSuffix

	_, globalCIDR, _ := net.ParseCIDR("192.168.0.0/16")

	subnets := map[string]*net.IPNet{}
	_, subnets[c.spec.AvailabilityZones[0]], _ = net.ParseCIDR("192.168.64.0/18")
	_, subnets[c.spec.AvailabilityZones[1]], _ = net.ParseCIDR("192.168.128.0/18")
	_, subnets[c.spec.AvailabilityZones[2]], _ = net.ParseCIDR("192.168.192.0/18")

	c.addResourcesForVPC(globalCIDR, subnets)
	c.addResourcesForIAM()
	c.addResourcesForControlPlane("1.10")

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

func (c *ClusterResourceSet) addResourcesForControlPlane(version string) {
	c.newResource("ControlPlane", &gfn.AWSEKSCluster{
		Name:    gfn.NewString(c.spec.ClusterName),
		RoleArn: gfn.MakeFnGetAttString("ServiceRole.Arn"),
		Version: gfn.NewString(version),
		ResourcesVpcConfig: &gfn.AWSEKSCluster_ResourcesVpcConfig{
			SubnetIds:        c.subnets,
			SecurityGroupIds: c.securityGroups,
		},
	})

	c.rs.newOutputFromAtt(cfnOutputClusterCertificateAuthorityData, "ControlPlane.CertificateAuthorityData", false)
	c.rs.newOutputFromAtt(cfnOutputClusterEndpoint, "ControlPlane.Endpoint", true)
	c.rs.newOutputFromAtt(cfnOutputClusterARN, "ControlPlane.Arn", true)
}

// GetAllOutputs collects all outputs of the cluster
func (c *ClusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return c.rs.GetAllOutputs(stack, c.spec)
}
