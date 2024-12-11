package evidently

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Experiment_OnlineAbConfigObject AWS CloudFormation Resource (AWS::Evidently::Experiment.OnlineAbConfigObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-onlineabconfigobject.html
type Experiment_OnlineAbConfigObject struct {

	// ControlTreatmentName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-onlineabconfigobject.html#cfn-evidently-experiment-onlineabconfigobject-controltreatmentname
	ControlTreatmentName *types.Value `json:"ControlTreatmentName,omitempty"`

	// TreatmentWeights AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-onlineabconfigobject.html#cfn-evidently-experiment-onlineabconfigobject-treatmentweights
	TreatmentWeights []Experiment_TreatmentToWeight `json:"TreatmentWeights,omitempty"`

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
func (r *Experiment_OnlineAbConfigObject) AWSCloudFormationType() string {
	return "AWS::Evidently::Experiment.OnlineAbConfigObject"
}
