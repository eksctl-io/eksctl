package imagebuilder

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ContainerRecipe_InstanceConfiguration AWS CloudFormation Resource (AWS::ImageBuilder::ContainerRecipe.InstanceConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-containerrecipe-instanceconfiguration.html
type ContainerRecipe_InstanceConfiguration struct {

	// BlockDeviceMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-containerrecipe-instanceconfiguration.html#cfn-imagebuilder-containerrecipe-instanceconfiguration-blockdevicemappings
	BlockDeviceMappings []ContainerRecipe_InstanceBlockDeviceMapping `json:"BlockDeviceMappings,omitempty"`

	// Image AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-containerrecipe-instanceconfiguration.html#cfn-imagebuilder-containerrecipe-instanceconfiguration-image
	Image *types.Value `json:"Image,omitempty"`

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
func (r *ContainerRecipe_InstanceConfiguration) AWSCloudFormationType() string {
	return "AWS::ImageBuilder::ContainerRecipe.InstanceConfiguration"
}
