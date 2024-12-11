package s3objectlambda

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AccessPoint_ObjectLambdaConfiguration AWS CloudFormation Resource (AWS::S3ObjectLambda::AccessPoint.ObjectLambdaConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3objectlambda-accesspoint-objectlambdaconfiguration.html
type AccessPoint_ObjectLambdaConfiguration struct {

	// AllowedFeatures AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3objectlambda-accesspoint-objectlambdaconfiguration.html#cfn-s3objectlambda-accesspoint-objectlambdaconfiguration-allowedfeatures
	AllowedFeatures *types.Value `json:"AllowedFeatures,omitempty"`

	// CloudWatchMetricsEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3objectlambda-accesspoint-objectlambdaconfiguration.html#cfn-s3objectlambda-accesspoint-objectlambdaconfiguration-cloudwatchmetricsenabled
	CloudWatchMetricsEnabled *types.Value `json:"CloudWatchMetricsEnabled,omitempty"`

	// SupportingAccessPoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3objectlambda-accesspoint-objectlambdaconfiguration.html#cfn-s3objectlambda-accesspoint-objectlambdaconfiguration-supportingaccesspoint
	SupportingAccessPoint *types.Value `json:"SupportingAccessPoint,omitempty"`

	// TransformationConfigurations AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3objectlambda-accesspoint-objectlambdaconfiguration.html#cfn-s3objectlambda-accesspoint-objectlambdaconfiguration-transformationconfigurations
	TransformationConfigurations []AccessPoint_TransformationConfiguration `json:"TransformationConfigurations,omitempty"`

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
func (r *AccessPoint_ObjectLambdaConfiguration) AWSCloudFormationType() string {
	return "AWS::S3ObjectLambda::AccessPoint.ObjectLambdaConfiguration"
}
