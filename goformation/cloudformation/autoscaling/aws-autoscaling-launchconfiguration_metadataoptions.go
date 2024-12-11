package autoscaling

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LaunchConfiguration_MetadataOptions AWS CloudFormation Resource (AWS::AutoScaling::LaunchConfiguration.MetadataOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-launchconfig-metadataoptions.html
type LaunchConfiguration_MetadataOptions struct {

	// HttpEndpoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-launchconfig-metadataoptions.html#cfn-autoscaling-launchconfig-metadataoptions-httpendpoint
	HttpEndpoint *types.Value `json:"HttpEndpoint,omitempty"`

	// HttpPutResponseHopLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-launchconfig-metadataoptions.html#cfn-autoscaling-launchconfig-metadataoptions-httpputresponsehoplimit
	HttpPutResponseHopLimit *types.Value `json:"HttpPutResponseHopLimit,omitempty"`

	// HttpTokens AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-launchconfig-metadataoptions.html#cfn-autoscaling-launchconfig-metadataoptions-httptokens
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
func (r *LaunchConfiguration_MetadataOptions) AWSCloudFormationType() string {
	return "AWS::AutoScaling::LaunchConfiguration.MetadataOptions"
}
