package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_DynamoDBAction AWS CloudFormation Resource (AWS::IoT::TopicRule.DynamoDBAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html
type AWSIoTTopicRule_DynamoDBAction struct {

	// HashKeyField AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-hashkeyfield
	HashKeyField *Value `json:"HashKeyField,omitempty"`

	// HashKeyType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-hashkeytype
	HashKeyType *Value `json:"HashKeyType,omitempty"`

	// HashKeyValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-hashkeyvalue
	HashKeyValue *Value `json:"HashKeyValue,omitempty"`

	// PayloadField AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-payloadfield
	PayloadField *Value `json:"PayloadField,omitempty"`

	// RangeKeyField AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-rangekeyfield
	RangeKeyField *Value `json:"RangeKeyField,omitempty"`

	// RangeKeyType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-rangekeytype
	RangeKeyType *Value `json:"RangeKeyType,omitempty"`

	// RangeKeyValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-rangekeyvalue
	RangeKeyValue *Value `json:"RangeKeyValue,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// TableName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-dynamodbaction.html#cfn-iot-topicrule-dynamodbaction-tablename
	TableName *Value `json:"TableName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_DynamoDBAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.DynamoDBAction"
}

func (r *AWSIoTTopicRule_DynamoDBAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
