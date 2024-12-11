package wisdom

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// KnowledgeBase_ServerSideEncryptionConfiguration AWS CloudFormation Resource (AWS::Wisdom::KnowledgeBase.ServerSideEncryptionConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wisdom-knowledgebase-serversideencryptionconfiguration.html
type KnowledgeBase_ServerSideEncryptionConfiguration struct {

	// KmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wisdom-knowledgebase-serversideencryptionconfiguration.html#cfn-wisdom-knowledgebase-serversideencryptionconfiguration-kmskeyid
	KmsKeyId *types.Value `json:"KmsKeyId,omitempty"`

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
func (r *KnowledgeBase_ServerSideEncryptionConfiguration) AWSCloudFormationType() string {
	return "AWS::Wisdom::KnowledgeBase.ServerSideEncryptionConfiguration"
}
