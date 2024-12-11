package greengrassv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ComponentVersion_LambdaContainerParams AWS CloudFormation Resource (AWS::GreengrassV2::ComponentVersion.LambdaContainerParams)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdacontainerparams.html
type ComponentVersion_LambdaContainerParams struct {

	// Devices AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdacontainerparams.html#cfn-greengrassv2-componentversion-lambdacontainerparams-devices
	Devices []ComponentVersion_LambdaDeviceMount `json:"Devices,omitempty"`

	// MemorySizeInKB AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdacontainerparams.html#cfn-greengrassv2-componentversion-lambdacontainerparams-memorysizeinkb
	MemorySizeInKB *types.Value `json:"MemorySizeInKB,omitempty"`

	// MountROSysfs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdacontainerparams.html#cfn-greengrassv2-componentversion-lambdacontainerparams-mountrosysfs
	MountROSysfs *types.Value `json:"MountROSysfs,omitempty"`

	// Volumes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdacontainerparams.html#cfn-greengrassv2-componentversion-lambdacontainerparams-volumes
	Volumes []ComponentVersion_LambdaVolumeMount `json:"Volumes,omitempty"`

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
func (r *ComponentVersion_LambdaContainerParams) AWSCloudFormationType() string {
	return "AWS::GreengrassV2::ComponentVersion.LambdaContainerParams"
}
