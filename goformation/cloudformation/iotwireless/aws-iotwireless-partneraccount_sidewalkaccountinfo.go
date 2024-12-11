package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PartnerAccount_SidewalkAccountInfo AWS CloudFormation Resource (AWS::IoTWireless::PartnerAccount.SidewalkAccountInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-partneraccount-sidewalkaccountinfo.html
type PartnerAccount_SidewalkAccountInfo struct {

	// AppServerPrivateKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-partneraccount-sidewalkaccountinfo.html#cfn-iotwireless-partneraccount-sidewalkaccountinfo-appserverprivatekey
	AppServerPrivateKey *types.Value `json:"AppServerPrivateKey,omitempty"`

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
func (r *PartnerAccount_SidewalkAccountInfo) AWSCloudFormationType() string {
	return "AWS::IoTWireless::PartnerAccount.SidewalkAccountInfo"
}
