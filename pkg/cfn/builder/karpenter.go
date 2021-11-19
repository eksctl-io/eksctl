package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// KarpenterResourceSet stores the resource information of the Karpenter stack
type KarpenterResourceSet struct {
	rs          *resourceSet
	clusterSpec *api.ClusterConfig
	iamAPI      iamiface.IAMAPI
	// instanceProfileARN *gfnt.Value
}

// NewKarpenterResourceSet returns a resource set for a Karpenter embedded in a cluster config
func NewKarpenterResourceSet(iamAPI iamiface.IAMAPI, spec *api.ClusterConfig) *KarpenterResourceSet {
	return &KarpenterResourceSet{
		rs:          newResourceSet(),
		clusterSpec: spec,
		iamAPI:      iamAPI,
	}
}

// AddAllResources adds all the information about the nodegroup to the resource set
func (k *KarpenterResourceSet) AddAllResources() error {
	k.rs.template.Description = fmt.Sprintf("Karpenter Stack %s", templateDescriptionSuffix)
	return k.addResourcesForKarpenter()
}

// RenderJSON returns the rendered JSON
func (k *KarpenterResourceSet) RenderJSON() ([]byte, error) {
	return k.rs.renderJSON()
}

// Template returns the CloudFormation template
func (k *KarpenterResourceSet) Template() gfn.Template {
	return *k.rs.template
}

func (k *KarpenterResourceSet) newResource(name string, resource gfn.Resource) *gfnt.Value {
	return k.rs.newResource(name, resource)
}

func (k *KarpenterResourceSet) addResourcesForKarpenter() error {
	cfr := karpenterResource(k.clusterSpec)
	k.newResource("Karpenter", cfr)
	return nil
}

func karpenterResource(cfg *api.ClusterConfig) *awsCloudFormationResource {
	return &awsCloudFormationResource{}
}

// WithIAM implements the ResourceSet interface
func (k *KarpenterResourceSet) WithIAM() bool {
	// eksctl does not support passing pre-created IAM instance roles to Managed Nodes,
	// so the IAM capability is always required
	return true
}

// WithNamedIAM implements the ResourceSet interface
func (k *KarpenterResourceSet) WithNamedIAM() bool {
	//return k.nodeGroup.IAM.InstanceRoleName != ""
	return false
}

// GetAllOutputs collects all outputs of the nodegroup
func (k *KarpenterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return k.rs.GetAllOutputs(stack)
}
