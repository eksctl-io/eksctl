package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ClientVpnEndpoint_FederatedAuthenticationRequest AWS CloudFormation Resource (AWS::EC2::ClientVpnEndpoint.FederatedAuthenticationRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-federatedauthenticationrequest.html
type ClientVpnEndpoint_FederatedAuthenticationRequest struct {

	// SAMLProviderArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-federatedauthenticationrequest.html#cfn-ec2-clientvpnendpoint-federatedauthenticationrequest-samlproviderarn
	SAMLProviderArn *types.Value `json:"SAMLProviderArn,omitempty"`

	// SelfServiceSAMLProviderArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-clientvpnendpoint-federatedauthenticationrequest.html#cfn-ec2-clientvpnendpoint-federatedauthenticationrequest-selfservicesamlproviderarn
	SelfServiceSAMLProviderArn *types.Value `json:"SelfServiceSAMLProviderArn,omitempty"`

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
func (r *ClientVpnEndpoint_FederatedAuthenticationRequest) AWSCloudFormationType() string {
	return "AWS::EC2::ClientVpnEndpoint.FederatedAuthenticationRequest"
}
