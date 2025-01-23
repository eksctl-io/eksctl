package acmpca

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// CertificateAuthority_AccessDescription AWS CloudFormation Resource (AWS::ACMPCA::CertificateAuthority.AccessDescription)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-accessdescription.html
type CertificateAuthority_AccessDescription struct {

	// AccessLocation AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-accessdescription.html#cfn-acmpca-certificateauthority-accessdescription-accesslocation
	AccessLocation *CertificateAuthority_GeneralName `json:"AccessLocation,omitempty"`

	// AccessMethod AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-accessdescription.html#cfn-acmpca-certificateauthority-accessdescription-accessmethod
	AccessMethod *CertificateAuthority_AccessMethod `json:"AccessMethod,omitempty"`

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
func (r *CertificateAuthority_AccessDescription) AWSCloudFormationType() string {
	return "AWS::ACMPCA::CertificateAuthority.AccessDescription"
}
