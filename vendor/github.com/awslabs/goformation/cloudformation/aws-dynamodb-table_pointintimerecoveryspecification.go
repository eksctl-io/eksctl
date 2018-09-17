package cloudformation

import (
	"encoding/json"
)

// AWSDynamoDBTable_PointInTimeRecoverySpecification AWS CloudFormation Resource (AWS::DynamoDB::Table.PointInTimeRecoverySpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-table-pointintimerecoveryspecification.html
type AWSDynamoDBTable_PointInTimeRecoverySpecification struct {

	// PointInTimeRecoveryEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-table-pointintimerecoveryspecification.html#cfn-dynamodb-table-pointintimerecoveryspecification-pointintimerecoveryenabled
	PointInTimeRecoveryEnabled *Value `json:"PointInTimeRecoveryEnabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDynamoDBTable_PointInTimeRecoverySpecification) AWSCloudFormationType() string {
	return "AWS::DynamoDB::Table.PointInTimeRecoverySpecification"
}

func (r *AWSDynamoDBTable_PointInTimeRecoverySpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
