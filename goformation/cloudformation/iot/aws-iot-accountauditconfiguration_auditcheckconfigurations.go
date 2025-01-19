package iot

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// AccountAuditConfiguration_AuditCheckConfigurations AWS CloudFormation Resource (AWS::IoT::AccountAuditConfiguration.AuditCheckConfigurations)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html
type AccountAuditConfiguration_AuditCheckConfigurations struct {

	// AuthenticatedCognitoRoleOverlyPermissiveCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-authenticatedcognitoroleoverlypermissivecheck
	AuthenticatedCognitoRoleOverlyPermissiveCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"AuthenticatedCognitoRoleOverlyPermissiveCheck,omitempty"`

	// CaCertificateExpiringCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-cacertificateexpiringcheck
	CaCertificateExpiringCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"CaCertificateExpiringCheck,omitempty"`

	// CaCertificateKeyQualityCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-cacertificatekeyqualitycheck
	CaCertificateKeyQualityCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"CaCertificateKeyQualityCheck,omitempty"`

	// ConflictingClientIdsCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-conflictingclientidscheck
	ConflictingClientIdsCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"ConflictingClientIdsCheck,omitempty"`

	// DeviceCertificateExpiringCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-devicecertificateexpiringcheck
	DeviceCertificateExpiringCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"DeviceCertificateExpiringCheck,omitempty"`

	// DeviceCertificateKeyQualityCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-devicecertificatekeyqualitycheck
	DeviceCertificateKeyQualityCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"DeviceCertificateKeyQualityCheck,omitempty"`

	// DeviceCertificateSharedCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-devicecertificatesharedcheck
	DeviceCertificateSharedCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"DeviceCertificateSharedCheck,omitempty"`

	// IotPolicyOverlyPermissiveCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-iotpolicyoverlypermissivecheck
	IotPolicyOverlyPermissiveCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"IotPolicyOverlyPermissiveCheck,omitempty"`

	// IotRoleAliasAllowsAccessToUnusedServicesCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-iotrolealiasallowsaccesstounusedservicescheck
	IotRoleAliasAllowsAccessToUnusedServicesCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"IotRoleAliasAllowsAccessToUnusedServicesCheck,omitempty"`

	// IotRoleAliasOverlyPermissiveCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-iotrolealiasoverlypermissivecheck
	IotRoleAliasOverlyPermissiveCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"IotRoleAliasOverlyPermissiveCheck,omitempty"`

	// LoggingDisabledCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-loggingdisabledcheck
	LoggingDisabledCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"LoggingDisabledCheck,omitempty"`

	// RevokedCaCertificateStillActiveCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-revokedcacertificatestillactivecheck
	RevokedCaCertificateStillActiveCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"RevokedCaCertificateStillActiveCheck,omitempty"`

	// RevokedDeviceCertificateStillActiveCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-revokeddevicecertificatestillactivecheck
	RevokedDeviceCertificateStillActiveCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"RevokedDeviceCertificateStillActiveCheck,omitempty"`

	// UnauthenticatedCognitoRoleOverlyPermissiveCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-accountauditconfiguration-auditcheckconfigurations.html#cfn-iot-accountauditconfiguration-auditcheckconfigurations-unauthenticatedcognitoroleoverlypermissivecheck
	UnauthenticatedCognitoRoleOverlyPermissiveCheck *AccountAuditConfiguration_AuditCheckConfiguration `json:"UnauthenticatedCognitoRoleOverlyPermissiveCheck,omitempty"`

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
func (r *AccountAuditConfiguration_AuditCheckConfigurations) AWSCloudFormationType() string {
	return "AWS::IoT::AccountAuditConfiguration.AuditCheckConfigurations"
}
