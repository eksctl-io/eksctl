package cloudformation

// AWSAppSyncGraphQLApi_LogConfig AWS CloudFormation Resource (AWS::AppSync::GraphQLApi.LogConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-logconfig.html
type AWSAppSyncGraphQLApi_LogConfig struct {

	// CloudWatchLogsRoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-logconfig.html#cfn-appsync-graphqlapi-logconfig-cloudwatchlogsrolearn
	CloudWatchLogsRoleArn *StringIntrinsic `json:"CloudWatchLogsRoleArn,omitempty"`

	// FieldLogLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-logconfig.html#cfn-appsync-graphqlapi-logconfig-fieldloglevel
	FieldLogLevel *StringIntrinsic `json:"FieldLogLevel,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncGraphQLApi_LogConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::GraphQLApi.LogConfig"
}
