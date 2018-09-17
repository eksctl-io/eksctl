package cloudformation

import (
	"encoding/json"
)

// AWSSSMPatchBaseline_Rule AWS CloudFormation Resource (AWS::SSM::PatchBaseline.Rule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-rule.html
type AWSSSMPatchBaseline_Rule struct {

	// ApproveAfterDays AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-rule.html#cfn-ssm-patchbaseline-rule-approveafterdays
	ApproveAfterDays *Value `json:"ApproveAfterDays,omitempty"`

	// ComplianceLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-rule.html#cfn-ssm-patchbaseline-rule-compliancelevel
	ComplianceLevel *Value `json:"ComplianceLevel,omitempty"`

	// EnableNonSecurity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-rule.html#cfn-ssm-patchbaseline-rule-enablenonsecurity
	EnableNonSecurity *Value `json:"EnableNonSecurity,omitempty"`

	// PatchFilterGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-rule.html#cfn-ssm-patchbaseline-rule-patchfiltergroup
	PatchFilterGroup *AWSSSMPatchBaseline_PatchFilterGroup `json:"PatchFilterGroup,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMPatchBaseline_Rule) AWSCloudFormationType() string {
	return "AWS::SSM::PatchBaseline.Rule"
}

func (r *AWSSSMPatchBaseline_Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
