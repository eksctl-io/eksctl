package evidently

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Experiment_TreatmentToWeight AWS CloudFormation Resource (AWS::Evidently::Experiment.TreatmentToWeight)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmenttoweight.html
type Experiment_TreatmentToWeight struct {

	// SplitWeight AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmenttoweight.html#cfn-evidently-experiment-treatmenttoweight-splitweight
	SplitWeight *types.Value `json:"SplitWeight"`

	// Treatment AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmenttoweight.html#cfn-evidently-experiment-treatmenttoweight-treatment
	Treatment *types.Value `json:"Treatment,omitempty"`

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
func (r *Experiment_TreatmentToWeight) AWSCloudFormationType() string {
	return "AWS::Evidently::Experiment.TreatmentToWeight"
}
