package nimblestudio

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StudioComponent_ComputeFarmConfiguration AWS CloudFormation Resource (AWS::NimbleStudio::StudioComponent.ComputeFarmConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-computefarmconfiguration.html
type StudioComponent_ComputeFarmConfiguration struct {

	// ActiveDirectoryUser AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-computefarmconfiguration.html#cfn-nimblestudio-studiocomponent-computefarmconfiguration-activedirectoryuser
	ActiveDirectoryUser *types.Value `json:"ActiveDirectoryUser,omitempty"`

	// Endpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-computefarmconfiguration.html#cfn-nimblestudio-studiocomponent-computefarmconfiguration-endpoint
	Endpoint *types.Value `json:"Endpoint,omitempty"`

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
func (r *StudioComponent_ComputeFarmConfiguration) AWSCloudFormationType() string {
	return "AWS::NimbleStudio::StudioComponent.ComputeFarmConfiguration"
}
