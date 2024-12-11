package apprunner

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Service_SourceConfiguration AWS CloudFormation Resource (AWS::AppRunner::Service.SourceConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-sourceconfiguration.html
type Service_SourceConfiguration struct {

	// AuthenticationConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-sourceconfiguration.html#cfn-apprunner-service-sourceconfiguration-authenticationconfiguration
	AuthenticationConfiguration *Service_AuthenticationConfiguration `json:"AuthenticationConfiguration,omitempty"`

	// AutoDeploymentsEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-sourceconfiguration.html#cfn-apprunner-service-sourceconfiguration-autodeploymentsenabled
	AutoDeploymentsEnabled *types.Value `json:"AutoDeploymentsEnabled,omitempty"`

	// CodeRepository AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-sourceconfiguration.html#cfn-apprunner-service-sourceconfiguration-coderepository
	CodeRepository *Service_CodeRepository `json:"CodeRepository,omitempty"`

	// ImageRepository AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apprunner-service-sourceconfiguration.html#cfn-apprunner-service-sourceconfiguration-imagerepository
	ImageRepository *Service_ImageRepository `json:"ImageRepository,omitempty"`

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
func (r *Service_SourceConfiguration) AWSCloudFormationType() string {
	return "AWS::AppRunner::Service.SourceConfiguration"
}
