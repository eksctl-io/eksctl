package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinitionVersion_Environment AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion.Environment)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-environment.html
type AWSGreengrassFunctionDefinitionVersion_Environment struct {

	// AccessSysfs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-environment.html#cfn-greengrass-functiondefinitionversion-environment-accesssysfs
	AccessSysfs *Value `json:"AccessSysfs,omitempty"`

	// Execution AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-environment.html#cfn-greengrass-functiondefinitionversion-environment-execution
	Execution *AWSGreengrassFunctionDefinitionVersion_Execution `json:"Execution,omitempty"`

	// ResourceAccessPolicies AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-environment.html#cfn-greengrass-functiondefinitionversion-environment-resourceaccesspolicies
	ResourceAccessPolicies []AWSGreengrassFunctionDefinitionVersion_ResourceAccessPolicy `json:"ResourceAccessPolicies,omitempty"`

	// Variables AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-environment.html#cfn-greengrass-functiondefinitionversion-environment-variables
	Variables interface{} `json:"Variables,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion_Environment) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion.Environment"
}

func (r *AWSGreengrassFunctionDefinitionVersion_Environment) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
