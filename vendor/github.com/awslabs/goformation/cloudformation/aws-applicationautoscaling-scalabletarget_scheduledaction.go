package cloudformation

// AWSApplicationAutoScalingScalableTarget_ScheduledAction AWS CloudFormation Resource (AWS::ApplicationAutoScaling::ScalableTarget.ScheduledAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scheduledaction.html
type AWSApplicationAutoScalingScalableTarget_ScheduledAction struct {

	// EndTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scheduledaction.html#cfn-applicationautoscaling-scalabletarget-scheduledaction-endtime
	EndTime string `json:"EndTime,omitempty"`

	// ScalableTargetAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scheduledaction.html#cfn-applicationautoscaling-scalabletarget-scheduledaction-scalabletargetaction
	ScalableTargetAction *AWSApplicationAutoScalingScalableTarget_ScalableTargetAction `json:"ScalableTargetAction,omitempty"`

	// Schedule AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scheduledaction.html#cfn-applicationautoscaling-scalabletarget-scheduledaction-schedule
	Schedule *StringIntrinsic `json:"Schedule,omitempty"`

	// ScheduledActionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scheduledaction.html#cfn-applicationautoscaling-scalabletarget-scheduledaction-scheduledactionname
	ScheduledActionName *StringIntrinsic `json:"ScheduledActionName,omitempty"`

	// StartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalabletarget-scheduledaction.html#cfn-applicationautoscaling-scalabletarget-scheduledaction-starttime
	StartTime string `json:"StartTime,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApplicationAutoScalingScalableTarget_ScheduledAction) AWSCloudFormationType() string {
	return "AWS::ApplicationAutoScaling::ScalableTarget.ScheduledAction"
}
