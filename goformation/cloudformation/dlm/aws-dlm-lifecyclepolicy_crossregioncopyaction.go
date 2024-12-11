package dlm

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LifecyclePolicy_CrossRegionCopyAction AWS CloudFormation Resource (AWS::DLM::LifecyclePolicy.CrossRegionCopyAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-crossregioncopyaction.html
type LifecyclePolicy_CrossRegionCopyAction struct {

	// EncryptionConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-crossregioncopyaction.html#cfn-dlm-lifecyclepolicy-crossregioncopyaction-encryptionconfiguration
	EncryptionConfiguration *LifecyclePolicy_EncryptionConfiguration `json:"EncryptionConfiguration,omitempty"`

	// RetainRule AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-crossregioncopyaction.html#cfn-dlm-lifecyclepolicy-crossregioncopyaction-retainrule
	RetainRule *LifecyclePolicy_CrossRegionCopyRetainRule `json:"RetainRule,omitempty"`

	// Target AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-crossregioncopyaction.html#cfn-dlm-lifecyclepolicy-crossregioncopyaction-target
	Target *types.Value `json:"Target,omitempty"`

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
func (r *LifecyclePolicy_CrossRegionCopyAction) AWSCloudFormationType() string {
	return "AWS::DLM::LifecyclePolicy.CrossRegionCopyAction"
}
