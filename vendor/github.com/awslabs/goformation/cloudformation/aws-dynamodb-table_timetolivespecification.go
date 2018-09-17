package cloudformation

import (
	"encoding/json"
)

// AWSDynamoDBTable_TimeToLiveSpecification AWS CloudFormation Resource (AWS::DynamoDB::Table.TimeToLiveSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-timetolivespecification.html
type AWSDynamoDBTable_TimeToLiveSpecification struct {

	// AttributeName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-timetolivespecification.html#cfn-dynamodb-timetolivespecification-attributename
	AttributeName *Value `json:"AttributeName,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-timetolivespecification.html#cfn-dynamodb-timetolivespecification-enabled
	Enabled *Value `json:"Enabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDynamoDBTable_TimeToLiveSpecification) AWSCloudFormationType() string {
	return "AWS::DynamoDB::Table.TimeToLiveSpecification"
}

func (r *AWSDynamoDBTable_TimeToLiveSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
