package emr

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Cluster_InstanceFleetProvisioningSpecifications AWS CloudFormation Resource (AWS::EMR::Cluster.InstanceFleetProvisioningSpecifications)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancefleetprovisioningspecifications.html
type Cluster_InstanceFleetProvisioningSpecifications struct {

	// OnDemandSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancefleetprovisioningspecifications.html#cfn-elasticmapreduce-cluster-instancefleetprovisioningspecifications-ondemandspecification
	OnDemandSpecification *Cluster_OnDemandProvisioningSpecification `json:"OnDemandSpecification,omitempty"`

	// SpotSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-instancefleetprovisioningspecifications.html#cfn-elasticmapreduce-cluster-instancefleetprovisioningspecifications-spotspecification
	SpotSpecification *Cluster_SpotProvisioningSpecification `json:"SpotSpecification,omitempty"`

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
func (r *Cluster_InstanceFleetProvisioningSpecifications) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.InstanceFleetProvisioningSpecifications"
}
