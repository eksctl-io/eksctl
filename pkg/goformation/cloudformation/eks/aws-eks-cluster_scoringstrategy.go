package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

// Cluster_ScoringStrategy AWS CloudFormation Resource (AWS::EKS::Cluster.ScoringStrategy)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-scoringstrategy.html
type Cluster_ScoringStrategy struct {

	// Type AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-scoringstrategy.html#cfn-eks-cluster-scoringstrategy-type
	Type *types.Value `json:"Type,omitempty"`

	// Resources AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-scoringstrategy.html#cfn-eks-cluster-scoringstrategy-resources
	Resources []Cluster_ResourceWeight `json:"Resources,omitempty"`

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
func (r *Cluster_ScoringStrategy) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.ScoringStrategy"
}
