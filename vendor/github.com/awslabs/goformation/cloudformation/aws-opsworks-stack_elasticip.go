package cloudformation

import (
	"encoding/json"
)

// AWSOpsWorksStack_ElasticIp AWS CloudFormation Resource (AWS::OpsWorks::Stack.ElasticIp)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-stack-elasticip.html
type AWSOpsWorksStack_ElasticIp struct {

	// Ip AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-stack-elasticip.html#cfn-opsworks-stack-elasticip-ip
	Ip *Value `json:"Ip,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-stack-elasticip.html#cfn-opsworks-stack-elasticip-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksStack_ElasticIp) AWSCloudFormationType() string {
	return "AWS::OpsWorks::Stack.ElasticIp"
}

func (r *AWSOpsWorksStack_ElasticIp) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
