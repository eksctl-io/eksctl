package s3

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// StorageLens_BucketLevel AWS CloudFormation Resource (AWS::S3::StorageLens.BucketLevel)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-bucketlevel.html
type StorageLens_BucketLevel struct {

	// ActivityMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-bucketlevel.html#cfn-s3-storagelens-bucketlevel-activitymetrics
	ActivityMetrics *StorageLens_ActivityMetrics `json:"ActivityMetrics,omitempty"`

	// AdvancedCostOptimizationMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-bucketlevel.html#cfn-s3-storagelens-bucketlevel-advancedcostoptimizationmetrics
	AdvancedCostOptimizationMetrics *StorageLens_AdvancedCostOptimizationMetrics `json:"AdvancedCostOptimizationMetrics,omitempty"`

	// AdvancedDataProtectionMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-bucketlevel.html#cfn-s3-storagelens-bucketlevel-advanceddataprotectionmetrics
	AdvancedDataProtectionMetrics *StorageLens_AdvancedDataProtectionMetrics `json:"AdvancedDataProtectionMetrics,omitempty"`

	// DetailedStatusCodesMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-bucketlevel.html#cfn-s3-storagelens-bucketlevel-detailedstatuscodesmetrics
	DetailedStatusCodesMetrics *StorageLens_DetailedStatusCodesMetrics `json:"DetailedStatusCodesMetrics,omitempty"`

	// PrefixLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-bucketlevel.html#cfn-s3-storagelens-bucketlevel-prefixlevel
	PrefixLevel *StorageLens_PrefixLevel `json:"PrefixLevel,omitempty"`

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
func (r *StorageLens_BucketLevel) AWSCloudFormationType() string {
	return "AWS::S3::StorageLens.BucketLevel"
}
