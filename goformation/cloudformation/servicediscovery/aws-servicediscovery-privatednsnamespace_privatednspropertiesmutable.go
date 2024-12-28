package servicediscovery

import (
	"goformation/v4/cloudformation/policies"
)

// PrivateDnsNamespace_PrivateDnsPropertiesMutable AWS CloudFormation Resource (AWS::ServiceDiscovery::PrivateDnsNamespace.PrivateDnsPropertiesMutable)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-privatednsnamespace-privatednspropertiesmutable.html
type PrivateDnsNamespace_PrivateDnsPropertiesMutable struct {

	// SOA AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-privatednsnamespace-privatednspropertiesmutable.html#cfn-servicediscovery-privatednsnamespace-privatednspropertiesmutable-soa
	SOA *PrivateDnsNamespace_SOA `json:"SOA,omitempty"`

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
func (r *PrivateDnsNamespace_PrivateDnsPropertiesMutable) AWSCloudFormationType() string {
	return "AWS::ServiceDiscovery::PrivateDnsNamespace.PrivateDnsPropertiesMutable"
}
