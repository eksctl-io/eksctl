package cloudformation

// AWSCodePipelineWebhook_WebhookFilterRule AWS CloudFormation Resource (AWS::CodePipeline::Webhook.WebhookFilterRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-webhook-webhookfilterrule.html
type AWSCodePipelineWebhook_WebhookFilterRule struct {

	// JsonPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-webhook-webhookfilterrule.html#cfn-codepipeline-webhook-webhookfilterrule-jsonpath
	JsonPath *StringIntrinsic `json:"JsonPath,omitempty"`

	// MatchEquals AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-webhook-webhookfilterrule.html#cfn-codepipeline-webhook-webhookfilterrule-matchequals
	MatchEquals *StringIntrinsic `json:"MatchEquals,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodePipelineWebhook_WebhookFilterRule) AWSCloudFormationType() string {
	return "AWS::CodePipeline::Webhook.WebhookFilterRule"
}
