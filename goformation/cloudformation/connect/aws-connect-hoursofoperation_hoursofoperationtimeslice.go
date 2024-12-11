package connect

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// HoursOfOperation_HoursOfOperationTimeSlice AWS CloudFormation Resource (AWS::Connect::HoursOfOperation.HoursOfOperationTimeSlice)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-hoursofoperation-hoursofoperationtimeslice.html
type HoursOfOperation_HoursOfOperationTimeSlice struct {

	// Hours AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-hoursofoperation-hoursofoperationtimeslice.html#cfn-connect-hoursofoperation-hoursofoperationtimeslice-hours
	Hours *types.Value `json:"Hours"`

	// Minutes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-hoursofoperation-hoursofoperationtimeslice.html#cfn-connect-hoursofoperation-hoursofoperationtimeslice-minutes
	Minutes *types.Value `json:"Minutes"`

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
func (r *HoursOfOperation_HoursOfOperationTimeSlice) AWSCloudFormationType() string {
	return "AWS::Connect::HoursOfOperation.HoursOfOperationTimeSlice"
}
