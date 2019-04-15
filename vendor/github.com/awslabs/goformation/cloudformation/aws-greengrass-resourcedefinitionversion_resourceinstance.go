package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassResourceDefinitionVersion_ResourceInstance AWS CloudFormation Resource (AWS::Greengrass::ResourceDefinitionVersion.ResourceInstance)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourceinstance.html
type AWSGreengrassResourceDefinitionVersion_ResourceInstance struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourceinstance.html#cfn-greengrass-resourcedefinitionversion-resourceinstance-id
	Id *Value `json:"Id,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourceinstance.html#cfn-greengrass-resourcedefinitionversion-resourceinstance-name
	Name *Value `json:"Name,omitempty"`

	// ResourceDataContainer AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-resourcedefinitionversion-resourceinstance.html#cfn-greengrass-resourcedefinitionversion-resourceinstance-resourcedatacontainer
	ResourceDataContainer *AWSGreengrassResourceDefinitionVersion_ResourceDataContainer `json:"ResourceDataContainer,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassResourceDefinitionVersion_ResourceInstance) AWSCloudFormationType() string {
	return "AWS::Greengrass::ResourceDefinitionVersion.ResourceInstance"
}

func (r *AWSGreengrassResourceDefinitionVersion_ResourceInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
