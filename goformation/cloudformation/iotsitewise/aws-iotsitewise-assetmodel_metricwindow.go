package iotsitewise

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// AssetModel_MetricWindow AWS CloudFormation Resource (AWS::IoTSiteWise::AssetModel.MetricWindow)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-metricwindow.html
type AssetModel_MetricWindow struct {

	// Tumbling AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-metricwindow.html#cfn-iotsitewise-assetmodel-metricwindow-tumbling
	Tumbling *AssetModel_TumblingWindow `json:"Tumbling,omitempty"`

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
func (r *AssetModel_MetricWindow) AWSCloudFormationType() string {
	return "AWS::IoTSiteWise::AssetModel.MetricWindow"
}
