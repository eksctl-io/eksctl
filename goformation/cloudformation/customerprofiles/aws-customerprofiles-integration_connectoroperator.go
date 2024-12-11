package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Integration_ConnectorOperator AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.ConnectorOperator)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-connectoroperator.html
type Integration_ConnectorOperator struct {

	// Marketo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-connectoroperator.html#cfn-customerprofiles-integration-connectoroperator-marketo
	Marketo *types.Value `json:"Marketo,omitempty"`

	// S3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-connectoroperator.html#cfn-customerprofiles-integration-connectoroperator-s3
	S3 *types.Value `json:"S3,omitempty"`

	// Salesforce AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-connectoroperator.html#cfn-customerprofiles-integration-connectoroperator-salesforce
	Salesforce *types.Value `json:"Salesforce,omitempty"`

	// ServiceNow AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-connectoroperator.html#cfn-customerprofiles-integration-connectoroperator-servicenow
	ServiceNow *types.Value `json:"ServiceNow,omitempty"`

	// Zendesk AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-connectoroperator.html#cfn-customerprofiles-integration-connectoroperator-zendesk
	Zendesk *types.Value `json:"Zendesk,omitempty"`

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
func (r *Integration_ConnectorOperator) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.ConnectorOperator"
}
