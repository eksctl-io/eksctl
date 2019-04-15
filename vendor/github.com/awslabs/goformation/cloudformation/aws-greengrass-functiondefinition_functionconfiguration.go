package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_FunctionConfiguration AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.FunctionConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html
type AWSGreengrassFunctionDefinition_FunctionConfiguration struct {

	// EncodingType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html#cfn-greengrass-functiondefinition-functionconfiguration-encodingtype
	EncodingType *Value `json:"EncodingType,omitempty"`

	// Environment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html#cfn-greengrass-functiondefinition-functionconfiguration-environment
	Environment *AWSGreengrassFunctionDefinition_Environment `json:"Environment,omitempty"`

	// ExecArgs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html#cfn-greengrass-functiondefinition-functionconfiguration-execargs
	ExecArgs *Value `json:"ExecArgs,omitempty"`

	// Executable AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html#cfn-greengrass-functiondefinition-functionconfiguration-executable
	Executable *Value `json:"Executable,omitempty"`

	// MemorySize AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html#cfn-greengrass-functiondefinition-functionconfiguration-memorysize
	MemorySize *Value `json:"MemorySize,omitempty"`

	// Pinned AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html#cfn-greengrass-functiondefinition-functionconfiguration-pinned
	Pinned *Value `json:"Pinned,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functionconfiguration.html#cfn-greengrass-functiondefinition-functionconfiguration-timeout
	Timeout *Value `json:"Timeout,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_FunctionConfiguration) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.FunctionConfiguration"
}

func (r *AWSGreengrassFunctionDefinition_FunctionConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
