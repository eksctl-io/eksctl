package apprunner

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Service_ImageRepository AWS CloudFormation Resource (AWS::AppRunner::Service.ImageRepository)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-imagerepository.html
type Service_ImageRepository struct {

	// ImageConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-imagerepository.html#cfn-apprunner-service-imagerepository-imageconfiguration
	ImageConfiguration *Service_ImageConfiguration `json:"ImageConfiguration,omitempty"`

	// ImageIdentifier AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-imagerepository.html#cfn-apprunner-service-imagerepository-imageidentifier
	ImageIdentifier *types.Value `json:"ImageIdentifier,omitempty"`

	// ImageRepositoryType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-imagerepository.html#cfn-apprunner-service-imagerepository-imagerepositorytype
	ImageRepositoryType *types.Value `json:"ImageRepositoryType,omitempty"`

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
func (r *Service_ImageRepository) AWSCloudFormationType() string {
	return "AWS::AppRunner::Service.ImageRepository"
}
