package cloudformation

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StackSet_AutoDeployment AWS CloudFormation Resource (AWS::CloudFormation::StackSet.AutoDeployment)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-autodeployment.html
type StackSet_AutoDeployment struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-autodeployment.html#cfn-cloudformation-stackset-autodeployment-enabled
	Enabled *types.Value `json:"Enabled,omitempty"`

	// RetainStacksOnAccountRemoval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-stackset-autodeployment.html#cfn-cloudformation-stackset-autodeployment-retainstacksonaccountremoval
	RetainStacksOnAccountRemoval *types.Value `json:"RetainStacksOnAccountRemoval,omitempty"`

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
func (r *StackSet_AutoDeployment) AWSCloudFormationType() string {
	return "AWS::CloudFormation::StackSet.AutoDeployment"
}
