package eks

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Nodegroup_UpdateConfig AWS CloudFormation Resource (AWS::EKS::Nodegroup.UpdateConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-updateconfig.html
type Nodegroup_UpdateConfig struct {

	// MaxUnavailable AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-updateconfig.html#cfn-eks-nodegroup-updateconfig-maxunavailable
	MaxUnavailable *types.Value `json:"MaxUnavailable,omitempty"`

	// MaxUnavailablePercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-updateconfig.html#cfn-eks-nodegroup-updateconfig-maxunavailablepercentage
	MaxUnavailablePercentage *types.Value `json:"MaxUnavailablePercentage,omitempty"`

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
func (r *Nodegroup_UpdateConfig) AWSCloudFormationType() string {
	return "AWS::EKS::Nodegroup.UpdateConfig"
}
