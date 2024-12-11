package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_M3u8Settings AWS CloudFormation Resource (AWS::MediaLive::Channel.M3u8Settings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html
type Channel_M3u8Settings struct {

	// AudioFramesPerPes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-audioframesperpes
	AudioFramesPerPes *types.Value `json:"AudioFramesPerPes,omitempty"`

	// AudioPids AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-audiopids
	AudioPids *types.Value `json:"AudioPids,omitempty"`

	// EcmPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-ecmpid
	EcmPid *types.Value `json:"EcmPid,omitempty"`

	// NielsenId3Behavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-nielsenid3behavior
	NielsenId3Behavior *types.Value `json:"NielsenId3Behavior,omitempty"`

	// PatInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-patinterval
	PatInterval *types.Value `json:"PatInterval,omitempty"`

	// PcrControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-pcrcontrol
	PcrControl *types.Value `json:"PcrControl,omitempty"`

	// PcrPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-pcrperiod
	PcrPeriod *types.Value `json:"PcrPeriod,omitempty"`

	// PcrPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-pcrpid
	PcrPid *types.Value `json:"PcrPid,omitempty"`

	// PmtInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-pmtinterval
	PmtInterval *types.Value `json:"PmtInterval,omitempty"`

	// PmtPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-pmtpid
	PmtPid *types.Value `json:"PmtPid,omitempty"`

	// ProgramNum AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-programnum
	ProgramNum *types.Value `json:"ProgramNum,omitempty"`

	// Scte35Behavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-scte35behavior
	Scte35Behavior *types.Value `json:"Scte35Behavior,omitempty"`

	// Scte35Pid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-scte35pid
	Scte35Pid *types.Value `json:"Scte35Pid,omitempty"`

	// TimedMetadataBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-timedmetadatabehavior
	TimedMetadataBehavior *types.Value `json:"TimedMetadataBehavior,omitempty"`

	// TimedMetadataPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-timedmetadatapid
	TimedMetadataPid *types.Value `json:"TimedMetadataPid,omitempty"`

	// TransportStreamId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-transportstreamid
	TransportStreamId *types.Value `json:"TransportStreamId,omitempty"`

	// VideoPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m3u8settings.html#cfn-medialive-channel-m3u8settings-videopid
	VideoPid *types.Value `json:"VideoPid,omitempty"`

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
func (r *Channel_M3u8Settings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.M3u8Settings"
}
