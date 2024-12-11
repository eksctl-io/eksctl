package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_SingularConnectorProfileCredentials AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.SingularConnectorProfileCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-singularconnectorprofilecredentials.html
type ConnectorProfile_SingularConnectorProfileCredentials struct {

	// ApiKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-singularconnectorprofilecredentials.html#cfn-appflow-connectorprofile-singularconnectorprofilecredentials-apikey
	ApiKey *types.Value `json:"ApiKey,omitempty"`

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
func (r *ConnectorProfile_SingularConnectorProfileCredentials) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.SingularConnectorProfileCredentials"
}
