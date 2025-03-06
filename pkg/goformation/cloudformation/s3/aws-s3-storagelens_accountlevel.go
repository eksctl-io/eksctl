package s3

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// StorageLens_AccountLevel AWS CloudFormation Resource (AWS::S3::StorageLens.AccountLevel)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-accountlevel.html
type StorageLens_AccountLevel struct {

	// ActivityMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-accountlevel.html#cfn-s3-storagelens-accountlevel-activitymetrics
	ActivityMetrics *StorageLens_ActivityMetrics `json:"ActivityMetrics,omitempty"`

	// AdvancedCostOptimizationMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-accountlevel.html#cfn-s3-storagelens-accountlevel-advancedcostoptimizationmetrics
	AdvancedCostOptimizationMetrics *StorageLens_AdvancedCostOptimizationMetrics `json:"AdvancedCostOptimizationMetrics,omitempty"`

	// AdvancedDataProtectionMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-accountlevel.html#cfn-s3-storagelens-accountlevel-advanceddataprotectionmetrics
	AdvancedDataProtectionMetrics *StorageLens_AdvancedDataProtectionMetrics `json:"AdvancedDataProtectionMetrics,omitempty"`

	// BucketLevel AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-accountlevel.html#cfn-s3-storagelens-accountlevel-bucketlevel
	BucketLevel *StorageLens_BucketLevel `json:"BucketLevel,omitempty"`

	// DetailedStatusCodesMetrics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-accountlevel.html#cfn-s3-storagelens-accountlevel-detailedstatuscodesmetrics
	DetailedStatusCodesMetrics *StorageLens_DetailedStatusCodesMetrics `json:"DetailedStatusCodesMetrics,omitempty"`

	// StorageLensGroupLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-accountlevel.html#cfn-s3-storagelens-accountlevel-storagelensgrouplevel
	StorageLensGroupLevel *StorageLens_StorageLensGroupLevel `json:"StorageLensGroupLevel,omitempty"`

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
func (r *StorageLens_AccountLevel) AWSCloudFormationType() string {
	return "AWS::S3::StorageLens.AccountLevel"
}
