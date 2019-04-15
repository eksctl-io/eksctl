package cloudformation

import (
	"encoding/json"
)

// AWSIoT1ClickProject_DeviceTemplate AWS CloudFormation Resource (AWS::IoT1Click::Project.DeviceTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot1click-project-devicetemplate.html
type AWSIoT1ClickProject_DeviceTemplate struct {

	// CallbackOverrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot1click-project-devicetemplate.html#cfn-iot1click-project-devicetemplate-callbackoverrides
	CallbackOverrides interface{} `json:"CallbackOverrides,omitempty"`

	// DeviceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot1click-project-devicetemplate.html#cfn-iot1click-project-devicetemplate-devicetype
	DeviceType *Value `json:"DeviceType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoT1ClickProject_DeviceTemplate) AWSCloudFormationType() string {
	return "AWS::IoT1Click::Project.DeviceTemplate"
}

func (r *AWSIoT1ClickProject_DeviceTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
