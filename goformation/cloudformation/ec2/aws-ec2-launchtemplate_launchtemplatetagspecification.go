package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/cloudformation"
	"goformation/v4/cloudformation/policies"
)

// LaunchTemplate_LaunchTemplateTagSpecification AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.LaunchTemplateTagSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatetagspecification.html
type LaunchTemplate_LaunchTemplateTagSpecification struct {

	// ResourceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatetagspecification.html#cfn-ec2-launchtemplate-launchtemplatetagspecification-resourcetype
	ResourceType *types.Value `json:"ResourceType,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatetagspecification.html#cfn-ec2-launchtemplate-launchtemplatetagspecification-tags
	Tags []cloudformation.Tag `json:"Tags,omitempty"`

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
func (r *LaunchTemplate_LaunchTemplateTagSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.LaunchTemplateTagSpecification"
}
