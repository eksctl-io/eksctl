package greengrassv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ComponentVersion_LambdaDeviceMount AWS CloudFormation Resource (AWS::GreengrassV2::ComponentVersion.LambdaDeviceMount)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdadevicemount.html
type ComponentVersion_LambdaDeviceMount struct {

	// AddGroupOwner AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdadevicemount.html#cfn-greengrassv2-componentversion-lambdadevicemount-addgroupowner
	AddGroupOwner *types.Value `json:"AddGroupOwner,omitempty"`

	// Path AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdadevicemount.html#cfn-greengrassv2-componentversion-lambdadevicemount-path
	Path *types.Value `json:"Path,omitempty"`

	// Permission AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdadevicemount.html#cfn-greengrassv2-componentversion-lambdadevicemount-permission
	Permission *types.Value `json:"Permission,omitempty"`

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
func (r *ComponentVersion_LambdaDeviceMount) AWSCloudFormationType() string {
	return "AWS::GreengrassV2::ComponentVersion.LambdaDeviceMount"
}
