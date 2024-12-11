package signer

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SigningProfile_SignatureValidityPeriod AWS CloudFormation Resource (AWS::Signer::SigningProfile.SignatureValidityPeriod)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-signer-signingprofile-signaturevalidityperiod.html
type SigningProfile_SignatureValidityPeriod struct {

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-signer-signingprofile-signaturevalidityperiod.html#cfn-signer-signingprofile-signaturevalidityperiod-type
	Type *types.Value `json:"Type,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-signer-signingprofile-signaturevalidityperiod.html#cfn-signer-signingprofile-signaturevalidityperiod-value
	Value *types.Value `json:"Value,omitempty"`

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
func (r *SigningProfile_SignatureValidityPeriod) AWSCloudFormationType() string {
	return "AWS::Signer::SigningProfile.SignatureValidityPeriod"
}
