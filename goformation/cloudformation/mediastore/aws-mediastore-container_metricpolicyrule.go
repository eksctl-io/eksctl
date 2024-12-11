package mediastore

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Container_MetricPolicyRule AWS CloudFormation Resource (AWS::MediaStore::Container.MetricPolicyRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediastore-container-metricpolicyrule.html
type Container_MetricPolicyRule struct {

	// ObjectGroup AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediastore-container-metricpolicyrule.html#cfn-mediastore-container-metricpolicyrule-objectgroup
	ObjectGroup *types.Value `json:"ObjectGroup,omitempty"`

	// ObjectGroupName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediastore-container-metricpolicyrule.html#cfn-mediastore-container-metricpolicyrule-objectgroupname
	ObjectGroupName *types.Value `json:"ObjectGroupName,omitempty"`

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
func (r *Container_MetricPolicyRule) AWSCloudFormationType() string {
	return "AWS::MediaStore::Container.MetricPolicyRule"
}
