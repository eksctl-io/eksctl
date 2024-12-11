package healthlake

import (
	"goformation/v4/cloudformation/policies"
)

// FHIRDatastore_SseConfiguration AWS CloudFormation Resource (AWS::HealthLake::FHIRDatastore.SseConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-healthlake-fhirdatastore-sseconfiguration.html
type FHIRDatastore_SseConfiguration struct {

	// KmsEncryptionConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-healthlake-fhirdatastore-sseconfiguration.html#cfn-healthlake-fhirdatastore-sseconfiguration-kmsencryptionconfig
	KmsEncryptionConfig *FHIRDatastore_KmsEncryptionConfig `json:"KmsEncryptionConfig,omitempty"`

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
func (r *FHIRDatastore_SseConfiguration) AWSCloudFormationType() string {
	return "AWS::HealthLake::FHIRDatastore.SseConfiguration"
}
