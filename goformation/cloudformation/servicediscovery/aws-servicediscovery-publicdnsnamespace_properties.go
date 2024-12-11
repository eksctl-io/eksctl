package servicediscovery

import (
	"goformation/v4/cloudformation/policies"
)

// PublicDnsNamespace_Properties AWS CloudFormation Resource (AWS::ServiceDiscovery::PublicDnsNamespace.Properties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-publicdnsnamespace-properties.html
type PublicDnsNamespace_Properties struct {

	// DnsProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-publicdnsnamespace-properties.html#cfn-servicediscovery-publicdnsnamespace-properties-dnsproperties
	DnsProperties *PublicDnsNamespace_PublicDnsPropertiesMutable `json:"DnsProperties,omitempty"`

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
func (r *PublicDnsNamespace_Properties) AWSCloudFormationType() string {
	return "AWS::ServiceDiscovery::PublicDnsNamespace.Properties"
}
