package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingLoadBalancer_AccessLoggingPolicy AWS CloudFormation Resource (AWS::ElasticLoadBalancing::LoadBalancer.AccessLoggingPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-accessloggingpolicy.html
type AWSElasticLoadBalancingLoadBalancer_AccessLoggingPolicy struct {

	// EmitInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-accessloggingpolicy.html#cfn-elb-accessloggingpolicy-emitinterval
	EmitInterval *Value `json:"EmitInterval,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-accessloggingpolicy.html#cfn-elb-accessloggingpolicy-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// S3BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-accessloggingpolicy.html#cfn-elb-accessloggingpolicy-s3bucketname
	S3BucketName *Value `json:"S3BucketName,omitempty"`

	// S3BucketPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-accessloggingpolicy.html#cfn-elb-accessloggingpolicy-s3bucketprefix
	S3BucketPrefix *Value `json:"S3BucketPrefix,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingLoadBalancer_AccessLoggingPolicy) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer.AccessLoggingPolicy"
}

func (r *AWSElasticLoadBalancingLoadBalancer_AccessLoggingPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
