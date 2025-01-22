package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Index_ValueImportanceItems AWS CloudFormation Resource (AWS::Kendra::Index.ValueImportanceItems)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-valueimportanceitems.html
type Index_ValueImportanceItems struct {

	// ValueImportanceItems AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-valueimportanceitems.html#cfn-kendra-index-valueimportanceitems-valueimportanceitems
	ValueImportanceItems []Index_ValueImportanceItem `json:"ValueImportanceItems,omitempty"`

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
func (r *Index_ValueImportanceItems) AWSCloudFormationType() string {
	return "AWS::Kendra::Index.ValueImportanceItems"
}
