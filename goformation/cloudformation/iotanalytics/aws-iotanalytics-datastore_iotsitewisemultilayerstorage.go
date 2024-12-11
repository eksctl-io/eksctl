package iotanalytics

import (
	"goformation/v4/cloudformation/policies"
)

// Datastore_IotSiteWiseMultiLayerStorage AWS CloudFormation Resource (AWS::IoTAnalytics::Datastore.IotSiteWiseMultiLayerStorage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-iotsitewisemultilayerstorage.html
type Datastore_IotSiteWiseMultiLayerStorage struct {

	// CustomerManagedS3Storage AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-datastore-iotsitewisemultilayerstorage.html#cfn-iotanalytics-datastore-iotsitewisemultilayerstorage-customermanageds3storage
	CustomerManagedS3Storage *Datastore_CustomerManagedS3Storage `json:"CustomerManagedS3Storage,omitempty"`

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
func (r *Datastore_IotSiteWiseMultiLayerStorage) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Datastore.IotSiteWiseMultiLayerStorage"
}
