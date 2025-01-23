package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Domain_JupyterServerAppSettings AWS CloudFormation Resource (AWS::SageMaker::Domain.JupyterServerAppSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-jupyterserverappsettings.html
type Domain_JupyterServerAppSettings struct {

	// DefaultResourceSpec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-domain-jupyterserverappsettings.html#cfn-sagemaker-domain-jupyterserverappsettings-defaultresourcespec
	DefaultResourceSpec *Domain_ResourceSpec `json:"DefaultResourceSpec,omitempty"`

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
func (r *Domain_JupyterServerAppSettings) AWSCloudFormationType() string {
	return "AWS::SageMaker::Domain.JupyterServerAppSettings"
}
