package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Integration_FlowDefinition AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.FlowDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-flowdefinition.html
type Integration_FlowDefinition struct {

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-flowdefinition.html#cfn-customerprofiles-integration-flowdefinition-description
	Description *types.Value `json:"Description,omitempty"`

	// FlowName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-flowdefinition.html#cfn-customerprofiles-integration-flowdefinition-flowname
	FlowName *types.Value `json:"FlowName,omitempty"`

	// KmsArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-flowdefinition.html#cfn-customerprofiles-integration-flowdefinition-kmsarn
	KmsArn *types.Value `json:"KmsArn,omitempty"`

	// SourceFlowConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-flowdefinition.html#cfn-customerprofiles-integration-flowdefinition-sourceflowconfig
	SourceFlowConfig *Integration_SourceFlowConfig `json:"SourceFlowConfig,omitempty"`

	// Tasks AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-flowdefinition.html#cfn-customerprofiles-integration-flowdefinition-tasks
	Tasks []Integration_Task `json:"Tasks,omitempty"`

	// TriggerConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-flowdefinition.html#cfn-customerprofiles-integration-flowdefinition-triggerconfig
	TriggerConfig *Integration_TriggerConfig `json:"TriggerConfig,omitempty"`

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
func (r *Integration_FlowDefinition) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.FlowDefinition"
}
