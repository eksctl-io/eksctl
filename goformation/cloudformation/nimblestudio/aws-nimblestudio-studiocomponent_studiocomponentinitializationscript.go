package nimblestudio

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StudioComponent_StudioComponentInitializationScript AWS CloudFormation Resource (AWS::NimbleStudio::StudioComponent.StudioComponentInitializationScript)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentinitializationscript.html
type StudioComponent_StudioComponentInitializationScript struct {

	// LaunchProfileProtocolVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentinitializationscript.html#cfn-nimblestudio-studiocomponent-studiocomponentinitializationscript-launchprofileprotocolversion
	LaunchProfileProtocolVersion *types.Value `json:"LaunchProfileProtocolVersion,omitempty"`

	// Platform AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentinitializationscript.html#cfn-nimblestudio-studiocomponent-studiocomponentinitializationscript-platform
	Platform *types.Value `json:"Platform,omitempty"`

	// RunContext AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentinitializationscript.html#cfn-nimblestudio-studiocomponent-studiocomponentinitializationscript-runcontext
	RunContext *types.Value `json:"RunContext,omitempty"`

	// Script AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-studiocomponentinitializationscript.html#cfn-nimblestudio-studiocomponent-studiocomponentinitializationscript-script
	Script *types.Value `json:"Script,omitempty"`

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
func (r *StudioComponent_StudioComponentInitializationScript) AWSCloudFormationType() string {
	return "AWS::NimbleStudio::StudioComponent.StudioComponentInitializationScript"
}
