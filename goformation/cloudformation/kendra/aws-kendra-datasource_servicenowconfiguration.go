package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_ServiceNowConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.ServiceNowConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-servicenowconfiguration.html
type DataSource_ServiceNowConfiguration struct {

	// AuthenticationType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-servicenowconfiguration.html#cfn-kendra-datasource-servicenowconfiguration-authenticationtype
	AuthenticationType *types.Value `json:"AuthenticationType,omitempty"`

	// HostUrl AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-servicenowconfiguration.html#cfn-kendra-datasource-servicenowconfiguration-hosturl
	HostUrl *types.Value `json:"HostUrl,omitempty"`

	// KnowledgeArticleConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-servicenowconfiguration.html#cfn-kendra-datasource-servicenowconfiguration-knowledgearticleconfiguration
	KnowledgeArticleConfiguration *DataSource_ServiceNowKnowledgeArticleConfiguration `json:"KnowledgeArticleConfiguration,omitempty"`

	// SecretArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-servicenowconfiguration.html#cfn-kendra-datasource-servicenowconfiguration-secretarn
	SecretArn *types.Value `json:"SecretArn,omitempty"`

	// ServiceCatalogConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-servicenowconfiguration.html#cfn-kendra-datasource-servicenowconfiguration-servicecatalogconfiguration
	ServiceCatalogConfiguration *DataSource_ServiceNowServiceCatalogConfiguration `json:"ServiceCatalogConfiguration,omitempty"`

	// ServiceNowBuildVersion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-servicenowconfiguration.html#cfn-kendra-datasource-servicenowconfiguration-servicenowbuildversion
	ServiceNowBuildVersion *types.Value `json:"ServiceNowBuildVersion,omitempty"`

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
func (r *DataSource_ServiceNowConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ServiceNowConfiguration"
}
