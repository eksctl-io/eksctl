package cloudformation

import (
	"encoding/json"
)

// AWSEMRCluster_InstanceTypeConfig AWS CloudFormation Resource (AWS::EMR::Cluster.InstanceTypeConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancetypeconfig.html
type AWSEMRCluster_InstanceTypeConfig struct {

	// BidPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancetypeconfig.html#cfn-elasticmapreduce-cluster-instancetypeconfig-bidprice
	BidPrice *Value `json:"BidPrice,omitempty"`

	// BidPriceAsPercentageOfOnDemandPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancetypeconfig.html#cfn-elasticmapreduce-cluster-instancetypeconfig-bidpriceaspercentageofondemandprice
	BidPriceAsPercentageOfOnDemandPrice *Value `json:"BidPriceAsPercentageOfOnDemandPrice,omitempty"`

	// Configurations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancetypeconfig.html#cfn-elasticmapreduce-cluster-instancetypeconfig-configurations
	Configurations []AWSEMRCluster_Configuration `json:"Configurations,omitempty"`

	// EbsConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancetypeconfig.html#cfn-elasticmapreduce-cluster-instancetypeconfig-ebsconfiguration
	EbsConfiguration *AWSEMRCluster_EbsConfiguration `json:"EbsConfiguration,omitempty"`

	// InstanceType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancetypeconfig.html#cfn-elasticmapreduce-cluster-instancetypeconfig-instancetype
	InstanceType *Value `json:"InstanceType,omitempty"`

	// WeightedCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancetypeconfig.html#cfn-elasticmapreduce-cluster-instancetypeconfig-weightedcapacity
	WeightedCapacity *Value `json:"WeightedCapacity,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRCluster_InstanceTypeConfig) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.InstanceTypeConfig"
}

func (r *AWSEMRCluster_InstanceTypeConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
