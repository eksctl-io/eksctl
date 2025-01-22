package nimblestudio

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// StudioComponent_StudioComponentConfiguration AWS CloudFormation Resource (AWS::NimbleStudio::StudioComponent.StudioComponentConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentconfiguration.html
type StudioComponent_StudioComponentConfiguration struct {

	// ActiveDirectoryConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentconfiguration.html#cfn-nimblestudio-studiocomponent-studiocomponentconfiguration-activedirectoryconfiguration
	ActiveDirectoryConfiguration *StudioComponent_ActiveDirectoryConfiguration `json:"ActiveDirectoryConfiguration,omitempty"`

	// ComputeFarmConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentconfiguration.html#cfn-nimblestudio-studiocomponent-studiocomponentconfiguration-computefarmconfiguration
	ComputeFarmConfiguration *StudioComponent_ComputeFarmConfiguration `json:"ComputeFarmConfiguration,omitempty"`

	// LicenseServiceConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentconfiguration.html#cfn-nimblestudio-studiocomponent-studiocomponentconfiguration-licenseserviceconfiguration
	LicenseServiceConfiguration *StudioComponent_LicenseServiceConfiguration `json:"LicenseServiceConfiguration,omitempty"`

	// SharedFileSystemConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentconfiguration.html#cfn-nimblestudio-studiocomponent-studiocomponentconfiguration-sharedfilesystemconfiguration
	SharedFileSystemConfiguration *StudioComponent_SharedFileSystemConfiguration `json:"SharedFileSystemConfiguration,omitempty"`

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
func (r *StudioComponent_StudioComponentConfiguration) AWSCloudFormationType() string {
	return "AWS::NimbleStudio::StudioComponent.StudioComponentConfiguration"
}
