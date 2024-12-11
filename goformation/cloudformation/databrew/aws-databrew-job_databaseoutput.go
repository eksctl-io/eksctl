package databrew

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Job_DatabaseOutput AWS CloudFormation Resource (AWS::DataBrew::Job.DatabaseOutput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-databaseoutput.html
type Job_DatabaseOutput struct {

	// DatabaseOptions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-databaseoutput.html#cfn-databrew-job-databaseoutput-databaseoptions
	DatabaseOptions *Job_DatabaseTableOutputOptions `json:"DatabaseOptions,omitempty"`

	// DatabaseOutputMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-databaseoutput.html#cfn-databrew-job-databaseoutput-databaseoutputmode
	DatabaseOutputMode *types.Value `json:"DatabaseOutputMode,omitempty"`

	// GlueConnectionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-databaseoutput.html#cfn-databrew-job-databaseoutput-glueconnectionname
	GlueConnectionName *types.Value `json:"GlueConnectionName,omitempty"`

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
func (r *Job_DatabaseOutput) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job.DatabaseOutput"
}
