package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// NetworkInsightsAnalysis_AnalysisSecurityGroupRule AWS CloudFormation Resource (AWS::EC2::NetworkInsightsAnalysis.AnalysisSecurityGroupRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysissecuritygrouprule.html
type NetworkInsightsAnalysis_AnalysisSecurityGroupRule struct {

	// Cidr AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysissecuritygrouprule.html#cfn-ec2-networkinsightsanalysis-analysissecuritygrouprule-cidr
	Cidr *types.Value `json:"Cidr,omitempty"`

	// Direction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysissecuritygrouprule.html#cfn-ec2-networkinsightsanalysis-analysissecuritygrouprule-direction
	Direction *types.Value `json:"Direction,omitempty"`

	// PortRange AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysissecuritygrouprule.html#cfn-ec2-networkinsightsanalysis-analysissecuritygrouprule-portrange
	PortRange *NetworkInsightsAnalysis_PortRange `json:"PortRange,omitempty"`

	// PrefixListId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysissecuritygrouprule.html#cfn-ec2-networkinsightsanalysis-analysissecuritygrouprule-prefixlistid
	PrefixListId *types.Value `json:"PrefixListId,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysissecuritygrouprule.html#cfn-ec2-networkinsightsanalysis-analysissecuritygrouprule-protocol
	Protocol *types.Value `json:"Protocol,omitempty"`

	// SecurityGroupId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-analysissecuritygrouprule.html#cfn-ec2-networkinsightsanalysis-analysissecuritygrouprule-securitygroupid
	SecurityGroupId *types.Value `json:"SecurityGroupId,omitempty"`

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
func (r *NetworkInsightsAnalysis_AnalysisSecurityGroupRule) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkInsightsAnalysis.AnalysisSecurityGroupRule"
}
