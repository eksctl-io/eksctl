package cloudformation

import (
	"encoding/json"
)

// AWSGluePartition_SerdeInfo AWS CloudFormation Resource (AWS::Glue::Partition.SerdeInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-serdeinfo.html
type AWSGluePartition_SerdeInfo struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-serdeinfo.html#cfn-glue-partition-serdeinfo-name
	Name *Value `json:"Name,omitempty"`

	// Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-serdeinfo.html#cfn-glue-partition-serdeinfo-parameters
	Parameters interface{} `json:"Parameters,omitempty"`

	// SerializationLibrary AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-partition-serdeinfo.html#cfn-glue-partition-serdeinfo-serializationlibrary
	SerializationLibrary *Value `json:"SerializationLibrary,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGluePartition_SerdeInfo) AWSCloudFormationType() string {
	return "AWS::Glue::Partition.SerdeInfo"
}

func (r *AWSGluePartition_SerdeInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
