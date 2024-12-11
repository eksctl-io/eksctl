package imagebuilder

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Image_ImageTestsConfiguration AWS CloudFormation Resource (AWS::ImageBuilder::Image.ImageTestsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-image-imagetestsconfiguration.html
type Image_ImageTestsConfiguration struct {

	// ImageTestsEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-image-imagetestsconfiguration.html#cfn-imagebuilder-image-imagetestsconfiguration-imagetestsenabled
	ImageTestsEnabled *types.Value `json:"ImageTestsEnabled,omitempty"`

	// TimeoutMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-image-imagetestsconfiguration.html#cfn-imagebuilder-image-imagetestsconfiguration-timeoutminutes
	TimeoutMinutes *types.Value `json:"TimeoutMinutes,omitempty"`

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
func (r *Image_ImageTestsConfiguration) AWSCloudFormationType() string {
	return "AWS::ImageBuilder::Image.ImageTestsConfiguration"
}
