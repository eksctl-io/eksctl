package wisdom

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// KnowledgeBase_RenderingConfiguration AWS CloudFormation Resource (AWS::Wisdom::KnowledgeBase.RenderingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wisdom-knowledgebase-renderingconfiguration.html
type KnowledgeBase_RenderingConfiguration struct {

	// TemplateUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wisdom-knowledgebase-renderingconfiguration.html#cfn-wisdom-knowledgebase-renderingconfiguration-templateuri
	TemplateUri *types.Value `json:"TemplateUri,omitempty"`

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
func (r *KnowledgeBase_RenderingConfiguration) AWSCloudFormationType() string {
	return "AWS::Wisdom::KnowledgeBase.RenderingConfiguration"
}
