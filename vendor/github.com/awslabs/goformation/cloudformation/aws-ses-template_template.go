package cloudformation

import (
	"encoding/json"
)

// AWSSESTemplate_Template AWS CloudFormation Resource (AWS::SES::Template.Template)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-template-template.html
type AWSSESTemplate_Template struct {

	// HtmlPart AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-template-template.html#cfn-ses-template-template-htmlpart
	HtmlPart *Value `json:"HtmlPart,omitempty"`

	// SubjectPart AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-template-template.html#cfn-ses-template-template-subjectpart
	SubjectPart *Value `json:"SubjectPart,omitempty"`

	// TemplateName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-template-template.html#cfn-ses-template-template-templatename
	TemplateName *Value `json:"TemplateName,omitempty"`

	// TextPart AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-template-template.html#cfn-ses-template-template-textpart
	TextPart *Value `json:"TextPart,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESTemplate_Template) AWSCloudFormationType() string {
	return "AWS::SES::Template.Template"
}

func (r *AWSSESTemplate_Template) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
