package cloudformation

// AWSAppSyncDataSource_LambdaConfig AWS CloudFormation Resource (AWS::AppSync::DataSource.LambdaConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-lambdaconfig.html
type AWSAppSyncDataSource_LambdaConfig struct {

	// LambdaFunctionArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-datasource-lambdaconfig.html#cfn-appsync-datasource-lambdaconfig-lambdafunctionarn
	LambdaFunctionArn *StringIntrinsic `json:"LambdaFunctionArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncDataSource_LambdaConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::DataSource.LambdaConfig"
}
