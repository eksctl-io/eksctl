package builder

import (
	"encoding/base64"
	"fmt"

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
	provider             api.ClusterProvider
	supportsManagedNodes bool
	vpc                  *gfnt.Value
	subnets              map[api.SubnetTopology][]*gfnt.Value
	securityGroups       []*gfnt.Value
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig, supportsManagedNodes bool, existingStack *gjson.Result) *ClusterResourceSet {
	if existingStack != nil {
		unsetExistingResources(existingStack, spec)
	}
	return &ClusterResourceSet{
		rs:                   newResourceSet(),
		spec:                 spec,
		provider:             provider,
		supportsManagedNodes: supportsManagedNodes,
	}
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

// HasManagedNodesSG reports whether the stack has the security group required for communication between
// managed and unmanaged nodegroups
func HasManagedNodesSG(stackResources *gjson.Result) bool {
	return stackResources.Get(cfnIngressClusterToNodeSGResource).Exists()
}

func (c *ClusterResourceSet) newResource(name string, resource gfn.Resource) *gfnt.Value {
	return c.rs.newResource(name, resource)
}

func (c *ClusterResourceSet) addResourcesForControlPlane() {
	clusterVPC := &gfneks.Cluster_ResourcesVpcConfig{
		SecurityGroupIds: gfnt.NewSlice(c.securityGroups...),
	}
	var subnetIds []*gfnt.Value
	for topology := range c.subnets {
		subnetIds = append(subnetIds, c.subnets[topology]...)
	}
	clusterVPC.SubnetIds = gfnt.NewSlice(subnetIds...)

	serviceRoleARN := gfnt.MakeFnGetAttString("ServiceRole", "Arn")
	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRoleARN) {
		serviceRoleARN = gfnt.NewString(*c.spec.IAM.ServiceRoleARN)
	}

	var encryptionConfigs []gfneks.Cluster_EncryptionConfig
	if c.spec.SecretsEncryption != nil && c.spec.SecretsEncryption.KeyARN != nil {
		encryptionConfigs = []gfneks.Cluster_EncryptionConfig{
			{
				Resources: gfnt.NewSlice(gfnt.NewString("secrets")),
				Provider: &gfneks.Cluster_Provider{
					KeyArn: gfnt.NewString(*c.spec.SecretsEncryption.KeyARN),
				},
			},
		}
	}

	c.newResource("ControlPlane", &gfneks.Cluster{
		Name:               gfnt.NewString(c.spec.Metadata.Name),
		RoleArn:            serviceRoleARN,
		Version:            gfnt.NewString(c.spec.Metadata.Version),
		ResourcesVpcConfig: clusterVPC,
		EncryptionConfig:   encryptionConfigs,
	})

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
	_ = AddResourcesForFargate(c.rs, c.spec)
}

// GetAllOutputs collects all outputs of the cluster
func (c *ClusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return c.rs.GetAllOutputs(stack)
}
