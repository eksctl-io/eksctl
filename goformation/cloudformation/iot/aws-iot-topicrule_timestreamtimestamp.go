package iot

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// TopicRule_TimestreamTimestamp AWS CloudFormation Resource (AWS::IoT::TopicRule.TimestreamTimestamp)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-timestreamtimestamp.html
type TopicRule_TimestreamTimestamp struct {

	// Unit AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-timestreamtimestamp.html#cfn-iot-topicrule-timestreamtimestamp-unit
	Unit *types.Value `json:"Unit,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-timestreamtimestamp.html#cfn-iot-topicrule-timestreamtimestamp-value
	Value *types.Value `json:"Value,omitempty"`

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
func (r *TopicRule_TimestreamTimestamp) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.TimestreamTimestamp"
}
