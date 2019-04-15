package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassFunctionDefinition_ResourceAccessPolicy AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.ResourceAccessPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-resourceaccesspolicy.html
type AWSGreengrassFunctionDefinition_ResourceAccessPolicy struct {

	// Permission AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-resourceaccesspolicy.html#cfn-greengrass-functiondefinition-resourceaccesspolicy-permission
	Permission *Value `json:"Permission,omitempty"`

	// ResourceId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-resourceaccesspolicy.html#cfn-greengrass-functiondefinition-resourceaccesspolicy-resourceid
	ResourceId *Value `json:"ResourceId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition_ResourceAccessPolicy) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.ResourceAccessPolicy"
}

func (r *AWSGreengrassFunctionDefinition_ResourceAccessPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
