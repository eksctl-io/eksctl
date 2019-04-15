package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinitionVersion_Function AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion.Function)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-function.html
type AWSGreengrassFunctionDefinitionVersion_Function struct {

	// FunctionArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-function.html#cfn-greengrass-functiondefinitionversion-function-functionarn
	FunctionArn *Value `json:"FunctionArn,omitempty"`

	// FunctionConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-function.html#cfn-greengrass-functiondefinitionversion-function-functionconfiguration
	FunctionConfiguration *AWSGreengrassFunctionDefinitionVersion_FunctionConfiguration `json:"FunctionConfiguration,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-function.html#cfn-greengrass-functiondefinitionversion-function-id
	Id *Value `json:"Id,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion_Function) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion.Function"
}

func (r *AWSGreengrassFunctionDefinitionVersion_Function) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
