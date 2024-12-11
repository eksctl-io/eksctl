package networkfirewall

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// RuleGroup_TCPFlagField AWS CloudFormation Resource (AWS::NetworkFirewall::RuleGroup.TCPFlagField)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-tcpflagfield.html
type RuleGroup_TCPFlagField struct {

	// Flags AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-tcpflagfield.html#cfn-networkfirewall-rulegroup-tcpflagfield-flags
	Flags *types.Value `json:"Flags,omitempty"`

	// Masks AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-tcpflagfield.html#cfn-networkfirewall-rulegroup-tcpflagfield-masks
	Masks *types.Value `json:"Masks,omitempty"`

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
func (r *RuleGroup_TCPFlagField) AWSCloudFormationType() string {
	return "AWS::NetworkFirewall::RuleGroup.TCPFlagField"
}
