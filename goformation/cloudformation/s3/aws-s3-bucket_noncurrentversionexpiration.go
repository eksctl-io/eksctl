package s3

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Bucket_NoncurrentVersionExpiration AWS CloudFormation Resource (AWS::S3::Bucket.NoncurrentVersionExpiration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-lifecycleconfig-rule-noncurrentversionexpiration.html
type Bucket_NoncurrentVersionExpiration struct {

	// NewerNoncurrentVersions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-lifecycleconfig-rule-noncurrentversionexpiration.html#cfn-s3-bucket-lifecycleconfig-rule-noncurrentversionexpiration-newernoncurrentversions
	NewerNoncurrentVersions *types.Value `json:"NewerNoncurrentVersions,omitempty"`

	// NoncurrentDays AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-lifecycleconfig-rule-noncurrentversionexpiration.html#cfn-s3-bucket-lifecycleconfig-rule-noncurrentversionexpiration-noncurrentdays
	NoncurrentDays *types.Value `json:"NoncurrentDays"`

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
func (r *Bucket_NoncurrentVersionExpiration) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.NoncurrentVersionExpiration"
}
