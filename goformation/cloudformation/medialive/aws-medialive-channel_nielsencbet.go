package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_NielsenCBET AWS CloudFormation Resource (AWS::MediaLive::Channel.NielsenCBET)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsencbet.html
type Channel_NielsenCBET struct {

	// CbetCheckDigitString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsencbet.html#cfn-medialive-channel-nielsencbet-cbetcheckdigitstring
	CbetCheckDigitString *types.Value `json:"CbetCheckDigitString,omitempty"`

	// CbetStepaside AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsencbet.html#cfn-medialive-channel-nielsencbet-cbetstepaside
	CbetStepaside *types.Value `json:"CbetStepaside,omitempty"`

	// Csid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsencbet.html#cfn-medialive-channel-nielsencbet-csid
	Csid *types.Value `json:"Csid,omitempty"`

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
func (r *Channel_NielsenCBET) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.NielsenCBET"
}
