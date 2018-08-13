package cloudformation

import (
	"encoding/json"
)

// AWSEMRInstanceFleetConfig_Configuration AWS CloudFormation Resource (AWS::EMR::InstanceFleetConfig.Configuration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-configuration.html
type AWSEMRInstanceFleetConfig_Configuration struct {

	// Classification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-configuration.html#cfn-elasticmapreduce-instancefleetconfig-configuration-classification
	Classification *Value `json:"Classification,omitempty"`

	// ConfigurationProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-configuration.html#cfn-elasticmapreduce-instancefleetconfig-configuration-configurationproperties
	ConfigurationProperties map[string]*Value `json:"ConfigurationProperties,omitempty"`

	// Configurations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-configuration.html#cfn-elasticmapreduce-instancefleetconfig-configuration-configurations
	Configurations []AWSEMRInstanceFleetConfig_Configuration `json:"Configurations,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRInstanceFleetConfig_Configuration) AWSCloudFormationType() string {
	return "AWS::EMR::InstanceFleetConfig.Configuration"
}

func (r *AWSEMRInstanceFleetConfig_Configuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
