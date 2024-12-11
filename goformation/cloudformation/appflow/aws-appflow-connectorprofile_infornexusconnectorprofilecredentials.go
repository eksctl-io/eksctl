package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectorProfile_InforNexusConnectorProfileCredentials AWS CloudFormation Resource (AWS::AppFlow::ConnectorProfile.InforNexusConnectorProfileCredentials)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-infornexusconnectorprofilecredentials.html
type ConnectorProfile_InforNexusConnectorProfileCredentials struct {

	// AccessKeyId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-infornexusconnectorprofilecredentials.html#cfn-appflow-connectorprofile-infornexusconnectorprofilecredentials-accesskeyid
	AccessKeyId *types.Value `json:"AccessKeyId,omitempty"`

	// Datakey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-infornexusconnectorprofilecredentials.html#cfn-appflow-connectorprofile-infornexusconnectorprofilecredentials-datakey
	Datakey *types.Value `json:"Datakey,omitempty"`

	// SecretAccessKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-infornexusconnectorprofilecredentials.html#cfn-appflow-connectorprofile-infornexusconnectorprofilecredentials-secretaccesskey
	SecretAccessKey *types.Value `json:"SecretAccessKey,omitempty"`

	// UserId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-connectorprofile-infornexusconnectorprofilecredentials.html#cfn-appflow-connectorprofile-infornexusconnectorprofilecredentials-userid
	UserId *types.Value `json:"UserId,omitempty"`

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
func (r *ConnectorProfile_InforNexusConnectorProfileCredentials) AWSCloudFormationType() string {
	return "AWS::AppFlow::ConnectorProfile.InforNexusConnectorProfileCredentials"
}
