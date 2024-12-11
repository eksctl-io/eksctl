package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// EC2Fleet_SpotOptionsRequest AWS CloudFormation Resource (AWS::EC2::EC2Fleet.SpotOptionsRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html
type EC2Fleet_SpotOptionsRequest struct {

	// AllocationStrategy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-allocationstrategy
	AllocationStrategy *types.Value `json:"AllocationStrategy,omitempty"`

	// InstanceInterruptionBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-instanceinterruptionbehavior
	InstanceInterruptionBehavior *types.Value `json:"InstanceInterruptionBehavior,omitempty"`

	// InstancePoolsToUseCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-instancepoolstousecount
	InstancePoolsToUseCount *types.Value `json:"InstancePoolsToUseCount,omitempty"`

	// MaintenanceStrategies AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-maintenancestrategies
	MaintenanceStrategies *EC2Fleet_MaintenanceStrategies `json:"MaintenanceStrategies,omitempty"`

	// MaxTotalPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-maxtotalprice
	MaxTotalPrice *types.Value `json:"MaxTotalPrice,omitempty"`

	// MinTargetCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-mintargetcapacity
	MinTargetCapacity *types.Value `json:"MinTargetCapacity,omitempty"`

	// SingleAvailabilityZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-singleavailabilityzone
	SingleAvailabilityZone *types.Value `json:"SingleAvailabilityZone,omitempty"`

	// SingleInstanceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-singleinstancetype
	SingleInstanceType *types.Value `json:"SingleInstanceType,omitempty"`

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
func (r *EC2Fleet_SpotOptionsRequest) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.SpotOptionsRequest"
}
