package cloudformation

import (
	"encoding/json"
)

// AWSGluePartition_Column AWS CloudFormation Resource (AWS::Glue::Partition.Column)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-column.html
type AWSGluePartition_Column struct {

	// Comment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-column.html#cfn-glue-partition-column-comment
	Comment *Value `json:"Comment,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-column.html#cfn-glue-partition-column-name
	Name *Value `json:"Name,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-column.html#cfn-glue-partition-column-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGluePartition_Column) AWSCloudFormationType() string {
	return "AWS::Glue::Partition.Column"
}

func (r *AWSGluePartition_Column) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
