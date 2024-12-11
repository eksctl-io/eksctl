package nimblestudio

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StudioComponent_LicenseServiceConfiguration AWS CloudFormation Resource (AWS::NimbleStudio::StudioComponent.LicenseServiceConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-licenseserviceconfiguration.html
type StudioComponent_LicenseServiceConfiguration struct {

	// Endpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-studiocomponent-licenseserviceconfiguration.html#cfn-nimblestudio-studiocomponent-licenseserviceconfiguration-endpoint
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
func (r *StudioComponent_LicenseServiceConfiguration) AWSCloudFormationType() string {
	return "AWS::NimbleStudio::StudioComponent.LicenseServiceConfiguration"
}
