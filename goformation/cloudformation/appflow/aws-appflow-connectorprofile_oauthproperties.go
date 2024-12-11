package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_OAuthProperties AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.OAuthProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-oauthproperties.html
type ConnectorProfile_OAuthProperties struct {

	// AuthCodeUrl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-oauthproperties.html#cfn-appflow-connectorprofile-oauthproperties-authcodeurl
	AuthCodeUrl *types.Value `json:"AuthCodeUrl,omitempty"`

	// OAuthScopes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-oauthproperties.html#cfn-appflow-connectorprofile-oauthproperties-oauthscopes
	OAuthScopes *types.Value `json:"OAuthScopes,omitempty"`

	// TokenUrl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-oauthproperties.html#cfn-appflow-connectorprofile-oauthproperties-tokenurl
	TokenUrl *types.Value `json:"TokenUrl,omitempty"`

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
func (r *ConnectorProfile_OAuthProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.OAuthProperties"
}
