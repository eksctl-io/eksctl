package cloudformation

import (
	"encoding/json"
)

// AWSGreengrassCoreDefinitionVersion_Core AWS CloudFormation Resource (AWS::Greengrass::CoreDefinitionVersion.Core)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-coredefinitionversion-core.html
type AWSGreengrassCoreDefinitionVersion_Core struct {

	// CertificateArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-coredefinitionversion-core.html#cfn-greengrass-coredefinitionversion-core-certificatearn
	CertificateArn *Value `json:"CertificateArn,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-coredefinitionversion-core.html#cfn-greengrass-coredefinitionversion-core-id
	Id *Value `json:"Id,omitempty"`

	// SyncShadow AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-coredefinitionversion-core.html#cfn-greengrass-coredefinitionversion-core-syncshadow
	SyncShadow *Value `json:"SyncShadow,omitempty"`

	// ThingArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-coredefinitionversion-core.html#cfn-greengrass-coredefinitionversion-core-thingarn
	ThingArn *Value `json:"ThingArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassCoreDefinitionVersion_Core) AWSCloudFormationType() string {
	return "AWS::Greengrass::CoreDefinitionVersion.Core"
}

func (r *AWSGreengrassCoreDefinitionVersion_Core) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
