package imagebuilder

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DistributionConfiguration_Distribution AWS CloudFormation Resource (AWS::ImageBuilder::DistributionConfiguration.Distribution)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-distributionconfiguration-distribution.html
type DistributionConfiguration_Distribution struct {

	// AmiDistributionConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-distributionconfiguration-distribution.html#cfn-imagebuilder-distributionconfiguration-distribution-amidistributionconfiguration
	AmiDistributionConfiguration interface{} `json:"AmiDistributionConfiguration,omitempty"`

	// ContainerDistributionConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-distributionconfiguration-distribution.html#cfn-imagebuilder-distributionconfiguration-distribution-containerdistributionconfiguration
	ContainerDistributionConfiguration interface{} `json:"ContainerDistributionConfiguration,omitempty"`

	// LaunchTemplateConfigurations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-distributionconfiguration-distribution.html#cfn-imagebuilder-distributionconfiguration-distribution-launchtemplateconfigurations
	LaunchTemplateConfigurations []DistributionConfiguration_LaunchTemplateConfiguration `json:"LaunchTemplateConfigurations,omitempty"`

	// LicenseConfigurationArns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-distributionconfiguration-distribution.html#cfn-imagebuilder-distributionconfiguration-distribution-licenseconfigurationarns
	LicenseConfigurationArns *types.Value `json:"LicenseConfigurationArns,omitempty"`

	// Region AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-distributionconfiguration-distribution.html#cfn-imagebuilder-distributionconfiguration-distribution-region
	Region *types.Value `json:"Region,omitempty"`

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
func (r *DistributionConfiguration_Distribution) AWSCloudFormationType() string {
	return "AWS::ImageBuilder::DistributionConfiguration.Distribution"
}
