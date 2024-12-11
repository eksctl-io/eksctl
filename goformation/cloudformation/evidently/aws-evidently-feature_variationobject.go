package evidently

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Feature_VariationObject AWS CloudFormation Resource (AWS::Evidently::Feature.VariationObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-feature-variationobject.html
type Feature_VariationObject struct {

	// BooleanValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-feature-variationobject.html#cfn-evidently-feature-variationobject-booleanvalue
	BooleanValue *types.Value `json:"BooleanValue,omitempty"`

	// DoubleValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-feature-variationobject.html#cfn-evidently-feature-variationobject-doublevalue
	DoubleValue *types.Value `json:"DoubleValue,omitempty"`

	// LongValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-feature-variationobject.html#cfn-evidently-feature-variationobject-longvalue
	LongValue *types.Value `json:"LongValue,omitempty"`

	// StringValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-feature-variationobject.html#cfn-evidently-feature-variationobject-stringvalue
	StringValue *types.Value `json:"StringValue,omitempty"`

	// VariationName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-feature-variationobject.html#cfn-evidently-feature-variationobject-variationname
	VariationName *types.Value `json:"VariationName,omitempty"`

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
func (r *Feature_VariationObject) AWSCloudFormationType() string {
	return "AWS::Evidently::Feature.VariationObject"
}
