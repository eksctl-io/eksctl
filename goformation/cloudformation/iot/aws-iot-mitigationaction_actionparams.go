package iot

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// MitigationAction_ActionParams AWS CloudFormation Resource (AWS::IoT::MitigationAction.ActionParams)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-actionparams.html
type MitigationAction_ActionParams struct {

	// AddThingsToThingGroupParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-actionparams.html#cfn-iot-mitigationaction-actionparams-addthingstothinggroupparams
	AddThingsToThingGroupParams *MitigationAction_AddThingsToThingGroupParams `json:"AddThingsToThingGroupParams,omitempty"`

	// EnableIoTLoggingParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-actionparams.html#cfn-iot-mitigationaction-actionparams-enableiotloggingparams
	EnableIoTLoggingParams *MitigationAction_EnableIoTLoggingParams `json:"EnableIoTLoggingParams,omitempty"`

	// PublishFindingToSnsParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-actionparams.html#cfn-iot-mitigationaction-actionparams-publishfindingtosnsparams
	PublishFindingToSnsParams *MitigationAction_PublishFindingToSnsParams `json:"PublishFindingToSnsParams,omitempty"`

	// ReplaceDefaultPolicyVersionParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-actionparams.html#cfn-iot-mitigationaction-actionparams-replacedefaultpolicyversionparams
	ReplaceDefaultPolicyVersionParams *MitigationAction_ReplaceDefaultPolicyVersionParams `json:"ReplaceDefaultPolicyVersionParams,omitempty"`

	// UpdateCACertificateParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-actionparams.html#cfn-iot-mitigationaction-actionparams-updatecacertificateparams
	UpdateCACertificateParams *MitigationAction_UpdateCACertificateParams `json:"UpdateCACertificateParams,omitempty"`

	// UpdateDeviceCertificateParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-actionparams.html#cfn-iot-mitigationaction-actionparams-updatedevicecertificateparams
	UpdateDeviceCertificateParams *MitigationAction_UpdateDeviceCertificateParams `json:"UpdateDeviceCertificateParams,omitempty"`

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
func (r *MitigationAction_ActionParams) AWSCloudFormationType() string {
	return "AWS::IoT::MitigationAction.ActionParams"
}
