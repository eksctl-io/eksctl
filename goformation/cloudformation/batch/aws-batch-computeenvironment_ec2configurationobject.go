package batch

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ComputeEnvironment_Ec2ConfigurationObject AWS CloudFormation Resource (AWS::Batch::ComputeEnvironment.Ec2ConfigurationObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-computeenvironment-ec2configurationobject.html
type ComputeEnvironment_Ec2ConfigurationObject struct {

	// ImageIdOverride AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-computeenvironment-ec2configurationobject.html#cfn-batch-computeenvironment-ec2configurationobject-imageidoverride
	ImageIdOverride *types.Value `json:"ImageIdOverride,omitempty"`

	// ImageType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-computeenvironment-ec2configurationobject.html#cfn-batch-computeenvironment-ec2configurationobject-imagetype
	ImageType *types.Value `json:"ImageType,omitempty"`

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
func (r *ComputeEnvironment_Ec2ConfigurationObject) AWSCloudFormationType() string {
	return "AWS::Batch::ComputeEnvironment.Ec2ConfigurationObject"
}
