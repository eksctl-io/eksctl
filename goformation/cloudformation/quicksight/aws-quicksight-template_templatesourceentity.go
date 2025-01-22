package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Template_TemplateSourceEntity AWS CloudFormation Resource (AWS::QuickSight::Template.TemplateSourceEntity)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templatesourceentity.html
type Template_TemplateSourceEntity struct {

	// SourceAnalysis AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templatesourceentity.html#cfn-quicksight-template-templatesourceentity-sourceanalysis
	SourceAnalysis *Template_TemplateSourceAnalysis `json:"SourceAnalysis,omitempty"`

	// SourceTemplate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templatesourceentity.html#cfn-quicksight-template-templatesourceentity-sourcetemplate
	SourceTemplate *Template_TemplateSourceTemplate `json:"SourceTemplate,omitempty"`

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
func (r *Template_TemplateSourceEntity) AWSCloudFormationType() string {
	return "AWS::QuickSight::Template.TemplateSourceEntity"
}
