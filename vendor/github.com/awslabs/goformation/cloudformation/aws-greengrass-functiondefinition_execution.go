package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_Execution AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.Execution)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-execution.html
type AWSGreengrassFunctionDefinition_Execution struct {

	// IsolationMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-execution.html#cfn-greengrass-functiondefinition-execution-isolationmode
	IsolationMode *Value `json:"IsolationMode,omitempty"`

	// RunAs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-execution.html#cfn-greengrass-functiondefinition-execution-runas
	RunAs *AWSGreengrassFunctionDefinition_RunAs `json:"RunAs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_Execution) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.Execution"
}

func (r *AWSGreengrassFunctionDefinition_Execution) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
