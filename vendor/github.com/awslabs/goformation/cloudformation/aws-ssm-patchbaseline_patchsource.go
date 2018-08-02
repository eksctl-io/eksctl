package cloudformation

// AWSSSMPatchBaseline_PatchSource AWS CloudFormation Resource (AWS::SSM::PatchBaseline.PatchSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-patchsource.html
type AWSSSMPatchBaseline_PatchSource struct {

	// Configuration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-patchsource.html#cfn-ssm-patchbaseline-patchsource-configuration
	Configuration *StringIntrinsic `json:"Configuration,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-patchsource.html#cfn-ssm-patchbaseline-patchsource-name
	Name *StringIntrinsic `json:"Name,omitempty"`

	// Products AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-patchbaseline-patchsource.html#cfn-ssm-patchbaseline-patchsource-products
	Products []*StringIntrinsic `json:"Products,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMPatchBaseline_PatchSource) AWSCloudFormationType() string {
	return "AWS::SSM::PatchBaseline.PatchSource"
}
