package acmpca

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Certificate_PolicyQualifierInfo AWS CloudFormation Resource (AWS::ACMPCA::Certificate.PolicyQualifierInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyqualifierinfo.html
type Certificate_PolicyQualifierInfo struct {

	// PolicyQualifierId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyqualifierinfo.html#cfn-acmpca-certificate-policyqualifierinfo-policyqualifierid
	PolicyQualifierId *types.Value `json:"PolicyQualifierId,omitempty"`

	// Qualifier AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyqualifierinfo.html#cfn-acmpca-certificate-policyqualifierinfo-qualifier
	Qualifier *Certificate_Qualifier `json:"Qualifier,omitempty"`

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
func (r *Certificate_PolicyQualifierInfo) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.PolicyQualifierInfo"
}
