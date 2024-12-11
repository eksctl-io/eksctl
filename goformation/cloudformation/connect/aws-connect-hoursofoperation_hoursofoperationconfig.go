package connect

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// HoursOfOperation_HoursOfOperationConfig AWS CloudFormation Resource (AWS::Connect::HoursOfOperation.HoursOfOperationConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-hoursofoperation-hoursofoperationconfig.html
type HoursOfOperation_HoursOfOperationConfig struct {

	// Day AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-hoursofoperation-hoursofoperationconfig.html#cfn-connect-hoursofoperation-hoursofoperationconfig-day
	Day *types.Value `json:"Day,omitempty"`

	// EndTime AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-hoursofoperation-hoursofoperationconfig.html#cfn-connect-hoursofoperation-hoursofoperationconfig-endtime
	EndTime *HoursOfOperation_HoursOfOperationTimeSlice `json:"EndTime,omitempty"`

	// StartTime AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-hoursofoperation-hoursofoperationconfig.html#cfn-connect-hoursofoperation-hoursofoperationconfig-starttime
	StartTime *HoursOfOperation_HoursOfOperationTimeSlice `json:"StartTime,omitempty"`

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
func (r *HoursOfOperation_HoursOfOperationConfig) AWSCloudFormationType() string {
	return "AWS::Connect::HoursOfOperation.HoursOfOperationConfig"
}
