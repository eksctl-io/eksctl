package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ClientVpnEndpoint_CertificateAuthenticationRequest AWS CloudFormation Resource (AWS::EC2::ClientVpnEndpoint.CertificateAuthenticationRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-certificateauthenticationrequest.html
type ClientVpnEndpoint_CertificateAuthenticationRequest struct {

	// ClientRootCertificateChainArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-certificateauthenticationrequest.html#cfn-ec2-clientvpnendpoint-certificateauthenticationrequest-clientrootcertificatechainarn
	ClientRootCertificateChainArn *types.Value `json:"ClientRootCertificateChainArn,omitempty"`

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
func (r *ClientVpnEndpoint_CertificateAuthenticationRequest) AWSCloudFormationType() string {
	return "AWS::EC2::ClientVpnEndpoint.CertificateAuthenticationRequest"
}
