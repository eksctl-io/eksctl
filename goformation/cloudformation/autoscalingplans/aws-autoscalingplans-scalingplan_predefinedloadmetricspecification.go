package autoscalingplans

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ScalingPlan_PredefinedLoadMetricSpecification AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.PredefinedLoadMetricSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedloadmetricspecification.html
type ScalingPlan_PredefinedLoadMetricSpecification struct {

	// PredefinedLoadMetricType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedloadmetricspecification.html#cfn-autoscalingplans-scalingplan-predefinedloadmetricspecification-predefinedloadmetrictype
	PredefinedLoadMetricType *types.Value `json:"PredefinedLoadMetricType,omitempty"`

	// ResourceLabel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedloadmetricspecification.html#cfn-autoscalingplans-scalingplan-predefinedloadmetricspecification-resourcelabel
	ResourceLabel *types.Value `json:"ResourceLabel,omitempty"`

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
func (r *ScalingPlan_PredefinedLoadMetricSpecification) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.PredefinedLoadMetricSpecification"
}
