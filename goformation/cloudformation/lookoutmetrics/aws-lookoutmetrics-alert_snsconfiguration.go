package lookoutmetrics

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Alert_SNSConfiguration AWS CloudFormation Resource (AWS::LookoutMetrics::Alert.SNSConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-alert-snsconfiguration.html
type Alert_SNSConfiguration struct {

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-alert-snsconfiguration.html#cfn-lookoutmetrics-alert-snsconfiguration-rolearn
	RoleArn *types.Value `json:"RoleArn,omitempty"`

	// SnsTopicArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-alert-snsconfiguration.html#cfn-lookoutmetrics-alert-snsconfiguration-snstopicarn
	SnsTopicArn *types.Value `json:"SnsTopicArn,omitempty"`

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
func (r *Alert_SNSConfiguration) AWSCloudFormationType() string {
	return "AWS::LookoutMetrics::Alert.SNSConfiguration"
}
