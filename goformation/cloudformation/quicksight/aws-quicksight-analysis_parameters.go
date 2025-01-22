package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Analysis_Parameters AWS CloudFormation Resource (AWS::QuickSight::Analysis.Parameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-parameters.html
type Analysis_Parameters struct {

	// DateTimeParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-parameters.html#cfn-quicksight-analysis-parameters-datetimeparameters
	DateTimeParameters []Analysis_DateTimeParameter `json:"DateTimeParameters,omitempty"`

	// DecimalParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-parameters.html#cfn-quicksight-analysis-parameters-decimalparameters
	DecimalParameters []Analysis_DecimalParameter `json:"DecimalParameters,omitempty"`

	// IntegerParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-parameters.html#cfn-quicksight-analysis-parameters-integerparameters
	IntegerParameters []Analysis_IntegerParameter `json:"IntegerParameters,omitempty"`

	// StringParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-analysis-parameters.html#cfn-quicksight-analysis-parameters-stringparameters
	StringParameters []Analysis_StringParameter `json:"StringParameters,omitempty"`

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
func (r *Analysis_Parameters) AWSCloudFormationType() string {
	return "AWS::QuickSight::Analysis.Parameters"
}
