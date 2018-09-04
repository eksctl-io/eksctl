package cloudformation

// AWSCodePipelinePipeline_EncryptionKey AWS CloudFormation Resource (AWS::CodePipeline::Pipeline.EncryptionKey)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-artifactstore-encryptionkey.html
type AWSCodePipelinePipeline_EncryptionKey struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-artifactstore-encryptionkey.html#cfn-codepipeline-pipeline-artifactstore-encryptionkey-id
	Id *StringIntrinsic `json:"Id,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-artifactstore-encryptionkey.html#cfn-codepipeline-pipeline-artifactstore-encryptionkey-type
	Type *StringIntrinsic `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodePipelinePipeline_EncryptionKey) AWSCloudFormationType() string {
	return "AWS::CodePipeline::Pipeline.EncryptionKey"
}
