package lambda

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CodeSigningConfig_AllowedPublishers AWS CloudFormation Resource (AWS::Lambda::CodeSigningConfig.AllowedPublishers)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-codesigningconfig-allowedpublishers.html
type CodeSigningConfig_AllowedPublishers struct {

	// SigningProfileVersionArns AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-codesigningconfig-allowedpublishers.html#cfn-lambda-codesigningconfig-allowedpublishers-signingprofileversionarns
	SigningProfileVersionArns *types.Value `json:"SigningProfileVersionArns,omitempty"`

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
func (r *CodeSigningConfig_AllowedPublishers) AWSCloudFormationType() string {
	return "AWS::Lambda::CodeSigningConfig.AllowedPublishers"
}
