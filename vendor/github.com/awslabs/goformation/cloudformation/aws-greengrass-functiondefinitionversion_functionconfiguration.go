package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinitionVersion_FunctionConfiguration AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion.FunctionConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html
type AWSGreengrassFunctionDefinitionVersion_FunctionConfiguration struct {

	// EncodingType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html#cfn-greengrass-functiondefinitionversion-functionconfiguration-encodingtype
	EncodingType *Value `json:"EncodingType,omitempty"`

	// Environment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html#cfn-greengrass-functiondefinitionversion-functionconfiguration-environment
	Environment *AWSGreengrassFunctionDefinitionVersion_Environment `json:"Environment,omitempty"`

	// ExecArgs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html#cfn-greengrass-functiondefinitionversion-functionconfiguration-execargs
	ExecArgs *Value `json:"ExecArgs,omitempty"`

	// Executable AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html#cfn-greengrass-functiondefinitionversion-functionconfiguration-executable
	Executable *Value `json:"Executable,omitempty"`

	// MemorySize AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html#cfn-greengrass-functiondefinitionversion-functionconfiguration-memorysize
	MemorySize *Value `json:"MemorySize,omitempty"`

	// Pinned AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html#cfn-greengrass-functiondefinitionversion-functionconfiguration-pinned
	Pinned *Value `json:"Pinned,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-functionconfiguration.html#cfn-greengrass-functiondefinitionversion-functionconfiguration-timeout
	Timeout *Value `json:"Timeout,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion_FunctionConfiguration) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion.FunctionConfiguration"
}

func (r *AWSGreengrassFunctionDefinitionVersion_FunctionConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
