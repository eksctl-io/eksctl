package cloudformation

// AWSDynamoDBTable_PointInTimeRecoverySpecification AWS CloudFormation Resource (AWS::DynamoDB::Table.PointInTimeRecoverySpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-table-pointintimerecoveryspecification.html
type AWSDynamoDBTable_PointInTimeRecoverySpecification struct {

	// PointInTimeRecoveryEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-table-pointintimerecoveryspecification.html#cfn-dynamodb-table-pointintimerecoveryspecification-pointintimerecoveryenabled
	PointInTimeRecoveryEnabled bool `json:"PointInTimeRecoveryEnabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDynamoDBTable_PointInTimeRecoverySpecification) AWSCloudFormationType() string {
	return "AWS::DynamoDB::Table.PointInTimeRecoverySpecification"
}
