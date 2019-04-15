package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_S3MachineLearningModelResourceData AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.S3MachineLearningModelResourceData)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-s3machinelearningmodelresourcedata.html
type AWSGreengrassResourceDefinitionVersion_S3MachineLearningModelResourceData struct {

	// DestinationPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-s3machinelearningmodelresourcedata.html#cfn-greengrass-resourcedefinitionversion-s3machinelearningmodelresourcedata-destinationpath
	DestinationPath *Value `json:"DestinationPath,omitempty"`

	// S3Uri AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-s3machinelearningmodelresourcedata.html#cfn-greengrass-resourcedefinitionversion-s3machinelearningmodelresourcedata-s3uri
	S3Uri *Value `json:"S3Uri,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_S3MachineLearningModelResourceData) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.S3MachineLearningModelResourceData"
}

func (r *AWSGreengrassResourceDefinitionVersion_S3MachineLearningModelResourceData) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
