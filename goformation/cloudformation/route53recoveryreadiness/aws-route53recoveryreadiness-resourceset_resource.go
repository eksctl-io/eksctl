package route53recoveryreadiness

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResourceSet_Resource AWS CloudFormation Resource (AWS::Route53RecoveryReadiness::ResourceSet.Resource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoveryreadiness-resourceset-resource.html
type ResourceSet_Resource struct {

	// ComponentId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoveryreadiness-resourceset-resource.html#cfn-route53recoveryreadiness-resourceset-resource-componentid
	ComponentId *types.Value `json:"ComponentId,omitempty"`

	// DnsTargetResource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoveryreadiness-resourceset-resource.html#cfn-route53recoveryreadiness-resourceset-resource-dnstargetresource
	DnsTargetResource *ResourceSet_DNSTargetResource `json:"DnsTargetResource,omitempty"`

	// ReadinessScopes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoveryreadiness-resourceset-resource.html#cfn-route53recoveryreadiness-resourceset-resource-readinessscopes
	ReadinessScopes *types.Value `json:"ReadinessScopes,omitempty"`

	// ResourceArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53recoveryreadiness-resourceset-resource.html#cfn-route53recoveryreadiness-resourceset-resource-resourcearn
	ResourceArn *types.Value `json:"ResourceArn,omitempty"`

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
func (r *ResourceSet_Resource) AWSCloudFormationType() string {
	return "AWS::Route53RecoveryReadiness::ResourceSet.Resource"
}
