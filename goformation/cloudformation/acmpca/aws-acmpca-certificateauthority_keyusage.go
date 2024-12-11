package acmpca

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CertificateAuthority_KeyUsage AWS CloudFormation Resource (AWS::ACMPCA::CertificateAuthority.KeyUsage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html
type CertificateAuthority_KeyUsage struct {

	// CRLSign AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-crlsign
	CRLSign *types.Value `json:"CRLSign,omitempty"`

	// DataEncipherment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-dataencipherment
	DataEncipherment *types.Value `json:"DataEncipherment,omitempty"`

	// DecipherOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-decipheronly
	DecipherOnly *types.Value `json:"DecipherOnly,omitempty"`

	// DigitalSignature AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-digitalsignature
	DigitalSignature *types.Value `json:"DigitalSignature,omitempty"`

	// EncipherOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-encipheronly
	EncipherOnly *types.Value `json:"EncipherOnly,omitempty"`

	// KeyAgreement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-keyagreement
	KeyAgreement *types.Value `json:"KeyAgreement,omitempty"`

	// KeyCertSign AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-keycertsign
	KeyCertSign *types.Value `json:"KeyCertSign,omitempty"`

	// KeyEncipherment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-keyencipherment
	KeyEncipherment *types.Value `json:"KeyEncipherment,omitempty"`

	// NonRepudiation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-nonrepudiation
	NonRepudiation *types.Value `json:"NonRepudiation,omitempty"`

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
func (r *CertificateAuthority_KeyUsage) AWSCloudFormationType() string {
	return "AWS::ACMPCA::CertificateAuthority.KeyUsage"
}
