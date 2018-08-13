package cloudformation

import (
	"encoding/json"
)

// AWSOpsWorksApp_SslConfiguration AWS CloudFormation Resource (AWS::OpsWorks::App.SslConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-app-sslconfiguration.html
type AWSOpsWorksApp_SslConfiguration struct {

	// Certificate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-app-sslconfiguration.html#cfn-opsworks-app-sslconfig-certificate
	Certificate *Value `json:"Certificate,omitempty"`

	// Chain AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-app-sslconfiguration.html#cfn-opsworks-app-sslconfig-chain
	Chain *Value `json:"Chain,omitempty"`

	// PrivateKey AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-app-sslconfiguration.html#cfn-opsworks-app-sslconfig-privatekey
	PrivateKey *Value `json:"PrivateKey,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksApp_SslConfiguration) AWSCloudFormationType() string {
	return "AWS::OpsWorks::App.SslConfiguration"
}

func (r *AWSOpsWorksApp_SslConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
