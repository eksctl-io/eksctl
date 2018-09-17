package cloudformation

import (
	"encoding/json"
)

// AWSCodePipelineWebhook_WebhookAuthConfiguration AWS CloudFormation Resource (AWS::CodePipeline::Webhook.WebhookAuthConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-webhook-webhookauthconfiguration.html
type AWSCodePipelineWebhook_WebhookAuthConfiguration struct {

	// AllowedIPRange AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-webhook-webhookauthconfiguration.html#cfn-codepipeline-webhook-webhookauthconfiguration-allowediprange
	AllowedIPRange *Value `json:"AllowedIPRange,omitempty"`

	// SecretToken AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-webhook-webhookauthconfiguration.html#cfn-codepipeline-webhook-webhookauthconfiguration-secrettoken
	SecretToken *Value `json:"SecretToken,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodePipelineWebhook_WebhookAuthConfiguration) AWSCloudFormationType() string {
	return "AWS::CodePipeline::Webhook.WebhookAuthConfiguration"
}

func (r *AWSCodePipelineWebhook_WebhookAuthConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
