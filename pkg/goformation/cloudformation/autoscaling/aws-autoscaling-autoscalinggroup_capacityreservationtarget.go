package autoscaling

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// AutoScalingGroup_CapacityReservationTarget AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.CapacityReservationTarget)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-capacityreservationtarget.html
type AutoScalingGroup_CapacityReservationTarget struct {

	// CapacityReservationIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-capacityreservationtarget.html#cfn-autoscaling-autoscalinggroup-capacityreservationtarget-capacityreservationids
	CapacityReservationIds *types.Value `json:"CapacityReservationIds,omitempty"`

	// CapacityReservationResourceGroupArns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-capacityreservationtarget.html#cfn-autoscaling-autoscalinggroup-capacityreservationtarget-capacityreservationresourcegrouparns
	CapacityReservationResourceGroupArns *types.Value `json:"CapacityReservationResourceGroupArns,omitempty"`

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
func (r *AutoScalingGroup_CapacityReservationTarget) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.CapacityReservationTarget"
}
