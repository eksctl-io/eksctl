package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Subnet_PrivateDnsNameOptionsOnLaunch AWS CloudFormation Resource (AWS::EC2::Subnet.PrivateDnsNameOptionsOnLaunch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-subnet-privatednsnameoptionsonlaunch.html
type Subnet_PrivateDnsNameOptionsOnLaunch struct {

	// EnableResourceNameDnsAAAARecord AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-subnet-privatednsnameoptionsonlaunch.html#cfn-ec2-subnet-privatednsnameoptionsonlaunch-enableresourcenamednsaaaarecord
	EnableResourceNameDnsAAAARecord *types.Value `json:"EnableResourceNameDnsAAAARecord,omitempty"`

	// EnableResourceNameDnsARecord AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-subnet-privatednsnameoptionsonlaunch.html#cfn-ec2-subnet-privatednsnameoptionsonlaunch-enableresourcenamednsarecord
	EnableResourceNameDnsARecord *types.Value `json:"EnableResourceNameDnsARecord,omitempty"`

	// HostnameType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-subnet-privatednsnameoptionsonlaunch.html#cfn-ec2-subnet-privatednsnameoptionsonlaunch-hostnametype
	HostnameType *types.Value `json:"HostnameType,omitempty"`

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
func (r *Subnet_PrivateDnsNameOptionsOnLaunch) AWSCloudFormationType() string {
	return "AWS::EC2::Subnet.PrivateDnsNameOptionsOnLaunch"
}
