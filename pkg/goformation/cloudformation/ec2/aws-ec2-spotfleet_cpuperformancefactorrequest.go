package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// SpotFleet_CpuPerformanceFactorRequest AWS CloudFormation Resource (AWS::EC2::SpotFleet.CpuPerformanceFactorRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-cpuperformancefactorrequest.html
type SpotFleet_CpuPerformanceFactorRequest struct {

	// References AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-cpuperformancefactorrequest.html#cfn-ec2-spotfleet-cpuperformancefactorrequest-references
	References []SpotFleet_PerformanceFactorReferenceRequest `json:"References,omitempty"`

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
func (r *SpotFleet_CpuPerformanceFactorRequest) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.CpuPerformanceFactorRequest"
}
