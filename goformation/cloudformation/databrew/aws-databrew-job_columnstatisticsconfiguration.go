package databrew

import (
	"goformation/v4/cloudformation/policies"
)

// Job_ColumnStatisticsConfiguration AWS CloudFormation Resource (AWS::DataBrew::Job.ColumnStatisticsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-columnstatisticsconfiguration.html
type Job_ColumnStatisticsConfiguration struct {

	// Selectors AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-columnstatisticsconfiguration.html#cfn-databrew-job-columnstatisticsconfiguration-selectors
	Selectors []Job_ColumnSelector `json:"Selectors,omitempty"`

	// Statistics AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-columnstatisticsconfiguration.html#cfn-databrew-job-columnstatisticsconfiguration-statistics
	Statistics *Job_StatisticsConfiguration `json:"Statistics,omitempty"`

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
func (r *Job_ColumnStatisticsConfiguration) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job.ColumnStatisticsConfiguration"
}
