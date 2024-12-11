package databrew

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Job_DatabaseTableOutputOptions AWS CloudFormation Resource (AWS::DataBrew::Job.DatabaseTableOutputOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-databasetableoutputoptions.html
type Job_DatabaseTableOutputOptions struct {

	// TableName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-databasetableoutputoptions.html#cfn-databrew-job-databasetableoutputoptions-tablename
	TableName *types.Value `json:"TableName,omitempty"`

	// TempDirectory AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-databasetableoutputoptions.html#cfn-databrew-job-databasetableoutputoptions-tempdirectory
	TempDirectory *Job_S3Location `json:"TempDirectory,omitempty"`

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
func (r *Job_DatabaseTableOutputOptions) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job.DatabaseTableOutputOptions"
}
