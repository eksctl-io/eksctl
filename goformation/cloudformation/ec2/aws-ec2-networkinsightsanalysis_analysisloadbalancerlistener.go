package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// NetworkInsightsAnalysis_AnalysisLoadBalancerListener AWS CloudFormation Resource (AWS::EC2::NetworkInsightsAnalysis.AnalysisLoadBalancerListener)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysisloadbalancerlistener.html
type NetworkInsightsAnalysis_AnalysisLoadBalancerListener struct {

	// InstancePort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysisloadbalancerlistener.html#cfn-ec2-networkinsightsanalysis-analysisloadbalancerlistener-instanceport
	InstancePort *types.Value `json:"InstancePort,omitempty"`

	// LoadBalancerPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysisloadbalancerlistener.html#cfn-ec2-networkinsightsanalysis-analysisloadbalancerlistener-loadbalancerport
	LoadBalancerPort *types.Value `json:"LoadBalancerPort,omitempty"`

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
func (r *NetworkInsightsAnalysis_AnalysisLoadBalancerListener) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkInsightsAnalysis.AnalysisLoadBalancerListener"
}
