package cloudformation

// AWSDynamoDBTable_SSESpecification AWS CloudFormation Resource (AWS::DynamoDB::Table.SSESpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-table-ssespecification.html
type AWSDynamoDBTable_SSESpecification struct {

	// SSEEnabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-table-ssespecification.html#cfn-dynamodb-table-ssespecification-sseenabled
	SSEEnabled bool `json:"SSEEnabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDynamoDBTable_SSESpecification) AWSCloudFormationType() string {
	return "AWS::DynamoDB::Table.SSESpecification"
}
