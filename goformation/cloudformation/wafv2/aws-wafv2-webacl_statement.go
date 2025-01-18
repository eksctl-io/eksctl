package wafv2

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// WebACL_Statement AWS CloudFormation Resource (AWS::WAFv2::WebACL.Statement)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html
type WebACL_Statement struct {

	// AndStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-andstatement
	AndStatement *WebACL_AndStatement `json:"AndStatement,omitempty"`

	// ByteMatchStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-bytematchstatement
	ByteMatchStatement *WebACL_ByteMatchStatement `json:"ByteMatchStatement,omitempty"`

	// GeoMatchStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-geomatchstatement
	GeoMatchStatement *WebACL_GeoMatchStatement `json:"GeoMatchStatement,omitempty"`

	// IPSetReferenceStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-ipsetreferencestatement
	IPSetReferenceStatement *WebACL_IPSetReferenceStatement `json:"IPSetReferenceStatement,omitempty"`

	// LabelMatchStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-labelmatchstatement
	LabelMatchStatement *WebACL_LabelMatchStatement `json:"LabelMatchStatement,omitempty"`

	// ManagedRuleGroupStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-managedrulegroupstatement
	ManagedRuleGroupStatement *WebACL_ManagedRuleGroupStatement `json:"ManagedRuleGroupStatement,omitempty"`

	// NotStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-notstatement
	NotStatement *WebACL_NotStatement `json:"NotStatement,omitempty"`

	// OrStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-orstatement
	OrStatement *WebACL_OrStatement `json:"OrStatement,omitempty"`

	// RateBasedStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-ratebasedstatement
	RateBasedStatement *WebACL_RateBasedStatement `json:"RateBasedStatement,omitempty"`

	// RegexPatternSetReferenceStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-regexpatternsetreferencestatement
	RegexPatternSetReferenceStatement *WebACL_RegexPatternSetReferenceStatement `json:"RegexPatternSetReferenceStatement,omitempty"`

	// RuleGroupReferenceStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-rulegroupreferencestatement
	RuleGroupReferenceStatement *WebACL_RuleGroupReferenceStatement `json:"RuleGroupReferenceStatement,omitempty"`

	// SizeConstraintStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-sizeconstraintstatement
	SizeConstraintStatement *WebACL_SizeConstraintStatement `json:"SizeConstraintStatement,omitempty"`

	// SqliMatchStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-sqlimatchstatement
	SqliMatchStatement *WebACL_SqliMatchStatement `json:"SqliMatchStatement,omitempty"`

	// XssMatchStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-statement.html#cfn-wafv2-webacl-statement-xssmatchstatement
	XssMatchStatement *WebACL_XssMatchStatement `json:"XssMatchStatement,omitempty"`

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
func (r *WebACL_Statement) AWSCloudFormationType() string {
	return "AWS::WAFv2::WebACL.Statement"
}
