package cloudformation

import (
	"encoding/json"
)

// AWSEMRCluster_StepConfig AWS CloudFormation Resource (AWS::EMR::Cluster.StepConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-stepconfig.html
type AWSEMRCluster_StepConfig struct {

	// ActionOnFailure AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-stepconfig.html#cfn-elasticmapreduce-cluster-stepconfig-actiononfailure
	ActionOnFailure *Value `json:"ActionOnFailure,omitempty"`

	// HadoopJarStep AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-stepconfig.html#cfn-elasticmapreduce-cluster-stepconfig-hadoopjarstep
	HadoopJarStep *AWSEMRCluster_HadoopJarStepConfig `json:"HadoopJarStep,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-stepconfig.html#cfn-elasticmapreduce-cluster-stepconfig-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRCluster_StepConfig) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.StepConfig"
}

func (r *AWSEMRCluster_StepConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
