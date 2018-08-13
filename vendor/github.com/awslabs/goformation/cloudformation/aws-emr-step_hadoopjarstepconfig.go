package cloudformation

import (
	"encoding/json"
)

// AWSEMRStep_HadoopJarStepConfig AWS CloudFormation Resource (AWS::EMR::Step.HadoopJarStepConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-step-hadoopjarstepconfig.html
type AWSEMRStep_HadoopJarStepConfig struct {

	// Args AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-step-hadoopjarstepconfig.html#cfn-elasticmapreduce-step-hadoopjarstepconfig-args
	Args []*Value `json:"Args,omitempty"`

	// Jar AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-step-hadoopjarstepconfig.html#cfn-elasticmapreduce-step-hadoopjarstepconfig-jar
	Jar *Value `json:"Jar,omitempty"`

	// MainClass AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-step-hadoopjarstepconfig.html#cfn-elasticmapreduce-step-hadoopjarstepconfig-mainclass
	MainClass *Value `json:"MainClass,omitempty"`

	// StepProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-step-hadoopjarstepconfig.html#cfn-elasticmapreduce-step-hadoopjarstepconfig-stepproperties
	StepProperties []AWSEMRStep_KeyValue `json:"StepProperties,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRStep_HadoopJarStepConfig) AWSCloudFormationType() string {
	return "AWS::EMR::Step.HadoopJarStepConfig"
}

func (r *AWSEMRStep_HadoopJarStepConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
