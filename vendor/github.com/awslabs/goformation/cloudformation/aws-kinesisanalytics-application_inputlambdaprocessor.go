package cloudformation

// AWSKinesisAnalyticsApplication_InputLambdaProcessor AWS CloudFormation Resource (AWS::KinesisAnalytics::Application.InputLambdaProcessor)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputlambdaprocessor.html
type AWSKinesisAnalyticsApplication_InputLambdaProcessor struct {

	// ResourceARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputlambdaprocessor.html#cfn-kinesisanalytics-application-inputlambdaprocessor-resourcearn
	ResourceARN *StringIntrinsic `json:"ResourceARN,omitempty"`

	// RoleARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputlambdaprocessor.html#cfn-kinesisanalytics-application-inputlambdaprocessor-rolearn
	RoleARN *StringIntrinsic `json:"RoleARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplication_InputLambdaProcessor) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::Application.InputLambdaProcessor"
}
