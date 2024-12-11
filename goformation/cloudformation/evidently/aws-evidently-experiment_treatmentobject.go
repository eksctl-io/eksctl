package evidently

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Experiment_TreatmentObject AWS CloudFormation Resource (AWS::Evidently::Experiment.TreatmentObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmentobject.html
type Experiment_TreatmentObject struct {

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmentobject.html#cfn-evidently-experiment-treatmentobject-description
	Description *types.Value `json:"Description,omitempty"`

	// Feature AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmentobject.html#cfn-evidently-experiment-treatmentobject-feature
	Feature *types.Value `json:"Feature,omitempty"`

	// TreatmentName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmentobject.html#cfn-evidently-experiment-treatmentobject-treatmentname
	TreatmentName *types.Value `json:"TreatmentName,omitempty"`

	// Variation AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-treatmentobject.html#cfn-evidently-experiment-treatmentobject-variation
	Variation *types.Value `json:"Variation,omitempty"`

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
func (r *Experiment_TreatmentObject) AWSCloudFormationType() string {
	return "AWS::Evidently::Experiment.TreatmentObject"
}
