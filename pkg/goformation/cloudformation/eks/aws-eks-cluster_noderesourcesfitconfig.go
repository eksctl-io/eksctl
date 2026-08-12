package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// Cluster_NodeResourcesFitConfig AWS CloudFormation Resource (AWS::EKS::Cluster.NodeResourcesFitConfig)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-noderesourcesfitconfig.html
type Cluster_NodeResourcesFitConfig struct {

	// ScoringStrategy AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-noderesourcesfitconfig.html#cfn-eks-cluster-noderesourcesfitconfig-scoringstrategy
	ScoringStrategy *Cluster_ScoringStrategy `json:"ScoringStrategy,omitempty"`

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
func (r *Cluster_NodeResourcesFitConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.NodeResourcesFitConfig"
}
