package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// VerifiedAccessInstance_VerifiedAccessLogs AWS CloudFormation Resource (AWS::EC2::VerifiedAccessInstance.VerifiedAccessLogs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessinstance-verifiedaccesslogs.html
type VerifiedAccessInstance_VerifiedAccessLogs struct {

	// CloudWatchLogs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessinstance-verifiedaccesslogs.html#cfn-ec2-verifiedaccessinstance-verifiedaccesslogs-cloudwatchlogs
	CloudWatchLogs *VerifiedAccessInstance_CloudWatchLogs `json:"CloudWatchLogs,omitempty"`

	// IncludeTrustContext AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessinstance-verifiedaccesslogs.html#cfn-ec2-verifiedaccessinstance-verifiedaccesslogs-includetrustcontext
	IncludeTrustContext *types.Value `json:"IncludeTrustContext,omitempty"`

	// KinesisDataFirehose AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessinstance-verifiedaccesslogs.html#cfn-ec2-verifiedaccessinstance-verifiedaccesslogs-kinesisdatafirehose
	KinesisDataFirehose *VerifiedAccessInstance_KinesisDataFirehose `json:"KinesisDataFirehose,omitempty"`

	// LogVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessinstance-verifiedaccesslogs.html#cfn-ec2-verifiedaccessinstance-verifiedaccesslogs-logversion
	LogVersion *types.Value `json:"LogVersion,omitempty"`

	// S3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-verifiedaccessinstance-verifiedaccesslogs.html#cfn-ec2-verifiedaccessinstance-verifiedaccesslogs-s3
	S3 *VerifiedAccessInstance_S3 `json:"S3,omitempty"`

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
func (r *VerifiedAccessInstance_VerifiedAccessLogs) AWSCloudFormationType() string {
	return "AWS::EC2::VerifiedAccessInstance.VerifiedAccessLogs"
}
