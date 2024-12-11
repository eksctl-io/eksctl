package databrew

import (
	"goformation/v4/cloudformation/policies"
)

// Job_ProfileConfiguration AWS CloudFormation Resource (AWS::DataBrew::Job.ProfileConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-profileconfiguration.html
type Job_ProfileConfiguration struct {

	// ColumnStatisticsConfigurations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-profileconfiguration.html#cfn-databrew-job-profileconfiguration-columnstatisticsconfigurations
	ColumnStatisticsConfigurations []Job_ColumnStatisticsConfiguration `json:"ColumnStatisticsConfigurations,omitempty"`

	// DatasetStatisticsConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-profileconfiguration.html#cfn-databrew-job-profileconfiguration-datasetstatisticsconfiguration
	DatasetStatisticsConfiguration *Job_StatisticsConfiguration `json:"DatasetStatisticsConfiguration,omitempty"`

	// EntityDetectorConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-profileconfiguration.html#cfn-databrew-job-profileconfiguration-entitydetectorconfiguration
	EntityDetectorConfiguration *Job_EntityDetectorConfiguration `json:"EntityDetectorConfiguration,omitempty"`

	// ProfileColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-profileconfiguration.html#cfn-databrew-job-profileconfiguration-profilecolumns
	ProfileColumns []Job_ColumnSelector `json:"ProfileColumns,omitempty"`

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
func (r *Job_ProfileConfiguration) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job.ProfileConfiguration"
}
