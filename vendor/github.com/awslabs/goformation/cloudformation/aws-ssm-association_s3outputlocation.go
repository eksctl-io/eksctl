package cloudformation

import (
	"encoding/json"
)

// AWSSSMAssociation_S3OutputLocation AWS CloudFormation Resource (AWS::SSM::Association.S3OutputLocation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-s3outputlocation.html
type AWSSSMAssociation_S3OutputLocation struct {

	// OutputS3BucketName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-s3outputlocation.html#cfn-ssm-association-s3outputlocation-outputs3bucketname
	OutputS3BucketName *Value `json:"OutputS3BucketName,omitempty"`

	// OutputS3KeyPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-s3outputlocation.html#cfn-ssm-association-s3outputlocation-outputs3keyprefix
	OutputS3KeyPrefix *Value `json:"OutputS3KeyPrefix,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMAssociation_S3OutputLocation) AWSCloudFormationType() string {
	return "AWS::SSM::Association.S3OutputLocation"
}

func (r *AWSSSMAssociation_S3OutputLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
