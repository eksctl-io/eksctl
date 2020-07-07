package builder

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	gfn "github.com/weaveworks/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs                   *resourceSet
	spec                 *api.ClusterConfig
	provider             api.ClusterProvider
	supportsManagedNodes bool
	vpcResourceSet       *VPCResourceSet
	securityGroups       []*gfn.Value
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig, supportsManagedNodes bool, existingStack *gjson.Result) *ClusterResourceSet {
	if existingStack != nil {
		unsetExistingResources(existingStack, spec)
	}
	rs := newResourceSet()
	return &ClusterResourceSet{
		rs:                   rs,
		spec:                 spec,
		provider:             provider,
		supportsManagedNodes: supportsManagedNodes,
		vpcResourceSet:       NewVPCResourceSet(rs, spec, provider),
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

	if err := c.spec.HasSufficientSubnets(); err != nil {
		return err
	}

	vpcResource, err := c.vpcResourceSet.AddResources()
	if err != nil {
		return errors.Wrap(err, "error adding VPC resources")
	}

	c.vpcResourceSet.AddOutputs()
	clusterSG := c.addResourcesForSecurityGroups(vpcResource)

	if privateCluster := c.spec.PrivateCluster; privateCluster.Enabled {
		vpcEndpointResourceSet := NewVPCEndpointResourceSet(c.provider, c.rs, c.spec, vpcResource.VPC, vpcResource.SubnetDetails.Private, clusterSG.ClusterSharedNode)

		if err := vpcEndpointResourceSet.AddResources(); err != nil {
			return errors.Wrap(err, "error adding resources for VPC endpoints")
		}
	}

	c.addResourcesForIAM()
	c.addResourcesForControlPlane(vpcResource.SubnetDetails)

	if len(c.spec.FargateProfiles) > 0 {
		c.addResourcesForFargate()
	}

	c.rs.defineOutput(outputs.ClusterStackName, gfn.RefStackName, false, func(v string) error {
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

// HasManagedNodesSG reports whether the stack has the security group required for communication between
// managed and unmanaged nodegroups
func HasManagedNodesSG(stackResources *gjson.Result) bool {
	return stackResources.Get(cfnIngressClusterToNodeSGResource).Exists()
}

func (c *ClusterResourceSet) newResource(name string, resource interface{}) *gfn.Value {
	return c.rs.newResource(name, resource)
}

// TODO use goformation after support is out
type awsEKSClusterKMS struct {
	*awsEKSCluster   `json:",inline"`
	EncryptionConfig []*encryptionConfig `json:"EncryptionConfig,omitempty"`
}

func (e *awsEKSClusterKMS) MarshalJSON() ([]byte, error) {
	type Properties awsEKSClusterKMS
	val, err := json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       "AWS::EKS::Cluster",
		Properties: Properties(*e),
	})
	return val, err
}

type encryptionProvider struct {
	KeyArn string `json:"KeyArn"`
}

type encryptionConfig struct {
	Provider  *encryptionProvider `json:"Provider"`
	Resources []string            `json:"Resources"`
}

type awsEKSCluster gfn.AWSEKSCluster

func (c *ClusterResourceSet) addResourcesForControlPlane(subnetDetails *subnetDetails) {
	clusterVPC := &gfn.AWSEKSCluster_ResourcesVpcConfig{
		SecurityGroupIds: c.securityGroups,
	}

	clusterVPC.SubnetIds = append(clusterVPC.SubnetIds, subnetDetails.PublicSubnetRefs()...)
	clusterVPC.SubnetIds = append(clusterVPC.SubnetIds, subnetDetails.PrivateSubnetRefs()...)

	serviceRoleARN := gfn.MakeFnGetAttString("ServiceRole.Arn")
	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRoleARN) {
		serviceRoleARN = gfn.NewString(*c.spec.IAM.ServiceRoleARN)
	}

	var encryptionConfigs []*encryptionConfig
	if c.spec.SecretsEncryption != nil && c.spec.SecretsEncryption.KeyARN != nil {
		encryptionConfigs = []*encryptionConfig{
			{
				Resources: []string{"secrets"},
				Provider: &encryptionProvider{
					KeyArn: *c.spec.SecretsEncryption.KeyARN,
				},
			},
		}
	}

	c.newResource("ControlPlane", &awsEKSClusterKMS{
		awsEKSCluster: &awsEKSCluster{
			Name:               gfn.NewString(c.spec.Metadata.Name),
			RoleArn:            serviceRoleARN,
			Version:            gfn.NewString(c.spec.Metadata.Version),
			ResourcesVpcConfig: clusterVPC,
		},
		EncryptionConfig: encryptionConfigs,
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

func (c *ClusterResourceSet) addResourcesForFargate() {
	_ = AddResourcesForFargate(c.rs, c.spec)
}

// GetAllOutputs collects all outputs of the cluster
func (c *ClusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return c.rs.GetAllOutputs(stack)
}
