package appintegrations

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// EventIntegration_EventIntegrationAssociation AWS CloudFormation Resource (AWS::AppIntegrations::EventIntegration.EventIntegrationAssociation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appintegrations-eventintegration-eventintegrationassociation.html
type EventIntegration_EventIntegrationAssociation struct {

	// ClientAssociationMetadata AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appintegrations-eventintegration-eventintegrationassociation.html#cfn-appintegrations-eventintegration-eventintegrationassociation-clientassociationmetadata
	ClientAssociationMetadata []EventIntegration_Metadata `json:"ClientAssociationMetadata,omitempty"`

	// ClientId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appintegrations-eventintegration-eventintegrationassociation.html#cfn-appintegrations-eventintegration-eventintegrationassociation-clientid
	ClientId *types.Value `json:"ClientId,omitempty"`

	// EventBridgeRuleName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appintegrations-eventintegration-eventintegrationassociation.html#cfn-appintegrations-eventintegration-eventintegrationassociation-eventbridgerulename
	EventBridgeRuleName *types.Value `json:"EventBridgeRuleName,omitempty"`

	// EventIntegrationAssociationArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appintegrations-eventintegration-eventintegrationassociation.html#cfn-appintegrations-eventintegration-eventintegrationassociation-eventintegrationassociationarn
	EventIntegrationAssociationArn *types.Value `json:"EventIntegrationAssociationArn,omitempty"`

	// EventIntegrationAssociationId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appintegrations-eventintegration-eventintegrationassociation.html#cfn-appintegrations-eventintegration-eventintegrationassociation-eventintegrationassociationid
	EventIntegrationAssociationId *types.Value `json:"EventIntegrationAssociationId,omitempty"`

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
func (r *EventIntegration_EventIntegrationAssociation) AWSCloudFormationType() string {
	return "AWS::AppIntegrations::EventIntegration.EventIntegrationAssociation"
}
