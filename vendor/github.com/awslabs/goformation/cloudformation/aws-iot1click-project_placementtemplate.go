package cloudformation

import (
	"encoding/json"
)

// AWSIoT1ClickProject_PlacementTemplate AWS CloudFormation Resource (AWS::IoT1Click::Project.PlacementTemplate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot1click-project-placementtemplate.html
type AWSIoT1ClickProject_PlacementTemplate struct {

	// DefaultAttributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot1click-project-placementtemplate.html#cfn-iot1click-project-placementtemplate-defaultattributes
	DefaultAttributes interface{} `json:"DefaultAttributes,omitempty"`

	// DeviceTemplates AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot1click-project-placementtemplate.html#cfn-iot1click-project-placementtemplate-devicetemplates
	DeviceTemplates interface{} `json:"DeviceTemplates,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoT1ClickProject_PlacementTemplate) AWSCloudFormationType() string {
	return "AWS::IoT1Click::Project.PlacementTemplate"
}

func (r *AWSIoT1ClickProject_PlacementTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
