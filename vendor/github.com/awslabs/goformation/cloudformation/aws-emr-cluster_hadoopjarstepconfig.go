package cloudformation

import (
	"encoding/json"
)

// AWSEMRCluster_HadoopJarStepConfig AWS CloudFormation Resource (AWS::EMR::Cluster.HadoopJarStepConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-hadoopjarstepconfig.html
type AWSEMRCluster_HadoopJarStepConfig struct {

	// Args AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-hadoopjarstepconfig.html#cfn-elasticmapreduce-cluster-hadoopjarstepconfig-args
	Args []*Value `json:"Args,omitempty"`

	// Jar AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-hadoopjarstepconfig.html#cfn-elasticmapreduce-cluster-hadoopjarstepconfig-jar
	Jar *Value `json:"Jar,omitempty"`

	// MainClass AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-hadoopjarstepconfig.html#cfn-elasticmapreduce-cluster-hadoopjarstepconfig-mainclass
	MainClass *Value `json:"MainClass,omitempty"`

	// StepProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-hadoopjarstepconfig.html#cfn-elasticmapreduce-cluster-hadoopjarstepconfig-stepproperties
	StepProperties []AWSEMRCluster_KeyValue `json:"StepProperties,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRCluster_HadoopJarStepConfig) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.HadoopJarStepConfig"
}

func (r *AWSEMRCluster_HadoopJarStepConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
