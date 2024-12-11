package route53recoverycontrol

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SafetyRule_GatingRule AWS CloudFormation Resource (AWS::Route53RecoveryControl::SafetyRule.GatingRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoverycontrol-safetyrule-gatingrule.html
type SafetyRule_GatingRule struct {

	// GatingControls AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoverycontrol-safetyrule-gatingrule.html#cfn-route53recoverycontrol-safetyrule-gatingrule-gatingcontrols
	GatingControls *types.Value `json:"GatingControls,omitempty"`

	// TargetControls AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoverycontrol-safetyrule-gatingrule.html#cfn-route53recoverycontrol-safetyrule-gatingrule-targetcontrols
	TargetControls *types.Value `json:"TargetControls,omitempty"`

	// WaitPeriodMs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoverycontrol-safetyrule-gatingrule.html#cfn-route53recoverycontrol-safetyrule-gatingrule-waitperiodms
	WaitPeriodMs *types.Value `json:"WaitPeriodMs"`

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
func (r *SafetyRule_GatingRule) AWSCloudFormationType() string {
	return "AWS::Route53RecoveryControl::SafetyRule.GatingRule"
}
