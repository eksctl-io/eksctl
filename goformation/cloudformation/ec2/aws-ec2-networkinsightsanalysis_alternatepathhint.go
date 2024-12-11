package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// NetworkInsightsAnalysis_AlternatePathHint AWS CloudFormation Resource (AWS::EC2::NetworkInsightsAnalysis.AlternatePathHint)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-alternatepathhint.html
type NetworkInsightsAnalysis_AlternatePathHint struct {

	// ComponentArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-alternatepathhint.html#cfn-ec2-networkinsightsanalysis-alternatepathhint-componentarn
	ComponentArn *types.Value `json:"ComponentArn,omitempty"`

	// ComponentId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-alternatepathhint.html#cfn-ec2-networkinsightsanalysis-alternatepathhint-componentid
	ComponentId *types.Value `json:"ComponentId,omitempty"`

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
func (r *NetworkInsightsAnalysis_AlternatePathHint) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkInsightsAnalysis.AlternatePathHint"
}
