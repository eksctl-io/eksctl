package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

// Cluster_ResourceWeight AWS CloudFormation Resource (AWS::EKS::Cluster.ResourceWeight)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-resourceweight.html
type Cluster_ResourceWeight struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-resourceweight.html#cfn-eks-cluster-resourceweight-name
	Name *types.Value `json:"Name,omitempty"`

	// Weight AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-resourceweight.html#cfn-eks-cluster-resourceweight-weight
	Weight *types.Value `json:"Weight,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Cluster_ResourceWeight) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.ResourceWeight"
}
