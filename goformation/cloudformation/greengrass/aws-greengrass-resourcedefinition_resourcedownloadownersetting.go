package greengrass

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResourceDefinition_ResourceDownloadOwnerSetting AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinition.ResourceDownloadOwnerSetting)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinition-resourcedownloadownersetting.html
type ResourceDefinition_ResourceDownloadOwnerSetting struct {

	// GroupOwner AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinition-resourcedownloadownersetting.html#cfn-greengrass-resourcedefinition-resourcedownloadownersetting-groupowner
	GroupOwner *types.Value `json:"GroupOwner,omitempty"`

	// GroupPermission AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinition-resourcedownloadownersetting.html#cfn-greengrass-resourcedefinition-resourcedownloadownersetting-grouppermission
	GroupPermission *types.Value `json:"GroupPermission,omitempty"`

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
func (r *ResourceDefinition_ResourceDownloadOwnerSetting) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinition.ResourceDownloadOwnerSetting"
}
