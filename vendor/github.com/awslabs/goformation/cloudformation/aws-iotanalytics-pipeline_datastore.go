package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_Datastore AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.Datastore)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-datastore.html
type AWSIoTAnalyticsPipeline_Datastore struct {

	// DatastoreName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-datastore.html#cfn-iotanalytics-pipeline-datastore-datastorename
	DatastoreName *Value `json:"DatastoreName,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-datastore.html#cfn-iotanalytics-pipeline-datastore-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_Datastore) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.Datastore"
}

func (r *AWSIoTAnalyticsPipeline_Datastore) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
