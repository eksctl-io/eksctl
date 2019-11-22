package builder

import (
	"encoding/base64"
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs                   *resourceSet
	spec                 *api.ClusterConfig
	provider             api.ClusterProvider
	supportsManagedNodes bool
	vpc                  *gfn.Value
	subnets              map[api.SubnetTopology][]*gfn.Value
	securityGroups       []*gfn.Value
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig, supportsManagedNodes bool) *ClusterResourceSet {
	return &ClusterResourceSet{
		rs:                   newResourceSet(),
		spec:                 spec,
		provider:             provider,
		supportsManagedNodes: supportsManagedNodes,
	}
}

// AddAllResources adds all the information about the cluster to the resource set
func (c *ClusterResourceSet) AddAllResources() error {
	dedicatedVPC := c.spec.VPC.ID == ""

	if err := c.spec.HasSufficientSubnets(); err != nil {
		return err
	}

	if dedicatedVPC {
		if err := c.addResourcesForVPC(); err != nil {
			return errors.Wrap(err, "error adding VPC resources")
		}
	} else {
		c.importResourcesForVPC()
	}
	c.addOutputsForVPC()

	c.addResourcesForSecurityGroups()
	c.addResourcesForIAM()
	c.addResourcesForControlPlane()

	c.rs.defineOutput(outputs.ClusterStackName, gfn.RefStackName, false, func(v string) error {
		if c.spec.Status == nil {
			c.spec.Status = &api.ClusterStatus{}
		}
		c.spec.Status.StackName = v
		return nil
	})

	c.rs.template.Description = fmt.Sprintf(
		"%s (dedicated VPC: %v, dedicated IAM: %v) %s",
		clusterTemplateDescription,
		dedicatedVPC, c.rs.withIAM,
		templateDescriptionSuffix)

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
	for topology := range c.subnets {
		clusterVPC.SubnetIds = append(clusterVPC.SubnetIds, c.subnets[topology]...)
	}

	serviceRoleARN := gfn.MakeFnGetAttString("ServiceRole.Arn")
	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRoleARN) {
		serviceRoleARN = gfn.NewString(*c.spec.IAM.ServiceRoleARN)
	}

	c.newResource("ControlPlane", &gfn.AWSEKSCluster{
		Name:               gfn.NewString(c.spec.Metadata.Name),
		RoleArn:            serviceRoleARN,
		Version:            gfn.NewString(c.spec.Metadata.Version),
		ResourcesVpcConfig: clusterVPC,
	})

	if c.spec.Status == nil {
		c.spec.Status = &api.ClusterStatus{}
	}

	c.rs.defineOutputFromAtt(outputs.ClusterCertificateAuthorityData, "ControlPlane.CertificateAuthorityData", false, func(v string) error {
		caData, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return errors.Wrap(err, "decoding certificate authority data")
		}
		c.spec.Status.CertificateAuthorityData = caData
		return nil
	})
	c.rs.defineOutputFromAtt(outputs.ClusterEndpoint, "ControlPlane.Endpoint", true, func(v string) error {
		c.spec.Status.Endpoint = v
		return nil
	})
	c.rs.defineOutputFromAtt(outputs.ClusterARN, "ControlPlane.Arn", true, func(v string) error {
		c.spec.Status.ARN = v
		return nil
	})

	if c.supportsManagedNodes {
		// This exports the cluster security group ID that EKS creates by default. To enable communication between both
		// managed and unmanaged nodegroups, they must share a security group.
		// EKS attaches this to Managed Nodegroups by default, but we need to add this for unmanaged nodegroups.
		// This exported value is imported in the CloudFormation resource for unmanaged nodegroups
		c.rs.defineOutputFromAtt(outputs.ClusterDefaultSecurityGroup, "ControlPlane.ClusterSecurityGroupId",
			true, func(s string) error {
				return nil
			})
	}
}

// GetAllOutputs collects all outputs of the cluster
func (c *ClusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return c.rs.GetAllOutputs(stack)
}
