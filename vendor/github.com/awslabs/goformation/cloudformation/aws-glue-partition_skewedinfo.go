package cloudformation

import (
	"encoding/json"
)

// AWSGluePartition_SkewedInfo AWS CloudFormation Resource (AWS::Glue::Partition.SkewedInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-skewedinfo.html
type AWSGluePartition_SkewedInfo struct {

	// SkewedColumnNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-skewedinfo.html#cfn-glue-partition-skewedinfo-skewedcolumnnames
	SkewedColumnNames []*Value `json:"SkewedColumnNames,omitempty"`

	// SkewedColumnValueLocationMaps AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-skewedinfo.html#cfn-glue-partition-skewedinfo-skewedcolumnvaluelocationmaps
	SkewedColumnValueLocationMaps interface{} `json:"SkewedColumnValueLocationMaps,omitempty"`

	// SkewedColumnValues AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-skewedinfo.html#cfn-glue-partition-skewedinfo-skewedcolumnvalues
	SkewedColumnValues []*Value `json:"SkewedColumnValues,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGluePartition_SkewedInfo) AWSCloudFormationType() string {
	return "AWS::Glue::Partition.SkewedInfo"
}

func (r *AWSGluePartition_SkewedInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
