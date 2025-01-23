package appflow

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ConnectorProfile_ConnectorProfileConfig AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.ConnectorProfileConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-connectorprofileconfig.html
type ConnectorProfile_ConnectorProfileConfig struct {

	// ConnectorProfileCredentials AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-connectorprofileconfig.html#cfn-appflow-connectorprofile-connectorprofileconfig-connectorprofilecredentials
	ConnectorProfileCredentials *ConnectorProfile_ConnectorProfileCredentials `json:"ConnectorProfileCredentials,omitempty"`

	// ConnectorProfileProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-connectorprofileconfig.html#cfn-appflow-connectorprofile-connectorprofileconfig-connectorprofileproperties
	ConnectorProfileProperties *ConnectorProfile_ConnectorProfileProperties `json:"ConnectorProfileProperties,omitempty"`

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
func (r *ConnectorProfile_ConnectorProfileConfig) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.ConnectorProfileConfig"
}
