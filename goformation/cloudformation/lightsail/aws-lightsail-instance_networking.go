package lightsail

import (
	"goformation/v4/cloudformation/policies"
)

// Instance_Networking AWS CloudFormation Resource (AWS::Lightsail::Instance.Networking)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-networking.html
type Instance_Networking struct {

	// MonthlyTransfer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-networking.html#cfn-lightsail-instance-networking-monthlytransfer
	MonthlyTransfer *Instance_MonthlyTransfer `json:"MonthlyTransfer,omitempty"`

	// Ports AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-networking.html#cfn-lightsail-instance-networking-ports
	Ports []Instance_Port `json:"Ports,omitempty"`

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
func (r *Instance_Networking) AWSCloudFormationType() string {
	return "AWS::Lightsail::Instance.Networking"
}
