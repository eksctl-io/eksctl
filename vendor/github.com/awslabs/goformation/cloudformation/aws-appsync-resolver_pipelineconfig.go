package cloudformation

import (
	"encoding/json"
)

// AWSAppSyncResolver_PipelineConfig AWS CloudFormation Resource (AWS::AppSync::Resolver.PipelineConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-resolver-pipelineconfig.html
type AWSAppSyncResolver_PipelineConfig struct {

	// Functions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appsync-resolver-pipelineconfig.html#cfn-appsync-resolver-pipelineconfig-functions
	Functions []*Value `json:"Functions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncResolver_PipelineConfig) AWSCloudFormationType() string {
	return "AWS::AppSync::Resolver.PipelineConfig"
}

func (r *AWSAppSyncResolver_PipelineConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
