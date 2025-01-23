package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_S3Parameters AWS CloudFormation Resource (AWS::QuickSight::DataSource.S3Parameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-s3parameters.html
type DataSource_S3Parameters struct {

	// ManifestFileLocation AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-s3parameters.html#cfn-quicksight-datasource-s3parameters-manifestfilelocation
	ManifestFileLocation *DataSource_ManifestFileLocation `json:"ManifestFileLocation,omitempty"`

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
func (r *DataSource_S3Parameters) AWSCloudFormationType() string {
	return "AWS::QuickSight::DataSource.S3Parameters"
}
