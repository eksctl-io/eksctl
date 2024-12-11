package ask

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Skill_SkillPackage AWS CloudFormation Resource (Alexa::ASK::Skill.SkillPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html
type Skill_SkillPackage struct {

	// Overrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-overrides
	Overrides *Skill_Overrides `json:"Overrides,omitempty"`

	// S3Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3bucket
	S3Bucket *types.Value `json:"S3Bucket,omitempty"`

	// S3BucketRole AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3bucketrole
	S3BucketRole *types.Value `json:"S3BucketRole,omitempty"`

	// S3Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3key
	S3Key *types.Value `json:"S3Key,omitempty"`

	// S3ObjectVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3objectversion
	S3ObjectVersion *types.Value `json:"S3ObjectVersion,omitempty"`

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
func (r *Skill_SkillPackage) AWSCloudFormationType() string {
	return "Alexa::ASK::Skill.SkillPackage"
}
