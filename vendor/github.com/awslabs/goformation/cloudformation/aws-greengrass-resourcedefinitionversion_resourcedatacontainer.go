package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_ResourceDataContainer AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.ResourceDataContainer)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourcedatacontainer.html
type AWSGreengrassResourceDefinitionVersion_ResourceDataContainer struct {

	// LocalDeviceResourceData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourcedatacontainer.html#cfn-greengrass-resourcedefinitionversion-resourcedatacontainer-localdeviceresourcedata
	LocalDeviceResourceData *AWSGreengrassResourceDefinitionVersion_LocalDeviceResourceData `json:"LocalDeviceResourceData,omitempty"`

	// LocalVolumeResourceData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourcedatacontainer.html#cfn-greengrass-resourcedefinitionversion-resourcedatacontainer-localvolumeresourcedata
	LocalVolumeResourceData *AWSGreengrassResourceDefinitionVersion_LocalVolumeResourceData `json:"LocalVolumeResourceData,omitempty"`

	// S3MachineLearningModelResourceData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourcedatacontainer.html#cfn-greengrass-resourcedefinitionversion-resourcedatacontainer-s3machinelearningmodelresourcedata
	S3MachineLearningModelResourceData *AWSGreengrassResourceDefinitionVersion_S3MachineLearningModelResourceData `json:"S3MachineLearningModelResourceData,omitempty"`

	// SageMakerMachineLearningModelResourceData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourcedatacontainer.html#cfn-greengrass-resourcedefinitionversion-resourcedatacontainer-sagemakermachinelearningmodelresourcedata
	SageMakerMachineLearningModelResourceData *AWSGreengrassResourceDefinitionVersion_SageMakerMachineLearningModelResourceData `json:"SageMakerMachineLearningModelResourceData,omitempty"`

	// SecretsManagerSecretResourceData AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourcedatacontainer.html#cfn-greengrass-resourcedefinitionversion-resourcedatacontainer-secretsmanagersecretresourcedata
	SecretsManagerSecretResourceData *AWSGreengrassResourceDefinitionVersion_SecretsManagerSecretResourceData `json:"SecretsManagerSecretResourceData,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_ResourceDataContainer) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.ResourceDataContainer"
}

func (r *AWSGreengrassResourceDefinitionVersion_ResourceDataContainer) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
