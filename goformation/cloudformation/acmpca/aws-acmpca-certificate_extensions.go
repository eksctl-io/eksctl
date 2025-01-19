package acmpca

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Certificate_Extensions AWS CloudFormation Resource (AWS::ACMPCA::Certificate.Extensions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extensions.html
type Certificate_Extensions struct {

	// CertificatePolicies AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extensions.html#cfn-acmpca-certificate-extensions-certificatepolicies
	CertificatePolicies []Certificate_PolicyInformation `json:"CertificatePolicies,omitempty"`

	// ExtendedKeyUsage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extensions.html#cfn-acmpca-certificate-extensions-extendedkeyusage
	ExtendedKeyUsage []Certificate_ExtendedKeyUsage `json:"ExtendedKeyUsage,omitempty"`

	// KeyUsage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extensions.html#cfn-acmpca-certificate-extensions-keyusage
	KeyUsage *Certificate_KeyUsage `json:"KeyUsage,omitempty"`

	// SubjectAlternativeNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-extensions.html#cfn-acmpca-certificate-extensions-subjectalternativenames
	SubjectAlternativeNames []Certificate_GeneralName `json:"SubjectAlternativeNames,omitempty"`

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
func (r *Certificate_Extensions) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.Extensions"
}
