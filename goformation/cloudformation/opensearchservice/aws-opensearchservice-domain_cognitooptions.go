package opensearchservice

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Domain_CognitoOptions AWS CloudFormation Resource (AWS::OpenSearchService::Domain.CognitoOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opensearchservice-domain-cognitooptions.html
type Domain_CognitoOptions struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opensearchservice-domain-cognitooptions.html#cfn-opensearchservice-domain-cognitooptions-enabled
	Enabled *types.Value `json:"Enabled,omitempty"`

	// IdentityPoolId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opensearchservice-domain-cognitooptions.html#cfn-opensearchservice-domain-cognitooptions-identitypoolid
	IdentityPoolId *types.Value `json:"IdentityPoolId,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opensearchservice-domain-cognitooptions.html#cfn-opensearchservice-domain-cognitooptions-rolearn
	RoleArn *types.Value `json:"RoleArn,omitempty"`

	// UserPoolId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opensearchservice-domain-cognitooptions.html#cfn-opensearchservice-domain-cognitooptions-userpoolid
	UserPoolId *types.Value `json:"UserPoolId,omitempty"`

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
func (r *Domain_CognitoOptions) AWSCloudFormationType() string {
	return "AWS::OpenSearchService::Domain.CognitoOptions"
}
