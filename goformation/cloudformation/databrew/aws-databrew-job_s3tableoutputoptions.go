package databrew

import (
	"goformation/v4/cloudformation/policies"
)

// Job_S3TableOutputOptions AWS CloudFormation Resource (AWS::DataBrew::Job.S3TableOutputOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-s3tableoutputoptions.html
type Job_S3TableOutputOptions struct {

	// Location AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-s3tableoutputoptions.html#cfn-databrew-job-s3tableoutputoptions-location
	Location *Job_S3Location `json:"Location,omitempty"`

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
func (r *Job_S3TableOutputOptions) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job.S3TableOutputOptions"
}
