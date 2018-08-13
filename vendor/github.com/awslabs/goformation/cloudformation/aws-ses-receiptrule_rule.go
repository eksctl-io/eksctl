package cloudformation

import (
	"encoding/json"
)

// AWSSESReceiptRule_Rule AWS CloudFormation Resource (AWS::SES::ReceiptRule.Rule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-rule.html
type AWSSESReceiptRule_Rule struct {

	// Actions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-rule.html#cfn-ses-receiptrule-rule-actions
	Actions []AWSSESReceiptRule_Action `json:"Actions,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-rule.html#cfn-ses-receiptrule-rule-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-rule.html#cfn-ses-receiptrule-rule-name
	Name *Value `json:"Name,omitempty"`

	// Recipients AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-rule.html#cfn-ses-receiptrule-rule-recipients
	Recipients []*Value `json:"Recipients,omitempty"`

	// ScanEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-rule.html#cfn-ses-receiptrule-rule-scanenabled
	ScanEnabled *Value `json:"ScanEnabled,omitempty"`

	// TlsPolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-rule.html#cfn-ses-receiptrule-rule-tlspolicy
	TlsPolicy *Value `json:"TlsPolicy,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_Rule) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.Rule"
}

func (r *AWSSESReceiptRule_Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
