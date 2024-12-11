package mediastore

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Container_MetricPolicy AWS CloudFormation Resource (AWS::MediaStore::Container.MetricPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediastore-container-metricpolicy.html
type Container_MetricPolicy struct {

	// ContainerLevelMetrics AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediastore-container-metricpolicy.html#cfn-mediastore-container-metricpolicy-containerlevelmetrics
	ContainerLevelMetrics *types.Value `json:"ContainerLevelMetrics,omitempty"`

	// MetricPolicyRules AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediastore-container-metricpolicy.html#cfn-mediastore-container-metricpolicy-metricpolicyrules
	MetricPolicyRules []Container_MetricPolicyRule `json:"MetricPolicyRules,omitempty"`

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
func (r *Container_MetricPolicy) AWSCloudFormationType() string {
	return "AWS::MediaStore::Container.MetricPolicy"
}
