package lookoutmetrics

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AnomalyDetector_S3SourceConfig AWS CloudFormation Resource (AWS::LookoutMetrics::AnomalyDetector.S3SourceConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-s3sourceconfig.html
type AnomalyDetector_S3SourceConfig struct {

	// FileFormatDescriptor AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-s3sourceconfig.html#cfn-lookoutmetrics-anomalydetector-s3sourceconfig-fileformatdescriptor
	FileFormatDescriptor *AnomalyDetector_FileFormatDescriptor `json:"FileFormatDescriptor,omitempty"`

	// HistoricalDataPathList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-s3sourceconfig.html#cfn-lookoutmetrics-anomalydetector-s3sourceconfig-historicaldatapathlist
	HistoricalDataPathList *types.Value `json:"HistoricalDataPathList,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-s3sourceconfig.html#cfn-lookoutmetrics-anomalydetector-s3sourceconfig-rolearn
	RoleArn *types.Value `json:"RoleArn,omitempty"`

	// TemplatedPathList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-s3sourceconfig.html#cfn-lookoutmetrics-anomalydetector-s3sourceconfig-templatedpathlist
	TemplatedPathList *types.Value `json:"TemplatedPathList,omitempty"`

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
func (r *AnomalyDetector_S3SourceConfig) AWSCloudFormationType() string {
	return "AWS::LookoutMetrics::AnomalyDetector.S3SourceConfig"
}
