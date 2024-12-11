package quicksight

import (
	"goformation/v4/cloudformation/policies"
)

// DataSource_DataSourceParameters AWS CloudFormation Resource (AWS::QuickSight::DataSource.DataSourceParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html
type DataSource_DataSourceParameters struct {

	// AmazonElasticsearchParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-amazonelasticsearchparameters
	AmazonElasticsearchParameters *DataSource_AmazonElasticsearchParameters `json:"AmazonElasticsearchParameters,omitempty"`

	// AmazonOpenSearchParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-amazonopensearchparameters
	AmazonOpenSearchParameters *DataSource_AmazonOpenSearchParameters `json:"AmazonOpenSearchParameters,omitempty"`

	// AthenaParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-athenaparameters
	AthenaParameters *DataSource_AthenaParameters `json:"AthenaParameters,omitempty"`

	// AuroraParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-auroraparameters
	AuroraParameters *DataSource_AuroraParameters `json:"AuroraParameters,omitempty"`

	// AuroraPostgreSqlParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-aurorapostgresqlparameters
	AuroraPostgreSqlParameters *DataSource_AuroraPostgreSqlParameters `json:"AuroraPostgreSqlParameters,omitempty"`

	// MariaDbParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-mariadbparameters
	MariaDbParameters *DataSource_MariaDbParameters `json:"MariaDbParameters,omitempty"`

	// MySqlParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-mysqlparameters
	MySqlParameters *DataSource_MySqlParameters `json:"MySqlParameters,omitempty"`

	// OracleParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-oracleparameters
	OracleParameters *DataSource_OracleParameters `json:"OracleParameters,omitempty"`

	// PostgreSqlParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-postgresqlparameters
	PostgreSqlParameters *DataSource_PostgreSqlParameters `json:"PostgreSqlParameters,omitempty"`

	// PrestoParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-prestoparameters
	PrestoParameters *DataSource_PrestoParameters `json:"PrestoParameters,omitempty"`

	// RdsParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-rdsparameters
	RdsParameters *DataSource_RdsParameters `json:"RdsParameters,omitempty"`

	// RedshiftParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-redshiftparameters
	RedshiftParameters *DataSource_RedshiftParameters `json:"RedshiftParameters,omitempty"`

	// S3Parameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-s3parameters
	S3Parameters *DataSource_S3Parameters `json:"S3Parameters,omitempty"`

	// SnowflakeParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-snowflakeparameters
	SnowflakeParameters *DataSource_SnowflakeParameters `json:"SnowflakeParameters,omitempty"`

	// SparkParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-sparkparameters
	SparkParameters *DataSource_SparkParameters `json:"SparkParameters,omitempty"`

	// SqlServerParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-sqlserverparameters
	SqlServerParameters *DataSource_SqlServerParameters `json:"SqlServerParameters,omitempty"`

	// TeradataParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-datasource-datasourceparameters.html#cfn-quicksight-datasource-datasourceparameters-teradataparameters
	TeradataParameters *DataSource_TeradataParameters `json:"TeradataParameters,omitempty"`

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
func (r *DataSource_DataSourceParameters) AWSCloudFormationType() string {
	return "AWS::QuickSight::DataSource.DataSourceParameters"
}
