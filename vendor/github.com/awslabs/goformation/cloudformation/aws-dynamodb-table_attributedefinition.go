package cloudformation

// AWSDynamoDBTable_AttributeDefinition AWS CloudFormation Resource (AWS::DynamoDB::Table.AttributeDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-attributedef.html
type AWSDynamoDBTable_AttributeDefinition struct {

	// AttributeName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-attributedef.html#cfn-dynamodb-attributedef-attributename
	AttributeName *StringIntrinsic `json:"AttributeName,omitempty"`

	// AttributeType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-attributedef.html#cfn-dynamodb-attributedef-attributename-attributetype
	AttributeType *StringIntrinsic `json:"AttributeType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDynamoDBTable_AttributeDefinition) AWSCloudFormationType() string {
	return "AWS::DynamoDB::Table.AttributeDefinition"
}
