package iot

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// AccountAuditConfiguration_AuditNotificationTargetConfigurations AWS CloudFormation Resource (AWS::IoT::AccountAuditConfiguration.AuditNotificationTargetConfigurations)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditnotificationtargetconfigurations.html
type AccountAuditConfiguration_AuditNotificationTargetConfigurations struct {

	// Sns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditnotificationtargetconfigurations.html#cfn-iot-accountauditconfiguration-auditnotificationtargetconfigurations-sns
	Sns *AccountAuditConfiguration_AuditNotificationTarget `json:"Sns,omitempty"`

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
func (r *AccountAuditConfiguration_AuditNotificationTargetConfigurations) AWSCloudFormationType() string {
	return "AWS::IoT::AccountAuditConfiguration.AuditNotificationTargetConfigurations"
}
