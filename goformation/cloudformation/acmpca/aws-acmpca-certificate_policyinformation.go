package acmpca

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Certificate_PolicyInformation AWS CloudFormation Resource (AWS::ACMPCA::Certificate.PolicyInformation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyinformation.html
type Certificate_PolicyInformation struct {

	// CertPolicyId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyinformation.html#cfn-acmpca-certificate-policyinformation-certpolicyid
	CertPolicyId *types.Value `json:"CertPolicyId,omitempty"`

	// PolicyQualifiers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyinformation.html#cfn-acmpca-certificate-policyinformation-policyqualifiers
	PolicyQualifiers []Certificate_PolicyQualifierInfo `json:"PolicyQualifiers,omitempty"`

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
func (r *Certificate_PolicyInformation) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.PolicyInformation"
}
