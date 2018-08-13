package cloudformation

import (
	"encoding/json"
)

// AWSBatchJobDefinition_VolumesHost AWS CloudFormation Resource (AWS::Batch::JobDefinition.VolumesHost)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-volumeshost.html
type AWSBatchJobDefinition_VolumesHost struct {

	// SourcePath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-volumeshost.html#cfn-batch-jobdefinition-volumeshost-sourcepath
	SourcePath *Value `json:"SourcePath,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchJobDefinition_VolumesHost) AWSCloudFormationType() string {
	return "AWS::Batch::JobDefinition.VolumesHost"
}

func (r *AWSBatchJobDefinition_VolumesHost) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
