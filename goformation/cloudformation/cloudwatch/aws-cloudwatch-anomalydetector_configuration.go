package cloudwatch

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AnomalyDetector_Configuration AWS CloudFormation Resource (AWS::CloudWatch::AnomalyDetector.Configuration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-anomalydetector-configuration.html
type AnomalyDetector_Configuration struct {

	// ExcludedTimeRanges AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-anomalydetector-configuration.html#cfn-cloudwatch-anomalydetector-configuration-excludedtimeranges
	ExcludedTimeRanges []AnomalyDetector_Range `json:"ExcludedTimeRanges,omitempty"`

	// MetricTimeZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-anomalydetector-configuration.html#cfn-cloudwatch-anomalydetector-configuration-metrictimezone
	MetricTimeZone *types.Value `json:"MetricTimeZone,omitempty"`

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
func (r *AnomalyDetector_Configuration) AWSCloudFormationType() string {
	return "AWS::CloudWatch::AnomalyDetector.Configuration"
}
