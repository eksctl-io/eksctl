package eks

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

// Cluster_ServiceNodePortRange AWS CloudFormation Resource (AWS::EKS::Cluster.ServiceNodePortRange)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-servicenodeportrange.html
type Cluster_ServiceNodePortRange struct {

	// MinPort AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-servicenodeportrange.html#cfn-eks-cluster-servicenodeportrange-minport
	MinPort *types.Value `json:"MinPort,omitempty"`

	// MaxPort AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-properties-eks-cluster-servicenodeportrange.html#cfn-eks-cluster-servicenodeportrange-maxport
	MaxPort *types.Value `json:"MaxPort,omitempty"`

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
func (r *Cluster_ServiceNodePortRange) AWSCloudFormationType() string {
	return "AWS::EKS::Cluster.ServiceNodePortRange"
}
