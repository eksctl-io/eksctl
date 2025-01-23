package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Dashboard_Parameters AWS CloudFormation Resource (AWS::QuickSight::Dashboard.Parameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-parameters.html
type Dashboard_Parameters struct {

	// DateTimeParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-parameters.html#cfn-quicksight-dashboard-parameters-datetimeparameters
	DateTimeParameters []Dashboard_DateTimeParameter `json:"DateTimeParameters,omitempty"`

	// DecimalParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-parameters.html#cfn-quicksight-dashboard-parameters-decimalparameters
	DecimalParameters []Dashboard_DecimalParameter `json:"DecimalParameters,omitempty"`

	// IntegerParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-parameters.html#cfn-quicksight-dashboard-parameters-integerparameters
	IntegerParameters []Dashboard_IntegerParameter `json:"IntegerParameters,omitempty"`

	// StringParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-parameters.html#cfn-quicksight-dashboard-parameters-stringparameters
	StringParameters []Dashboard_StringParameter `json:"StringParameters,omitempty"`

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
func (r *Dashboard_Parameters) AWSCloudFormationType() string {
	return "AWS::QuickSight::Dashboard.Parameters"
}
