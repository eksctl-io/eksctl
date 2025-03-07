package autoscaling

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// ScalingPolicy_PredictiveScalingCustomizedCapacityMetric AWS CloudFormation Resource (AWS::AutoScaling::ScalingPolicy.PredictiveScalingCustomizedCapacityMetric)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingcustomizedcapacitymetric.html
type ScalingPolicy_PredictiveScalingCustomizedCapacityMetric struct {

	// MetricDataQueries AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingcustomizedcapacitymetric.html#cfn-autoscaling-scalingpolicy-predictivescalingcustomizedcapacitymetric-metricdataqueries
	MetricDataQueries []ScalingPolicy_MetricDataQuery `json:"MetricDataQueries,omitempty"`

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
func (r *ScalingPolicy_PredictiveScalingCustomizedCapacityMetric) AWSCloudFormationType() string {
	return "AWS::AutoScaling::ScalingPolicy.PredictiveScalingCustomizedCapacityMetric"
}
