package emr

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Cluster_KerberosAttributes AWS CloudFormation Resource (AWS::EMR::Cluster.KerberosAttributes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-kerberosattributes.html
type Cluster_KerberosAttributes struct {

	// ADDomainJoinPassword AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-kerberosattributes.html#cfn-elasticmapreduce-cluster-kerberosattributes-addomainjoinpassword
	ADDomainJoinPassword *types.Value `json:"ADDomainJoinPassword,omitempty"`

	// ADDomainJoinUser AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-kerberosattributes.html#cfn-elasticmapreduce-cluster-kerberosattributes-addomainjoinuser
	ADDomainJoinUser *types.Value `json:"ADDomainJoinUser,omitempty"`

	// CrossRealmTrustPrincipalPassword AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-kerberosattributes.html#cfn-elasticmapreduce-cluster-kerberosattributes-crossrealmtrustprincipalpassword
	CrossRealmTrustPrincipalPassword *types.Value `json:"CrossRealmTrustPrincipalPassword,omitempty"`

	// KdcAdminPassword AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-kerberosattributes.html#cfn-elasticmapreduce-cluster-kerberosattributes-kdcadminpassword
	KdcAdminPassword *types.Value `json:"KdcAdminPassword,omitempty"`

	// Realm AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-kerberosattributes.html#cfn-elasticmapreduce-cluster-kerberosattributes-realm
	Realm *types.Value `json:"Realm,omitempty"`

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
func (r *Cluster_KerberosAttributes) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.KerberosAttributes"
}
