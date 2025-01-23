package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_SalesforceCustomKnowledgeArticleTypeConfigurationList AWS CloudFormation Resource (AWS::Kendra::DataSource.SalesforceCustomKnowledgeArticleTypeConfigurationList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcecustomknowledgearticletypeconfigurationlist.html
type DataSource_SalesforceCustomKnowledgeArticleTypeConfigurationList struct {

	// SalesforceCustomKnowledgeArticleTypeConfigurationList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-salesforcecustomknowledgearticletypeconfigurationlist.html#cfn-kendra-datasource-salesforcecustomknowledgearticletypeconfigurationlist-salesforcecustomknowledgearticletypeconfigurationlist
	SalesforceCustomKnowledgeArticleTypeConfigurationList []DataSource_SalesforceCustomKnowledgeArticleTypeConfiguration `json:"SalesforceCustomKnowledgeArticleTypeConfigurationList,omitempty"`

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
func (r *DataSource_SalesforceCustomKnowledgeArticleTypeConfigurationList) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.SalesforceCustomKnowledgeArticleTypeConfigurationList"
}
