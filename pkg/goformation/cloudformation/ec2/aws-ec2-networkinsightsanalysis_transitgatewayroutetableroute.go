package ec2

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// NetworkInsightsAnalysis_TransitGatewayRouteTableRoute AWS CloudFormation Resource (AWS::EC2::NetworkInsightsAnalysis.TransitGatewayRouteTableRoute)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html
type NetworkInsightsAnalysis_TransitGatewayRouteTableRoute struct {

	// AttachmentId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html#cfn-ec2-networkinsightsanalysis-transitgatewayroutetableroute-attachmentid
	AttachmentId *types.Value `json:"AttachmentId,omitempty"`

	// DestinationCidr AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html#cfn-ec2-networkinsightsanalysis-transitgatewayroutetableroute-destinationcidr
	DestinationCidr *types.Value `json:"DestinationCidr,omitempty"`

	// PrefixListId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html#cfn-ec2-networkinsightsanalysis-transitgatewayroutetableroute-prefixlistid
	PrefixListId *types.Value `json:"PrefixListId,omitempty"`

	// ResourceId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html#cfn-ec2-networkinsightsanalysis-transitgatewayroutetableroute-resourceid
	ResourceId *types.Value `json:"ResourceId,omitempty"`

	// ResourceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html#cfn-ec2-networkinsightsanalysis-transitgatewayroutetableroute-resourcetype
	ResourceType *types.Value `json:"ResourceType,omitempty"`

	// RouteOrigin AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html#cfn-ec2-networkinsightsanalysis-transitgatewayroutetableroute-routeorigin
	RouteOrigin *types.Value `json:"RouteOrigin,omitempty"`

	// State AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkinsightsanalysis-transitgatewayroutetableroute.html#cfn-ec2-networkinsightsanalysis-transitgatewayroutetableroute-state
	State *types.Value `json:"State,omitempty"`

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
func (r *NetworkInsightsAnalysis_TransitGatewayRouteTableRoute) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkInsightsAnalysis.TransitGatewayRouteTableRoute"
}
