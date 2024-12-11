package serverless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// StateMachine_EventBridgeRuleEvent AWS CloudFormation Resource (AWS::Serverless::StateMachine.EventBridgeRuleEvent)
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-statemachine-cloudwatchevent.html
type StateMachine_EventBridgeRuleEvent struct {

	// EventBusName AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-statemachine-cloudwatchevent.html
	EventBusName *types.Value `json:"EventBusName,omitempty"`

	// Input AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-statemachine-cloudwatchevent.html
	Input *types.Value `json:"Input,omitempty"`

	// InputPath AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-statemachine-cloudwatchevent.html
	InputPath *types.Value `json:"InputPath,omitempty"`

	// Pattern AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AmazonCloudWatch/latest/events/CloudWatchEventsandEventPatterns.html
	Pattern interface{} `json:"Pattern,omitempty"`

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
func (r *StateMachine_EventBridgeRuleEvent) AWSCloudFormationType() string {
	return "AWS::Serverless::StateMachine.EventBridgeRuleEvent"
}
