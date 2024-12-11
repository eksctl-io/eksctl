package acmpca

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CertificateAuthority_OcspConfiguration AWS CloudFormation Resource (AWS::ACMPCA::CertificateAuthority.OcspConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-ocspconfiguration.html
type CertificateAuthority_OcspConfiguration struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-ocspconfiguration.html#cfn-acmpca-certificateauthority-ocspconfiguration-enabled
	Enabled *types.Value `json:"Enabled,omitempty"`

	// OcspCustomCname AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-ocspconfiguration.html#cfn-acmpca-certificateauthority-ocspconfiguration-ocspcustomcname
	OcspCustomCname *types.Value `json:"OcspCustomCname,omitempty"`

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
func (r *CertificateAuthority_OcspConfiguration) AWSCloudFormationType() string {
	return "AWS::ACMPCA::CertificateAuthority.OcspConfiguration"
}
