package s3

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StorageLens_SelectionCriteria AWS CloudFormation Resource (AWS::S3::StorageLens.SelectionCriteria)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-selectioncriteria.html
type StorageLens_SelectionCriteria struct {

	// Delimiter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-selectioncriteria.html#cfn-s3-storagelens-selectioncriteria-delimiter
	Delimiter *types.Value `json:"Delimiter,omitempty"`

	// MaxDepth AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-selectioncriteria.html#cfn-s3-storagelens-selectioncriteria-maxdepth
	MaxDepth *types.Value `json:"MaxDepth,omitempty"`

	// MinStorageBytesPercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-selectioncriteria.html#cfn-s3-storagelens-selectioncriteria-minstoragebytespercentage
	MinStorageBytesPercentage *types.Value `json:"MinStorageBytesPercentage,omitempty"`

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
func (r *StorageLens_SelectionCriteria) AWSCloudFormationType() string {
	return "AWS::S3::StorageLens.SelectionCriteria"
}
