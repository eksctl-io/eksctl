package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_DefaultConfig AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.DefaultConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-defaultconfig.html
type AWSGreengrassFunctionDefinition_DefaultConfig struct {

	// Execution AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-defaultconfig.html#cfn-greengrass-functiondefinition-defaultconfig-execution
	Execution *AWSGreengrassFunctionDefinition_Execution `json:"Execution,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_DefaultConfig) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.DefaultConfig"
}

func (r *AWSGreengrassFunctionDefinition_DefaultConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
