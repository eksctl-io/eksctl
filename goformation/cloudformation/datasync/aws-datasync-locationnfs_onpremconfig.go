package datasync

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LocationNFS_OnPremConfig AWS CloudFormation Resource (AWS::DataSync::LocationNFS.OnPremConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-locationnfs-onpremconfig.html
type LocationNFS_OnPremConfig struct {

	// AgentArns AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-locationnfs-onpremconfig.html#cfn-datasync-locationnfs-onpremconfig-agentarns
	AgentArns *types.Value `json:"AgentArns,omitempty"`

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
func (r *LocationNFS_OnPremConfig) AWSCloudFormationType() string {
	return "AWS::DataSync::LocationNFS.OnPremConfig"
}
