package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_SageMakerMachineLearningModelResourceData AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.SageMakerMachineLearningModelResourceData)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-sagemakermachinelearningmodelresourcedata.html
type AWSGreengrassResourceDefinitionVersion_SageMakerMachineLearningModelResourceData struct {

	// DestinationPath AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-sagemakermachinelearningmodelresourcedata.html#cfn-greengrass-resourcedefinitionversion-sagemakermachinelearningmodelresourcedata-destinationpath
	DestinationPath *Value `json:"DestinationPath,omitempty"`

	// SageMakerJobArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-sagemakermachinelearningmodelresourcedata.html#cfn-greengrass-resourcedefinitionversion-sagemakermachinelearningmodelresourcedata-sagemakerjobarn
	SageMakerJobArn *Value `json:"SageMakerJobArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_SageMakerMachineLearningModelResourceData) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.SageMakerMachineLearningModelResourceData"
}

func (r *AWSGreengrassResourceDefinitionVersion_SageMakerMachineLearningModelResourceData) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
