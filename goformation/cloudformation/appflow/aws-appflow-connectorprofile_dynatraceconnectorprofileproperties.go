package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_DynatraceConnectorProfileProperties AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.DynatraceConnectorProfileProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-dynatraceconnectorprofileproperties.html
type ConnectorProfile_DynatraceConnectorProfileProperties struct {

	// InstanceUrl AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-dynatraceconnectorprofileproperties.html#cfn-appflow-connectorprofile-dynatraceconnectorprofileproperties-instanceurl
	InstanceUrl *types.Value `json:"InstanceUrl,omitempty"`

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
func (r *ConnectorProfile_DynatraceConnectorProfileProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.DynatraceConnectorProfileProperties"
}
