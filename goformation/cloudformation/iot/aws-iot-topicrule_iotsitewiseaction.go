package iot

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// TopicRule_IotSiteWiseAction AWS CloudFormation Resource (AWS::IoT::TopicRule.IotSiteWiseAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-iotsitewiseaction.html
type TopicRule_IotSiteWiseAction struct {

	// PutAssetPropertyValueEntries AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-iotsitewiseaction.html#cfn-iot-topicrule-iotsitewiseaction-putassetpropertyvalueentries
	PutAssetPropertyValueEntries []TopicRule_PutAssetPropertyValueEntry `json:"PutAssetPropertyValueEntries,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-iotsitewiseaction.html#cfn-iot-topicrule-iotsitewiseaction-rolearn
	RoleArn *types.Value `json:"RoleArn,omitempty"`

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
func (r *TopicRule_IotSiteWiseAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.IotSiteWiseAction"
}
