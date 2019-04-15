package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinitionVersion_DefaultConfig AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion.DefaultConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-defaultconfig.html
type AWSGreengrassFunctionDefinitionVersion_DefaultConfig struct {

	// Execution AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-defaultconfig.html#cfn-greengrass-functiondefinitionversion-defaultconfig-execution
	Execution *AWSGreengrassFunctionDefinitionVersion_Execution `json:"Execution,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion_DefaultConfig) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion.DefaultConfig"
}

func (r *AWSGreengrassFunctionDefinitionVersion_DefaultConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
