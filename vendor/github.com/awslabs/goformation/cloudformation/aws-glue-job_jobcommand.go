package cloudformation

import (
	"encoding/json"
)

// AWSGlueJob_JobCommand AWS CloudFormation Resource (AWS::Glue::Job.JobCommand)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-job-jobcommand.html
type AWSGlueJob_JobCommand struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-job-jobcommand.html#cfn-glue-job-jobcommand-name
	Name *Value `json:"Name,omitempty"`

	// ScriptLocation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-job-jobcommand.html#cfn-glue-job-jobcommand-scriptlocation
	ScriptLocation *Value `json:"ScriptLocation,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueJob_JobCommand) AWSCloudFormationType() string {
	return "AWS::Glue::Job.JobCommand"
}

func (r *AWSGlueJob_JobCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
