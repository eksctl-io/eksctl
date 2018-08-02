package cloudformation

// AWSKinesisAnalyticsApplication_InputProcessingConfiguration AWS CloudFormation Resource (AWS::KinesisAnalytics::Application.InputProcessingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputprocessingconfiguration.html
type AWSKinesisAnalyticsApplication_InputProcessingConfiguration struct {

	// InputLambdaProcessor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-application-inputprocessingconfiguration.html#cfn-kinesisanalytics-application-inputprocessingconfiguration-inputlambdaprocessor
	InputLambdaProcessor *AWSKinesisAnalyticsApplication_InputLambdaProcessor `json:"InputLambdaProcessor,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsApplication_InputProcessingConfiguration) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::Application.InputProcessingConfiguration"
}
