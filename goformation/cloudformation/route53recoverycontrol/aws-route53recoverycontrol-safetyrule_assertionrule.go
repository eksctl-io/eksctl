package route53recoverycontrol

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// SafetyRule_AssertionRule AWS CloudFormation Resource (AWS::Route53RecoveryControl::SafetyRule.AssertionRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoverycontrol-safetyrule-assertionrule.html
type SafetyRule_AssertionRule struct {

	// AssertedControls AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoverycontrol-safetyrule-assertionrule.html#cfn-route53recoverycontrol-safetyrule-assertionrule-assertedcontrols
	AssertedControls *types.Value `json:"AssertedControls,omitempty"`

	// WaitPeriodMs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoverycontrol-safetyrule-assertionrule.html#cfn-route53recoverycontrol-safetyrule-assertionrule-waitperiodms
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
func (r *SafetyRule_AssertionRule) AWSCloudFormationType() string {
	return "AWS::Route53RecoveryControl::SafetyRule.AssertionRule"
}
