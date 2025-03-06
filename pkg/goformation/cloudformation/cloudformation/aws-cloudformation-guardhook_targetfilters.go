package cloudformation

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// GuardHook_TargetFilters AWS CloudFormation Resource (AWS::CloudFormation::GuardHook.TargetFilters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-guardhook-targetfilters.html
type GuardHook_TargetFilters struct {

	// Actions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-guardhook-targetfilters.html#cfn-cloudformation-guardhook-targetfilters-actions
	Actions *types.Value `json:"Actions,omitempty"`

	// InvocationPoints AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-guardhook-targetfilters.html#cfn-cloudformation-guardhook-targetfilters-invocationpoints
	InvocationPoints *types.Value `json:"InvocationPoints,omitempty"`

	// TargetNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-guardhook-targetfilters.html#cfn-cloudformation-guardhook-targetfilters-targetnames
	TargetNames *types.Value `json:"TargetNames,omitempty"`

	// Targets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudformation-guardhook-targetfilters.html#cfn-cloudformation-guardhook-targetfilters-targets
	Targets []GuardHook_HookTarget `json:"Targets,omitempty"`

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
func (r *GuardHook_TargetFilters) AWSCloudFormationType() string {
	return "AWS::CloudFormation::GuardHook.TargetFilters"
}
