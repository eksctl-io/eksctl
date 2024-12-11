package wisdom

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// KnowledgeBase_AppIntegrationsConfiguration AWS CloudFormation Resource (AWS::Wisdom::KnowledgeBase.AppIntegrationsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wisdom-knowledgebase-appintegrationsconfiguration.html
type KnowledgeBase_AppIntegrationsConfiguration struct {

	// AppIntegrationArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wisdom-knowledgebase-appintegrationsconfiguration.html#cfn-wisdom-knowledgebase-appintegrationsconfiguration-appintegrationarn
	AppIntegrationArn *types.Value `json:"AppIntegrationArn,omitempty"`

	// ObjectFields AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wisdom-knowledgebase-appintegrationsconfiguration.html#cfn-wisdom-knowledgebase-appintegrationsconfiguration-objectfields
	ObjectFields *types.Value `json:"ObjectFields,omitempty"`

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
func (r *KnowledgeBase_AppIntegrationsConfiguration) AWSCloudFormationType() string {
	return "AWS::Wisdom::KnowledgeBase.AppIntegrationsConfiguration"
}
