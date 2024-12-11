package databrew

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Job_EntityDetectorConfiguration AWS CloudFormation Resource (AWS::DataBrew::Job.EntityDetectorConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-entitydetectorconfiguration.html
type Job_EntityDetectorConfiguration struct {

	// AllowedStatistics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-entitydetectorconfiguration.html#cfn-databrew-job-entitydetectorconfiguration-allowedstatistics
	AllowedStatistics *Job_AllowedStatistics `json:"AllowedStatistics,omitempty"`

	// EntityTypes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-job-entitydetectorconfiguration.html#cfn-databrew-job-entitydetectorconfiguration-entitytypes
	EntityTypes *types.Value `json:"EntityTypes,omitempty"`

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
func (r *Job_EntityDetectorConfiguration) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job.EntityDetectorConfiguration"
}
