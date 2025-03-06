package autoscaling

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// AutoScalingGroup_BaselinePerformanceFactorsRequest AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.BaselinePerformanceFactorsRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-baselineperformancefactorsrequest.html
type AutoScalingGroup_BaselinePerformanceFactorsRequest struct {

	// Cpu AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-baselineperformancefactorsrequest.html#cfn-autoscaling-autoscalinggroup-baselineperformancefactorsrequest-cpu
	Cpu *AutoScalingGroup_CpuPerformanceFactorRequest `json:"Cpu,omitempty"`

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
func (r *AutoScalingGroup_BaselinePerformanceFactorsRequest) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.BaselinePerformanceFactorsRequest"
}
