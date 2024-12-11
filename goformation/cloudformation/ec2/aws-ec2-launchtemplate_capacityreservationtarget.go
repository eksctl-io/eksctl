package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LaunchTemplate_CapacityReservationTarget AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.CapacityReservationTarget)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-capacityreservationtarget.html
type LaunchTemplate_CapacityReservationTarget struct {

	// CapacityReservationId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-capacityreservationtarget.html#cfn-ec2-launchtemplate-capacityreservationtarget-capacityreservationid
	CapacityReservationId *types.Value `json:"CapacityReservationId,omitempty"`

	// CapacityReservationResourceGroupArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-capacityreservationtarget.html#cfn-ec2-launchtemplate-capacityreservationtarget-capacityreservationresourcegrouparn
	CapacityReservationResourceGroupArn *types.Value `json:"CapacityReservationResourceGroupArn,omitempty"`

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
func (r *LaunchTemplate_CapacityReservationTarget) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.CapacityReservationTarget"
}
