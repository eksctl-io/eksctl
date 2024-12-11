package athena

import (
	"goformation/v4/cloudformation/cloudformation"
	"goformation/v4/cloudformation/policies"
)

// DataCatalog_Tags AWS CloudFormation Resource (AWS::Athena::DataCatalog.Tags)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-datacatalog-tags.html
type DataCatalog_Tags struct {

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-athena-datacatalog-tags.html#cfn-athena-datacatalog-tags-tags
	Tags []cloudformation.Tag `json:"Tags,omitempty"`

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
func (r *DataCatalog_Tags) AWSCloudFormationType() string {
	return "AWS::Athena::DataCatalog.Tags"
}
