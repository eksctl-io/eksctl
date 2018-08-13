package cloudformation

import (
	"encoding/json"
)

// AWSApplicationAutoScalingScalableTarget_ScalableTargetAction AWS CloudFormation Resource (AWS::ApplicationAutoScaling::ScalableTarget.ScalableTargetAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scalabletargetaction.html
type AWSApplicationAutoScalingScalableTarget_ScalableTargetAction struct {

	// MaxCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scalabletargetaction.html#cfn-applicationautoscaling-scalabletarget-scalabletargetaction-maxcapacity
	MaxCapacity *Value `json:"MaxCapacity,omitempty"`

	// MinCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scalabletargetaction.html#cfn-applicationautoscaling-scalabletarget-scalabletargetaction-mincapacity
	MinCapacity *Value `json:"MinCapacity,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApplicationAutoScalingScalableTarget_ScalableTargetAction) AWSCloudFormationType() string {
	return "AWS::ApplicationAutoScaling::ScalableTarget.ScalableTargetAction"
}

func (r *AWSApplicationAutoScalingScalableTarget_ScalableTargetAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
