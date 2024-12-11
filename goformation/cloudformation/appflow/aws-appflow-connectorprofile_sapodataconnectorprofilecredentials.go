package appflow

import (
	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_SAPODataConnectorProfileCredentials AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.SAPODataConnectorProfileCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-sapodataconnectorprofilecredentials.html
type ConnectorProfile_SAPODataConnectorProfileCredentials struct {

	// BasicAuthCredentials AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-sapodataconnectorprofilecredentials.html#cfn-appflow-connectorprofile-sapodataconnectorprofilecredentials-basicauthcredentials
	BasicAuthCredentials interface{} `json:"BasicAuthCredentials,omitempty"`

	// OAuthCredentials AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-sapodataconnectorprofilecredentials.html#cfn-appflow-connectorprofile-sapodataconnectorprofilecredentials-oauthcredentials
	OAuthCredentials interface{} `json:"OAuthCredentials,omitempty"`

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
func (r *ConnectorProfile_SAPODataConnectorProfileCredentials) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.SAPODataConnectorProfileCredentials"
}
