package elasticache

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CacheCluster_KinesisFirehoseDestinationDetails AWS CloudFormation Resource (AWS::ElastiCache::CacheCluster.KinesisFirehoseDestinationDetails)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-cachecluster-kinesisfirehosedestinationdetails.html
type CacheCluster_KinesisFirehoseDestinationDetails struct {

	// DeliveryStream AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-cachecluster-kinesisfirehosedestinationdetails.html#cfn-elasticache-cachecluster-kinesisfirehosedestinationdetails-deliverystream
	DeliveryStream *types.Value `json:"DeliveryStream,omitempty"`

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
func (r *CacheCluster_KinesisFirehoseDestinationDetails) AWSCloudFormationType() string {
	return "AWS::ElastiCache::CacheCluster.KinesisFirehoseDestinationDetails"
}
