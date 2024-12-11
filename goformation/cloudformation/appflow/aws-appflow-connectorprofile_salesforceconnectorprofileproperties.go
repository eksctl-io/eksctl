package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_SalesforceConnectorProfileProperties AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.SalesforceConnectorProfileProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofileproperties.html
type ConnectorProfile_SalesforceConnectorProfileProperties struct {

	// InstanceUrl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofileproperties.html#cfn-appflow-connectorprofile-salesforceconnectorprofileproperties-instanceurl
	InstanceUrl *types.Value `json:"InstanceUrl,omitempty"`

	// isSandboxEnvironment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-salesforceconnectorprofileproperties.html#cfn-appflow-connectorprofile-salesforceconnectorprofileproperties-issandboxenvironment
	isSandboxEnvironment *types.Value `json:"isSandboxEnvironment,omitempty"`

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
func (r *ConnectorProfile_SalesforceConnectorProfileProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.SalesforceConnectorProfileProperties"
}
