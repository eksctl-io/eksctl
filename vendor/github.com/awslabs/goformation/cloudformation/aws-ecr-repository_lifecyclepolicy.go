package cloudformation

// AWSECRRepository_LifecyclePolicy AWS CloudFormation Resource (AWS::ECR::Repository.LifecyclePolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecr-repository-lifecyclepolicy.html
type AWSECRRepository_LifecyclePolicy struct {

	// LifecyclePolicyText AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecr-repository-lifecyclepolicy.html#cfn-ecr-repository-lifecyclepolicy-lifecyclepolicytext
	LifecyclePolicyText *StringIntrinsic `json:"LifecyclePolicyText,omitempty"`

	// RegistryId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecr-repository-lifecyclepolicy.html#cfn-ecr-repository-lifecyclepolicy-registryid
	RegistryId *StringIntrinsic `json:"RegistryId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSECRRepository_LifecyclePolicy) AWSCloudFormationType() string {
	return "AWS::ECR::Repository.LifecyclePolicy"
}
