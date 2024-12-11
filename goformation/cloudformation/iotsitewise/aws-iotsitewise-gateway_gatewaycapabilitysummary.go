package iotsitewise

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Gateway_GatewayCapabilitySummary AWS CloudFormation Resource (AWS::IoTSiteWise::Gateway.GatewayCapabilitySummary)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-gateway-gatewaycapabilitysummary.html
type Gateway_GatewayCapabilitySummary struct {

	// CapabilityConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-gateway-gatewaycapabilitysummary.html#cfn-iotsitewise-gateway-gatewaycapabilitysummary-capabilityconfiguration
	CapabilityConfiguration *types.Value `json:"CapabilityConfiguration,omitempty"`

	// CapabilityNamespace AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-gateway-gatewaycapabilitysummary.html#cfn-iotsitewise-gateway-gatewaycapabilitysummary-capabilitynamespace
	CapabilityNamespace *types.Value `json:"CapabilityNamespace,omitempty"`

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
func (r *Gateway_GatewayCapabilitySummary) AWSCloudFormationType() string {
	return "AWS::IoTSiteWise::Gateway.GatewayCapabilitySummary"
}
