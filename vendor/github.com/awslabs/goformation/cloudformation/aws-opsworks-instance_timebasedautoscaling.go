package cloudformation

import (
	"encoding/json"
)

// AWSOpsWorksInstance_TimeBasedAutoScaling AWS CloudFormation Resource (AWS::OpsWorks::Instance.TimeBasedAutoScaling)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html
type AWSOpsWorksInstance_TimeBasedAutoScaling struct {

	// Friday AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html#cfn-opsworks-instance-timebasedautoscaling-friday
	Friday map[string]*Value `json:"Friday,omitempty"`

	// Monday AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html#cfn-opsworks-instance-timebasedautoscaling-monday
	Monday map[string]*Value `json:"Monday,omitempty"`

	// Saturday AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html#cfn-opsworks-instance-timebasedautoscaling-saturday
	Saturday map[string]*Value `json:"Saturday,omitempty"`

	// Sunday AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html#cfn-opsworks-instance-timebasedautoscaling-sunday
	Sunday map[string]*Value `json:"Sunday,omitempty"`

	// Thursday AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html#cfn-opsworks-instance-timebasedautoscaling-thursday
	Thursday map[string]*Value `json:"Thursday,omitempty"`

	// Tuesday AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html#cfn-opsworks-instance-timebasedautoscaling-tuesday
	Tuesday map[string]*Value `json:"Tuesday,omitempty"`

	// Wednesday AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-instance-timebasedautoscaling.html#cfn-opsworks-instance-timebasedautoscaling-wednesday
	Wednesday map[string]*Value `json:"Wednesday,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksInstance_TimeBasedAutoScaling) AWSCloudFormationType() string {
	return "AWS::OpsWorks::Instance.TimeBasedAutoScaling"
}

func (r *AWSOpsWorksInstance_TimeBasedAutoScaling) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
