package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_TrendmicroConnectorProfileCredentials AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.TrendmicroConnectorProfileCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-trendmicroconnectorprofilecredentials.html
type ConnectorProfile_TrendmicroConnectorProfileCredentials struct {

	// ApiSecretKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-trendmicroconnectorprofilecredentials.html#cfn-appflow-connectorprofile-trendmicroconnectorprofilecredentials-apisecretkey
	ApiSecretKey *types.Value `json:"ApiSecretKey,omitempty"`

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
func (r *ConnectorProfile_TrendmicroConnectorProfileCredentials) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.TrendmicroConnectorProfileCredentials"
}
