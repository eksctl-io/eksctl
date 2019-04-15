package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_Function AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.Function)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-function.html
type AWSGreengrassFunctionDefinition_Function struct {

	// FunctionArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-function.html#cfn-greengrass-functiondefinition-function-functionarn
	FunctionArn *Value `json:"FunctionArn,omitempty"`

	// FunctionConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-function.html#cfn-greengrass-functiondefinition-function-functionconfiguration
	FunctionConfiguration *AWSGreengrassFunctionDefinition_FunctionConfiguration `json:"FunctionConfiguration,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-function.html#cfn-greengrass-functiondefinition-function-id
	Id *Value `json:"Id,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_Function) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.Function"
}

func (r *AWSGreengrassFunctionDefinition_Function) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
