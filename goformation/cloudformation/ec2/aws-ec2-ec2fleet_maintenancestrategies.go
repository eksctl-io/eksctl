package ec2

import (
	"goformation/v4/cloudformation/policies"
)

// EC2Fleet_MaintenanceStrategies AWS CloudFormation Resource (AWS::EC2::EC2Fleet.MaintenanceStrategies)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-maintenancestrategies.html
type EC2Fleet_MaintenanceStrategies struct {

	// CapacityRebalance AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-maintenancestrategies.html#cfn-ec2-ec2fleet-maintenancestrategies-capacityrebalance
	CapacityRebalance *EC2Fleet_CapacityRebalance `json:"CapacityRebalance,omitempty"`

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
func (r *EC2Fleet_MaintenanceStrategies) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.MaintenanceStrategies"
}
