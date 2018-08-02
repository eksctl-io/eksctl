package cloudformation

// AWSCognitoUserPool_SmsConfiguration AWS CloudFormation Resource (AWS::Cognito::UserPool.SmsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-smsconfiguration.html
type AWSCognitoUserPool_SmsConfiguration struct {

	// ExternalId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-smsconfiguration.html#cfn-cognito-userpool-smsconfiguration-externalid
	ExternalId *StringIntrinsic `json:"ExternalId,omitempty"`

	// SnsCallerArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-smsconfiguration.html#cfn-cognito-userpool-smsconfiguration-snscallerarn
	SnsCallerArn *StringIntrinsic `json:"SnsCallerArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCognitoUserPool_SmsConfiguration) AWSCloudFormationType() string {
	return "AWS::Cognito::UserPool.SmsConfiguration"
}
