package ec2

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// SpotFleet_SpotMaintenanceStrategies AWS CloudFormation Resource (AWS::EC2::SpotFleet.SpotMaintenanceStrategies)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotmaintenancestrategies.html
type SpotFleet_SpotMaintenanceStrategies struct {

	// CapacityRebalance AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotmaintenancestrategies.html#cfn-ec2-spotfleet-spotmaintenancestrategies-capacityrebalance
	CapacityRebalance *SpotFleet_SpotCapacityRebalance `json:"CapacityRebalance,omitempty"`

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
func (r *SpotFleet_SpotMaintenanceStrategies) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.SpotMaintenanceStrategies"
}
