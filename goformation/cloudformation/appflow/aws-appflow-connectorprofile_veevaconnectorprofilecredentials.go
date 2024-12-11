package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_VeevaConnectorProfileCredentials AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.VeevaConnectorProfileCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-veevaconnectorprofilecredentials.html
type ConnectorProfile_VeevaConnectorProfileCredentials struct {

	// Password AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-veevaconnectorprofilecredentials.html#cfn-appflow-connectorprofile-veevaconnectorprofilecredentials-password
	Password *types.Value `json:"Password,omitempty"`

	// Username AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-veevaconnectorprofilecredentials.html#cfn-appflow-connectorprofile-veevaconnectorprofilecredentials-username
	Username *types.Value `json:"Username,omitempty"`

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
func (r *ConnectorProfile_VeevaConnectorProfileCredentials) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.VeevaConnectorProfileCredentials"
}
