package s3

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StorageLens_PrefixLevelStorageMetrics AWS CloudFormation Resource (AWS::S3::StorageLens.PrefixLevelStorageMetrics)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-prefixlevelstoragemetrics.html
type StorageLens_PrefixLevelStorageMetrics struct {

	// IsEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-prefixlevelstoragemetrics.html#cfn-s3-storagelens-prefixlevelstoragemetrics-isenabled
	IsEnabled *types.Value `json:"IsEnabled,omitempty"`

	// SelectionCriteria AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-prefixlevelstoragemetrics.html#cfn-s3-storagelens-prefixlevelstoragemetrics-selectioncriteria
	SelectionCriteria *StorageLens_SelectionCriteria `json:"SelectionCriteria,omitempty"`

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
func (r *StorageLens_PrefixLevelStorageMetrics) AWSCloudFormationType() string {
	return "AWS::S3::StorageLens.PrefixLevelStorageMetrics"
}
