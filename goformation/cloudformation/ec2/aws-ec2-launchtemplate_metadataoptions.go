package ec2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LaunchTemplate_MetadataOptions AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.MetadataOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-metadataoptions.html
type LaunchTemplate_MetadataOptions struct {

	// HttpEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-metadataoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-metadataoptions-httpendpoint
	HttpEndpoint *types.Value `json:"HttpEndpoint,omitempty"`

	// HttpProtocolIpv6 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-metadataoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-metadataoptions-httpprotocolipv6
	HttpProtocolIpv6 *types.Value `json:"HttpProtocolIpv6,omitempty"`

	// HttpPutResponseHopLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-metadataoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-metadataoptions-httpputresponsehoplimit
	HttpPutResponseHopLimit *types.Value `json:"HttpPutResponseHopLimit,omitempty"`

	// HttpTokens AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-metadataoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-metadataoptions-httptokens
	HttpTokens *types.Value `json:"HttpTokens,omitempty"`

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
func (r *LaunchTemplate_MetadataOptions) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.MetadataOptions"
}
