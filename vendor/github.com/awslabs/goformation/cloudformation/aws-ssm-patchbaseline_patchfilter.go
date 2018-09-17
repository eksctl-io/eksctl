package cloudformation

import (
	"encoding/json"
)

// AWSSSMPatchBaseline_PatchFilter AWS CloudFormation Resource (AWS::SSM::PatchBaseline.PatchFilter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-patchfilter.html
type AWSSSMPatchBaseline_PatchFilter struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-patchfilter.html#cfn-ssm-patchbaseline-patchfilter-key
	Key *Value `json:"Key,omitempty"`

	// Values AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-patchfilter.html#cfn-ssm-patchbaseline-patchfilter-values
	Values []*Value `json:"Values,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMPatchBaseline_PatchFilter) AWSCloudFormationType() string {
	return "AWS::SSM::PatchBaseline.PatchFilter"
}

func (r *AWSSSMPatchBaseline_PatchFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
