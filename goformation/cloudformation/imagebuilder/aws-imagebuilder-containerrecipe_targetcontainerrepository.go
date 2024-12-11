package imagebuilder

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ContainerRecipe_TargetContainerRepository AWS CloudFormation Resource (AWS::ImageBuilder::ContainerRecipe.TargetContainerRepository)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-containerrecipe-targetcontainerrepository.html
type ContainerRecipe_TargetContainerRepository struct {

	// RepositoryName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-containerrecipe-targetcontainerrepository.html#cfn-imagebuilder-containerrecipe-targetcontainerrepository-repositoryname
	RepositoryName *types.Value `json:"RepositoryName,omitempty"`

	// Service AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-containerrecipe-targetcontainerrepository.html#cfn-imagebuilder-containerrecipe-targetcontainerrepository-service
	Service *types.Value `json:"Service,omitempty"`

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
func (r *ContainerRecipe_TargetContainerRepository) AWSCloudFormationType() string {
	return "AWS::ImageBuilder::ContainerRecipe.TargetContainerRepository"
}
