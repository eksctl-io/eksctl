package cloudformation

import (
	"encoding/json"
)

// AWSGlueTable_TableInput AWS CloudFormation Resource (AWS::Glue::Table.TableInput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html
type AWSGlueTable_TableInput struct {

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-description
	Description *Value `json:"Description,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-name
	Name *Value `json:"Name,omitempty"`

	// Owner AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-owner
	Owner *Value `json:"Owner,omitempty"`

	// Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-parameters
	Parameters interface{} `json:"Parameters,omitempty"`

	// PartitionKeys AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-partitionkeys
	PartitionKeys []AWSGlueTable_Column `json:"PartitionKeys,omitempty"`

	// Retention AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-retention
	Retention *Value `json:"Retention,omitempty"`

	// StorageDescriptor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-storagedescriptor
	StorageDescriptor *AWSGlueTable_StorageDescriptor `json:"StorageDescriptor,omitempty"`

	// TableType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-tabletype
	TableType *Value `json:"TableType,omitempty"`

	// ViewExpandedText AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-viewexpandedtext
	ViewExpandedText *Value `json:"ViewExpandedText,omitempty"`

	// ViewOriginalText AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-table-tableinput.html#cfn-glue-table-tableinput-vieworiginaltext
	ViewOriginalText *Value `json:"ViewOriginalText,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueTable_TableInput) AWSCloudFormationType() string {
	return "AWS::Glue::Table.TableInput"
}

func (r *AWSGlueTable_TableInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
