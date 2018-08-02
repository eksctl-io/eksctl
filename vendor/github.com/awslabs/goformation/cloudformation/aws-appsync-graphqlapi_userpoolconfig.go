package cloudformation

// AWSAppSyncGraphQLApi_UserPoolConfig AWS CloudFormation Resource (AWS::AppSync::GraphQLApi.UserPoolConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-userpoolconfig.html
type AWSAppSyncGraphQLApi_UserPoolConfig struct {

	// AppIdClientRegex AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-userpoolconfig.html#cfn-appsync-graphqlapi-userpoolconfig-appidclientregex
	AppIdClientRegex *StringIntrinsic `json:"AppIdClientRegex,omitempty"`

	// AwsRegion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-userpoolconfig.html#cfn-appsync-graphqlapi-userpoolconfig-awsregion
	AwsRegion *StringIntrinsic `json:"AwsRegion,omitempty"`

	// DefaultAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-userpoolconfig.html#cfn-appsync-graphqlapi-userpoolconfig-defaultaction
	DefaultAction *StringIntrinsic `json:"DefaultAction,omitempty"`

	// UserPoolId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-graphqlapi-userpoolconfig.html#cfn-appsync-graphqlapi-userpoolconfig-userpoolid
	UserPoolId *StringIntrinsic `json:"UserPoolId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncGraphQLApi_UserPoolConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::GraphQLApi.UserPoolConfig"
}
