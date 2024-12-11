package ssm

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Association_S3OutputLocation AWS CloudFormation Resource (AWS::SSM::Association.S3OutputLocation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-s3outputlocation.html
type Association_S3OutputLocation struct {

	// OutputS3BucketName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-s3outputlocation.html#cfn-ssm-association-s3outputlocation-outputs3bucketname
	OutputS3BucketName *types.Value `json:"OutputS3BucketName,omitempty"`

	// OutputS3KeyPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-s3outputlocation.html#cfn-ssm-association-s3outputlocation-outputs3keyprefix
	OutputS3KeyPrefix *types.Value `json:"OutputS3KeyPrefix,omitempty"`

	// OutputS3Region AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-s3outputlocation.html#cfn-ssm-association-s3outputlocation-outputs3region
	OutputS3Region *types.Value `json:"OutputS3Region,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Association_S3OutputLocation) AWSCloudFormationType() string {
	return "AWS::SSM::Association.S3OutputLocation"
}
