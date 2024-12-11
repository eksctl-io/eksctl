package serverless

import (
	"bytes"
	"encoding/json"
	"fmt"

	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StateMachine AWS CloudFormation Resource (AWS::Serverless::StateMachine)
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
type StateMachine struct {

	// Definition AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Definition interface{} `json:"Definition,omitempty"`

	// DefinitionSubstitutions AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	DefinitionSubstitutions map[string]*types.Value `json:"DefinitionSubstitutions,omitempty"`

	// DefinitionUri AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	DefinitionUri *StateMachine_DefinitionUri `json:"DefinitionUri,omitempty"`

	// Events AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Events map[string]StateMachine_EventSource `json:"Events,omitempty"`

	// Logging AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Logging *StateMachine_LoggingConfiguration `json:"Logging,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Name *types.Value `json:"Name,omitempty"`

	// Policies AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Policies *StateMachine_Policies `json:"Policies,omitempty"`

	// Role AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Role *types.Value `json:"Role,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Tags map[string]*types.Value `json:"Tags,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
	Type *types.Value `json:"Type,omitempty"`

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
func (r *StateMachine) AWSCloudFormationType() string {
	return "AWS::Serverless::StateMachine"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r StateMachine) MarshalJSON() ([]byte, error) {
	type Properties StateMachine
	return json.Marshal(&struct {
		Type                string
		Properties          Properties
		DependsOn           []string                     `json:"DependsOn,omitempty"`
		Metadata            map[string]interface{}       `json:"Metadata,omitempty"`
		DeletionPolicy      policies.DeletionPolicy      `json:"DeletionPolicy,omitempty"`
		UpdateReplacePolicy policies.UpdateReplacePolicy `json:"UpdateReplacePolicy,omitempty"`
		Condition           string                       `json:"Condition,omitempty"`
	}{
		Type:                r.AWSCloudFormationType(),
		Properties:          (Properties)(r),
		DependsOn:           r.AWSCloudFormationDependsOn,
		Metadata:            r.AWSCloudFormationMetadata,
		DeletionPolicy:      r.AWSCloudFormationDeletionPolicy,
		UpdateReplacePolicy: r.AWSCloudFormationUpdateReplacePolicy,
		Condition:           r.AWSCloudFormationCondition,
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *StateMachine) UnmarshalJSON(b []byte) error {
	type Properties StateMachine
	res := &struct {
		Type                string
		Properties          *Properties
		DependsOn           []string
		Metadata            map[string]interface{}
		DeletionPolicy      string
		UpdateReplacePolicy string
		Condition           string
	}{}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = StateMachine(*res.Properties)
	}
	if res.DependsOn != nil {
		r.AWSCloudFormationDependsOn = res.DependsOn
	}
	if res.Metadata != nil {
		r.AWSCloudFormationMetadata = res.Metadata
	}
	if res.DeletionPolicy != "" {
		r.AWSCloudFormationDeletionPolicy = policies.DeletionPolicy(res.DeletionPolicy)
	}
	if res.UpdateReplacePolicy != "" {
		r.AWSCloudFormationUpdateReplacePolicy = policies.UpdateReplacePolicy(res.UpdateReplacePolicy)
	}
	if res.Condition != "" {
		r.AWSCloudFormationCondition = res.Condition
	}
	return nil
}
