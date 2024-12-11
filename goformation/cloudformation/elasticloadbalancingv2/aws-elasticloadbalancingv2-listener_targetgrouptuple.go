package elasticloadbalancingv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Listener_TargetGroupTuple AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::Listener.TargetGroupTuple)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-targetgrouptuple.html
type Listener_TargetGroupTuple struct {

	// TargetGroupArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-targetgrouptuple.html#cfn-elasticloadbalancingv2-listener-targetgrouptuple-targetgrouparn
	TargetGroupArn *types.Value `json:"TargetGroupArn,omitempty"`

	// Weight AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-targetgrouptuple.html#cfn-elasticloadbalancingv2-listener-targetgrouptuple-weight
	Weight *types.Value `json:"Weight,omitempty"`

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
func (r *Listener_TargetGroupTuple) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::Listener.TargetGroupTuple"
}
