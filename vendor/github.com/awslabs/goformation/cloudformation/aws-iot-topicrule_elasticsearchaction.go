package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_ElasticsearchAction AWS CloudFormation Resource (AWS::IoT::TopicRule.ElasticsearchAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-elasticsearchaction.html
type AWSIoTTopicRule_ElasticsearchAction struct {

	// Endpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-elasticsearchaction.html#cfn-iot-topicrule-elasticsearchaction-endpoint
	Endpoint *Value `json:"Endpoint,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-elasticsearchaction.html#cfn-iot-topicrule-elasticsearchaction-id
	Id *Value `json:"Id,omitempty"`

	// Index AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-elasticsearchaction.html#cfn-iot-topicrule-elasticsearchaction-index
	Index *Value `json:"Index,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-elasticsearchaction.html#cfn-iot-topicrule-elasticsearchaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-elasticsearchaction.html#cfn-iot-topicrule-elasticsearchaction-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_ElasticsearchAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.ElasticsearchAction"
}

func (r *AWSIoTTopicRule_ElasticsearchAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
