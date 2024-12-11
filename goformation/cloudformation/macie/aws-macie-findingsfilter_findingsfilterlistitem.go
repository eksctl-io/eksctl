package macie

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// FindingsFilter_FindingsFilterListItem AWS CloudFormation Resource (AWS::Macie::FindingsFilter.FindingsFilterListItem)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-macie-findingsfilter-findingsfilterlistitem.html
type FindingsFilter_FindingsFilterListItem struct {

	// Id AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-macie-findingsfilter-findingsfilterlistitem.html#cfn-macie-findingsfilter-findingsfilterlistitem-id
	Id *types.Value `json:"Id,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-macie-findingsfilter-findingsfilterlistitem.html#cfn-macie-findingsfilter-findingsfilterlistitem-name
	Name *types.Value `json:"Name,omitempty"`

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
func (r *FindingsFilter_FindingsFilterListItem) AWSCloudFormationType() string {
	return "AWS::Macie::FindingsFilter.FindingsFilterListItem"
}
