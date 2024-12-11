package databrew

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Job_StatisticsConfiguration AWS CloudFormation Resource (AWS::DataBrew::Job.StatisticsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-statisticsconfiguration.html
type Job_StatisticsConfiguration struct {

	// IncludedStatistics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-statisticsconfiguration.html#cfn-databrew-job-statisticsconfiguration-includedstatistics
	IncludedStatistics *types.Value `json:"IncludedStatistics,omitempty"`

	// Overrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-statisticsconfiguration.html#cfn-databrew-job-statisticsconfiguration-overrides
	Overrides []Job_StatisticOverride `json:"Overrides,omitempty"`

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
func (r *Job_StatisticsConfiguration) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job.StatisticsConfiguration"
}
