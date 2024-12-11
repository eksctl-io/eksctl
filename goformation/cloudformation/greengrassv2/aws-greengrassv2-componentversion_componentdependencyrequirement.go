package greengrassv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ComponentVersion_ComponentDependencyRequirement AWS CloudFormation Resource (AWS::GreengrassV2::ComponentVersion.ComponentDependencyRequirement)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-componentdependencyrequirement.html
type ComponentVersion_ComponentDependencyRequirement struct {

	// DependencyType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-componentdependencyrequirement.html#cfn-greengrassv2-componentversion-componentdependencyrequirement-dependencytype
	DependencyType *types.Value `json:"DependencyType,omitempty"`

	// VersionRequirement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-componentdependencyrequirement.html#cfn-greengrassv2-componentversion-componentdependencyrequirement-versionrequirement
	VersionRequirement *types.Value `json:"VersionRequirement,omitempty"`

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
func (r *ComponentVersion_ComponentDependencyRequirement) AWSCloudFormationType() string {
	return "AWS::GreengrassV2::ComponentVersion.ComponentDependencyRequirement"
}
