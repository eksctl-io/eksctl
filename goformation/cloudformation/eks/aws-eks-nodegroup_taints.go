package eks

import "goformation/v4/cloudformation/types"

// Nodegroup_Taints AWS CloudFormation Resource (AWS::EKS::Nodegroup.Taints)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-nodegroup-taint.html
type Nodegroup_Taints struct {
	// TODO documentation
	Key    *types.Value `json:"Key,omitempty"`
	Value  *types.Value `json:"Value,omitempty"`
	Effect *types.Value `json:"Effect,omitempty"`
}
