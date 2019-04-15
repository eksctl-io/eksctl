package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_RunAs AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.RunAs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-runas.html
type AWSGreengrassFunctionDefinition_RunAs struct {

	// Gid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-runas.html#cfn-greengrass-functiondefinition-runas-gid
	Gid *Value `json:"Gid,omitempty"`

	// Uid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-runas.html#cfn-greengrass-functiondefinition-runas-uid
	Uid *Value `json:"Uid,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_RunAs) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.RunAs"
}

func (r *AWSGreengrassFunctionDefinition_RunAs) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
