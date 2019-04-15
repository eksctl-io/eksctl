package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_Environment AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.Environment)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-environment.html
type AWSGreengrassFunctionDefinition_Environment struct {

	// AccessSysfs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-environment.html#cfn-greengrass-functiondefinition-environment-accesssysfs
	AccessSysfs *Value `json:"AccessSysfs,omitempty"`

	// Execution AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-environment.html#cfn-greengrass-functiondefinition-environment-execution
	Execution *AWSGreengrassFunctionDefinition_Execution `json:"Execution,omitempty"`

	// ResourceAccessPolicies AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-environment.html#cfn-greengrass-functiondefinition-environment-resourceaccesspolicies
	ResourceAccessPolicies []AWSGreengrassFunctionDefinition_ResourceAccessPolicy `json:"ResourceAccessPolicies,omitempty"`

	// Variables AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-environment.html#cfn-greengrass-functiondefinition-environment-variables
	Variables interface{} `json:"Variables,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_Environment) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.Environment"
}

func (r *AWSGreengrassFunctionDefinition_Environment) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
