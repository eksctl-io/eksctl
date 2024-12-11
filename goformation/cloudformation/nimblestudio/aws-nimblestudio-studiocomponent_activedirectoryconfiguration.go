package nimblestudio

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StudioComponent_ActiveDirectoryConfiguration AWS CloudFormation Resource (AWS::NimbleStudio::StudioComponent.ActiveDirectoryConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-activedirectoryconfiguration.html
type StudioComponent_ActiveDirectoryConfiguration struct {

	// ComputerAttributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-activedirectoryconfiguration.html#cfn-nimblestudio-studiocomponent-activedirectoryconfiguration-computerattributes
	ComputerAttributes []StudioComponent_ActiveDirectoryComputerAttribute `json:"ComputerAttributes,omitempty"`

	// DirectoryId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-activedirectoryconfiguration.html#cfn-nimblestudio-studiocomponent-activedirectoryconfiguration-directoryid
	DirectoryId *types.Value `json:"DirectoryId,omitempty"`

	// OrganizationalUnitDistinguishedName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-activedirectoryconfiguration.html#cfn-nimblestudio-studiocomponent-activedirectoryconfiguration-organizationalunitdistinguishedname
	OrganizationalUnitDistinguishedName *types.Value `json:"OrganizationalUnitDistinguishedName,omitempty"`

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
func (r *StudioComponent_ActiveDirectoryConfiguration) AWSCloudFormationType() string {
	return "AWS::NimbleStudio::StudioComponent.ActiveDirectoryConfiguration"
}
