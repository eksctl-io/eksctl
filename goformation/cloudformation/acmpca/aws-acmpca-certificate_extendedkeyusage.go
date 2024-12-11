package acmpca

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Certificate_ExtendedKeyUsage AWS CloudFormation Resource (AWS::ACMPCA::Certificate.ExtendedKeyUsage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extendedkeyusage.html
type Certificate_ExtendedKeyUsage struct {

	// ExtendedKeyUsageObjectIdentifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extendedkeyusage.html#cfn-acmpca-certificate-extendedkeyusage-extendedkeyusageobjectidentifier
	ExtendedKeyUsageObjectIdentifier *types.Value `json:"ExtendedKeyUsageObjectIdentifier,omitempty"`

	// ExtendedKeyUsageType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extendedkeyusage.html#cfn-acmpca-certificate-extendedkeyusage-extendedkeyusagetype
	ExtendedKeyUsageType *types.Value `json:"ExtendedKeyUsageType,omitempty"`

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
func (r *Certificate_ExtendedKeyUsage) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.ExtendedKeyUsage"
}
