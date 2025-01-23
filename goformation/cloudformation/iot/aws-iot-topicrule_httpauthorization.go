package iot

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// TopicRule_HttpAuthorization AWS CloudFormation Resource (AWS::IoT::TopicRule.HttpAuthorization)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-httpauthorization.html
type TopicRule_HttpAuthorization struct {

	// Sigv4 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-httpauthorization.html#cfn-iot-topicrule-httpauthorization-sigv4
	Sigv4 *TopicRule_SigV4Authorization `json:"Sigv4,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *TopicRule_HttpAuthorization) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.HttpAuthorization"
}
