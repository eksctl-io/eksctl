package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_SalesforceConnectorProfileCredentials AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.SalesforceConnectorProfileCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofilecredentials.html
type ConnectorProfile_SalesforceConnectorProfileCredentials struct {

	// AccessToken AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofilecredentials.html#cfn-appflow-connectorprofile-salesforceconnectorprofilecredentials-accesstoken
	AccessToken *types.Value `json:"AccessToken,omitempty"`

	// ClientCredentialsArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofilecredentials.html#cfn-appflow-connectorprofile-salesforceconnectorprofilecredentials-clientcredentialsarn
	ClientCredentialsArn *types.Value `json:"ClientCredentialsArn,omitempty"`

	// ConnectorOAuthRequest AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofilecredentials.html#cfn-appflow-connectorprofile-salesforceconnectorprofilecredentials-connectoroauthrequest
	ConnectorOAuthRequest *ConnectorProfile_ConnectorOAuthRequest `json:"ConnectorOAuthRequest,omitempty"`

	// RefreshToken AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofilecredentials.html#cfn-appflow-connectorprofile-salesforceconnectorprofilecredentials-refreshtoken
	RefreshToken *types.Value `json:"RefreshToken,omitempty"`

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
func (r *ConnectorProfile_SalesforceConnectorProfileCredentials) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.SalesforceConnectorProfileCredentials"
}
