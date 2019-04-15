package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinitionVersion_Execution AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion.Execution)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-execution.html
type AWSGreengrassFunctionDefinitionVersion_Execution struct {

	// IsolationMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-execution.html#cfn-greengrass-functiondefinitionversion-execution-isolationmode
	IsolationMode *Value `json:"IsolationMode,omitempty"`

	// RunAs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-execution.html#cfn-greengrass-functiondefinitionversion-execution-runas
	RunAs *AWSGreengrassFunctionDefinitionVersion_RunAs `json:"RunAs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion_Execution) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion.Execution"
}

func (r *AWSGreengrassFunctionDefinitionVersion_Execution) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
