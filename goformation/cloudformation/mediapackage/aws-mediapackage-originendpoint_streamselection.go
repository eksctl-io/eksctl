package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// OriginEndpoint_StreamSelection AWS CloudFormation Resource (AWS::MediaPackage::OriginEndpoint.StreamSelection)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-streamselection.html
type OriginEndpoint_StreamSelection struct {

	// MaxVideoBitsPerSecond AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-streamselection.html#cfn-mediapackage-originendpoint-streamselection-maxvideobitspersecond
	MaxVideoBitsPerSecond *types.Value `json:"MaxVideoBitsPerSecond,omitempty"`

	// MinVideoBitsPerSecond AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-streamselection.html#cfn-mediapackage-originendpoint-streamselection-minvideobitspersecond
	MinVideoBitsPerSecond *types.Value `json:"MinVideoBitsPerSecond,omitempty"`

	// StreamOrder AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-streamselection.html#cfn-mediapackage-originendpoint-streamselection-streamorder
	StreamOrder *types.Value `json:"StreamOrder,omitempty"`

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
func (r *OriginEndpoint_StreamSelection) AWSCloudFormationType() string {
	return "AWS::MediaPackage::OriginEndpoint.StreamSelection"
}
