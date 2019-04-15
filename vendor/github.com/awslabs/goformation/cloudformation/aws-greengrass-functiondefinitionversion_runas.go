package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinitionVersion_RunAs AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion.RunAs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-runas.html
type AWSGreengrassFunctionDefinitionVersion_RunAs struct {

	// Gid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-runas.html#cfn-greengrass-functiondefinitionversion-runas-gid
	Gid *Value `json:"Gid,omitempty"`

	// Uid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinitionversion-runas.html#cfn-greengrass-functiondefinitionversion-runas-uid
	Uid *Value `json:"Uid,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion_RunAs) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion.RunAs"
}

func (r *AWSGreengrassFunctionDefinitionVersion_RunAs) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
