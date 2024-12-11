package eks

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// IdentityProviderConfig_OidcIdentityProviderConfig AWS CloudFormation Resource (AWS::EKS::IdentityProviderConfig.OidcIdentityProviderConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html
type IdentityProviderConfig_OidcIdentityProviderConfig struct {

	// ClientId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html#cfn-eks-identityproviderconfig-oidcidentityproviderconfig-clientid
	ClientId *types.Value `json:"ClientId,omitempty"`

	// GroupsClaim AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html#cfn-eks-identityproviderconfig-oidcidentityproviderconfig-groupsclaim
	GroupsClaim *types.Value `json:"GroupsClaim,omitempty"`

	// GroupsPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html#cfn-eks-identityproviderconfig-oidcidentityproviderconfig-groupsprefix
	GroupsPrefix *types.Value `json:"GroupsPrefix,omitempty"`

	// IssuerUrl AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html#cfn-eks-identityproviderconfig-oidcidentityproviderconfig-issuerurl
	IssuerUrl *types.Value `json:"IssuerUrl,omitempty"`

	// RequiredClaims AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html#cfn-eks-identityproviderconfig-oidcidentityproviderconfig-requiredclaims
	RequiredClaims []IdentityProviderConfig_RequiredClaim `json:"RequiredClaims,omitempty"`

	// UsernameClaim AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html#cfn-eks-identityproviderconfig-oidcidentityproviderconfig-usernameclaim
	UsernameClaim *types.Value `json:"UsernameClaim,omitempty"`

	// UsernamePrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-identityproviderconfig-oidcidentityproviderconfig.html#cfn-eks-identityproviderconfig-oidcidentityproviderconfig-usernameprefix
	UsernamePrefix *types.Value `json:"UsernamePrefix,omitempty"`

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
func (r *IdentityProviderConfig_OidcIdentityProviderConfig) AWSCloudFormationType() string {
	return "AWS::EKS::IdentityProviderConfig.OidcIdentityProviderConfig"
}
