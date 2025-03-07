package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// EC2Fleet_BaselinePerformanceFactorsRequest AWS CloudFormation Resource (AWS::EC2::EC2Fleet.BaselinePerformanceFactorsRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-baselineperformancefactorsrequest.html
type EC2Fleet_BaselinePerformanceFactorsRequest struct {

	// Cpu AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-baselineperformancefactorsrequest.html#cfn-ec2-ec2fleet-baselineperformancefactorsrequest-cpu
	Cpu *EC2Fleet_CpuPerformanceFactorRequest `json:"Cpu,omitempty"`

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
func (r *EC2Fleet_BaselinePerformanceFactorsRequest) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.BaselinePerformanceFactorsRequest"
}
