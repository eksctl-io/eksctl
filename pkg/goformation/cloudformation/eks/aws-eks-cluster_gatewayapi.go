package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// Cluster_GatewayApi AWS CloudFormation Resource (AWS::EKS::Cluster.GatewayApi)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-gatewayapi.html
type Cluster_GatewayApi struct {

	// Crds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-eks-cluster-gatewayapi.html#cfn-eks-cluster-gatewayapi-crds
	Crds *Cluster_GatewayApiCrds `json:"Crds,omitempty"`

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
func (r *Cluster_GatewayApi) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.GatewayApi"
}
