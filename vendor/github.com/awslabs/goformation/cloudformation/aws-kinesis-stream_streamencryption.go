package cloudformation

// AWSKinesisStream_StreamEncryption AWS CloudFormation Resource (AWS::Kinesis::Stream.StreamEncryption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesis-stream-streamencryption.html
type AWSKinesisStream_StreamEncryption struct {

	// EncryptionType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesis-stream-streamencryption.html#cfn-kinesis-stream-streamencryption-encryptiontype
	EncryptionType string `json:"EncryptionType,omitempty"`

	// KeyId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesis-stream-streamencryption.html#cfn-kinesis-stream-streamencryption-keyid
	KeyId string `json:"KeyId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisStream_StreamEncryption) AWSCloudFormationType() string {
	return "AWS::Kinesis::Stream.StreamEncryption"
}
