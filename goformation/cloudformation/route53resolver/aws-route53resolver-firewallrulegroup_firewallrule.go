package route53resolver

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// FirewallRuleGroup_FirewallRule AWS CloudFormation Resource (AWS::Route53Resolver::FirewallRuleGroup.FirewallRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html
type FirewallRuleGroup_FirewallRule struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html#cfn-route53resolver-firewallrulegroup-firewallrule-action
	Action *types.Value `json:"Action,omitempty"`

	// BlockOverrideDnsType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html#cfn-route53resolver-firewallrulegroup-firewallrule-blockoverridednstype
	BlockOverrideDnsType *types.Value `json:"BlockOverrideDnsType,omitempty"`

	// BlockOverrideDomain AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html#cfn-route53resolver-firewallrulegroup-firewallrule-blockoverridedomain
	BlockOverrideDomain *types.Value `json:"BlockOverrideDomain,omitempty"`

	// BlockOverrideTtl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html#cfn-route53resolver-firewallrulegroup-firewallrule-blockoverridettl
	BlockOverrideTtl *types.Value `json:"BlockOverrideTtl,omitempty"`

	// BlockResponse AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html#cfn-route53resolver-firewallrulegroup-firewallrule-blockresponse
	BlockResponse *types.Value `json:"BlockResponse,omitempty"`

	// FirewallDomainListId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html#cfn-route53resolver-firewallrulegroup-firewallrule-firewalldomainlistid
	FirewallDomainListId *types.Value `json:"FirewallDomainListId,omitempty"`

	// Priority AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-firewallrulegroup-firewallrule.html#cfn-route53resolver-firewallrulegroup-firewallrule-priority
	Priority *types.Value `json:"Priority"`

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
func (r *FirewallRuleGroup_FirewallRule) AWSCloudFormationType() string {
	return "AWS::Route53Resolver::FirewallRuleGroup.FirewallRule"
}
