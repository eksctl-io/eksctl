package cloudformation

import (
	"encoding/json"
)

// AlexaASKSkill_SkillPackage AWS CloudFormation Resource (Alexa::ASK::Skill.SkillPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html
type AlexaASKSkill_SkillPackage struct {

	// Overrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-overrides
	Overrides *AlexaASKSkill_Overrides `json:"Overrides,omitempty"`

	// S3Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3bucket
	S3Bucket *Value `json:"S3Bucket,omitempty"`

	// S3BucketRole AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3bucketrole
	S3BucketRole *Value `json:"S3BucketRole,omitempty"`

	// S3Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3key
	S3Key *Value `json:"S3Key,omitempty"`

	// S3ObjectVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-skillpackage.html#cfn-ask-skill-skillpackage-s3objectversion
	S3ObjectVersion *Value `json:"S3ObjectVersion,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AlexaASKSkill_SkillPackage) AWSCloudFormationType() string {
	return "Alexa::ASK::Skill.SkillPackage"
}

func (r *AlexaASKSkill_SkillPackage) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
