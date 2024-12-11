package autoscalingplans

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ScalingPlan_CustomizedScalingMetricSpecification AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.CustomizedScalingMetricSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-customizedscalingmetricspecification.html
type ScalingPlan_CustomizedScalingMetricSpecification struct {

	// Dimensions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-customizedscalingmetricspecification.html#cfn-autoscalingplans-scalingplan-customizedscalingmetricspecification-dimensions
	Dimensions []ScalingPlan_MetricDimension `json:"Dimensions,omitempty"`

	// MetricName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-customizedscalingmetricspecification.html#cfn-autoscalingplans-scalingplan-customizedscalingmetricspecification-metricname
	MetricName *types.Value `json:"MetricName,omitempty"`

	// Namespace AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-customizedscalingmetricspecification.html#cfn-autoscalingplans-scalingplan-customizedscalingmetricspecification-namespace
	Namespace *types.Value `json:"Namespace,omitempty"`

	// Statistic AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-customizedscalingmetricspecification.html#cfn-autoscalingplans-scalingplan-customizedscalingmetricspecification-statistic
	Statistic *types.Value `json:"Statistic,omitempty"`

	// Unit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-customizedscalingmetricspecification.html#cfn-autoscalingplans-scalingplan-customizedscalingmetricspecification-unit
	Unit *types.Value `json:"Unit,omitempty"`

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
func (r *ScalingPlan_CustomizedScalingMetricSpecification) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.CustomizedScalingMetricSpecification"
}
