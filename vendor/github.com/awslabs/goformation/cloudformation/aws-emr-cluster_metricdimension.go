package cloudformation

import (
	"encoding/json"
)

// AWSEMRCluster_MetricDimension AWS CloudFormation Resource (AWS::EMR::Cluster.MetricDimension)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-metricdimension.html
type AWSEMRCluster_MetricDimension struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-metricdimension.html#cfn-elasticmapreduce-cluster-metricdimension-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-metricdimension.html#cfn-elasticmapreduce-cluster-metricdimension-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRCluster_MetricDimension) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.MetricDimension"
}

func (r *AWSEMRCluster_MetricDimension) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
