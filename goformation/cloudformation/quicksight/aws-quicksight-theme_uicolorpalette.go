package quicksight

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Theme_UIColorPalette AWS CloudFormation Resource (AWS::QuickSight::Theme.UIColorPalette)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html
type Theme_UIColorPalette struct {

	// Accent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-accent
	Accent *types.Value `json:"Accent,omitempty"`

	// AccentForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-accentforeground
	AccentForeground *types.Value `json:"AccentForeground,omitempty"`

	// Danger AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-danger
	Danger *types.Value `json:"Danger,omitempty"`

	// DangerForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-dangerforeground
	DangerForeground *types.Value `json:"DangerForeground,omitempty"`

	// Dimension AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-dimension
	Dimension *types.Value `json:"Dimension,omitempty"`

	// DimensionForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-dimensionforeground
	DimensionForeground *types.Value `json:"DimensionForeground,omitempty"`

	// Measure AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-measure
	Measure *types.Value `json:"Measure,omitempty"`

	// MeasureForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-measureforeground
	MeasureForeground *types.Value `json:"MeasureForeground,omitempty"`

	// PrimaryBackground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-primarybackground
	PrimaryBackground *types.Value `json:"PrimaryBackground,omitempty"`

	// PrimaryForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-primaryforeground
	PrimaryForeground *types.Value `json:"PrimaryForeground,omitempty"`

	// SecondaryBackground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-secondarybackground
	SecondaryBackground *types.Value `json:"SecondaryBackground,omitempty"`

	// SecondaryForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-secondaryforeground
	SecondaryForeground *types.Value `json:"SecondaryForeground,omitempty"`

	// Success AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-success
	Success *types.Value `json:"Success,omitempty"`

	// SuccessForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-successforeground
	SuccessForeground *types.Value `json:"SuccessForeground,omitempty"`

	// Warning AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-warning
	Warning *types.Value `json:"Warning,omitempty"`

	// WarningForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-warningforeground
	WarningForeground *types.Value `json:"WarningForeground,omitempty"`

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
func (r *Theme_UIColorPalette) AWSCloudFormationType() string {
	return "AWS::QuickSight::Theme.UIColorPalette"
}
