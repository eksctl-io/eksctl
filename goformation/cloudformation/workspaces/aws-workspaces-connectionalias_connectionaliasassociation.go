package workspaces

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ConnectionAlias_ConnectionAliasAssociation AWS CloudFormation Resource (AWS::WorkSpaces::ConnectionAlias.ConnectionAliasAssociation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-connectionalias-connectionaliasassociation.html
type ConnectionAlias_ConnectionAliasAssociation struct {

	// AssociatedAccountId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-connectionalias-connectionaliasassociation.html#cfn-workspaces-connectionalias-connectionaliasassociation-associatedaccountid
	AssociatedAccountId *types.Value `json:"AssociatedAccountId,omitempty"`

	// AssociationStatus AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-connectionalias-connectionaliasassociation.html#cfn-workspaces-connectionalias-connectionaliasassociation-associationstatus
	AssociationStatus *types.Value `json:"AssociationStatus,omitempty"`

	// ConnectionIdentifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-connectionalias-connectionaliasassociation.html#cfn-workspaces-connectionalias-connectionaliasassociation-connectionidentifier
	ConnectionIdentifier *types.Value `json:"ConnectionIdentifier,omitempty"`

	// ResourceId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-workspaces-connectionalias-connectionaliasassociation.html#cfn-workspaces-connectionalias-connectionaliasassociation-resourceid
	ResourceId *types.Value `json:"ResourceId,omitempty"`

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
func (r *ConnectionAlias_ConnectionAliasAssociation) AWSCloudFormationType() string {
	return "AWS::WorkSpaces::ConnectionAlias.ConnectionAliasAssociation"
}
