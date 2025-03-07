package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VerifiedAccessTrustProvider_NativeApplicationOidcOptions AWS CloudFormation Resource (AWS::EC2::VerifiedAccessTrustProvider.NativeApplicationOidcOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html
type VerifiedAccessTrustProvider_NativeApplicationOidcOptions struct {

	// AuthorizationEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-authorizationendpoint
	AuthorizationEndpoint *types.Value `json:"AuthorizationEndpoint,omitempty"`

	// ClientId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-clientid
	ClientId *types.Value `json:"ClientId,omitempty"`

	// ClientSecret AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-clientsecret
	ClientSecret *types.Value `json:"ClientSecret,omitempty"`

	// Issuer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-issuer
	Issuer *types.Value `json:"Issuer,omitempty"`

	// PublicSigningKeyEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-publicsigningkeyendpoint
	PublicSigningKeyEndpoint *types.Value `json:"PublicSigningKeyEndpoint,omitempty"`

	// Scope AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-scope
	Scope *types.Value `json:"Scope,omitempty"`

	// TokenEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-tokenendpoint
	TokenEndpoint *types.Value `json:"TokenEndpoint,omitempty"`

	// UserInfoEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions.html#cfn-ec2-verifiedaccesstrustprovider-nativeapplicationoidcoptions-userinfoendpoint
	UserInfoEndpoint *types.Value `json:"UserInfoEndpoint,omitempty"`

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
func (r *VerifiedAccessTrustProvider_NativeApplicationOidcOptions) AWSCloudFormationType() string {
	return "AWS::EC2::VerifiedAccessTrustProvider.NativeApplicationOidcOptions"
}
