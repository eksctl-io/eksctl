package servicediscovery

import (
	"goformation/v4/cloudformation/policies"
)

// PublicDnsNamespace_PublicDnsPropertiesMutable AWS CloudFormation Resource (AWS::ServiceDiscovery::PublicDnsNamespace.PublicDnsPropertiesMutable)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-publicdnsnamespace-publicdnspropertiesmutable.html
type PublicDnsNamespace_PublicDnsPropertiesMutable struct {

	// SOA AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-publicdnsnamespace-publicdnspropertiesmutable.html#cfn-servicediscovery-publicdnsnamespace-publicdnspropertiesmutable-soa
	SOA *PublicDnsNamespace_SOA `json:"SOA,omitempty"`

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
func (r *PublicDnsNamespace_PublicDnsPropertiesMutable) AWSCloudFormationType() string {
	return "AWS::ServiceDiscovery::PublicDnsNamespace.PublicDnsPropertiesMutable"
}
