package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// OriginEndpoint_CmafEncryption AWS CloudFormation Resource (AWS::MediaPackage::OriginEndpoint.CmafEncryption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafencryption.html
type OriginEndpoint_CmafEncryption struct {

	// ConstantInitializationVector AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafencryption.html#cfn-mediapackage-originendpoint-cmafencryption-constantinitializationvector
	ConstantInitializationVector *types.Value `json:"ConstantInitializationVector,omitempty"`

	// KeyRotationIntervalSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafencryption.html#cfn-mediapackage-originendpoint-cmafencryption-keyrotationintervalseconds
	KeyRotationIntervalSeconds *types.Value `json:"KeyRotationIntervalSeconds,omitempty"`

	// SpekeKeyProvider AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-cmafencryption.html#cfn-mediapackage-originendpoint-cmafencryption-spekekeyprovider
	SpekeKeyProvider *OriginEndpoint_SpekeKeyProvider `json:"SpekeKeyProvider,omitempty"`

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
func (r *OriginEndpoint_CmafEncryption) AWSCloudFormationType() string {
	return "AWS::MediaPackage::OriginEndpoint.CmafEncryption"
}
