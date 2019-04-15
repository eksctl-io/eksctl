package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_SecretsManagerSecretResourceData AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.SecretsManagerSecretResourceData)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-secretsmanagersecretresourcedata.html
type AWSGreengrassResourceDefinitionVersion_SecretsManagerSecretResourceData struct {

	// ARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-secretsmanagersecretresourcedata.html#cfn-greengrass-resourcedefinitionversion-secretsmanagersecretresourcedata-arn
	ARN *Value `json:"ARN,omitempty"`

	// AdditionalStagingLabelsToDownload AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-secretsmanagersecretresourcedata.html#cfn-greengrass-resourcedefinitionversion-secretsmanagersecretresourcedata-additionalstaginglabelstodownload
	AdditionalStagingLabelsToDownload []*Value `json:"AdditionalStagingLabelsToDownload,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_SecretsManagerSecretResourceData) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.SecretsManagerSecretResourceData"
}

func (r *AWSGreengrassResourceDefinitionVersion_SecretsManagerSecretResourceData) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
