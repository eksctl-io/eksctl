package cloudformation

import (
	"encoding/json"
)

// AWSEMRCluster_KeyValue AWS CloudFormation Resource (AWS::EMR::Cluster.KeyValue)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-keyvalue.html
type AWSEMRCluster_KeyValue struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-keyvalue.html#cfn-elasticmapreduce-cluster-keyvalue-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-keyvalue.html#cfn-elasticmapreduce-cluster-keyvalue-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRCluster_KeyValue) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.KeyValue"
}

func (r *AWSEMRCluster_KeyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
