package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Integration_TriggerConfig AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.TriggerConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-triggerconfig.html
type Integration_TriggerConfig struct {

	// TriggerProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-triggerconfig.html#cfn-customerprofiles-integration-triggerconfig-triggerproperties
	TriggerProperties *Integration_TriggerProperties `json:"TriggerProperties,omitempty"`

	// TriggerType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-triggerconfig.html#cfn-customerprofiles-integration-triggerconfig-triggertype
	TriggerType *types.Value `json:"TriggerType,omitempty"`

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
func (r *Integration_TriggerConfig) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.TriggerConfig"
}
