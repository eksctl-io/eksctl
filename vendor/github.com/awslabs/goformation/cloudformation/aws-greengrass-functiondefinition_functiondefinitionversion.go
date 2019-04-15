package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_FunctionDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.FunctionDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functiondefinitionversion.html
type AWSGreengrassFunctionDefinition_FunctionDefinitionVersion struct {

	// DefaultConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functiondefinitionversion.html#cfn-greengrass-functiondefinition-functiondefinitionversion-defaultconfig
	DefaultConfig *AWSGreengrassFunctionDefinition_DefaultConfig `json:"DefaultConfig,omitempty"`

	// Functions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-functiondefinitionversion.html#cfn-greengrass-functiondefinition-functiondefinitionversion-functions
	Functions []AWSGreengrassFunctionDefinition_Function `json:"Functions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_FunctionDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.FunctionDefinitionVersion"
}

func (r *AWSGreengrassFunctionDefinition_FunctionDefinitionVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
