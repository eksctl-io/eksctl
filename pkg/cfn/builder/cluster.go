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

type clusterResourceSet struct {
	rs             *resourceSet
	spec           *api.ClusterConfig
	subnets        []string
	securityGroups []string
}

func NewClusterResourceSet(spec *api.ClusterConfig) *clusterResourceSet {
	return &clusterResourceSet{
		rs:   newResourceSet(),
		spec: spec,
	}
}

func (c *clusterResourceSet) AddAllResources() error {
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

	c.rs.newOutput(cfnOutputClusterStackName, refStackName, false)

	return nil
}

func (c *clusterResourceSet) RenderJSON() ([]byte, error) {
	return c.rs.renderJSON()
}

func (c *clusterResourceSet) newResource(name string, resource interface{}) string {
	return c.rs.newResource(name, resource)
}

func (c *clusterResourceSet) addResourcesForControlPlane(version string) {
	c.newResource("ControlPlane", &gfn.AWSEKSCluster{
		Name:    c.rs.newStringParameter(ParamClusterName, ""),
		RoleArn: gfn.GetAtt("ServiceRole", "Arn"),
		Version: version,
		ResourcesVpcConfig: &gfn.AWSEKSCluster_ResourcesVpcConfig{
			SubnetIds:        c.subnets,
			SecurityGroupIds: c.securityGroups,
		},
	})

	c.rs.newOutputFromAtt(cfnOutputClusterCertificateAuthorityData, "ControlPlane", "CertificateAuthorityData", false)
	c.rs.newOutputFromAtt(cfnOutputClusterEndpoint, "ControlPlane", "Endpoint", true)
	c.rs.newOutputFromAtt(cfnOutputClusterARN, "ControlPlane", "Arn", true)
}

func (c *clusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return c.rs.GetAllOutputs(stack, c.spec)
}
