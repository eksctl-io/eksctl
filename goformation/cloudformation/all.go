package cloudformation

import (
	"fmt"

	"goformation/v4/cloudformation/accessanalyzer"
	"goformation/v4/cloudformation/acmpca"
	"goformation/v4/cloudformation/amazonmq"
	"goformation/v4/cloudformation/amplify"
	"goformation/v4/cloudformation/amplifyuibuilder"
	"goformation/v4/cloudformation/apigateway"
	"goformation/v4/cloudformation/apigatewayv2"
	"goformation/v4/cloudformation/appconfig"
	"goformation/v4/cloudformation/appflow"
	"goformation/v4/cloudformation/appintegrations"
	"goformation/v4/cloudformation/applicationautoscaling"
	"goformation/v4/cloudformation/applicationinsights"
	"goformation/v4/cloudformation/appmesh"
	"goformation/v4/cloudformation/apprunner"
	"goformation/v4/cloudformation/appstream"
	"goformation/v4/cloudformation/appsync"
	"goformation/v4/cloudformation/aps"
	"goformation/v4/cloudformation/ask"
	"goformation/v4/cloudformation/athena"
	"goformation/v4/cloudformation/auditmanager"
	"goformation/v4/cloudformation/autoscaling"
	"goformation/v4/cloudformation/autoscalingplans"
	"goformation/v4/cloudformation/backup"
	"goformation/v4/cloudformation/batch"
	"goformation/v4/cloudformation/budgets"
	"goformation/v4/cloudformation/cassandra"
	"goformation/v4/cloudformation/ce"
	"goformation/v4/cloudformation/certificatemanager"
	"goformation/v4/cloudformation/chatbot"
	"goformation/v4/cloudformation/cloud9"
	"goformation/v4/cloudformation/cloudformation"
	"goformation/v4/cloudformation/cloudfront"
	"goformation/v4/cloudformation/cloudtrail"
	"goformation/v4/cloudformation/cloudwatch"
	"goformation/v4/cloudformation/codeartifact"
	"goformation/v4/cloudformation/codebuild"
	"goformation/v4/cloudformation/codecommit"
	"goformation/v4/cloudformation/codedeploy"
	"goformation/v4/cloudformation/codeguruprofiler"
	"goformation/v4/cloudformation/codegurureviewer"
	"goformation/v4/cloudformation/codepipeline"
	"goformation/v4/cloudformation/codestar"
	"goformation/v4/cloudformation/codestarconnections"
	"goformation/v4/cloudformation/codestarnotifications"
	"goformation/v4/cloudformation/cognito"
	"goformation/v4/cloudformation/config"
	"goformation/v4/cloudformation/connect"
	"goformation/v4/cloudformation/cur"
	"goformation/v4/cloudformation/customerprofiles"
	"goformation/v4/cloudformation/databrew"
	"goformation/v4/cloudformation/datapipeline"
	"goformation/v4/cloudformation/datasync"
	"goformation/v4/cloudformation/dax"
	"goformation/v4/cloudformation/detective"
	"goformation/v4/cloudformation/devopsguru"
	"goformation/v4/cloudformation/directoryservice"
	"goformation/v4/cloudformation/dlm"
	"goformation/v4/cloudformation/dms"
	"goformation/v4/cloudformation/docdb"
	"goformation/v4/cloudformation/dynamodb"
	"goformation/v4/cloudformation/ec2"
	"goformation/v4/cloudformation/ecr"
	"goformation/v4/cloudformation/ecs"
	"goformation/v4/cloudformation/efs"
	"goformation/v4/cloudformation/eks"
	"goformation/v4/cloudformation/elasticache"
	"goformation/v4/cloudformation/elasticbeanstalk"
	"goformation/v4/cloudformation/elasticloadbalancing"
	"goformation/v4/cloudformation/elasticloadbalancingv2"
	"goformation/v4/cloudformation/elasticsearch"
	"goformation/v4/cloudformation/emr"
	"goformation/v4/cloudformation/emrcontainers"
	"goformation/v4/cloudformation/events"
	"goformation/v4/cloudformation/eventschemas"
	"goformation/v4/cloudformation/evidently"
	"goformation/v4/cloudformation/finspace"
	"goformation/v4/cloudformation/fis"
	"goformation/v4/cloudformation/fms"
	"goformation/v4/cloudformation/frauddetector"
	"goformation/v4/cloudformation/fsx"
	"goformation/v4/cloudformation/gamelift"
	"goformation/v4/cloudformation/globalaccelerator"
	"goformation/v4/cloudformation/glue"
	"goformation/v4/cloudformation/greengrass"
	"goformation/v4/cloudformation/greengrassv2"
	"goformation/v4/cloudformation/groundstation"
	"goformation/v4/cloudformation/guardduty"
	"goformation/v4/cloudformation/healthlake"
	"goformation/v4/cloudformation/iam"
	"goformation/v4/cloudformation/imagebuilder"
	"goformation/v4/cloudformation/inspector"
	"goformation/v4/cloudformation/iot"
	"goformation/v4/cloudformation/iot1click"
	"goformation/v4/cloudformation/iotanalytics"
	"goformation/v4/cloudformation/iotcoredeviceadvisor"
	"goformation/v4/cloudformation/iotevents"
	"goformation/v4/cloudformation/iotfleethub"
	"goformation/v4/cloudformation/iotsitewise"
	"goformation/v4/cloudformation/iotthingsgraph"
	"goformation/v4/cloudformation/iotwireless"
	"goformation/v4/cloudformation/ivs"
	"goformation/v4/cloudformation/kendra"
	"goformation/v4/cloudformation/kinesis"
	"goformation/v4/cloudformation/kinesisanalytics"
	"goformation/v4/cloudformation/kinesisanalyticsv2"
	"goformation/v4/cloudformation/kinesisfirehose"
	"goformation/v4/cloudformation/kms"
	"goformation/v4/cloudformation/lakeformation"
	"goformation/v4/cloudformation/lambda"
	"goformation/v4/cloudformation/licensemanager"
	"goformation/v4/cloudformation/lightsail"
	"goformation/v4/cloudformation/location"
	"goformation/v4/cloudformation/logs"
	"goformation/v4/cloudformation/lookoutequipment"
	"goformation/v4/cloudformation/lookoutmetrics"
	"goformation/v4/cloudformation/lookoutvision"
	"goformation/v4/cloudformation/macie"
	"goformation/v4/cloudformation/managedblockchain"
	"goformation/v4/cloudformation/mediaconnect"
	"goformation/v4/cloudformation/mediaconvert"
	"goformation/v4/cloudformation/medialive"
	"goformation/v4/cloudformation/mediapackage"
	"goformation/v4/cloudformation/mediastore"
	"goformation/v4/cloudformation/memorydb"
	"goformation/v4/cloudformation/msk"
	"goformation/v4/cloudformation/mwaa"
	"goformation/v4/cloudformation/neptune"
	"goformation/v4/cloudformation/networkfirewall"
	"goformation/v4/cloudformation/networkmanager"
	"goformation/v4/cloudformation/nimblestudio"
	"goformation/v4/cloudformation/opensearchservice"
	"goformation/v4/cloudformation/opsworks"
	"goformation/v4/cloudformation/opsworkscm"
	"goformation/v4/cloudformation/panorama"
	"goformation/v4/cloudformation/pinpoint"
	"goformation/v4/cloudformation/pinpointemail"
	"goformation/v4/cloudformation/qldb"
	"goformation/v4/cloudformation/quicksight"
	"goformation/v4/cloudformation/ram"
	"goformation/v4/cloudformation/rds"
	"goformation/v4/cloudformation/redshift"
	"goformation/v4/cloudformation/refactorspaces"
	"goformation/v4/cloudformation/rekognition"
	"goformation/v4/cloudformation/resiliencehub"
	"goformation/v4/cloudformation/resourcegroups"
	"goformation/v4/cloudformation/robomaker"
	"goformation/v4/cloudformation/route53"
	"goformation/v4/cloudformation/route53recoverycontrol"
	"goformation/v4/cloudformation/route53recoveryreadiness"
	"goformation/v4/cloudformation/route53resolver"
	"goformation/v4/cloudformation/rum"
	"goformation/v4/cloudformation/s3"
	"goformation/v4/cloudformation/s3objectlambda"
	"goformation/v4/cloudformation/s3outposts"
	"goformation/v4/cloudformation/sagemaker"
	"goformation/v4/cloudformation/sdb"
	"goformation/v4/cloudformation/secretsmanager"
	"goformation/v4/cloudformation/securityhub"
	"goformation/v4/cloudformation/serverless"
	"goformation/v4/cloudformation/servicecatalog"
	"goformation/v4/cloudformation/servicecatalogappregistry"
	"goformation/v4/cloudformation/servicediscovery"
	"goformation/v4/cloudformation/ses"
	"goformation/v4/cloudformation/signer"
	"goformation/v4/cloudformation/sns"
	"goformation/v4/cloudformation/sqs"
	"goformation/v4/cloudformation/ssm"
	"goformation/v4/cloudformation/ssmcontacts"
	"goformation/v4/cloudformation/ssmincidents"
	"goformation/v4/cloudformation/sso"
	"goformation/v4/cloudformation/stepfunctions"
	"goformation/v4/cloudformation/synthetics"
	"goformation/v4/cloudformation/timestream"
	"goformation/v4/cloudformation/transfer"
	"goformation/v4/cloudformation/waf"
	"goformation/v4/cloudformation/wafregional"
	"goformation/v4/cloudformation/wafv2"
	"goformation/v4/cloudformation/wisdom"
	"goformation/v4/cloudformation/workspaces"
	"goformation/v4/cloudformation/xray"
)

// AllResources fetches an iterable map all CloudFormation and SAM resources
func AllResources() map[string]Resource {
	return map[string]Resource{
		"AWS::ACMPCA::Certificate":                                    &acmpca.Certificate{},
		"AWS::ACMPCA::CertificateAuthority":                           &acmpca.CertificateAuthority{},
		"AWS::ACMPCA::CertificateAuthorityActivation":                 &acmpca.CertificateAuthorityActivation{},
		"AWS::ACMPCA::Permission":                                     &acmpca.Permission{},
		"AWS::APS::RuleGroupsNamespace":                               &aps.RuleGroupsNamespace{},
		"AWS::APS::Workspace":                                         &aps.Workspace{},
		"AWS::AccessAnalyzer::Analyzer":                               &accessanalyzer.Analyzer{},
		"AWS::AmazonMQ::Broker":                                       &amazonmq.Broker{},
		"AWS::AmazonMQ::Configuration":                                &amazonmq.Configuration{},
		"AWS::AmazonMQ::ConfigurationAssociation":                     &amazonmq.ConfigurationAssociation{},
		"AWS::Amplify::App":                                           &amplify.App{},
		"AWS::Amplify::Branch":                                        &amplify.Branch{},
		"AWS::Amplify::Domain":                                        &amplify.Domain{},
		"AWS::AmplifyUIBuilder::Component":                            &amplifyuibuilder.Component{},
		"AWS::AmplifyUIBuilder::Theme":                                &amplifyuibuilder.Theme{},
		"AWS::ApiGateway::Account":                                    &apigateway.Account{},
		"AWS::ApiGateway::ApiKey":                                     &apigateway.ApiKey{},
		"AWS::ApiGateway::Authorizer":                                 &apigateway.Authorizer{},
		"AWS::ApiGateway::BasePathMapping":                            &apigateway.BasePathMapping{},
		"AWS::ApiGateway::ClientCertificate":                          &apigateway.ClientCertificate{},
		"AWS::ApiGateway::Deployment":                                 &apigateway.Deployment{},
		"AWS::ApiGateway::DocumentationPart":                          &apigateway.DocumentationPart{},
		"AWS::ApiGateway::DocumentationVersion":                       &apigateway.DocumentationVersion{},
		"AWS::ApiGateway::DomainName":                                 &apigateway.DomainName{},
		"AWS::ApiGateway::GatewayResponse":                            &apigateway.GatewayResponse{},
		"AWS::ApiGateway::Method":                                     &apigateway.Method{},
		"AWS::ApiGateway::Model":                                      &apigateway.Model{},
		"AWS::ApiGateway::RequestValidator":                           &apigateway.RequestValidator{},
		"AWS::ApiGateway::Resource":                                   &apigateway.Resource{},
		"AWS::ApiGateway::RestApi":                                    &apigateway.RestApi{},
		"AWS::ApiGateway::Stage":                                      &apigateway.Stage{},
		"AWS::ApiGateway::UsagePlan":                                  &apigateway.UsagePlan{},
		"AWS::ApiGateway::UsagePlanKey":                               &apigateway.UsagePlanKey{},
		"AWS::ApiGateway::VpcLink":                                    &apigateway.VpcLink{},
		"AWS::ApiGatewayV2::Api":                                      &apigatewayv2.Api{},
		"AWS::ApiGatewayV2::ApiGatewayManagedOverrides":               &apigatewayv2.ApiGatewayManagedOverrides{},
		"AWS::ApiGatewayV2::ApiMapping":                               &apigatewayv2.ApiMapping{},
		"AWS::ApiGatewayV2::Authorizer":                               &apigatewayv2.Authorizer{},
		"AWS::ApiGatewayV2::Deployment":                               &apigatewayv2.Deployment{},
		"AWS::ApiGatewayV2::DomainName":                               &apigatewayv2.DomainName{},
		"AWS::ApiGatewayV2::Integration":                              &apigatewayv2.Integration{},
		"AWS::ApiGatewayV2::IntegrationResponse":                      &apigatewayv2.IntegrationResponse{},
		"AWS::ApiGatewayV2::Model":                                    &apigatewayv2.Model{},
		"AWS::ApiGatewayV2::Route":                                    &apigatewayv2.Route{},
		"AWS::ApiGatewayV2::RouteResponse":                            &apigatewayv2.RouteResponse{},
		"AWS::ApiGatewayV2::Stage":                                    &apigatewayv2.Stage{},
		"AWS::ApiGatewayV2::VpcLink":                                  &apigatewayv2.VpcLink{},
		"AWS::AppConfig::Application":                                 &appconfig.Application{},
		"AWS::AppConfig::ConfigurationProfile":                        &appconfig.ConfigurationProfile{},
		"AWS::AppConfig::Deployment":                                  &appconfig.Deployment{},
		"AWS::AppConfig::DeploymentStrategy":                          &appconfig.DeploymentStrategy{},
		"AWS::AppConfig::Environment":                                 &appconfig.Environment{},
		"AWS::AppConfig::HostedConfigurationVersion":                  &appconfig.HostedConfigurationVersion{},
		"AWS::AppFlow::ConnectorProfile":                              &appflow.ConnectorProfile{},
		"AWS::AppFlow::Flow":                                          &appflow.Flow{},
		"AWS::AppIntegrations::EventIntegration":                      &appintegrations.EventIntegration{},
		"AWS::AppMesh::GatewayRoute":                                  &appmesh.GatewayRoute{},
		"AWS::AppMesh::Mesh":                                          &appmesh.Mesh{},
		"AWS::AppMesh::Route":                                         &appmesh.Route{},
		"AWS::AppMesh::VirtualGateway":                                &appmesh.VirtualGateway{},
		"AWS::AppMesh::VirtualNode":                                   &appmesh.VirtualNode{},
		"AWS::AppMesh::VirtualRouter":                                 &appmesh.VirtualRouter{},
		"AWS::AppMesh::VirtualService":                                &appmesh.VirtualService{},
		"AWS::AppRunner::Service":                                     &apprunner.Service{},
		"AWS::AppStream::AppBlock":                                    &appstream.AppBlock{},
		"AWS::AppStream::Application":                                 &appstream.Application{},
		"AWS::AppStream::ApplicationFleetAssociation":                 &appstream.ApplicationFleetAssociation{},
		"AWS::AppStream::DirectoryConfig":                             &appstream.DirectoryConfig{},
		"AWS::AppStream::Fleet":                                       &appstream.Fleet{},
		"AWS::AppStream::ImageBuilder":                                &appstream.ImageBuilder{},
		"AWS::AppStream::Stack":                                       &appstream.Stack{},
		"AWS::AppStream::StackFleetAssociation":                       &appstream.StackFleetAssociation{},
		"AWS::AppStream::StackUserAssociation":                        &appstream.StackUserAssociation{},
		"AWS::AppStream::User":                                        &appstream.User{},
		"AWS::AppSync::ApiCache":                                      &appsync.ApiCache{},
		"AWS::AppSync::ApiKey":                                        &appsync.ApiKey{},
		"AWS::AppSync::DataSource":                                    &appsync.DataSource{},
		"AWS::AppSync::FunctionConfiguration":                         &appsync.FunctionConfiguration{},
		"AWS::AppSync::GraphQLApi":                                    &appsync.GraphQLApi{},
		"AWS::AppSync::GraphQLSchema":                                 &appsync.GraphQLSchema{},
		"AWS::AppSync::Resolver":                                      &appsync.Resolver{},
		"AWS::ApplicationAutoScaling::ScalableTarget":                 &applicationautoscaling.ScalableTarget{},
		"AWS::ApplicationAutoScaling::ScalingPolicy":                  &applicationautoscaling.ScalingPolicy{},
		"AWS::ApplicationInsights::Application":                       &applicationinsights.Application{},
		"AWS::Athena::DataCatalog":                                    &athena.DataCatalog{},
		"AWS::Athena::NamedQuery":                                     &athena.NamedQuery{},
		"AWS::Athena::PreparedStatement":                              &athena.PreparedStatement{},
		"AWS::Athena::WorkGroup":                                      &athena.WorkGroup{},
		"AWS::AuditManager::Assessment":                               &auditmanager.Assessment{},
		"AWS::AutoScaling::AutoScalingGroup":                          &autoscaling.AutoScalingGroup{},
		"AWS::AutoScaling::LaunchConfiguration":                       &autoscaling.LaunchConfiguration{},
		"AWS::AutoScaling::LifecycleHook":                             &autoscaling.LifecycleHook{},
		"AWS::AutoScaling::ScalingPolicy":                             &autoscaling.ScalingPolicy{},
		"AWS::AutoScaling::ScheduledAction":                           &autoscaling.ScheduledAction{},
		"AWS::AutoScaling::WarmPool":                                  &autoscaling.WarmPool{},
		"AWS::AutoScalingPlans::ScalingPlan":                          &autoscalingplans.ScalingPlan{},
		"AWS::Backup::BackupPlan":                                     &backup.BackupPlan{},
		"AWS::Backup::BackupSelection":                                &backup.BackupSelection{},
		"AWS::Backup::BackupVault":                                    &backup.BackupVault{},
		"AWS::Backup::Framework":                                      &backup.Framework{},
		"AWS::Backup::ReportPlan":                                     &backup.ReportPlan{},
		"AWS::Batch::ComputeEnvironment":                              &batch.ComputeEnvironment{},
		"AWS::Batch::JobDefinition":                                   &batch.JobDefinition{},
		"AWS::Batch::JobQueue":                                        &batch.JobQueue{},
		"AWS::Batch::SchedulingPolicy":                                &batch.SchedulingPolicy{},
		"AWS::Budgets::Budget":                                        &budgets.Budget{},
		"AWS::Budgets::BudgetsAction":                                 &budgets.BudgetsAction{},
		"AWS::CE::AnomalyMonitor":                                     &ce.AnomalyMonitor{},
		"AWS::CE::AnomalySubscription":                                &ce.AnomalySubscription{},
		"AWS::CE::CostCategory":                                       &ce.CostCategory{},
		"AWS::CUR::ReportDefinition":                                  &cur.ReportDefinition{},
		"AWS::Cassandra::Keyspace":                                    &cassandra.Keyspace{},
		"AWS::Cassandra::Table":                                       &cassandra.Table{},
		"AWS::CertificateManager::Account":                            &certificatemanager.Account{},
		"AWS::CertificateManager::Certificate":                        &certificatemanager.Certificate{},
		"AWS::Chatbot::SlackChannelConfiguration":                     &chatbot.SlackChannelConfiguration{},
		"AWS::Cloud9::EnvironmentEC2":                                 &cloud9.EnvironmentEC2{},
		"AWS::CloudFormation::CustomResource":                         &cloudformation.CustomResource{},
		"AWS::CloudFormation::Macro":                                  &cloudformation.Macro{},
		"AWS::CloudFormation::ModuleDefaultVersion":                   &cloudformation.ModuleDefaultVersion{},
		"AWS::CloudFormation::ModuleVersion":                          &cloudformation.ModuleVersion{},
		"AWS::CloudFormation::PublicTypeVersion":                      &cloudformation.PublicTypeVersion{},
		"AWS::CloudFormation::Publisher":                              &cloudformation.Publisher{},
		"AWS::CloudFormation::ResourceDefaultVersion":                 &cloudformation.ResourceDefaultVersion{},
		"AWS::CloudFormation::ResourceVersion":                        &cloudformation.ResourceVersion{},
		"AWS::CloudFormation::Stack":                                  &cloudformation.Stack{},
		"AWS::CloudFormation::StackSet":                               &cloudformation.StackSet{},
		"AWS::CloudFormation::TypeActivation":                         &cloudformation.TypeActivation{},
		"AWS::CloudFormation::WaitCondition":                          &cloudformation.WaitCondition{},
		"AWS::CloudFormation::WaitConditionHandle":                    &cloudformation.WaitConditionHandle{},
		"AWS::CloudFront::CachePolicy":                                &cloudfront.CachePolicy{},
		"AWS::CloudFront::CloudFrontOriginAccessIdentity":             &cloudfront.CloudFrontOriginAccessIdentity{},
		"AWS::CloudFront::Distribution":                               &cloudfront.Distribution{},
		"AWS::CloudFront::Function":                                   &cloudfront.Function{},
		"AWS::CloudFront::KeyGroup":                                   &cloudfront.KeyGroup{},
		"AWS::CloudFront::OriginRequestPolicy":                        &cloudfront.OriginRequestPolicy{},
		"AWS::CloudFront::PublicKey":                                  &cloudfront.PublicKey{},
		"AWS::CloudFront::RealtimeLogConfig":                          &cloudfront.RealtimeLogConfig{},
		"AWS::CloudFront::ResponseHeadersPolicy":                      &cloudfront.ResponseHeadersPolicy{},
		"AWS::CloudFront::StreamingDistribution":                      &cloudfront.StreamingDistribution{},
		"AWS::CloudTrail::Trail":                                      &cloudtrail.Trail{},
		"AWS::CloudWatch::Alarm":                                      &cloudwatch.Alarm{},
		"AWS::CloudWatch::AnomalyDetector":                            &cloudwatch.AnomalyDetector{},
		"AWS::CloudWatch::CompositeAlarm":                             &cloudwatch.CompositeAlarm{},
		"AWS::CloudWatch::Dashboard":                                  &cloudwatch.Dashboard{},
		"AWS::CloudWatch::InsightRule":                                &cloudwatch.InsightRule{},
		"AWS::CloudWatch::MetricStream":                               &cloudwatch.MetricStream{},
		"AWS::CodeArtifact::Domain":                                   &codeartifact.Domain{},
		"AWS::CodeArtifact::Repository":                               &codeartifact.Repository{},
		"AWS::CodeBuild::Project":                                     &codebuild.Project{},
		"AWS::CodeBuild::ReportGroup":                                 &codebuild.ReportGroup{},
		"AWS::CodeBuild::SourceCredential":                            &codebuild.SourceCredential{},
		"AWS::CodeCommit::Repository":                                 &codecommit.Repository{},
		"AWS::CodeDeploy::Application":                                &codedeploy.Application{},
		"AWS::CodeDeploy::DeploymentConfig":                           &codedeploy.DeploymentConfig{},
		"AWS::CodeDeploy::DeploymentGroup":                            &codedeploy.DeploymentGroup{},
		"AWS::CodeGuruProfiler::ProfilingGroup":                       &codeguruprofiler.ProfilingGroup{},
		"AWS::CodeGuruReviewer::RepositoryAssociation":                &codegurureviewer.RepositoryAssociation{},
		"AWS::CodePipeline::CustomActionType":                         &codepipeline.CustomActionType{},
		"AWS::CodePipeline::Pipeline":                                 &codepipeline.Pipeline{},
		"AWS::CodePipeline::Webhook":                                  &codepipeline.Webhook{},
		"AWS::CodeStar::GitHubRepository":                             &codestar.GitHubRepository{},
		"AWS::CodeStarConnections::Connection":                        &codestarconnections.Connection{},
		"AWS::CodeStarNotifications::NotificationRule":                &codestarnotifications.NotificationRule{},
		"AWS::Cognito::IdentityPool":                                  &cognito.IdentityPool{},
		"AWS::Cognito::IdentityPoolRoleAttachment":                    &cognito.IdentityPoolRoleAttachment{},
		"AWS::Cognito::UserPool":                                      &cognito.UserPool{},
		"AWS::Cognito::UserPoolClient":                                &cognito.UserPoolClient{},
		"AWS::Cognito::UserPoolDomain":                                &cognito.UserPoolDomain{},
		"AWS::Cognito::UserPoolGroup":                                 &cognito.UserPoolGroup{},
		"AWS::Cognito::UserPoolIdentityProvider":                      &cognito.UserPoolIdentityProvider{},
		"AWS::Cognito::UserPoolResourceServer":                        &cognito.UserPoolResourceServer{},
		"AWS::Cognito::UserPoolRiskConfigurationAttachment":           &cognito.UserPoolRiskConfigurationAttachment{},
		"AWS::Cognito::UserPoolUICustomizationAttachment":             &cognito.UserPoolUICustomizationAttachment{},
		"AWS::Cognito::UserPoolUser":                                  &cognito.UserPoolUser{},
		"AWS::Cognito::UserPoolUserToGroupAttachment":                 &cognito.UserPoolUserToGroupAttachment{},
		"AWS::Config::AggregationAuthorization":                       &config.AggregationAuthorization{},
		"AWS::Config::ConfigRule":                                     &config.ConfigRule{},
		"AWS::Config::ConfigurationAggregator":                        &config.ConfigurationAggregator{},
		"AWS::Config::ConfigurationRecorder":                          &config.ConfigurationRecorder{},
		"AWS::Config::ConformancePack":                                &config.ConformancePack{},
		"AWS::Config::DeliveryChannel":                                &config.DeliveryChannel{},
		"AWS::Config::OrganizationConfigRule":                         &config.OrganizationConfigRule{},
		"AWS::Config::OrganizationConformancePack":                    &config.OrganizationConformancePack{},
		"AWS::Config::RemediationConfiguration":                       &config.RemediationConfiguration{},
		"AWS::Config::StoredQuery":                                    &config.StoredQuery{},
		"AWS::Connect::ContactFlow":                                   &connect.ContactFlow{},
		"AWS::Connect::ContactFlowModule":                             &connect.ContactFlowModule{},
		"AWS::Connect::HoursOfOperation":                              &connect.HoursOfOperation{},
		"AWS::Connect::QuickConnect":                                  &connect.QuickConnect{},
		"AWS::Connect::User":                                          &connect.User{},
		"AWS::Connect::UserHierarchyGroup":                            &connect.UserHierarchyGroup{},
		"AWS::CustomerProfiles::Domain":                               &customerprofiles.Domain{},
		"AWS::CustomerProfiles::Integration":                          &customerprofiles.Integration{},
		"AWS::CustomerProfiles::ObjectType":                           &customerprofiles.ObjectType{},
		"AWS::DAX::Cluster":                                           &dax.Cluster{},
		"AWS::DAX::ParameterGroup":                                    &dax.ParameterGroup{},
		"AWS::DAX::SubnetGroup":                                       &dax.SubnetGroup{},
		"AWS::DLM::LifecyclePolicy":                                   &dlm.LifecyclePolicy{},
		"AWS::DMS::Certificate":                                       &dms.Certificate{},
		"AWS::DMS::Endpoint":                                          &dms.Endpoint{},
		"AWS::DMS::EventSubscription":                                 &dms.EventSubscription{},
		"AWS::DMS::ReplicationInstance":                               &dms.ReplicationInstance{},
		"AWS::DMS::ReplicationSubnetGroup":                            &dms.ReplicationSubnetGroup{},
		"AWS::DMS::ReplicationTask":                                   &dms.ReplicationTask{},
		"AWS::DataBrew::Dataset":                                      &databrew.Dataset{},
		"AWS::DataBrew::Job":                                          &databrew.Job{},
		"AWS::DataBrew::Project":                                      &databrew.Project{},
		"AWS::DataBrew::Recipe":                                       &databrew.Recipe{},
		"AWS::DataBrew::Ruleset":                                      &databrew.Ruleset{},
		"AWS::DataBrew::Schedule":                                     &databrew.Schedule{},
		"AWS::DataPipeline::Pipeline":                                 &datapipeline.Pipeline{},
		"AWS::DataSync::Agent":                                        &datasync.Agent{},
		"AWS::DataSync::LocationEFS":                                  &datasync.LocationEFS{},
		"AWS::DataSync::LocationFSxWindows":                           &datasync.LocationFSxWindows{},
		"AWS::DataSync::LocationHDFS":                                 &datasync.LocationHDFS{},
		"AWS::DataSync::LocationNFS":                                  &datasync.LocationNFS{},
		"AWS::DataSync::LocationObjectStorage":                        &datasync.LocationObjectStorage{},
		"AWS::DataSync::LocationS3":                                   &datasync.LocationS3{},
		"AWS::DataSync::LocationSMB":                                  &datasync.LocationSMB{},
		"AWS::DataSync::Task":                                         &datasync.Task{},
		"AWS::Detective::Graph":                                       &detective.Graph{},
		"AWS::Detective::MemberInvitation":                            &detective.MemberInvitation{},
		"AWS::DevOpsGuru::NotificationChannel":                        &devopsguru.NotificationChannel{},
		"AWS::DevOpsGuru::ResourceCollection":                         &devopsguru.ResourceCollection{},
		"AWS::DirectoryService::MicrosoftAD":                          &directoryservice.MicrosoftAD{},
		"AWS::DirectoryService::SimpleAD":                             &directoryservice.SimpleAD{},
		"AWS::DocDB::DBCluster":                                       &docdb.DBCluster{},
		"AWS::DocDB::DBClusterParameterGroup":                         &docdb.DBClusterParameterGroup{},
		"AWS::DocDB::DBInstance":                                      &docdb.DBInstance{},
		"AWS::DocDB::DBSubnetGroup":                                   &docdb.DBSubnetGroup{},
		"AWS::DynamoDB::GlobalTable":                                  &dynamodb.GlobalTable{},
		"AWS::DynamoDB::Table":                                        &dynamodb.Table{},
		"AWS::EC2::CapacityReservation":                               &ec2.CapacityReservation{},
		"AWS::EC2::CapacityReservationFleet":                          &ec2.CapacityReservationFleet{},
		"AWS::EC2::CarrierGateway":                                    &ec2.CarrierGateway{},
		"AWS::EC2::ClientVpnAuthorizationRule":                        &ec2.ClientVpnAuthorizationRule{},
		"AWS::EC2::ClientVpnEndpoint":                                 &ec2.ClientVpnEndpoint{},
		"AWS::EC2::ClientVpnRoute":                                    &ec2.ClientVpnRoute{},
		"AWS::EC2::ClientVpnTargetNetworkAssociation":                 &ec2.ClientVpnTargetNetworkAssociation{},
		"AWS::EC2::CustomerGateway":                                   &ec2.CustomerGateway{},
		"AWS::EC2::DHCPOptions":                                       &ec2.DHCPOptions{},
		"AWS::EC2::EC2Fleet":                                          &ec2.EC2Fleet{},
		"AWS::EC2::EIP":                                               &ec2.EIP{},
		"AWS::EC2::EIPAssociation":                                    &ec2.EIPAssociation{},
		"AWS::EC2::EgressOnlyInternetGateway":                         &ec2.EgressOnlyInternetGateway{},
		"AWS::EC2::EnclaveCertificateIamRoleAssociation":              &ec2.EnclaveCertificateIamRoleAssociation{},
		"AWS::EC2::FlowLog":                                           &ec2.FlowLog{},
		"AWS::EC2::GatewayRouteTableAssociation":                      &ec2.GatewayRouteTableAssociation{},
		"AWS::EC2::Host":                                              &ec2.Host{},
		"AWS::EC2::IPAM":                                              &ec2.IPAM{},
		"AWS::EC2::IPAMAllocation":                                    &ec2.IPAMAllocation{},
		"AWS::EC2::IPAMPool":                                          &ec2.IPAMPool{},
		"AWS::EC2::IPAMScope":                                         &ec2.IPAMScope{},
		"AWS::EC2::Instance":                                          &ec2.Instance{},
		"AWS::EC2::InternetGateway":                                   &ec2.InternetGateway{},
		"AWS::EC2::LaunchTemplate":                                    &ec2.LaunchTemplate{},
		"AWS::EC2::LocalGatewayRoute":                                 &ec2.LocalGatewayRoute{},
		"AWS::EC2::LocalGatewayRouteTableVPCAssociation":              &ec2.LocalGatewayRouteTableVPCAssociation{},
		"AWS::EC2::NatGateway":                                        &ec2.NatGateway{},
		"AWS::EC2::NetworkAcl":                                        &ec2.NetworkAcl{},
		"AWS::EC2::NetworkAclEntry":                                   &ec2.NetworkAclEntry{},
		"AWS::EC2::NetworkInsightsAnalysis":                           &ec2.NetworkInsightsAnalysis{},
		"AWS::EC2::NetworkInsightsPath":                               &ec2.NetworkInsightsPath{},
		"AWS::EC2::NetworkInterface":                                  &ec2.NetworkInterface{},
		"AWS::EC2::NetworkInterfaceAttachment":                        &ec2.NetworkInterfaceAttachment{},
		"AWS::EC2::NetworkInterfacePermission":                        &ec2.NetworkInterfacePermission{},
		"AWS::EC2::PlacementGroup":                                    &ec2.PlacementGroup{},
		"AWS::EC2::PrefixList":                                        &ec2.PrefixList{},
		"AWS::EC2::Route":                                             &ec2.Route{},
		"AWS::EC2::RouteTable":                                        &ec2.RouteTable{},
		"AWS::EC2::SecurityGroup":                                     &ec2.SecurityGroup{},
		"AWS::EC2::SecurityGroupEgress":                               &ec2.SecurityGroupEgress{},
		"AWS::EC2::SecurityGroupIngress":                              &ec2.SecurityGroupIngress{},
		"AWS::EC2::SpotFleet":                                         &ec2.SpotFleet{},
		"AWS::EC2::Subnet":                                            &ec2.Subnet{},
		"AWS::EC2::SubnetCidrBlock":                                   &ec2.SubnetCidrBlock{},
		"AWS::EC2::SubnetNetworkAclAssociation":                       &ec2.SubnetNetworkAclAssociation{},
		"AWS::EC2::SubnetRouteTableAssociation":                       &ec2.SubnetRouteTableAssociation{},
		"AWS::EC2::TrafficMirrorFilter":                               &ec2.TrafficMirrorFilter{},
		"AWS::EC2::TrafficMirrorFilterRule":                           &ec2.TrafficMirrorFilterRule{},
		"AWS::EC2::TrafficMirrorSession":                              &ec2.TrafficMirrorSession{},
		"AWS::EC2::TrafficMirrorTarget":                               &ec2.TrafficMirrorTarget{},
		"AWS::EC2::TransitGateway":                                    &ec2.TransitGateway{},
		"AWS::EC2::TransitGatewayAttachment":                          &ec2.TransitGatewayAttachment{},
		"AWS::EC2::TransitGatewayConnect":                             &ec2.TransitGatewayConnect{},
		"AWS::EC2::TransitGatewayMulticastDomain":                     &ec2.TransitGatewayMulticastDomain{},
		"AWS::EC2::TransitGatewayMulticastDomainAssociation":          &ec2.TransitGatewayMulticastDomainAssociation{},
		"AWS::EC2::TransitGatewayMulticastGroupMember":                &ec2.TransitGatewayMulticastGroupMember{},
		"AWS::EC2::TransitGatewayMulticastGroupSource":                &ec2.TransitGatewayMulticastGroupSource{},
		"AWS::EC2::TransitGatewayPeeringAttachment":                   &ec2.TransitGatewayPeeringAttachment{},
		"AWS::EC2::TransitGatewayRoute":                               &ec2.TransitGatewayRoute{},
		"AWS::EC2::TransitGatewayRouteTable":                          &ec2.TransitGatewayRouteTable{},
		"AWS::EC2::TransitGatewayRouteTableAssociation":               &ec2.TransitGatewayRouteTableAssociation{},
		"AWS::EC2::TransitGatewayRouteTablePropagation":               &ec2.TransitGatewayRouteTablePropagation{},
		"AWS::EC2::TransitGatewayVpcAttachment":                       &ec2.TransitGatewayVpcAttachment{},
		"AWS::EC2::VPC":                                               &ec2.VPC{},
		"AWS::EC2::VPCCidrBlock":                                      &ec2.VPCCidrBlock{},
		"AWS::EC2::VPCDHCPOptionsAssociation":                         &ec2.VPCDHCPOptionsAssociation{},
		"AWS::EC2::VPCEndpoint":                                       &ec2.VPCEndpoint{},
		"AWS::EC2::VPCEndpointConnectionNotification":                 &ec2.VPCEndpointConnectionNotification{},
		"AWS::EC2::VPCEndpointService":                                &ec2.VPCEndpointService{},
		"AWS::EC2::VPCEndpointServicePermissions":                     &ec2.VPCEndpointServicePermissions{},
		"AWS::EC2::VPCGatewayAttachment":                              &ec2.VPCGatewayAttachment{},
		"AWS::EC2::VPCPeeringConnection":                              &ec2.VPCPeeringConnection{},
		"AWS::EC2::VPNConnection":                                     &ec2.VPNConnection{},
		"AWS::EC2::VPNConnectionRoute":                                &ec2.VPNConnectionRoute{},
		"AWS::EC2::VPNGateway":                                        &ec2.VPNGateway{},
		"AWS::EC2::VPNGatewayRoutePropagation":                        &ec2.VPNGatewayRoutePropagation{},
		"AWS::EC2::Volume":                                            &ec2.Volume{},
		"AWS::EC2::VolumeAttachment":                                  &ec2.VolumeAttachment{},
		"AWS::ECR::PublicRepository":                                  &ecr.PublicRepository{},
		"AWS::ECR::RegistryPolicy":                                    &ecr.RegistryPolicy{},
		"AWS::ECR::ReplicationConfiguration":                          &ecr.ReplicationConfiguration{},
		"AWS::ECR::Repository":                                        &ecr.Repository{},
		"AWS::ECS::CapacityProvider":                                  &ecs.CapacityProvider{},
		"AWS::ECS::Cluster":                                           &ecs.Cluster{},
		"AWS::ECS::ClusterCapacityProviderAssociations":               &ecs.ClusterCapacityProviderAssociations{},
		"AWS::ECS::PrimaryTaskSet":                                    &ecs.PrimaryTaskSet{},
		"AWS::ECS::Service":                                           &ecs.Service{},
		"AWS::ECS::TaskDefinition":                                    &ecs.TaskDefinition{},
		"AWS::ECS::TaskSet":                                           &ecs.TaskSet{},
		"AWS::EFS::AccessPoint":                                       &efs.AccessPoint{},
		"AWS::EFS::FileSystem":                                        &efs.FileSystem{},
		"AWS::EFS::MountTarget":                                       &efs.MountTarget{},
		"AWS::EKS::Addon":                                             &eks.Addon{},
		"AWS::EKS::Cluster":                                           &eks.Cluster{},
		"AWS::EKS::AccessEntry":                                       &eks.AccessEntry{},
		"AWS::EKS::FargateProfile":                                    &eks.FargateProfile{},
		"AWS::EKS::Nodegroup":                                         &eks.Nodegroup{},
		"AWS::EMR::Cluster":                                           &emr.Cluster{},
		"AWS::EMR::InstanceFleetConfig":                               &emr.InstanceFleetConfig{},
		"AWS::EMR::InstanceGroupConfig":                               &emr.InstanceGroupConfig{},
		"AWS::EMR::SecurityConfiguration":                             &emr.SecurityConfiguration{},
		"AWS::EMR::Step":                                              &emr.Step{},
		"AWS::EMR::Studio":                                            &emr.Studio{},
		"AWS::EMR::StudioSessionMapping":                              &emr.StudioSessionMapping{},
		"AWS::EMRContainers::VirtualCluster":                          &emrcontainers.VirtualCluster{},
		"AWS::ElastiCache::CacheCluster":                              &elasticache.CacheCluster{},
		"AWS::ElastiCache::GlobalReplicationGroup":                    &elasticache.GlobalReplicationGroup{},
		"AWS::ElastiCache::ParameterGroup":                            &elasticache.ParameterGroup{},
		"AWS::ElastiCache::ReplicationGroup":                          &elasticache.ReplicationGroup{},
		"AWS::ElastiCache::SecurityGroup":                             &elasticache.SecurityGroup{},
		"AWS::ElastiCache::SecurityGroupIngress":                      &elasticache.SecurityGroupIngress{},
		"AWS::ElastiCache::SubnetGroup":                               &elasticache.SubnetGroup{},
		"AWS::ElastiCache::User":                                      &elasticache.User{},
		"AWS::ElastiCache::UserGroup":                                 &elasticache.UserGroup{},
		"AWS::ElasticBeanstalk::Application":                          &elasticbeanstalk.Application{},
		"AWS::ElasticBeanstalk::ApplicationVersion":                   &elasticbeanstalk.ApplicationVersion{},
		"AWS::ElasticBeanstalk::ConfigurationTemplate":                &elasticbeanstalk.ConfigurationTemplate{},
		"AWS::ElasticBeanstalk::Environment":                          &elasticbeanstalk.Environment{},
		"AWS::ElasticLoadBalancing::LoadBalancer":                     &elasticloadbalancing.LoadBalancer{},
		"AWS::ElasticLoadBalancingV2::Listener":                       &elasticloadbalancingv2.Listener{},
		"AWS::ElasticLoadBalancingV2::ListenerCertificate":            &elasticloadbalancingv2.ListenerCertificate{},
		"AWS::ElasticLoadBalancingV2::ListenerRule":                   &elasticloadbalancingv2.ListenerRule{},
		"AWS::ElasticLoadBalancingV2::LoadBalancer":                   &elasticloadbalancingv2.LoadBalancer{},
		"AWS::ElasticLoadBalancingV2::TargetGroup":                    &elasticloadbalancingv2.TargetGroup{},
		"AWS::Elasticsearch::Domain":                                  &elasticsearch.Domain{},
		"AWS::EventSchemas::Discoverer":                               &eventschemas.Discoverer{},
		"AWS::EventSchemas::Registry":                                 &eventschemas.Registry{},
		"AWS::EventSchemas::RegistryPolicy":                           &eventschemas.RegistryPolicy{},
		"AWS::EventSchemas::Schema":                                   &eventschemas.Schema{},
		"AWS::Events::ApiDestination":                                 &events.ApiDestination{},
		"AWS::Events::Archive":                                        &events.Archive{},
		"AWS::Events::Connection":                                     &events.Connection{},
		"AWS::Events::EventBus":                                       &events.EventBus{},
		"AWS::Events::EventBusPolicy":                                 &events.EventBusPolicy{},
		"AWS::Events::Rule":                                           &events.Rule{},
		"AWS::Evidently::Experiment":                                  &evidently.Experiment{},
		"AWS::Evidently::Feature":                                     &evidently.Feature{},
		"AWS::Evidently::Launch":                                      &evidently.Launch{},
		"AWS::Evidently::Project":                                     &evidently.Project{},
		"AWS::FIS::ExperimentTemplate":                                &fis.ExperimentTemplate{},
		"AWS::FMS::NotificationChannel":                               &fms.NotificationChannel{},
		"AWS::FMS::Policy":                                            &fms.Policy{},
		"AWS::FSx::FileSystem":                                        &fsx.FileSystem{},
		"AWS::FinSpace::Environment":                                  &finspace.Environment{},
		"AWS::FraudDetector::Detector":                                &frauddetector.Detector{},
		"AWS::FraudDetector::EntityType":                              &frauddetector.EntityType{},
		"AWS::FraudDetector::EventType":                               &frauddetector.EventType{},
		"AWS::FraudDetector::Label":                                   &frauddetector.Label{},
		"AWS::FraudDetector::Outcome":                                 &frauddetector.Outcome{},
		"AWS::FraudDetector::Variable":                                &frauddetector.Variable{},
		"AWS::GameLift::Alias":                                        &gamelift.Alias{},
		"AWS::GameLift::Build":                                        &gamelift.Build{},
		"AWS::GameLift::Fleet":                                        &gamelift.Fleet{},
		"AWS::GameLift::GameServerGroup":                              &gamelift.GameServerGroup{},
		"AWS::GameLift::GameSessionQueue":                             &gamelift.GameSessionQueue{},
		"AWS::GameLift::MatchmakingConfiguration":                     &gamelift.MatchmakingConfiguration{},
		"AWS::GameLift::MatchmakingRuleSet":                           &gamelift.MatchmakingRuleSet{},
		"AWS::GameLift::Script":                                       &gamelift.Script{},
		"AWS::GlobalAccelerator::Accelerator":                         &globalaccelerator.Accelerator{},
		"AWS::GlobalAccelerator::EndpointGroup":                       &globalaccelerator.EndpointGroup{},
		"AWS::GlobalAccelerator::Listener":                            &globalaccelerator.Listener{},
		"AWS::Glue::Classifier":                                       &glue.Classifier{},
		"AWS::Glue::Connection":                                       &glue.Connection{},
		"AWS::Glue::Crawler":                                          &glue.Crawler{},
		"AWS::Glue::DataCatalogEncryptionSettings":                    &glue.DataCatalogEncryptionSettings{},
		"AWS::Glue::Database":                                         &glue.Database{},
		"AWS::Glue::DevEndpoint":                                      &glue.DevEndpoint{},
		"AWS::Glue::Job":                                              &glue.Job{},
		"AWS::Glue::MLTransform":                                      &glue.MLTransform{},
		"AWS::Glue::Partition":                                        &glue.Partition{},
		"AWS::Glue::Registry":                                         &glue.Registry{},
		"AWS::Glue::Schema":                                           &glue.Schema{},
		"AWS::Glue::SchemaVersion":                                    &glue.SchemaVersion{},
		"AWS::Glue::SchemaVersionMetadata":                            &glue.SchemaVersionMetadata{},
		"AWS::Glue::SecurityConfiguration":                            &glue.SecurityConfiguration{},
		"AWS::Glue::Table":                                            &glue.Table{},
		"AWS::Glue::Trigger":                                          &glue.Trigger{},
		"AWS::Glue::Workflow":                                         &glue.Workflow{},
		"AWS::Greengrass::ConnectorDefinition":                        &greengrass.ConnectorDefinition{},
		"AWS::Greengrass::ConnectorDefinitionVersion":                 &greengrass.ConnectorDefinitionVersion{},
		"AWS::Greengrass::CoreDefinition":                             &greengrass.CoreDefinition{},
		"AWS::Greengrass::CoreDefinitionVersion":                      &greengrass.CoreDefinitionVersion{},
		"AWS::Greengrass::DeviceDefinition":                           &greengrass.DeviceDefinition{},
		"AWS::Greengrass::DeviceDefinitionVersion":                    &greengrass.DeviceDefinitionVersion{},
		"AWS::Greengrass::FunctionDefinition":                         &greengrass.FunctionDefinition{},
		"AWS::Greengrass::FunctionDefinitionVersion":                  &greengrass.FunctionDefinitionVersion{},
		"AWS::Greengrass::Group":                                      &greengrass.Group{},
		"AWS::Greengrass::GroupVersion":                               &greengrass.GroupVersion{},
		"AWS::Greengrass::LoggerDefinition":                           &greengrass.LoggerDefinition{},
		"AWS::Greengrass::LoggerDefinitionVersion":                    &greengrass.LoggerDefinitionVersion{},
		"AWS::Greengrass::ResourceDefinition":                         &greengrass.ResourceDefinition{},
		"AWS::Greengrass::ResourceDefinitionVersion":                  &greengrass.ResourceDefinitionVersion{},
		"AWS::Greengrass::SubscriptionDefinition":                     &greengrass.SubscriptionDefinition{},
		"AWS::Greengrass::SubscriptionDefinitionVersion":              &greengrass.SubscriptionDefinitionVersion{},
		"AWS::GreengrassV2::ComponentVersion":                         &greengrassv2.ComponentVersion{},
		"AWS::GroundStation::Config":                                  &groundstation.Config{},
		"AWS::GroundStation::DataflowEndpointGroup":                   &groundstation.DataflowEndpointGroup{},
		"AWS::GroundStation::MissionProfile":                          &groundstation.MissionProfile{},
		"AWS::GuardDuty::Detector":                                    &guardduty.Detector{},
		"AWS::GuardDuty::Filter":                                      &guardduty.Filter{},
		"AWS::GuardDuty::IPSet":                                       &guardduty.IPSet{},
		"AWS::GuardDuty::Master":                                      &guardduty.Master{},
		"AWS::GuardDuty::Member":                                      &guardduty.Member{},
		"AWS::GuardDuty::ThreatIntelSet":                              &guardduty.ThreatIntelSet{},
		"AWS::HealthLake::FHIRDatastore":                              &healthlake.FHIRDatastore{},
		"AWS::IAM::AccessKey":                                         &iam.AccessKey{},
		"AWS::IAM::Group":                                             &iam.Group{},
		"AWS::IAM::InstanceProfile":                                   &iam.InstanceProfile{},
		"AWS::IAM::ManagedPolicy":                                     &iam.ManagedPolicy{},
		"AWS::IAM::OIDCProvider":                                      &iam.OIDCProvider{},
		"AWS::IAM::Policy":                                            &iam.Policy{},
		"AWS::IAM::Role":                                              &iam.Role{},
		"AWS::IAM::SAMLProvider":                                      &iam.SAMLProvider{},
		"AWS::IAM::ServerCertificate":                                 &iam.ServerCertificate{},
		"AWS::IAM::ServiceLinkedRole":                                 &iam.ServiceLinkedRole{},
		"AWS::IAM::User":                                              &iam.User{},
		"AWS::IAM::UserToGroupAddition":                               &iam.UserToGroupAddition{},
		"AWS::IAM::VirtualMFADevice":                                  &iam.VirtualMFADevice{},
		"AWS::IVS::Channel":                                           &ivs.Channel{},
		"AWS::IVS::PlaybackKeyPair":                                   &ivs.PlaybackKeyPair{},
		"AWS::IVS::RecordingConfiguration":                            &ivs.RecordingConfiguration{},
		"AWS::IVS::StreamKey":                                         &ivs.StreamKey{},
		"AWS::ImageBuilder::Component":                                &imagebuilder.Component{},
		"AWS::ImageBuilder::ContainerRecipe":                          &imagebuilder.ContainerRecipe{},
		"AWS::ImageBuilder::DistributionConfiguration":                &imagebuilder.DistributionConfiguration{},
		"AWS::ImageBuilder::Image":                                    &imagebuilder.Image{},
		"AWS::ImageBuilder::ImagePipeline":                            &imagebuilder.ImagePipeline{},
		"AWS::ImageBuilder::ImageRecipe":                              &imagebuilder.ImageRecipe{},
		"AWS::ImageBuilder::InfrastructureConfiguration":              &imagebuilder.InfrastructureConfiguration{},
		"AWS::Inspector::AssessmentTarget":                            &inspector.AssessmentTarget{},
		"AWS::Inspector::AssessmentTemplate":                          &inspector.AssessmentTemplate{},
		"AWS::Inspector::ResourceGroup":                               &inspector.ResourceGroup{},
		"AWS::IoT1Click::Device":                                      &iot1click.Device{},
		"AWS::IoT1Click::Placement":                                   &iot1click.Placement{},
		"AWS::IoT1Click::Project":                                     &iot1click.Project{},
		"AWS::IoT::AccountAuditConfiguration":                         &iot.AccountAuditConfiguration{},
		"AWS::IoT::Authorizer":                                        &iot.Authorizer{},
		"AWS::IoT::Certificate":                                       &iot.Certificate{},
		"AWS::IoT::CustomMetric":                                      &iot.CustomMetric{},
		"AWS::IoT::Dimension":                                         &iot.Dimension{},
		"AWS::IoT::DomainConfiguration":                               &iot.DomainConfiguration{},
		"AWS::IoT::FleetMetric":                                       &iot.FleetMetric{},
		"AWS::IoT::JobTemplate":                                       &iot.JobTemplate{},
		"AWS::IoT::Logging":                                           &iot.Logging{},
		"AWS::IoT::MitigationAction":                                  &iot.MitigationAction{},
		"AWS::IoT::Policy":                                            &iot.Policy{},
		"AWS::IoT::PolicyPrincipalAttachment":                         &iot.PolicyPrincipalAttachment{},
		"AWS::IoT::ProvisioningTemplate":                              &iot.ProvisioningTemplate{},
		"AWS::IoT::ResourceSpecificLogging":                           &iot.ResourceSpecificLogging{},
		"AWS::IoT::ScheduledAudit":                                    &iot.ScheduledAudit{},
		"AWS::IoT::SecurityProfile":                                   &iot.SecurityProfile{},
		"AWS::IoT::Thing":                                             &iot.Thing{},
		"AWS::IoT::ThingPrincipalAttachment":                          &iot.ThingPrincipalAttachment{},
		"AWS::IoT::TopicRule":                                         &iot.TopicRule{},
		"AWS::IoT::TopicRuleDestination":                              &iot.TopicRuleDestination{},
		"AWS::IoTAnalytics::Channel":                                  &iotanalytics.Channel{},
		"AWS::IoTAnalytics::Dataset":                                  &iotanalytics.Dataset{},
		"AWS::IoTAnalytics::Datastore":                                &iotanalytics.Datastore{},
		"AWS::IoTAnalytics::Pipeline":                                 &iotanalytics.Pipeline{},
		"AWS::IoTCoreDeviceAdvisor::SuiteDefinition":                  &iotcoredeviceadvisor.SuiteDefinition{},
		"AWS::IoTEvents::DetectorModel":                               &iotevents.DetectorModel{},
		"AWS::IoTEvents::Input":                                       &iotevents.Input{},
		"AWS::IoTFleetHub::Application":                               &iotfleethub.Application{},
		"AWS::IoTSiteWise::AccessPolicy":                              &iotsitewise.AccessPolicy{},
		"AWS::IoTSiteWise::Asset":                                     &iotsitewise.Asset{},
		"AWS::IoTSiteWise::AssetModel":                                &iotsitewise.AssetModel{},
		"AWS::IoTSiteWise::Dashboard":                                 &iotsitewise.Dashboard{},
		"AWS::IoTSiteWise::Gateway":                                   &iotsitewise.Gateway{},
		"AWS::IoTSiteWise::Portal":                                    &iotsitewise.Portal{},
		"AWS::IoTSiteWise::Project":                                   &iotsitewise.Project{},
		"AWS::IoTThingsGraph::FlowTemplate":                           &iotthingsgraph.FlowTemplate{},
		"AWS::IoTWireless::Destination":                               &iotwireless.Destination{},
		"AWS::IoTWireless::DeviceProfile":                             &iotwireless.DeviceProfile{},
		"AWS::IoTWireless::FuotaTask":                                 &iotwireless.FuotaTask{},
		"AWS::IoTWireless::MulticastGroup":                            &iotwireless.MulticastGroup{},
		"AWS::IoTWireless::PartnerAccount":                            &iotwireless.PartnerAccount{},
		"AWS::IoTWireless::ServiceProfile":                            &iotwireless.ServiceProfile{},
		"AWS::IoTWireless::TaskDefinition":                            &iotwireless.TaskDefinition{},
		"AWS::IoTWireless::WirelessDevice":                            &iotwireless.WirelessDevice{},
		"AWS::IoTWireless::WirelessGateway":                           &iotwireless.WirelessGateway{},
		"AWS::KMS::Alias":                                             &kms.Alias{},
		"AWS::KMS::Key":                                               &kms.Key{},
		"AWS::KMS::ReplicaKey":                                        &kms.ReplicaKey{},
		"AWS::Kendra::DataSource":                                     &kendra.DataSource{},
		"AWS::Kendra::Faq":                                            &kendra.Faq{},
		"AWS::Kendra::Index":                                          &kendra.Index{},
		"AWS::Kinesis::Stream":                                        &kinesis.Stream{},
		"AWS::Kinesis::StreamConsumer":                                &kinesis.StreamConsumer{},
		"AWS::KinesisAnalytics::Application":                          &kinesisanalytics.Application{},
		"AWS::KinesisAnalytics::ApplicationOutput":                    &kinesisanalytics.ApplicationOutput{},
		"AWS::KinesisAnalytics::ApplicationReferenceDataSource":       &kinesisanalytics.ApplicationReferenceDataSource{},
		"AWS::KinesisAnalyticsV2::Application":                        &kinesisanalyticsv2.Application{},
		"AWS::KinesisAnalyticsV2::ApplicationCloudWatchLoggingOption": &kinesisanalyticsv2.ApplicationCloudWatchLoggingOption{},
		"AWS::KinesisAnalyticsV2::ApplicationOutput":                  &kinesisanalyticsv2.ApplicationOutput{},
		"AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource":     &kinesisanalyticsv2.ApplicationReferenceDataSource{},
		"AWS::KinesisFirehose::DeliveryStream":                        &kinesisfirehose.DeliveryStream{},
		"AWS::LakeFormation::DataLakeSettings":                        &lakeformation.DataLakeSettings{},
		"AWS::LakeFormation::Permissions":                             &lakeformation.Permissions{},
		"AWS::LakeFormation::Resource":                                &lakeformation.Resource{},
		"AWS::Lambda::Alias":                                          &lambda.Alias{},
		"AWS::Lambda::CodeSigningConfig":                              &lambda.CodeSigningConfig{},
		"AWS::Lambda::EventInvokeConfig":                              &lambda.EventInvokeConfig{},
		"AWS::Lambda::EventSourceMapping":                             &lambda.EventSourceMapping{},
		"AWS::Lambda::Function":                                       &lambda.Function{},
		"AWS::Lambda::LayerVersion":                                   &lambda.LayerVersion{},
		"AWS::Lambda::LayerVersionPermission":                         &lambda.LayerVersionPermission{},
		"AWS::Lambda::Permission":                                     &lambda.Permission{},
		"AWS::Lambda::Version":                                        &lambda.Version{},
		"AWS::LicenseManager::Grant":                                  &licensemanager.Grant{},
		"AWS::LicenseManager::License":                                &licensemanager.License{},
		"AWS::Lightsail::Database":                                    &lightsail.Database{},
		"AWS::Lightsail::Disk":                                        &lightsail.Disk{},
		"AWS::Lightsail::Instance":                                    &lightsail.Instance{},
		"AWS::Lightsail::StaticIp":                                    &lightsail.StaticIp{},
		"AWS::Location::GeofenceCollection":                           &location.GeofenceCollection{},
		"AWS::Location::Map":                                          &location.Map{},
		"AWS::Location::PlaceIndex":                                   &location.PlaceIndex{},
		"AWS::Location::RouteCalculator":                              &location.RouteCalculator{},
		"AWS::Location::Tracker":                                      &location.Tracker{},
		"AWS::Location::TrackerConsumer":                              &location.TrackerConsumer{},
		"AWS::Logs::Destination":                                      &logs.Destination{},
		"AWS::Logs::LogGroup":                                         &logs.LogGroup{},
		"AWS::Logs::LogStream":                                        &logs.LogStream{},
		"AWS::Logs::MetricFilter":                                     &logs.MetricFilter{},
		"AWS::Logs::QueryDefinition":                                  &logs.QueryDefinition{},
		"AWS::Logs::ResourcePolicy":                                   &logs.ResourcePolicy{},
		"AWS::Logs::SubscriptionFilter":                               &logs.SubscriptionFilter{},
		"AWS::LookoutEquipment::InferenceScheduler":                   &lookoutequipment.InferenceScheduler{},
		"AWS::LookoutMetrics::Alert":                                  &lookoutmetrics.Alert{},
		"AWS::LookoutMetrics::AnomalyDetector":                        &lookoutmetrics.AnomalyDetector{},
		"AWS::LookoutVision::Project":                                 &lookoutvision.Project{},
		"AWS::MSK::Cluster":                                           &msk.Cluster{},
		"AWS::MWAA::Environment":                                      &mwaa.Environment{},
		"AWS::Macie::CustomDataIdentifier":                            &macie.CustomDataIdentifier{},
		"AWS::Macie::FindingsFilter":                                  &macie.FindingsFilter{},
		"AWS::Macie::Session":                                         &macie.Session{},
		"AWS::ManagedBlockchain::Member":                              &managedblockchain.Member{},
		"AWS::ManagedBlockchain::Node":                                &managedblockchain.Node{},
		"AWS::MediaConnect::Flow":                                     &mediaconnect.Flow{},
		"AWS::MediaConnect::FlowEntitlement":                          &mediaconnect.FlowEntitlement{},
		"AWS::MediaConnect::FlowOutput":                               &mediaconnect.FlowOutput{},
		"AWS::MediaConnect::FlowSource":                               &mediaconnect.FlowSource{},
		"AWS::MediaConnect::FlowVpcInterface":                         &mediaconnect.FlowVpcInterface{},
		"AWS::MediaConvert::JobTemplate":                              &mediaconvert.JobTemplate{},
		"AWS::MediaConvert::Preset":                                   &mediaconvert.Preset{},
		"AWS::MediaConvert::Queue":                                    &mediaconvert.Queue{},
		"AWS::MediaLive::Channel":                                     &medialive.Channel{},
		"AWS::MediaLive::Input":                                       &medialive.Input{},
		"AWS::MediaLive::InputSecurityGroup":                          &medialive.InputSecurityGroup{},
		"AWS::MediaPackage::Asset":                                    &mediapackage.Asset{},
		"AWS::MediaPackage::Channel":                                  &mediapackage.Channel{},
		"AWS::MediaPackage::OriginEndpoint":                           &mediapackage.OriginEndpoint{},
		"AWS::MediaPackage::PackagingConfiguration":                   &mediapackage.PackagingConfiguration{},
		"AWS::MediaPackage::PackagingGroup":                           &mediapackage.PackagingGroup{},
		"AWS::MediaStore::Container":                                  &mediastore.Container{},
		"AWS::MemoryDB::ACL":                                          &memorydb.ACL{},
		"AWS::MemoryDB::Cluster":                                      &memorydb.Cluster{},
		"AWS::MemoryDB::ParameterGroup":                               &memorydb.ParameterGroup{},
		"AWS::MemoryDB::SubnetGroup":                                  &memorydb.SubnetGroup{},
		"AWS::MemoryDB::User":                                         &memorydb.User{},
		"AWS::Neptune::DBCluster":                                     &neptune.DBCluster{},
		"AWS::Neptune::DBClusterParameterGroup":                       &neptune.DBClusterParameterGroup{},
		"AWS::Neptune::DBInstance":                                    &neptune.DBInstance{},
		"AWS::Neptune::DBParameterGroup":                              &neptune.DBParameterGroup{},
		"AWS::Neptune::DBSubnetGroup":                                 &neptune.DBSubnetGroup{},
		"AWS::NetworkFirewall::Firewall":                              &networkfirewall.Firewall{},
		"AWS::NetworkFirewall::FirewallPolicy":                        &networkfirewall.FirewallPolicy{},
		"AWS::NetworkFirewall::LoggingConfiguration":                  &networkfirewall.LoggingConfiguration{},
		"AWS::NetworkFirewall::RuleGroup":                             &networkfirewall.RuleGroup{},
		"AWS::NetworkManager::CustomerGatewayAssociation":             &networkmanager.CustomerGatewayAssociation{},
		"AWS::NetworkManager::Device":                                 &networkmanager.Device{},
		"AWS::NetworkManager::GlobalNetwork":                          &networkmanager.GlobalNetwork{},
		"AWS::NetworkManager::Link":                                   &networkmanager.Link{},
		"AWS::NetworkManager::LinkAssociation":                        &networkmanager.LinkAssociation{},
		"AWS::NetworkManager::Site":                                   &networkmanager.Site{},
		"AWS::NetworkManager::TransitGatewayRegistration":             &networkmanager.TransitGatewayRegistration{},
		"AWS::NimbleStudio::LaunchProfile":                            &nimblestudio.LaunchProfile{},
		"AWS::NimbleStudio::StreamingImage":                           &nimblestudio.StreamingImage{},
		"AWS::NimbleStudio::Studio":                                   &nimblestudio.Studio{},
		"AWS::NimbleStudio::StudioComponent":                          &nimblestudio.StudioComponent{},
		"AWS::OpenSearchService::Domain":                              &opensearchservice.Domain{},
		"AWS::OpsWorks::App":                                          &opsworks.App{},
		"AWS::OpsWorks::ElasticLoadBalancerAttachment":                &opsworks.ElasticLoadBalancerAttachment{},
		"AWS::OpsWorks::Instance":                                     &opsworks.Instance{},
		"AWS::OpsWorks::Layer":                                        &opsworks.Layer{},
		"AWS::OpsWorks::Stack":                                        &opsworks.Stack{},
		"AWS::OpsWorks::UserProfile":                                  &opsworks.UserProfile{},
		"AWS::OpsWorks::Volume":                                       &opsworks.Volume{},
		"AWS::OpsWorksCM::Server":                                     &opsworkscm.Server{},
		"AWS::Panorama::ApplicationInstance":                          &panorama.ApplicationInstance{},
		"AWS::Panorama::Package":                                      &panorama.Package{},
		"AWS::Panorama::PackageVersion":                               &panorama.PackageVersion{},
		"AWS::Pinpoint::ADMChannel":                                   &pinpoint.ADMChannel{},
		"AWS::Pinpoint::APNSChannel":                                  &pinpoint.APNSChannel{},
		"AWS::Pinpoint::APNSSandboxChannel":                           &pinpoint.APNSSandboxChannel{},
		"AWS::Pinpoint::APNSVoipChannel":                              &pinpoint.APNSVoipChannel{},
		"AWS::Pinpoint::APNSVoipSandboxChannel":                       &pinpoint.APNSVoipSandboxChannel{},
		"AWS::Pinpoint::App":                                          &pinpoint.App{},
		"AWS::Pinpoint::ApplicationSettings":                          &pinpoint.ApplicationSettings{},
		"AWS::Pinpoint::BaiduChannel":                                 &pinpoint.BaiduChannel{},
		"AWS::Pinpoint::Campaign":                                     &pinpoint.Campaign{},
		"AWS::Pinpoint::EmailChannel":                                 &pinpoint.EmailChannel{},
		"AWS::Pinpoint::EmailTemplate":                                &pinpoint.EmailTemplate{},
		"AWS::Pinpoint::EventStream":                                  &pinpoint.EventStream{},
		"AWS::Pinpoint::GCMChannel":                                   &pinpoint.GCMChannel{},
		"AWS::Pinpoint::InAppTemplate":                                &pinpoint.InAppTemplate{},
		"AWS::Pinpoint::PushTemplate":                                 &pinpoint.PushTemplate{},
		"AWS::Pinpoint::SMSChannel":                                   &pinpoint.SMSChannel{},
		"AWS::Pinpoint::Segment":                                      &pinpoint.Segment{},
		"AWS::Pinpoint::SmsTemplate":                                  &pinpoint.SmsTemplate{},
		"AWS::Pinpoint::VoiceChannel":                                 &pinpoint.VoiceChannel{},
		"AWS::PinpointEmail::ConfigurationSet":                        &pinpointemail.ConfigurationSet{},
		"AWS::PinpointEmail::ConfigurationSetEventDestination":        &pinpointemail.ConfigurationSetEventDestination{},
		"AWS::PinpointEmail::DedicatedIpPool":                         &pinpointemail.DedicatedIpPool{},
		"AWS::PinpointEmail::Identity":                                &pinpointemail.Identity{},
		"AWS::QLDB::Ledger":                                           &qldb.Ledger{},
		"AWS::QLDB::Stream":                                           &qldb.Stream{},
		"AWS::QuickSight::Analysis":                                   &quicksight.Analysis{},
		"AWS::QuickSight::Dashboard":                                  &quicksight.Dashboard{},
		"AWS::QuickSight::DataSet":                                    &quicksight.DataSet{},
		"AWS::QuickSight::DataSource":                                 &quicksight.DataSource{},
		"AWS::QuickSight::Template":                                   &quicksight.Template{},
		"AWS::QuickSight::Theme":                                      &quicksight.Theme{},
		"AWS::RAM::ResourceShare":                                     &ram.ResourceShare{},
		"AWS::RDS::DBCluster":                                         &rds.DBCluster{},
		"AWS::RDS::DBClusterParameterGroup":                           &rds.DBClusterParameterGroup{},
		"AWS::RDS::DBInstance":                                        &rds.DBInstance{},
		"AWS::RDS::DBParameterGroup":                                  &rds.DBParameterGroup{},
		"AWS::RDS::DBProxy":                                           &rds.DBProxy{},
		"AWS::RDS::DBProxyEndpoint":                                   &rds.DBProxyEndpoint{},
		"AWS::RDS::DBProxyTargetGroup":                                &rds.DBProxyTargetGroup{},
		"AWS::RDS::DBSecurityGroup":                                   &rds.DBSecurityGroup{},
		"AWS::RDS::DBSecurityGroupIngress":                            &rds.DBSecurityGroupIngress{},
		"AWS::RDS::DBSubnetGroup":                                     &rds.DBSubnetGroup{},
		"AWS::RDS::EventSubscription":                                 &rds.EventSubscription{},
		"AWS::RDS::GlobalCluster":                                     &rds.GlobalCluster{},
		"AWS::RDS::OptionGroup":                                       &rds.OptionGroup{},
		"AWS::RUM::AppMonitor":                                        &rum.AppMonitor{},
		"AWS::Redshift::Cluster":                                      &redshift.Cluster{},
		"AWS::Redshift::ClusterParameterGroup":                        &redshift.ClusterParameterGroup{},
		"AWS::Redshift::ClusterSecurityGroup":                         &redshift.ClusterSecurityGroup{},
		"AWS::Redshift::ClusterSecurityGroupIngress":                  &redshift.ClusterSecurityGroupIngress{},
		"AWS::Redshift::ClusterSubnetGroup":                           &redshift.ClusterSubnetGroup{},
		"AWS::Redshift::EndpointAccess":                               &redshift.EndpointAccess{},
		"AWS::Redshift::EndpointAuthorization":                        &redshift.EndpointAuthorization{},
		"AWS::Redshift::EventSubscription":                            &redshift.EventSubscription{},
		"AWS::Redshift::ScheduledAction":                              &redshift.ScheduledAction{},
		"AWS::RefactorSpaces::Application":                            &refactorspaces.Application{},
		"AWS::RefactorSpaces::Environment":                            &refactorspaces.Environment{},
		"AWS::RefactorSpaces::Route":                                  &refactorspaces.Route{},
		"AWS::RefactorSpaces::Service":                                &refactorspaces.Service{},
		"AWS::Rekognition::Project":                                   &rekognition.Project{},
		"AWS::ResilienceHub::App":                                     &resiliencehub.App{},
		"AWS::ResilienceHub::ResiliencyPolicy":                        &resiliencehub.ResiliencyPolicy{},
		"AWS::ResourceGroups::Group":                                  &resourcegroups.Group{},
		"AWS::RoboMaker::Fleet":                                       &robomaker.Fleet{},
		"AWS::RoboMaker::Robot":                                       &robomaker.Robot{},
		"AWS::RoboMaker::RobotApplication":                            &robomaker.RobotApplication{},
		"AWS::RoboMaker::RobotApplicationVersion":                     &robomaker.RobotApplicationVersion{},
		"AWS::RoboMaker::SimulationApplication":                       &robomaker.SimulationApplication{},
		"AWS::RoboMaker::SimulationApplicationVersion":                &robomaker.SimulationApplicationVersion{},
		"AWS::Route53::DNSSEC":                                        &route53.DNSSEC{},
		"AWS::Route53::HealthCheck":                                   &route53.HealthCheck{},
		"AWS::Route53::HostedZone":                                    &route53.HostedZone{},
		"AWS::Route53::KeySigningKey":                                 &route53.KeySigningKey{},
		"AWS::Route53::RecordSet":                                     &route53.RecordSet{},
		"AWS::Route53::RecordSetGroup":                                &route53.RecordSetGroup{},
		"AWS::Route53RecoveryControl::Cluster":                        &route53recoverycontrol.Cluster{},
		"AWS::Route53RecoveryControl::ControlPanel":                   &route53recoverycontrol.ControlPanel{},
		"AWS::Route53RecoveryControl::RoutingControl":                 &route53recoverycontrol.RoutingControl{},
		"AWS::Route53RecoveryControl::SafetyRule":                     &route53recoverycontrol.SafetyRule{},
		"AWS::Route53RecoveryReadiness::Cell":                         &route53recoveryreadiness.Cell{},
		"AWS::Route53RecoveryReadiness::ReadinessCheck":               &route53recoveryreadiness.ReadinessCheck{},
		"AWS::Route53RecoveryReadiness::RecoveryGroup":                &route53recoveryreadiness.RecoveryGroup{},
		"AWS::Route53RecoveryReadiness::ResourceSet":                  &route53recoveryreadiness.ResourceSet{},
		"AWS::Route53Resolver::FirewallDomainList":                    &route53resolver.FirewallDomainList{},
		"AWS::Route53Resolver::FirewallRuleGroup":                     &route53resolver.FirewallRuleGroup{},
		"AWS::Route53Resolver::FirewallRuleGroupAssociation":          &route53resolver.FirewallRuleGroupAssociation{},
		"AWS::Route53Resolver::ResolverConfig":                        &route53resolver.ResolverConfig{},
		"AWS::Route53Resolver::ResolverDNSSECConfig":                  &route53resolver.ResolverDNSSECConfig{},
		"AWS::Route53Resolver::ResolverEndpoint":                      &route53resolver.ResolverEndpoint{},
		"AWS::Route53Resolver::ResolverQueryLoggingConfig":            &route53resolver.ResolverQueryLoggingConfig{},
		"AWS::Route53Resolver::ResolverQueryLoggingConfigAssociation": &route53resolver.ResolverQueryLoggingConfigAssociation{},
		"AWS::Route53Resolver::ResolverRule":                          &route53resolver.ResolverRule{},
		"AWS::Route53Resolver::ResolverRuleAssociation":               &route53resolver.ResolverRuleAssociation{},
		"AWS::S3::AccessPoint":                                        &s3.AccessPoint{},
		"AWS::S3::Bucket":                                             &s3.Bucket{},
		"AWS::S3::BucketPolicy":                                       &s3.BucketPolicy{},
		"AWS::S3::MultiRegionAccessPoint":                             &s3.MultiRegionAccessPoint{},
		"AWS::S3::MultiRegionAccessPointPolicy":                       &s3.MultiRegionAccessPointPolicy{},
		"AWS::S3::StorageLens":                                        &s3.StorageLens{},
		"AWS::S3ObjectLambda::AccessPoint":                            &s3objectlambda.AccessPoint{},
		"AWS::S3ObjectLambda::AccessPointPolicy":                      &s3objectlambda.AccessPointPolicy{},
		"AWS::S3Outposts::AccessPoint":                                &s3outposts.AccessPoint{},
		"AWS::S3Outposts::Bucket":                                     &s3outposts.Bucket{},
		"AWS::S3Outposts::BucketPolicy":                               &s3outposts.BucketPolicy{},
		"AWS::S3Outposts::Endpoint":                                   &s3outposts.Endpoint{},
		"AWS::SDB::Domain":                                            &sdb.Domain{},
		"AWS::SES::ConfigurationSet":                                  &ses.ConfigurationSet{},
		"AWS::SES::ConfigurationSetEventDestination":                  &ses.ConfigurationSetEventDestination{},
		"AWS::SES::ContactList":                                       &ses.ContactList{},
		"AWS::SES::ReceiptFilter":                                     &ses.ReceiptFilter{},
		"AWS::SES::ReceiptRule":                                       &ses.ReceiptRule{},
		"AWS::SES::ReceiptRuleSet":                                    &ses.ReceiptRuleSet{},
		"AWS::SES::Template":                                          &ses.Template{},
		"AWS::SNS::Subscription":                                      &sns.Subscription{},
		"AWS::SNS::Topic":                                             &sns.Topic{},
		"AWS::SNS::TopicPolicy":                                       &sns.TopicPolicy{},
		"AWS::SQS::Queue":                                             &sqs.Queue{},
		"AWS::SQS::QueuePolicy":                                       &sqs.QueuePolicy{},
		"AWS::SSM::Association":                                       &ssm.Association{},
		"AWS::SSM::Document":                                          &ssm.Document{},
		"AWS::SSM::MaintenanceWindow":                                 &ssm.MaintenanceWindow{},
		"AWS::SSM::MaintenanceWindowTarget":                           &ssm.MaintenanceWindowTarget{},
		"AWS::SSM::MaintenanceWindowTask":                             &ssm.MaintenanceWindowTask{},
		"AWS::SSM::Parameter":                                         &ssm.Parameter{},
		"AWS::SSM::PatchBaseline":                                     &ssm.PatchBaseline{},
		"AWS::SSM::ResourceDataSync":                                  &ssm.ResourceDataSync{},
		"AWS::SSMContacts::Contact":                                   &ssmcontacts.Contact{},
		"AWS::SSMContacts::ContactChannel":                            &ssmcontacts.ContactChannel{},
		"AWS::SSMIncidents::ReplicationSet":                           &ssmincidents.ReplicationSet{},
		"AWS::SSMIncidents::ResponsePlan":                             &ssmincidents.ResponsePlan{},
		"AWS::SSO::Assignment":                                        &sso.Assignment{},
		"AWS::SSO::InstanceAccessControlAttributeConfiguration":       &sso.InstanceAccessControlAttributeConfiguration{},
		"AWS::SSO::PermissionSet":                                     &sso.PermissionSet{},
		"AWS::SageMaker::App":                                         &sagemaker.App{},
		"AWS::SageMaker::AppImageConfig":                              &sagemaker.AppImageConfig{},
		"AWS::SageMaker::CodeRepository":                              &sagemaker.CodeRepository{},
		"AWS::SageMaker::DataQualityJobDefinition":                    &sagemaker.DataQualityJobDefinition{},
		"AWS::SageMaker::Device":                                      &sagemaker.Device{},
		"AWS::SageMaker::DeviceFleet":                                 &sagemaker.DeviceFleet{},
		"AWS::SageMaker::Domain":                                      &sagemaker.Domain{},
		"AWS::SageMaker::Endpoint":                                    &sagemaker.Endpoint{},
		"AWS::SageMaker::EndpointConfig":                              &sagemaker.EndpointConfig{},
		"AWS::SageMaker::FeatureGroup":                                &sagemaker.FeatureGroup{},
		"AWS::SageMaker::Image":                                       &sagemaker.Image{},
		"AWS::SageMaker::ImageVersion":                                &sagemaker.ImageVersion{},
		"AWS::SageMaker::Model":                                       &sagemaker.Model{},
		"AWS::SageMaker::ModelBiasJobDefinition":                      &sagemaker.ModelBiasJobDefinition{},
		"AWS::SageMaker::ModelExplainabilityJobDefinition":            &sagemaker.ModelExplainabilityJobDefinition{},
		"AWS::SageMaker::ModelPackageGroup":                           &sagemaker.ModelPackageGroup{},
		"AWS::SageMaker::ModelQualityJobDefinition":                   &sagemaker.ModelQualityJobDefinition{},
		"AWS::SageMaker::MonitoringSchedule":                          &sagemaker.MonitoringSchedule{},
		"AWS::SageMaker::NotebookInstance":                            &sagemaker.NotebookInstance{},
		"AWS::SageMaker::NotebookInstanceLifecycleConfig":             &sagemaker.NotebookInstanceLifecycleConfig{},
		"AWS::SageMaker::Pipeline":                                    &sagemaker.Pipeline{},
		"AWS::SageMaker::Project":                                     &sagemaker.Project{},
		"AWS::SageMaker::UserProfile":                                 &sagemaker.UserProfile{},
		"AWS::SageMaker::Workteam":                                    &sagemaker.Workteam{},
		"AWS::SecretsManager::ResourcePolicy":                         &secretsmanager.ResourcePolicy{},
		"AWS::SecretsManager::RotationSchedule":                       &secretsmanager.RotationSchedule{},
		"AWS::SecretsManager::Secret":                                 &secretsmanager.Secret{},
		"AWS::SecretsManager::SecretTargetAttachment":                 &secretsmanager.SecretTargetAttachment{},
		"AWS::SecurityHub::Hub":                                       &securityhub.Hub{},
		"AWS::Serverless::Api":                                        &serverless.Api{},
		"AWS::Serverless::Application":                                &serverless.Application{},
		"AWS::Serverless::Function":                                   &serverless.Function{},
		"AWS::Serverless::LayerVersion":                               &serverless.LayerVersion{},
		"AWS::Serverless::SimpleTable":                                &serverless.SimpleTable{},
		"AWS::Serverless::StateMachine":                               &serverless.StateMachine{},
		"AWS::ServiceCatalog::AcceptedPortfolioShare":                 &servicecatalog.AcceptedPortfolioShare{},
		"AWS::ServiceCatalog::CloudFormationProduct":                  &servicecatalog.CloudFormationProduct{},
		"AWS::ServiceCatalog::CloudFormationProvisionedProduct":       &servicecatalog.CloudFormationProvisionedProduct{},
		"AWS::ServiceCatalog::LaunchNotificationConstraint":           &servicecatalog.LaunchNotificationConstraint{},
		"AWS::ServiceCatalog::LaunchRoleConstraint":                   &servicecatalog.LaunchRoleConstraint{},
		"AWS::ServiceCatalog::LaunchTemplateConstraint":               &servicecatalog.LaunchTemplateConstraint{},
		"AWS::ServiceCatalog::Portfolio":                              &servicecatalog.Portfolio{},
		"AWS::ServiceCatalog::PortfolioPrincipalAssociation":          &servicecatalog.PortfolioPrincipalAssociation{},
		"AWS::ServiceCatalog::PortfolioProductAssociation":            &servicecatalog.PortfolioProductAssociation{},
		"AWS::ServiceCatalog::PortfolioShare":                         &servicecatalog.PortfolioShare{},
		"AWS::ServiceCatalog::ResourceUpdateConstraint":               &servicecatalog.ResourceUpdateConstraint{},
		"AWS::ServiceCatalog::ServiceAction":                          &servicecatalog.ServiceAction{},
		"AWS::ServiceCatalog::ServiceActionAssociation":               &servicecatalog.ServiceActionAssociation{},
		"AWS::ServiceCatalog::StackSetConstraint":                     &servicecatalog.StackSetConstraint{},
		"AWS::ServiceCatalog::TagOption":                              &servicecatalog.TagOption{},
		"AWS::ServiceCatalog::TagOptionAssociation":                   &servicecatalog.TagOptionAssociation{},
		"AWS::ServiceCatalogAppRegistry::Application":                 &servicecatalogappregistry.Application{},
		"AWS::ServiceCatalogAppRegistry::AttributeGroup":              &servicecatalogappregistry.AttributeGroup{},
		"AWS::ServiceCatalogAppRegistry::AttributeGroupAssociation":   &servicecatalogappregistry.AttributeGroupAssociation{},
		"AWS::ServiceCatalogAppRegistry::ResourceAssociation":         &servicecatalogappregistry.ResourceAssociation{},
		"AWS::ServiceDiscovery::HttpNamespace":                        &servicediscovery.HttpNamespace{},
		"AWS::ServiceDiscovery::Instance":                             &servicediscovery.Instance{},
		"AWS::ServiceDiscovery::PrivateDnsNamespace":                  &servicediscovery.PrivateDnsNamespace{},
		"AWS::ServiceDiscovery::PublicDnsNamespace":                   &servicediscovery.PublicDnsNamespace{},
		"AWS::ServiceDiscovery::Service":                              &servicediscovery.Service{},
		"AWS::Signer::ProfilePermission":                              &signer.ProfilePermission{},
		"AWS::Signer::SigningProfile":                                 &signer.SigningProfile{},
		"AWS::StepFunctions::Activity":                                &stepfunctions.Activity{},
		"AWS::StepFunctions::StateMachine":                            &stepfunctions.StateMachine{},
		"AWS::Synthetics::Canary":                                     &synthetics.Canary{},
		"AWS::Timestream::Database":                                   &timestream.Database{},
		"AWS::Timestream::ScheduledQuery":                             &timestream.ScheduledQuery{},
		"AWS::Timestream::Table":                                      &timestream.Table{},
		"AWS::Transfer::Server":                                       &transfer.Server{},
		"AWS::Transfer::User":                                         &transfer.User{},
		"AWS::Transfer::Workflow":                                     &transfer.Workflow{},
		"AWS::WAF::ByteMatchSet":                                      &waf.ByteMatchSet{},
		"AWS::WAF::IPSet":                                             &waf.IPSet{},
		"AWS::WAF::Rule":                                              &waf.Rule{},
		"AWS::WAF::SizeConstraintSet":                                 &waf.SizeConstraintSet{},
		"AWS::WAF::SqlInjectionMatchSet":                              &waf.SqlInjectionMatchSet{},
		"AWS::WAF::WebACL":                                            &waf.WebACL{},
		"AWS::WAF::XssMatchSet":                                       &waf.XssMatchSet{},
		"AWS::WAFRegional::ByteMatchSet":                              &wafregional.ByteMatchSet{},
		"AWS::WAFRegional::GeoMatchSet":                               &wafregional.GeoMatchSet{},
		"AWS::WAFRegional::IPSet":                                     &wafregional.IPSet{},
		"AWS::WAFRegional::RateBasedRule":                             &wafregional.RateBasedRule{},
		"AWS::WAFRegional::RegexPatternSet":                           &wafregional.RegexPatternSet{},
		"AWS::WAFRegional::Rule":                                      &wafregional.Rule{},
		"AWS::WAFRegional::SizeConstraintSet":                         &wafregional.SizeConstraintSet{},
		"AWS::WAFRegional::SqlInjectionMatchSet":                      &wafregional.SqlInjectionMatchSet{},
		"AWS::WAFRegional::WebACL":                                    &wafregional.WebACL{},
		"AWS::WAFRegional::WebACLAssociation":                         &wafregional.WebACLAssociation{},
		"AWS::WAFRegional::XssMatchSet":                               &wafregional.XssMatchSet{},
		"AWS::WAFv2::IPSet":                                           &wafv2.IPSet{},
		"AWS::WAFv2::LoggingConfiguration":                            &wafv2.LoggingConfiguration{},
		"AWS::WAFv2::RegexPatternSet":                                 &wafv2.RegexPatternSet{},
		"AWS::WAFv2::RuleGroup":                                       &wafv2.RuleGroup{},
		"AWS::WAFv2::WebACL":                                          &wafv2.WebACL{},
		"AWS::WAFv2::WebACLAssociation":                               &wafv2.WebACLAssociation{},
		"AWS::Wisdom::Assistant":                                      &wisdom.Assistant{},
		"AWS::Wisdom::AssistantAssociation":                           &wisdom.AssistantAssociation{},
		"AWS::Wisdom::KnowledgeBase":                                  &wisdom.KnowledgeBase{},
		"AWS::WorkSpaces::ConnectionAlias":                            &workspaces.ConnectionAlias{},
		"AWS::WorkSpaces::Workspace":                                  &workspaces.Workspace{},
		"AWS::XRay::Group":                                            &xray.Group{},
		"AWS::XRay::SamplingRule":                                     &xray.SamplingRule{},
		"Alexa::ASK::Skill":                                           &ask.Skill{},
	}
}

// GetAllACMPCACertificateResources retrieves all acmpca.Certificate items from an AWS CloudFormation template
func (t *Template) GetAllACMPCACertificateResources() map[string]*acmpca.Certificate {
	results := map[string]*acmpca.Certificate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *acmpca.Certificate:
			results[name] = resource
		}
	}
	return results
}

// GetACMPCACertificateWithName retrieves all acmpca.Certificate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetACMPCACertificateWithName(name string) (*acmpca.Certificate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *acmpca.Certificate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type acmpca.Certificate not found", name)
}

// GetAllACMPCACertificateAuthorityResources retrieves all acmpca.CertificateAuthority items from an AWS CloudFormation template
func (t *Template) GetAllACMPCACertificateAuthorityResources() map[string]*acmpca.CertificateAuthority {
	results := map[string]*acmpca.CertificateAuthority{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *acmpca.CertificateAuthority:
			results[name] = resource
		}
	}
	return results
}

// GetACMPCACertificateAuthorityWithName retrieves all acmpca.CertificateAuthority items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetACMPCACertificateAuthorityWithName(name string) (*acmpca.CertificateAuthority, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *acmpca.CertificateAuthority:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type acmpca.CertificateAuthority not found", name)
}

// GetAllACMPCACertificateAuthorityActivationResources retrieves all acmpca.CertificateAuthorityActivation items from an AWS CloudFormation template
func (t *Template) GetAllACMPCACertificateAuthorityActivationResources() map[string]*acmpca.CertificateAuthorityActivation {
	results := map[string]*acmpca.CertificateAuthorityActivation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *acmpca.CertificateAuthorityActivation:
			results[name] = resource
		}
	}
	return results
}

// GetACMPCACertificateAuthorityActivationWithName retrieves all acmpca.CertificateAuthorityActivation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetACMPCACertificateAuthorityActivationWithName(name string) (*acmpca.CertificateAuthorityActivation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *acmpca.CertificateAuthorityActivation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type acmpca.CertificateAuthorityActivation not found", name)
}

// GetAllACMPCAPermissionResources retrieves all acmpca.Permission items from an AWS CloudFormation template
func (t *Template) GetAllACMPCAPermissionResources() map[string]*acmpca.Permission {
	results := map[string]*acmpca.Permission{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *acmpca.Permission:
			results[name] = resource
		}
	}
	return results
}

// GetACMPCAPermissionWithName retrieves all acmpca.Permission items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetACMPCAPermissionWithName(name string) (*acmpca.Permission, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *acmpca.Permission:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type acmpca.Permission not found", name)
}

// GetAllAPSRuleGroupsNamespaceResources retrieves all aps.RuleGroupsNamespace items from an AWS CloudFormation template
func (t *Template) GetAllAPSRuleGroupsNamespaceResources() map[string]*aps.RuleGroupsNamespace {
	results := map[string]*aps.RuleGroupsNamespace{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *aps.RuleGroupsNamespace:
			results[name] = resource
		}
	}
	return results
}

// GetAPSRuleGroupsNamespaceWithName retrieves all aps.RuleGroupsNamespace items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAPSRuleGroupsNamespaceWithName(name string) (*aps.RuleGroupsNamespace, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *aps.RuleGroupsNamespace:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type aps.RuleGroupsNamespace not found", name)
}

// GetAllAPSWorkspaceResources retrieves all aps.Workspace items from an AWS CloudFormation template
func (t *Template) GetAllAPSWorkspaceResources() map[string]*aps.Workspace {
	results := map[string]*aps.Workspace{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *aps.Workspace:
			results[name] = resource
		}
	}
	return results
}

// GetAPSWorkspaceWithName retrieves all aps.Workspace items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAPSWorkspaceWithName(name string) (*aps.Workspace, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *aps.Workspace:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type aps.Workspace not found", name)
}

// GetAllAccessAnalyzerAnalyzerResources retrieves all accessanalyzer.Analyzer items from an AWS CloudFormation template
func (t *Template) GetAllAccessAnalyzerAnalyzerResources() map[string]*accessanalyzer.Analyzer {
	results := map[string]*accessanalyzer.Analyzer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *accessanalyzer.Analyzer:
			results[name] = resource
		}
	}
	return results
}

// GetAccessAnalyzerAnalyzerWithName retrieves all accessanalyzer.Analyzer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAccessAnalyzerAnalyzerWithName(name string) (*accessanalyzer.Analyzer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *accessanalyzer.Analyzer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type accessanalyzer.Analyzer not found", name)
}

// GetAllAmazonMQBrokerResources retrieves all amazonmq.Broker items from an AWS CloudFormation template
func (t *Template) GetAllAmazonMQBrokerResources() map[string]*amazonmq.Broker {
	results := map[string]*amazonmq.Broker{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amazonmq.Broker:
			results[name] = resource
		}
	}
	return results
}

// GetAmazonMQBrokerWithName retrieves all amazonmq.Broker items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmazonMQBrokerWithName(name string) (*amazonmq.Broker, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amazonmq.Broker:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amazonmq.Broker not found", name)
}

// GetAllAmazonMQConfigurationResources retrieves all amazonmq.Configuration items from an AWS CloudFormation template
func (t *Template) GetAllAmazonMQConfigurationResources() map[string]*amazonmq.Configuration {
	results := map[string]*amazonmq.Configuration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amazonmq.Configuration:
			results[name] = resource
		}
	}
	return results
}

// GetAmazonMQConfigurationWithName retrieves all amazonmq.Configuration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmazonMQConfigurationWithName(name string) (*amazonmq.Configuration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amazonmq.Configuration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amazonmq.Configuration not found", name)
}

// GetAllAmazonMQConfigurationAssociationResources retrieves all amazonmq.ConfigurationAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAmazonMQConfigurationAssociationResources() map[string]*amazonmq.ConfigurationAssociation {
	results := map[string]*amazonmq.ConfigurationAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amazonmq.ConfigurationAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetAmazonMQConfigurationAssociationWithName retrieves all amazonmq.ConfigurationAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmazonMQConfigurationAssociationWithName(name string) (*amazonmq.ConfigurationAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amazonmq.ConfigurationAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amazonmq.ConfigurationAssociation not found", name)
}

// GetAllAmplifyAppResources retrieves all amplify.App items from an AWS CloudFormation template
func (t *Template) GetAllAmplifyAppResources() map[string]*amplify.App {
	results := map[string]*amplify.App{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amplify.App:
			results[name] = resource
		}
	}
	return results
}

// GetAmplifyAppWithName retrieves all amplify.App items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmplifyAppWithName(name string) (*amplify.App, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amplify.App:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amplify.App not found", name)
}

// GetAllAmplifyBranchResources retrieves all amplify.Branch items from an AWS CloudFormation template
func (t *Template) GetAllAmplifyBranchResources() map[string]*amplify.Branch {
	results := map[string]*amplify.Branch{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amplify.Branch:
			results[name] = resource
		}
	}
	return results
}

// GetAmplifyBranchWithName retrieves all amplify.Branch items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmplifyBranchWithName(name string) (*amplify.Branch, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amplify.Branch:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amplify.Branch not found", name)
}

// GetAllAmplifyDomainResources retrieves all amplify.Domain items from an AWS CloudFormation template
func (t *Template) GetAllAmplifyDomainResources() map[string]*amplify.Domain {
	results := map[string]*amplify.Domain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amplify.Domain:
			results[name] = resource
		}
	}
	return results
}

// GetAmplifyDomainWithName retrieves all amplify.Domain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmplifyDomainWithName(name string) (*amplify.Domain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amplify.Domain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amplify.Domain not found", name)
}

// GetAllAmplifyUIBuilderComponentResources retrieves all amplifyuibuilder.Component items from an AWS CloudFormation template
func (t *Template) GetAllAmplifyUIBuilderComponentResources() map[string]*amplifyuibuilder.Component {
	results := map[string]*amplifyuibuilder.Component{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amplifyuibuilder.Component:
			results[name] = resource
		}
	}
	return results
}

// GetAmplifyUIBuilderComponentWithName retrieves all amplifyuibuilder.Component items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmplifyUIBuilderComponentWithName(name string) (*amplifyuibuilder.Component, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amplifyuibuilder.Component:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amplifyuibuilder.Component not found", name)
}

// GetAllAmplifyUIBuilderThemeResources retrieves all amplifyuibuilder.Theme items from an AWS CloudFormation template
func (t *Template) GetAllAmplifyUIBuilderThemeResources() map[string]*amplifyuibuilder.Theme {
	results := map[string]*amplifyuibuilder.Theme{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *amplifyuibuilder.Theme:
			results[name] = resource
		}
	}
	return results
}

// GetAmplifyUIBuilderThemeWithName retrieves all amplifyuibuilder.Theme items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAmplifyUIBuilderThemeWithName(name string) (*amplifyuibuilder.Theme, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *amplifyuibuilder.Theme:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type amplifyuibuilder.Theme not found", name)
}

// GetAllApiGatewayAccountResources retrieves all apigateway.Account items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayAccountResources() map[string]*apigateway.Account {
	results := map[string]*apigateway.Account{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.Account:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayAccountWithName retrieves all apigateway.Account items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayAccountWithName(name string) (*apigateway.Account, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.Account:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.Account not found", name)
}

// GetAllApiGatewayApiKeyResources retrieves all apigateway.ApiKey items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayApiKeyResources() map[string]*apigateway.ApiKey {
	results := map[string]*apigateway.ApiKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.ApiKey:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayApiKeyWithName retrieves all apigateway.ApiKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayApiKeyWithName(name string) (*apigateway.ApiKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.ApiKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.ApiKey not found", name)
}

// GetAllApiGatewayAuthorizerResources retrieves all apigateway.Authorizer items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayAuthorizerResources() map[string]*apigateway.Authorizer {
	results := map[string]*apigateway.Authorizer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.Authorizer:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayAuthorizerWithName retrieves all apigateway.Authorizer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayAuthorizerWithName(name string) (*apigateway.Authorizer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.Authorizer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.Authorizer not found", name)
}

// GetAllApiGatewayBasePathMappingResources retrieves all apigateway.BasePathMapping items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayBasePathMappingResources() map[string]*apigateway.BasePathMapping {
	results := map[string]*apigateway.BasePathMapping{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.BasePathMapping:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayBasePathMappingWithName retrieves all apigateway.BasePathMapping items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayBasePathMappingWithName(name string) (*apigateway.BasePathMapping, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.BasePathMapping:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.BasePathMapping not found", name)
}

// GetAllApiGatewayClientCertificateResources retrieves all apigateway.ClientCertificate items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayClientCertificateResources() map[string]*apigateway.ClientCertificate {
	results := map[string]*apigateway.ClientCertificate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.ClientCertificate:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayClientCertificateWithName retrieves all apigateway.ClientCertificate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayClientCertificateWithName(name string) (*apigateway.ClientCertificate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.ClientCertificate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.ClientCertificate not found", name)
}

// GetAllApiGatewayDeploymentResources retrieves all apigateway.Deployment items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayDeploymentResources() map[string]*apigateway.Deployment {
	results := map[string]*apigateway.Deployment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.Deployment:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayDeploymentWithName retrieves all apigateway.Deployment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayDeploymentWithName(name string) (*apigateway.Deployment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.Deployment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.Deployment not found", name)
}

// GetAllApiGatewayDocumentationPartResources retrieves all apigateway.DocumentationPart items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayDocumentationPartResources() map[string]*apigateway.DocumentationPart {
	results := map[string]*apigateway.DocumentationPart{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.DocumentationPart:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayDocumentationPartWithName retrieves all apigateway.DocumentationPart items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayDocumentationPartWithName(name string) (*apigateway.DocumentationPart, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.DocumentationPart:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.DocumentationPart not found", name)
}

// GetAllApiGatewayDocumentationVersionResources retrieves all apigateway.DocumentationVersion items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayDocumentationVersionResources() map[string]*apigateway.DocumentationVersion {
	results := map[string]*apigateway.DocumentationVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.DocumentationVersion:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayDocumentationVersionWithName retrieves all apigateway.DocumentationVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayDocumentationVersionWithName(name string) (*apigateway.DocumentationVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.DocumentationVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.DocumentationVersion not found", name)
}

// GetAllApiGatewayDomainNameResources retrieves all apigateway.DomainName items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayDomainNameResources() map[string]*apigateway.DomainName {
	results := map[string]*apigateway.DomainName{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.DomainName:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayDomainNameWithName retrieves all apigateway.DomainName items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayDomainNameWithName(name string) (*apigateway.DomainName, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.DomainName:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.DomainName not found", name)
}

// GetAllApiGatewayGatewayResponseResources retrieves all apigateway.GatewayResponse items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayGatewayResponseResources() map[string]*apigateway.GatewayResponse {
	results := map[string]*apigateway.GatewayResponse{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.GatewayResponse:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayGatewayResponseWithName retrieves all apigateway.GatewayResponse items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayGatewayResponseWithName(name string) (*apigateway.GatewayResponse, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.GatewayResponse:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.GatewayResponse not found", name)
}

// GetAllApiGatewayMethodResources retrieves all apigateway.Method items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayMethodResources() map[string]*apigateway.Method {
	results := map[string]*apigateway.Method{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.Method:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayMethodWithName retrieves all apigateway.Method items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayMethodWithName(name string) (*apigateway.Method, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.Method:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.Method not found", name)
}

// GetAllApiGatewayModelResources retrieves all apigateway.Model items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayModelResources() map[string]*apigateway.Model {
	results := map[string]*apigateway.Model{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.Model:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayModelWithName retrieves all apigateway.Model items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayModelWithName(name string) (*apigateway.Model, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.Model:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.Model not found", name)
}

// GetAllApiGatewayRequestValidatorResources retrieves all apigateway.RequestValidator items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayRequestValidatorResources() map[string]*apigateway.RequestValidator {
	results := map[string]*apigateway.RequestValidator{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.RequestValidator:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayRequestValidatorWithName retrieves all apigateway.RequestValidator items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayRequestValidatorWithName(name string) (*apigateway.RequestValidator, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.RequestValidator:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.RequestValidator not found", name)
}

// GetAllApiGatewayResourceResources retrieves all apigateway.Resource items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayResourceResources() map[string]*apigateway.Resource {
	results := map[string]*apigateway.Resource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.Resource:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayResourceWithName retrieves all apigateway.Resource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayResourceWithName(name string) (*apigateway.Resource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.Resource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.Resource not found", name)
}

// GetAllApiGatewayRestApiResources retrieves all apigateway.RestApi items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayRestApiResources() map[string]*apigateway.RestApi {
	results := map[string]*apigateway.RestApi{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.RestApi:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayRestApiWithName retrieves all apigateway.RestApi items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayRestApiWithName(name string) (*apigateway.RestApi, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.RestApi:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.RestApi not found", name)
}

// GetAllApiGatewayStageResources retrieves all apigateway.Stage items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayStageResources() map[string]*apigateway.Stage {
	results := map[string]*apigateway.Stage{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.Stage:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayStageWithName retrieves all apigateway.Stage items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayStageWithName(name string) (*apigateway.Stage, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.Stage:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.Stage not found", name)
}

// GetAllApiGatewayUsagePlanResources retrieves all apigateway.UsagePlan items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayUsagePlanResources() map[string]*apigateway.UsagePlan {
	results := map[string]*apigateway.UsagePlan{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.UsagePlan:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayUsagePlanWithName retrieves all apigateway.UsagePlan items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayUsagePlanWithName(name string) (*apigateway.UsagePlan, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.UsagePlan:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.UsagePlan not found", name)
}

// GetAllApiGatewayUsagePlanKeyResources retrieves all apigateway.UsagePlanKey items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayUsagePlanKeyResources() map[string]*apigateway.UsagePlanKey {
	results := map[string]*apigateway.UsagePlanKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.UsagePlanKey:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayUsagePlanKeyWithName retrieves all apigateway.UsagePlanKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayUsagePlanKeyWithName(name string) (*apigateway.UsagePlanKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.UsagePlanKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.UsagePlanKey not found", name)
}

// GetAllApiGatewayVpcLinkResources retrieves all apigateway.VpcLink items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayVpcLinkResources() map[string]*apigateway.VpcLink {
	results := map[string]*apigateway.VpcLink{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigateway.VpcLink:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayVpcLinkWithName retrieves all apigateway.VpcLink items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayVpcLinkWithName(name string) (*apigateway.VpcLink, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigateway.VpcLink:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigateway.VpcLink not found", name)
}

// GetAllApiGatewayV2ApiResources retrieves all apigatewayv2.Api items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2ApiResources() map[string]*apigatewayv2.Api {
	results := map[string]*apigatewayv2.Api{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.Api:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2ApiWithName retrieves all apigatewayv2.Api items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2ApiWithName(name string) (*apigatewayv2.Api, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.Api:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.Api not found", name)
}

// GetAllApiGatewayV2ApiGatewayManagedOverridesResources retrieves all apigatewayv2.ApiGatewayManagedOverrides items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2ApiGatewayManagedOverridesResources() map[string]*apigatewayv2.ApiGatewayManagedOverrides {
	results := map[string]*apigatewayv2.ApiGatewayManagedOverrides{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.ApiGatewayManagedOverrides:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2ApiGatewayManagedOverridesWithName retrieves all apigatewayv2.ApiGatewayManagedOverrides items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2ApiGatewayManagedOverridesWithName(name string) (*apigatewayv2.ApiGatewayManagedOverrides, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.ApiGatewayManagedOverrides:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.ApiGatewayManagedOverrides not found", name)
}

// GetAllApiGatewayV2ApiMappingResources retrieves all apigatewayv2.ApiMapping items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2ApiMappingResources() map[string]*apigatewayv2.ApiMapping {
	results := map[string]*apigatewayv2.ApiMapping{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.ApiMapping:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2ApiMappingWithName retrieves all apigatewayv2.ApiMapping items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2ApiMappingWithName(name string) (*apigatewayv2.ApiMapping, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.ApiMapping:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.ApiMapping not found", name)
}

// GetAllApiGatewayV2AuthorizerResources retrieves all apigatewayv2.Authorizer items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2AuthorizerResources() map[string]*apigatewayv2.Authorizer {
	results := map[string]*apigatewayv2.Authorizer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.Authorizer:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2AuthorizerWithName retrieves all apigatewayv2.Authorizer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2AuthorizerWithName(name string) (*apigatewayv2.Authorizer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.Authorizer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.Authorizer not found", name)
}

// GetAllApiGatewayV2DeploymentResources retrieves all apigatewayv2.Deployment items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2DeploymentResources() map[string]*apigatewayv2.Deployment {
	results := map[string]*apigatewayv2.Deployment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.Deployment:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2DeploymentWithName retrieves all apigatewayv2.Deployment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2DeploymentWithName(name string) (*apigatewayv2.Deployment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.Deployment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.Deployment not found", name)
}

// GetAllApiGatewayV2DomainNameResources retrieves all apigatewayv2.DomainName items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2DomainNameResources() map[string]*apigatewayv2.DomainName {
	results := map[string]*apigatewayv2.DomainName{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.DomainName:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2DomainNameWithName retrieves all apigatewayv2.DomainName items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2DomainNameWithName(name string) (*apigatewayv2.DomainName, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.DomainName:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.DomainName not found", name)
}

// GetAllApiGatewayV2IntegrationResources retrieves all apigatewayv2.Integration items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2IntegrationResources() map[string]*apigatewayv2.Integration {
	results := map[string]*apigatewayv2.Integration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.Integration:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2IntegrationWithName retrieves all apigatewayv2.Integration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2IntegrationWithName(name string) (*apigatewayv2.Integration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.Integration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.Integration not found", name)
}

// GetAllApiGatewayV2IntegrationResponseResources retrieves all apigatewayv2.IntegrationResponse items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2IntegrationResponseResources() map[string]*apigatewayv2.IntegrationResponse {
	results := map[string]*apigatewayv2.IntegrationResponse{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.IntegrationResponse:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2IntegrationResponseWithName retrieves all apigatewayv2.IntegrationResponse items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2IntegrationResponseWithName(name string) (*apigatewayv2.IntegrationResponse, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.IntegrationResponse:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.IntegrationResponse not found", name)
}

// GetAllApiGatewayV2ModelResources retrieves all apigatewayv2.Model items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2ModelResources() map[string]*apigatewayv2.Model {
	results := map[string]*apigatewayv2.Model{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.Model:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2ModelWithName retrieves all apigatewayv2.Model items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2ModelWithName(name string) (*apigatewayv2.Model, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.Model:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.Model not found", name)
}

// GetAllApiGatewayV2RouteResources retrieves all apigatewayv2.Route items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2RouteResources() map[string]*apigatewayv2.Route {
	results := map[string]*apigatewayv2.Route{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.Route:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2RouteWithName retrieves all apigatewayv2.Route items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2RouteWithName(name string) (*apigatewayv2.Route, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.Route:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.Route not found", name)
}

// GetAllApiGatewayV2RouteResponseResources retrieves all apigatewayv2.RouteResponse items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2RouteResponseResources() map[string]*apigatewayv2.RouteResponse {
	results := map[string]*apigatewayv2.RouteResponse{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.RouteResponse:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2RouteResponseWithName retrieves all apigatewayv2.RouteResponse items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2RouteResponseWithName(name string) (*apigatewayv2.RouteResponse, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.RouteResponse:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.RouteResponse not found", name)
}

// GetAllApiGatewayV2StageResources retrieves all apigatewayv2.Stage items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2StageResources() map[string]*apigatewayv2.Stage {
	results := map[string]*apigatewayv2.Stage{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.Stage:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2StageWithName retrieves all apigatewayv2.Stage items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2StageWithName(name string) (*apigatewayv2.Stage, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.Stage:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.Stage not found", name)
}

// GetAllApiGatewayV2VpcLinkResources retrieves all apigatewayv2.VpcLink items from an AWS CloudFormation template
func (t *Template) GetAllApiGatewayV2VpcLinkResources() map[string]*apigatewayv2.VpcLink {
	results := map[string]*apigatewayv2.VpcLink{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apigatewayv2.VpcLink:
			results[name] = resource
		}
	}
	return results
}

// GetApiGatewayV2VpcLinkWithName retrieves all apigatewayv2.VpcLink items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApiGatewayV2VpcLinkWithName(name string) (*apigatewayv2.VpcLink, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apigatewayv2.VpcLink:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apigatewayv2.VpcLink not found", name)
}

// GetAllAppConfigApplicationResources retrieves all appconfig.Application items from an AWS CloudFormation template
func (t *Template) GetAllAppConfigApplicationResources() map[string]*appconfig.Application {
	results := map[string]*appconfig.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appconfig.Application:
			results[name] = resource
		}
	}
	return results
}

// GetAppConfigApplicationWithName retrieves all appconfig.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppConfigApplicationWithName(name string) (*appconfig.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appconfig.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appconfig.Application not found", name)
}

// GetAllAppConfigConfigurationProfileResources retrieves all appconfig.ConfigurationProfile items from an AWS CloudFormation template
func (t *Template) GetAllAppConfigConfigurationProfileResources() map[string]*appconfig.ConfigurationProfile {
	results := map[string]*appconfig.ConfigurationProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appconfig.ConfigurationProfile:
			results[name] = resource
		}
	}
	return results
}

// GetAppConfigConfigurationProfileWithName retrieves all appconfig.ConfigurationProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppConfigConfigurationProfileWithName(name string) (*appconfig.ConfigurationProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appconfig.ConfigurationProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appconfig.ConfigurationProfile not found", name)
}

// GetAllAppConfigDeploymentResources retrieves all appconfig.Deployment items from an AWS CloudFormation template
func (t *Template) GetAllAppConfigDeploymentResources() map[string]*appconfig.Deployment {
	results := map[string]*appconfig.Deployment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appconfig.Deployment:
			results[name] = resource
		}
	}
	return results
}

// GetAppConfigDeploymentWithName retrieves all appconfig.Deployment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppConfigDeploymentWithName(name string) (*appconfig.Deployment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appconfig.Deployment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appconfig.Deployment not found", name)
}

// GetAllAppConfigDeploymentStrategyResources retrieves all appconfig.DeploymentStrategy items from an AWS CloudFormation template
func (t *Template) GetAllAppConfigDeploymentStrategyResources() map[string]*appconfig.DeploymentStrategy {
	results := map[string]*appconfig.DeploymentStrategy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appconfig.DeploymentStrategy:
			results[name] = resource
		}
	}
	return results
}

// GetAppConfigDeploymentStrategyWithName retrieves all appconfig.DeploymentStrategy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppConfigDeploymentStrategyWithName(name string) (*appconfig.DeploymentStrategy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appconfig.DeploymentStrategy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appconfig.DeploymentStrategy not found", name)
}

// GetAllAppConfigEnvironmentResources retrieves all appconfig.Environment items from an AWS CloudFormation template
func (t *Template) GetAllAppConfigEnvironmentResources() map[string]*appconfig.Environment {
	results := map[string]*appconfig.Environment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appconfig.Environment:
			results[name] = resource
		}
	}
	return results
}

// GetAppConfigEnvironmentWithName retrieves all appconfig.Environment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppConfigEnvironmentWithName(name string) (*appconfig.Environment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appconfig.Environment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appconfig.Environment not found", name)
}

// GetAllAppConfigHostedConfigurationVersionResources retrieves all appconfig.HostedConfigurationVersion items from an AWS CloudFormation template
func (t *Template) GetAllAppConfigHostedConfigurationVersionResources() map[string]*appconfig.HostedConfigurationVersion {
	results := map[string]*appconfig.HostedConfigurationVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appconfig.HostedConfigurationVersion:
			results[name] = resource
		}
	}
	return results
}

// GetAppConfigHostedConfigurationVersionWithName retrieves all appconfig.HostedConfigurationVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppConfigHostedConfigurationVersionWithName(name string) (*appconfig.HostedConfigurationVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appconfig.HostedConfigurationVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appconfig.HostedConfigurationVersion not found", name)
}

// GetAllAppFlowConnectorProfileResources retrieves all appflow.ConnectorProfile items from an AWS CloudFormation template
func (t *Template) GetAllAppFlowConnectorProfileResources() map[string]*appflow.ConnectorProfile {
	results := map[string]*appflow.ConnectorProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appflow.ConnectorProfile:
			results[name] = resource
		}
	}
	return results
}

// GetAppFlowConnectorProfileWithName retrieves all appflow.ConnectorProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppFlowConnectorProfileWithName(name string) (*appflow.ConnectorProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appflow.ConnectorProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appflow.ConnectorProfile not found", name)
}

// GetAllAppFlowFlowResources retrieves all appflow.Flow items from an AWS CloudFormation template
func (t *Template) GetAllAppFlowFlowResources() map[string]*appflow.Flow {
	results := map[string]*appflow.Flow{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appflow.Flow:
			results[name] = resource
		}
	}
	return results
}

// GetAppFlowFlowWithName retrieves all appflow.Flow items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppFlowFlowWithName(name string) (*appflow.Flow, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appflow.Flow:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appflow.Flow not found", name)
}

// GetAllAppIntegrationsEventIntegrationResources retrieves all appintegrations.EventIntegration items from an AWS CloudFormation template
func (t *Template) GetAllAppIntegrationsEventIntegrationResources() map[string]*appintegrations.EventIntegration {
	results := map[string]*appintegrations.EventIntegration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appintegrations.EventIntegration:
			results[name] = resource
		}
	}
	return results
}

// GetAppIntegrationsEventIntegrationWithName retrieves all appintegrations.EventIntegration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppIntegrationsEventIntegrationWithName(name string) (*appintegrations.EventIntegration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appintegrations.EventIntegration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appintegrations.EventIntegration not found", name)
}

// GetAllAppMeshGatewayRouteResources retrieves all appmesh.GatewayRoute items from an AWS CloudFormation template
func (t *Template) GetAllAppMeshGatewayRouteResources() map[string]*appmesh.GatewayRoute {
	results := map[string]*appmesh.GatewayRoute{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appmesh.GatewayRoute:
			results[name] = resource
		}
	}
	return results
}

// GetAppMeshGatewayRouteWithName retrieves all appmesh.GatewayRoute items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppMeshGatewayRouteWithName(name string) (*appmesh.GatewayRoute, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appmesh.GatewayRoute:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appmesh.GatewayRoute not found", name)
}

// GetAllAppMeshMeshResources retrieves all appmesh.Mesh items from an AWS CloudFormation template
func (t *Template) GetAllAppMeshMeshResources() map[string]*appmesh.Mesh {
	results := map[string]*appmesh.Mesh{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appmesh.Mesh:
			results[name] = resource
		}
	}
	return results
}

// GetAppMeshMeshWithName retrieves all appmesh.Mesh items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppMeshMeshWithName(name string) (*appmesh.Mesh, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appmesh.Mesh:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appmesh.Mesh not found", name)
}

// GetAllAppMeshRouteResources retrieves all appmesh.Route items from an AWS CloudFormation template
func (t *Template) GetAllAppMeshRouteResources() map[string]*appmesh.Route {
	results := map[string]*appmesh.Route{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appmesh.Route:
			results[name] = resource
		}
	}
	return results
}

// GetAppMeshRouteWithName retrieves all appmesh.Route items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppMeshRouteWithName(name string) (*appmesh.Route, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appmesh.Route:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appmesh.Route not found", name)
}

// GetAllAppMeshVirtualGatewayResources retrieves all appmesh.VirtualGateway items from an AWS CloudFormation template
func (t *Template) GetAllAppMeshVirtualGatewayResources() map[string]*appmesh.VirtualGateway {
	results := map[string]*appmesh.VirtualGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appmesh.VirtualGateway:
			results[name] = resource
		}
	}
	return results
}

// GetAppMeshVirtualGatewayWithName retrieves all appmesh.VirtualGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppMeshVirtualGatewayWithName(name string) (*appmesh.VirtualGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appmesh.VirtualGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appmesh.VirtualGateway not found", name)
}

// GetAllAppMeshVirtualNodeResources retrieves all appmesh.VirtualNode items from an AWS CloudFormation template
func (t *Template) GetAllAppMeshVirtualNodeResources() map[string]*appmesh.VirtualNode {
	results := map[string]*appmesh.VirtualNode{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appmesh.VirtualNode:
			results[name] = resource
		}
	}
	return results
}

// GetAppMeshVirtualNodeWithName retrieves all appmesh.VirtualNode items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppMeshVirtualNodeWithName(name string) (*appmesh.VirtualNode, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appmesh.VirtualNode:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appmesh.VirtualNode not found", name)
}

// GetAllAppMeshVirtualRouterResources retrieves all appmesh.VirtualRouter items from an AWS CloudFormation template
func (t *Template) GetAllAppMeshVirtualRouterResources() map[string]*appmesh.VirtualRouter {
	results := map[string]*appmesh.VirtualRouter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appmesh.VirtualRouter:
			results[name] = resource
		}
	}
	return results
}

// GetAppMeshVirtualRouterWithName retrieves all appmesh.VirtualRouter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppMeshVirtualRouterWithName(name string) (*appmesh.VirtualRouter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appmesh.VirtualRouter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appmesh.VirtualRouter not found", name)
}

// GetAllAppMeshVirtualServiceResources retrieves all appmesh.VirtualService items from an AWS CloudFormation template
func (t *Template) GetAllAppMeshVirtualServiceResources() map[string]*appmesh.VirtualService {
	results := map[string]*appmesh.VirtualService{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appmesh.VirtualService:
			results[name] = resource
		}
	}
	return results
}

// GetAppMeshVirtualServiceWithName retrieves all appmesh.VirtualService items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppMeshVirtualServiceWithName(name string) (*appmesh.VirtualService, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appmesh.VirtualService:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appmesh.VirtualService not found", name)
}

// GetAllAppRunnerServiceResources retrieves all apprunner.Service items from an AWS CloudFormation template
func (t *Template) GetAllAppRunnerServiceResources() map[string]*apprunner.Service {
	results := map[string]*apprunner.Service{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *apprunner.Service:
			results[name] = resource
		}
	}
	return results
}

// GetAppRunnerServiceWithName retrieves all apprunner.Service items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppRunnerServiceWithName(name string) (*apprunner.Service, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *apprunner.Service:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type apprunner.Service not found", name)
}

// GetAllAppStreamAppBlockResources retrieves all appstream.AppBlock items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamAppBlockResources() map[string]*appstream.AppBlock {
	results := map[string]*appstream.AppBlock{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.AppBlock:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamAppBlockWithName retrieves all appstream.AppBlock items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamAppBlockWithName(name string) (*appstream.AppBlock, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.AppBlock:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.AppBlock not found", name)
}

// GetAllAppStreamApplicationResources retrieves all appstream.Application items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamApplicationResources() map[string]*appstream.Application {
	results := map[string]*appstream.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.Application:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamApplicationWithName retrieves all appstream.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamApplicationWithName(name string) (*appstream.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.Application not found", name)
}

// GetAllAppStreamApplicationFleetAssociationResources retrieves all appstream.ApplicationFleetAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamApplicationFleetAssociationResources() map[string]*appstream.ApplicationFleetAssociation {
	results := map[string]*appstream.ApplicationFleetAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.ApplicationFleetAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamApplicationFleetAssociationWithName retrieves all appstream.ApplicationFleetAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamApplicationFleetAssociationWithName(name string) (*appstream.ApplicationFleetAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.ApplicationFleetAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.ApplicationFleetAssociation not found", name)
}

// GetAllAppStreamDirectoryConfigResources retrieves all appstream.DirectoryConfig items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamDirectoryConfigResources() map[string]*appstream.DirectoryConfig {
	results := map[string]*appstream.DirectoryConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.DirectoryConfig:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamDirectoryConfigWithName retrieves all appstream.DirectoryConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamDirectoryConfigWithName(name string) (*appstream.DirectoryConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.DirectoryConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.DirectoryConfig not found", name)
}

// GetAllAppStreamFleetResources retrieves all appstream.Fleet items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamFleetResources() map[string]*appstream.Fleet {
	results := map[string]*appstream.Fleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.Fleet:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamFleetWithName retrieves all appstream.Fleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamFleetWithName(name string) (*appstream.Fleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.Fleet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.Fleet not found", name)
}

// GetAllAppStreamImageBuilderResources retrieves all appstream.ImageBuilder items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamImageBuilderResources() map[string]*appstream.ImageBuilder {
	results := map[string]*appstream.ImageBuilder{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.ImageBuilder:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamImageBuilderWithName retrieves all appstream.ImageBuilder items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamImageBuilderWithName(name string) (*appstream.ImageBuilder, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.ImageBuilder:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.ImageBuilder not found", name)
}

// GetAllAppStreamStackResources retrieves all appstream.Stack items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamStackResources() map[string]*appstream.Stack {
	results := map[string]*appstream.Stack{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.Stack:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamStackWithName retrieves all appstream.Stack items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamStackWithName(name string) (*appstream.Stack, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.Stack:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.Stack not found", name)
}

// GetAllAppStreamStackFleetAssociationResources retrieves all appstream.StackFleetAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamStackFleetAssociationResources() map[string]*appstream.StackFleetAssociation {
	results := map[string]*appstream.StackFleetAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.StackFleetAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamStackFleetAssociationWithName retrieves all appstream.StackFleetAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamStackFleetAssociationWithName(name string) (*appstream.StackFleetAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.StackFleetAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.StackFleetAssociation not found", name)
}

// GetAllAppStreamStackUserAssociationResources retrieves all appstream.StackUserAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamStackUserAssociationResources() map[string]*appstream.StackUserAssociation {
	results := map[string]*appstream.StackUserAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.StackUserAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamStackUserAssociationWithName retrieves all appstream.StackUserAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamStackUserAssociationWithName(name string) (*appstream.StackUserAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.StackUserAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.StackUserAssociation not found", name)
}

// GetAllAppStreamUserResources retrieves all appstream.User items from an AWS CloudFormation template
func (t *Template) GetAllAppStreamUserResources() map[string]*appstream.User {
	results := map[string]*appstream.User{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appstream.User:
			results[name] = resource
		}
	}
	return results
}

// GetAppStreamUserWithName retrieves all appstream.User items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppStreamUserWithName(name string) (*appstream.User, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appstream.User:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appstream.User not found", name)
}

// GetAllAppSyncApiCacheResources retrieves all appsync.ApiCache items from an AWS CloudFormation template
func (t *Template) GetAllAppSyncApiCacheResources() map[string]*appsync.ApiCache {
	results := map[string]*appsync.ApiCache{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appsync.ApiCache:
			results[name] = resource
		}
	}
	return results
}

// GetAppSyncApiCacheWithName retrieves all appsync.ApiCache items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppSyncApiCacheWithName(name string) (*appsync.ApiCache, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appsync.ApiCache:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appsync.ApiCache not found", name)
}

// GetAllAppSyncApiKeyResources retrieves all appsync.ApiKey items from an AWS CloudFormation template
func (t *Template) GetAllAppSyncApiKeyResources() map[string]*appsync.ApiKey {
	results := map[string]*appsync.ApiKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appsync.ApiKey:
			results[name] = resource
		}
	}
	return results
}

// GetAppSyncApiKeyWithName retrieves all appsync.ApiKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppSyncApiKeyWithName(name string) (*appsync.ApiKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appsync.ApiKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appsync.ApiKey not found", name)
}

// GetAllAppSyncDataSourceResources retrieves all appsync.DataSource items from an AWS CloudFormation template
func (t *Template) GetAllAppSyncDataSourceResources() map[string]*appsync.DataSource {
	results := map[string]*appsync.DataSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appsync.DataSource:
			results[name] = resource
		}
	}
	return results
}

// GetAppSyncDataSourceWithName retrieves all appsync.DataSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppSyncDataSourceWithName(name string) (*appsync.DataSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appsync.DataSource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appsync.DataSource not found", name)
}

// GetAllAppSyncFunctionConfigurationResources retrieves all appsync.FunctionConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllAppSyncFunctionConfigurationResources() map[string]*appsync.FunctionConfiguration {
	results := map[string]*appsync.FunctionConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appsync.FunctionConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetAppSyncFunctionConfigurationWithName retrieves all appsync.FunctionConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppSyncFunctionConfigurationWithName(name string) (*appsync.FunctionConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appsync.FunctionConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appsync.FunctionConfiguration not found", name)
}

// GetAllAppSyncGraphQLApiResources retrieves all appsync.GraphQLApi items from an AWS CloudFormation template
func (t *Template) GetAllAppSyncGraphQLApiResources() map[string]*appsync.GraphQLApi {
	results := map[string]*appsync.GraphQLApi{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appsync.GraphQLApi:
			results[name] = resource
		}
	}
	return results
}

// GetAppSyncGraphQLApiWithName retrieves all appsync.GraphQLApi items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppSyncGraphQLApiWithName(name string) (*appsync.GraphQLApi, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appsync.GraphQLApi:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appsync.GraphQLApi not found", name)
}

// GetAllAppSyncGraphQLSchemaResources retrieves all appsync.GraphQLSchema items from an AWS CloudFormation template
func (t *Template) GetAllAppSyncGraphQLSchemaResources() map[string]*appsync.GraphQLSchema {
	results := map[string]*appsync.GraphQLSchema{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appsync.GraphQLSchema:
			results[name] = resource
		}
	}
	return results
}

// GetAppSyncGraphQLSchemaWithName retrieves all appsync.GraphQLSchema items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppSyncGraphQLSchemaWithName(name string) (*appsync.GraphQLSchema, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appsync.GraphQLSchema:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appsync.GraphQLSchema not found", name)
}

// GetAllAppSyncResolverResources retrieves all appsync.Resolver items from an AWS CloudFormation template
func (t *Template) GetAllAppSyncResolverResources() map[string]*appsync.Resolver {
	results := map[string]*appsync.Resolver{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *appsync.Resolver:
			results[name] = resource
		}
	}
	return results
}

// GetAppSyncResolverWithName retrieves all appsync.Resolver items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAppSyncResolverWithName(name string) (*appsync.Resolver, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *appsync.Resolver:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type appsync.Resolver not found", name)
}

// GetAllApplicationAutoScalingScalableTargetResources retrieves all applicationautoscaling.ScalableTarget items from an AWS CloudFormation template
func (t *Template) GetAllApplicationAutoScalingScalableTargetResources() map[string]*applicationautoscaling.ScalableTarget {
	results := map[string]*applicationautoscaling.ScalableTarget{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *applicationautoscaling.ScalableTarget:
			results[name] = resource
		}
	}
	return results
}

// GetApplicationAutoScalingScalableTargetWithName retrieves all applicationautoscaling.ScalableTarget items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApplicationAutoScalingScalableTargetWithName(name string) (*applicationautoscaling.ScalableTarget, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *applicationautoscaling.ScalableTarget:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type applicationautoscaling.ScalableTarget not found", name)
}

// GetAllApplicationAutoScalingScalingPolicyResources retrieves all applicationautoscaling.ScalingPolicy items from an AWS CloudFormation template
func (t *Template) GetAllApplicationAutoScalingScalingPolicyResources() map[string]*applicationautoscaling.ScalingPolicy {
	results := map[string]*applicationautoscaling.ScalingPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *applicationautoscaling.ScalingPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetApplicationAutoScalingScalingPolicyWithName retrieves all applicationautoscaling.ScalingPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApplicationAutoScalingScalingPolicyWithName(name string) (*applicationautoscaling.ScalingPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *applicationautoscaling.ScalingPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type applicationautoscaling.ScalingPolicy not found", name)
}

// GetAllApplicationInsightsApplicationResources retrieves all applicationinsights.Application items from an AWS CloudFormation template
func (t *Template) GetAllApplicationInsightsApplicationResources() map[string]*applicationinsights.Application {
	results := map[string]*applicationinsights.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *applicationinsights.Application:
			results[name] = resource
		}
	}
	return results
}

// GetApplicationInsightsApplicationWithName retrieves all applicationinsights.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetApplicationInsightsApplicationWithName(name string) (*applicationinsights.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *applicationinsights.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type applicationinsights.Application not found", name)
}

// GetAllAthenaDataCatalogResources retrieves all athena.DataCatalog items from an AWS CloudFormation template
func (t *Template) GetAllAthenaDataCatalogResources() map[string]*athena.DataCatalog {
	results := map[string]*athena.DataCatalog{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *athena.DataCatalog:
			results[name] = resource
		}
	}
	return results
}

// GetAthenaDataCatalogWithName retrieves all athena.DataCatalog items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAthenaDataCatalogWithName(name string) (*athena.DataCatalog, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *athena.DataCatalog:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type athena.DataCatalog not found", name)
}

// GetAllAthenaNamedQueryResources retrieves all athena.NamedQuery items from an AWS CloudFormation template
func (t *Template) GetAllAthenaNamedQueryResources() map[string]*athena.NamedQuery {
	results := map[string]*athena.NamedQuery{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *athena.NamedQuery:
			results[name] = resource
		}
	}
	return results
}

// GetAthenaNamedQueryWithName retrieves all athena.NamedQuery items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAthenaNamedQueryWithName(name string) (*athena.NamedQuery, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *athena.NamedQuery:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type athena.NamedQuery not found", name)
}

// GetAllAthenaPreparedStatementResources retrieves all athena.PreparedStatement items from an AWS CloudFormation template
func (t *Template) GetAllAthenaPreparedStatementResources() map[string]*athena.PreparedStatement {
	results := map[string]*athena.PreparedStatement{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *athena.PreparedStatement:
			results[name] = resource
		}
	}
	return results
}

// GetAthenaPreparedStatementWithName retrieves all athena.PreparedStatement items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAthenaPreparedStatementWithName(name string) (*athena.PreparedStatement, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *athena.PreparedStatement:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type athena.PreparedStatement not found", name)
}

// GetAllAthenaWorkGroupResources retrieves all athena.WorkGroup items from an AWS CloudFormation template
func (t *Template) GetAllAthenaWorkGroupResources() map[string]*athena.WorkGroup {
	results := map[string]*athena.WorkGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *athena.WorkGroup:
			results[name] = resource
		}
	}
	return results
}

// GetAthenaWorkGroupWithName retrieves all athena.WorkGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAthenaWorkGroupWithName(name string) (*athena.WorkGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *athena.WorkGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type athena.WorkGroup not found", name)
}

// GetAllAuditManagerAssessmentResources retrieves all auditmanager.Assessment items from an AWS CloudFormation template
func (t *Template) GetAllAuditManagerAssessmentResources() map[string]*auditmanager.Assessment {
	results := map[string]*auditmanager.Assessment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *auditmanager.Assessment:
			results[name] = resource
		}
	}
	return results
}

// GetAuditManagerAssessmentWithName retrieves all auditmanager.Assessment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAuditManagerAssessmentWithName(name string) (*auditmanager.Assessment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *auditmanager.Assessment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type auditmanager.Assessment not found", name)
}

// GetAllAutoScalingAutoScalingGroupResources retrieves all autoscaling.AutoScalingGroup items from an AWS CloudFormation template
func (t *Template) GetAllAutoScalingAutoScalingGroupResources() map[string]*autoscaling.AutoScalingGroup {
	results := map[string]*autoscaling.AutoScalingGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *autoscaling.AutoScalingGroup:
			results[name] = resource
		}
	}
	return results
}

// GetAutoScalingAutoScalingGroupWithName retrieves all autoscaling.AutoScalingGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAutoScalingAutoScalingGroupWithName(name string) (*autoscaling.AutoScalingGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *autoscaling.AutoScalingGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type autoscaling.AutoScalingGroup not found", name)
}

// GetAllAutoScalingLaunchConfigurationResources retrieves all autoscaling.LaunchConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllAutoScalingLaunchConfigurationResources() map[string]*autoscaling.LaunchConfiguration {
	results := map[string]*autoscaling.LaunchConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *autoscaling.LaunchConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetAutoScalingLaunchConfigurationWithName retrieves all autoscaling.LaunchConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAutoScalingLaunchConfigurationWithName(name string) (*autoscaling.LaunchConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *autoscaling.LaunchConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type autoscaling.LaunchConfiguration not found", name)
}

// GetAllAutoScalingLifecycleHookResources retrieves all autoscaling.LifecycleHook items from an AWS CloudFormation template
func (t *Template) GetAllAutoScalingLifecycleHookResources() map[string]*autoscaling.LifecycleHook {
	results := map[string]*autoscaling.LifecycleHook{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *autoscaling.LifecycleHook:
			results[name] = resource
		}
	}
	return results
}

// GetAutoScalingLifecycleHookWithName retrieves all autoscaling.LifecycleHook items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAutoScalingLifecycleHookWithName(name string) (*autoscaling.LifecycleHook, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *autoscaling.LifecycleHook:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type autoscaling.LifecycleHook not found", name)
}

// GetAllAutoScalingScalingPolicyResources retrieves all autoscaling.ScalingPolicy items from an AWS CloudFormation template
func (t *Template) GetAllAutoScalingScalingPolicyResources() map[string]*autoscaling.ScalingPolicy {
	results := map[string]*autoscaling.ScalingPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *autoscaling.ScalingPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetAutoScalingScalingPolicyWithName retrieves all autoscaling.ScalingPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAutoScalingScalingPolicyWithName(name string) (*autoscaling.ScalingPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *autoscaling.ScalingPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type autoscaling.ScalingPolicy not found", name)
}

// GetAllAutoScalingScheduledActionResources retrieves all autoscaling.ScheduledAction items from an AWS CloudFormation template
func (t *Template) GetAllAutoScalingScheduledActionResources() map[string]*autoscaling.ScheduledAction {
	results := map[string]*autoscaling.ScheduledAction{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *autoscaling.ScheduledAction:
			results[name] = resource
		}
	}
	return results
}

// GetAutoScalingScheduledActionWithName retrieves all autoscaling.ScheduledAction items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAutoScalingScheduledActionWithName(name string) (*autoscaling.ScheduledAction, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *autoscaling.ScheduledAction:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type autoscaling.ScheduledAction not found", name)
}

// GetAllAutoScalingWarmPoolResources retrieves all autoscaling.WarmPool items from an AWS CloudFormation template
func (t *Template) GetAllAutoScalingWarmPoolResources() map[string]*autoscaling.WarmPool {
	results := map[string]*autoscaling.WarmPool{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *autoscaling.WarmPool:
			results[name] = resource
		}
	}
	return results
}

// GetAutoScalingWarmPoolWithName retrieves all autoscaling.WarmPool items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAutoScalingWarmPoolWithName(name string) (*autoscaling.WarmPool, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *autoscaling.WarmPool:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type autoscaling.WarmPool not found", name)
}

// GetAllAutoScalingPlansScalingPlanResources retrieves all autoscalingplans.ScalingPlan items from an AWS CloudFormation template
func (t *Template) GetAllAutoScalingPlansScalingPlanResources() map[string]*autoscalingplans.ScalingPlan {
	results := map[string]*autoscalingplans.ScalingPlan{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *autoscalingplans.ScalingPlan:
			results[name] = resource
		}
	}
	return results
}

// GetAutoScalingPlansScalingPlanWithName retrieves all autoscalingplans.ScalingPlan items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAutoScalingPlansScalingPlanWithName(name string) (*autoscalingplans.ScalingPlan, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *autoscalingplans.ScalingPlan:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type autoscalingplans.ScalingPlan not found", name)
}

// GetAllBackupBackupPlanResources retrieves all backup.BackupPlan items from an AWS CloudFormation template
func (t *Template) GetAllBackupBackupPlanResources() map[string]*backup.BackupPlan {
	results := map[string]*backup.BackupPlan{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *backup.BackupPlan:
			results[name] = resource
		}
	}
	return results
}

// GetBackupBackupPlanWithName retrieves all backup.BackupPlan items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBackupBackupPlanWithName(name string) (*backup.BackupPlan, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *backup.BackupPlan:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type backup.BackupPlan not found", name)
}

// GetAllBackupBackupSelectionResources retrieves all backup.BackupSelection items from an AWS CloudFormation template
func (t *Template) GetAllBackupBackupSelectionResources() map[string]*backup.BackupSelection {
	results := map[string]*backup.BackupSelection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *backup.BackupSelection:
			results[name] = resource
		}
	}
	return results
}

// GetBackupBackupSelectionWithName retrieves all backup.BackupSelection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBackupBackupSelectionWithName(name string) (*backup.BackupSelection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *backup.BackupSelection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type backup.BackupSelection not found", name)
}

// GetAllBackupBackupVaultResources retrieves all backup.BackupVault items from an AWS CloudFormation template
func (t *Template) GetAllBackupBackupVaultResources() map[string]*backup.BackupVault {
	results := map[string]*backup.BackupVault{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *backup.BackupVault:
			results[name] = resource
		}
	}
	return results
}

// GetBackupBackupVaultWithName retrieves all backup.BackupVault items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBackupBackupVaultWithName(name string) (*backup.BackupVault, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *backup.BackupVault:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type backup.BackupVault not found", name)
}

// GetAllBackupFrameworkResources retrieves all backup.Framework items from an AWS CloudFormation template
func (t *Template) GetAllBackupFrameworkResources() map[string]*backup.Framework {
	results := map[string]*backup.Framework{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *backup.Framework:
			results[name] = resource
		}
	}
	return results
}

// GetBackupFrameworkWithName retrieves all backup.Framework items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBackupFrameworkWithName(name string) (*backup.Framework, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *backup.Framework:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type backup.Framework not found", name)
}

// GetAllBackupReportPlanResources retrieves all backup.ReportPlan items from an AWS CloudFormation template
func (t *Template) GetAllBackupReportPlanResources() map[string]*backup.ReportPlan {
	results := map[string]*backup.ReportPlan{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *backup.ReportPlan:
			results[name] = resource
		}
	}
	return results
}

// GetBackupReportPlanWithName retrieves all backup.ReportPlan items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBackupReportPlanWithName(name string) (*backup.ReportPlan, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *backup.ReportPlan:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type backup.ReportPlan not found", name)
}

// GetAllBatchComputeEnvironmentResources retrieves all batch.ComputeEnvironment items from an AWS CloudFormation template
func (t *Template) GetAllBatchComputeEnvironmentResources() map[string]*batch.ComputeEnvironment {
	results := map[string]*batch.ComputeEnvironment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *batch.ComputeEnvironment:
			results[name] = resource
		}
	}
	return results
}

// GetBatchComputeEnvironmentWithName retrieves all batch.ComputeEnvironment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBatchComputeEnvironmentWithName(name string) (*batch.ComputeEnvironment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *batch.ComputeEnvironment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type batch.ComputeEnvironment not found", name)
}

// GetAllBatchJobDefinitionResources retrieves all batch.JobDefinition items from an AWS CloudFormation template
func (t *Template) GetAllBatchJobDefinitionResources() map[string]*batch.JobDefinition {
	results := map[string]*batch.JobDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *batch.JobDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetBatchJobDefinitionWithName retrieves all batch.JobDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBatchJobDefinitionWithName(name string) (*batch.JobDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *batch.JobDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type batch.JobDefinition not found", name)
}

// GetAllBatchJobQueueResources retrieves all batch.JobQueue items from an AWS CloudFormation template
func (t *Template) GetAllBatchJobQueueResources() map[string]*batch.JobQueue {
	results := map[string]*batch.JobQueue{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *batch.JobQueue:
			results[name] = resource
		}
	}
	return results
}

// GetBatchJobQueueWithName retrieves all batch.JobQueue items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBatchJobQueueWithName(name string) (*batch.JobQueue, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *batch.JobQueue:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type batch.JobQueue not found", name)
}

// GetAllBatchSchedulingPolicyResources retrieves all batch.SchedulingPolicy items from an AWS CloudFormation template
func (t *Template) GetAllBatchSchedulingPolicyResources() map[string]*batch.SchedulingPolicy {
	results := map[string]*batch.SchedulingPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *batch.SchedulingPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetBatchSchedulingPolicyWithName retrieves all batch.SchedulingPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBatchSchedulingPolicyWithName(name string) (*batch.SchedulingPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *batch.SchedulingPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type batch.SchedulingPolicy not found", name)
}

// GetAllBudgetsBudgetResources retrieves all budgets.Budget items from an AWS CloudFormation template
func (t *Template) GetAllBudgetsBudgetResources() map[string]*budgets.Budget {
	results := map[string]*budgets.Budget{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *budgets.Budget:
			results[name] = resource
		}
	}
	return results
}

// GetBudgetsBudgetWithName retrieves all budgets.Budget items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBudgetsBudgetWithName(name string) (*budgets.Budget, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *budgets.Budget:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type budgets.Budget not found", name)
}

// GetAllBudgetsBudgetsActionResources retrieves all budgets.BudgetsAction items from an AWS CloudFormation template
func (t *Template) GetAllBudgetsBudgetsActionResources() map[string]*budgets.BudgetsAction {
	results := map[string]*budgets.BudgetsAction{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *budgets.BudgetsAction:
			results[name] = resource
		}
	}
	return results
}

// GetBudgetsBudgetsActionWithName retrieves all budgets.BudgetsAction items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetBudgetsBudgetsActionWithName(name string) (*budgets.BudgetsAction, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *budgets.BudgetsAction:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type budgets.BudgetsAction not found", name)
}

// GetAllCEAnomalyMonitorResources retrieves all ce.AnomalyMonitor items from an AWS CloudFormation template
func (t *Template) GetAllCEAnomalyMonitorResources() map[string]*ce.AnomalyMonitor {
	results := map[string]*ce.AnomalyMonitor{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ce.AnomalyMonitor:
			results[name] = resource
		}
	}
	return results
}

// GetCEAnomalyMonitorWithName retrieves all ce.AnomalyMonitor items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCEAnomalyMonitorWithName(name string) (*ce.AnomalyMonitor, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ce.AnomalyMonitor:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ce.AnomalyMonitor not found", name)
}

// GetAllCEAnomalySubscriptionResources retrieves all ce.AnomalySubscription items from an AWS CloudFormation template
func (t *Template) GetAllCEAnomalySubscriptionResources() map[string]*ce.AnomalySubscription {
	results := map[string]*ce.AnomalySubscription{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ce.AnomalySubscription:
			results[name] = resource
		}
	}
	return results
}

// GetCEAnomalySubscriptionWithName retrieves all ce.AnomalySubscription items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCEAnomalySubscriptionWithName(name string) (*ce.AnomalySubscription, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ce.AnomalySubscription:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ce.AnomalySubscription not found", name)
}

// GetAllCECostCategoryResources retrieves all ce.CostCategory items from an AWS CloudFormation template
func (t *Template) GetAllCECostCategoryResources() map[string]*ce.CostCategory {
	results := map[string]*ce.CostCategory{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ce.CostCategory:
			results[name] = resource
		}
	}
	return results
}

// GetCECostCategoryWithName retrieves all ce.CostCategory items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCECostCategoryWithName(name string) (*ce.CostCategory, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ce.CostCategory:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ce.CostCategory not found", name)
}

// GetAllCURReportDefinitionResources retrieves all cur.ReportDefinition items from an AWS CloudFormation template
func (t *Template) GetAllCURReportDefinitionResources() map[string]*cur.ReportDefinition {
	results := map[string]*cur.ReportDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cur.ReportDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetCURReportDefinitionWithName retrieves all cur.ReportDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCURReportDefinitionWithName(name string) (*cur.ReportDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cur.ReportDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cur.ReportDefinition not found", name)
}

// GetAllCassandraKeyspaceResources retrieves all cassandra.Keyspace items from an AWS CloudFormation template
func (t *Template) GetAllCassandraKeyspaceResources() map[string]*cassandra.Keyspace {
	results := map[string]*cassandra.Keyspace{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cassandra.Keyspace:
			results[name] = resource
		}
	}
	return results
}

// GetCassandraKeyspaceWithName retrieves all cassandra.Keyspace items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCassandraKeyspaceWithName(name string) (*cassandra.Keyspace, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cassandra.Keyspace:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cassandra.Keyspace not found", name)
}

// GetAllCassandraTableResources retrieves all cassandra.Table items from an AWS CloudFormation template
func (t *Template) GetAllCassandraTableResources() map[string]*cassandra.Table {
	results := map[string]*cassandra.Table{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cassandra.Table:
			results[name] = resource
		}
	}
	return results
}

// GetCassandraTableWithName retrieves all cassandra.Table items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCassandraTableWithName(name string) (*cassandra.Table, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cassandra.Table:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cassandra.Table not found", name)
}

// GetAllCertificateManagerAccountResources retrieves all certificatemanager.Account items from an AWS CloudFormation template
func (t *Template) GetAllCertificateManagerAccountResources() map[string]*certificatemanager.Account {
	results := map[string]*certificatemanager.Account{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *certificatemanager.Account:
			results[name] = resource
		}
	}
	return results
}

// GetCertificateManagerAccountWithName retrieves all certificatemanager.Account items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCertificateManagerAccountWithName(name string) (*certificatemanager.Account, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *certificatemanager.Account:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type certificatemanager.Account not found", name)
}

// GetAllCertificateManagerCertificateResources retrieves all certificatemanager.Certificate items from an AWS CloudFormation template
func (t *Template) GetAllCertificateManagerCertificateResources() map[string]*certificatemanager.Certificate {
	results := map[string]*certificatemanager.Certificate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *certificatemanager.Certificate:
			results[name] = resource
		}
	}
	return results
}

// GetCertificateManagerCertificateWithName retrieves all certificatemanager.Certificate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCertificateManagerCertificateWithName(name string) (*certificatemanager.Certificate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *certificatemanager.Certificate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type certificatemanager.Certificate not found", name)
}

// GetAllChatbotSlackChannelConfigurationResources retrieves all chatbot.SlackChannelConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllChatbotSlackChannelConfigurationResources() map[string]*chatbot.SlackChannelConfiguration {
	results := map[string]*chatbot.SlackChannelConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *chatbot.SlackChannelConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetChatbotSlackChannelConfigurationWithName retrieves all chatbot.SlackChannelConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetChatbotSlackChannelConfigurationWithName(name string) (*chatbot.SlackChannelConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *chatbot.SlackChannelConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type chatbot.SlackChannelConfiguration not found", name)
}

// GetAllCloud9EnvironmentEC2Resources retrieves all cloud9.EnvironmentEC2 items from an AWS CloudFormation template
func (t *Template) GetAllCloud9EnvironmentEC2Resources() map[string]*cloud9.EnvironmentEC2 {
	results := map[string]*cloud9.EnvironmentEC2{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloud9.EnvironmentEC2:
			results[name] = resource
		}
	}
	return results
}

// GetCloud9EnvironmentEC2WithName retrieves all cloud9.EnvironmentEC2 items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloud9EnvironmentEC2WithName(name string) (*cloud9.EnvironmentEC2, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloud9.EnvironmentEC2:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloud9.EnvironmentEC2 not found", name)
}

// GetAllCloudFormationCustomResourceResources retrieves all cloudformation.CustomResource items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationCustomResourceResources() map[string]*cloudformation.CustomResource {
	results := map[string]*cloudformation.CustomResource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.CustomResource:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationCustomResourceWithName retrieves all cloudformation.CustomResource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationCustomResourceWithName(name string) (*cloudformation.CustomResource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.CustomResource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.CustomResource not found", name)
}

// GetAllCloudFormationMacroResources retrieves all cloudformation.Macro items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationMacroResources() map[string]*cloudformation.Macro {
	results := map[string]*cloudformation.Macro{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.Macro:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationMacroWithName retrieves all cloudformation.Macro items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationMacroWithName(name string) (*cloudformation.Macro, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.Macro:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.Macro not found", name)
}

// GetAllCloudFormationModuleDefaultVersionResources retrieves all cloudformation.ModuleDefaultVersion items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationModuleDefaultVersionResources() map[string]*cloudformation.ModuleDefaultVersion {
	results := map[string]*cloudformation.ModuleDefaultVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.ModuleDefaultVersion:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationModuleDefaultVersionWithName retrieves all cloudformation.ModuleDefaultVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationModuleDefaultVersionWithName(name string) (*cloudformation.ModuleDefaultVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.ModuleDefaultVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.ModuleDefaultVersion not found", name)
}

// GetAllCloudFormationModuleVersionResources retrieves all cloudformation.ModuleVersion items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationModuleVersionResources() map[string]*cloudformation.ModuleVersion {
	results := map[string]*cloudformation.ModuleVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.ModuleVersion:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationModuleVersionWithName retrieves all cloudformation.ModuleVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationModuleVersionWithName(name string) (*cloudformation.ModuleVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.ModuleVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.ModuleVersion not found", name)
}

// GetAllCloudFormationPublicTypeVersionResources retrieves all cloudformation.PublicTypeVersion items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationPublicTypeVersionResources() map[string]*cloudformation.PublicTypeVersion {
	results := map[string]*cloudformation.PublicTypeVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.PublicTypeVersion:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationPublicTypeVersionWithName retrieves all cloudformation.PublicTypeVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationPublicTypeVersionWithName(name string) (*cloudformation.PublicTypeVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.PublicTypeVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.PublicTypeVersion not found", name)
}

// GetAllCloudFormationPublisherResources retrieves all cloudformation.Publisher items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationPublisherResources() map[string]*cloudformation.Publisher {
	results := map[string]*cloudformation.Publisher{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.Publisher:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationPublisherWithName retrieves all cloudformation.Publisher items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationPublisherWithName(name string) (*cloudformation.Publisher, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.Publisher:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.Publisher not found", name)
}

// GetAllCloudFormationResourceDefaultVersionResources retrieves all cloudformation.ResourceDefaultVersion items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationResourceDefaultVersionResources() map[string]*cloudformation.ResourceDefaultVersion {
	results := map[string]*cloudformation.ResourceDefaultVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.ResourceDefaultVersion:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationResourceDefaultVersionWithName retrieves all cloudformation.ResourceDefaultVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationResourceDefaultVersionWithName(name string) (*cloudformation.ResourceDefaultVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.ResourceDefaultVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.ResourceDefaultVersion not found", name)
}

// GetAllCloudFormationResourceVersionResources retrieves all cloudformation.ResourceVersion items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationResourceVersionResources() map[string]*cloudformation.ResourceVersion {
	results := map[string]*cloudformation.ResourceVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.ResourceVersion:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationResourceVersionWithName retrieves all cloudformation.ResourceVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationResourceVersionWithName(name string) (*cloudformation.ResourceVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.ResourceVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.ResourceVersion not found", name)
}

// GetAllCloudFormationStackResources retrieves all cloudformation.Stack items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationStackResources() map[string]*cloudformation.Stack {
	results := map[string]*cloudformation.Stack{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.Stack:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationStackWithName retrieves all cloudformation.Stack items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationStackWithName(name string) (*cloudformation.Stack, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.Stack:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.Stack not found", name)
}

// GetAllCloudFormationStackSetResources retrieves all cloudformation.StackSet items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationStackSetResources() map[string]*cloudformation.StackSet {
	results := map[string]*cloudformation.StackSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.StackSet:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationStackSetWithName retrieves all cloudformation.StackSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationStackSetWithName(name string) (*cloudformation.StackSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.StackSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.StackSet not found", name)
}

// GetAllCloudFormationTypeActivationResources retrieves all cloudformation.TypeActivation items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationTypeActivationResources() map[string]*cloudformation.TypeActivation {
	results := map[string]*cloudformation.TypeActivation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.TypeActivation:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationTypeActivationWithName retrieves all cloudformation.TypeActivation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationTypeActivationWithName(name string) (*cloudformation.TypeActivation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.TypeActivation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.TypeActivation not found", name)
}

// GetAllCloudFormationWaitConditionResources retrieves all cloudformation.WaitCondition items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationWaitConditionResources() map[string]*cloudformation.WaitCondition {
	results := map[string]*cloudformation.WaitCondition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.WaitCondition:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationWaitConditionWithName retrieves all cloudformation.WaitCondition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationWaitConditionWithName(name string) (*cloudformation.WaitCondition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.WaitCondition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.WaitCondition not found", name)
}

// GetAllCloudFormationWaitConditionHandleResources retrieves all cloudformation.WaitConditionHandle items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationWaitConditionHandleResources() map[string]*cloudformation.WaitConditionHandle {
	results := map[string]*cloudformation.WaitConditionHandle{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.WaitConditionHandle:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationWaitConditionHandleWithName retrieves all cloudformation.WaitConditionHandle items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationWaitConditionHandleWithName(name string) (*cloudformation.WaitConditionHandle, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.WaitConditionHandle:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.WaitConditionHandle not found", name)
}

// GetAllCloudFrontCachePolicyResources retrieves all cloudfront.CachePolicy items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontCachePolicyResources() map[string]*cloudfront.CachePolicy {
	results := map[string]*cloudfront.CachePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.CachePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontCachePolicyWithName retrieves all cloudfront.CachePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontCachePolicyWithName(name string) (*cloudfront.CachePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.CachePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.CachePolicy not found", name)
}

// GetAllCloudFrontCloudFrontOriginAccessIdentityResources retrieves all cloudfront.CloudFrontOriginAccessIdentity items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontCloudFrontOriginAccessIdentityResources() map[string]*cloudfront.CloudFrontOriginAccessIdentity {
	results := map[string]*cloudfront.CloudFrontOriginAccessIdentity{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.CloudFrontOriginAccessIdentity:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontCloudFrontOriginAccessIdentityWithName retrieves all cloudfront.CloudFrontOriginAccessIdentity items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontCloudFrontOriginAccessIdentityWithName(name string) (*cloudfront.CloudFrontOriginAccessIdentity, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.CloudFrontOriginAccessIdentity:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.CloudFrontOriginAccessIdentity not found", name)
}

// GetAllCloudFrontDistributionResources retrieves all cloudfront.Distribution items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontDistributionResources() map[string]*cloudfront.Distribution {
	results := map[string]*cloudfront.Distribution{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.Distribution:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontDistributionWithName retrieves all cloudfront.Distribution items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontDistributionWithName(name string) (*cloudfront.Distribution, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.Distribution:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.Distribution not found", name)
}

// GetAllCloudFrontFunctionResources retrieves all cloudfront.Function items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontFunctionResources() map[string]*cloudfront.Function {
	results := map[string]*cloudfront.Function{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.Function:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontFunctionWithName retrieves all cloudfront.Function items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontFunctionWithName(name string) (*cloudfront.Function, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.Function:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.Function not found", name)
}

// GetAllCloudFrontKeyGroupResources retrieves all cloudfront.KeyGroup items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontKeyGroupResources() map[string]*cloudfront.KeyGroup {
	results := map[string]*cloudfront.KeyGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.KeyGroup:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontKeyGroupWithName retrieves all cloudfront.KeyGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontKeyGroupWithName(name string) (*cloudfront.KeyGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.KeyGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.KeyGroup not found", name)
}

// GetAllCloudFrontOriginRequestPolicyResources retrieves all cloudfront.OriginRequestPolicy items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontOriginRequestPolicyResources() map[string]*cloudfront.OriginRequestPolicy {
	results := map[string]*cloudfront.OriginRequestPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.OriginRequestPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontOriginRequestPolicyWithName retrieves all cloudfront.OriginRequestPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontOriginRequestPolicyWithName(name string) (*cloudfront.OriginRequestPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.OriginRequestPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.OriginRequestPolicy not found", name)
}

// GetAllCloudFrontPublicKeyResources retrieves all cloudfront.PublicKey items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontPublicKeyResources() map[string]*cloudfront.PublicKey {
	results := map[string]*cloudfront.PublicKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.PublicKey:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontPublicKeyWithName retrieves all cloudfront.PublicKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontPublicKeyWithName(name string) (*cloudfront.PublicKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.PublicKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.PublicKey not found", name)
}

// GetAllCloudFrontRealtimeLogConfigResources retrieves all cloudfront.RealtimeLogConfig items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontRealtimeLogConfigResources() map[string]*cloudfront.RealtimeLogConfig {
	results := map[string]*cloudfront.RealtimeLogConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.RealtimeLogConfig:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontRealtimeLogConfigWithName retrieves all cloudfront.RealtimeLogConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontRealtimeLogConfigWithName(name string) (*cloudfront.RealtimeLogConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.RealtimeLogConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.RealtimeLogConfig not found", name)
}

// GetAllCloudFrontResponseHeadersPolicyResources retrieves all cloudfront.ResponseHeadersPolicy items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontResponseHeadersPolicyResources() map[string]*cloudfront.ResponseHeadersPolicy {
	results := map[string]*cloudfront.ResponseHeadersPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.ResponseHeadersPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontResponseHeadersPolicyWithName retrieves all cloudfront.ResponseHeadersPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontResponseHeadersPolicyWithName(name string) (*cloudfront.ResponseHeadersPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.ResponseHeadersPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.ResponseHeadersPolicy not found", name)
}

// GetAllCloudFrontStreamingDistributionResources retrieves all cloudfront.StreamingDistribution items from an AWS CloudFormation template
func (t *Template) GetAllCloudFrontStreamingDistributionResources() map[string]*cloudfront.StreamingDistribution {
	results := map[string]*cloudfront.StreamingDistribution{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudfront.StreamingDistribution:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFrontStreamingDistributionWithName retrieves all cloudfront.StreamingDistribution items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFrontStreamingDistributionWithName(name string) (*cloudfront.StreamingDistribution, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudfront.StreamingDistribution:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudfront.StreamingDistribution not found", name)
}

// GetAllCloudTrailTrailResources retrieves all cloudtrail.Trail items from an AWS CloudFormation template
func (t *Template) GetAllCloudTrailTrailResources() map[string]*cloudtrail.Trail {
	results := map[string]*cloudtrail.Trail{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudtrail.Trail:
			results[name] = resource
		}
	}
	return results
}

// GetCloudTrailTrailWithName retrieves all cloudtrail.Trail items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudTrailTrailWithName(name string) (*cloudtrail.Trail, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudtrail.Trail:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudtrail.Trail not found", name)
}

// GetAllCloudWatchAlarmResources retrieves all cloudwatch.Alarm items from an AWS CloudFormation template
func (t *Template) GetAllCloudWatchAlarmResources() map[string]*cloudwatch.Alarm {
	results := map[string]*cloudwatch.Alarm{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudwatch.Alarm:
			results[name] = resource
		}
	}
	return results
}

// GetCloudWatchAlarmWithName retrieves all cloudwatch.Alarm items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudWatchAlarmWithName(name string) (*cloudwatch.Alarm, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudwatch.Alarm:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudwatch.Alarm not found", name)
}

// GetAllCloudWatchAnomalyDetectorResources retrieves all cloudwatch.AnomalyDetector items from an AWS CloudFormation template
func (t *Template) GetAllCloudWatchAnomalyDetectorResources() map[string]*cloudwatch.AnomalyDetector {
	results := map[string]*cloudwatch.AnomalyDetector{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudwatch.AnomalyDetector:
			results[name] = resource
		}
	}
	return results
}

// GetCloudWatchAnomalyDetectorWithName retrieves all cloudwatch.AnomalyDetector items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudWatchAnomalyDetectorWithName(name string) (*cloudwatch.AnomalyDetector, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudwatch.AnomalyDetector:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudwatch.AnomalyDetector not found", name)
}

// GetAllCloudWatchCompositeAlarmResources retrieves all cloudwatch.CompositeAlarm items from an AWS CloudFormation template
func (t *Template) GetAllCloudWatchCompositeAlarmResources() map[string]*cloudwatch.CompositeAlarm {
	results := map[string]*cloudwatch.CompositeAlarm{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudwatch.CompositeAlarm:
			results[name] = resource
		}
	}
	return results
}

// GetCloudWatchCompositeAlarmWithName retrieves all cloudwatch.CompositeAlarm items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudWatchCompositeAlarmWithName(name string) (*cloudwatch.CompositeAlarm, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudwatch.CompositeAlarm:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudwatch.CompositeAlarm not found", name)
}

// GetAllCloudWatchDashboardResources retrieves all cloudwatch.Dashboard items from an AWS CloudFormation template
func (t *Template) GetAllCloudWatchDashboardResources() map[string]*cloudwatch.Dashboard {
	results := map[string]*cloudwatch.Dashboard{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudwatch.Dashboard:
			results[name] = resource
		}
	}
	return results
}

// GetCloudWatchDashboardWithName retrieves all cloudwatch.Dashboard items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudWatchDashboardWithName(name string) (*cloudwatch.Dashboard, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudwatch.Dashboard:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudwatch.Dashboard not found", name)
}

// GetAllCloudWatchInsightRuleResources retrieves all cloudwatch.InsightRule items from an AWS CloudFormation template
func (t *Template) GetAllCloudWatchInsightRuleResources() map[string]*cloudwatch.InsightRule {
	results := map[string]*cloudwatch.InsightRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudwatch.InsightRule:
			results[name] = resource
		}
	}
	return results
}

// GetCloudWatchInsightRuleWithName retrieves all cloudwatch.InsightRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudWatchInsightRuleWithName(name string) (*cloudwatch.InsightRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudwatch.InsightRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudwatch.InsightRule not found", name)
}

// GetAllCloudWatchMetricStreamResources retrieves all cloudwatch.MetricStream items from an AWS CloudFormation template
func (t *Template) GetAllCloudWatchMetricStreamResources() map[string]*cloudwatch.MetricStream {
	results := map[string]*cloudwatch.MetricStream{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudwatch.MetricStream:
			results[name] = resource
		}
	}
	return results
}

// GetCloudWatchMetricStreamWithName retrieves all cloudwatch.MetricStream items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudWatchMetricStreamWithName(name string) (*cloudwatch.MetricStream, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudwatch.MetricStream:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudwatch.MetricStream not found", name)
}

// GetAllCodeArtifactDomainResources retrieves all codeartifact.Domain items from an AWS CloudFormation template
func (t *Template) GetAllCodeArtifactDomainResources() map[string]*codeartifact.Domain {
	results := map[string]*codeartifact.Domain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codeartifact.Domain:
			results[name] = resource
		}
	}
	return results
}

// GetCodeArtifactDomainWithName retrieves all codeartifact.Domain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeArtifactDomainWithName(name string) (*codeartifact.Domain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codeartifact.Domain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codeartifact.Domain not found", name)
}

// GetAllCodeArtifactRepositoryResources retrieves all codeartifact.Repository items from an AWS CloudFormation template
func (t *Template) GetAllCodeArtifactRepositoryResources() map[string]*codeartifact.Repository {
	results := map[string]*codeartifact.Repository{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codeartifact.Repository:
			results[name] = resource
		}
	}
	return results
}

// GetCodeArtifactRepositoryWithName retrieves all codeartifact.Repository items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeArtifactRepositoryWithName(name string) (*codeartifact.Repository, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codeartifact.Repository:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codeartifact.Repository not found", name)
}

// GetAllCodeBuildProjectResources retrieves all codebuild.Project items from an AWS CloudFormation template
func (t *Template) GetAllCodeBuildProjectResources() map[string]*codebuild.Project {
	results := map[string]*codebuild.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codebuild.Project:
			results[name] = resource
		}
	}
	return results
}

// GetCodeBuildProjectWithName retrieves all codebuild.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeBuildProjectWithName(name string) (*codebuild.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codebuild.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codebuild.Project not found", name)
}

// GetAllCodeBuildReportGroupResources retrieves all codebuild.ReportGroup items from an AWS CloudFormation template
func (t *Template) GetAllCodeBuildReportGroupResources() map[string]*codebuild.ReportGroup {
	results := map[string]*codebuild.ReportGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codebuild.ReportGroup:
			results[name] = resource
		}
	}
	return results
}

// GetCodeBuildReportGroupWithName retrieves all codebuild.ReportGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeBuildReportGroupWithName(name string) (*codebuild.ReportGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codebuild.ReportGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codebuild.ReportGroup not found", name)
}

// GetAllCodeBuildSourceCredentialResources retrieves all codebuild.SourceCredential items from an AWS CloudFormation template
func (t *Template) GetAllCodeBuildSourceCredentialResources() map[string]*codebuild.SourceCredential {
	results := map[string]*codebuild.SourceCredential{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codebuild.SourceCredential:
			results[name] = resource
		}
	}
	return results
}

// GetCodeBuildSourceCredentialWithName retrieves all codebuild.SourceCredential items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeBuildSourceCredentialWithName(name string) (*codebuild.SourceCredential, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codebuild.SourceCredential:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codebuild.SourceCredential not found", name)
}

// GetAllCodeCommitRepositoryResources retrieves all codecommit.Repository items from an AWS CloudFormation template
func (t *Template) GetAllCodeCommitRepositoryResources() map[string]*codecommit.Repository {
	results := map[string]*codecommit.Repository{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codecommit.Repository:
			results[name] = resource
		}
	}
	return results
}

// GetCodeCommitRepositoryWithName retrieves all codecommit.Repository items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeCommitRepositoryWithName(name string) (*codecommit.Repository, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codecommit.Repository:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codecommit.Repository not found", name)
}

// GetAllCodeDeployApplicationResources retrieves all codedeploy.Application items from an AWS CloudFormation template
func (t *Template) GetAllCodeDeployApplicationResources() map[string]*codedeploy.Application {
	results := map[string]*codedeploy.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codedeploy.Application:
			results[name] = resource
		}
	}
	return results
}

// GetCodeDeployApplicationWithName retrieves all codedeploy.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeDeployApplicationWithName(name string) (*codedeploy.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codedeploy.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codedeploy.Application not found", name)
}

// GetAllCodeDeployDeploymentConfigResources retrieves all codedeploy.DeploymentConfig items from an AWS CloudFormation template
func (t *Template) GetAllCodeDeployDeploymentConfigResources() map[string]*codedeploy.DeploymentConfig {
	results := map[string]*codedeploy.DeploymentConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codedeploy.DeploymentConfig:
			results[name] = resource
		}
	}
	return results
}

// GetCodeDeployDeploymentConfigWithName retrieves all codedeploy.DeploymentConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeDeployDeploymentConfigWithName(name string) (*codedeploy.DeploymentConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codedeploy.DeploymentConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codedeploy.DeploymentConfig not found", name)
}

// GetAllCodeDeployDeploymentGroupResources retrieves all codedeploy.DeploymentGroup items from an AWS CloudFormation template
func (t *Template) GetAllCodeDeployDeploymentGroupResources() map[string]*codedeploy.DeploymentGroup {
	results := map[string]*codedeploy.DeploymentGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codedeploy.DeploymentGroup:
			results[name] = resource
		}
	}
	return results
}

// GetCodeDeployDeploymentGroupWithName retrieves all codedeploy.DeploymentGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeDeployDeploymentGroupWithName(name string) (*codedeploy.DeploymentGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codedeploy.DeploymentGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codedeploy.DeploymentGroup not found", name)
}

// GetAllCodeGuruProfilerProfilingGroupResources retrieves all codeguruprofiler.ProfilingGroup items from an AWS CloudFormation template
func (t *Template) GetAllCodeGuruProfilerProfilingGroupResources() map[string]*codeguruprofiler.ProfilingGroup {
	results := map[string]*codeguruprofiler.ProfilingGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codeguruprofiler.ProfilingGroup:
			results[name] = resource
		}
	}
	return results
}

// GetCodeGuruProfilerProfilingGroupWithName retrieves all codeguruprofiler.ProfilingGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeGuruProfilerProfilingGroupWithName(name string) (*codeguruprofiler.ProfilingGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codeguruprofiler.ProfilingGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codeguruprofiler.ProfilingGroup not found", name)
}

// GetAllCodeGuruReviewerRepositoryAssociationResources retrieves all codegurureviewer.RepositoryAssociation items from an AWS CloudFormation template
func (t *Template) GetAllCodeGuruReviewerRepositoryAssociationResources() map[string]*codegurureviewer.RepositoryAssociation {
	results := map[string]*codegurureviewer.RepositoryAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codegurureviewer.RepositoryAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetCodeGuruReviewerRepositoryAssociationWithName retrieves all codegurureviewer.RepositoryAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeGuruReviewerRepositoryAssociationWithName(name string) (*codegurureviewer.RepositoryAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codegurureviewer.RepositoryAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codegurureviewer.RepositoryAssociation not found", name)
}

// GetAllCodePipelineCustomActionTypeResources retrieves all codepipeline.CustomActionType items from an AWS CloudFormation template
func (t *Template) GetAllCodePipelineCustomActionTypeResources() map[string]*codepipeline.CustomActionType {
	results := map[string]*codepipeline.CustomActionType{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codepipeline.CustomActionType:
			results[name] = resource
		}
	}
	return results
}

// GetCodePipelineCustomActionTypeWithName retrieves all codepipeline.CustomActionType items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodePipelineCustomActionTypeWithName(name string) (*codepipeline.CustomActionType, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codepipeline.CustomActionType:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codepipeline.CustomActionType not found", name)
}

// GetAllCodePipelinePipelineResources retrieves all codepipeline.Pipeline items from an AWS CloudFormation template
func (t *Template) GetAllCodePipelinePipelineResources() map[string]*codepipeline.Pipeline {
	results := map[string]*codepipeline.Pipeline{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codepipeline.Pipeline:
			results[name] = resource
		}
	}
	return results
}

// GetCodePipelinePipelineWithName retrieves all codepipeline.Pipeline items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodePipelinePipelineWithName(name string) (*codepipeline.Pipeline, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codepipeline.Pipeline:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codepipeline.Pipeline not found", name)
}

// GetAllCodePipelineWebhookResources retrieves all codepipeline.Webhook items from an AWS CloudFormation template
func (t *Template) GetAllCodePipelineWebhookResources() map[string]*codepipeline.Webhook {
	results := map[string]*codepipeline.Webhook{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codepipeline.Webhook:
			results[name] = resource
		}
	}
	return results
}

// GetCodePipelineWebhookWithName retrieves all codepipeline.Webhook items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodePipelineWebhookWithName(name string) (*codepipeline.Webhook, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codepipeline.Webhook:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codepipeline.Webhook not found", name)
}

// GetAllCodeStarGitHubRepositoryResources retrieves all codestar.GitHubRepository items from an AWS CloudFormation template
func (t *Template) GetAllCodeStarGitHubRepositoryResources() map[string]*codestar.GitHubRepository {
	results := map[string]*codestar.GitHubRepository{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codestar.GitHubRepository:
			results[name] = resource
		}
	}
	return results
}

// GetCodeStarGitHubRepositoryWithName retrieves all codestar.GitHubRepository items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeStarGitHubRepositoryWithName(name string) (*codestar.GitHubRepository, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codestar.GitHubRepository:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codestar.GitHubRepository not found", name)
}

// GetAllCodeStarConnectionsConnectionResources retrieves all codestarconnections.Connection items from an AWS CloudFormation template
func (t *Template) GetAllCodeStarConnectionsConnectionResources() map[string]*codestarconnections.Connection {
	results := map[string]*codestarconnections.Connection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codestarconnections.Connection:
			results[name] = resource
		}
	}
	return results
}

// GetCodeStarConnectionsConnectionWithName retrieves all codestarconnections.Connection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeStarConnectionsConnectionWithName(name string) (*codestarconnections.Connection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codestarconnections.Connection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codestarconnections.Connection not found", name)
}

// GetAllCodeStarNotificationsNotificationRuleResources retrieves all codestarnotifications.NotificationRule items from an AWS CloudFormation template
func (t *Template) GetAllCodeStarNotificationsNotificationRuleResources() map[string]*codestarnotifications.NotificationRule {
	results := map[string]*codestarnotifications.NotificationRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *codestarnotifications.NotificationRule:
			results[name] = resource
		}
	}
	return results
}

// GetCodeStarNotificationsNotificationRuleWithName retrieves all codestarnotifications.NotificationRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCodeStarNotificationsNotificationRuleWithName(name string) (*codestarnotifications.NotificationRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *codestarnotifications.NotificationRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type codestarnotifications.NotificationRule not found", name)
}

// GetAllCognitoIdentityPoolResources retrieves all cognito.IdentityPool items from an AWS CloudFormation template
func (t *Template) GetAllCognitoIdentityPoolResources() map[string]*cognito.IdentityPool {
	results := map[string]*cognito.IdentityPool{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.IdentityPool:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoIdentityPoolWithName retrieves all cognito.IdentityPool items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoIdentityPoolWithName(name string) (*cognito.IdentityPool, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.IdentityPool:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.IdentityPool not found", name)
}

// GetAllCognitoIdentityPoolRoleAttachmentResources retrieves all cognito.IdentityPoolRoleAttachment items from an AWS CloudFormation template
func (t *Template) GetAllCognitoIdentityPoolRoleAttachmentResources() map[string]*cognito.IdentityPoolRoleAttachment {
	results := map[string]*cognito.IdentityPoolRoleAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.IdentityPoolRoleAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoIdentityPoolRoleAttachmentWithName retrieves all cognito.IdentityPoolRoleAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoIdentityPoolRoleAttachmentWithName(name string) (*cognito.IdentityPoolRoleAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.IdentityPoolRoleAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.IdentityPoolRoleAttachment not found", name)
}

// GetAllCognitoUserPoolResources retrieves all cognito.UserPool items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolResources() map[string]*cognito.UserPool {
	results := map[string]*cognito.UserPool{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPool:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolWithName retrieves all cognito.UserPool items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolWithName(name string) (*cognito.UserPool, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPool:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPool not found", name)
}

// GetAllCognitoUserPoolClientResources retrieves all cognito.UserPoolClient items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolClientResources() map[string]*cognito.UserPoolClient {
	results := map[string]*cognito.UserPoolClient{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolClient:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolClientWithName retrieves all cognito.UserPoolClient items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolClientWithName(name string) (*cognito.UserPoolClient, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolClient:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolClient not found", name)
}

// GetAllCognitoUserPoolDomainResources retrieves all cognito.UserPoolDomain items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolDomainResources() map[string]*cognito.UserPoolDomain {
	results := map[string]*cognito.UserPoolDomain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolDomain:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolDomainWithName retrieves all cognito.UserPoolDomain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolDomainWithName(name string) (*cognito.UserPoolDomain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolDomain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolDomain not found", name)
}

// GetAllCognitoUserPoolGroupResources retrieves all cognito.UserPoolGroup items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolGroupResources() map[string]*cognito.UserPoolGroup {
	results := map[string]*cognito.UserPoolGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolGroup:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolGroupWithName retrieves all cognito.UserPoolGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolGroupWithName(name string) (*cognito.UserPoolGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolGroup not found", name)
}

// GetAllCognitoUserPoolIdentityProviderResources retrieves all cognito.UserPoolIdentityProvider items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolIdentityProviderResources() map[string]*cognito.UserPoolIdentityProvider {
	results := map[string]*cognito.UserPoolIdentityProvider{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolIdentityProvider:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolIdentityProviderWithName retrieves all cognito.UserPoolIdentityProvider items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolIdentityProviderWithName(name string) (*cognito.UserPoolIdentityProvider, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolIdentityProvider:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolIdentityProvider not found", name)
}

// GetAllCognitoUserPoolResourceServerResources retrieves all cognito.UserPoolResourceServer items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolResourceServerResources() map[string]*cognito.UserPoolResourceServer {
	results := map[string]*cognito.UserPoolResourceServer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolResourceServer:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolResourceServerWithName retrieves all cognito.UserPoolResourceServer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolResourceServerWithName(name string) (*cognito.UserPoolResourceServer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolResourceServer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolResourceServer not found", name)
}

// GetAllCognitoUserPoolRiskConfigurationAttachmentResources retrieves all cognito.UserPoolRiskConfigurationAttachment items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolRiskConfigurationAttachmentResources() map[string]*cognito.UserPoolRiskConfigurationAttachment {
	results := map[string]*cognito.UserPoolRiskConfigurationAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolRiskConfigurationAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolRiskConfigurationAttachmentWithName retrieves all cognito.UserPoolRiskConfigurationAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolRiskConfigurationAttachmentWithName(name string) (*cognito.UserPoolRiskConfigurationAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolRiskConfigurationAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolRiskConfigurationAttachment not found", name)
}

// GetAllCognitoUserPoolUICustomizationAttachmentResources retrieves all cognito.UserPoolUICustomizationAttachment items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolUICustomizationAttachmentResources() map[string]*cognito.UserPoolUICustomizationAttachment {
	results := map[string]*cognito.UserPoolUICustomizationAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolUICustomizationAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolUICustomizationAttachmentWithName retrieves all cognito.UserPoolUICustomizationAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolUICustomizationAttachmentWithName(name string) (*cognito.UserPoolUICustomizationAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolUICustomizationAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolUICustomizationAttachment not found", name)
}

// GetAllCognitoUserPoolUserResources retrieves all cognito.UserPoolUser items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolUserResources() map[string]*cognito.UserPoolUser {
	results := map[string]*cognito.UserPoolUser{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolUser:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolUserWithName retrieves all cognito.UserPoolUser items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolUserWithName(name string) (*cognito.UserPoolUser, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolUser:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolUser not found", name)
}

// GetAllCognitoUserPoolUserToGroupAttachmentResources retrieves all cognito.UserPoolUserToGroupAttachment items from an AWS CloudFormation template
func (t *Template) GetAllCognitoUserPoolUserToGroupAttachmentResources() map[string]*cognito.UserPoolUserToGroupAttachment {
	results := map[string]*cognito.UserPoolUserToGroupAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cognito.UserPoolUserToGroupAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetCognitoUserPoolUserToGroupAttachmentWithName retrieves all cognito.UserPoolUserToGroupAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCognitoUserPoolUserToGroupAttachmentWithName(name string) (*cognito.UserPoolUserToGroupAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cognito.UserPoolUserToGroupAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cognito.UserPoolUserToGroupAttachment not found", name)
}

// GetAllConfigAggregationAuthorizationResources retrieves all config.AggregationAuthorization items from an AWS CloudFormation template
func (t *Template) GetAllConfigAggregationAuthorizationResources() map[string]*config.AggregationAuthorization {
	results := map[string]*config.AggregationAuthorization{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.AggregationAuthorization:
			results[name] = resource
		}
	}
	return results
}

// GetConfigAggregationAuthorizationWithName retrieves all config.AggregationAuthorization items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigAggregationAuthorizationWithName(name string) (*config.AggregationAuthorization, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.AggregationAuthorization:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.AggregationAuthorization not found", name)
}

// GetAllConfigConfigRuleResources retrieves all config.ConfigRule items from an AWS CloudFormation template
func (t *Template) GetAllConfigConfigRuleResources() map[string]*config.ConfigRule {
	results := map[string]*config.ConfigRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.ConfigRule:
			results[name] = resource
		}
	}
	return results
}

// GetConfigConfigRuleWithName retrieves all config.ConfigRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigConfigRuleWithName(name string) (*config.ConfigRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.ConfigRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.ConfigRule not found", name)
}

// GetAllConfigConfigurationAggregatorResources retrieves all config.ConfigurationAggregator items from an AWS CloudFormation template
func (t *Template) GetAllConfigConfigurationAggregatorResources() map[string]*config.ConfigurationAggregator {
	results := map[string]*config.ConfigurationAggregator{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.ConfigurationAggregator:
			results[name] = resource
		}
	}
	return results
}

// GetConfigConfigurationAggregatorWithName retrieves all config.ConfigurationAggregator items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigConfigurationAggregatorWithName(name string) (*config.ConfigurationAggregator, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.ConfigurationAggregator:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.ConfigurationAggregator not found", name)
}

// GetAllConfigConfigurationRecorderResources retrieves all config.ConfigurationRecorder items from an AWS CloudFormation template
func (t *Template) GetAllConfigConfigurationRecorderResources() map[string]*config.ConfigurationRecorder {
	results := map[string]*config.ConfigurationRecorder{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.ConfigurationRecorder:
			results[name] = resource
		}
	}
	return results
}

// GetConfigConfigurationRecorderWithName retrieves all config.ConfigurationRecorder items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigConfigurationRecorderWithName(name string) (*config.ConfigurationRecorder, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.ConfigurationRecorder:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.ConfigurationRecorder not found", name)
}

// GetAllConfigConformancePackResources retrieves all config.ConformancePack items from an AWS CloudFormation template
func (t *Template) GetAllConfigConformancePackResources() map[string]*config.ConformancePack {
	results := map[string]*config.ConformancePack{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.ConformancePack:
			results[name] = resource
		}
	}
	return results
}

// GetConfigConformancePackWithName retrieves all config.ConformancePack items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigConformancePackWithName(name string) (*config.ConformancePack, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.ConformancePack:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.ConformancePack not found", name)
}

// GetAllConfigDeliveryChannelResources retrieves all config.DeliveryChannel items from an AWS CloudFormation template
func (t *Template) GetAllConfigDeliveryChannelResources() map[string]*config.DeliveryChannel {
	results := map[string]*config.DeliveryChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.DeliveryChannel:
			results[name] = resource
		}
	}
	return results
}

// GetConfigDeliveryChannelWithName retrieves all config.DeliveryChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigDeliveryChannelWithName(name string) (*config.DeliveryChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.DeliveryChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.DeliveryChannel not found", name)
}

// GetAllConfigOrganizationConfigRuleResources retrieves all config.OrganizationConfigRule items from an AWS CloudFormation template
func (t *Template) GetAllConfigOrganizationConfigRuleResources() map[string]*config.OrganizationConfigRule {
	results := map[string]*config.OrganizationConfigRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.OrganizationConfigRule:
			results[name] = resource
		}
	}
	return results
}

// GetConfigOrganizationConfigRuleWithName retrieves all config.OrganizationConfigRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigOrganizationConfigRuleWithName(name string) (*config.OrganizationConfigRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.OrganizationConfigRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.OrganizationConfigRule not found", name)
}

// GetAllConfigOrganizationConformancePackResources retrieves all config.OrganizationConformancePack items from an AWS CloudFormation template
func (t *Template) GetAllConfigOrganizationConformancePackResources() map[string]*config.OrganizationConformancePack {
	results := map[string]*config.OrganizationConformancePack{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.OrganizationConformancePack:
			results[name] = resource
		}
	}
	return results
}

// GetConfigOrganizationConformancePackWithName retrieves all config.OrganizationConformancePack items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigOrganizationConformancePackWithName(name string) (*config.OrganizationConformancePack, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.OrganizationConformancePack:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.OrganizationConformancePack not found", name)
}

// GetAllConfigRemediationConfigurationResources retrieves all config.RemediationConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllConfigRemediationConfigurationResources() map[string]*config.RemediationConfiguration {
	results := map[string]*config.RemediationConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.RemediationConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetConfigRemediationConfigurationWithName retrieves all config.RemediationConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigRemediationConfigurationWithName(name string) (*config.RemediationConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.RemediationConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.RemediationConfiguration not found", name)
}

// GetAllConfigStoredQueryResources retrieves all config.StoredQuery items from an AWS CloudFormation template
func (t *Template) GetAllConfigStoredQueryResources() map[string]*config.StoredQuery {
	results := map[string]*config.StoredQuery{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *config.StoredQuery:
			results[name] = resource
		}
	}
	return results
}

// GetConfigStoredQueryWithName retrieves all config.StoredQuery items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConfigStoredQueryWithName(name string) (*config.StoredQuery, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *config.StoredQuery:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type config.StoredQuery not found", name)
}

// GetAllConnectContactFlowResources retrieves all connect.ContactFlow items from an AWS CloudFormation template
func (t *Template) GetAllConnectContactFlowResources() map[string]*connect.ContactFlow {
	results := map[string]*connect.ContactFlow{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *connect.ContactFlow:
			results[name] = resource
		}
	}
	return results
}

// GetConnectContactFlowWithName retrieves all connect.ContactFlow items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConnectContactFlowWithName(name string) (*connect.ContactFlow, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *connect.ContactFlow:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type connect.ContactFlow not found", name)
}

// GetAllConnectContactFlowModuleResources retrieves all connect.ContactFlowModule items from an AWS CloudFormation template
func (t *Template) GetAllConnectContactFlowModuleResources() map[string]*connect.ContactFlowModule {
	results := map[string]*connect.ContactFlowModule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *connect.ContactFlowModule:
			results[name] = resource
		}
	}
	return results
}

// GetConnectContactFlowModuleWithName retrieves all connect.ContactFlowModule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConnectContactFlowModuleWithName(name string) (*connect.ContactFlowModule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *connect.ContactFlowModule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type connect.ContactFlowModule not found", name)
}

// GetAllConnectHoursOfOperationResources retrieves all connect.HoursOfOperation items from an AWS CloudFormation template
func (t *Template) GetAllConnectHoursOfOperationResources() map[string]*connect.HoursOfOperation {
	results := map[string]*connect.HoursOfOperation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *connect.HoursOfOperation:
			results[name] = resource
		}
	}
	return results
}

// GetConnectHoursOfOperationWithName retrieves all connect.HoursOfOperation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConnectHoursOfOperationWithName(name string) (*connect.HoursOfOperation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *connect.HoursOfOperation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type connect.HoursOfOperation not found", name)
}

// GetAllConnectQuickConnectResources retrieves all connect.QuickConnect items from an AWS CloudFormation template
func (t *Template) GetAllConnectQuickConnectResources() map[string]*connect.QuickConnect {
	results := map[string]*connect.QuickConnect{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *connect.QuickConnect:
			results[name] = resource
		}
	}
	return results
}

// GetConnectQuickConnectWithName retrieves all connect.QuickConnect items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConnectQuickConnectWithName(name string) (*connect.QuickConnect, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *connect.QuickConnect:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type connect.QuickConnect not found", name)
}

// GetAllConnectUserResources retrieves all connect.User items from an AWS CloudFormation template
func (t *Template) GetAllConnectUserResources() map[string]*connect.User {
	results := map[string]*connect.User{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *connect.User:
			results[name] = resource
		}
	}
	return results
}

// GetConnectUserWithName retrieves all connect.User items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConnectUserWithName(name string) (*connect.User, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *connect.User:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type connect.User not found", name)
}

// GetAllConnectUserHierarchyGroupResources retrieves all connect.UserHierarchyGroup items from an AWS CloudFormation template
func (t *Template) GetAllConnectUserHierarchyGroupResources() map[string]*connect.UserHierarchyGroup {
	results := map[string]*connect.UserHierarchyGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *connect.UserHierarchyGroup:
			results[name] = resource
		}
	}
	return results
}

// GetConnectUserHierarchyGroupWithName retrieves all connect.UserHierarchyGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetConnectUserHierarchyGroupWithName(name string) (*connect.UserHierarchyGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *connect.UserHierarchyGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type connect.UserHierarchyGroup not found", name)
}

// GetAllCustomerProfilesDomainResources retrieves all customerprofiles.Domain items from an AWS CloudFormation template
func (t *Template) GetAllCustomerProfilesDomainResources() map[string]*customerprofiles.Domain {
	results := map[string]*customerprofiles.Domain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *customerprofiles.Domain:
			results[name] = resource
		}
	}
	return results
}

// GetCustomerProfilesDomainWithName retrieves all customerprofiles.Domain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCustomerProfilesDomainWithName(name string) (*customerprofiles.Domain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *customerprofiles.Domain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type customerprofiles.Domain not found", name)
}

// GetAllCustomerProfilesIntegrationResources retrieves all customerprofiles.Integration items from an AWS CloudFormation template
func (t *Template) GetAllCustomerProfilesIntegrationResources() map[string]*customerprofiles.Integration {
	results := map[string]*customerprofiles.Integration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *customerprofiles.Integration:
			results[name] = resource
		}
	}
	return results
}

// GetCustomerProfilesIntegrationWithName retrieves all customerprofiles.Integration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCustomerProfilesIntegrationWithName(name string) (*customerprofiles.Integration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *customerprofiles.Integration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type customerprofiles.Integration not found", name)
}

// GetAllCustomerProfilesObjectTypeResources retrieves all customerprofiles.ObjectType items from an AWS CloudFormation template
func (t *Template) GetAllCustomerProfilesObjectTypeResources() map[string]*customerprofiles.ObjectType {
	results := map[string]*customerprofiles.ObjectType{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *customerprofiles.ObjectType:
			results[name] = resource
		}
	}
	return results
}

// GetCustomerProfilesObjectTypeWithName retrieves all customerprofiles.ObjectType items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCustomerProfilesObjectTypeWithName(name string) (*customerprofiles.ObjectType, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *customerprofiles.ObjectType:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type customerprofiles.ObjectType not found", name)
}

// GetAllDAXClusterResources retrieves all dax.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllDAXClusterResources() map[string]*dax.Cluster {
	results := map[string]*dax.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dax.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetDAXClusterWithName retrieves all dax.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDAXClusterWithName(name string) (*dax.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dax.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dax.Cluster not found", name)
}

// GetAllDAXParameterGroupResources retrieves all dax.ParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllDAXParameterGroupResources() map[string]*dax.ParameterGroup {
	results := map[string]*dax.ParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dax.ParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetDAXParameterGroupWithName retrieves all dax.ParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDAXParameterGroupWithName(name string) (*dax.ParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dax.ParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dax.ParameterGroup not found", name)
}

// GetAllDAXSubnetGroupResources retrieves all dax.SubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllDAXSubnetGroupResources() map[string]*dax.SubnetGroup {
	results := map[string]*dax.SubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dax.SubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetDAXSubnetGroupWithName retrieves all dax.SubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDAXSubnetGroupWithName(name string) (*dax.SubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dax.SubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dax.SubnetGroup not found", name)
}

// GetAllDLMLifecyclePolicyResources retrieves all dlm.LifecyclePolicy items from an AWS CloudFormation template
func (t *Template) GetAllDLMLifecyclePolicyResources() map[string]*dlm.LifecyclePolicy {
	results := map[string]*dlm.LifecyclePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dlm.LifecyclePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetDLMLifecyclePolicyWithName retrieves all dlm.LifecyclePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDLMLifecyclePolicyWithName(name string) (*dlm.LifecyclePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dlm.LifecyclePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dlm.LifecyclePolicy not found", name)
}

// GetAllDMSCertificateResources retrieves all dms.Certificate items from an AWS CloudFormation template
func (t *Template) GetAllDMSCertificateResources() map[string]*dms.Certificate {
	results := map[string]*dms.Certificate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dms.Certificate:
			results[name] = resource
		}
	}
	return results
}

// GetDMSCertificateWithName retrieves all dms.Certificate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDMSCertificateWithName(name string) (*dms.Certificate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dms.Certificate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dms.Certificate not found", name)
}

// GetAllDMSEndpointResources retrieves all dms.Endpoint items from an AWS CloudFormation template
func (t *Template) GetAllDMSEndpointResources() map[string]*dms.Endpoint {
	results := map[string]*dms.Endpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dms.Endpoint:
			results[name] = resource
		}
	}
	return results
}

// GetDMSEndpointWithName retrieves all dms.Endpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDMSEndpointWithName(name string) (*dms.Endpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dms.Endpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dms.Endpoint not found", name)
}

// GetAllDMSEventSubscriptionResources retrieves all dms.EventSubscription items from an AWS CloudFormation template
func (t *Template) GetAllDMSEventSubscriptionResources() map[string]*dms.EventSubscription {
	results := map[string]*dms.EventSubscription{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dms.EventSubscription:
			results[name] = resource
		}
	}
	return results
}

// GetDMSEventSubscriptionWithName retrieves all dms.EventSubscription items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDMSEventSubscriptionWithName(name string) (*dms.EventSubscription, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dms.EventSubscription:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dms.EventSubscription not found", name)
}

// GetAllDMSReplicationInstanceResources retrieves all dms.ReplicationInstance items from an AWS CloudFormation template
func (t *Template) GetAllDMSReplicationInstanceResources() map[string]*dms.ReplicationInstance {
	results := map[string]*dms.ReplicationInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dms.ReplicationInstance:
			results[name] = resource
		}
	}
	return results
}

// GetDMSReplicationInstanceWithName retrieves all dms.ReplicationInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDMSReplicationInstanceWithName(name string) (*dms.ReplicationInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dms.ReplicationInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dms.ReplicationInstance not found", name)
}

// GetAllDMSReplicationSubnetGroupResources retrieves all dms.ReplicationSubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllDMSReplicationSubnetGroupResources() map[string]*dms.ReplicationSubnetGroup {
	results := map[string]*dms.ReplicationSubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dms.ReplicationSubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetDMSReplicationSubnetGroupWithName retrieves all dms.ReplicationSubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDMSReplicationSubnetGroupWithName(name string) (*dms.ReplicationSubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dms.ReplicationSubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dms.ReplicationSubnetGroup not found", name)
}

// GetAllDMSReplicationTaskResources retrieves all dms.ReplicationTask items from an AWS CloudFormation template
func (t *Template) GetAllDMSReplicationTaskResources() map[string]*dms.ReplicationTask {
	results := map[string]*dms.ReplicationTask{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dms.ReplicationTask:
			results[name] = resource
		}
	}
	return results
}

// GetDMSReplicationTaskWithName retrieves all dms.ReplicationTask items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDMSReplicationTaskWithName(name string) (*dms.ReplicationTask, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dms.ReplicationTask:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dms.ReplicationTask not found", name)
}

// GetAllDataBrewDatasetResources retrieves all databrew.Dataset items from an AWS CloudFormation template
func (t *Template) GetAllDataBrewDatasetResources() map[string]*databrew.Dataset {
	results := map[string]*databrew.Dataset{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *databrew.Dataset:
			results[name] = resource
		}
	}
	return results
}

// GetDataBrewDatasetWithName retrieves all databrew.Dataset items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataBrewDatasetWithName(name string) (*databrew.Dataset, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *databrew.Dataset:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type databrew.Dataset not found", name)
}

// GetAllDataBrewJobResources retrieves all databrew.Job items from an AWS CloudFormation template
func (t *Template) GetAllDataBrewJobResources() map[string]*databrew.Job {
	results := map[string]*databrew.Job{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *databrew.Job:
			results[name] = resource
		}
	}
	return results
}

// GetDataBrewJobWithName retrieves all databrew.Job items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataBrewJobWithName(name string) (*databrew.Job, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *databrew.Job:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type databrew.Job not found", name)
}

// GetAllDataBrewProjectResources retrieves all databrew.Project items from an AWS CloudFormation template
func (t *Template) GetAllDataBrewProjectResources() map[string]*databrew.Project {
	results := map[string]*databrew.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *databrew.Project:
			results[name] = resource
		}
	}
	return results
}

// GetDataBrewProjectWithName retrieves all databrew.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataBrewProjectWithName(name string) (*databrew.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *databrew.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type databrew.Project not found", name)
}

// GetAllDataBrewRecipeResources retrieves all databrew.Recipe items from an AWS CloudFormation template
func (t *Template) GetAllDataBrewRecipeResources() map[string]*databrew.Recipe {
	results := map[string]*databrew.Recipe{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *databrew.Recipe:
			results[name] = resource
		}
	}
	return results
}

// GetDataBrewRecipeWithName retrieves all databrew.Recipe items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataBrewRecipeWithName(name string) (*databrew.Recipe, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *databrew.Recipe:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type databrew.Recipe not found", name)
}

// GetAllDataBrewRulesetResources retrieves all databrew.Ruleset items from an AWS CloudFormation template
func (t *Template) GetAllDataBrewRulesetResources() map[string]*databrew.Ruleset {
	results := map[string]*databrew.Ruleset{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *databrew.Ruleset:
			results[name] = resource
		}
	}
	return results
}

// GetDataBrewRulesetWithName retrieves all databrew.Ruleset items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataBrewRulesetWithName(name string) (*databrew.Ruleset, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *databrew.Ruleset:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type databrew.Ruleset not found", name)
}

// GetAllDataBrewScheduleResources retrieves all databrew.Schedule items from an AWS CloudFormation template
func (t *Template) GetAllDataBrewScheduleResources() map[string]*databrew.Schedule {
	results := map[string]*databrew.Schedule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *databrew.Schedule:
			results[name] = resource
		}
	}
	return results
}

// GetDataBrewScheduleWithName retrieves all databrew.Schedule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataBrewScheduleWithName(name string) (*databrew.Schedule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *databrew.Schedule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type databrew.Schedule not found", name)
}

// GetAllDataPipelinePipelineResources retrieves all datapipeline.Pipeline items from an AWS CloudFormation template
func (t *Template) GetAllDataPipelinePipelineResources() map[string]*datapipeline.Pipeline {
	results := map[string]*datapipeline.Pipeline{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datapipeline.Pipeline:
			results[name] = resource
		}
	}
	return results
}

// GetDataPipelinePipelineWithName retrieves all datapipeline.Pipeline items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataPipelinePipelineWithName(name string) (*datapipeline.Pipeline, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datapipeline.Pipeline:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datapipeline.Pipeline not found", name)
}

// GetAllDataSyncAgentResources retrieves all datasync.Agent items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncAgentResources() map[string]*datasync.Agent {
	results := map[string]*datasync.Agent{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.Agent:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncAgentWithName retrieves all datasync.Agent items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncAgentWithName(name string) (*datasync.Agent, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.Agent:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.Agent not found", name)
}

// GetAllDataSyncLocationEFSResources retrieves all datasync.LocationEFS items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncLocationEFSResources() map[string]*datasync.LocationEFS {
	results := map[string]*datasync.LocationEFS{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.LocationEFS:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncLocationEFSWithName retrieves all datasync.LocationEFS items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncLocationEFSWithName(name string) (*datasync.LocationEFS, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.LocationEFS:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.LocationEFS not found", name)
}

// GetAllDataSyncLocationFSxWindowsResources retrieves all datasync.LocationFSxWindows items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncLocationFSxWindowsResources() map[string]*datasync.LocationFSxWindows {
	results := map[string]*datasync.LocationFSxWindows{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.LocationFSxWindows:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncLocationFSxWindowsWithName retrieves all datasync.LocationFSxWindows items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncLocationFSxWindowsWithName(name string) (*datasync.LocationFSxWindows, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.LocationFSxWindows:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.LocationFSxWindows not found", name)
}

// GetAllDataSyncLocationHDFSResources retrieves all datasync.LocationHDFS items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncLocationHDFSResources() map[string]*datasync.LocationHDFS {
	results := map[string]*datasync.LocationHDFS{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.LocationHDFS:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncLocationHDFSWithName retrieves all datasync.LocationHDFS items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncLocationHDFSWithName(name string) (*datasync.LocationHDFS, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.LocationHDFS:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.LocationHDFS not found", name)
}

// GetAllDataSyncLocationNFSResources retrieves all datasync.LocationNFS items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncLocationNFSResources() map[string]*datasync.LocationNFS {
	results := map[string]*datasync.LocationNFS{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.LocationNFS:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncLocationNFSWithName retrieves all datasync.LocationNFS items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncLocationNFSWithName(name string) (*datasync.LocationNFS, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.LocationNFS:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.LocationNFS not found", name)
}

// GetAllDataSyncLocationObjectStorageResources retrieves all datasync.LocationObjectStorage items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncLocationObjectStorageResources() map[string]*datasync.LocationObjectStorage {
	results := map[string]*datasync.LocationObjectStorage{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.LocationObjectStorage:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncLocationObjectStorageWithName retrieves all datasync.LocationObjectStorage items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncLocationObjectStorageWithName(name string) (*datasync.LocationObjectStorage, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.LocationObjectStorage:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.LocationObjectStorage not found", name)
}

// GetAllDataSyncLocationS3Resources retrieves all datasync.LocationS3 items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncLocationS3Resources() map[string]*datasync.LocationS3 {
	results := map[string]*datasync.LocationS3{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.LocationS3:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncLocationS3WithName retrieves all datasync.LocationS3 items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncLocationS3WithName(name string) (*datasync.LocationS3, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.LocationS3:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.LocationS3 not found", name)
}

// GetAllDataSyncLocationSMBResources retrieves all datasync.LocationSMB items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncLocationSMBResources() map[string]*datasync.LocationSMB {
	results := map[string]*datasync.LocationSMB{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.LocationSMB:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncLocationSMBWithName retrieves all datasync.LocationSMB items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncLocationSMBWithName(name string) (*datasync.LocationSMB, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.LocationSMB:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.LocationSMB not found", name)
}

// GetAllDataSyncTaskResources retrieves all datasync.Task items from an AWS CloudFormation template
func (t *Template) GetAllDataSyncTaskResources() map[string]*datasync.Task {
	results := map[string]*datasync.Task{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *datasync.Task:
			results[name] = resource
		}
	}
	return results
}

// GetDataSyncTaskWithName retrieves all datasync.Task items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDataSyncTaskWithName(name string) (*datasync.Task, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *datasync.Task:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type datasync.Task not found", name)
}

// GetAllDetectiveGraphResources retrieves all detective.Graph items from an AWS CloudFormation template
func (t *Template) GetAllDetectiveGraphResources() map[string]*detective.Graph {
	results := map[string]*detective.Graph{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *detective.Graph:
			results[name] = resource
		}
	}
	return results
}

// GetDetectiveGraphWithName retrieves all detective.Graph items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDetectiveGraphWithName(name string) (*detective.Graph, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *detective.Graph:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type detective.Graph not found", name)
}

// GetAllDetectiveMemberInvitationResources retrieves all detective.MemberInvitation items from an AWS CloudFormation template
func (t *Template) GetAllDetectiveMemberInvitationResources() map[string]*detective.MemberInvitation {
	results := map[string]*detective.MemberInvitation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *detective.MemberInvitation:
			results[name] = resource
		}
	}
	return results
}

// GetDetectiveMemberInvitationWithName retrieves all detective.MemberInvitation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDetectiveMemberInvitationWithName(name string) (*detective.MemberInvitation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *detective.MemberInvitation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type detective.MemberInvitation not found", name)
}

// GetAllDevOpsGuruNotificationChannelResources retrieves all devopsguru.NotificationChannel items from an AWS CloudFormation template
func (t *Template) GetAllDevOpsGuruNotificationChannelResources() map[string]*devopsguru.NotificationChannel {
	results := map[string]*devopsguru.NotificationChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *devopsguru.NotificationChannel:
			results[name] = resource
		}
	}
	return results
}

// GetDevOpsGuruNotificationChannelWithName retrieves all devopsguru.NotificationChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDevOpsGuruNotificationChannelWithName(name string) (*devopsguru.NotificationChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *devopsguru.NotificationChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type devopsguru.NotificationChannel not found", name)
}

// GetAllDevOpsGuruResourceCollectionResources retrieves all devopsguru.ResourceCollection items from an AWS CloudFormation template
func (t *Template) GetAllDevOpsGuruResourceCollectionResources() map[string]*devopsguru.ResourceCollection {
	results := map[string]*devopsguru.ResourceCollection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *devopsguru.ResourceCollection:
			results[name] = resource
		}
	}
	return results
}

// GetDevOpsGuruResourceCollectionWithName retrieves all devopsguru.ResourceCollection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDevOpsGuruResourceCollectionWithName(name string) (*devopsguru.ResourceCollection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *devopsguru.ResourceCollection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type devopsguru.ResourceCollection not found", name)
}

// GetAllDirectoryServiceMicrosoftADResources retrieves all directoryservice.MicrosoftAD items from an AWS CloudFormation template
func (t *Template) GetAllDirectoryServiceMicrosoftADResources() map[string]*directoryservice.MicrosoftAD {
	results := map[string]*directoryservice.MicrosoftAD{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *directoryservice.MicrosoftAD:
			results[name] = resource
		}
	}
	return results
}

// GetDirectoryServiceMicrosoftADWithName retrieves all directoryservice.MicrosoftAD items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDirectoryServiceMicrosoftADWithName(name string) (*directoryservice.MicrosoftAD, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *directoryservice.MicrosoftAD:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type directoryservice.MicrosoftAD not found", name)
}

// GetAllDirectoryServiceSimpleADResources retrieves all directoryservice.SimpleAD items from an AWS CloudFormation template
func (t *Template) GetAllDirectoryServiceSimpleADResources() map[string]*directoryservice.SimpleAD {
	results := map[string]*directoryservice.SimpleAD{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *directoryservice.SimpleAD:
			results[name] = resource
		}
	}
	return results
}

// GetDirectoryServiceSimpleADWithName retrieves all directoryservice.SimpleAD items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDirectoryServiceSimpleADWithName(name string) (*directoryservice.SimpleAD, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *directoryservice.SimpleAD:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type directoryservice.SimpleAD not found", name)
}

// GetAllDocDBDBClusterResources retrieves all docdb.DBCluster items from an AWS CloudFormation template
func (t *Template) GetAllDocDBDBClusterResources() map[string]*docdb.DBCluster {
	results := map[string]*docdb.DBCluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *docdb.DBCluster:
			results[name] = resource
		}
	}
	return results
}

// GetDocDBDBClusterWithName retrieves all docdb.DBCluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDocDBDBClusterWithName(name string) (*docdb.DBCluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *docdb.DBCluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type docdb.DBCluster not found", name)
}

// GetAllDocDBDBClusterParameterGroupResources retrieves all docdb.DBClusterParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllDocDBDBClusterParameterGroupResources() map[string]*docdb.DBClusterParameterGroup {
	results := map[string]*docdb.DBClusterParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *docdb.DBClusterParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetDocDBDBClusterParameterGroupWithName retrieves all docdb.DBClusterParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDocDBDBClusterParameterGroupWithName(name string) (*docdb.DBClusterParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *docdb.DBClusterParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type docdb.DBClusterParameterGroup not found", name)
}

// GetAllDocDBDBInstanceResources retrieves all docdb.DBInstance items from an AWS CloudFormation template
func (t *Template) GetAllDocDBDBInstanceResources() map[string]*docdb.DBInstance {
	results := map[string]*docdb.DBInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *docdb.DBInstance:
			results[name] = resource
		}
	}
	return results
}

// GetDocDBDBInstanceWithName retrieves all docdb.DBInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDocDBDBInstanceWithName(name string) (*docdb.DBInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *docdb.DBInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type docdb.DBInstance not found", name)
}

// GetAllDocDBDBSubnetGroupResources retrieves all docdb.DBSubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllDocDBDBSubnetGroupResources() map[string]*docdb.DBSubnetGroup {
	results := map[string]*docdb.DBSubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *docdb.DBSubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetDocDBDBSubnetGroupWithName retrieves all docdb.DBSubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDocDBDBSubnetGroupWithName(name string) (*docdb.DBSubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *docdb.DBSubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type docdb.DBSubnetGroup not found", name)
}

// GetAllDynamoDBGlobalTableResources retrieves all dynamodb.GlobalTable items from an AWS CloudFormation template
func (t *Template) GetAllDynamoDBGlobalTableResources() map[string]*dynamodb.GlobalTable {
	results := map[string]*dynamodb.GlobalTable{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dynamodb.GlobalTable:
			results[name] = resource
		}
	}
	return results
}

// GetDynamoDBGlobalTableWithName retrieves all dynamodb.GlobalTable items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDynamoDBGlobalTableWithName(name string) (*dynamodb.GlobalTable, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dynamodb.GlobalTable:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dynamodb.GlobalTable not found", name)
}

// GetAllDynamoDBTableResources retrieves all dynamodb.Table items from an AWS CloudFormation template
func (t *Template) GetAllDynamoDBTableResources() map[string]*dynamodb.Table {
	results := map[string]*dynamodb.Table{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *dynamodb.Table:
			results[name] = resource
		}
	}
	return results
}

// GetDynamoDBTableWithName retrieves all dynamodb.Table items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetDynamoDBTableWithName(name string) (*dynamodb.Table, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *dynamodb.Table:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type dynamodb.Table not found", name)
}

// GetAllEC2CapacityReservationResources retrieves all ec2.CapacityReservation items from an AWS CloudFormation template
func (t *Template) GetAllEC2CapacityReservationResources() map[string]*ec2.CapacityReservation {
	results := map[string]*ec2.CapacityReservation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.CapacityReservation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2CapacityReservationWithName retrieves all ec2.CapacityReservation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2CapacityReservationWithName(name string) (*ec2.CapacityReservation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.CapacityReservation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.CapacityReservation not found", name)
}

// GetAllEC2CapacityReservationFleetResources retrieves all ec2.CapacityReservationFleet items from an AWS CloudFormation template
func (t *Template) GetAllEC2CapacityReservationFleetResources() map[string]*ec2.CapacityReservationFleet {
	results := map[string]*ec2.CapacityReservationFleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.CapacityReservationFleet:
			results[name] = resource
		}
	}
	return results
}

// GetEC2CapacityReservationFleetWithName retrieves all ec2.CapacityReservationFleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2CapacityReservationFleetWithName(name string) (*ec2.CapacityReservationFleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.CapacityReservationFleet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.CapacityReservationFleet not found", name)
}

// GetAllEC2CarrierGatewayResources retrieves all ec2.CarrierGateway items from an AWS CloudFormation template
func (t *Template) GetAllEC2CarrierGatewayResources() map[string]*ec2.CarrierGateway {
	results := map[string]*ec2.CarrierGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.CarrierGateway:
			results[name] = resource
		}
	}
	return results
}

// GetEC2CarrierGatewayWithName retrieves all ec2.CarrierGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2CarrierGatewayWithName(name string) (*ec2.CarrierGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.CarrierGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.CarrierGateway not found", name)
}

// GetAllEC2ClientVpnAuthorizationRuleResources retrieves all ec2.ClientVpnAuthorizationRule items from an AWS CloudFormation template
func (t *Template) GetAllEC2ClientVpnAuthorizationRuleResources() map[string]*ec2.ClientVpnAuthorizationRule {
	results := map[string]*ec2.ClientVpnAuthorizationRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnAuthorizationRule:
			results[name] = resource
		}
	}
	return results
}

// GetEC2ClientVpnAuthorizationRuleWithName retrieves all ec2.ClientVpnAuthorizationRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2ClientVpnAuthorizationRuleWithName(name string) (*ec2.ClientVpnAuthorizationRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnAuthorizationRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.ClientVpnAuthorizationRule not found", name)
}

// GetAllEC2ClientVpnEndpointResources retrieves all ec2.ClientVpnEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllEC2ClientVpnEndpointResources() map[string]*ec2.ClientVpnEndpoint {
	results := map[string]*ec2.ClientVpnEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetEC2ClientVpnEndpointWithName retrieves all ec2.ClientVpnEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2ClientVpnEndpointWithName(name string) (*ec2.ClientVpnEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.ClientVpnEndpoint not found", name)
}

// GetAllEC2ClientVpnRouteResources retrieves all ec2.ClientVpnRoute items from an AWS CloudFormation template
func (t *Template) GetAllEC2ClientVpnRouteResources() map[string]*ec2.ClientVpnRoute {
	results := map[string]*ec2.ClientVpnRoute{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnRoute:
			results[name] = resource
		}
	}
	return results
}

// GetEC2ClientVpnRouteWithName retrieves all ec2.ClientVpnRoute items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2ClientVpnRouteWithName(name string) (*ec2.ClientVpnRoute, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnRoute:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.ClientVpnRoute not found", name)
}

// GetAllEC2ClientVpnTargetNetworkAssociationResources retrieves all ec2.ClientVpnTargetNetworkAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2ClientVpnTargetNetworkAssociationResources() map[string]*ec2.ClientVpnTargetNetworkAssociation {
	results := map[string]*ec2.ClientVpnTargetNetworkAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnTargetNetworkAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2ClientVpnTargetNetworkAssociationWithName retrieves all ec2.ClientVpnTargetNetworkAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2ClientVpnTargetNetworkAssociationWithName(name string) (*ec2.ClientVpnTargetNetworkAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.ClientVpnTargetNetworkAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.ClientVpnTargetNetworkAssociation not found", name)
}

// GetAllEC2CustomerGatewayResources retrieves all ec2.CustomerGateway items from an AWS CloudFormation template
func (t *Template) GetAllEC2CustomerGatewayResources() map[string]*ec2.CustomerGateway {
	results := map[string]*ec2.CustomerGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.CustomerGateway:
			results[name] = resource
		}
	}
	return results
}

// GetEC2CustomerGatewayWithName retrieves all ec2.CustomerGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2CustomerGatewayWithName(name string) (*ec2.CustomerGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.CustomerGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.CustomerGateway not found", name)
}

// GetAllEC2DHCPOptionsResources retrieves all ec2.DHCPOptions items from an AWS CloudFormation template
func (t *Template) GetAllEC2DHCPOptionsResources() map[string]*ec2.DHCPOptions {
	results := map[string]*ec2.DHCPOptions{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.DHCPOptions:
			results[name] = resource
		}
	}
	return results
}

// GetEC2DHCPOptionsWithName retrieves all ec2.DHCPOptions items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2DHCPOptionsWithName(name string) (*ec2.DHCPOptions, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.DHCPOptions:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.DHCPOptions not found", name)
}

// GetAllEC2EC2FleetResources retrieves all ec2.EC2Fleet items from an AWS CloudFormation template
func (t *Template) GetAllEC2EC2FleetResources() map[string]*ec2.EC2Fleet {
	results := map[string]*ec2.EC2Fleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.EC2Fleet:
			results[name] = resource
		}
	}
	return results
}

// GetEC2EC2FleetWithName retrieves all ec2.EC2Fleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2EC2FleetWithName(name string) (*ec2.EC2Fleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.EC2Fleet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.EC2Fleet not found", name)
}

// GetAllEC2EIPResources retrieves all ec2.EIP items from an AWS CloudFormation template
func (t *Template) GetAllEC2EIPResources() map[string]*ec2.EIP {
	results := map[string]*ec2.EIP{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.EIP:
			results[name] = resource
		}
	}
	return results
}

// GetEC2EIPWithName retrieves all ec2.EIP items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2EIPWithName(name string) (*ec2.EIP, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.EIP:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.EIP not found", name)
}

// GetAllEC2EIPAssociationResources retrieves all ec2.EIPAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2EIPAssociationResources() map[string]*ec2.EIPAssociation {
	results := map[string]*ec2.EIPAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.EIPAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2EIPAssociationWithName retrieves all ec2.EIPAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2EIPAssociationWithName(name string) (*ec2.EIPAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.EIPAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.EIPAssociation not found", name)
}

// GetAllEC2EgressOnlyInternetGatewayResources retrieves all ec2.EgressOnlyInternetGateway items from an AWS CloudFormation template
func (t *Template) GetAllEC2EgressOnlyInternetGatewayResources() map[string]*ec2.EgressOnlyInternetGateway {
	results := map[string]*ec2.EgressOnlyInternetGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.EgressOnlyInternetGateway:
			results[name] = resource
		}
	}
	return results
}

// GetEC2EgressOnlyInternetGatewayWithName retrieves all ec2.EgressOnlyInternetGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2EgressOnlyInternetGatewayWithName(name string) (*ec2.EgressOnlyInternetGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.EgressOnlyInternetGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.EgressOnlyInternetGateway not found", name)
}

// GetAllEC2EnclaveCertificateIamRoleAssociationResources retrieves all ec2.EnclaveCertificateIamRoleAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2EnclaveCertificateIamRoleAssociationResources() map[string]*ec2.EnclaveCertificateIamRoleAssociation {
	results := map[string]*ec2.EnclaveCertificateIamRoleAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.EnclaveCertificateIamRoleAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2EnclaveCertificateIamRoleAssociationWithName retrieves all ec2.EnclaveCertificateIamRoleAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2EnclaveCertificateIamRoleAssociationWithName(name string) (*ec2.EnclaveCertificateIamRoleAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.EnclaveCertificateIamRoleAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.EnclaveCertificateIamRoleAssociation not found", name)
}

// GetAllEC2FlowLogResources retrieves all ec2.FlowLog items from an AWS CloudFormation template
func (t *Template) GetAllEC2FlowLogResources() map[string]*ec2.FlowLog {
	results := map[string]*ec2.FlowLog{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.FlowLog:
			results[name] = resource
		}
	}
	return results
}

// GetEC2FlowLogWithName retrieves all ec2.FlowLog items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2FlowLogWithName(name string) (*ec2.FlowLog, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.FlowLog:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.FlowLog not found", name)
}

// GetAllEC2GatewayRouteTableAssociationResources retrieves all ec2.GatewayRouteTableAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2GatewayRouteTableAssociationResources() map[string]*ec2.GatewayRouteTableAssociation {
	results := map[string]*ec2.GatewayRouteTableAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.GatewayRouteTableAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2GatewayRouteTableAssociationWithName retrieves all ec2.GatewayRouteTableAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2GatewayRouteTableAssociationWithName(name string) (*ec2.GatewayRouteTableAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.GatewayRouteTableAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.GatewayRouteTableAssociation not found", name)
}

// GetAllEC2HostResources retrieves all ec2.Host items from an AWS CloudFormation template
func (t *Template) GetAllEC2HostResources() map[string]*ec2.Host {
	results := map[string]*ec2.Host{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.Host:
			results[name] = resource
		}
	}
	return results
}

// GetEC2HostWithName retrieves all ec2.Host items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2HostWithName(name string) (*ec2.Host, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.Host:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.Host not found", name)
}

// GetAllEC2IPAMResources retrieves all ec2.IPAM items from an AWS CloudFormation template
func (t *Template) GetAllEC2IPAMResources() map[string]*ec2.IPAM {
	results := map[string]*ec2.IPAM{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.IPAM:
			results[name] = resource
		}
	}
	return results
}

// GetEC2IPAMWithName retrieves all ec2.IPAM items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2IPAMWithName(name string) (*ec2.IPAM, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.IPAM:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.IPAM not found", name)
}

// GetAllEC2IPAMAllocationResources retrieves all ec2.IPAMAllocation items from an AWS CloudFormation template
func (t *Template) GetAllEC2IPAMAllocationResources() map[string]*ec2.IPAMAllocation {
	results := map[string]*ec2.IPAMAllocation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.IPAMAllocation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2IPAMAllocationWithName retrieves all ec2.IPAMAllocation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2IPAMAllocationWithName(name string) (*ec2.IPAMAllocation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.IPAMAllocation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.IPAMAllocation not found", name)
}

// GetAllEC2IPAMPoolResources retrieves all ec2.IPAMPool items from an AWS CloudFormation template
func (t *Template) GetAllEC2IPAMPoolResources() map[string]*ec2.IPAMPool {
	results := map[string]*ec2.IPAMPool{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.IPAMPool:
			results[name] = resource
		}
	}
	return results
}

// GetEC2IPAMPoolWithName retrieves all ec2.IPAMPool items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2IPAMPoolWithName(name string) (*ec2.IPAMPool, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.IPAMPool:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.IPAMPool not found", name)
}

// GetAllEC2IPAMScopeResources retrieves all ec2.IPAMScope items from an AWS CloudFormation template
func (t *Template) GetAllEC2IPAMScopeResources() map[string]*ec2.IPAMScope {
	results := map[string]*ec2.IPAMScope{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.IPAMScope:
			results[name] = resource
		}
	}
	return results
}

// GetEC2IPAMScopeWithName retrieves all ec2.IPAMScope items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2IPAMScopeWithName(name string) (*ec2.IPAMScope, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.IPAMScope:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.IPAMScope not found", name)
}

// GetAllEC2InstanceResources retrieves all ec2.Instance items from an AWS CloudFormation template
func (t *Template) GetAllEC2InstanceResources() map[string]*ec2.Instance {
	results := map[string]*ec2.Instance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.Instance:
			results[name] = resource
		}
	}
	return results
}

// GetEC2InstanceWithName retrieves all ec2.Instance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2InstanceWithName(name string) (*ec2.Instance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.Instance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.Instance not found", name)
}

// GetAllEC2InternetGatewayResources retrieves all ec2.InternetGateway items from an AWS CloudFormation template
func (t *Template) GetAllEC2InternetGatewayResources() map[string]*ec2.InternetGateway {
	results := map[string]*ec2.InternetGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.InternetGateway:
			results[name] = resource
		}
	}
	return results
}

// GetEC2InternetGatewayWithName retrieves all ec2.InternetGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2InternetGatewayWithName(name string) (*ec2.InternetGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.InternetGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.InternetGateway not found", name)
}

// GetAllEC2LaunchTemplateResources retrieves all ec2.LaunchTemplate items from an AWS CloudFormation template
func (t *Template) GetAllEC2LaunchTemplateResources() map[string]*ec2.LaunchTemplate {
	results := map[string]*ec2.LaunchTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.LaunchTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetEC2LaunchTemplateWithName retrieves all ec2.LaunchTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2LaunchTemplateWithName(name string) (*ec2.LaunchTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.LaunchTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.LaunchTemplate not found", name)
}

// GetAllEC2LocalGatewayRouteResources retrieves all ec2.LocalGatewayRoute items from an AWS CloudFormation template
func (t *Template) GetAllEC2LocalGatewayRouteResources() map[string]*ec2.LocalGatewayRoute {
	results := map[string]*ec2.LocalGatewayRoute{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRoute:
			results[name] = resource
		}
	}
	return results
}

// GetEC2LocalGatewayRouteWithName retrieves all ec2.LocalGatewayRoute items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2LocalGatewayRouteWithName(name string) (*ec2.LocalGatewayRoute, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRoute:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.LocalGatewayRoute not found", name)
}

// GetAllEC2LocalGatewayRouteTableVPCAssociationResources retrieves all ec2.LocalGatewayRouteTableVPCAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2LocalGatewayRouteTableVPCAssociationResources() map[string]*ec2.LocalGatewayRouteTableVPCAssociation {
	results := map[string]*ec2.LocalGatewayRouteTableVPCAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRouteTableVPCAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2LocalGatewayRouteTableVPCAssociationWithName retrieves all ec2.LocalGatewayRouteTableVPCAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2LocalGatewayRouteTableVPCAssociationWithName(name string) (*ec2.LocalGatewayRouteTableVPCAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRouteTableVPCAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.LocalGatewayRouteTableVPCAssociation not found", name)
}

// GetAllEC2NatGatewayResources retrieves all ec2.NatGateway items from an AWS CloudFormation template
func (t *Template) GetAllEC2NatGatewayResources() map[string]*ec2.NatGateway {
	results := map[string]*ec2.NatGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NatGateway:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NatGatewayWithName retrieves all ec2.NatGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NatGatewayWithName(name string) (*ec2.NatGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NatGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NatGateway not found", name)
}

// GetAllEC2NetworkAclResources retrieves all ec2.NetworkAcl items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkAclResources() map[string]*ec2.NetworkAcl {
	results := map[string]*ec2.NetworkAcl{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkAcl:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkAclWithName retrieves all ec2.NetworkAcl items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkAclWithName(name string) (*ec2.NetworkAcl, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkAcl:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkAcl not found", name)
}

// GetAllEC2NetworkAclEntryResources retrieves all ec2.NetworkAclEntry items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkAclEntryResources() map[string]*ec2.NetworkAclEntry {
	results := map[string]*ec2.NetworkAclEntry{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkAclEntry:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkAclEntryWithName retrieves all ec2.NetworkAclEntry items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkAclEntryWithName(name string) (*ec2.NetworkAclEntry, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkAclEntry:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkAclEntry not found", name)
}

// GetAllEC2NetworkInsightsAnalysisResources retrieves all ec2.NetworkInsightsAnalysis items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkInsightsAnalysisResources() map[string]*ec2.NetworkInsightsAnalysis {
	results := map[string]*ec2.NetworkInsightsAnalysis{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsAnalysis:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkInsightsAnalysisWithName retrieves all ec2.NetworkInsightsAnalysis items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkInsightsAnalysisWithName(name string) (*ec2.NetworkInsightsAnalysis, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsAnalysis:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkInsightsAnalysis not found", name)
}

// GetAllEC2NetworkInsightsPathResources retrieves all ec2.NetworkInsightsPath items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkInsightsPathResources() map[string]*ec2.NetworkInsightsPath {
	results := map[string]*ec2.NetworkInsightsPath{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsPath:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkInsightsPathWithName retrieves all ec2.NetworkInsightsPath items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkInsightsPathWithName(name string) (*ec2.NetworkInsightsPath, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsPath:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkInsightsPath not found", name)
}

// GetAllEC2NetworkInterfaceResources retrieves all ec2.NetworkInterface items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkInterfaceResources() map[string]*ec2.NetworkInterface {
	results := map[string]*ec2.NetworkInterface{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkInterface:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkInterfaceWithName retrieves all ec2.NetworkInterface items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkInterfaceWithName(name string) (*ec2.NetworkInterface, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkInterface:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkInterface not found", name)
}

// GetAllEC2NetworkInterfaceAttachmentResources retrieves all ec2.NetworkInterfaceAttachment items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkInterfaceAttachmentResources() map[string]*ec2.NetworkInterfaceAttachment {
	results := map[string]*ec2.NetworkInterfaceAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkInterfaceAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkInterfaceAttachmentWithName retrieves all ec2.NetworkInterfaceAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkInterfaceAttachmentWithName(name string) (*ec2.NetworkInterfaceAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkInterfaceAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkInterfaceAttachment not found", name)
}

// GetAllEC2NetworkInterfacePermissionResources retrieves all ec2.NetworkInterfacePermission items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkInterfacePermissionResources() map[string]*ec2.NetworkInterfacePermission {
	results := map[string]*ec2.NetworkInterfacePermission{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkInterfacePermission:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkInterfacePermissionWithName retrieves all ec2.NetworkInterfacePermission items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkInterfacePermissionWithName(name string) (*ec2.NetworkInterfacePermission, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkInterfacePermission:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkInterfacePermission not found", name)
}

// GetAllEC2PlacementGroupResources retrieves all ec2.PlacementGroup items from an AWS CloudFormation template
func (t *Template) GetAllEC2PlacementGroupResources() map[string]*ec2.PlacementGroup {
	results := map[string]*ec2.PlacementGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.PlacementGroup:
			results[name] = resource
		}
	}
	return results
}

// GetEC2PlacementGroupWithName retrieves all ec2.PlacementGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2PlacementGroupWithName(name string) (*ec2.PlacementGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.PlacementGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.PlacementGroup not found", name)
}

// GetAllEC2PrefixListResources retrieves all ec2.PrefixList items from an AWS CloudFormation template
func (t *Template) GetAllEC2PrefixListResources() map[string]*ec2.PrefixList {
	results := map[string]*ec2.PrefixList{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.PrefixList:
			results[name] = resource
		}
	}
	return results
}

// GetEC2PrefixListWithName retrieves all ec2.PrefixList items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2PrefixListWithName(name string) (*ec2.PrefixList, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.PrefixList:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.PrefixList not found", name)
}

// GetAllEC2RouteResources retrieves all ec2.Route items from an AWS CloudFormation template
func (t *Template) GetAllEC2RouteResources() map[string]*ec2.Route {
	results := map[string]*ec2.Route{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.Route:
			results[name] = resource
		}
	}
	return results
}

// GetEC2RouteWithName retrieves all ec2.Route items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2RouteWithName(name string) (*ec2.Route, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.Route:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.Route not found", name)
}

// GetAllEC2RouteTableResources retrieves all ec2.RouteTable items from an AWS CloudFormation template
func (t *Template) GetAllEC2RouteTableResources() map[string]*ec2.RouteTable {
	results := map[string]*ec2.RouteTable{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.RouteTable:
			results[name] = resource
		}
	}
	return results
}

// GetEC2RouteTableWithName retrieves all ec2.RouteTable items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2RouteTableWithName(name string) (*ec2.RouteTable, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.RouteTable:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.RouteTable not found", name)
}

// GetAllEC2SecurityGroupResources retrieves all ec2.SecurityGroup items from an AWS CloudFormation template
func (t *Template) GetAllEC2SecurityGroupResources() map[string]*ec2.SecurityGroup {
	results := map[string]*ec2.SecurityGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroup:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SecurityGroupWithName retrieves all ec2.SecurityGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SecurityGroupWithName(name string) (*ec2.SecurityGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SecurityGroup not found", name)
}

// GetAllEC2SecurityGroupEgressResources retrieves all ec2.SecurityGroupEgress items from an AWS CloudFormation template
func (t *Template) GetAllEC2SecurityGroupEgressResources() map[string]*ec2.SecurityGroupEgress {
	results := map[string]*ec2.SecurityGroupEgress{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroupEgress:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SecurityGroupEgressWithName retrieves all ec2.SecurityGroupEgress items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SecurityGroupEgressWithName(name string) (*ec2.SecurityGroupEgress, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroupEgress:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SecurityGroupEgress not found", name)
}

// GetAllEC2SecurityGroupIngressResources retrieves all ec2.SecurityGroupIngress items from an AWS CloudFormation template
func (t *Template) GetAllEC2SecurityGroupIngressResources() map[string]*ec2.SecurityGroupIngress {
	results := map[string]*ec2.SecurityGroupIngress{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroupIngress:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SecurityGroupIngressWithName retrieves all ec2.SecurityGroupIngress items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SecurityGroupIngressWithName(name string) (*ec2.SecurityGroupIngress, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroupIngress:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SecurityGroupIngress not found", name)
}

// GetAllEC2SpotFleetResources retrieves all ec2.SpotFleet items from an AWS CloudFormation template
func (t *Template) GetAllEC2SpotFleetResources() map[string]*ec2.SpotFleet {
	results := map[string]*ec2.SpotFleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SpotFleet:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SpotFleetWithName retrieves all ec2.SpotFleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SpotFleetWithName(name string) (*ec2.SpotFleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SpotFleet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SpotFleet not found", name)
}

// GetAllEC2SubnetResources retrieves all ec2.Subnet items from an AWS CloudFormation template
func (t *Template) GetAllEC2SubnetResources() map[string]*ec2.Subnet {
	results := map[string]*ec2.Subnet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.Subnet:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SubnetWithName retrieves all ec2.Subnet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SubnetWithName(name string) (*ec2.Subnet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.Subnet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.Subnet not found", name)
}

// GetAllEC2SubnetCidrBlockResources retrieves all ec2.SubnetCidrBlock items from an AWS CloudFormation template
func (t *Template) GetAllEC2SubnetCidrBlockResources() map[string]*ec2.SubnetCidrBlock {
	results := map[string]*ec2.SubnetCidrBlock{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SubnetCidrBlock:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SubnetCidrBlockWithName retrieves all ec2.SubnetCidrBlock items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SubnetCidrBlockWithName(name string) (*ec2.SubnetCidrBlock, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SubnetCidrBlock:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SubnetCidrBlock not found", name)
}

// GetAllEC2SubnetNetworkAclAssociationResources retrieves all ec2.SubnetNetworkAclAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2SubnetNetworkAclAssociationResources() map[string]*ec2.SubnetNetworkAclAssociation {
	results := map[string]*ec2.SubnetNetworkAclAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SubnetNetworkAclAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SubnetNetworkAclAssociationWithName retrieves all ec2.SubnetNetworkAclAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SubnetNetworkAclAssociationWithName(name string) (*ec2.SubnetNetworkAclAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SubnetNetworkAclAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SubnetNetworkAclAssociation not found", name)
}

// GetAllEC2SubnetRouteTableAssociationResources retrieves all ec2.SubnetRouteTableAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2SubnetRouteTableAssociationResources() map[string]*ec2.SubnetRouteTableAssociation {
	results := map[string]*ec2.SubnetRouteTableAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SubnetRouteTableAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SubnetRouteTableAssociationWithName retrieves all ec2.SubnetRouteTableAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SubnetRouteTableAssociationWithName(name string) (*ec2.SubnetRouteTableAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SubnetRouteTableAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SubnetRouteTableAssociation not found", name)
}

// GetAllEC2TrafficMirrorFilterResources retrieves all ec2.TrafficMirrorFilter items from an AWS CloudFormation template
func (t *Template) GetAllEC2TrafficMirrorFilterResources() map[string]*ec2.TrafficMirrorFilter {
	results := map[string]*ec2.TrafficMirrorFilter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorFilter:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TrafficMirrorFilterWithName retrieves all ec2.TrafficMirrorFilter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TrafficMirrorFilterWithName(name string) (*ec2.TrafficMirrorFilter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorFilter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TrafficMirrorFilter not found", name)
}

// GetAllEC2TrafficMirrorFilterRuleResources retrieves all ec2.TrafficMirrorFilterRule items from an AWS CloudFormation template
func (t *Template) GetAllEC2TrafficMirrorFilterRuleResources() map[string]*ec2.TrafficMirrorFilterRule {
	results := map[string]*ec2.TrafficMirrorFilterRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorFilterRule:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TrafficMirrorFilterRuleWithName retrieves all ec2.TrafficMirrorFilterRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TrafficMirrorFilterRuleWithName(name string) (*ec2.TrafficMirrorFilterRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorFilterRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TrafficMirrorFilterRule not found", name)
}

// GetAllEC2TrafficMirrorSessionResources retrieves all ec2.TrafficMirrorSession items from an AWS CloudFormation template
func (t *Template) GetAllEC2TrafficMirrorSessionResources() map[string]*ec2.TrafficMirrorSession {
	results := map[string]*ec2.TrafficMirrorSession{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorSession:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TrafficMirrorSessionWithName retrieves all ec2.TrafficMirrorSession items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TrafficMirrorSessionWithName(name string) (*ec2.TrafficMirrorSession, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorSession:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TrafficMirrorSession not found", name)
}

// GetAllEC2TrafficMirrorTargetResources retrieves all ec2.TrafficMirrorTarget items from an AWS CloudFormation template
func (t *Template) GetAllEC2TrafficMirrorTargetResources() map[string]*ec2.TrafficMirrorTarget {
	results := map[string]*ec2.TrafficMirrorTarget{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorTarget:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TrafficMirrorTargetWithName retrieves all ec2.TrafficMirrorTarget items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TrafficMirrorTargetWithName(name string) (*ec2.TrafficMirrorTarget, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TrafficMirrorTarget:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TrafficMirrorTarget not found", name)
}

// GetAllEC2TransitGatewayResources retrieves all ec2.TransitGateway items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayResources() map[string]*ec2.TransitGateway {
	results := map[string]*ec2.TransitGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGateway:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayWithName retrieves all ec2.TransitGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayWithName(name string) (*ec2.TransitGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGateway not found", name)
}

// GetAllEC2TransitGatewayAttachmentResources retrieves all ec2.TransitGatewayAttachment items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayAttachmentResources() map[string]*ec2.TransitGatewayAttachment {
	results := map[string]*ec2.TransitGatewayAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayAttachmentWithName retrieves all ec2.TransitGatewayAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayAttachmentWithName(name string) (*ec2.TransitGatewayAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayAttachment not found", name)
}

// GetAllEC2TransitGatewayConnectResources retrieves all ec2.TransitGatewayConnect items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayConnectResources() map[string]*ec2.TransitGatewayConnect {
	results := map[string]*ec2.TransitGatewayConnect{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayConnect:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayConnectWithName retrieves all ec2.TransitGatewayConnect items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayConnectWithName(name string) (*ec2.TransitGatewayConnect, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayConnect:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayConnect not found", name)
}

// GetAllEC2TransitGatewayMulticastDomainResources retrieves all ec2.TransitGatewayMulticastDomain items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayMulticastDomainResources() map[string]*ec2.TransitGatewayMulticastDomain {
	results := map[string]*ec2.TransitGatewayMulticastDomain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastDomain:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayMulticastDomainWithName retrieves all ec2.TransitGatewayMulticastDomain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayMulticastDomainWithName(name string) (*ec2.TransitGatewayMulticastDomain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastDomain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayMulticastDomain not found", name)
}

// GetAllEC2TransitGatewayMulticastDomainAssociationResources retrieves all ec2.TransitGatewayMulticastDomainAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayMulticastDomainAssociationResources() map[string]*ec2.TransitGatewayMulticastDomainAssociation {
	results := map[string]*ec2.TransitGatewayMulticastDomainAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastDomainAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayMulticastDomainAssociationWithName retrieves all ec2.TransitGatewayMulticastDomainAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayMulticastDomainAssociationWithName(name string) (*ec2.TransitGatewayMulticastDomainAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastDomainAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayMulticastDomainAssociation not found", name)
}

// GetAllEC2TransitGatewayMulticastGroupMemberResources retrieves all ec2.TransitGatewayMulticastGroupMember items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayMulticastGroupMemberResources() map[string]*ec2.TransitGatewayMulticastGroupMember {
	results := map[string]*ec2.TransitGatewayMulticastGroupMember{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastGroupMember:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayMulticastGroupMemberWithName retrieves all ec2.TransitGatewayMulticastGroupMember items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayMulticastGroupMemberWithName(name string) (*ec2.TransitGatewayMulticastGroupMember, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastGroupMember:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayMulticastGroupMember not found", name)
}

// GetAllEC2TransitGatewayMulticastGroupSourceResources retrieves all ec2.TransitGatewayMulticastGroupSource items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayMulticastGroupSourceResources() map[string]*ec2.TransitGatewayMulticastGroupSource {
	results := map[string]*ec2.TransitGatewayMulticastGroupSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastGroupSource:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayMulticastGroupSourceWithName retrieves all ec2.TransitGatewayMulticastGroupSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayMulticastGroupSourceWithName(name string) (*ec2.TransitGatewayMulticastGroupSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayMulticastGroupSource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayMulticastGroupSource not found", name)
}

// GetAllEC2TransitGatewayPeeringAttachmentResources retrieves all ec2.TransitGatewayPeeringAttachment items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayPeeringAttachmentResources() map[string]*ec2.TransitGatewayPeeringAttachment {
	results := map[string]*ec2.TransitGatewayPeeringAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayPeeringAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayPeeringAttachmentWithName retrieves all ec2.TransitGatewayPeeringAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayPeeringAttachmentWithName(name string) (*ec2.TransitGatewayPeeringAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayPeeringAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayPeeringAttachment not found", name)
}

// GetAllEC2TransitGatewayRouteResources retrieves all ec2.TransitGatewayRoute items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayRouteResources() map[string]*ec2.TransitGatewayRoute {
	results := map[string]*ec2.TransitGatewayRoute{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRoute:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayRouteWithName retrieves all ec2.TransitGatewayRoute items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayRouteWithName(name string) (*ec2.TransitGatewayRoute, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRoute:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayRoute not found", name)
}

// GetAllEC2TransitGatewayRouteTableResources retrieves all ec2.TransitGatewayRouteTable items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayRouteTableResources() map[string]*ec2.TransitGatewayRouteTable {
	results := map[string]*ec2.TransitGatewayRouteTable{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRouteTable:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayRouteTableWithName retrieves all ec2.TransitGatewayRouteTable items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayRouteTableWithName(name string) (*ec2.TransitGatewayRouteTable, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRouteTable:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayRouteTable not found", name)
}

// GetAllEC2TransitGatewayRouteTableAssociationResources retrieves all ec2.TransitGatewayRouteTableAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayRouteTableAssociationResources() map[string]*ec2.TransitGatewayRouteTableAssociation {
	results := map[string]*ec2.TransitGatewayRouteTableAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRouteTableAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayRouteTableAssociationWithName retrieves all ec2.TransitGatewayRouteTableAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayRouteTableAssociationWithName(name string) (*ec2.TransitGatewayRouteTableAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRouteTableAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayRouteTableAssociation not found", name)
}

// GetAllEC2TransitGatewayRouteTablePropagationResources retrieves all ec2.TransitGatewayRouteTablePropagation items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayRouteTablePropagationResources() map[string]*ec2.TransitGatewayRouteTablePropagation {
	results := map[string]*ec2.TransitGatewayRouteTablePropagation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRouteTablePropagation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayRouteTablePropagationWithName retrieves all ec2.TransitGatewayRouteTablePropagation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayRouteTablePropagationWithName(name string) (*ec2.TransitGatewayRouteTablePropagation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayRouteTablePropagation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayRouteTablePropagation not found", name)
}

// GetAllEC2TransitGatewayVpcAttachmentResources retrieves all ec2.TransitGatewayVpcAttachment items from an AWS CloudFormation template
func (t *Template) GetAllEC2TransitGatewayVpcAttachmentResources() map[string]*ec2.TransitGatewayVpcAttachment {
	results := map[string]*ec2.TransitGatewayVpcAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayVpcAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetEC2TransitGatewayVpcAttachmentWithName retrieves all ec2.TransitGatewayVpcAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2TransitGatewayVpcAttachmentWithName(name string) (*ec2.TransitGatewayVpcAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.TransitGatewayVpcAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.TransitGatewayVpcAttachment not found", name)
}

// GetAllEC2VPCResources retrieves all ec2.VPC items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCResources() map[string]*ec2.VPC {
	results := map[string]*ec2.VPC{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPC:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCWithName retrieves all ec2.VPC items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCWithName(name string) (*ec2.VPC, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPC:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPC not found", name)
}

// GetAllEC2VPCCidrBlockResources retrieves all ec2.VPCCidrBlock items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCCidrBlockResources() map[string]*ec2.VPCCidrBlock {
	results := map[string]*ec2.VPCCidrBlock{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCCidrBlock:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCCidrBlockWithName retrieves all ec2.VPCCidrBlock items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCCidrBlockWithName(name string) (*ec2.VPCCidrBlock, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCCidrBlock:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCCidrBlock not found", name)
}

// GetAllEC2VPCDHCPOptionsAssociationResources retrieves all ec2.VPCDHCPOptionsAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCDHCPOptionsAssociationResources() map[string]*ec2.VPCDHCPOptionsAssociation {
	results := map[string]*ec2.VPCDHCPOptionsAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCDHCPOptionsAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCDHCPOptionsAssociationWithName retrieves all ec2.VPCDHCPOptionsAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCDHCPOptionsAssociationWithName(name string) (*ec2.VPCDHCPOptionsAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCDHCPOptionsAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCDHCPOptionsAssociation not found", name)
}

// GetAllEC2VPCEndpointResources retrieves all ec2.VPCEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCEndpointResources() map[string]*ec2.VPCEndpoint {
	results := map[string]*ec2.VPCEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCEndpointWithName retrieves all ec2.VPCEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCEndpointWithName(name string) (*ec2.VPCEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCEndpoint not found", name)
}

// GetAllEC2VPCEndpointConnectionNotificationResources retrieves all ec2.VPCEndpointConnectionNotification items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCEndpointConnectionNotificationResources() map[string]*ec2.VPCEndpointConnectionNotification {
	results := map[string]*ec2.VPCEndpointConnectionNotification{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpointConnectionNotification:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCEndpointConnectionNotificationWithName retrieves all ec2.VPCEndpointConnectionNotification items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCEndpointConnectionNotificationWithName(name string) (*ec2.VPCEndpointConnectionNotification, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpointConnectionNotification:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCEndpointConnectionNotification not found", name)
}

// GetAllEC2VPCEndpointServiceResources retrieves all ec2.VPCEndpointService items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCEndpointServiceResources() map[string]*ec2.VPCEndpointService {
	results := map[string]*ec2.VPCEndpointService{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpointService:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCEndpointServiceWithName retrieves all ec2.VPCEndpointService items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCEndpointServiceWithName(name string) (*ec2.VPCEndpointService, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpointService:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCEndpointService not found", name)
}

// GetAllEC2VPCEndpointServicePermissionsResources retrieves all ec2.VPCEndpointServicePermissions items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCEndpointServicePermissionsResources() map[string]*ec2.VPCEndpointServicePermissions {
	results := map[string]*ec2.VPCEndpointServicePermissions{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpointServicePermissions:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCEndpointServicePermissionsWithName retrieves all ec2.VPCEndpointServicePermissions items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCEndpointServicePermissionsWithName(name string) (*ec2.VPCEndpointServicePermissions, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCEndpointServicePermissions:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCEndpointServicePermissions not found", name)
}

// GetAllEC2VPCGatewayAttachmentResources retrieves all ec2.VPCGatewayAttachment items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCGatewayAttachmentResources() map[string]*ec2.VPCGatewayAttachment {
	results := map[string]*ec2.VPCGatewayAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCGatewayAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCGatewayAttachmentWithName retrieves all ec2.VPCGatewayAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCGatewayAttachmentWithName(name string) (*ec2.VPCGatewayAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCGatewayAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCGatewayAttachment not found", name)
}

// GetAllEC2VPCPeeringConnectionResources retrieves all ec2.VPCPeeringConnection items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCPeeringConnectionResources() map[string]*ec2.VPCPeeringConnection {
	results := map[string]*ec2.VPCPeeringConnection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCPeeringConnection:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCPeeringConnectionWithName retrieves all ec2.VPCPeeringConnection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCPeeringConnectionWithName(name string) (*ec2.VPCPeeringConnection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCPeeringConnection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCPeeringConnection not found", name)
}

// GetAllEC2VPNConnectionResources retrieves all ec2.VPNConnection items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPNConnectionResources() map[string]*ec2.VPNConnection {
	results := map[string]*ec2.VPNConnection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPNConnection:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPNConnectionWithName retrieves all ec2.VPNConnection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPNConnectionWithName(name string) (*ec2.VPNConnection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPNConnection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPNConnection not found", name)
}

// GetAllEC2VPNConnectionRouteResources retrieves all ec2.VPNConnectionRoute items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPNConnectionRouteResources() map[string]*ec2.VPNConnectionRoute {
	results := map[string]*ec2.VPNConnectionRoute{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPNConnectionRoute:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPNConnectionRouteWithName retrieves all ec2.VPNConnectionRoute items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPNConnectionRouteWithName(name string) (*ec2.VPNConnectionRoute, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPNConnectionRoute:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPNConnectionRoute not found", name)
}

// GetAllEC2VPNGatewayResources retrieves all ec2.VPNGateway items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPNGatewayResources() map[string]*ec2.VPNGateway {
	results := map[string]*ec2.VPNGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPNGateway:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPNGatewayWithName retrieves all ec2.VPNGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPNGatewayWithName(name string) (*ec2.VPNGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPNGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPNGateway not found", name)
}

// GetAllEC2VPNGatewayRoutePropagationResources retrieves all ec2.VPNGatewayRoutePropagation items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPNGatewayRoutePropagationResources() map[string]*ec2.VPNGatewayRoutePropagation {
	results := map[string]*ec2.VPNGatewayRoutePropagation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPNGatewayRoutePropagation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPNGatewayRoutePropagationWithName retrieves all ec2.VPNGatewayRoutePropagation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPNGatewayRoutePropagationWithName(name string) (*ec2.VPNGatewayRoutePropagation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPNGatewayRoutePropagation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPNGatewayRoutePropagation not found", name)
}

// GetAllEC2VolumeResources retrieves all ec2.Volume items from an AWS CloudFormation template
func (t *Template) GetAllEC2VolumeResources() map[string]*ec2.Volume {
	results := map[string]*ec2.Volume{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.Volume:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VolumeWithName retrieves all ec2.Volume items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VolumeWithName(name string) (*ec2.Volume, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.Volume:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.Volume not found", name)
}

// GetAllEC2VolumeAttachmentResources retrieves all ec2.VolumeAttachment items from an AWS CloudFormation template
func (t *Template) GetAllEC2VolumeAttachmentResources() map[string]*ec2.VolumeAttachment {
	results := map[string]*ec2.VolumeAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VolumeAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VolumeAttachmentWithName retrieves all ec2.VolumeAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VolumeAttachmentWithName(name string) (*ec2.VolumeAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VolumeAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VolumeAttachment not found", name)
}

// GetAllECRPublicRepositoryResources retrieves all ecr.PublicRepository items from an AWS CloudFormation template
func (t *Template) GetAllECRPublicRepositoryResources() map[string]*ecr.PublicRepository {
	results := map[string]*ecr.PublicRepository{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecr.PublicRepository:
			results[name] = resource
		}
	}
	return results
}

// GetECRPublicRepositoryWithName retrieves all ecr.PublicRepository items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECRPublicRepositoryWithName(name string) (*ecr.PublicRepository, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecr.PublicRepository:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecr.PublicRepository not found", name)
}

// GetAllECRRegistryPolicyResources retrieves all ecr.RegistryPolicy items from an AWS CloudFormation template
func (t *Template) GetAllECRRegistryPolicyResources() map[string]*ecr.RegistryPolicy {
	results := map[string]*ecr.RegistryPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecr.RegistryPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetECRRegistryPolicyWithName retrieves all ecr.RegistryPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECRRegistryPolicyWithName(name string) (*ecr.RegistryPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecr.RegistryPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecr.RegistryPolicy not found", name)
}

// GetAllECRReplicationConfigurationResources retrieves all ecr.ReplicationConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllECRReplicationConfigurationResources() map[string]*ecr.ReplicationConfiguration {
	results := map[string]*ecr.ReplicationConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecr.ReplicationConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetECRReplicationConfigurationWithName retrieves all ecr.ReplicationConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECRReplicationConfigurationWithName(name string) (*ecr.ReplicationConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecr.ReplicationConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecr.ReplicationConfiguration not found", name)
}

// GetAllECRRepositoryResources retrieves all ecr.Repository items from an AWS CloudFormation template
func (t *Template) GetAllECRRepositoryResources() map[string]*ecr.Repository {
	results := map[string]*ecr.Repository{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecr.Repository:
			results[name] = resource
		}
	}
	return results
}

// GetECRRepositoryWithName retrieves all ecr.Repository items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECRRepositoryWithName(name string) (*ecr.Repository, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecr.Repository:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecr.Repository not found", name)
}

// GetAllECSCapacityProviderResources retrieves all ecs.CapacityProvider items from an AWS CloudFormation template
func (t *Template) GetAllECSCapacityProviderResources() map[string]*ecs.CapacityProvider {
	results := map[string]*ecs.CapacityProvider{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecs.CapacityProvider:
			results[name] = resource
		}
	}
	return results
}

// GetECSCapacityProviderWithName retrieves all ecs.CapacityProvider items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECSCapacityProviderWithName(name string) (*ecs.CapacityProvider, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecs.CapacityProvider:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecs.CapacityProvider not found", name)
}

// GetAllECSClusterResources retrieves all ecs.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllECSClusterResources() map[string]*ecs.Cluster {
	results := map[string]*ecs.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecs.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetECSClusterWithName retrieves all ecs.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECSClusterWithName(name string) (*ecs.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecs.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecs.Cluster not found", name)
}

// GetAllECSClusterCapacityProviderAssociationsResources retrieves all ecs.ClusterCapacityProviderAssociations items from an AWS CloudFormation template
func (t *Template) GetAllECSClusterCapacityProviderAssociationsResources() map[string]*ecs.ClusterCapacityProviderAssociations {
	results := map[string]*ecs.ClusterCapacityProviderAssociations{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecs.ClusterCapacityProviderAssociations:
			results[name] = resource
		}
	}
	return results
}

// GetECSClusterCapacityProviderAssociationsWithName retrieves all ecs.ClusterCapacityProviderAssociations items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECSClusterCapacityProviderAssociationsWithName(name string) (*ecs.ClusterCapacityProviderAssociations, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecs.ClusterCapacityProviderAssociations:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecs.ClusterCapacityProviderAssociations not found", name)
}

// GetAllECSPrimaryTaskSetResources retrieves all ecs.PrimaryTaskSet items from an AWS CloudFormation template
func (t *Template) GetAllECSPrimaryTaskSetResources() map[string]*ecs.PrimaryTaskSet {
	results := map[string]*ecs.PrimaryTaskSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecs.PrimaryTaskSet:
			results[name] = resource
		}
	}
	return results
}

// GetECSPrimaryTaskSetWithName retrieves all ecs.PrimaryTaskSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECSPrimaryTaskSetWithName(name string) (*ecs.PrimaryTaskSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecs.PrimaryTaskSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecs.PrimaryTaskSet not found", name)
}

// GetAllECSServiceResources retrieves all ecs.Service items from an AWS CloudFormation template
func (t *Template) GetAllECSServiceResources() map[string]*ecs.Service {
	results := map[string]*ecs.Service{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecs.Service:
			results[name] = resource
		}
	}
	return results
}

// GetECSServiceWithName retrieves all ecs.Service items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECSServiceWithName(name string) (*ecs.Service, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecs.Service:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecs.Service not found", name)
}

// GetAllECSTaskDefinitionResources retrieves all ecs.TaskDefinition items from an AWS CloudFormation template
func (t *Template) GetAllECSTaskDefinitionResources() map[string]*ecs.TaskDefinition {
	results := map[string]*ecs.TaskDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecs.TaskDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetECSTaskDefinitionWithName retrieves all ecs.TaskDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECSTaskDefinitionWithName(name string) (*ecs.TaskDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecs.TaskDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecs.TaskDefinition not found", name)
}

// GetAllECSTaskSetResources retrieves all ecs.TaskSet items from an AWS CloudFormation template
func (t *Template) GetAllECSTaskSetResources() map[string]*ecs.TaskSet {
	results := map[string]*ecs.TaskSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecs.TaskSet:
			results[name] = resource
		}
	}
	return results
}

// GetECSTaskSetWithName retrieves all ecs.TaskSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECSTaskSetWithName(name string) (*ecs.TaskSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecs.TaskSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecs.TaskSet not found", name)
}

// GetAllEFSAccessPointResources retrieves all efs.AccessPoint items from an AWS CloudFormation template
func (t *Template) GetAllEFSAccessPointResources() map[string]*efs.AccessPoint {
	results := map[string]*efs.AccessPoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *efs.AccessPoint:
			results[name] = resource
		}
	}
	return results
}

// GetEFSAccessPointWithName retrieves all efs.AccessPoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEFSAccessPointWithName(name string) (*efs.AccessPoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *efs.AccessPoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type efs.AccessPoint not found", name)
}

// GetAllEFSFileSystemResources retrieves all efs.FileSystem items from an AWS CloudFormation template
func (t *Template) GetAllEFSFileSystemResources() map[string]*efs.FileSystem {
	results := map[string]*efs.FileSystem{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *efs.FileSystem:
			results[name] = resource
		}
	}
	return results
}

// GetEFSFileSystemWithName retrieves all efs.FileSystem items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEFSFileSystemWithName(name string) (*efs.FileSystem, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *efs.FileSystem:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type efs.FileSystem not found", name)
}

// GetAllEFSMountTargetResources retrieves all efs.MountTarget items from an AWS CloudFormation template
func (t *Template) GetAllEFSMountTargetResources() map[string]*efs.MountTarget {
	results := map[string]*efs.MountTarget{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *efs.MountTarget:
			results[name] = resource
		}
	}
	return results
}

// GetEFSMountTargetWithName retrieves all efs.MountTarget items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEFSMountTargetWithName(name string) (*efs.MountTarget, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *efs.MountTarget:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type efs.MountTarget not found", name)
}

// GetAllEKSAddonResources retrieves all eks.Addon items from an AWS CloudFormation template
func (t *Template) GetAllEKSAddonResources() map[string]*eks.Addon {
	results := map[string]*eks.Addon{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eks.Addon:
			results[name] = resource
		}
	}
	return results
}

// GetEKSAddonWithName retrieves all eks.Addon items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEKSAddonWithName(name string) (*eks.Addon, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eks.Addon:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eks.Addon not found", name)
}

// GetAllEKSClusterResources retrieves all eks.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllEKSClusterResources() map[string]*eks.Cluster {
	results := map[string]*eks.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eks.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetEKSClusterWithName retrieves all eks.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEKSClusterWithName(name string) (*eks.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eks.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eks.Cluster not found", name)
}

// GetAllEKSFargateProfileResources retrieves all eks.FargateProfile items from an AWS CloudFormation template
func (t *Template) GetAllEKSFargateProfileResources() map[string]*eks.FargateProfile {
	results := map[string]*eks.FargateProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eks.FargateProfile:
			results[name] = resource
		}
	}
	return results
}

// GetEKSFargateProfileWithName retrieves all eks.FargateProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEKSFargateProfileWithName(name string) (*eks.FargateProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eks.FargateProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eks.FargateProfile not found", name)
}

// GetAllEKSNodegroupResources retrieves all eks.Nodegroup items from an AWS CloudFormation template
func (t *Template) GetAllEKSNodegroupResources() map[string]*eks.Nodegroup {
	results := map[string]*eks.Nodegroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eks.Nodegroup:
			results[name] = resource
		}
	}
	return results
}

// GetEKSNodegroupWithName retrieves all eks.Nodegroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEKSNodegroupWithName(name string) (*eks.Nodegroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eks.Nodegroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eks.Nodegroup not found", name)
}

// GetAllEMRClusterResources retrieves all emr.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllEMRClusterResources() map[string]*emr.Cluster {
	results := map[string]*emr.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emr.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetEMRClusterWithName retrieves all emr.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRClusterWithName(name string) (*emr.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emr.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emr.Cluster not found", name)
}

// GetAllEMRInstanceFleetConfigResources retrieves all emr.InstanceFleetConfig items from an AWS CloudFormation template
func (t *Template) GetAllEMRInstanceFleetConfigResources() map[string]*emr.InstanceFleetConfig {
	results := map[string]*emr.InstanceFleetConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emr.InstanceFleetConfig:
			results[name] = resource
		}
	}
	return results
}

// GetEMRInstanceFleetConfigWithName retrieves all emr.InstanceFleetConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRInstanceFleetConfigWithName(name string) (*emr.InstanceFleetConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emr.InstanceFleetConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emr.InstanceFleetConfig not found", name)
}

// GetAllEMRInstanceGroupConfigResources retrieves all emr.InstanceGroupConfig items from an AWS CloudFormation template
func (t *Template) GetAllEMRInstanceGroupConfigResources() map[string]*emr.InstanceGroupConfig {
	results := map[string]*emr.InstanceGroupConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emr.InstanceGroupConfig:
			results[name] = resource
		}
	}
	return results
}

// GetEMRInstanceGroupConfigWithName retrieves all emr.InstanceGroupConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRInstanceGroupConfigWithName(name string) (*emr.InstanceGroupConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emr.InstanceGroupConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emr.InstanceGroupConfig not found", name)
}

// GetAllEMRSecurityConfigurationResources retrieves all emr.SecurityConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllEMRSecurityConfigurationResources() map[string]*emr.SecurityConfiguration {
	results := map[string]*emr.SecurityConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emr.SecurityConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetEMRSecurityConfigurationWithName retrieves all emr.SecurityConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRSecurityConfigurationWithName(name string) (*emr.SecurityConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emr.SecurityConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emr.SecurityConfiguration not found", name)
}

// GetAllEMRStepResources retrieves all emr.Step items from an AWS CloudFormation template
func (t *Template) GetAllEMRStepResources() map[string]*emr.Step {
	results := map[string]*emr.Step{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emr.Step:
			results[name] = resource
		}
	}
	return results
}

// GetEMRStepWithName retrieves all emr.Step items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRStepWithName(name string) (*emr.Step, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emr.Step:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emr.Step not found", name)
}

// GetAllEMRStudioResources retrieves all emr.Studio items from an AWS CloudFormation template
func (t *Template) GetAllEMRStudioResources() map[string]*emr.Studio {
	results := map[string]*emr.Studio{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emr.Studio:
			results[name] = resource
		}
	}
	return results
}

// GetEMRStudioWithName retrieves all emr.Studio items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRStudioWithName(name string) (*emr.Studio, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emr.Studio:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emr.Studio not found", name)
}

// GetAllEMRStudioSessionMappingResources retrieves all emr.StudioSessionMapping items from an AWS CloudFormation template
func (t *Template) GetAllEMRStudioSessionMappingResources() map[string]*emr.StudioSessionMapping {
	results := map[string]*emr.StudioSessionMapping{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emr.StudioSessionMapping:
			results[name] = resource
		}
	}
	return results
}

// GetEMRStudioSessionMappingWithName retrieves all emr.StudioSessionMapping items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRStudioSessionMappingWithName(name string) (*emr.StudioSessionMapping, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emr.StudioSessionMapping:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emr.StudioSessionMapping not found", name)
}

// GetAllEMRContainersVirtualClusterResources retrieves all emrcontainers.VirtualCluster items from an AWS CloudFormation template
func (t *Template) GetAllEMRContainersVirtualClusterResources() map[string]*emrcontainers.VirtualCluster {
	results := map[string]*emrcontainers.VirtualCluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *emrcontainers.VirtualCluster:
			results[name] = resource
		}
	}
	return results
}

// GetEMRContainersVirtualClusterWithName retrieves all emrcontainers.VirtualCluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEMRContainersVirtualClusterWithName(name string) (*emrcontainers.VirtualCluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *emrcontainers.VirtualCluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type emrcontainers.VirtualCluster not found", name)
}

// GetAllElastiCacheCacheClusterResources retrieves all elasticache.CacheCluster items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheCacheClusterResources() map[string]*elasticache.CacheCluster {
	results := map[string]*elasticache.CacheCluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.CacheCluster:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheCacheClusterWithName retrieves all elasticache.CacheCluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheCacheClusterWithName(name string) (*elasticache.CacheCluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.CacheCluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.CacheCluster not found", name)
}

// GetAllElastiCacheGlobalReplicationGroupResources retrieves all elasticache.GlobalReplicationGroup items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheGlobalReplicationGroupResources() map[string]*elasticache.GlobalReplicationGroup {
	results := map[string]*elasticache.GlobalReplicationGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.GlobalReplicationGroup:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheGlobalReplicationGroupWithName retrieves all elasticache.GlobalReplicationGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheGlobalReplicationGroupWithName(name string) (*elasticache.GlobalReplicationGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.GlobalReplicationGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.GlobalReplicationGroup not found", name)
}

// GetAllElastiCacheParameterGroupResources retrieves all elasticache.ParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheParameterGroupResources() map[string]*elasticache.ParameterGroup {
	results := map[string]*elasticache.ParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.ParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheParameterGroupWithName retrieves all elasticache.ParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheParameterGroupWithName(name string) (*elasticache.ParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.ParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.ParameterGroup not found", name)
}

// GetAllElastiCacheReplicationGroupResources retrieves all elasticache.ReplicationGroup items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheReplicationGroupResources() map[string]*elasticache.ReplicationGroup {
	results := map[string]*elasticache.ReplicationGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.ReplicationGroup:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheReplicationGroupWithName retrieves all elasticache.ReplicationGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheReplicationGroupWithName(name string) (*elasticache.ReplicationGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.ReplicationGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.ReplicationGroup not found", name)
}

// GetAllElastiCacheSecurityGroupResources retrieves all elasticache.SecurityGroup items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheSecurityGroupResources() map[string]*elasticache.SecurityGroup {
	results := map[string]*elasticache.SecurityGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.SecurityGroup:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheSecurityGroupWithName retrieves all elasticache.SecurityGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheSecurityGroupWithName(name string) (*elasticache.SecurityGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.SecurityGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.SecurityGroup not found", name)
}

// GetAllElastiCacheSecurityGroupIngressResources retrieves all elasticache.SecurityGroupIngress items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheSecurityGroupIngressResources() map[string]*elasticache.SecurityGroupIngress {
	results := map[string]*elasticache.SecurityGroupIngress{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.SecurityGroupIngress:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheSecurityGroupIngressWithName retrieves all elasticache.SecurityGroupIngress items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheSecurityGroupIngressWithName(name string) (*elasticache.SecurityGroupIngress, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.SecurityGroupIngress:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.SecurityGroupIngress not found", name)
}

// GetAllElastiCacheSubnetGroupResources retrieves all elasticache.SubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheSubnetGroupResources() map[string]*elasticache.SubnetGroup {
	results := map[string]*elasticache.SubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.SubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheSubnetGroupWithName retrieves all elasticache.SubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheSubnetGroupWithName(name string) (*elasticache.SubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.SubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.SubnetGroup not found", name)
}

// GetAllElastiCacheUserResources retrieves all elasticache.User items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheUserResources() map[string]*elasticache.User {
	results := map[string]*elasticache.User{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.User:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheUserWithName retrieves all elasticache.User items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheUserWithName(name string) (*elasticache.User, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.User:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.User not found", name)
}

// GetAllElastiCacheUserGroupResources retrieves all elasticache.UserGroup items from an AWS CloudFormation template
func (t *Template) GetAllElastiCacheUserGroupResources() map[string]*elasticache.UserGroup {
	results := map[string]*elasticache.UserGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticache.UserGroup:
			results[name] = resource
		}
	}
	return results
}

// GetElastiCacheUserGroupWithName retrieves all elasticache.UserGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElastiCacheUserGroupWithName(name string) (*elasticache.UserGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticache.UserGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticache.UserGroup not found", name)
}

// GetAllElasticBeanstalkApplicationResources retrieves all elasticbeanstalk.Application items from an AWS CloudFormation template
func (t *Template) GetAllElasticBeanstalkApplicationResources() map[string]*elasticbeanstalk.Application {
	results := map[string]*elasticbeanstalk.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.Application:
			results[name] = resource
		}
	}
	return results
}

// GetElasticBeanstalkApplicationWithName retrieves all elasticbeanstalk.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticBeanstalkApplicationWithName(name string) (*elasticbeanstalk.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticbeanstalk.Application not found", name)
}

// GetAllElasticBeanstalkApplicationVersionResources retrieves all elasticbeanstalk.ApplicationVersion items from an AWS CloudFormation template
func (t *Template) GetAllElasticBeanstalkApplicationVersionResources() map[string]*elasticbeanstalk.ApplicationVersion {
	results := map[string]*elasticbeanstalk.ApplicationVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.ApplicationVersion:
			results[name] = resource
		}
	}
	return results
}

// GetElasticBeanstalkApplicationVersionWithName retrieves all elasticbeanstalk.ApplicationVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticBeanstalkApplicationVersionWithName(name string) (*elasticbeanstalk.ApplicationVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.ApplicationVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticbeanstalk.ApplicationVersion not found", name)
}

// GetAllElasticBeanstalkConfigurationTemplateResources retrieves all elasticbeanstalk.ConfigurationTemplate items from an AWS CloudFormation template
func (t *Template) GetAllElasticBeanstalkConfigurationTemplateResources() map[string]*elasticbeanstalk.ConfigurationTemplate {
	results := map[string]*elasticbeanstalk.ConfigurationTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.ConfigurationTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetElasticBeanstalkConfigurationTemplateWithName retrieves all elasticbeanstalk.ConfigurationTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticBeanstalkConfigurationTemplateWithName(name string) (*elasticbeanstalk.ConfigurationTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.ConfigurationTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticbeanstalk.ConfigurationTemplate not found", name)
}

// GetAllElasticBeanstalkEnvironmentResources retrieves all elasticbeanstalk.Environment items from an AWS CloudFormation template
func (t *Template) GetAllElasticBeanstalkEnvironmentResources() map[string]*elasticbeanstalk.Environment {
	results := map[string]*elasticbeanstalk.Environment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.Environment:
			results[name] = resource
		}
	}
	return results
}

// GetElasticBeanstalkEnvironmentWithName retrieves all elasticbeanstalk.Environment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticBeanstalkEnvironmentWithName(name string) (*elasticbeanstalk.Environment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticbeanstalk.Environment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticbeanstalk.Environment not found", name)
}

// GetAllElasticLoadBalancingLoadBalancerResources retrieves all elasticloadbalancing.LoadBalancer items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingLoadBalancerResources() map[string]*elasticloadbalancing.LoadBalancer {
	results := map[string]*elasticloadbalancing.LoadBalancer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancing.LoadBalancer:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingLoadBalancerWithName retrieves all elasticloadbalancing.LoadBalancer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingLoadBalancerWithName(name string) (*elasticloadbalancing.LoadBalancer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancing.LoadBalancer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancing.LoadBalancer not found", name)
}

// GetAllElasticLoadBalancingV2ListenerResources retrieves all elasticloadbalancingv2.Listener items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingV2ListenerResources() map[string]*elasticloadbalancingv2.Listener {
	results := map[string]*elasticloadbalancingv2.Listener{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.Listener:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingV2ListenerWithName retrieves all elasticloadbalancingv2.Listener items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingV2ListenerWithName(name string) (*elasticloadbalancingv2.Listener, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.Listener:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancingv2.Listener not found", name)
}

// GetAllElasticLoadBalancingV2ListenerCertificateResources retrieves all elasticloadbalancingv2.ListenerCertificate items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingV2ListenerCertificateResources() map[string]*elasticloadbalancingv2.ListenerCertificate {
	results := map[string]*elasticloadbalancingv2.ListenerCertificate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.ListenerCertificate:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingV2ListenerCertificateWithName retrieves all elasticloadbalancingv2.ListenerCertificate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingV2ListenerCertificateWithName(name string) (*elasticloadbalancingv2.ListenerCertificate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.ListenerCertificate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancingv2.ListenerCertificate not found", name)
}

// GetAllElasticLoadBalancingV2ListenerRuleResources retrieves all elasticloadbalancingv2.ListenerRule items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingV2ListenerRuleResources() map[string]*elasticloadbalancingv2.ListenerRule {
	results := map[string]*elasticloadbalancingv2.ListenerRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.ListenerRule:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingV2ListenerRuleWithName retrieves all elasticloadbalancingv2.ListenerRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingV2ListenerRuleWithName(name string) (*elasticloadbalancingv2.ListenerRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.ListenerRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancingv2.ListenerRule not found", name)
}

// GetAllElasticLoadBalancingV2LoadBalancerResources retrieves all elasticloadbalancingv2.LoadBalancer items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingV2LoadBalancerResources() map[string]*elasticloadbalancingv2.LoadBalancer {
	results := map[string]*elasticloadbalancingv2.LoadBalancer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.LoadBalancer:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingV2LoadBalancerWithName retrieves all elasticloadbalancingv2.LoadBalancer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingV2LoadBalancerWithName(name string) (*elasticloadbalancingv2.LoadBalancer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.LoadBalancer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancingv2.LoadBalancer not found", name)
}

// GetAllElasticLoadBalancingV2TargetGroupResources retrieves all elasticloadbalancingv2.TargetGroup items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingV2TargetGroupResources() map[string]*elasticloadbalancingv2.TargetGroup {
	results := map[string]*elasticloadbalancingv2.TargetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.TargetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingV2TargetGroupWithName retrieves all elasticloadbalancingv2.TargetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingV2TargetGroupWithName(name string) (*elasticloadbalancingv2.TargetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.TargetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancingv2.TargetGroup not found", name)
}

// GetAllElasticsearchDomainResources retrieves all elasticsearch.Domain items from an AWS CloudFormation template
func (t *Template) GetAllElasticsearchDomainResources() map[string]*elasticsearch.Domain {
	results := map[string]*elasticsearch.Domain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticsearch.Domain:
			results[name] = resource
		}
	}
	return results
}

// GetElasticsearchDomainWithName retrieves all elasticsearch.Domain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticsearchDomainWithName(name string) (*elasticsearch.Domain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticsearch.Domain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticsearch.Domain not found", name)
}

// GetAllEventSchemasDiscovererResources retrieves all eventschemas.Discoverer items from an AWS CloudFormation template
func (t *Template) GetAllEventSchemasDiscovererResources() map[string]*eventschemas.Discoverer {
	results := map[string]*eventschemas.Discoverer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eventschemas.Discoverer:
			results[name] = resource
		}
	}
	return results
}

// GetEventSchemasDiscovererWithName retrieves all eventschemas.Discoverer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventSchemasDiscovererWithName(name string) (*eventschemas.Discoverer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eventschemas.Discoverer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eventschemas.Discoverer not found", name)
}

// GetAllEventSchemasRegistryResources retrieves all eventschemas.Registry items from an AWS CloudFormation template
func (t *Template) GetAllEventSchemasRegistryResources() map[string]*eventschemas.Registry {
	results := map[string]*eventschemas.Registry{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eventschemas.Registry:
			results[name] = resource
		}
	}
	return results
}

// GetEventSchemasRegistryWithName retrieves all eventschemas.Registry items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventSchemasRegistryWithName(name string) (*eventschemas.Registry, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eventschemas.Registry:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eventschemas.Registry not found", name)
}

// GetAllEventSchemasRegistryPolicyResources retrieves all eventschemas.RegistryPolicy items from an AWS CloudFormation template
func (t *Template) GetAllEventSchemasRegistryPolicyResources() map[string]*eventschemas.RegistryPolicy {
	results := map[string]*eventschemas.RegistryPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eventschemas.RegistryPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetEventSchemasRegistryPolicyWithName retrieves all eventschemas.RegistryPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventSchemasRegistryPolicyWithName(name string) (*eventschemas.RegistryPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eventschemas.RegistryPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eventschemas.RegistryPolicy not found", name)
}

// GetAllEventSchemasSchemaResources retrieves all eventschemas.Schema items from an AWS CloudFormation template
func (t *Template) GetAllEventSchemasSchemaResources() map[string]*eventschemas.Schema {
	results := map[string]*eventschemas.Schema{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eventschemas.Schema:
			results[name] = resource
		}
	}
	return results
}

// GetEventSchemasSchemaWithName retrieves all eventschemas.Schema items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventSchemasSchemaWithName(name string) (*eventschemas.Schema, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eventschemas.Schema:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eventschemas.Schema not found", name)
}

// GetAllEventsApiDestinationResources retrieves all events.ApiDestination items from an AWS CloudFormation template
func (t *Template) GetAllEventsApiDestinationResources() map[string]*events.ApiDestination {
	results := map[string]*events.ApiDestination{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *events.ApiDestination:
			results[name] = resource
		}
	}
	return results
}

// GetEventsApiDestinationWithName retrieves all events.ApiDestination items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventsApiDestinationWithName(name string) (*events.ApiDestination, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *events.ApiDestination:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type events.ApiDestination not found", name)
}

// GetAllEventsArchiveResources retrieves all events.Archive items from an AWS CloudFormation template
func (t *Template) GetAllEventsArchiveResources() map[string]*events.Archive {
	results := map[string]*events.Archive{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *events.Archive:
			results[name] = resource
		}
	}
	return results
}

// GetEventsArchiveWithName retrieves all events.Archive items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventsArchiveWithName(name string) (*events.Archive, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *events.Archive:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type events.Archive not found", name)
}

// GetAllEventsConnectionResources retrieves all events.Connection items from an AWS CloudFormation template
func (t *Template) GetAllEventsConnectionResources() map[string]*events.Connection {
	results := map[string]*events.Connection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *events.Connection:
			results[name] = resource
		}
	}
	return results
}

// GetEventsConnectionWithName retrieves all events.Connection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventsConnectionWithName(name string) (*events.Connection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *events.Connection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type events.Connection not found", name)
}

// GetAllEventsEventBusResources retrieves all events.EventBus items from an AWS CloudFormation template
func (t *Template) GetAllEventsEventBusResources() map[string]*events.EventBus {
	results := map[string]*events.EventBus{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *events.EventBus:
			results[name] = resource
		}
	}
	return results
}

// GetEventsEventBusWithName retrieves all events.EventBus items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventsEventBusWithName(name string) (*events.EventBus, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *events.EventBus:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type events.EventBus not found", name)
}

// GetAllEventsEventBusPolicyResources retrieves all events.EventBusPolicy items from an AWS CloudFormation template
func (t *Template) GetAllEventsEventBusPolicyResources() map[string]*events.EventBusPolicy {
	results := map[string]*events.EventBusPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *events.EventBusPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetEventsEventBusPolicyWithName retrieves all events.EventBusPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventsEventBusPolicyWithName(name string) (*events.EventBusPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *events.EventBusPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type events.EventBusPolicy not found", name)
}

// GetAllEventsRuleResources retrieves all events.Rule items from an AWS CloudFormation template
func (t *Template) GetAllEventsRuleResources() map[string]*events.Rule {
	results := map[string]*events.Rule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *events.Rule:
			results[name] = resource
		}
	}
	return results
}

// GetEventsRuleWithName retrieves all events.Rule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventsRuleWithName(name string) (*events.Rule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *events.Rule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type events.Rule not found", name)
}

// GetAllEvidentlyExperimentResources retrieves all evidently.Experiment items from an AWS CloudFormation template
func (t *Template) GetAllEvidentlyExperimentResources() map[string]*evidently.Experiment {
	results := map[string]*evidently.Experiment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *evidently.Experiment:
			results[name] = resource
		}
	}
	return results
}

// GetEvidentlyExperimentWithName retrieves all evidently.Experiment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEvidentlyExperimentWithName(name string) (*evidently.Experiment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *evidently.Experiment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type evidently.Experiment not found", name)
}

// GetAllEvidentlyFeatureResources retrieves all evidently.Feature items from an AWS CloudFormation template
func (t *Template) GetAllEvidentlyFeatureResources() map[string]*evidently.Feature {
	results := map[string]*evidently.Feature{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *evidently.Feature:
			results[name] = resource
		}
	}
	return results
}

// GetEvidentlyFeatureWithName retrieves all evidently.Feature items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEvidentlyFeatureWithName(name string) (*evidently.Feature, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *evidently.Feature:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type evidently.Feature not found", name)
}

// GetAllEvidentlyLaunchResources retrieves all evidently.Launch items from an AWS CloudFormation template
func (t *Template) GetAllEvidentlyLaunchResources() map[string]*evidently.Launch {
	results := map[string]*evidently.Launch{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *evidently.Launch:
			results[name] = resource
		}
	}
	return results
}

// GetEvidentlyLaunchWithName retrieves all evidently.Launch items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEvidentlyLaunchWithName(name string) (*evidently.Launch, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *evidently.Launch:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type evidently.Launch not found", name)
}

// GetAllEvidentlyProjectResources retrieves all evidently.Project items from an AWS CloudFormation template
func (t *Template) GetAllEvidentlyProjectResources() map[string]*evidently.Project {
	results := map[string]*evidently.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *evidently.Project:
			results[name] = resource
		}
	}
	return results
}

// GetEvidentlyProjectWithName retrieves all evidently.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEvidentlyProjectWithName(name string) (*evidently.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *evidently.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type evidently.Project not found", name)
}

// GetAllFISExperimentTemplateResources retrieves all fis.ExperimentTemplate items from an AWS CloudFormation template
func (t *Template) GetAllFISExperimentTemplateResources() map[string]*fis.ExperimentTemplate {
	results := map[string]*fis.ExperimentTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *fis.ExperimentTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetFISExperimentTemplateWithName retrieves all fis.ExperimentTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFISExperimentTemplateWithName(name string) (*fis.ExperimentTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *fis.ExperimentTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type fis.ExperimentTemplate not found", name)
}

// GetAllFMSNotificationChannelResources retrieves all fms.NotificationChannel items from an AWS CloudFormation template
func (t *Template) GetAllFMSNotificationChannelResources() map[string]*fms.NotificationChannel {
	results := map[string]*fms.NotificationChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *fms.NotificationChannel:
			results[name] = resource
		}
	}
	return results
}

// GetFMSNotificationChannelWithName retrieves all fms.NotificationChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFMSNotificationChannelWithName(name string) (*fms.NotificationChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *fms.NotificationChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type fms.NotificationChannel not found", name)
}

// GetAllFMSPolicyResources retrieves all fms.Policy items from an AWS CloudFormation template
func (t *Template) GetAllFMSPolicyResources() map[string]*fms.Policy {
	results := map[string]*fms.Policy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *fms.Policy:
			results[name] = resource
		}
	}
	return results
}

// GetFMSPolicyWithName retrieves all fms.Policy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFMSPolicyWithName(name string) (*fms.Policy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *fms.Policy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type fms.Policy not found", name)
}

// GetAllFSxFileSystemResources retrieves all fsx.FileSystem items from an AWS CloudFormation template
func (t *Template) GetAllFSxFileSystemResources() map[string]*fsx.FileSystem {
	results := map[string]*fsx.FileSystem{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *fsx.FileSystem:
			results[name] = resource
		}
	}
	return results
}

// GetFSxFileSystemWithName retrieves all fsx.FileSystem items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFSxFileSystemWithName(name string) (*fsx.FileSystem, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *fsx.FileSystem:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type fsx.FileSystem not found", name)
}

// GetAllFinSpaceEnvironmentResources retrieves all finspace.Environment items from an AWS CloudFormation template
func (t *Template) GetAllFinSpaceEnvironmentResources() map[string]*finspace.Environment {
	results := map[string]*finspace.Environment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *finspace.Environment:
			results[name] = resource
		}
	}
	return results
}

// GetFinSpaceEnvironmentWithName retrieves all finspace.Environment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFinSpaceEnvironmentWithName(name string) (*finspace.Environment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *finspace.Environment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type finspace.Environment not found", name)
}

// GetAllFraudDetectorDetectorResources retrieves all frauddetector.Detector items from an AWS CloudFormation template
func (t *Template) GetAllFraudDetectorDetectorResources() map[string]*frauddetector.Detector {
	results := map[string]*frauddetector.Detector{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *frauddetector.Detector:
			results[name] = resource
		}
	}
	return results
}

// GetFraudDetectorDetectorWithName retrieves all frauddetector.Detector items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFraudDetectorDetectorWithName(name string) (*frauddetector.Detector, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *frauddetector.Detector:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type frauddetector.Detector not found", name)
}

// GetAllFraudDetectorEntityTypeResources retrieves all frauddetector.EntityType items from an AWS CloudFormation template
func (t *Template) GetAllFraudDetectorEntityTypeResources() map[string]*frauddetector.EntityType {
	results := map[string]*frauddetector.EntityType{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *frauddetector.EntityType:
			results[name] = resource
		}
	}
	return results
}

// GetFraudDetectorEntityTypeWithName retrieves all frauddetector.EntityType items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFraudDetectorEntityTypeWithName(name string) (*frauddetector.EntityType, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *frauddetector.EntityType:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type frauddetector.EntityType not found", name)
}

// GetAllFraudDetectorEventTypeResources retrieves all frauddetector.EventType items from an AWS CloudFormation template
func (t *Template) GetAllFraudDetectorEventTypeResources() map[string]*frauddetector.EventType {
	results := map[string]*frauddetector.EventType{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *frauddetector.EventType:
			results[name] = resource
		}
	}
	return results
}

// GetFraudDetectorEventTypeWithName retrieves all frauddetector.EventType items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFraudDetectorEventTypeWithName(name string) (*frauddetector.EventType, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *frauddetector.EventType:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type frauddetector.EventType not found", name)
}

// GetAllFraudDetectorLabelResources retrieves all frauddetector.Label items from an AWS CloudFormation template
func (t *Template) GetAllFraudDetectorLabelResources() map[string]*frauddetector.Label {
	results := map[string]*frauddetector.Label{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *frauddetector.Label:
			results[name] = resource
		}
	}
	return results
}

// GetFraudDetectorLabelWithName retrieves all frauddetector.Label items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFraudDetectorLabelWithName(name string) (*frauddetector.Label, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *frauddetector.Label:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type frauddetector.Label not found", name)
}

// GetAllFraudDetectorOutcomeResources retrieves all frauddetector.Outcome items from an AWS CloudFormation template
func (t *Template) GetAllFraudDetectorOutcomeResources() map[string]*frauddetector.Outcome {
	results := map[string]*frauddetector.Outcome{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *frauddetector.Outcome:
			results[name] = resource
		}
	}
	return results
}

// GetFraudDetectorOutcomeWithName retrieves all frauddetector.Outcome items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFraudDetectorOutcomeWithName(name string) (*frauddetector.Outcome, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *frauddetector.Outcome:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type frauddetector.Outcome not found", name)
}

// GetAllFraudDetectorVariableResources retrieves all frauddetector.Variable items from an AWS CloudFormation template
func (t *Template) GetAllFraudDetectorVariableResources() map[string]*frauddetector.Variable {
	results := map[string]*frauddetector.Variable{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *frauddetector.Variable:
			results[name] = resource
		}
	}
	return results
}

// GetFraudDetectorVariableWithName retrieves all frauddetector.Variable items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetFraudDetectorVariableWithName(name string) (*frauddetector.Variable, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *frauddetector.Variable:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type frauddetector.Variable not found", name)
}

// GetAllGameLiftAliasResources retrieves all gamelift.Alias items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftAliasResources() map[string]*gamelift.Alias {
	results := map[string]*gamelift.Alias{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.Alias:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftAliasWithName retrieves all gamelift.Alias items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftAliasWithName(name string) (*gamelift.Alias, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.Alias:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.Alias not found", name)
}

// GetAllGameLiftBuildResources retrieves all gamelift.Build items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftBuildResources() map[string]*gamelift.Build {
	results := map[string]*gamelift.Build{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.Build:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftBuildWithName retrieves all gamelift.Build items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftBuildWithName(name string) (*gamelift.Build, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.Build:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.Build not found", name)
}

// GetAllGameLiftFleetResources retrieves all gamelift.Fleet items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftFleetResources() map[string]*gamelift.Fleet {
	results := map[string]*gamelift.Fleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.Fleet:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftFleetWithName retrieves all gamelift.Fleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftFleetWithName(name string) (*gamelift.Fleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.Fleet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.Fleet not found", name)
}

// GetAllGameLiftGameServerGroupResources retrieves all gamelift.GameServerGroup items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftGameServerGroupResources() map[string]*gamelift.GameServerGroup {
	results := map[string]*gamelift.GameServerGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.GameServerGroup:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftGameServerGroupWithName retrieves all gamelift.GameServerGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftGameServerGroupWithName(name string) (*gamelift.GameServerGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.GameServerGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.GameServerGroup not found", name)
}

// GetAllGameLiftGameSessionQueueResources retrieves all gamelift.GameSessionQueue items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftGameSessionQueueResources() map[string]*gamelift.GameSessionQueue {
	results := map[string]*gamelift.GameSessionQueue{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.GameSessionQueue:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftGameSessionQueueWithName retrieves all gamelift.GameSessionQueue items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftGameSessionQueueWithName(name string) (*gamelift.GameSessionQueue, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.GameSessionQueue:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.GameSessionQueue not found", name)
}

// GetAllGameLiftMatchmakingConfigurationResources retrieves all gamelift.MatchmakingConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftMatchmakingConfigurationResources() map[string]*gamelift.MatchmakingConfiguration {
	results := map[string]*gamelift.MatchmakingConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.MatchmakingConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftMatchmakingConfigurationWithName retrieves all gamelift.MatchmakingConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftMatchmakingConfigurationWithName(name string) (*gamelift.MatchmakingConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.MatchmakingConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.MatchmakingConfiguration not found", name)
}

// GetAllGameLiftMatchmakingRuleSetResources retrieves all gamelift.MatchmakingRuleSet items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftMatchmakingRuleSetResources() map[string]*gamelift.MatchmakingRuleSet {
	results := map[string]*gamelift.MatchmakingRuleSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.MatchmakingRuleSet:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftMatchmakingRuleSetWithName retrieves all gamelift.MatchmakingRuleSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftMatchmakingRuleSetWithName(name string) (*gamelift.MatchmakingRuleSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.MatchmakingRuleSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.MatchmakingRuleSet not found", name)
}

// GetAllGameLiftScriptResources retrieves all gamelift.Script items from an AWS CloudFormation template
func (t *Template) GetAllGameLiftScriptResources() map[string]*gamelift.Script {
	results := map[string]*gamelift.Script{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *gamelift.Script:
			results[name] = resource
		}
	}
	return results
}

// GetGameLiftScriptWithName retrieves all gamelift.Script items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGameLiftScriptWithName(name string) (*gamelift.Script, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *gamelift.Script:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type gamelift.Script not found", name)
}

// GetAllGlobalAcceleratorAcceleratorResources retrieves all globalaccelerator.Accelerator items from an AWS CloudFormation template
func (t *Template) GetAllGlobalAcceleratorAcceleratorResources() map[string]*globalaccelerator.Accelerator {
	results := map[string]*globalaccelerator.Accelerator{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *globalaccelerator.Accelerator:
			results[name] = resource
		}
	}
	return results
}

// GetGlobalAcceleratorAcceleratorWithName retrieves all globalaccelerator.Accelerator items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlobalAcceleratorAcceleratorWithName(name string) (*globalaccelerator.Accelerator, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *globalaccelerator.Accelerator:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type globalaccelerator.Accelerator not found", name)
}

// GetAllGlobalAcceleratorEndpointGroupResources retrieves all globalaccelerator.EndpointGroup items from an AWS CloudFormation template
func (t *Template) GetAllGlobalAcceleratorEndpointGroupResources() map[string]*globalaccelerator.EndpointGroup {
	results := map[string]*globalaccelerator.EndpointGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *globalaccelerator.EndpointGroup:
			results[name] = resource
		}
	}
	return results
}

// GetGlobalAcceleratorEndpointGroupWithName retrieves all globalaccelerator.EndpointGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlobalAcceleratorEndpointGroupWithName(name string) (*globalaccelerator.EndpointGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *globalaccelerator.EndpointGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type globalaccelerator.EndpointGroup not found", name)
}

// GetAllGlobalAcceleratorListenerResources retrieves all globalaccelerator.Listener items from an AWS CloudFormation template
func (t *Template) GetAllGlobalAcceleratorListenerResources() map[string]*globalaccelerator.Listener {
	results := map[string]*globalaccelerator.Listener{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *globalaccelerator.Listener:
			results[name] = resource
		}
	}
	return results
}

// GetGlobalAcceleratorListenerWithName retrieves all globalaccelerator.Listener items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlobalAcceleratorListenerWithName(name string) (*globalaccelerator.Listener, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *globalaccelerator.Listener:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type globalaccelerator.Listener not found", name)
}

// GetAllGlueClassifierResources retrieves all glue.Classifier items from an AWS CloudFormation template
func (t *Template) GetAllGlueClassifierResources() map[string]*glue.Classifier {
	results := map[string]*glue.Classifier{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Classifier:
			results[name] = resource
		}
	}
	return results
}

// GetGlueClassifierWithName retrieves all glue.Classifier items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueClassifierWithName(name string) (*glue.Classifier, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Classifier:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Classifier not found", name)
}

// GetAllGlueConnectionResources retrieves all glue.Connection items from an AWS CloudFormation template
func (t *Template) GetAllGlueConnectionResources() map[string]*glue.Connection {
	results := map[string]*glue.Connection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Connection:
			results[name] = resource
		}
	}
	return results
}

// GetGlueConnectionWithName retrieves all glue.Connection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueConnectionWithName(name string) (*glue.Connection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Connection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Connection not found", name)
}

// GetAllGlueCrawlerResources retrieves all glue.Crawler items from an AWS CloudFormation template
func (t *Template) GetAllGlueCrawlerResources() map[string]*glue.Crawler {
	results := map[string]*glue.Crawler{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Crawler:
			results[name] = resource
		}
	}
	return results
}

// GetGlueCrawlerWithName retrieves all glue.Crawler items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueCrawlerWithName(name string) (*glue.Crawler, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Crawler:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Crawler not found", name)
}

// GetAllGlueDataCatalogEncryptionSettingsResources retrieves all glue.DataCatalogEncryptionSettings items from an AWS CloudFormation template
func (t *Template) GetAllGlueDataCatalogEncryptionSettingsResources() map[string]*glue.DataCatalogEncryptionSettings {
	results := map[string]*glue.DataCatalogEncryptionSettings{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.DataCatalogEncryptionSettings:
			results[name] = resource
		}
	}
	return results
}

// GetGlueDataCatalogEncryptionSettingsWithName retrieves all glue.DataCatalogEncryptionSettings items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueDataCatalogEncryptionSettingsWithName(name string) (*glue.DataCatalogEncryptionSettings, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.DataCatalogEncryptionSettings:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.DataCatalogEncryptionSettings not found", name)
}

// GetAllGlueDatabaseResources retrieves all glue.Database items from an AWS CloudFormation template
func (t *Template) GetAllGlueDatabaseResources() map[string]*glue.Database {
	results := map[string]*glue.Database{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Database:
			results[name] = resource
		}
	}
	return results
}

// GetGlueDatabaseWithName retrieves all glue.Database items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueDatabaseWithName(name string) (*glue.Database, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Database:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Database not found", name)
}

// GetAllGlueDevEndpointResources retrieves all glue.DevEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllGlueDevEndpointResources() map[string]*glue.DevEndpoint {
	results := map[string]*glue.DevEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.DevEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetGlueDevEndpointWithName retrieves all glue.DevEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueDevEndpointWithName(name string) (*glue.DevEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.DevEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.DevEndpoint not found", name)
}

// GetAllGlueJobResources retrieves all glue.Job items from an AWS CloudFormation template
func (t *Template) GetAllGlueJobResources() map[string]*glue.Job {
	results := map[string]*glue.Job{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Job:
			results[name] = resource
		}
	}
	return results
}

// GetGlueJobWithName retrieves all glue.Job items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueJobWithName(name string) (*glue.Job, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Job:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Job not found", name)
}

// GetAllGlueMLTransformResources retrieves all glue.MLTransform items from an AWS CloudFormation template
func (t *Template) GetAllGlueMLTransformResources() map[string]*glue.MLTransform {
	results := map[string]*glue.MLTransform{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.MLTransform:
			results[name] = resource
		}
	}
	return results
}

// GetGlueMLTransformWithName retrieves all glue.MLTransform items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueMLTransformWithName(name string) (*glue.MLTransform, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.MLTransform:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.MLTransform not found", name)
}

// GetAllGluePartitionResources retrieves all glue.Partition items from an AWS CloudFormation template
func (t *Template) GetAllGluePartitionResources() map[string]*glue.Partition {
	results := map[string]*glue.Partition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Partition:
			results[name] = resource
		}
	}
	return results
}

// GetGluePartitionWithName retrieves all glue.Partition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGluePartitionWithName(name string) (*glue.Partition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Partition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Partition not found", name)
}

// GetAllGlueRegistryResources retrieves all glue.Registry items from an AWS CloudFormation template
func (t *Template) GetAllGlueRegistryResources() map[string]*glue.Registry {
	results := map[string]*glue.Registry{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Registry:
			results[name] = resource
		}
	}
	return results
}

// GetGlueRegistryWithName retrieves all glue.Registry items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueRegistryWithName(name string) (*glue.Registry, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Registry:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Registry not found", name)
}

// GetAllGlueSchemaResources retrieves all glue.Schema items from an AWS CloudFormation template
func (t *Template) GetAllGlueSchemaResources() map[string]*glue.Schema {
	results := map[string]*glue.Schema{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Schema:
			results[name] = resource
		}
	}
	return results
}

// GetGlueSchemaWithName retrieves all glue.Schema items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueSchemaWithName(name string) (*glue.Schema, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Schema:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Schema not found", name)
}

// GetAllGlueSchemaVersionResources retrieves all glue.SchemaVersion items from an AWS CloudFormation template
func (t *Template) GetAllGlueSchemaVersionResources() map[string]*glue.SchemaVersion {
	results := map[string]*glue.SchemaVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.SchemaVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGlueSchemaVersionWithName retrieves all glue.SchemaVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueSchemaVersionWithName(name string) (*glue.SchemaVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.SchemaVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.SchemaVersion not found", name)
}

// GetAllGlueSchemaVersionMetadataResources retrieves all glue.SchemaVersionMetadata items from an AWS CloudFormation template
func (t *Template) GetAllGlueSchemaVersionMetadataResources() map[string]*glue.SchemaVersionMetadata {
	results := map[string]*glue.SchemaVersionMetadata{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.SchemaVersionMetadata:
			results[name] = resource
		}
	}
	return results
}

// GetGlueSchemaVersionMetadataWithName retrieves all glue.SchemaVersionMetadata items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueSchemaVersionMetadataWithName(name string) (*glue.SchemaVersionMetadata, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.SchemaVersionMetadata:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.SchemaVersionMetadata not found", name)
}

// GetAllGlueSecurityConfigurationResources retrieves all glue.SecurityConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllGlueSecurityConfigurationResources() map[string]*glue.SecurityConfiguration {
	results := map[string]*glue.SecurityConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.SecurityConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetGlueSecurityConfigurationWithName retrieves all glue.SecurityConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueSecurityConfigurationWithName(name string) (*glue.SecurityConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.SecurityConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.SecurityConfiguration not found", name)
}

// GetAllGlueTableResources retrieves all glue.Table items from an AWS CloudFormation template
func (t *Template) GetAllGlueTableResources() map[string]*glue.Table {
	results := map[string]*glue.Table{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Table:
			results[name] = resource
		}
	}
	return results
}

// GetGlueTableWithName retrieves all glue.Table items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueTableWithName(name string) (*glue.Table, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Table:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Table not found", name)
}

// GetAllGlueTriggerResources retrieves all glue.Trigger items from an AWS CloudFormation template
func (t *Template) GetAllGlueTriggerResources() map[string]*glue.Trigger {
	results := map[string]*glue.Trigger{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Trigger:
			results[name] = resource
		}
	}
	return results
}

// GetGlueTriggerWithName retrieves all glue.Trigger items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueTriggerWithName(name string) (*glue.Trigger, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Trigger:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Trigger not found", name)
}

// GetAllGlueWorkflowResources retrieves all glue.Workflow items from an AWS CloudFormation template
func (t *Template) GetAllGlueWorkflowResources() map[string]*glue.Workflow {
	results := map[string]*glue.Workflow{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *glue.Workflow:
			results[name] = resource
		}
	}
	return results
}

// GetGlueWorkflowWithName retrieves all glue.Workflow items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGlueWorkflowWithName(name string) (*glue.Workflow, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *glue.Workflow:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type glue.Workflow not found", name)
}

// GetAllGreengrassConnectorDefinitionResources retrieves all greengrass.ConnectorDefinition items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassConnectorDefinitionResources() map[string]*greengrass.ConnectorDefinition {
	results := map[string]*greengrass.ConnectorDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.ConnectorDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassConnectorDefinitionWithName retrieves all greengrass.ConnectorDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassConnectorDefinitionWithName(name string) (*greengrass.ConnectorDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.ConnectorDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.ConnectorDefinition not found", name)
}

// GetAllGreengrassConnectorDefinitionVersionResources retrieves all greengrass.ConnectorDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassConnectorDefinitionVersionResources() map[string]*greengrass.ConnectorDefinitionVersion {
	results := map[string]*greengrass.ConnectorDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.ConnectorDefinitionVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassConnectorDefinitionVersionWithName retrieves all greengrass.ConnectorDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassConnectorDefinitionVersionWithName(name string) (*greengrass.ConnectorDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.ConnectorDefinitionVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.ConnectorDefinitionVersion not found", name)
}

// GetAllGreengrassCoreDefinitionResources retrieves all greengrass.CoreDefinition items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassCoreDefinitionResources() map[string]*greengrass.CoreDefinition {
	results := map[string]*greengrass.CoreDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.CoreDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassCoreDefinitionWithName retrieves all greengrass.CoreDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassCoreDefinitionWithName(name string) (*greengrass.CoreDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.CoreDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.CoreDefinition not found", name)
}

// GetAllGreengrassCoreDefinitionVersionResources retrieves all greengrass.CoreDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassCoreDefinitionVersionResources() map[string]*greengrass.CoreDefinitionVersion {
	results := map[string]*greengrass.CoreDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.CoreDefinitionVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassCoreDefinitionVersionWithName retrieves all greengrass.CoreDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassCoreDefinitionVersionWithName(name string) (*greengrass.CoreDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.CoreDefinitionVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.CoreDefinitionVersion not found", name)
}

// GetAllGreengrassDeviceDefinitionResources retrieves all greengrass.DeviceDefinition items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassDeviceDefinitionResources() map[string]*greengrass.DeviceDefinition {
	results := map[string]*greengrass.DeviceDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.DeviceDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassDeviceDefinitionWithName retrieves all greengrass.DeviceDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassDeviceDefinitionWithName(name string) (*greengrass.DeviceDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.DeviceDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.DeviceDefinition not found", name)
}

// GetAllGreengrassDeviceDefinitionVersionResources retrieves all greengrass.DeviceDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassDeviceDefinitionVersionResources() map[string]*greengrass.DeviceDefinitionVersion {
	results := map[string]*greengrass.DeviceDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.DeviceDefinitionVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassDeviceDefinitionVersionWithName retrieves all greengrass.DeviceDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassDeviceDefinitionVersionWithName(name string) (*greengrass.DeviceDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.DeviceDefinitionVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.DeviceDefinitionVersion not found", name)
}

// GetAllGreengrassFunctionDefinitionResources retrieves all greengrass.FunctionDefinition items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassFunctionDefinitionResources() map[string]*greengrass.FunctionDefinition {
	results := map[string]*greengrass.FunctionDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.FunctionDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassFunctionDefinitionWithName retrieves all greengrass.FunctionDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassFunctionDefinitionWithName(name string) (*greengrass.FunctionDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.FunctionDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.FunctionDefinition not found", name)
}

// GetAllGreengrassFunctionDefinitionVersionResources retrieves all greengrass.FunctionDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassFunctionDefinitionVersionResources() map[string]*greengrass.FunctionDefinitionVersion {
	results := map[string]*greengrass.FunctionDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.FunctionDefinitionVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassFunctionDefinitionVersionWithName retrieves all greengrass.FunctionDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassFunctionDefinitionVersionWithName(name string) (*greengrass.FunctionDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.FunctionDefinitionVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.FunctionDefinitionVersion not found", name)
}

// GetAllGreengrassGroupResources retrieves all greengrass.Group items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassGroupResources() map[string]*greengrass.Group {
	results := map[string]*greengrass.Group{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.Group:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassGroupWithName retrieves all greengrass.Group items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassGroupWithName(name string) (*greengrass.Group, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.Group:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.Group not found", name)
}

// GetAllGreengrassGroupVersionResources retrieves all greengrass.GroupVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassGroupVersionResources() map[string]*greengrass.GroupVersion {
	results := map[string]*greengrass.GroupVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.GroupVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassGroupVersionWithName retrieves all greengrass.GroupVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassGroupVersionWithName(name string) (*greengrass.GroupVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.GroupVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.GroupVersion not found", name)
}

// GetAllGreengrassLoggerDefinitionResources retrieves all greengrass.LoggerDefinition items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassLoggerDefinitionResources() map[string]*greengrass.LoggerDefinition {
	results := map[string]*greengrass.LoggerDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.LoggerDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassLoggerDefinitionWithName retrieves all greengrass.LoggerDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassLoggerDefinitionWithName(name string) (*greengrass.LoggerDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.LoggerDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.LoggerDefinition not found", name)
}

// GetAllGreengrassLoggerDefinitionVersionResources retrieves all greengrass.LoggerDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassLoggerDefinitionVersionResources() map[string]*greengrass.LoggerDefinitionVersion {
	results := map[string]*greengrass.LoggerDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.LoggerDefinitionVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassLoggerDefinitionVersionWithName retrieves all greengrass.LoggerDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassLoggerDefinitionVersionWithName(name string) (*greengrass.LoggerDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.LoggerDefinitionVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.LoggerDefinitionVersion not found", name)
}

// GetAllGreengrassResourceDefinitionResources retrieves all greengrass.ResourceDefinition items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassResourceDefinitionResources() map[string]*greengrass.ResourceDefinition {
	results := map[string]*greengrass.ResourceDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.ResourceDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassResourceDefinitionWithName retrieves all greengrass.ResourceDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassResourceDefinitionWithName(name string) (*greengrass.ResourceDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.ResourceDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.ResourceDefinition not found", name)
}

// GetAllGreengrassResourceDefinitionVersionResources retrieves all greengrass.ResourceDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassResourceDefinitionVersionResources() map[string]*greengrass.ResourceDefinitionVersion {
	results := map[string]*greengrass.ResourceDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.ResourceDefinitionVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassResourceDefinitionVersionWithName retrieves all greengrass.ResourceDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassResourceDefinitionVersionWithName(name string) (*greengrass.ResourceDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.ResourceDefinitionVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.ResourceDefinitionVersion not found", name)
}

// GetAllGreengrassSubscriptionDefinitionResources retrieves all greengrass.SubscriptionDefinition items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassSubscriptionDefinitionResources() map[string]*greengrass.SubscriptionDefinition {
	results := map[string]*greengrass.SubscriptionDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.SubscriptionDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassSubscriptionDefinitionWithName retrieves all greengrass.SubscriptionDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassSubscriptionDefinitionWithName(name string) (*greengrass.SubscriptionDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.SubscriptionDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.SubscriptionDefinition not found", name)
}

// GetAllGreengrassSubscriptionDefinitionVersionResources retrieves all greengrass.SubscriptionDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassSubscriptionDefinitionVersionResources() map[string]*greengrass.SubscriptionDefinitionVersion {
	results := map[string]*greengrass.SubscriptionDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrass.SubscriptionDefinitionVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassSubscriptionDefinitionVersionWithName retrieves all greengrass.SubscriptionDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassSubscriptionDefinitionVersionWithName(name string) (*greengrass.SubscriptionDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrass.SubscriptionDefinitionVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrass.SubscriptionDefinitionVersion not found", name)
}

// GetAllGreengrassV2ComponentVersionResources retrieves all greengrassv2.ComponentVersion items from an AWS CloudFormation template
func (t *Template) GetAllGreengrassV2ComponentVersionResources() map[string]*greengrassv2.ComponentVersion {
	results := map[string]*greengrassv2.ComponentVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *greengrassv2.ComponentVersion:
			results[name] = resource
		}
	}
	return results
}

// GetGreengrassV2ComponentVersionWithName retrieves all greengrassv2.ComponentVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGreengrassV2ComponentVersionWithName(name string) (*greengrassv2.ComponentVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *greengrassv2.ComponentVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type greengrassv2.ComponentVersion not found", name)
}

// GetAllGroundStationConfigResources retrieves all groundstation.Config items from an AWS CloudFormation template
func (t *Template) GetAllGroundStationConfigResources() map[string]*groundstation.Config {
	results := map[string]*groundstation.Config{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *groundstation.Config:
			results[name] = resource
		}
	}
	return results
}

// GetGroundStationConfigWithName retrieves all groundstation.Config items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGroundStationConfigWithName(name string) (*groundstation.Config, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *groundstation.Config:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type groundstation.Config not found", name)
}

// GetAllGroundStationDataflowEndpointGroupResources retrieves all groundstation.DataflowEndpointGroup items from an AWS CloudFormation template
func (t *Template) GetAllGroundStationDataflowEndpointGroupResources() map[string]*groundstation.DataflowEndpointGroup {
	results := map[string]*groundstation.DataflowEndpointGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *groundstation.DataflowEndpointGroup:
			results[name] = resource
		}
	}
	return results
}

// GetGroundStationDataflowEndpointGroupWithName retrieves all groundstation.DataflowEndpointGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGroundStationDataflowEndpointGroupWithName(name string) (*groundstation.DataflowEndpointGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *groundstation.DataflowEndpointGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type groundstation.DataflowEndpointGroup not found", name)
}

// GetAllGroundStationMissionProfileResources retrieves all groundstation.MissionProfile items from an AWS CloudFormation template
func (t *Template) GetAllGroundStationMissionProfileResources() map[string]*groundstation.MissionProfile {
	results := map[string]*groundstation.MissionProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *groundstation.MissionProfile:
			results[name] = resource
		}
	}
	return results
}

// GetGroundStationMissionProfileWithName retrieves all groundstation.MissionProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGroundStationMissionProfileWithName(name string) (*groundstation.MissionProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *groundstation.MissionProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type groundstation.MissionProfile not found", name)
}

// GetAllGuardDutyDetectorResources retrieves all guardduty.Detector items from an AWS CloudFormation template
func (t *Template) GetAllGuardDutyDetectorResources() map[string]*guardduty.Detector {
	results := map[string]*guardduty.Detector{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *guardduty.Detector:
			results[name] = resource
		}
	}
	return results
}

// GetGuardDutyDetectorWithName retrieves all guardduty.Detector items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGuardDutyDetectorWithName(name string) (*guardduty.Detector, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *guardduty.Detector:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type guardduty.Detector not found", name)
}

// GetAllGuardDutyFilterResources retrieves all guardduty.Filter items from an AWS CloudFormation template
func (t *Template) GetAllGuardDutyFilterResources() map[string]*guardduty.Filter {
	results := map[string]*guardduty.Filter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *guardduty.Filter:
			results[name] = resource
		}
	}
	return results
}

// GetGuardDutyFilterWithName retrieves all guardduty.Filter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGuardDutyFilterWithName(name string) (*guardduty.Filter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *guardduty.Filter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type guardduty.Filter not found", name)
}

// GetAllGuardDutyIPSetResources retrieves all guardduty.IPSet items from an AWS CloudFormation template
func (t *Template) GetAllGuardDutyIPSetResources() map[string]*guardduty.IPSet {
	results := map[string]*guardduty.IPSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *guardduty.IPSet:
			results[name] = resource
		}
	}
	return results
}

// GetGuardDutyIPSetWithName retrieves all guardduty.IPSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGuardDutyIPSetWithName(name string) (*guardduty.IPSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *guardduty.IPSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type guardduty.IPSet not found", name)
}

// GetAllGuardDutyMasterResources retrieves all guardduty.Master items from an AWS CloudFormation template
func (t *Template) GetAllGuardDutyMasterResources() map[string]*guardduty.Master {
	results := map[string]*guardduty.Master{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *guardduty.Master:
			results[name] = resource
		}
	}
	return results
}

// GetGuardDutyMasterWithName retrieves all guardduty.Master items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGuardDutyMasterWithName(name string) (*guardduty.Master, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *guardduty.Master:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type guardduty.Master not found", name)
}

// GetAllGuardDutyMemberResources retrieves all guardduty.Member items from an AWS CloudFormation template
func (t *Template) GetAllGuardDutyMemberResources() map[string]*guardduty.Member {
	results := map[string]*guardduty.Member{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *guardduty.Member:
			results[name] = resource
		}
	}
	return results
}

// GetGuardDutyMemberWithName retrieves all guardduty.Member items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGuardDutyMemberWithName(name string) (*guardduty.Member, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *guardduty.Member:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type guardduty.Member not found", name)
}

// GetAllGuardDutyThreatIntelSetResources retrieves all guardduty.ThreatIntelSet items from an AWS CloudFormation template
func (t *Template) GetAllGuardDutyThreatIntelSetResources() map[string]*guardduty.ThreatIntelSet {
	results := map[string]*guardduty.ThreatIntelSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *guardduty.ThreatIntelSet:
			results[name] = resource
		}
	}
	return results
}

// GetGuardDutyThreatIntelSetWithName retrieves all guardduty.ThreatIntelSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetGuardDutyThreatIntelSetWithName(name string) (*guardduty.ThreatIntelSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *guardduty.ThreatIntelSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type guardduty.ThreatIntelSet not found", name)
}

// GetAllHealthLakeFHIRDatastoreResources retrieves all healthlake.FHIRDatastore items from an AWS CloudFormation template
func (t *Template) GetAllHealthLakeFHIRDatastoreResources() map[string]*healthlake.FHIRDatastore {
	results := map[string]*healthlake.FHIRDatastore{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *healthlake.FHIRDatastore:
			results[name] = resource
		}
	}
	return results
}

// GetHealthLakeFHIRDatastoreWithName retrieves all healthlake.FHIRDatastore items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetHealthLakeFHIRDatastoreWithName(name string) (*healthlake.FHIRDatastore, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *healthlake.FHIRDatastore:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type healthlake.FHIRDatastore not found", name)
}

// GetAllIAMAccessKeyResources retrieves all iam.AccessKey items from an AWS CloudFormation template
func (t *Template) GetAllIAMAccessKeyResources() map[string]*iam.AccessKey {
	results := map[string]*iam.AccessKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.AccessKey:
			results[name] = resource
		}
	}
	return results
}

// GetIAMAccessKeyWithName retrieves all iam.AccessKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMAccessKeyWithName(name string) (*iam.AccessKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.AccessKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.AccessKey not found", name)
}

// GetAllIAMGroupResources retrieves all iam.Group items from an AWS CloudFormation template
func (t *Template) GetAllIAMGroupResources() map[string]*iam.Group {
	results := map[string]*iam.Group{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.Group:
			results[name] = resource
		}
	}
	return results
}

// GetIAMGroupWithName retrieves all iam.Group items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMGroupWithName(name string) (*iam.Group, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.Group:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.Group not found", name)
}

// GetAllIAMInstanceProfileResources retrieves all iam.InstanceProfile items from an AWS CloudFormation template
func (t *Template) GetAllIAMInstanceProfileResources() map[string]*iam.InstanceProfile {
	results := map[string]*iam.InstanceProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.InstanceProfile:
			results[name] = resource
		}
	}
	return results
}

// GetIAMInstanceProfileWithName retrieves all iam.InstanceProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMInstanceProfileWithName(name string) (*iam.InstanceProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.InstanceProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.InstanceProfile not found", name)
}

// GetAllIAMManagedPolicyResources retrieves all iam.ManagedPolicy items from an AWS CloudFormation template
func (t *Template) GetAllIAMManagedPolicyResources() map[string]*iam.ManagedPolicy {
	results := map[string]*iam.ManagedPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.ManagedPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetIAMManagedPolicyWithName retrieves all iam.ManagedPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMManagedPolicyWithName(name string) (*iam.ManagedPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.ManagedPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.ManagedPolicy not found", name)
}

// GetAllIAMOIDCProviderResources retrieves all iam.OIDCProvider items from an AWS CloudFormation template
func (t *Template) GetAllIAMOIDCProviderResources() map[string]*iam.OIDCProvider {
	results := map[string]*iam.OIDCProvider{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.OIDCProvider:
			results[name] = resource
		}
	}
	return results
}

// GetIAMOIDCProviderWithName retrieves all iam.OIDCProvider items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMOIDCProviderWithName(name string) (*iam.OIDCProvider, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.OIDCProvider:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.OIDCProvider not found", name)
}

// GetAllIAMPolicyResources retrieves all iam.Policy items from an AWS CloudFormation template
func (t *Template) GetAllIAMPolicyResources() map[string]*iam.Policy {
	results := map[string]*iam.Policy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.Policy:
			results[name] = resource
		}
	}
	return results
}

// GetIAMPolicyWithName retrieves all iam.Policy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMPolicyWithName(name string) (*iam.Policy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.Policy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.Policy not found", name)
}

// GetAllIAMRoleResources retrieves all iam.Role items from an AWS CloudFormation template
func (t *Template) GetAllIAMRoleResources() map[string]*iam.Role {
	results := map[string]*iam.Role{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.Role:
			results[name] = resource
		}
	}
	return results
}

// GetIAMRoleWithName retrieves all iam.Role items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMRoleWithName(name string) (*iam.Role, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.Role:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.Role not found", name)
}

// GetAllIAMSAMLProviderResources retrieves all iam.SAMLProvider items from an AWS CloudFormation template
func (t *Template) GetAllIAMSAMLProviderResources() map[string]*iam.SAMLProvider {
	results := map[string]*iam.SAMLProvider{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.SAMLProvider:
			results[name] = resource
		}
	}
	return results
}

// GetIAMSAMLProviderWithName retrieves all iam.SAMLProvider items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMSAMLProviderWithName(name string) (*iam.SAMLProvider, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.SAMLProvider:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.SAMLProvider not found", name)
}

// GetAllIAMServerCertificateResources retrieves all iam.ServerCertificate items from an AWS CloudFormation template
func (t *Template) GetAllIAMServerCertificateResources() map[string]*iam.ServerCertificate {
	results := map[string]*iam.ServerCertificate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.ServerCertificate:
			results[name] = resource
		}
	}
	return results
}

// GetIAMServerCertificateWithName retrieves all iam.ServerCertificate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMServerCertificateWithName(name string) (*iam.ServerCertificate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.ServerCertificate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.ServerCertificate not found", name)
}

// GetAllIAMServiceLinkedRoleResources retrieves all iam.ServiceLinkedRole items from an AWS CloudFormation template
func (t *Template) GetAllIAMServiceLinkedRoleResources() map[string]*iam.ServiceLinkedRole {
	results := map[string]*iam.ServiceLinkedRole{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.ServiceLinkedRole:
			results[name] = resource
		}
	}
	return results
}

// GetIAMServiceLinkedRoleWithName retrieves all iam.ServiceLinkedRole items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMServiceLinkedRoleWithName(name string) (*iam.ServiceLinkedRole, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.ServiceLinkedRole:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.ServiceLinkedRole not found", name)
}

// GetAllIAMUserResources retrieves all iam.User items from an AWS CloudFormation template
func (t *Template) GetAllIAMUserResources() map[string]*iam.User {
	results := map[string]*iam.User{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.User:
			results[name] = resource
		}
	}
	return results
}

// GetIAMUserWithName retrieves all iam.User items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMUserWithName(name string) (*iam.User, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.User:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.User not found", name)
}

// GetAllIAMUserToGroupAdditionResources retrieves all iam.UserToGroupAddition items from an AWS CloudFormation template
func (t *Template) GetAllIAMUserToGroupAdditionResources() map[string]*iam.UserToGroupAddition {
	results := map[string]*iam.UserToGroupAddition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.UserToGroupAddition:
			results[name] = resource
		}
	}
	return results
}

// GetIAMUserToGroupAdditionWithName retrieves all iam.UserToGroupAddition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMUserToGroupAdditionWithName(name string) (*iam.UserToGroupAddition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.UserToGroupAddition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.UserToGroupAddition not found", name)
}

// GetAllIAMVirtualMFADeviceResources retrieves all iam.VirtualMFADevice items from an AWS CloudFormation template
func (t *Template) GetAllIAMVirtualMFADeviceResources() map[string]*iam.VirtualMFADevice {
	results := map[string]*iam.VirtualMFADevice{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.VirtualMFADevice:
			results[name] = resource
		}
	}
	return results
}

// GetIAMVirtualMFADeviceWithName retrieves all iam.VirtualMFADevice items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMVirtualMFADeviceWithName(name string) (*iam.VirtualMFADevice, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.VirtualMFADevice:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.VirtualMFADevice not found", name)
}

// GetAllIVSChannelResources retrieves all ivs.Channel items from an AWS CloudFormation template
func (t *Template) GetAllIVSChannelResources() map[string]*ivs.Channel {
	results := map[string]*ivs.Channel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ivs.Channel:
			results[name] = resource
		}
	}
	return results
}

// GetIVSChannelWithName retrieves all ivs.Channel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIVSChannelWithName(name string) (*ivs.Channel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ivs.Channel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ivs.Channel not found", name)
}

// GetAllIVSPlaybackKeyPairResources retrieves all ivs.PlaybackKeyPair items from an AWS CloudFormation template
func (t *Template) GetAllIVSPlaybackKeyPairResources() map[string]*ivs.PlaybackKeyPair {
	results := map[string]*ivs.PlaybackKeyPair{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ivs.PlaybackKeyPair:
			results[name] = resource
		}
	}
	return results
}

// GetIVSPlaybackKeyPairWithName retrieves all ivs.PlaybackKeyPair items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIVSPlaybackKeyPairWithName(name string) (*ivs.PlaybackKeyPair, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ivs.PlaybackKeyPair:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ivs.PlaybackKeyPair not found", name)
}

// GetAllIVSRecordingConfigurationResources retrieves all ivs.RecordingConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllIVSRecordingConfigurationResources() map[string]*ivs.RecordingConfiguration {
	results := map[string]*ivs.RecordingConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ivs.RecordingConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetIVSRecordingConfigurationWithName retrieves all ivs.RecordingConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIVSRecordingConfigurationWithName(name string) (*ivs.RecordingConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ivs.RecordingConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ivs.RecordingConfiguration not found", name)
}

// GetAllIVSStreamKeyResources retrieves all ivs.StreamKey items from an AWS CloudFormation template
func (t *Template) GetAllIVSStreamKeyResources() map[string]*ivs.StreamKey {
	results := map[string]*ivs.StreamKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ivs.StreamKey:
			results[name] = resource
		}
	}
	return results
}

// GetIVSStreamKeyWithName retrieves all ivs.StreamKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIVSStreamKeyWithName(name string) (*ivs.StreamKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ivs.StreamKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ivs.StreamKey not found", name)
}

// GetAllImageBuilderComponentResources retrieves all imagebuilder.Component items from an AWS CloudFormation template
func (t *Template) GetAllImageBuilderComponentResources() map[string]*imagebuilder.Component {
	results := map[string]*imagebuilder.Component{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *imagebuilder.Component:
			results[name] = resource
		}
	}
	return results
}

// GetImageBuilderComponentWithName retrieves all imagebuilder.Component items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetImageBuilderComponentWithName(name string) (*imagebuilder.Component, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *imagebuilder.Component:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type imagebuilder.Component not found", name)
}

// GetAllImageBuilderContainerRecipeResources retrieves all imagebuilder.ContainerRecipe items from an AWS CloudFormation template
func (t *Template) GetAllImageBuilderContainerRecipeResources() map[string]*imagebuilder.ContainerRecipe {
	results := map[string]*imagebuilder.ContainerRecipe{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *imagebuilder.ContainerRecipe:
			results[name] = resource
		}
	}
	return results
}

// GetImageBuilderContainerRecipeWithName retrieves all imagebuilder.ContainerRecipe items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetImageBuilderContainerRecipeWithName(name string) (*imagebuilder.ContainerRecipe, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *imagebuilder.ContainerRecipe:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type imagebuilder.ContainerRecipe not found", name)
}

// GetAllImageBuilderDistributionConfigurationResources retrieves all imagebuilder.DistributionConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllImageBuilderDistributionConfigurationResources() map[string]*imagebuilder.DistributionConfiguration {
	results := map[string]*imagebuilder.DistributionConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *imagebuilder.DistributionConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetImageBuilderDistributionConfigurationWithName retrieves all imagebuilder.DistributionConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetImageBuilderDistributionConfigurationWithName(name string) (*imagebuilder.DistributionConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *imagebuilder.DistributionConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type imagebuilder.DistributionConfiguration not found", name)
}

// GetAllImageBuilderImageResources retrieves all imagebuilder.Image items from an AWS CloudFormation template
func (t *Template) GetAllImageBuilderImageResources() map[string]*imagebuilder.Image {
	results := map[string]*imagebuilder.Image{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *imagebuilder.Image:
			results[name] = resource
		}
	}
	return results
}

// GetImageBuilderImageWithName retrieves all imagebuilder.Image items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetImageBuilderImageWithName(name string) (*imagebuilder.Image, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *imagebuilder.Image:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type imagebuilder.Image not found", name)
}

// GetAllImageBuilderImagePipelineResources retrieves all imagebuilder.ImagePipeline items from an AWS CloudFormation template
func (t *Template) GetAllImageBuilderImagePipelineResources() map[string]*imagebuilder.ImagePipeline {
	results := map[string]*imagebuilder.ImagePipeline{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *imagebuilder.ImagePipeline:
			results[name] = resource
		}
	}
	return results
}

// GetImageBuilderImagePipelineWithName retrieves all imagebuilder.ImagePipeline items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetImageBuilderImagePipelineWithName(name string) (*imagebuilder.ImagePipeline, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *imagebuilder.ImagePipeline:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type imagebuilder.ImagePipeline not found", name)
}

// GetAllImageBuilderImageRecipeResources retrieves all imagebuilder.ImageRecipe items from an AWS CloudFormation template
func (t *Template) GetAllImageBuilderImageRecipeResources() map[string]*imagebuilder.ImageRecipe {
	results := map[string]*imagebuilder.ImageRecipe{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *imagebuilder.ImageRecipe:
			results[name] = resource
		}
	}
	return results
}

// GetImageBuilderImageRecipeWithName retrieves all imagebuilder.ImageRecipe items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetImageBuilderImageRecipeWithName(name string) (*imagebuilder.ImageRecipe, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *imagebuilder.ImageRecipe:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type imagebuilder.ImageRecipe not found", name)
}

// GetAllImageBuilderInfrastructureConfigurationResources retrieves all imagebuilder.InfrastructureConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllImageBuilderInfrastructureConfigurationResources() map[string]*imagebuilder.InfrastructureConfiguration {
	results := map[string]*imagebuilder.InfrastructureConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *imagebuilder.InfrastructureConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetImageBuilderInfrastructureConfigurationWithName retrieves all imagebuilder.InfrastructureConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetImageBuilderInfrastructureConfigurationWithName(name string) (*imagebuilder.InfrastructureConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *imagebuilder.InfrastructureConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type imagebuilder.InfrastructureConfiguration not found", name)
}

// GetAllInspectorAssessmentTargetResources retrieves all inspector.AssessmentTarget items from an AWS CloudFormation template
func (t *Template) GetAllInspectorAssessmentTargetResources() map[string]*inspector.AssessmentTarget {
	results := map[string]*inspector.AssessmentTarget{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *inspector.AssessmentTarget:
			results[name] = resource
		}
	}
	return results
}

// GetInspectorAssessmentTargetWithName retrieves all inspector.AssessmentTarget items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetInspectorAssessmentTargetWithName(name string) (*inspector.AssessmentTarget, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *inspector.AssessmentTarget:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type inspector.AssessmentTarget not found", name)
}

// GetAllInspectorAssessmentTemplateResources retrieves all inspector.AssessmentTemplate items from an AWS CloudFormation template
func (t *Template) GetAllInspectorAssessmentTemplateResources() map[string]*inspector.AssessmentTemplate {
	results := map[string]*inspector.AssessmentTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *inspector.AssessmentTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetInspectorAssessmentTemplateWithName retrieves all inspector.AssessmentTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetInspectorAssessmentTemplateWithName(name string) (*inspector.AssessmentTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *inspector.AssessmentTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type inspector.AssessmentTemplate not found", name)
}

// GetAllInspectorResourceGroupResources retrieves all inspector.ResourceGroup items from an AWS CloudFormation template
func (t *Template) GetAllInspectorResourceGroupResources() map[string]*inspector.ResourceGroup {
	results := map[string]*inspector.ResourceGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *inspector.ResourceGroup:
			results[name] = resource
		}
	}
	return results
}

// GetInspectorResourceGroupWithName retrieves all inspector.ResourceGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetInspectorResourceGroupWithName(name string) (*inspector.ResourceGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *inspector.ResourceGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type inspector.ResourceGroup not found", name)
}

// GetAllIoT1ClickDeviceResources retrieves all iot1click.Device items from an AWS CloudFormation template
func (t *Template) GetAllIoT1ClickDeviceResources() map[string]*iot1click.Device {
	results := map[string]*iot1click.Device{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot1click.Device:
			results[name] = resource
		}
	}
	return results
}

// GetIoT1ClickDeviceWithName retrieves all iot1click.Device items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoT1ClickDeviceWithName(name string) (*iot1click.Device, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot1click.Device:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot1click.Device not found", name)
}

// GetAllIoT1ClickPlacementResources retrieves all iot1click.Placement items from an AWS CloudFormation template
func (t *Template) GetAllIoT1ClickPlacementResources() map[string]*iot1click.Placement {
	results := map[string]*iot1click.Placement{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot1click.Placement:
			results[name] = resource
		}
	}
	return results
}

// GetIoT1ClickPlacementWithName retrieves all iot1click.Placement items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoT1ClickPlacementWithName(name string) (*iot1click.Placement, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot1click.Placement:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot1click.Placement not found", name)
}

// GetAllIoT1ClickProjectResources retrieves all iot1click.Project items from an AWS CloudFormation template
func (t *Template) GetAllIoT1ClickProjectResources() map[string]*iot1click.Project {
	results := map[string]*iot1click.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot1click.Project:
			results[name] = resource
		}
	}
	return results
}

// GetIoT1ClickProjectWithName retrieves all iot1click.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoT1ClickProjectWithName(name string) (*iot1click.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot1click.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot1click.Project not found", name)
}

// GetAllIoTAccountAuditConfigurationResources retrieves all iot.AccountAuditConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllIoTAccountAuditConfigurationResources() map[string]*iot.AccountAuditConfiguration {
	results := map[string]*iot.AccountAuditConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.AccountAuditConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetIoTAccountAuditConfigurationWithName retrieves all iot.AccountAuditConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTAccountAuditConfigurationWithName(name string) (*iot.AccountAuditConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.AccountAuditConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.AccountAuditConfiguration not found", name)
}

// GetAllIoTAuthorizerResources retrieves all iot.Authorizer items from an AWS CloudFormation template
func (t *Template) GetAllIoTAuthorizerResources() map[string]*iot.Authorizer {
	results := map[string]*iot.Authorizer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.Authorizer:
			results[name] = resource
		}
	}
	return results
}

// GetIoTAuthorizerWithName retrieves all iot.Authorizer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTAuthorizerWithName(name string) (*iot.Authorizer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.Authorizer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.Authorizer not found", name)
}

// GetAllIoTCertificateResources retrieves all iot.Certificate items from an AWS CloudFormation template
func (t *Template) GetAllIoTCertificateResources() map[string]*iot.Certificate {
	results := map[string]*iot.Certificate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.Certificate:
			results[name] = resource
		}
	}
	return results
}

// GetIoTCertificateWithName retrieves all iot.Certificate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTCertificateWithName(name string) (*iot.Certificate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.Certificate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.Certificate not found", name)
}

// GetAllIoTCustomMetricResources retrieves all iot.CustomMetric items from an AWS CloudFormation template
func (t *Template) GetAllIoTCustomMetricResources() map[string]*iot.CustomMetric {
	results := map[string]*iot.CustomMetric{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.CustomMetric:
			results[name] = resource
		}
	}
	return results
}

// GetIoTCustomMetricWithName retrieves all iot.CustomMetric items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTCustomMetricWithName(name string) (*iot.CustomMetric, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.CustomMetric:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.CustomMetric not found", name)
}

// GetAllIoTDimensionResources retrieves all iot.Dimension items from an AWS CloudFormation template
func (t *Template) GetAllIoTDimensionResources() map[string]*iot.Dimension {
	results := map[string]*iot.Dimension{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.Dimension:
			results[name] = resource
		}
	}
	return results
}

// GetIoTDimensionWithName retrieves all iot.Dimension items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTDimensionWithName(name string) (*iot.Dimension, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.Dimension:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.Dimension not found", name)
}

// GetAllIoTDomainConfigurationResources retrieves all iot.DomainConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllIoTDomainConfigurationResources() map[string]*iot.DomainConfiguration {
	results := map[string]*iot.DomainConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.DomainConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetIoTDomainConfigurationWithName retrieves all iot.DomainConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTDomainConfigurationWithName(name string) (*iot.DomainConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.DomainConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.DomainConfiguration not found", name)
}

// GetAllIoTFleetMetricResources retrieves all iot.FleetMetric items from an AWS CloudFormation template
func (t *Template) GetAllIoTFleetMetricResources() map[string]*iot.FleetMetric {
	results := map[string]*iot.FleetMetric{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.FleetMetric:
			results[name] = resource
		}
	}
	return results
}

// GetIoTFleetMetricWithName retrieves all iot.FleetMetric items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTFleetMetricWithName(name string) (*iot.FleetMetric, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.FleetMetric:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.FleetMetric not found", name)
}

// GetAllIoTJobTemplateResources retrieves all iot.JobTemplate items from an AWS CloudFormation template
func (t *Template) GetAllIoTJobTemplateResources() map[string]*iot.JobTemplate {
	results := map[string]*iot.JobTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.JobTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetIoTJobTemplateWithName retrieves all iot.JobTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTJobTemplateWithName(name string) (*iot.JobTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.JobTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.JobTemplate not found", name)
}

// GetAllIoTLoggingResources retrieves all iot.Logging items from an AWS CloudFormation template
func (t *Template) GetAllIoTLoggingResources() map[string]*iot.Logging {
	results := map[string]*iot.Logging{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.Logging:
			results[name] = resource
		}
	}
	return results
}

// GetIoTLoggingWithName retrieves all iot.Logging items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTLoggingWithName(name string) (*iot.Logging, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.Logging:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.Logging not found", name)
}

// GetAllIoTMitigationActionResources retrieves all iot.MitigationAction items from an AWS CloudFormation template
func (t *Template) GetAllIoTMitigationActionResources() map[string]*iot.MitigationAction {
	results := map[string]*iot.MitigationAction{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.MitigationAction:
			results[name] = resource
		}
	}
	return results
}

// GetIoTMitigationActionWithName retrieves all iot.MitigationAction items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTMitigationActionWithName(name string) (*iot.MitigationAction, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.MitigationAction:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.MitigationAction not found", name)
}

// GetAllIoTPolicyResources retrieves all iot.Policy items from an AWS CloudFormation template
func (t *Template) GetAllIoTPolicyResources() map[string]*iot.Policy {
	results := map[string]*iot.Policy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.Policy:
			results[name] = resource
		}
	}
	return results
}

// GetIoTPolicyWithName retrieves all iot.Policy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTPolicyWithName(name string) (*iot.Policy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.Policy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.Policy not found", name)
}

// GetAllIoTPolicyPrincipalAttachmentResources retrieves all iot.PolicyPrincipalAttachment items from an AWS CloudFormation template
func (t *Template) GetAllIoTPolicyPrincipalAttachmentResources() map[string]*iot.PolicyPrincipalAttachment {
	results := map[string]*iot.PolicyPrincipalAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.PolicyPrincipalAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetIoTPolicyPrincipalAttachmentWithName retrieves all iot.PolicyPrincipalAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTPolicyPrincipalAttachmentWithName(name string) (*iot.PolicyPrincipalAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.PolicyPrincipalAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.PolicyPrincipalAttachment not found", name)
}

// GetAllIoTProvisioningTemplateResources retrieves all iot.ProvisioningTemplate items from an AWS CloudFormation template
func (t *Template) GetAllIoTProvisioningTemplateResources() map[string]*iot.ProvisioningTemplate {
	results := map[string]*iot.ProvisioningTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.ProvisioningTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetIoTProvisioningTemplateWithName retrieves all iot.ProvisioningTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTProvisioningTemplateWithName(name string) (*iot.ProvisioningTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.ProvisioningTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.ProvisioningTemplate not found", name)
}

// GetAllIoTResourceSpecificLoggingResources retrieves all iot.ResourceSpecificLogging items from an AWS CloudFormation template
func (t *Template) GetAllIoTResourceSpecificLoggingResources() map[string]*iot.ResourceSpecificLogging {
	results := map[string]*iot.ResourceSpecificLogging{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.ResourceSpecificLogging:
			results[name] = resource
		}
	}
	return results
}

// GetIoTResourceSpecificLoggingWithName retrieves all iot.ResourceSpecificLogging items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTResourceSpecificLoggingWithName(name string) (*iot.ResourceSpecificLogging, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.ResourceSpecificLogging:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.ResourceSpecificLogging not found", name)
}

// GetAllIoTScheduledAuditResources retrieves all iot.ScheduledAudit items from an AWS CloudFormation template
func (t *Template) GetAllIoTScheduledAuditResources() map[string]*iot.ScheduledAudit {
	results := map[string]*iot.ScheduledAudit{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.ScheduledAudit:
			results[name] = resource
		}
	}
	return results
}

// GetIoTScheduledAuditWithName retrieves all iot.ScheduledAudit items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTScheduledAuditWithName(name string) (*iot.ScheduledAudit, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.ScheduledAudit:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.ScheduledAudit not found", name)
}

// GetAllIoTSecurityProfileResources retrieves all iot.SecurityProfile items from an AWS CloudFormation template
func (t *Template) GetAllIoTSecurityProfileResources() map[string]*iot.SecurityProfile {
	results := map[string]*iot.SecurityProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.SecurityProfile:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSecurityProfileWithName retrieves all iot.SecurityProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSecurityProfileWithName(name string) (*iot.SecurityProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.SecurityProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.SecurityProfile not found", name)
}

// GetAllIoTThingResources retrieves all iot.Thing items from an AWS CloudFormation template
func (t *Template) GetAllIoTThingResources() map[string]*iot.Thing {
	results := map[string]*iot.Thing{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.Thing:
			results[name] = resource
		}
	}
	return results
}

// GetIoTThingWithName retrieves all iot.Thing items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTThingWithName(name string) (*iot.Thing, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.Thing:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.Thing not found", name)
}

// GetAllIoTThingPrincipalAttachmentResources retrieves all iot.ThingPrincipalAttachment items from an AWS CloudFormation template
func (t *Template) GetAllIoTThingPrincipalAttachmentResources() map[string]*iot.ThingPrincipalAttachment {
	results := map[string]*iot.ThingPrincipalAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.ThingPrincipalAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetIoTThingPrincipalAttachmentWithName retrieves all iot.ThingPrincipalAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTThingPrincipalAttachmentWithName(name string) (*iot.ThingPrincipalAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.ThingPrincipalAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.ThingPrincipalAttachment not found", name)
}

// GetAllIoTTopicRuleResources retrieves all iot.TopicRule items from an AWS CloudFormation template
func (t *Template) GetAllIoTTopicRuleResources() map[string]*iot.TopicRule {
	results := map[string]*iot.TopicRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.TopicRule:
			results[name] = resource
		}
	}
	return results
}

// GetIoTTopicRuleWithName retrieves all iot.TopicRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTTopicRuleWithName(name string) (*iot.TopicRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.TopicRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.TopicRule not found", name)
}

// GetAllIoTTopicRuleDestinationResources retrieves all iot.TopicRuleDestination items from an AWS CloudFormation template
func (t *Template) GetAllIoTTopicRuleDestinationResources() map[string]*iot.TopicRuleDestination {
	results := map[string]*iot.TopicRuleDestination{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iot.TopicRuleDestination:
			results[name] = resource
		}
	}
	return results
}

// GetIoTTopicRuleDestinationWithName retrieves all iot.TopicRuleDestination items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTTopicRuleDestinationWithName(name string) (*iot.TopicRuleDestination, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iot.TopicRuleDestination:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iot.TopicRuleDestination not found", name)
}

// GetAllIoTAnalyticsChannelResources retrieves all iotanalytics.Channel items from an AWS CloudFormation template
func (t *Template) GetAllIoTAnalyticsChannelResources() map[string]*iotanalytics.Channel {
	results := map[string]*iotanalytics.Channel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotanalytics.Channel:
			results[name] = resource
		}
	}
	return results
}

// GetIoTAnalyticsChannelWithName retrieves all iotanalytics.Channel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTAnalyticsChannelWithName(name string) (*iotanalytics.Channel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotanalytics.Channel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotanalytics.Channel not found", name)
}

// GetAllIoTAnalyticsDatasetResources retrieves all iotanalytics.Dataset items from an AWS CloudFormation template
func (t *Template) GetAllIoTAnalyticsDatasetResources() map[string]*iotanalytics.Dataset {
	results := map[string]*iotanalytics.Dataset{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotanalytics.Dataset:
			results[name] = resource
		}
	}
	return results
}

// GetIoTAnalyticsDatasetWithName retrieves all iotanalytics.Dataset items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTAnalyticsDatasetWithName(name string) (*iotanalytics.Dataset, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotanalytics.Dataset:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotanalytics.Dataset not found", name)
}

// GetAllIoTAnalyticsDatastoreResources retrieves all iotanalytics.Datastore items from an AWS CloudFormation template
func (t *Template) GetAllIoTAnalyticsDatastoreResources() map[string]*iotanalytics.Datastore {
	results := map[string]*iotanalytics.Datastore{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotanalytics.Datastore:
			results[name] = resource
		}
	}
	return results
}

// GetIoTAnalyticsDatastoreWithName retrieves all iotanalytics.Datastore items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTAnalyticsDatastoreWithName(name string) (*iotanalytics.Datastore, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotanalytics.Datastore:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotanalytics.Datastore not found", name)
}

// GetAllIoTAnalyticsPipelineResources retrieves all iotanalytics.Pipeline items from an AWS CloudFormation template
func (t *Template) GetAllIoTAnalyticsPipelineResources() map[string]*iotanalytics.Pipeline {
	results := map[string]*iotanalytics.Pipeline{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotanalytics.Pipeline:
			results[name] = resource
		}
	}
	return results
}

// GetIoTAnalyticsPipelineWithName retrieves all iotanalytics.Pipeline items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTAnalyticsPipelineWithName(name string) (*iotanalytics.Pipeline, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotanalytics.Pipeline:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotanalytics.Pipeline not found", name)
}

// GetAllIoTCoreDeviceAdvisorSuiteDefinitionResources retrieves all iotcoredeviceadvisor.SuiteDefinition items from an AWS CloudFormation template
func (t *Template) GetAllIoTCoreDeviceAdvisorSuiteDefinitionResources() map[string]*iotcoredeviceadvisor.SuiteDefinition {
	results := map[string]*iotcoredeviceadvisor.SuiteDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotcoredeviceadvisor.SuiteDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetIoTCoreDeviceAdvisorSuiteDefinitionWithName retrieves all iotcoredeviceadvisor.SuiteDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTCoreDeviceAdvisorSuiteDefinitionWithName(name string) (*iotcoredeviceadvisor.SuiteDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotcoredeviceadvisor.SuiteDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotcoredeviceadvisor.SuiteDefinition not found", name)
}

// GetAllIoTEventsDetectorModelResources retrieves all iotevents.DetectorModel items from an AWS CloudFormation template
func (t *Template) GetAllIoTEventsDetectorModelResources() map[string]*iotevents.DetectorModel {
	results := map[string]*iotevents.DetectorModel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotevents.DetectorModel:
			results[name] = resource
		}
	}
	return results
}

// GetIoTEventsDetectorModelWithName retrieves all iotevents.DetectorModel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTEventsDetectorModelWithName(name string) (*iotevents.DetectorModel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotevents.DetectorModel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotevents.DetectorModel not found", name)
}

// GetAllIoTEventsInputResources retrieves all iotevents.Input items from an AWS CloudFormation template
func (t *Template) GetAllIoTEventsInputResources() map[string]*iotevents.Input {
	results := map[string]*iotevents.Input{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotevents.Input:
			results[name] = resource
		}
	}
	return results
}

// GetIoTEventsInputWithName retrieves all iotevents.Input items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTEventsInputWithName(name string) (*iotevents.Input, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotevents.Input:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotevents.Input not found", name)
}

// GetAllIoTFleetHubApplicationResources retrieves all iotfleethub.Application items from an AWS CloudFormation template
func (t *Template) GetAllIoTFleetHubApplicationResources() map[string]*iotfleethub.Application {
	results := map[string]*iotfleethub.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotfleethub.Application:
			results[name] = resource
		}
	}
	return results
}

// GetIoTFleetHubApplicationWithName retrieves all iotfleethub.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTFleetHubApplicationWithName(name string) (*iotfleethub.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotfleethub.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotfleethub.Application not found", name)
}

// GetAllIoTSiteWiseAccessPolicyResources retrieves all iotsitewise.AccessPolicy items from an AWS CloudFormation template
func (t *Template) GetAllIoTSiteWiseAccessPolicyResources() map[string]*iotsitewise.AccessPolicy {
	results := map[string]*iotsitewise.AccessPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotsitewise.AccessPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSiteWiseAccessPolicyWithName retrieves all iotsitewise.AccessPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSiteWiseAccessPolicyWithName(name string) (*iotsitewise.AccessPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotsitewise.AccessPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotsitewise.AccessPolicy not found", name)
}

// GetAllIoTSiteWiseAssetResources retrieves all iotsitewise.Asset items from an AWS CloudFormation template
func (t *Template) GetAllIoTSiteWiseAssetResources() map[string]*iotsitewise.Asset {
	results := map[string]*iotsitewise.Asset{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotsitewise.Asset:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSiteWiseAssetWithName retrieves all iotsitewise.Asset items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSiteWiseAssetWithName(name string) (*iotsitewise.Asset, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotsitewise.Asset:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotsitewise.Asset not found", name)
}

// GetAllIoTSiteWiseAssetModelResources retrieves all iotsitewise.AssetModel items from an AWS CloudFormation template
func (t *Template) GetAllIoTSiteWiseAssetModelResources() map[string]*iotsitewise.AssetModel {
	results := map[string]*iotsitewise.AssetModel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotsitewise.AssetModel:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSiteWiseAssetModelWithName retrieves all iotsitewise.AssetModel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSiteWiseAssetModelWithName(name string) (*iotsitewise.AssetModel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotsitewise.AssetModel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotsitewise.AssetModel not found", name)
}

// GetAllIoTSiteWiseDashboardResources retrieves all iotsitewise.Dashboard items from an AWS CloudFormation template
func (t *Template) GetAllIoTSiteWiseDashboardResources() map[string]*iotsitewise.Dashboard {
	results := map[string]*iotsitewise.Dashboard{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotsitewise.Dashboard:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSiteWiseDashboardWithName retrieves all iotsitewise.Dashboard items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSiteWiseDashboardWithName(name string) (*iotsitewise.Dashboard, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotsitewise.Dashboard:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotsitewise.Dashboard not found", name)
}

// GetAllIoTSiteWiseGatewayResources retrieves all iotsitewise.Gateway items from an AWS CloudFormation template
func (t *Template) GetAllIoTSiteWiseGatewayResources() map[string]*iotsitewise.Gateway {
	results := map[string]*iotsitewise.Gateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotsitewise.Gateway:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSiteWiseGatewayWithName retrieves all iotsitewise.Gateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSiteWiseGatewayWithName(name string) (*iotsitewise.Gateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotsitewise.Gateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotsitewise.Gateway not found", name)
}

// GetAllIoTSiteWisePortalResources retrieves all iotsitewise.Portal items from an AWS CloudFormation template
func (t *Template) GetAllIoTSiteWisePortalResources() map[string]*iotsitewise.Portal {
	results := map[string]*iotsitewise.Portal{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotsitewise.Portal:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSiteWisePortalWithName retrieves all iotsitewise.Portal items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSiteWisePortalWithName(name string) (*iotsitewise.Portal, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotsitewise.Portal:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotsitewise.Portal not found", name)
}

// GetAllIoTSiteWiseProjectResources retrieves all iotsitewise.Project items from an AWS CloudFormation template
func (t *Template) GetAllIoTSiteWiseProjectResources() map[string]*iotsitewise.Project {
	results := map[string]*iotsitewise.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotsitewise.Project:
			results[name] = resource
		}
	}
	return results
}

// GetIoTSiteWiseProjectWithName retrieves all iotsitewise.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTSiteWiseProjectWithName(name string) (*iotsitewise.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotsitewise.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotsitewise.Project not found", name)
}

// GetAllIoTThingsGraphFlowTemplateResources retrieves all iotthingsgraph.FlowTemplate items from an AWS CloudFormation template
func (t *Template) GetAllIoTThingsGraphFlowTemplateResources() map[string]*iotthingsgraph.FlowTemplate {
	results := map[string]*iotthingsgraph.FlowTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotthingsgraph.FlowTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetIoTThingsGraphFlowTemplateWithName retrieves all iotthingsgraph.FlowTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTThingsGraphFlowTemplateWithName(name string) (*iotthingsgraph.FlowTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotthingsgraph.FlowTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotthingsgraph.FlowTemplate not found", name)
}

// GetAllIoTWirelessDestinationResources retrieves all iotwireless.Destination items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessDestinationResources() map[string]*iotwireless.Destination {
	results := map[string]*iotwireless.Destination{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.Destination:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessDestinationWithName retrieves all iotwireless.Destination items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessDestinationWithName(name string) (*iotwireless.Destination, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.Destination:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.Destination not found", name)
}

// GetAllIoTWirelessDeviceProfileResources retrieves all iotwireless.DeviceProfile items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessDeviceProfileResources() map[string]*iotwireless.DeviceProfile {
	results := map[string]*iotwireless.DeviceProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.DeviceProfile:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessDeviceProfileWithName retrieves all iotwireless.DeviceProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessDeviceProfileWithName(name string) (*iotwireless.DeviceProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.DeviceProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.DeviceProfile not found", name)
}

// GetAllIoTWirelessFuotaTaskResources retrieves all iotwireless.FuotaTask items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessFuotaTaskResources() map[string]*iotwireless.FuotaTask {
	results := map[string]*iotwireless.FuotaTask{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.FuotaTask:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessFuotaTaskWithName retrieves all iotwireless.FuotaTask items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessFuotaTaskWithName(name string) (*iotwireless.FuotaTask, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.FuotaTask:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.FuotaTask not found", name)
}

// GetAllIoTWirelessMulticastGroupResources retrieves all iotwireless.MulticastGroup items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessMulticastGroupResources() map[string]*iotwireless.MulticastGroup {
	results := map[string]*iotwireless.MulticastGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.MulticastGroup:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessMulticastGroupWithName retrieves all iotwireless.MulticastGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessMulticastGroupWithName(name string) (*iotwireless.MulticastGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.MulticastGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.MulticastGroup not found", name)
}

// GetAllIoTWirelessPartnerAccountResources retrieves all iotwireless.PartnerAccount items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessPartnerAccountResources() map[string]*iotwireless.PartnerAccount {
	results := map[string]*iotwireless.PartnerAccount{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.PartnerAccount:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessPartnerAccountWithName retrieves all iotwireless.PartnerAccount items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessPartnerAccountWithName(name string) (*iotwireless.PartnerAccount, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.PartnerAccount:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.PartnerAccount not found", name)
}

// GetAllIoTWirelessServiceProfileResources retrieves all iotwireless.ServiceProfile items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessServiceProfileResources() map[string]*iotwireless.ServiceProfile {
	results := map[string]*iotwireless.ServiceProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.ServiceProfile:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessServiceProfileWithName retrieves all iotwireless.ServiceProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessServiceProfileWithName(name string) (*iotwireless.ServiceProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.ServiceProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.ServiceProfile not found", name)
}

// GetAllIoTWirelessTaskDefinitionResources retrieves all iotwireless.TaskDefinition items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessTaskDefinitionResources() map[string]*iotwireless.TaskDefinition {
	results := map[string]*iotwireless.TaskDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.TaskDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessTaskDefinitionWithName retrieves all iotwireless.TaskDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessTaskDefinitionWithName(name string) (*iotwireless.TaskDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.TaskDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.TaskDefinition not found", name)
}

// GetAllIoTWirelessWirelessDeviceResources retrieves all iotwireless.WirelessDevice items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessWirelessDeviceResources() map[string]*iotwireless.WirelessDevice {
	results := map[string]*iotwireless.WirelessDevice{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.WirelessDevice:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessWirelessDeviceWithName retrieves all iotwireless.WirelessDevice items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessWirelessDeviceWithName(name string) (*iotwireless.WirelessDevice, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.WirelessDevice:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.WirelessDevice not found", name)
}

// GetAllIoTWirelessWirelessGatewayResources retrieves all iotwireless.WirelessGateway items from an AWS CloudFormation template
func (t *Template) GetAllIoTWirelessWirelessGatewayResources() map[string]*iotwireless.WirelessGateway {
	results := map[string]*iotwireless.WirelessGateway{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iotwireless.WirelessGateway:
			results[name] = resource
		}
	}
	return results
}

// GetIoTWirelessWirelessGatewayWithName retrieves all iotwireless.WirelessGateway items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIoTWirelessWirelessGatewayWithName(name string) (*iotwireless.WirelessGateway, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iotwireless.WirelessGateway:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iotwireless.WirelessGateway not found", name)
}

// GetAllKMSAliasResources retrieves all kms.Alias items from an AWS CloudFormation template
func (t *Template) GetAllKMSAliasResources() map[string]*kms.Alias {
	results := map[string]*kms.Alias{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kms.Alias:
			results[name] = resource
		}
	}
	return results
}

// GetKMSAliasWithName retrieves all kms.Alias items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKMSAliasWithName(name string) (*kms.Alias, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kms.Alias:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kms.Alias not found", name)
}

// GetAllKMSKeyResources retrieves all kms.Key items from an AWS CloudFormation template
func (t *Template) GetAllKMSKeyResources() map[string]*kms.Key {
	results := map[string]*kms.Key{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kms.Key:
			results[name] = resource
		}
	}
	return results
}

// GetKMSKeyWithName retrieves all kms.Key items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKMSKeyWithName(name string) (*kms.Key, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kms.Key:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kms.Key not found", name)
}

// GetAllKMSReplicaKeyResources retrieves all kms.ReplicaKey items from an AWS CloudFormation template
func (t *Template) GetAllKMSReplicaKeyResources() map[string]*kms.ReplicaKey {
	results := map[string]*kms.ReplicaKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kms.ReplicaKey:
			results[name] = resource
		}
	}
	return results
}

// GetKMSReplicaKeyWithName retrieves all kms.ReplicaKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKMSReplicaKeyWithName(name string) (*kms.ReplicaKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kms.ReplicaKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kms.ReplicaKey not found", name)
}

// GetAllKendraDataSourceResources retrieves all kendra.DataSource items from an AWS CloudFormation template
func (t *Template) GetAllKendraDataSourceResources() map[string]*kendra.DataSource {
	results := map[string]*kendra.DataSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kendra.DataSource:
			results[name] = resource
		}
	}
	return results
}

// GetKendraDataSourceWithName retrieves all kendra.DataSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKendraDataSourceWithName(name string) (*kendra.DataSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kendra.DataSource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kendra.DataSource not found", name)
}

// GetAllKendraFaqResources retrieves all kendra.Faq items from an AWS CloudFormation template
func (t *Template) GetAllKendraFaqResources() map[string]*kendra.Faq {
	results := map[string]*kendra.Faq{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kendra.Faq:
			results[name] = resource
		}
	}
	return results
}

// GetKendraFaqWithName retrieves all kendra.Faq items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKendraFaqWithName(name string) (*kendra.Faq, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kendra.Faq:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kendra.Faq not found", name)
}

// GetAllKendraIndexResources retrieves all kendra.Index items from an AWS CloudFormation template
func (t *Template) GetAllKendraIndexResources() map[string]*kendra.Index {
	results := map[string]*kendra.Index{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kendra.Index:
			results[name] = resource
		}
	}
	return results
}

// GetKendraIndexWithName retrieves all kendra.Index items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKendraIndexWithName(name string) (*kendra.Index, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kendra.Index:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kendra.Index not found", name)
}

// GetAllKinesisStreamResources retrieves all kinesis.Stream items from an AWS CloudFormation template
func (t *Template) GetAllKinesisStreamResources() map[string]*kinesis.Stream {
	results := map[string]*kinesis.Stream{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesis.Stream:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisStreamWithName retrieves all kinesis.Stream items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisStreamWithName(name string) (*kinesis.Stream, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesis.Stream:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesis.Stream not found", name)
}

// GetAllKinesisStreamConsumerResources retrieves all kinesis.StreamConsumer items from an AWS CloudFormation template
func (t *Template) GetAllKinesisStreamConsumerResources() map[string]*kinesis.StreamConsumer {
	results := map[string]*kinesis.StreamConsumer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesis.StreamConsumer:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisStreamConsumerWithName retrieves all kinesis.StreamConsumer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisStreamConsumerWithName(name string) (*kinesis.StreamConsumer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesis.StreamConsumer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesis.StreamConsumer not found", name)
}

// GetAllKinesisAnalyticsApplicationResources retrieves all kinesisanalytics.Application items from an AWS CloudFormation template
func (t *Template) GetAllKinesisAnalyticsApplicationResources() map[string]*kinesisanalytics.Application {
	results := map[string]*kinesisanalytics.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisanalytics.Application:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisAnalyticsApplicationWithName retrieves all kinesisanalytics.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisAnalyticsApplicationWithName(name string) (*kinesisanalytics.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisanalytics.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisanalytics.Application not found", name)
}

// GetAllKinesisAnalyticsApplicationOutputResources retrieves all kinesisanalytics.ApplicationOutput items from an AWS CloudFormation template
func (t *Template) GetAllKinesisAnalyticsApplicationOutputResources() map[string]*kinesisanalytics.ApplicationOutput {
	results := map[string]*kinesisanalytics.ApplicationOutput{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisanalytics.ApplicationOutput:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisAnalyticsApplicationOutputWithName retrieves all kinesisanalytics.ApplicationOutput items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisAnalyticsApplicationOutputWithName(name string) (*kinesisanalytics.ApplicationOutput, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisanalytics.ApplicationOutput:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisanalytics.ApplicationOutput not found", name)
}

// GetAllKinesisAnalyticsApplicationReferenceDataSourceResources retrieves all kinesisanalytics.ApplicationReferenceDataSource items from an AWS CloudFormation template
func (t *Template) GetAllKinesisAnalyticsApplicationReferenceDataSourceResources() map[string]*kinesisanalytics.ApplicationReferenceDataSource {
	results := map[string]*kinesisanalytics.ApplicationReferenceDataSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisanalytics.ApplicationReferenceDataSource:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisAnalyticsApplicationReferenceDataSourceWithName retrieves all kinesisanalytics.ApplicationReferenceDataSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisAnalyticsApplicationReferenceDataSourceWithName(name string) (*kinesisanalytics.ApplicationReferenceDataSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisanalytics.ApplicationReferenceDataSource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisanalytics.ApplicationReferenceDataSource not found", name)
}

// GetAllKinesisAnalyticsV2ApplicationResources retrieves all kinesisanalyticsv2.Application items from an AWS CloudFormation template
func (t *Template) GetAllKinesisAnalyticsV2ApplicationResources() map[string]*kinesisanalyticsv2.Application {
	results := map[string]*kinesisanalyticsv2.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.Application:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisAnalyticsV2ApplicationWithName retrieves all kinesisanalyticsv2.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisAnalyticsV2ApplicationWithName(name string) (*kinesisanalyticsv2.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisanalyticsv2.Application not found", name)
}

// GetAllKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionResources retrieves all kinesisanalyticsv2.ApplicationCloudWatchLoggingOption items from an AWS CloudFormation template
func (t *Template) GetAllKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionResources() map[string]*kinesisanalyticsv2.ApplicationCloudWatchLoggingOption {
	results := map[string]*kinesisanalyticsv2.ApplicationCloudWatchLoggingOption{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.ApplicationCloudWatchLoggingOption:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionWithName retrieves all kinesisanalyticsv2.ApplicationCloudWatchLoggingOption items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionWithName(name string) (*kinesisanalyticsv2.ApplicationCloudWatchLoggingOption, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.ApplicationCloudWatchLoggingOption:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisanalyticsv2.ApplicationCloudWatchLoggingOption not found", name)
}

// GetAllKinesisAnalyticsV2ApplicationOutputResources retrieves all kinesisanalyticsv2.ApplicationOutput items from an AWS CloudFormation template
func (t *Template) GetAllKinesisAnalyticsV2ApplicationOutputResources() map[string]*kinesisanalyticsv2.ApplicationOutput {
	results := map[string]*kinesisanalyticsv2.ApplicationOutput{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.ApplicationOutput:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisAnalyticsV2ApplicationOutputWithName retrieves all kinesisanalyticsv2.ApplicationOutput items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisAnalyticsV2ApplicationOutputWithName(name string) (*kinesisanalyticsv2.ApplicationOutput, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.ApplicationOutput:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisanalyticsv2.ApplicationOutput not found", name)
}

// GetAllKinesisAnalyticsV2ApplicationReferenceDataSourceResources retrieves all kinesisanalyticsv2.ApplicationReferenceDataSource items from an AWS CloudFormation template
func (t *Template) GetAllKinesisAnalyticsV2ApplicationReferenceDataSourceResources() map[string]*kinesisanalyticsv2.ApplicationReferenceDataSource {
	results := map[string]*kinesisanalyticsv2.ApplicationReferenceDataSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.ApplicationReferenceDataSource:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisAnalyticsV2ApplicationReferenceDataSourceWithName retrieves all kinesisanalyticsv2.ApplicationReferenceDataSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisAnalyticsV2ApplicationReferenceDataSourceWithName(name string) (*kinesisanalyticsv2.ApplicationReferenceDataSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisanalyticsv2.ApplicationReferenceDataSource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisanalyticsv2.ApplicationReferenceDataSource not found", name)
}

// GetAllKinesisFirehoseDeliveryStreamResources retrieves all kinesisfirehose.DeliveryStream items from an AWS CloudFormation template
func (t *Template) GetAllKinesisFirehoseDeliveryStreamResources() map[string]*kinesisfirehose.DeliveryStream {
	results := map[string]*kinesisfirehose.DeliveryStream{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesisfirehose.DeliveryStream:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisFirehoseDeliveryStreamWithName retrieves all kinesisfirehose.DeliveryStream items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisFirehoseDeliveryStreamWithName(name string) (*kinesisfirehose.DeliveryStream, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesisfirehose.DeliveryStream:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesisfirehose.DeliveryStream not found", name)
}

// GetAllLakeFormationDataLakeSettingsResources retrieves all lakeformation.DataLakeSettings items from an AWS CloudFormation template
func (t *Template) GetAllLakeFormationDataLakeSettingsResources() map[string]*lakeformation.DataLakeSettings {
	results := map[string]*lakeformation.DataLakeSettings{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lakeformation.DataLakeSettings:
			results[name] = resource
		}
	}
	return results
}

// GetLakeFormationDataLakeSettingsWithName retrieves all lakeformation.DataLakeSettings items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLakeFormationDataLakeSettingsWithName(name string) (*lakeformation.DataLakeSettings, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lakeformation.DataLakeSettings:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lakeformation.DataLakeSettings not found", name)
}

// GetAllLakeFormationPermissionsResources retrieves all lakeformation.Permissions items from an AWS CloudFormation template
func (t *Template) GetAllLakeFormationPermissionsResources() map[string]*lakeformation.Permissions {
	results := map[string]*lakeformation.Permissions{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lakeformation.Permissions:
			results[name] = resource
		}
	}
	return results
}

// GetLakeFormationPermissionsWithName retrieves all lakeformation.Permissions items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLakeFormationPermissionsWithName(name string) (*lakeformation.Permissions, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lakeformation.Permissions:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lakeformation.Permissions not found", name)
}

// GetAllLakeFormationResourceResources retrieves all lakeformation.Resource items from an AWS CloudFormation template
func (t *Template) GetAllLakeFormationResourceResources() map[string]*lakeformation.Resource {
	results := map[string]*lakeformation.Resource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lakeformation.Resource:
			results[name] = resource
		}
	}
	return results
}

// GetLakeFormationResourceWithName retrieves all lakeformation.Resource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLakeFormationResourceWithName(name string) (*lakeformation.Resource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lakeformation.Resource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lakeformation.Resource not found", name)
}

// GetAllLambdaAliasResources retrieves all lambda.Alias items from an AWS CloudFormation template
func (t *Template) GetAllLambdaAliasResources() map[string]*lambda.Alias {
	results := map[string]*lambda.Alias{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.Alias:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaAliasWithName retrieves all lambda.Alias items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaAliasWithName(name string) (*lambda.Alias, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.Alias:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.Alias not found", name)
}

// GetAllLambdaCodeSigningConfigResources retrieves all lambda.CodeSigningConfig items from an AWS CloudFormation template
func (t *Template) GetAllLambdaCodeSigningConfigResources() map[string]*lambda.CodeSigningConfig {
	results := map[string]*lambda.CodeSigningConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.CodeSigningConfig:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaCodeSigningConfigWithName retrieves all lambda.CodeSigningConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaCodeSigningConfigWithName(name string) (*lambda.CodeSigningConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.CodeSigningConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.CodeSigningConfig not found", name)
}

// GetAllLambdaEventInvokeConfigResources retrieves all lambda.EventInvokeConfig items from an AWS CloudFormation template
func (t *Template) GetAllLambdaEventInvokeConfigResources() map[string]*lambda.EventInvokeConfig {
	results := map[string]*lambda.EventInvokeConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.EventInvokeConfig:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaEventInvokeConfigWithName retrieves all lambda.EventInvokeConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaEventInvokeConfigWithName(name string) (*lambda.EventInvokeConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.EventInvokeConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.EventInvokeConfig not found", name)
}

// GetAllLambdaEventSourceMappingResources retrieves all lambda.EventSourceMapping items from an AWS CloudFormation template
func (t *Template) GetAllLambdaEventSourceMappingResources() map[string]*lambda.EventSourceMapping {
	results := map[string]*lambda.EventSourceMapping{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.EventSourceMapping:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaEventSourceMappingWithName retrieves all lambda.EventSourceMapping items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaEventSourceMappingWithName(name string) (*lambda.EventSourceMapping, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.EventSourceMapping:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.EventSourceMapping not found", name)
}

// GetAllLambdaFunctionResources retrieves all lambda.Function items from an AWS CloudFormation template
func (t *Template) GetAllLambdaFunctionResources() map[string]*lambda.Function {
	results := map[string]*lambda.Function{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.Function:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaFunctionWithName retrieves all lambda.Function items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaFunctionWithName(name string) (*lambda.Function, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.Function:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.Function not found", name)
}

// GetAllLambdaLayerVersionResources retrieves all lambda.LayerVersion items from an AWS CloudFormation template
func (t *Template) GetAllLambdaLayerVersionResources() map[string]*lambda.LayerVersion {
	results := map[string]*lambda.LayerVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.LayerVersion:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaLayerVersionWithName retrieves all lambda.LayerVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaLayerVersionWithName(name string) (*lambda.LayerVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.LayerVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.LayerVersion not found", name)
}

// GetAllLambdaLayerVersionPermissionResources retrieves all lambda.LayerVersionPermission items from an AWS CloudFormation template
func (t *Template) GetAllLambdaLayerVersionPermissionResources() map[string]*lambda.LayerVersionPermission {
	results := map[string]*lambda.LayerVersionPermission{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.LayerVersionPermission:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaLayerVersionPermissionWithName retrieves all lambda.LayerVersionPermission items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaLayerVersionPermissionWithName(name string) (*lambda.LayerVersionPermission, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.LayerVersionPermission:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.LayerVersionPermission not found", name)
}

// GetAllLambdaPermissionResources retrieves all lambda.Permission items from an AWS CloudFormation template
func (t *Template) GetAllLambdaPermissionResources() map[string]*lambda.Permission {
	results := map[string]*lambda.Permission{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.Permission:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaPermissionWithName retrieves all lambda.Permission items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaPermissionWithName(name string) (*lambda.Permission, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.Permission:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.Permission not found", name)
}

// GetAllLambdaVersionResources retrieves all lambda.Version items from an AWS CloudFormation template
func (t *Template) GetAllLambdaVersionResources() map[string]*lambda.Version {
	results := map[string]*lambda.Version{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.Version:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaVersionWithName retrieves all lambda.Version items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaVersionWithName(name string) (*lambda.Version, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.Version:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.Version not found", name)
}

// GetAllLicenseManagerGrantResources retrieves all licensemanager.Grant items from an AWS CloudFormation template
func (t *Template) GetAllLicenseManagerGrantResources() map[string]*licensemanager.Grant {
	results := map[string]*licensemanager.Grant{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *licensemanager.Grant:
			results[name] = resource
		}
	}
	return results
}

// GetLicenseManagerGrantWithName retrieves all licensemanager.Grant items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLicenseManagerGrantWithName(name string) (*licensemanager.Grant, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *licensemanager.Grant:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type licensemanager.Grant not found", name)
}

// GetAllLicenseManagerLicenseResources retrieves all licensemanager.License items from an AWS CloudFormation template
func (t *Template) GetAllLicenseManagerLicenseResources() map[string]*licensemanager.License {
	results := map[string]*licensemanager.License{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *licensemanager.License:
			results[name] = resource
		}
	}
	return results
}

// GetLicenseManagerLicenseWithName retrieves all licensemanager.License items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLicenseManagerLicenseWithName(name string) (*licensemanager.License, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *licensemanager.License:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type licensemanager.License not found", name)
}

// GetAllLightsailDatabaseResources retrieves all lightsail.Database items from an AWS CloudFormation template
func (t *Template) GetAllLightsailDatabaseResources() map[string]*lightsail.Database {
	results := map[string]*lightsail.Database{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lightsail.Database:
			results[name] = resource
		}
	}
	return results
}

// GetLightsailDatabaseWithName retrieves all lightsail.Database items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLightsailDatabaseWithName(name string) (*lightsail.Database, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lightsail.Database:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lightsail.Database not found", name)
}

// GetAllLightsailDiskResources retrieves all lightsail.Disk items from an AWS CloudFormation template
func (t *Template) GetAllLightsailDiskResources() map[string]*lightsail.Disk {
	results := map[string]*lightsail.Disk{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lightsail.Disk:
			results[name] = resource
		}
	}
	return results
}

// GetLightsailDiskWithName retrieves all lightsail.Disk items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLightsailDiskWithName(name string) (*lightsail.Disk, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lightsail.Disk:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lightsail.Disk not found", name)
}

// GetAllLightsailInstanceResources retrieves all lightsail.Instance items from an AWS CloudFormation template
func (t *Template) GetAllLightsailInstanceResources() map[string]*lightsail.Instance {
	results := map[string]*lightsail.Instance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lightsail.Instance:
			results[name] = resource
		}
	}
	return results
}

// GetLightsailInstanceWithName retrieves all lightsail.Instance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLightsailInstanceWithName(name string) (*lightsail.Instance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lightsail.Instance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lightsail.Instance not found", name)
}

// GetAllLightsailStaticIpResources retrieves all lightsail.StaticIp items from an AWS CloudFormation template
func (t *Template) GetAllLightsailStaticIpResources() map[string]*lightsail.StaticIp {
	results := map[string]*lightsail.StaticIp{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lightsail.StaticIp:
			results[name] = resource
		}
	}
	return results
}

// GetLightsailStaticIpWithName retrieves all lightsail.StaticIp items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLightsailStaticIpWithName(name string) (*lightsail.StaticIp, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lightsail.StaticIp:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lightsail.StaticIp not found", name)
}

// GetAllLocationGeofenceCollectionResources retrieves all location.GeofenceCollection items from an AWS CloudFormation template
func (t *Template) GetAllLocationGeofenceCollectionResources() map[string]*location.GeofenceCollection {
	results := map[string]*location.GeofenceCollection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *location.GeofenceCollection:
			results[name] = resource
		}
	}
	return results
}

// GetLocationGeofenceCollectionWithName retrieves all location.GeofenceCollection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLocationGeofenceCollectionWithName(name string) (*location.GeofenceCollection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *location.GeofenceCollection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type location.GeofenceCollection not found", name)
}

// GetAllLocationMapResources retrieves all location.Map items from an AWS CloudFormation template
func (t *Template) GetAllLocationMapResources() map[string]*location.Map {
	results := map[string]*location.Map{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *location.Map:
			results[name] = resource
		}
	}
	return results
}

// GetLocationMapWithName retrieves all location.Map items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLocationMapWithName(name string) (*location.Map, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *location.Map:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type location.Map not found", name)
}

// GetAllLocationPlaceIndexResources retrieves all location.PlaceIndex items from an AWS CloudFormation template
func (t *Template) GetAllLocationPlaceIndexResources() map[string]*location.PlaceIndex {
	results := map[string]*location.PlaceIndex{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *location.PlaceIndex:
			results[name] = resource
		}
	}
	return results
}

// GetLocationPlaceIndexWithName retrieves all location.PlaceIndex items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLocationPlaceIndexWithName(name string) (*location.PlaceIndex, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *location.PlaceIndex:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type location.PlaceIndex not found", name)
}

// GetAllLocationRouteCalculatorResources retrieves all location.RouteCalculator items from an AWS CloudFormation template
func (t *Template) GetAllLocationRouteCalculatorResources() map[string]*location.RouteCalculator {
	results := map[string]*location.RouteCalculator{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *location.RouteCalculator:
			results[name] = resource
		}
	}
	return results
}

// GetLocationRouteCalculatorWithName retrieves all location.RouteCalculator items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLocationRouteCalculatorWithName(name string) (*location.RouteCalculator, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *location.RouteCalculator:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type location.RouteCalculator not found", name)
}

// GetAllLocationTrackerResources retrieves all location.Tracker items from an AWS CloudFormation template
func (t *Template) GetAllLocationTrackerResources() map[string]*location.Tracker {
	results := map[string]*location.Tracker{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *location.Tracker:
			results[name] = resource
		}
	}
	return results
}

// GetLocationTrackerWithName retrieves all location.Tracker items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLocationTrackerWithName(name string) (*location.Tracker, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *location.Tracker:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type location.Tracker not found", name)
}

// GetAllLocationTrackerConsumerResources retrieves all location.TrackerConsumer items from an AWS CloudFormation template
func (t *Template) GetAllLocationTrackerConsumerResources() map[string]*location.TrackerConsumer {
	results := map[string]*location.TrackerConsumer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *location.TrackerConsumer:
			results[name] = resource
		}
	}
	return results
}

// GetLocationTrackerConsumerWithName retrieves all location.TrackerConsumer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLocationTrackerConsumerWithName(name string) (*location.TrackerConsumer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *location.TrackerConsumer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type location.TrackerConsumer not found", name)
}

// GetAllLogsDestinationResources retrieves all logs.Destination items from an AWS CloudFormation template
func (t *Template) GetAllLogsDestinationResources() map[string]*logs.Destination {
	results := map[string]*logs.Destination{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *logs.Destination:
			results[name] = resource
		}
	}
	return results
}

// GetLogsDestinationWithName retrieves all logs.Destination items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLogsDestinationWithName(name string) (*logs.Destination, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *logs.Destination:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type logs.Destination not found", name)
}

// GetAllLogsLogGroupResources retrieves all logs.LogGroup items from an AWS CloudFormation template
func (t *Template) GetAllLogsLogGroupResources() map[string]*logs.LogGroup {
	results := map[string]*logs.LogGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *logs.LogGroup:
			results[name] = resource
		}
	}
	return results
}

// GetLogsLogGroupWithName retrieves all logs.LogGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLogsLogGroupWithName(name string) (*logs.LogGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *logs.LogGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type logs.LogGroup not found", name)
}

// GetAllLogsLogStreamResources retrieves all logs.LogStream items from an AWS CloudFormation template
func (t *Template) GetAllLogsLogStreamResources() map[string]*logs.LogStream {
	results := map[string]*logs.LogStream{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *logs.LogStream:
			results[name] = resource
		}
	}
	return results
}

// GetLogsLogStreamWithName retrieves all logs.LogStream items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLogsLogStreamWithName(name string) (*logs.LogStream, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *logs.LogStream:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type logs.LogStream not found", name)
}

// GetAllLogsMetricFilterResources retrieves all logs.MetricFilter items from an AWS CloudFormation template
func (t *Template) GetAllLogsMetricFilterResources() map[string]*logs.MetricFilter {
	results := map[string]*logs.MetricFilter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *logs.MetricFilter:
			results[name] = resource
		}
	}
	return results
}

// GetLogsMetricFilterWithName retrieves all logs.MetricFilter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLogsMetricFilterWithName(name string) (*logs.MetricFilter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *logs.MetricFilter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type logs.MetricFilter not found", name)
}

// GetAllLogsQueryDefinitionResources retrieves all logs.QueryDefinition items from an AWS CloudFormation template
func (t *Template) GetAllLogsQueryDefinitionResources() map[string]*logs.QueryDefinition {
	results := map[string]*logs.QueryDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *logs.QueryDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetLogsQueryDefinitionWithName retrieves all logs.QueryDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLogsQueryDefinitionWithName(name string) (*logs.QueryDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *logs.QueryDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type logs.QueryDefinition not found", name)
}

// GetAllLogsResourcePolicyResources retrieves all logs.ResourcePolicy items from an AWS CloudFormation template
func (t *Template) GetAllLogsResourcePolicyResources() map[string]*logs.ResourcePolicy {
	results := map[string]*logs.ResourcePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *logs.ResourcePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetLogsResourcePolicyWithName retrieves all logs.ResourcePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLogsResourcePolicyWithName(name string) (*logs.ResourcePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *logs.ResourcePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type logs.ResourcePolicy not found", name)
}

// GetAllLogsSubscriptionFilterResources retrieves all logs.SubscriptionFilter items from an AWS CloudFormation template
func (t *Template) GetAllLogsSubscriptionFilterResources() map[string]*logs.SubscriptionFilter {
	results := map[string]*logs.SubscriptionFilter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *logs.SubscriptionFilter:
			results[name] = resource
		}
	}
	return results
}

// GetLogsSubscriptionFilterWithName retrieves all logs.SubscriptionFilter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLogsSubscriptionFilterWithName(name string) (*logs.SubscriptionFilter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *logs.SubscriptionFilter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type logs.SubscriptionFilter not found", name)
}

// GetAllLookoutEquipmentInferenceSchedulerResources retrieves all lookoutequipment.InferenceScheduler items from an AWS CloudFormation template
func (t *Template) GetAllLookoutEquipmentInferenceSchedulerResources() map[string]*lookoutequipment.InferenceScheduler {
	results := map[string]*lookoutequipment.InferenceScheduler{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lookoutequipment.InferenceScheduler:
			results[name] = resource
		}
	}
	return results
}

// GetLookoutEquipmentInferenceSchedulerWithName retrieves all lookoutequipment.InferenceScheduler items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLookoutEquipmentInferenceSchedulerWithName(name string) (*lookoutequipment.InferenceScheduler, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lookoutequipment.InferenceScheduler:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lookoutequipment.InferenceScheduler not found", name)
}

// GetAllLookoutMetricsAlertResources retrieves all lookoutmetrics.Alert items from an AWS CloudFormation template
func (t *Template) GetAllLookoutMetricsAlertResources() map[string]*lookoutmetrics.Alert {
	results := map[string]*lookoutmetrics.Alert{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lookoutmetrics.Alert:
			results[name] = resource
		}
	}
	return results
}

// GetLookoutMetricsAlertWithName retrieves all lookoutmetrics.Alert items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLookoutMetricsAlertWithName(name string) (*lookoutmetrics.Alert, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lookoutmetrics.Alert:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lookoutmetrics.Alert not found", name)
}

// GetAllLookoutMetricsAnomalyDetectorResources retrieves all lookoutmetrics.AnomalyDetector items from an AWS CloudFormation template
func (t *Template) GetAllLookoutMetricsAnomalyDetectorResources() map[string]*lookoutmetrics.AnomalyDetector {
	results := map[string]*lookoutmetrics.AnomalyDetector{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lookoutmetrics.AnomalyDetector:
			results[name] = resource
		}
	}
	return results
}

// GetLookoutMetricsAnomalyDetectorWithName retrieves all lookoutmetrics.AnomalyDetector items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLookoutMetricsAnomalyDetectorWithName(name string) (*lookoutmetrics.AnomalyDetector, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lookoutmetrics.AnomalyDetector:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lookoutmetrics.AnomalyDetector not found", name)
}

// GetAllLookoutVisionProjectResources retrieves all lookoutvision.Project items from an AWS CloudFormation template
func (t *Template) GetAllLookoutVisionProjectResources() map[string]*lookoutvision.Project {
	results := map[string]*lookoutvision.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lookoutvision.Project:
			results[name] = resource
		}
	}
	return results
}

// GetLookoutVisionProjectWithName retrieves all lookoutvision.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLookoutVisionProjectWithName(name string) (*lookoutvision.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lookoutvision.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lookoutvision.Project not found", name)
}

// GetAllMSKClusterResources retrieves all msk.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllMSKClusterResources() map[string]*msk.Cluster {
	results := map[string]*msk.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *msk.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetMSKClusterWithName retrieves all msk.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMSKClusterWithName(name string) (*msk.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *msk.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type msk.Cluster not found", name)
}

// GetAllMWAAEnvironmentResources retrieves all mwaa.Environment items from an AWS CloudFormation template
func (t *Template) GetAllMWAAEnvironmentResources() map[string]*mwaa.Environment {
	results := map[string]*mwaa.Environment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mwaa.Environment:
			results[name] = resource
		}
	}
	return results
}

// GetMWAAEnvironmentWithName retrieves all mwaa.Environment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMWAAEnvironmentWithName(name string) (*mwaa.Environment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mwaa.Environment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mwaa.Environment not found", name)
}

// GetAllMacieCustomDataIdentifierResources retrieves all macie.CustomDataIdentifier items from an AWS CloudFormation template
func (t *Template) GetAllMacieCustomDataIdentifierResources() map[string]*macie.CustomDataIdentifier {
	results := map[string]*macie.CustomDataIdentifier{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *macie.CustomDataIdentifier:
			results[name] = resource
		}
	}
	return results
}

// GetMacieCustomDataIdentifierWithName retrieves all macie.CustomDataIdentifier items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMacieCustomDataIdentifierWithName(name string) (*macie.CustomDataIdentifier, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *macie.CustomDataIdentifier:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type macie.CustomDataIdentifier not found", name)
}

// GetAllMacieFindingsFilterResources retrieves all macie.FindingsFilter items from an AWS CloudFormation template
func (t *Template) GetAllMacieFindingsFilterResources() map[string]*macie.FindingsFilter {
	results := map[string]*macie.FindingsFilter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *macie.FindingsFilter:
			results[name] = resource
		}
	}
	return results
}

// GetMacieFindingsFilterWithName retrieves all macie.FindingsFilter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMacieFindingsFilterWithName(name string) (*macie.FindingsFilter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *macie.FindingsFilter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type macie.FindingsFilter not found", name)
}

// GetAllMacieSessionResources retrieves all macie.Session items from an AWS CloudFormation template
func (t *Template) GetAllMacieSessionResources() map[string]*macie.Session {
	results := map[string]*macie.Session{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *macie.Session:
			results[name] = resource
		}
	}
	return results
}

// GetMacieSessionWithName retrieves all macie.Session items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMacieSessionWithName(name string) (*macie.Session, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *macie.Session:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type macie.Session not found", name)
}

// GetAllManagedBlockchainMemberResources retrieves all managedblockchain.Member items from an AWS CloudFormation template
func (t *Template) GetAllManagedBlockchainMemberResources() map[string]*managedblockchain.Member {
	results := map[string]*managedblockchain.Member{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *managedblockchain.Member:
			results[name] = resource
		}
	}
	return results
}

// GetManagedBlockchainMemberWithName retrieves all managedblockchain.Member items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetManagedBlockchainMemberWithName(name string) (*managedblockchain.Member, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *managedblockchain.Member:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type managedblockchain.Member not found", name)
}

// GetAllManagedBlockchainNodeResources retrieves all managedblockchain.Node items from an AWS CloudFormation template
func (t *Template) GetAllManagedBlockchainNodeResources() map[string]*managedblockchain.Node {
	results := map[string]*managedblockchain.Node{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *managedblockchain.Node:
			results[name] = resource
		}
	}
	return results
}

// GetManagedBlockchainNodeWithName retrieves all managedblockchain.Node items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetManagedBlockchainNodeWithName(name string) (*managedblockchain.Node, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *managedblockchain.Node:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type managedblockchain.Node not found", name)
}

// GetAllMediaConnectFlowResources retrieves all mediaconnect.Flow items from an AWS CloudFormation template
func (t *Template) GetAllMediaConnectFlowResources() map[string]*mediaconnect.Flow {
	results := map[string]*mediaconnect.Flow{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconnect.Flow:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConnectFlowWithName retrieves all mediaconnect.Flow items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConnectFlowWithName(name string) (*mediaconnect.Flow, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconnect.Flow:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconnect.Flow not found", name)
}

// GetAllMediaConnectFlowEntitlementResources retrieves all mediaconnect.FlowEntitlement items from an AWS CloudFormation template
func (t *Template) GetAllMediaConnectFlowEntitlementResources() map[string]*mediaconnect.FlowEntitlement {
	results := map[string]*mediaconnect.FlowEntitlement{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowEntitlement:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConnectFlowEntitlementWithName retrieves all mediaconnect.FlowEntitlement items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConnectFlowEntitlementWithName(name string) (*mediaconnect.FlowEntitlement, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowEntitlement:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconnect.FlowEntitlement not found", name)
}

// GetAllMediaConnectFlowOutputResources retrieves all mediaconnect.FlowOutput items from an AWS CloudFormation template
func (t *Template) GetAllMediaConnectFlowOutputResources() map[string]*mediaconnect.FlowOutput {
	results := map[string]*mediaconnect.FlowOutput{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowOutput:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConnectFlowOutputWithName retrieves all mediaconnect.FlowOutput items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConnectFlowOutputWithName(name string) (*mediaconnect.FlowOutput, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowOutput:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconnect.FlowOutput not found", name)
}

// GetAllMediaConnectFlowSourceResources retrieves all mediaconnect.FlowSource items from an AWS CloudFormation template
func (t *Template) GetAllMediaConnectFlowSourceResources() map[string]*mediaconnect.FlowSource {
	results := map[string]*mediaconnect.FlowSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowSource:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConnectFlowSourceWithName retrieves all mediaconnect.FlowSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConnectFlowSourceWithName(name string) (*mediaconnect.FlowSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowSource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconnect.FlowSource not found", name)
}

// GetAllMediaConnectFlowVpcInterfaceResources retrieves all mediaconnect.FlowVpcInterface items from an AWS CloudFormation template
func (t *Template) GetAllMediaConnectFlowVpcInterfaceResources() map[string]*mediaconnect.FlowVpcInterface {
	results := map[string]*mediaconnect.FlowVpcInterface{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowVpcInterface:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConnectFlowVpcInterfaceWithName retrieves all mediaconnect.FlowVpcInterface items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConnectFlowVpcInterfaceWithName(name string) (*mediaconnect.FlowVpcInterface, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconnect.FlowVpcInterface:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconnect.FlowVpcInterface not found", name)
}

// GetAllMediaConvertJobTemplateResources retrieves all mediaconvert.JobTemplate items from an AWS CloudFormation template
func (t *Template) GetAllMediaConvertJobTemplateResources() map[string]*mediaconvert.JobTemplate {
	results := map[string]*mediaconvert.JobTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconvert.JobTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConvertJobTemplateWithName retrieves all mediaconvert.JobTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConvertJobTemplateWithName(name string) (*mediaconvert.JobTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconvert.JobTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconvert.JobTemplate not found", name)
}

// GetAllMediaConvertPresetResources retrieves all mediaconvert.Preset items from an AWS CloudFormation template
func (t *Template) GetAllMediaConvertPresetResources() map[string]*mediaconvert.Preset {
	results := map[string]*mediaconvert.Preset{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconvert.Preset:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConvertPresetWithName retrieves all mediaconvert.Preset items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConvertPresetWithName(name string) (*mediaconvert.Preset, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconvert.Preset:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconvert.Preset not found", name)
}

// GetAllMediaConvertQueueResources retrieves all mediaconvert.Queue items from an AWS CloudFormation template
func (t *Template) GetAllMediaConvertQueueResources() map[string]*mediaconvert.Queue {
	results := map[string]*mediaconvert.Queue{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediaconvert.Queue:
			results[name] = resource
		}
	}
	return results
}

// GetMediaConvertQueueWithName retrieves all mediaconvert.Queue items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaConvertQueueWithName(name string) (*mediaconvert.Queue, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediaconvert.Queue:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediaconvert.Queue not found", name)
}

// GetAllMediaLiveChannelResources retrieves all medialive.Channel items from an AWS CloudFormation template
func (t *Template) GetAllMediaLiveChannelResources() map[string]*medialive.Channel {
	results := map[string]*medialive.Channel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *medialive.Channel:
			results[name] = resource
		}
	}
	return results
}

// GetMediaLiveChannelWithName retrieves all medialive.Channel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaLiveChannelWithName(name string) (*medialive.Channel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *medialive.Channel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type medialive.Channel not found", name)
}

// GetAllMediaLiveInputResources retrieves all medialive.Input items from an AWS CloudFormation template
func (t *Template) GetAllMediaLiveInputResources() map[string]*medialive.Input {
	results := map[string]*medialive.Input{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *medialive.Input:
			results[name] = resource
		}
	}
	return results
}

// GetMediaLiveInputWithName retrieves all medialive.Input items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaLiveInputWithName(name string) (*medialive.Input, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *medialive.Input:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type medialive.Input not found", name)
}

// GetAllMediaLiveInputSecurityGroupResources retrieves all medialive.InputSecurityGroup items from an AWS CloudFormation template
func (t *Template) GetAllMediaLiveInputSecurityGroupResources() map[string]*medialive.InputSecurityGroup {
	results := map[string]*medialive.InputSecurityGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *medialive.InputSecurityGroup:
			results[name] = resource
		}
	}
	return results
}

// GetMediaLiveInputSecurityGroupWithName retrieves all medialive.InputSecurityGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaLiveInputSecurityGroupWithName(name string) (*medialive.InputSecurityGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *medialive.InputSecurityGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type medialive.InputSecurityGroup not found", name)
}

// GetAllMediaPackageAssetResources retrieves all mediapackage.Asset items from an AWS CloudFormation template
func (t *Template) GetAllMediaPackageAssetResources() map[string]*mediapackage.Asset {
	results := map[string]*mediapackage.Asset{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediapackage.Asset:
			results[name] = resource
		}
	}
	return results
}

// GetMediaPackageAssetWithName retrieves all mediapackage.Asset items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaPackageAssetWithName(name string) (*mediapackage.Asset, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediapackage.Asset:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediapackage.Asset not found", name)
}

// GetAllMediaPackageChannelResources retrieves all mediapackage.Channel items from an AWS CloudFormation template
func (t *Template) GetAllMediaPackageChannelResources() map[string]*mediapackage.Channel {
	results := map[string]*mediapackage.Channel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediapackage.Channel:
			results[name] = resource
		}
	}
	return results
}

// GetMediaPackageChannelWithName retrieves all mediapackage.Channel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaPackageChannelWithName(name string) (*mediapackage.Channel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediapackage.Channel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediapackage.Channel not found", name)
}

// GetAllMediaPackageOriginEndpointResources retrieves all mediapackage.OriginEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllMediaPackageOriginEndpointResources() map[string]*mediapackage.OriginEndpoint {
	results := map[string]*mediapackage.OriginEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediapackage.OriginEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetMediaPackageOriginEndpointWithName retrieves all mediapackage.OriginEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaPackageOriginEndpointWithName(name string) (*mediapackage.OriginEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediapackage.OriginEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediapackage.OriginEndpoint not found", name)
}

// GetAllMediaPackagePackagingConfigurationResources retrieves all mediapackage.PackagingConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllMediaPackagePackagingConfigurationResources() map[string]*mediapackage.PackagingConfiguration {
	results := map[string]*mediapackage.PackagingConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediapackage.PackagingConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetMediaPackagePackagingConfigurationWithName retrieves all mediapackage.PackagingConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaPackagePackagingConfigurationWithName(name string) (*mediapackage.PackagingConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediapackage.PackagingConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediapackage.PackagingConfiguration not found", name)
}

// GetAllMediaPackagePackagingGroupResources retrieves all mediapackage.PackagingGroup items from an AWS CloudFormation template
func (t *Template) GetAllMediaPackagePackagingGroupResources() map[string]*mediapackage.PackagingGroup {
	results := map[string]*mediapackage.PackagingGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediapackage.PackagingGroup:
			results[name] = resource
		}
	}
	return results
}

// GetMediaPackagePackagingGroupWithName retrieves all mediapackage.PackagingGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaPackagePackagingGroupWithName(name string) (*mediapackage.PackagingGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediapackage.PackagingGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediapackage.PackagingGroup not found", name)
}

// GetAllMediaStoreContainerResources retrieves all mediastore.Container items from an AWS CloudFormation template
func (t *Template) GetAllMediaStoreContainerResources() map[string]*mediastore.Container {
	results := map[string]*mediastore.Container{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *mediastore.Container:
			results[name] = resource
		}
	}
	return results
}

// GetMediaStoreContainerWithName retrieves all mediastore.Container items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMediaStoreContainerWithName(name string) (*mediastore.Container, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *mediastore.Container:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type mediastore.Container not found", name)
}

// GetAllMemoryDBACLResources retrieves all memorydb.ACL items from an AWS CloudFormation template
func (t *Template) GetAllMemoryDBACLResources() map[string]*memorydb.ACL {
	results := map[string]*memorydb.ACL{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *memorydb.ACL:
			results[name] = resource
		}
	}
	return results
}

// GetMemoryDBACLWithName retrieves all memorydb.ACL items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMemoryDBACLWithName(name string) (*memorydb.ACL, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *memorydb.ACL:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type memorydb.ACL not found", name)
}

// GetAllMemoryDBClusterResources retrieves all memorydb.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllMemoryDBClusterResources() map[string]*memorydb.Cluster {
	results := map[string]*memorydb.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *memorydb.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetMemoryDBClusterWithName retrieves all memorydb.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMemoryDBClusterWithName(name string) (*memorydb.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *memorydb.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type memorydb.Cluster not found", name)
}

// GetAllMemoryDBParameterGroupResources retrieves all memorydb.ParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllMemoryDBParameterGroupResources() map[string]*memorydb.ParameterGroup {
	results := map[string]*memorydb.ParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *memorydb.ParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetMemoryDBParameterGroupWithName retrieves all memorydb.ParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMemoryDBParameterGroupWithName(name string) (*memorydb.ParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *memorydb.ParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type memorydb.ParameterGroup not found", name)
}

// GetAllMemoryDBSubnetGroupResources retrieves all memorydb.SubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllMemoryDBSubnetGroupResources() map[string]*memorydb.SubnetGroup {
	results := map[string]*memorydb.SubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *memorydb.SubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetMemoryDBSubnetGroupWithName retrieves all memorydb.SubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMemoryDBSubnetGroupWithName(name string) (*memorydb.SubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *memorydb.SubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type memorydb.SubnetGroup not found", name)
}

// GetAllMemoryDBUserResources retrieves all memorydb.User items from an AWS CloudFormation template
func (t *Template) GetAllMemoryDBUserResources() map[string]*memorydb.User {
	results := map[string]*memorydb.User{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *memorydb.User:
			results[name] = resource
		}
	}
	return results
}

// GetMemoryDBUserWithName retrieves all memorydb.User items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetMemoryDBUserWithName(name string) (*memorydb.User, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *memorydb.User:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type memorydb.User not found", name)
}

// GetAllNeptuneDBClusterResources retrieves all neptune.DBCluster items from an AWS CloudFormation template
func (t *Template) GetAllNeptuneDBClusterResources() map[string]*neptune.DBCluster {
	results := map[string]*neptune.DBCluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *neptune.DBCluster:
			results[name] = resource
		}
	}
	return results
}

// GetNeptuneDBClusterWithName retrieves all neptune.DBCluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNeptuneDBClusterWithName(name string) (*neptune.DBCluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *neptune.DBCluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type neptune.DBCluster not found", name)
}

// GetAllNeptuneDBClusterParameterGroupResources retrieves all neptune.DBClusterParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllNeptuneDBClusterParameterGroupResources() map[string]*neptune.DBClusterParameterGroup {
	results := map[string]*neptune.DBClusterParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *neptune.DBClusterParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetNeptuneDBClusterParameterGroupWithName retrieves all neptune.DBClusterParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNeptuneDBClusterParameterGroupWithName(name string) (*neptune.DBClusterParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *neptune.DBClusterParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type neptune.DBClusterParameterGroup not found", name)
}

// GetAllNeptuneDBInstanceResources retrieves all neptune.DBInstance items from an AWS CloudFormation template
func (t *Template) GetAllNeptuneDBInstanceResources() map[string]*neptune.DBInstance {
	results := map[string]*neptune.DBInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *neptune.DBInstance:
			results[name] = resource
		}
	}
	return results
}

// GetNeptuneDBInstanceWithName retrieves all neptune.DBInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNeptuneDBInstanceWithName(name string) (*neptune.DBInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *neptune.DBInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type neptune.DBInstance not found", name)
}

// GetAllNeptuneDBParameterGroupResources retrieves all neptune.DBParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllNeptuneDBParameterGroupResources() map[string]*neptune.DBParameterGroup {
	results := map[string]*neptune.DBParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *neptune.DBParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetNeptuneDBParameterGroupWithName retrieves all neptune.DBParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNeptuneDBParameterGroupWithName(name string) (*neptune.DBParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *neptune.DBParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type neptune.DBParameterGroup not found", name)
}

// GetAllNeptuneDBSubnetGroupResources retrieves all neptune.DBSubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllNeptuneDBSubnetGroupResources() map[string]*neptune.DBSubnetGroup {
	results := map[string]*neptune.DBSubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *neptune.DBSubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetNeptuneDBSubnetGroupWithName retrieves all neptune.DBSubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNeptuneDBSubnetGroupWithName(name string) (*neptune.DBSubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *neptune.DBSubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type neptune.DBSubnetGroup not found", name)
}

// GetAllNetworkFirewallFirewallResources retrieves all networkfirewall.Firewall items from an AWS CloudFormation template
func (t *Template) GetAllNetworkFirewallFirewallResources() map[string]*networkfirewall.Firewall {
	results := map[string]*networkfirewall.Firewall{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkfirewall.Firewall:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkFirewallFirewallWithName retrieves all networkfirewall.Firewall items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkFirewallFirewallWithName(name string) (*networkfirewall.Firewall, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkfirewall.Firewall:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkfirewall.Firewall not found", name)
}

// GetAllNetworkFirewallFirewallPolicyResources retrieves all networkfirewall.FirewallPolicy items from an AWS CloudFormation template
func (t *Template) GetAllNetworkFirewallFirewallPolicyResources() map[string]*networkfirewall.FirewallPolicy {
	results := map[string]*networkfirewall.FirewallPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkfirewall.FirewallPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkFirewallFirewallPolicyWithName retrieves all networkfirewall.FirewallPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkFirewallFirewallPolicyWithName(name string) (*networkfirewall.FirewallPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkfirewall.FirewallPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkfirewall.FirewallPolicy not found", name)
}

// GetAllNetworkFirewallLoggingConfigurationResources retrieves all networkfirewall.LoggingConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllNetworkFirewallLoggingConfigurationResources() map[string]*networkfirewall.LoggingConfiguration {
	results := map[string]*networkfirewall.LoggingConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkfirewall.LoggingConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkFirewallLoggingConfigurationWithName retrieves all networkfirewall.LoggingConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkFirewallLoggingConfigurationWithName(name string) (*networkfirewall.LoggingConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkfirewall.LoggingConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkfirewall.LoggingConfiguration not found", name)
}

// GetAllNetworkFirewallRuleGroupResources retrieves all networkfirewall.RuleGroup items from an AWS CloudFormation template
func (t *Template) GetAllNetworkFirewallRuleGroupResources() map[string]*networkfirewall.RuleGroup {
	results := map[string]*networkfirewall.RuleGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkfirewall.RuleGroup:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkFirewallRuleGroupWithName retrieves all networkfirewall.RuleGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkFirewallRuleGroupWithName(name string) (*networkfirewall.RuleGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkfirewall.RuleGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkfirewall.RuleGroup not found", name)
}

// GetAllNetworkManagerCustomerGatewayAssociationResources retrieves all networkmanager.CustomerGatewayAssociation items from an AWS CloudFormation template
func (t *Template) GetAllNetworkManagerCustomerGatewayAssociationResources() map[string]*networkmanager.CustomerGatewayAssociation {
	results := map[string]*networkmanager.CustomerGatewayAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkmanager.CustomerGatewayAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkManagerCustomerGatewayAssociationWithName retrieves all networkmanager.CustomerGatewayAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkManagerCustomerGatewayAssociationWithName(name string) (*networkmanager.CustomerGatewayAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkmanager.CustomerGatewayAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkmanager.CustomerGatewayAssociation not found", name)
}

// GetAllNetworkManagerDeviceResources retrieves all networkmanager.Device items from an AWS CloudFormation template
func (t *Template) GetAllNetworkManagerDeviceResources() map[string]*networkmanager.Device {
	results := map[string]*networkmanager.Device{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkmanager.Device:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkManagerDeviceWithName retrieves all networkmanager.Device items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkManagerDeviceWithName(name string) (*networkmanager.Device, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkmanager.Device:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkmanager.Device not found", name)
}

// GetAllNetworkManagerGlobalNetworkResources retrieves all networkmanager.GlobalNetwork items from an AWS CloudFormation template
func (t *Template) GetAllNetworkManagerGlobalNetworkResources() map[string]*networkmanager.GlobalNetwork {
	results := map[string]*networkmanager.GlobalNetwork{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkmanager.GlobalNetwork:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkManagerGlobalNetworkWithName retrieves all networkmanager.GlobalNetwork items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkManagerGlobalNetworkWithName(name string) (*networkmanager.GlobalNetwork, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkmanager.GlobalNetwork:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkmanager.GlobalNetwork not found", name)
}

// GetAllNetworkManagerLinkResources retrieves all networkmanager.Link items from an AWS CloudFormation template
func (t *Template) GetAllNetworkManagerLinkResources() map[string]*networkmanager.Link {
	results := map[string]*networkmanager.Link{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkmanager.Link:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkManagerLinkWithName retrieves all networkmanager.Link items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkManagerLinkWithName(name string) (*networkmanager.Link, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkmanager.Link:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkmanager.Link not found", name)
}

// GetAllNetworkManagerLinkAssociationResources retrieves all networkmanager.LinkAssociation items from an AWS CloudFormation template
func (t *Template) GetAllNetworkManagerLinkAssociationResources() map[string]*networkmanager.LinkAssociation {
	results := map[string]*networkmanager.LinkAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkmanager.LinkAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkManagerLinkAssociationWithName retrieves all networkmanager.LinkAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkManagerLinkAssociationWithName(name string) (*networkmanager.LinkAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkmanager.LinkAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkmanager.LinkAssociation not found", name)
}

// GetAllNetworkManagerSiteResources retrieves all networkmanager.Site items from an AWS CloudFormation template
func (t *Template) GetAllNetworkManagerSiteResources() map[string]*networkmanager.Site {
	results := map[string]*networkmanager.Site{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkmanager.Site:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkManagerSiteWithName retrieves all networkmanager.Site items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkManagerSiteWithName(name string) (*networkmanager.Site, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkmanager.Site:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkmanager.Site not found", name)
}

// GetAllNetworkManagerTransitGatewayRegistrationResources retrieves all networkmanager.TransitGatewayRegistration items from an AWS CloudFormation template
func (t *Template) GetAllNetworkManagerTransitGatewayRegistrationResources() map[string]*networkmanager.TransitGatewayRegistration {
	results := map[string]*networkmanager.TransitGatewayRegistration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *networkmanager.TransitGatewayRegistration:
			results[name] = resource
		}
	}
	return results
}

// GetNetworkManagerTransitGatewayRegistrationWithName retrieves all networkmanager.TransitGatewayRegistration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNetworkManagerTransitGatewayRegistrationWithName(name string) (*networkmanager.TransitGatewayRegistration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *networkmanager.TransitGatewayRegistration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type networkmanager.TransitGatewayRegistration not found", name)
}

// GetAllNimbleStudioLaunchProfileResources retrieves all nimblestudio.LaunchProfile items from an AWS CloudFormation template
func (t *Template) GetAllNimbleStudioLaunchProfileResources() map[string]*nimblestudio.LaunchProfile {
	results := map[string]*nimblestudio.LaunchProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *nimblestudio.LaunchProfile:
			results[name] = resource
		}
	}
	return results
}

// GetNimbleStudioLaunchProfileWithName retrieves all nimblestudio.LaunchProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNimbleStudioLaunchProfileWithName(name string) (*nimblestudio.LaunchProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *nimblestudio.LaunchProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type nimblestudio.LaunchProfile not found", name)
}

// GetAllNimbleStudioStreamingImageResources retrieves all nimblestudio.StreamingImage items from an AWS CloudFormation template
func (t *Template) GetAllNimbleStudioStreamingImageResources() map[string]*nimblestudio.StreamingImage {
	results := map[string]*nimblestudio.StreamingImage{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *nimblestudio.StreamingImage:
			results[name] = resource
		}
	}
	return results
}

// GetNimbleStudioStreamingImageWithName retrieves all nimblestudio.StreamingImage items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNimbleStudioStreamingImageWithName(name string) (*nimblestudio.StreamingImage, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *nimblestudio.StreamingImage:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type nimblestudio.StreamingImage not found", name)
}

// GetAllNimbleStudioStudioResources retrieves all nimblestudio.Studio items from an AWS CloudFormation template
func (t *Template) GetAllNimbleStudioStudioResources() map[string]*nimblestudio.Studio {
	results := map[string]*nimblestudio.Studio{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *nimblestudio.Studio:
			results[name] = resource
		}
	}
	return results
}

// GetNimbleStudioStudioWithName retrieves all nimblestudio.Studio items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNimbleStudioStudioWithName(name string) (*nimblestudio.Studio, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *nimblestudio.Studio:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type nimblestudio.Studio not found", name)
}

// GetAllNimbleStudioStudioComponentResources retrieves all nimblestudio.StudioComponent items from an AWS CloudFormation template
func (t *Template) GetAllNimbleStudioStudioComponentResources() map[string]*nimblestudio.StudioComponent {
	results := map[string]*nimblestudio.StudioComponent{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *nimblestudio.StudioComponent:
			results[name] = resource
		}
	}
	return results
}

// GetNimbleStudioStudioComponentWithName retrieves all nimblestudio.StudioComponent items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetNimbleStudioStudioComponentWithName(name string) (*nimblestudio.StudioComponent, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *nimblestudio.StudioComponent:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type nimblestudio.StudioComponent not found", name)
}

// GetAllOpenSearchServiceDomainResources retrieves all opensearchservice.Domain items from an AWS CloudFormation template
func (t *Template) GetAllOpenSearchServiceDomainResources() map[string]*opensearchservice.Domain {
	results := map[string]*opensearchservice.Domain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opensearchservice.Domain:
			results[name] = resource
		}
	}
	return results
}

// GetOpenSearchServiceDomainWithName retrieves all opensearchservice.Domain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpenSearchServiceDomainWithName(name string) (*opensearchservice.Domain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opensearchservice.Domain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opensearchservice.Domain not found", name)
}

// GetAllOpsWorksAppResources retrieves all opsworks.App items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksAppResources() map[string]*opsworks.App {
	results := map[string]*opsworks.App{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworks.App:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksAppWithName retrieves all opsworks.App items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksAppWithName(name string) (*opsworks.App, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworks.App:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworks.App not found", name)
}

// GetAllOpsWorksElasticLoadBalancerAttachmentResources retrieves all opsworks.ElasticLoadBalancerAttachment items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksElasticLoadBalancerAttachmentResources() map[string]*opsworks.ElasticLoadBalancerAttachment {
	results := map[string]*opsworks.ElasticLoadBalancerAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworks.ElasticLoadBalancerAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksElasticLoadBalancerAttachmentWithName retrieves all opsworks.ElasticLoadBalancerAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksElasticLoadBalancerAttachmentWithName(name string) (*opsworks.ElasticLoadBalancerAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworks.ElasticLoadBalancerAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworks.ElasticLoadBalancerAttachment not found", name)
}

// GetAllOpsWorksInstanceResources retrieves all opsworks.Instance items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksInstanceResources() map[string]*opsworks.Instance {
	results := map[string]*opsworks.Instance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworks.Instance:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksInstanceWithName retrieves all opsworks.Instance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksInstanceWithName(name string) (*opsworks.Instance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworks.Instance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworks.Instance not found", name)
}

// GetAllOpsWorksLayerResources retrieves all opsworks.Layer items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksLayerResources() map[string]*opsworks.Layer {
	results := map[string]*opsworks.Layer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworks.Layer:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksLayerWithName retrieves all opsworks.Layer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksLayerWithName(name string) (*opsworks.Layer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworks.Layer:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworks.Layer not found", name)
}

// GetAllOpsWorksStackResources retrieves all opsworks.Stack items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksStackResources() map[string]*opsworks.Stack {
	results := map[string]*opsworks.Stack{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworks.Stack:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksStackWithName retrieves all opsworks.Stack items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksStackWithName(name string) (*opsworks.Stack, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworks.Stack:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworks.Stack not found", name)
}

// GetAllOpsWorksUserProfileResources retrieves all opsworks.UserProfile items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksUserProfileResources() map[string]*opsworks.UserProfile {
	results := map[string]*opsworks.UserProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworks.UserProfile:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksUserProfileWithName retrieves all opsworks.UserProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksUserProfileWithName(name string) (*opsworks.UserProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworks.UserProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworks.UserProfile not found", name)
}

// GetAllOpsWorksVolumeResources retrieves all opsworks.Volume items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksVolumeResources() map[string]*opsworks.Volume {
	results := map[string]*opsworks.Volume{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworks.Volume:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksVolumeWithName retrieves all opsworks.Volume items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksVolumeWithName(name string) (*opsworks.Volume, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworks.Volume:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworks.Volume not found", name)
}

// GetAllOpsWorksCMServerResources retrieves all opsworkscm.Server items from an AWS CloudFormation template
func (t *Template) GetAllOpsWorksCMServerResources() map[string]*opsworkscm.Server {
	results := map[string]*opsworkscm.Server{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *opsworkscm.Server:
			results[name] = resource
		}
	}
	return results
}

// GetOpsWorksCMServerWithName retrieves all opsworkscm.Server items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetOpsWorksCMServerWithName(name string) (*opsworkscm.Server, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *opsworkscm.Server:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type opsworkscm.Server not found", name)
}

// GetAllPanoramaApplicationInstanceResources retrieves all panorama.ApplicationInstance items from an AWS CloudFormation template
func (t *Template) GetAllPanoramaApplicationInstanceResources() map[string]*panorama.ApplicationInstance {
	results := map[string]*panorama.ApplicationInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *panorama.ApplicationInstance:
			results[name] = resource
		}
	}
	return results
}

// GetPanoramaApplicationInstanceWithName retrieves all panorama.ApplicationInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPanoramaApplicationInstanceWithName(name string) (*panorama.ApplicationInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *panorama.ApplicationInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type panorama.ApplicationInstance not found", name)
}

// GetAllPanoramaPackageResources retrieves all panorama.Package items from an AWS CloudFormation template
func (t *Template) GetAllPanoramaPackageResources() map[string]*panorama.Package {
	results := map[string]*panorama.Package{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *panorama.Package:
			results[name] = resource
		}
	}
	return results
}

// GetPanoramaPackageWithName retrieves all panorama.Package items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPanoramaPackageWithName(name string) (*panorama.Package, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *panorama.Package:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type panorama.Package not found", name)
}

// GetAllPanoramaPackageVersionResources retrieves all panorama.PackageVersion items from an AWS CloudFormation template
func (t *Template) GetAllPanoramaPackageVersionResources() map[string]*panorama.PackageVersion {
	results := map[string]*panorama.PackageVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *panorama.PackageVersion:
			results[name] = resource
		}
	}
	return results
}

// GetPanoramaPackageVersionWithName retrieves all panorama.PackageVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPanoramaPackageVersionWithName(name string) (*panorama.PackageVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *panorama.PackageVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type panorama.PackageVersion not found", name)
}

// GetAllPinpointADMChannelResources retrieves all pinpoint.ADMChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointADMChannelResources() map[string]*pinpoint.ADMChannel {
	results := map[string]*pinpoint.ADMChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.ADMChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointADMChannelWithName retrieves all pinpoint.ADMChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointADMChannelWithName(name string) (*pinpoint.ADMChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.ADMChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.ADMChannel not found", name)
}

// GetAllPinpointAPNSChannelResources retrieves all pinpoint.APNSChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointAPNSChannelResources() map[string]*pinpoint.APNSChannel {
	results := map[string]*pinpoint.APNSChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.APNSChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointAPNSChannelWithName retrieves all pinpoint.APNSChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointAPNSChannelWithName(name string) (*pinpoint.APNSChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.APNSChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.APNSChannel not found", name)
}

// GetAllPinpointAPNSSandboxChannelResources retrieves all pinpoint.APNSSandboxChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointAPNSSandboxChannelResources() map[string]*pinpoint.APNSSandboxChannel {
	results := map[string]*pinpoint.APNSSandboxChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.APNSSandboxChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointAPNSSandboxChannelWithName retrieves all pinpoint.APNSSandboxChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointAPNSSandboxChannelWithName(name string) (*pinpoint.APNSSandboxChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.APNSSandboxChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.APNSSandboxChannel not found", name)
}

// GetAllPinpointAPNSVoipChannelResources retrieves all pinpoint.APNSVoipChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointAPNSVoipChannelResources() map[string]*pinpoint.APNSVoipChannel {
	results := map[string]*pinpoint.APNSVoipChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.APNSVoipChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointAPNSVoipChannelWithName retrieves all pinpoint.APNSVoipChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointAPNSVoipChannelWithName(name string) (*pinpoint.APNSVoipChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.APNSVoipChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.APNSVoipChannel not found", name)
}

// GetAllPinpointAPNSVoipSandboxChannelResources retrieves all pinpoint.APNSVoipSandboxChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointAPNSVoipSandboxChannelResources() map[string]*pinpoint.APNSVoipSandboxChannel {
	results := map[string]*pinpoint.APNSVoipSandboxChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.APNSVoipSandboxChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointAPNSVoipSandboxChannelWithName retrieves all pinpoint.APNSVoipSandboxChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointAPNSVoipSandboxChannelWithName(name string) (*pinpoint.APNSVoipSandboxChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.APNSVoipSandboxChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.APNSVoipSandboxChannel not found", name)
}

// GetAllPinpointAppResources retrieves all pinpoint.App items from an AWS CloudFormation template
func (t *Template) GetAllPinpointAppResources() map[string]*pinpoint.App {
	results := map[string]*pinpoint.App{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.App:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointAppWithName retrieves all pinpoint.App items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointAppWithName(name string) (*pinpoint.App, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.App:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.App not found", name)
}

// GetAllPinpointApplicationSettingsResources retrieves all pinpoint.ApplicationSettings items from an AWS CloudFormation template
func (t *Template) GetAllPinpointApplicationSettingsResources() map[string]*pinpoint.ApplicationSettings {
	results := map[string]*pinpoint.ApplicationSettings{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.ApplicationSettings:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointApplicationSettingsWithName retrieves all pinpoint.ApplicationSettings items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointApplicationSettingsWithName(name string) (*pinpoint.ApplicationSettings, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.ApplicationSettings:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.ApplicationSettings not found", name)
}

// GetAllPinpointBaiduChannelResources retrieves all pinpoint.BaiduChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointBaiduChannelResources() map[string]*pinpoint.BaiduChannel {
	results := map[string]*pinpoint.BaiduChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.BaiduChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointBaiduChannelWithName retrieves all pinpoint.BaiduChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointBaiduChannelWithName(name string) (*pinpoint.BaiduChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.BaiduChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.BaiduChannel not found", name)
}

// GetAllPinpointCampaignResources retrieves all pinpoint.Campaign items from an AWS CloudFormation template
func (t *Template) GetAllPinpointCampaignResources() map[string]*pinpoint.Campaign {
	results := map[string]*pinpoint.Campaign{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.Campaign:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointCampaignWithName retrieves all pinpoint.Campaign items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointCampaignWithName(name string) (*pinpoint.Campaign, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.Campaign:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.Campaign not found", name)
}

// GetAllPinpointEmailChannelResources retrieves all pinpoint.EmailChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointEmailChannelResources() map[string]*pinpoint.EmailChannel {
	results := map[string]*pinpoint.EmailChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.EmailChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointEmailChannelWithName retrieves all pinpoint.EmailChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointEmailChannelWithName(name string) (*pinpoint.EmailChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.EmailChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.EmailChannel not found", name)
}

// GetAllPinpointEmailTemplateResources retrieves all pinpoint.EmailTemplate items from an AWS CloudFormation template
func (t *Template) GetAllPinpointEmailTemplateResources() map[string]*pinpoint.EmailTemplate {
	results := map[string]*pinpoint.EmailTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.EmailTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointEmailTemplateWithName retrieves all pinpoint.EmailTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointEmailTemplateWithName(name string) (*pinpoint.EmailTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.EmailTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.EmailTemplate not found", name)
}

// GetAllPinpointEventStreamResources retrieves all pinpoint.EventStream items from an AWS CloudFormation template
func (t *Template) GetAllPinpointEventStreamResources() map[string]*pinpoint.EventStream {
	results := map[string]*pinpoint.EventStream{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.EventStream:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointEventStreamWithName retrieves all pinpoint.EventStream items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointEventStreamWithName(name string) (*pinpoint.EventStream, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.EventStream:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.EventStream not found", name)
}

// GetAllPinpointGCMChannelResources retrieves all pinpoint.GCMChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointGCMChannelResources() map[string]*pinpoint.GCMChannel {
	results := map[string]*pinpoint.GCMChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.GCMChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointGCMChannelWithName retrieves all pinpoint.GCMChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointGCMChannelWithName(name string) (*pinpoint.GCMChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.GCMChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.GCMChannel not found", name)
}

// GetAllPinpointInAppTemplateResources retrieves all pinpoint.InAppTemplate items from an AWS CloudFormation template
func (t *Template) GetAllPinpointInAppTemplateResources() map[string]*pinpoint.InAppTemplate {
	results := map[string]*pinpoint.InAppTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.InAppTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointInAppTemplateWithName retrieves all pinpoint.InAppTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointInAppTemplateWithName(name string) (*pinpoint.InAppTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.InAppTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.InAppTemplate not found", name)
}

// GetAllPinpointPushTemplateResources retrieves all pinpoint.PushTemplate items from an AWS CloudFormation template
func (t *Template) GetAllPinpointPushTemplateResources() map[string]*pinpoint.PushTemplate {
	results := map[string]*pinpoint.PushTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.PushTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointPushTemplateWithName retrieves all pinpoint.PushTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointPushTemplateWithName(name string) (*pinpoint.PushTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.PushTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.PushTemplate not found", name)
}

// GetAllPinpointSMSChannelResources retrieves all pinpoint.SMSChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointSMSChannelResources() map[string]*pinpoint.SMSChannel {
	results := map[string]*pinpoint.SMSChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.SMSChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointSMSChannelWithName retrieves all pinpoint.SMSChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointSMSChannelWithName(name string) (*pinpoint.SMSChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.SMSChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.SMSChannel not found", name)
}

// GetAllPinpointSegmentResources retrieves all pinpoint.Segment items from an AWS CloudFormation template
func (t *Template) GetAllPinpointSegmentResources() map[string]*pinpoint.Segment {
	results := map[string]*pinpoint.Segment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.Segment:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointSegmentWithName retrieves all pinpoint.Segment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointSegmentWithName(name string) (*pinpoint.Segment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.Segment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.Segment not found", name)
}

// GetAllPinpointSmsTemplateResources retrieves all pinpoint.SmsTemplate items from an AWS CloudFormation template
func (t *Template) GetAllPinpointSmsTemplateResources() map[string]*pinpoint.SmsTemplate {
	results := map[string]*pinpoint.SmsTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.SmsTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointSmsTemplateWithName retrieves all pinpoint.SmsTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointSmsTemplateWithName(name string) (*pinpoint.SmsTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.SmsTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.SmsTemplate not found", name)
}

// GetAllPinpointVoiceChannelResources retrieves all pinpoint.VoiceChannel items from an AWS CloudFormation template
func (t *Template) GetAllPinpointVoiceChannelResources() map[string]*pinpoint.VoiceChannel {
	results := map[string]*pinpoint.VoiceChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpoint.VoiceChannel:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointVoiceChannelWithName retrieves all pinpoint.VoiceChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointVoiceChannelWithName(name string) (*pinpoint.VoiceChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpoint.VoiceChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpoint.VoiceChannel not found", name)
}

// GetAllPinpointEmailConfigurationSetResources retrieves all pinpointemail.ConfigurationSet items from an AWS CloudFormation template
func (t *Template) GetAllPinpointEmailConfigurationSetResources() map[string]*pinpointemail.ConfigurationSet {
	results := map[string]*pinpointemail.ConfigurationSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpointemail.ConfigurationSet:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointEmailConfigurationSetWithName retrieves all pinpointemail.ConfigurationSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointEmailConfigurationSetWithName(name string) (*pinpointemail.ConfigurationSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpointemail.ConfigurationSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpointemail.ConfigurationSet not found", name)
}

// GetAllPinpointEmailConfigurationSetEventDestinationResources retrieves all pinpointemail.ConfigurationSetEventDestination items from an AWS CloudFormation template
func (t *Template) GetAllPinpointEmailConfigurationSetEventDestinationResources() map[string]*pinpointemail.ConfigurationSetEventDestination {
	results := map[string]*pinpointemail.ConfigurationSetEventDestination{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpointemail.ConfigurationSetEventDestination:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointEmailConfigurationSetEventDestinationWithName retrieves all pinpointemail.ConfigurationSetEventDestination items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointEmailConfigurationSetEventDestinationWithName(name string) (*pinpointemail.ConfigurationSetEventDestination, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpointemail.ConfigurationSetEventDestination:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpointemail.ConfigurationSetEventDestination not found", name)
}

// GetAllPinpointEmailDedicatedIpPoolResources retrieves all pinpointemail.DedicatedIpPool items from an AWS CloudFormation template
func (t *Template) GetAllPinpointEmailDedicatedIpPoolResources() map[string]*pinpointemail.DedicatedIpPool {
	results := map[string]*pinpointemail.DedicatedIpPool{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpointemail.DedicatedIpPool:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointEmailDedicatedIpPoolWithName retrieves all pinpointemail.DedicatedIpPool items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointEmailDedicatedIpPoolWithName(name string) (*pinpointemail.DedicatedIpPool, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpointemail.DedicatedIpPool:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpointemail.DedicatedIpPool not found", name)
}

// GetAllPinpointEmailIdentityResources retrieves all pinpointemail.Identity items from an AWS CloudFormation template
func (t *Template) GetAllPinpointEmailIdentityResources() map[string]*pinpointemail.Identity {
	results := map[string]*pinpointemail.Identity{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *pinpointemail.Identity:
			results[name] = resource
		}
	}
	return results
}

// GetPinpointEmailIdentityWithName retrieves all pinpointemail.Identity items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetPinpointEmailIdentityWithName(name string) (*pinpointemail.Identity, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *pinpointemail.Identity:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type pinpointemail.Identity not found", name)
}

// GetAllQLDBLedgerResources retrieves all qldb.Ledger items from an AWS CloudFormation template
func (t *Template) GetAllQLDBLedgerResources() map[string]*qldb.Ledger {
	results := map[string]*qldb.Ledger{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *qldb.Ledger:
			results[name] = resource
		}
	}
	return results
}

// GetQLDBLedgerWithName retrieves all qldb.Ledger items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQLDBLedgerWithName(name string) (*qldb.Ledger, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *qldb.Ledger:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type qldb.Ledger not found", name)
}

// GetAllQLDBStreamResources retrieves all qldb.Stream items from an AWS CloudFormation template
func (t *Template) GetAllQLDBStreamResources() map[string]*qldb.Stream {
	results := map[string]*qldb.Stream{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *qldb.Stream:
			results[name] = resource
		}
	}
	return results
}

// GetQLDBStreamWithName retrieves all qldb.Stream items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQLDBStreamWithName(name string) (*qldb.Stream, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *qldb.Stream:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type qldb.Stream not found", name)
}

// GetAllQuickSightAnalysisResources retrieves all quicksight.Analysis items from an AWS CloudFormation template
func (t *Template) GetAllQuickSightAnalysisResources() map[string]*quicksight.Analysis {
	results := map[string]*quicksight.Analysis{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *quicksight.Analysis:
			results[name] = resource
		}
	}
	return results
}

// GetQuickSightAnalysisWithName retrieves all quicksight.Analysis items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQuickSightAnalysisWithName(name string) (*quicksight.Analysis, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *quicksight.Analysis:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type quicksight.Analysis not found", name)
}

// GetAllQuickSightDashboardResources retrieves all quicksight.Dashboard items from an AWS CloudFormation template
func (t *Template) GetAllQuickSightDashboardResources() map[string]*quicksight.Dashboard {
	results := map[string]*quicksight.Dashboard{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *quicksight.Dashboard:
			results[name] = resource
		}
	}
	return results
}

// GetQuickSightDashboardWithName retrieves all quicksight.Dashboard items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQuickSightDashboardWithName(name string) (*quicksight.Dashboard, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *quicksight.Dashboard:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type quicksight.Dashboard not found", name)
}

// GetAllQuickSightDataSetResources retrieves all quicksight.DataSet items from an AWS CloudFormation template
func (t *Template) GetAllQuickSightDataSetResources() map[string]*quicksight.DataSet {
	results := map[string]*quicksight.DataSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *quicksight.DataSet:
			results[name] = resource
		}
	}
	return results
}

// GetQuickSightDataSetWithName retrieves all quicksight.DataSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQuickSightDataSetWithName(name string) (*quicksight.DataSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *quicksight.DataSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type quicksight.DataSet not found", name)
}

// GetAllQuickSightDataSourceResources retrieves all quicksight.DataSource items from an AWS CloudFormation template
func (t *Template) GetAllQuickSightDataSourceResources() map[string]*quicksight.DataSource {
	results := map[string]*quicksight.DataSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *quicksight.DataSource:
			results[name] = resource
		}
	}
	return results
}

// GetQuickSightDataSourceWithName retrieves all quicksight.DataSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQuickSightDataSourceWithName(name string) (*quicksight.DataSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *quicksight.DataSource:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type quicksight.DataSource not found", name)
}

// GetAllQuickSightTemplateResources retrieves all quicksight.Template items from an AWS CloudFormation template
func (t *Template) GetAllQuickSightTemplateResources() map[string]*quicksight.Template {
	results := map[string]*quicksight.Template{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *quicksight.Template:
			results[name] = resource
		}
	}
	return results
}

// GetQuickSightTemplateWithName retrieves all quicksight.Template items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQuickSightTemplateWithName(name string) (*quicksight.Template, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *quicksight.Template:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type quicksight.Template not found", name)
}

// GetAllQuickSightThemeResources retrieves all quicksight.Theme items from an AWS CloudFormation template
func (t *Template) GetAllQuickSightThemeResources() map[string]*quicksight.Theme {
	results := map[string]*quicksight.Theme{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *quicksight.Theme:
			results[name] = resource
		}
	}
	return results
}

// GetQuickSightThemeWithName retrieves all quicksight.Theme items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetQuickSightThemeWithName(name string) (*quicksight.Theme, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *quicksight.Theme:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type quicksight.Theme not found", name)
}

// GetAllRAMResourceShareResources retrieves all ram.ResourceShare items from an AWS CloudFormation template
func (t *Template) GetAllRAMResourceShareResources() map[string]*ram.ResourceShare {
	results := map[string]*ram.ResourceShare{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ram.ResourceShare:
			results[name] = resource
		}
	}
	return results
}

// GetRAMResourceShareWithName retrieves all ram.ResourceShare items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRAMResourceShareWithName(name string) (*ram.ResourceShare, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ram.ResourceShare:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ram.ResourceShare not found", name)
}

// GetAllRDSDBClusterResources retrieves all rds.DBCluster items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBClusterResources() map[string]*rds.DBCluster {
	results := map[string]*rds.DBCluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBCluster:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBClusterWithName retrieves all rds.DBCluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBClusterWithName(name string) (*rds.DBCluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBCluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBCluster not found", name)
}

// GetAllRDSDBClusterParameterGroupResources retrieves all rds.DBClusterParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBClusterParameterGroupResources() map[string]*rds.DBClusterParameterGroup {
	results := map[string]*rds.DBClusterParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBClusterParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBClusterParameterGroupWithName retrieves all rds.DBClusterParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBClusterParameterGroupWithName(name string) (*rds.DBClusterParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBClusterParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBClusterParameterGroup not found", name)
}

// GetAllRDSDBInstanceResources retrieves all rds.DBInstance items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBInstanceResources() map[string]*rds.DBInstance {
	results := map[string]*rds.DBInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBInstance:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBInstanceWithName retrieves all rds.DBInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBInstanceWithName(name string) (*rds.DBInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBInstance not found", name)
}

// GetAllRDSDBParameterGroupResources retrieves all rds.DBParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBParameterGroupResources() map[string]*rds.DBParameterGroup {
	results := map[string]*rds.DBParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBParameterGroupWithName retrieves all rds.DBParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBParameterGroupWithName(name string) (*rds.DBParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBParameterGroup not found", name)
}

// GetAllRDSDBProxyResources retrieves all rds.DBProxy items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBProxyResources() map[string]*rds.DBProxy {
	results := map[string]*rds.DBProxy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBProxy:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBProxyWithName retrieves all rds.DBProxy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBProxyWithName(name string) (*rds.DBProxy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBProxy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBProxy not found", name)
}

// GetAllRDSDBProxyEndpointResources retrieves all rds.DBProxyEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBProxyEndpointResources() map[string]*rds.DBProxyEndpoint {
	results := map[string]*rds.DBProxyEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBProxyEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBProxyEndpointWithName retrieves all rds.DBProxyEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBProxyEndpointWithName(name string) (*rds.DBProxyEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBProxyEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBProxyEndpoint not found", name)
}

// GetAllRDSDBProxyTargetGroupResources retrieves all rds.DBProxyTargetGroup items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBProxyTargetGroupResources() map[string]*rds.DBProxyTargetGroup {
	results := map[string]*rds.DBProxyTargetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBProxyTargetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBProxyTargetGroupWithName retrieves all rds.DBProxyTargetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBProxyTargetGroupWithName(name string) (*rds.DBProxyTargetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBProxyTargetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBProxyTargetGroup not found", name)
}

// GetAllRDSDBSecurityGroupResources retrieves all rds.DBSecurityGroup items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBSecurityGroupResources() map[string]*rds.DBSecurityGroup {
	results := map[string]*rds.DBSecurityGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBSecurityGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBSecurityGroupWithName retrieves all rds.DBSecurityGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBSecurityGroupWithName(name string) (*rds.DBSecurityGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBSecurityGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBSecurityGroup not found", name)
}

// GetAllRDSDBSecurityGroupIngressResources retrieves all rds.DBSecurityGroupIngress items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBSecurityGroupIngressResources() map[string]*rds.DBSecurityGroupIngress {
	results := map[string]*rds.DBSecurityGroupIngress{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBSecurityGroupIngress:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBSecurityGroupIngressWithName retrieves all rds.DBSecurityGroupIngress items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBSecurityGroupIngressWithName(name string) (*rds.DBSecurityGroupIngress, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBSecurityGroupIngress:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBSecurityGroupIngress not found", name)
}

// GetAllRDSDBSubnetGroupResources retrieves all rds.DBSubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBSubnetGroupResources() map[string]*rds.DBSubnetGroup {
	results := map[string]*rds.DBSubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBSubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBSubnetGroupWithName retrieves all rds.DBSubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBSubnetGroupWithName(name string) (*rds.DBSubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBSubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBSubnetGroup not found", name)
}

// GetAllRDSEventSubscriptionResources retrieves all rds.EventSubscription items from an AWS CloudFormation template
func (t *Template) GetAllRDSEventSubscriptionResources() map[string]*rds.EventSubscription {
	results := map[string]*rds.EventSubscription{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.EventSubscription:
			results[name] = resource
		}
	}
	return results
}

// GetRDSEventSubscriptionWithName retrieves all rds.EventSubscription items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSEventSubscriptionWithName(name string) (*rds.EventSubscription, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.EventSubscription:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.EventSubscription not found", name)
}

// GetAllRDSGlobalClusterResources retrieves all rds.GlobalCluster items from an AWS CloudFormation template
func (t *Template) GetAllRDSGlobalClusterResources() map[string]*rds.GlobalCluster {
	results := map[string]*rds.GlobalCluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.GlobalCluster:
			results[name] = resource
		}
	}
	return results
}

// GetRDSGlobalClusterWithName retrieves all rds.GlobalCluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSGlobalClusterWithName(name string) (*rds.GlobalCluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.GlobalCluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.GlobalCluster not found", name)
}

// GetAllRDSOptionGroupResources retrieves all rds.OptionGroup items from an AWS CloudFormation template
func (t *Template) GetAllRDSOptionGroupResources() map[string]*rds.OptionGroup {
	results := map[string]*rds.OptionGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.OptionGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRDSOptionGroupWithName retrieves all rds.OptionGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSOptionGroupWithName(name string) (*rds.OptionGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.OptionGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.OptionGroup not found", name)
}

// GetAllRUMAppMonitorResources retrieves all rum.AppMonitor items from an AWS CloudFormation template
func (t *Template) GetAllRUMAppMonitorResources() map[string]*rum.AppMonitor {
	results := map[string]*rum.AppMonitor{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rum.AppMonitor:
			results[name] = resource
		}
	}
	return results
}

// GetRUMAppMonitorWithName retrieves all rum.AppMonitor items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRUMAppMonitorWithName(name string) (*rum.AppMonitor, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rum.AppMonitor:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rum.AppMonitor not found", name)
}

// GetAllRedshiftClusterResources retrieves all redshift.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftClusterResources() map[string]*redshift.Cluster {
	results := map[string]*redshift.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftClusterWithName retrieves all redshift.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftClusterWithName(name string) (*redshift.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.Cluster not found", name)
}

// GetAllRedshiftClusterParameterGroupResources retrieves all redshift.ClusterParameterGroup items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftClusterParameterGroupResources() map[string]*redshift.ClusterParameterGroup {
	results := map[string]*redshift.ClusterParameterGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.ClusterParameterGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftClusterParameterGroupWithName retrieves all redshift.ClusterParameterGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftClusterParameterGroupWithName(name string) (*redshift.ClusterParameterGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.ClusterParameterGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.ClusterParameterGroup not found", name)
}

// GetAllRedshiftClusterSecurityGroupResources retrieves all redshift.ClusterSecurityGroup items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftClusterSecurityGroupResources() map[string]*redshift.ClusterSecurityGroup {
	results := map[string]*redshift.ClusterSecurityGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.ClusterSecurityGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftClusterSecurityGroupWithName retrieves all redshift.ClusterSecurityGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftClusterSecurityGroupWithName(name string) (*redshift.ClusterSecurityGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.ClusterSecurityGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.ClusterSecurityGroup not found", name)
}

// GetAllRedshiftClusterSecurityGroupIngressResources retrieves all redshift.ClusterSecurityGroupIngress items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftClusterSecurityGroupIngressResources() map[string]*redshift.ClusterSecurityGroupIngress {
	results := map[string]*redshift.ClusterSecurityGroupIngress{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.ClusterSecurityGroupIngress:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftClusterSecurityGroupIngressWithName retrieves all redshift.ClusterSecurityGroupIngress items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftClusterSecurityGroupIngressWithName(name string) (*redshift.ClusterSecurityGroupIngress, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.ClusterSecurityGroupIngress:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.ClusterSecurityGroupIngress not found", name)
}

// GetAllRedshiftClusterSubnetGroupResources retrieves all redshift.ClusterSubnetGroup items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftClusterSubnetGroupResources() map[string]*redshift.ClusterSubnetGroup {
	results := map[string]*redshift.ClusterSubnetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.ClusterSubnetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftClusterSubnetGroupWithName retrieves all redshift.ClusterSubnetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftClusterSubnetGroupWithName(name string) (*redshift.ClusterSubnetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.ClusterSubnetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.ClusterSubnetGroup not found", name)
}

// GetAllRedshiftEndpointAccessResources retrieves all redshift.EndpointAccess items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftEndpointAccessResources() map[string]*redshift.EndpointAccess {
	results := map[string]*redshift.EndpointAccess{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.EndpointAccess:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftEndpointAccessWithName retrieves all redshift.EndpointAccess items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftEndpointAccessWithName(name string) (*redshift.EndpointAccess, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.EndpointAccess:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.EndpointAccess not found", name)
}

// GetAllRedshiftEndpointAuthorizationResources retrieves all redshift.EndpointAuthorization items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftEndpointAuthorizationResources() map[string]*redshift.EndpointAuthorization {
	results := map[string]*redshift.EndpointAuthorization{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.EndpointAuthorization:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftEndpointAuthorizationWithName retrieves all redshift.EndpointAuthorization items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftEndpointAuthorizationWithName(name string) (*redshift.EndpointAuthorization, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.EndpointAuthorization:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.EndpointAuthorization not found", name)
}

// GetAllRedshiftEventSubscriptionResources retrieves all redshift.EventSubscription items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftEventSubscriptionResources() map[string]*redshift.EventSubscription {
	results := map[string]*redshift.EventSubscription{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.EventSubscription:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftEventSubscriptionWithName retrieves all redshift.EventSubscription items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftEventSubscriptionWithName(name string) (*redshift.EventSubscription, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.EventSubscription:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.EventSubscription not found", name)
}

// GetAllRedshiftScheduledActionResources retrieves all redshift.ScheduledAction items from an AWS CloudFormation template
func (t *Template) GetAllRedshiftScheduledActionResources() map[string]*redshift.ScheduledAction {
	results := map[string]*redshift.ScheduledAction{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *redshift.ScheduledAction:
			results[name] = resource
		}
	}
	return results
}

// GetRedshiftScheduledActionWithName retrieves all redshift.ScheduledAction items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRedshiftScheduledActionWithName(name string) (*redshift.ScheduledAction, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *redshift.ScheduledAction:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type redshift.ScheduledAction not found", name)
}

// GetAllRefactorSpacesApplicationResources retrieves all refactorspaces.Application items from an AWS CloudFormation template
func (t *Template) GetAllRefactorSpacesApplicationResources() map[string]*refactorspaces.Application {
	results := map[string]*refactorspaces.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *refactorspaces.Application:
			results[name] = resource
		}
	}
	return results
}

// GetRefactorSpacesApplicationWithName retrieves all refactorspaces.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRefactorSpacesApplicationWithName(name string) (*refactorspaces.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *refactorspaces.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type refactorspaces.Application not found", name)
}

// GetAllRefactorSpacesEnvironmentResources retrieves all refactorspaces.Environment items from an AWS CloudFormation template
func (t *Template) GetAllRefactorSpacesEnvironmentResources() map[string]*refactorspaces.Environment {
	results := map[string]*refactorspaces.Environment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *refactorspaces.Environment:
			results[name] = resource
		}
	}
	return results
}

// GetRefactorSpacesEnvironmentWithName retrieves all refactorspaces.Environment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRefactorSpacesEnvironmentWithName(name string) (*refactorspaces.Environment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *refactorspaces.Environment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type refactorspaces.Environment not found", name)
}

// GetAllRefactorSpacesRouteResources retrieves all refactorspaces.Route items from an AWS CloudFormation template
func (t *Template) GetAllRefactorSpacesRouteResources() map[string]*refactorspaces.Route {
	results := map[string]*refactorspaces.Route{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *refactorspaces.Route:
			results[name] = resource
		}
	}
	return results
}

// GetRefactorSpacesRouteWithName retrieves all refactorspaces.Route items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRefactorSpacesRouteWithName(name string) (*refactorspaces.Route, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *refactorspaces.Route:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type refactorspaces.Route not found", name)
}

// GetAllRefactorSpacesServiceResources retrieves all refactorspaces.Service items from an AWS CloudFormation template
func (t *Template) GetAllRefactorSpacesServiceResources() map[string]*refactorspaces.Service {
	results := map[string]*refactorspaces.Service{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *refactorspaces.Service:
			results[name] = resource
		}
	}
	return results
}

// GetRefactorSpacesServiceWithName retrieves all refactorspaces.Service items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRefactorSpacesServiceWithName(name string) (*refactorspaces.Service, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *refactorspaces.Service:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type refactorspaces.Service not found", name)
}

// GetAllRekognitionProjectResources retrieves all rekognition.Project items from an AWS CloudFormation template
func (t *Template) GetAllRekognitionProjectResources() map[string]*rekognition.Project {
	results := map[string]*rekognition.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rekognition.Project:
			results[name] = resource
		}
	}
	return results
}

// GetRekognitionProjectWithName retrieves all rekognition.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRekognitionProjectWithName(name string) (*rekognition.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rekognition.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rekognition.Project not found", name)
}

// GetAllResilienceHubAppResources retrieves all resiliencehub.App items from an AWS CloudFormation template
func (t *Template) GetAllResilienceHubAppResources() map[string]*resiliencehub.App {
	results := map[string]*resiliencehub.App{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *resiliencehub.App:
			results[name] = resource
		}
	}
	return results
}

// GetResilienceHubAppWithName retrieves all resiliencehub.App items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetResilienceHubAppWithName(name string) (*resiliencehub.App, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *resiliencehub.App:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type resiliencehub.App not found", name)
}

// GetAllResilienceHubResiliencyPolicyResources retrieves all resiliencehub.ResiliencyPolicy items from an AWS CloudFormation template
func (t *Template) GetAllResilienceHubResiliencyPolicyResources() map[string]*resiliencehub.ResiliencyPolicy {
	results := map[string]*resiliencehub.ResiliencyPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *resiliencehub.ResiliencyPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetResilienceHubResiliencyPolicyWithName retrieves all resiliencehub.ResiliencyPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetResilienceHubResiliencyPolicyWithName(name string) (*resiliencehub.ResiliencyPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *resiliencehub.ResiliencyPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type resiliencehub.ResiliencyPolicy not found", name)
}

// GetAllResourceGroupsGroupResources retrieves all resourcegroups.Group items from an AWS CloudFormation template
func (t *Template) GetAllResourceGroupsGroupResources() map[string]*resourcegroups.Group {
	results := map[string]*resourcegroups.Group{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *resourcegroups.Group:
			results[name] = resource
		}
	}
	return results
}

// GetResourceGroupsGroupWithName retrieves all resourcegroups.Group items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetResourceGroupsGroupWithName(name string) (*resourcegroups.Group, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *resourcegroups.Group:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type resourcegroups.Group not found", name)
}

// GetAllRoboMakerFleetResources retrieves all robomaker.Fleet items from an AWS CloudFormation template
func (t *Template) GetAllRoboMakerFleetResources() map[string]*robomaker.Fleet {
	results := map[string]*robomaker.Fleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *robomaker.Fleet:
			results[name] = resource
		}
	}
	return results
}

// GetRoboMakerFleetWithName retrieves all robomaker.Fleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoboMakerFleetWithName(name string) (*robomaker.Fleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *robomaker.Fleet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type robomaker.Fleet not found", name)
}

// GetAllRoboMakerRobotResources retrieves all robomaker.Robot items from an AWS CloudFormation template
func (t *Template) GetAllRoboMakerRobotResources() map[string]*robomaker.Robot {
	results := map[string]*robomaker.Robot{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *robomaker.Robot:
			results[name] = resource
		}
	}
	return results
}

// GetRoboMakerRobotWithName retrieves all robomaker.Robot items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoboMakerRobotWithName(name string) (*robomaker.Robot, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *robomaker.Robot:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type robomaker.Robot not found", name)
}

// GetAllRoboMakerRobotApplicationResources retrieves all robomaker.RobotApplication items from an AWS CloudFormation template
func (t *Template) GetAllRoboMakerRobotApplicationResources() map[string]*robomaker.RobotApplication {
	results := map[string]*robomaker.RobotApplication{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *robomaker.RobotApplication:
			results[name] = resource
		}
	}
	return results
}

// GetRoboMakerRobotApplicationWithName retrieves all robomaker.RobotApplication items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoboMakerRobotApplicationWithName(name string) (*robomaker.RobotApplication, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *robomaker.RobotApplication:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type robomaker.RobotApplication not found", name)
}

// GetAllRoboMakerRobotApplicationVersionResources retrieves all robomaker.RobotApplicationVersion items from an AWS CloudFormation template
func (t *Template) GetAllRoboMakerRobotApplicationVersionResources() map[string]*robomaker.RobotApplicationVersion {
	results := map[string]*robomaker.RobotApplicationVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *robomaker.RobotApplicationVersion:
			results[name] = resource
		}
	}
	return results
}

// GetRoboMakerRobotApplicationVersionWithName retrieves all robomaker.RobotApplicationVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoboMakerRobotApplicationVersionWithName(name string) (*robomaker.RobotApplicationVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *robomaker.RobotApplicationVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type robomaker.RobotApplicationVersion not found", name)
}

// GetAllRoboMakerSimulationApplicationResources retrieves all robomaker.SimulationApplication items from an AWS CloudFormation template
func (t *Template) GetAllRoboMakerSimulationApplicationResources() map[string]*robomaker.SimulationApplication {
	results := map[string]*robomaker.SimulationApplication{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *robomaker.SimulationApplication:
			results[name] = resource
		}
	}
	return results
}

// GetRoboMakerSimulationApplicationWithName retrieves all robomaker.SimulationApplication items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoboMakerSimulationApplicationWithName(name string) (*robomaker.SimulationApplication, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *robomaker.SimulationApplication:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type robomaker.SimulationApplication not found", name)
}

// GetAllRoboMakerSimulationApplicationVersionResources retrieves all robomaker.SimulationApplicationVersion items from an AWS CloudFormation template
func (t *Template) GetAllRoboMakerSimulationApplicationVersionResources() map[string]*robomaker.SimulationApplicationVersion {
	results := map[string]*robomaker.SimulationApplicationVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *robomaker.SimulationApplicationVersion:
			results[name] = resource
		}
	}
	return results
}

// GetRoboMakerSimulationApplicationVersionWithName retrieves all robomaker.SimulationApplicationVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoboMakerSimulationApplicationVersionWithName(name string) (*robomaker.SimulationApplicationVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *robomaker.SimulationApplicationVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type robomaker.SimulationApplicationVersion not found", name)
}

// GetAllRoute53DNSSECResources retrieves all route53.DNSSEC items from an AWS CloudFormation template
func (t *Template) GetAllRoute53DNSSECResources() map[string]*route53.DNSSEC {
	results := map[string]*route53.DNSSEC{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53.DNSSEC:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53DNSSECWithName retrieves all route53.DNSSEC items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53DNSSECWithName(name string) (*route53.DNSSEC, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53.DNSSEC:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53.DNSSEC not found", name)
}

// GetAllRoute53HealthCheckResources retrieves all route53.HealthCheck items from an AWS CloudFormation template
func (t *Template) GetAllRoute53HealthCheckResources() map[string]*route53.HealthCheck {
	results := map[string]*route53.HealthCheck{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53.HealthCheck:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53HealthCheckWithName retrieves all route53.HealthCheck items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53HealthCheckWithName(name string) (*route53.HealthCheck, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53.HealthCheck:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53.HealthCheck not found", name)
}

// GetAllRoute53HostedZoneResources retrieves all route53.HostedZone items from an AWS CloudFormation template
func (t *Template) GetAllRoute53HostedZoneResources() map[string]*route53.HostedZone {
	results := map[string]*route53.HostedZone{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53.HostedZone:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53HostedZoneWithName retrieves all route53.HostedZone items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53HostedZoneWithName(name string) (*route53.HostedZone, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53.HostedZone:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53.HostedZone not found", name)
}

// GetAllRoute53KeySigningKeyResources retrieves all route53.KeySigningKey items from an AWS CloudFormation template
func (t *Template) GetAllRoute53KeySigningKeyResources() map[string]*route53.KeySigningKey {
	results := map[string]*route53.KeySigningKey{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53.KeySigningKey:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53KeySigningKeyWithName retrieves all route53.KeySigningKey items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53KeySigningKeyWithName(name string) (*route53.KeySigningKey, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53.KeySigningKey:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53.KeySigningKey not found", name)
}

// GetAllRoute53RecordSetResources retrieves all route53.RecordSet items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecordSetResources() map[string]*route53.RecordSet {
	results := map[string]*route53.RecordSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53.RecordSet:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecordSetWithName retrieves all route53.RecordSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecordSetWithName(name string) (*route53.RecordSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53.RecordSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53.RecordSet not found", name)
}

// GetAllRoute53RecordSetGroupResources retrieves all route53.RecordSetGroup items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecordSetGroupResources() map[string]*route53.RecordSetGroup {
	results := map[string]*route53.RecordSetGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53.RecordSetGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecordSetGroupWithName retrieves all route53.RecordSetGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecordSetGroupWithName(name string) (*route53.RecordSetGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53.RecordSetGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53.RecordSetGroup not found", name)
}

// GetAllRoute53RecoveryControlClusterResources retrieves all route53recoverycontrol.Cluster items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryControlClusterResources() map[string]*route53recoverycontrol.Cluster {
	results := map[string]*route53recoverycontrol.Cluster{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.Cluster:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryControlClusterWithName retrieves all route53recoverycontrol.Cluster items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryControlClusterWithName(name string) (*route53recoverycontrol.Cluster, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.Cluster:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoverycontrol.Cluster not found", name)
}

// GetAllRoute53RecoveryControlControlPanelResources retrieves all route53recoverycontrol.ControlPanel items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryControlControlPanelResources() map[string]*route53recoverycontrol.ControlPanel {
	results := map[string]*route53recoverycontrol.ControlPanel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.ControlPanel:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryControlControlPanelWithName retrieves all route53recoverycontrol.ControlPanel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryControlControlPanelWithName(name string) (*route53recoverycontrol.ControlPanel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.ControlPanel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoverycontrol.ControlPanel not found", name)
}

// GetAllRoute53RecoveryControlRoutingControlResources retrieves all route53recoverycontrol.RoutingControl items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryControlRoutingControlResources() map[string]*route53recoverycontrol.RoutingControl {
	results := map[string]*route53recoverycontrol.RoutingControl{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.RoutingControl:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryControlRoutingControlWithName retrieves all route53recoverycontrol.RoutingControl items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryControlRoutingControlWithName(name string) (*route53recoverycontrol.RoutingControl, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.RoutingControl:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoverycontrol.RoutingControl not found", name)
}

// GetAllRoute53RecoveryControlSafetyRuleResources retrieves all route53recoverycontrol.SafetyRule items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryControlSafetyRuleResources() map[string]*route53recoverycontrol.SafetyRule {
	results := map[string]*route53recoverycontrol.SafetyRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.SafetyRule:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryControlSafetyRuleWithName retrieves all route53recoverycontrol.SafetyRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryControlSafetyRuleWithName(name string) (*route53recoverycontrol.SafetyRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoverycontrol.SafetyRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoverycontrol.SafetyRule not found", name)
}

// GetAllRoute53RecoveryReadinessCellResources retrieves all route53recoveryreadiness.Cell items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryReadinessCellResources() map[string]*route53recoveryreadiness.Cell {
	results := map[string]*route53recoveryreadiness.Cell{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.Cell:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryReadinessCellWithName retrieves all route53recoveryreadiness.Cell items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryReadinessCellWithName(name string) (*route53recoveryreadiness.Cell, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.Cell:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoveryreadiness.Cell not found", name)
}

// GetAllRoute53RecoveryReadinessReadinessCheckResources retrieves all route53recoveryreadiness.ReadinessCheck items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryReadinessReadinessCheckResources() map[string]*route53recoveryreadiness.ReadinessCheck {
	results := map[string]*route53recoveryreadiness.ReadinessCheck{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.ReadinessCheck:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryReadinessReadinessCheckWithName retrieves all route53recoveryreadiness.ReadinessCheck items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryReadinessReadinessCheckWithName(name string) (*route53recoveryreadiness.ReadinessCheck, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.ReadinessCheck:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoveryreadiness.ReadinessCheck not found", name)
}

// GetAllRoute53RecoveryReadinessRecoveryGroupResources retrieves all route53recoveryreadiness.RecoveryGroup items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryReadinessRecoveryGroupResources() map[string]*route53recoveryreadiness.RecoveryGroup {
	results := map[string]*route53recoveryreadiness.RecoveryGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.RecoveryGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryReadinessRecoveryGroupWithName retrieves all route53recoveryreadiness.RecoveryGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryReadinessRecoveryGroupWithName(name string) (*route53recoveryreadiness.RecoveryGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.RecoveryGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoveryreadiness.RecoveryGroup not found", name)
}

// GetAllRoute53RecoveryReadinessResourceSetResources retrieves all route53recoveryreadiness.ResourceSet items from an AWS CloudFormation template
func (t *Template) GetAllRoute53RecoveryReadinessResourceSetResources() map[string]*route53recoveryreadiness.ResourceSet {
	results := map[string]*route53recoveryreadiness.ResourceSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.ResourceSet:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53RecoveryReadinessResourceSetWithName retrieves all route53recoveryreadiness.ResourceSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53RecoveryReadinessResourceSetWithName(name string) (*route53recoveryreadiness.ResourceSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53recoveryreadiness.ResourceSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53recoveryreadiness.ResourceSet not found", name)
}

// GetAllRoute53ResolverFirewallDomainListResources retrieves all route53resolver.FirewallDomainList items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverFirewallDomainListResources() map[string]*route53resolver.FirewallDomainList {
	results := map[string]*route53resolver.FirewallDomainList{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.FirewallDomainList:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverFirewallDomainListWithName retrieves all route53resolver.FirewallDomainList items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverFirewallDomainListWithName(name string) (*route53resolver.FirewallDomainList, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.FirewallDomainList:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.FirewallDomainList not found", name)
}

// GetAllRoute53ResolverFirewallRuleGroupResources retrieves all route53resolver.FirewallRuleGroup items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverFirewallRuleGroupResources() map[string]*route53resolver.FirewallRuleGroup {
	results := map[string]*route53resolver.FirewallRuleGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.FirewallRuleGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverFirewallRuleGroupWithName retrieves all route53resolver.FirewallRuleGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverFirewallRuleGroupWithName(name string) (*route53resolver.FirewallRuleGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.FirewallRuleGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.FirewallRuleGroup not found", name)
}

// GetAllRoute53ResolverFirewallRuleGroupAssociationResources retrieves all route53resolver.FirewallRuleGroupAssociation items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverFirewallRuleGroupAssociationResources() map[string]*route53resolver.FirewallRuleGroupAssociation {
	results := map[string]*route53resolver.FirewallRuleGroupAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.FirewallRuleGroupAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverFirewallRuleGroupAssociationWithName retrieves all route53resolver.FirewallRuleGroupAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverFirewallRuleGroupAssociationWithName(name string) (*route53resolver.FirewallRuleGroupAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.FirewallRuleGroupAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.FirewallRuleGroupAssociation not found", name)
}

// GetAllRoute53ResolverResolverConfigResources retrieves all route53resolver.ResolverConfig items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverResolverConfigResources() map[string]*route53resolver.ResolverConfig {
	results := map[string]*route53resolver.ResolverConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverConfig:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverResolverConfigWithName retrieves all route53resolver.ResolverConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverResolverConfigWithName(name string) (*route53resolver.ResolverConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.ResolverConfig not found", name)
}

// GetAllRoute53ResolverResolverDNSSECConfigResources retrieves all route53resolver.ResolverDNSSECConfig items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverResolverDNSSECConfigResources() map[string]*route53resolver.ResolverDNSSECConfig {
	results := map[string]*route53resolver.ResolverDNSSECConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverDNSSECConfig:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverResolverDNSSECConfigWithName retrieves all route53resolver.ResolverDNSSECConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverResolverDNSSECConfigWithName(name string) (*route53resolver.ResolverDNSSECConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverDNSSECConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.ResolverDNSSECConfig not found", name)
}

// GetAllRoute53ResolverResolverEndpointResources retrieves all route53resolver.ResolverEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverResolverEndpointResources() map[string]*route53resolver.ResolverEndpoint {
	results := map[string]*route53resolver.ResolverEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverResolverEndpointWithName retrieves all route53resolver.ResolverEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverResolverEndpointWithName(name string) (*route53resolver.ResolverEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.ResolverEndpoint not found", name)
}

// GetAllRoute53ResolverResolverQueryLoggingConfigResources retrieves all route53resolver.ResolverQueryLoggingConfig items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverResolverQueryLoggingConfigResources() map[string]*route53resolver.ResolverQueryLoggingConfig {
	results := map[string]*route53resolver.ResolverQueryLoggingConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverQueryLoggingConfig:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverResolverQueryLoggingConfigWithName retrieves all route53resolver.ResolverQueryLoggingConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverResolverQueryLoggingConfigWithName(name string) (*route53resolver.ResolverQueryLoggingConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverQueryLoggingConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.ResolverQueryLoggingConfig not found", name)
}

// GetAllRoute53ResolverResolverQueryLoggingConfigAssociationResources retrieves all route53resolver.ResolverQueryLoggingConfigAssociation items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverResolverQueryLoggingConfigAssociationResources() map[string]*route53resolver.ResolverQueryLoggingConfigAssociation {
	results := map[string]*route53resolver.ResolverQueryLoggingConfigAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverQueryLoggingConfigAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverResolverQueryLoggingConfigAssociationWithName retrieves all route53resolver.ResolverQueryLoggingConfigAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverResolverQueryLoggingConfigAssociationWithName(name string) (*route53resolver.ResolverQueryLoggingConfigAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverQueryLoggingConfigAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.ResolverQueryLoggingConfigAssociation not found", name)
}

// GetAllRoute53ResolverResolverRuleResources retrieves all route53resolver.ResolverRule items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverResolverRuleResources() map[string]*route53resolver.ResolverRule {
	results := map[string]*route53resolver.ResolverRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverRule:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverResolverRuleWithName retrieves all route53resolver.ResolverRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverResolverRuleWithName(name string) (*route53resolver.ResolverRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.ResolverRule not found", name)
}

// GetAllRoute53ResolverResolverRuleAssociationResources retrieves all route53resolver.ResolverRuleAssociation items from an AWS CloudFormation template
func (t *Template) GetAllRoute53ResolverResolverRuleAssociationResources() map[string]*route53resolver.ResolverRuleAssociation {
	results := map[string]*route53resolver.ResolverRuleAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverRuleAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53ResolverResolverRuleAssociationWithName retrieves all route53resolver.ResolverRuleAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53ResolverResolverRuleAssociationWithName(name string) (*route53resolver.ResolverRuleAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53resolver.ResolverRuleAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53resolver.ResolverRuleAssociation not found", name)
}

// GetAllS3AccessPointResources retrieves all s3.AccessPoint items from an AWS CloudFormation template
func (t *Template) GetAllS3AccessPointResources() map[string]*s3.AccessPoint {
	results := map[string]*s3.AccessPoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.AccessPoint:
			results[name] = resource
		}
	}
	return results
}

// GetS3AccessPointWithName retrieves all s3.AccessPoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3AccessPointWithName(name string) (*s3.AccessPoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.AccessPoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.AccessPoint not found", name)
}

// GetAllS3BucketResources retrieves all s3.Bucket items from an AWS CloudFormation template
func (t *Template) GetAllS3BucketResources() map[string]*s3.Bucket {
	results := map[string]*s3.Bucket{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.Bucket:
			results[name] = resource
		}
	}
	return results
}

// GetS3BucketWithName retrieves all s3.Bucket items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3BucketWithName(name string) (*s3.Bucket, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.Bucket:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.Bucket not found", name)
}

// GetAllS3BucketPolicyResources retrieves all s3.BucketPolicy items from an AWS CloudFormation template
func (t *Template) GetAllS3BucketPolicyResources() map[string]*s3.BucketPolicy {
	results := map[string]*s3.BucketPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.BucketPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetS3BucketPolicyWithName retrieves all s3.BucketPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3BucketPolicyWithName(name string) (*s3.BucketPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.BucketPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.BucketPolicy not found", name)
}

// GetAllS3MultiRegionAccessPointResources retrieves all s3.MultiRegionAccessPoint items from an AWS CloudFormation template
func (t *Template) GetAllS3MultiRegionAccessPointResources() map[string]*s3.MultiRegionAccessPoint {
	results := map[string]*s3.MultiRegionAccessPoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.MultiRegionAccessPoint:
			results[name] = resource
		}
	}
	return results
}

// GetS3MultiRegionAccessPointWithName retrieves all s3.MultiRegionAccessPoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3MultiRegionAccessPointWithName(name string) (*s3.MultiRegionAccessPoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.MultiRegionAccessPoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.MultiRegionAccessPoint not found", name)
}

// GetAllS3MultiRegionAccessPointPolicyResources retrieves all s3.MultiRegionAccessPointPolicy items from an AWS CloudFormation template
func (t *Template) GetAllS3MultiRegionAccessPointPolicyResources() map[string]*s3.MultiRegionAccessPointPolicy {
	results := map[string]*s3.MultiRegionAccessPointPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.MultiRegionAccessPointPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetS3MultiRegionAccessPointPolicyWithName retrieves all s3.MultiRegionAccessPointPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3MultiRegionAccessPointPolicyWithName(name string) (*s3.MultiRegionAccessPointPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.MultiRegionAccessPointPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.MultiRegionAccessPointPolicy not found", name)
}

// GetAllS3StorageLensResources retrieves all s3.StorageLens items from an AWS CloudFormation template
func (t *Template) GetAllS3StorageLensResources() map[string]*s3.StorageLens {
	results := map[string]*s3.StorageLens{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.StorageLens:
			results[name] = resource
		}
	}
	return results
}

// GetS3StorageLensWithName retrieves all s3.StorageLens items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3StorageLensWithName(name string) (*s3.StorageLens, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.StorageLens:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.StorageLens not found", name)
}

// GetAllS3ObjectLambdaAccessPointResources retrieves all s3objectlambda.AccessPoint items from an AWS CloudFormation template
func (t *Template) GetAllS3ObjectLambdaAccessPointResources() map[string]*s3objectlambda.AccessPoint {
	results := map[string]*s3objectlambda.AccessPoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3objectlambda.AccessPoint:
			results[name] = resource
		}
	}
	return results
}

// GetS3ObjectLambdaAccessPointWithName retrieves all s3objectlambda.AccessPoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3ObjectLambdaAccessPointWithName(name string) (*s3objectlambda.AccessPoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3objectlambda.AccessPoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3objectlambda.AccessPoint not found", name)
}

// GetAllS3ObjectLambdaAccessPointPolicyResources retrieves all s3objectlambda.AccessPointPolicy items from an AWS CloudFormation template
func (t *Template) GetAllS3ObjectLambdaAccessPointPolicyResources() map[string]*s3objectlambda.AccessPointPolicy {
	results := map[string]*s3objectlambda.AccessPointPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3objectlambda.AccessPointPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetS3ObjectLambdaAccessPointPolicyWithName retrieves all s3objectlambda.AccessPointPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3ObjectLambdaAccessPointPolicyWithName(name string) (*s3objectlambda.AccessPointPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3objectlambda.AccessPointPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3objectlambda.AccessPointPolicy not found", name)
}

// GetAllS3OutpostsAccessPointResources retrieves all s3outposts.AccessPoint items from an AWS CloudFormation template
func (t *Template) GetAllS3OutpostsAccessPointResources() map[string]*s3outposts.AccessPoint {
	results := map[string]*s3outposts.AccessPoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3outposts.AccessPoint:
			results[name] = resource
		}
	}
	return results
}

// GetS3OutpostsAccessPointWithName retrieves all s3outposts.AccessPoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3OutpostsAccessPointWithName(name string) (*s3outposts.AccessPoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3outposts.AccessPoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3outposts.AccessPoint not found", name)
}

// GetAllS3OutpostsBucketResources retrieves all s3outposts.Bucket items from an AWS CloudFormation template
func (t *Template) GetAllS3OutpostsBucketResources() map[string]*s3outposts.Bucket {
	results := map[string]*s3outposts.Bucket{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3outposts.Bucket:
			results[name] = resource
		}
	}
	return results
}

// GetS3OutpostsBucketWithName retrieves all s3outposts.Bucket items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3OutpostsBucketWithName(name string) (*s3outposts.Bucket, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3outposts.Bucket:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3outposts.Bucket not found", name)
}

// GetAllS3OutpostsBucketPolicyResources retrieves all s3outposts.BucketPolicy items from an AWS CloudFormation template
func (t *Template) GetAllS3OutpostsBucketPolicyResources() map[string]*s3outposts.BucketPolicy {
	results := map[string]*s3outposts.BucketPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3outposts.BucketPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetS3OutpostsBucketPolicyWithName retrieves all s3outposts.BucketPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3OutpostsBucketPolicyWithName(name string) (*s3outposts.BucketPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3outposts.BucketPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3outposts.BucketPolicy not found", name)
}

// GetAllS3OutpostsEndpointResources retrieves all s3outposts.Endpoint items from an AWS CloudFormation template
func (t *Template) GetAllS3OutpostsEndpointResources() map[string]*s3outposts.Endpoint {
	results := map[string]*s3outposts.Endpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3outposts.Endpoint:
			results[name] = resource
		}
	}
	return results
}

// GetS3OutpostsEndpointWithName retrieves all s3outposts.Endpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3OutpostsEndpointWithName(name string) (*s3outposts.Endpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3outposts.Endpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3outposts.Endpoint not found", name)
}

// GetAllSDBDomainResources retrieves all sdb.Domain items from an AWS CloudFormation template
func (t *Template) GetAllSDBDomainResources() map[string]*sdb.Domain {
	results := map[string]*sdb.Domain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sdb.Domain:
			results[name] = resource
		}
	}
	return results
}

// GetSDBDomainWithName retrieves all sdb.Domain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSDBDomainWithName(name string) (*sdb.Domain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sdb.Domain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sdb.Domain not found", name)
}

// GetAllSESConfigurationSetResources retrieves all ses.ConfigurationSet items from an AWS CloudFormation template
func (t *Template) GetAllSESConfigurationSetResources() map[string]*ses.ConfigurationSet {
	results := map[string]*ses.ConfigurationSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ses.ConfigurationSet:
			results[name] = resource
		}
	}
	return results
}

// GetSESConfigurationSetWithName retrieves all ses.ConfigurationSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSESConfigurationSetWithName(name string) (*ses.ConfigurationSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ses.ConfigurationSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ses.ConfigurationSet not found", name)
}

// GetAllSESConfigurationSetEventDestinationResources retrieves all ses.ConfigurationSetEventDestination items from an AWS CloudFormation template
func (t *Template) GetAllSESConfigurationSetEventDestinationResources() map[string]*ses.ConfigurationSetEventDestination {
	results := map[string]*ses.ConfigurationSetEventDestination{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ses.ConfigurationSetEventDestination:
			results[name] = resource
		}
	}
	return results
}

// GetSESConfigurationSetEventDestinationWithName retrieves all ses.ConfigurationSetEventDestination items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSESConfigurationSetEventDestinationWithName(name string) (*ses.ConfigurationSetEventDestination, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ses.ConfigurationSetEventDestination:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ses.ConfigurationSetEventDestination not found", name)
}

// GetAllSESContactListResources retrieves all ses.ContactList items from an AWS CloudFormation template
func (t *Template) GetAllSESContactListResources() map[string]*ses.ContactList {
	results := map[string]*ses.ContactList{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ses.ContactList:
			results[name] = resource
		}
	}
	return results
}

// GetSESContactListWithName retrieves all ses.ContactList items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSESContactListWithName(name string) (*ses.ContactList, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ses.ContactList:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ses.ContactList not found", name)
}

// GetAllSESReceiptFilterResources retrieves all ses.ReceiptFilter items from an AWS CloudFormation template
func (t *Template) GetAllSESReceiptFilterResources() map[string]*ses.ReceiptFilter {
	results := map[string]*ses.ReceiptFilter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ses.ReceiptFilter:
			results[name] = resource
		}
	}
	return results
}

// GetSESReceiptFilterWithName retrieves all ses.ReceiptFilter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSESReceiptFilterWithName(name string) (*ses.ReceiptFilter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ses.ReceiptFilter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ses.ReceiptFilter not found", name)
}

// GetAllSESReceiptRuleResources retrieves all ses.ReceiptRule items from an AWS CloudFormation template
func (t *Template) GetAllSESReceiptRuleResources() map[string]*ses.ReceiptRule {
	results := map[string]*ses.ReceiptRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ses.ReceiptRule:
			results[name] = resource
		}
	}
	return results
}

// GetSESReceiptRuleWithName retrieves all ses.ReceiptRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSESReceiptRuleWithName(name string) (*ses.ReceiptRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ses.ReceiptRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ses.ReceiptRule not found", name)
}

// GetAllSESReceiptRuleSetResources retrieves all ses.ReceiptRuleSet items from an AWS CloudFormation template
func (t *Template) GetAllSESReceiptRuleSetResources() map[string]*ses.ReceiptRuleSet {
	results := map[string]*ses.ReceiptRuleSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ses.ReceiptRuleSet:
			results[name] = resource
		}
	}
	return results
}

// GetSESReceiptRuleSetWithName retrieves all ses.ReceiptRuleSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSESReceiptRuleSetWithName(name string) (*ses.ReceiptRuleSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ses.ReceiptRuleSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ses.ReceiptRuleSet not found", name)
}

// GetAllSESTemplateResources retrieves all ses.Template items from an AWS CloudFormation template
func (t *Template) GetAllSESTemplateResources() map[string]*ses.Template {
	results := map[string]*ses.Template{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ses.Template:
			results[name] = resource
		}
	}
	return results
}

// GetSESTemplateWithName retrieves all ses.Template items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSESTemplateWithName(name string) (*ses.Template, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ses.Template:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ses.Template not found", name)
}

// GetAllSNSSubscriptionResources retrieves all sns.Subscription items from an AWS CloudFormation template
func (t *Template) GetAllSNSSubscriptionResources() map[string]*sns.Subscription {
	results := map[string]*sns.Subscription{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sns.Subscription:
			results[name] = resource
		}
	}
	return results
}

// GetSNSSubscriptionWithName retrieves all sns.Subscription items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSNSSubscriptionWithName(name string) (*sns.Subscription, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sns.Subscription:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sns.Subscription not found", name)
}

// GetAllSNSTopicResources retrieves all sns.Topic items from an AWS CloudFormation template
func (t *Template) GetAllSNSTopicResources() map[string]*sns.Topic {
	results := map[string]*sns.Topic{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sns.Topic:
			results[name] = resource
		}
	}
	return results
}

// GetSNSTopicWithName retrieves all sns.Topic items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSNSTopicWithName(name string) (*sns.Topic, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sns.Topic:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sns.Topic not found", name)
}

// GetAllSNSTopicPolicyResources retrieves all sns.TopicPolicy items from an AWS CloudFormation template
func (t *Template) GetAllSNSTopicPolicyResources() map[string]*sns.TopicPolicy {
	results := map[string]*sns.TopicPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sns.TopicPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetSNSTopicPolicyWithName retrieves all sns.TopicPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSNSTopicPolicyWithName(name string) (*sns.TopicPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sns.TopicPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sns.TopicPolicy not found", name)
}

// GetAllSQSQueueResources retrieves all sqs.Queue items from an AWS CloudFormation template
func (t *Template) GetAllSQSQueueResources() map[string]*sqs.Queue {
	results := map[string]*sqs.Queue{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sqs.Queue:
			results[name] = resource
		}
	}
	return results
}

// GetSQSQueueWithName retrieves all sqs.Queue items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSQSQueueWithName(name string) (*sqs.Queue, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sqs.Queue:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sqs.Queue not found", name)
}

// GetAllSQSQueuePolicyResources retrieves all sqs.QueuePolicy items from an AWS CloudFormation template
func (t *Template) GetAllSQSQueuePolicyResources() map[string]*sqs.QueuePolicy {
	results := map[string]*sqs.QueuePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sqs.QueuePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetSQSQueuePolicyWithName retrieves all sqs.QueuePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSQSQueuePolicyWithName(name string) (*sqs.QueuePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sqs.QueuePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sqs.QueuePolicy not found", name)
}

// GetAllSSMAssociationResources retrieves all ssm.Association items from an AWS CloudFormation template
func (t *Template) GetAllSSMAssociationResources() map[string]*ssm.Association {
	results := map[string]*ssm.Association{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.Association:
			results[name] = resource
		}
	}
	return results
}

// GetSSMAssociationWithName retrieves all ssm.Association items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMAssociationWithName(name string) (*ssm.Association, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.Association:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.Association not found", name)
}

// GetAllSSMDocumentResources retrieves all ssm.Document items from an AWS CloudFormation template
func (t *Template) GetAllSSMDocumentResources() map[string]*ssm.Document {
	results := map[string]*ssm.Document{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.Document:
			results[name] = resource
		}
	}
	return results
}

// GetSSMDocumentWithName retrieves all ssm.Document items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMDocumentWithName(name string) (*ssm.Document, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.Document:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.Document not found", name)
}

// GetAllSSMMaintenanceWindowResources retrieves all ssm.MaintenanceWindow items from an AWS CloudFormation template
func (t *Template) GetAllSSMMaintenanceWindowResources() map[string]*ssm.MaintenanceWindow {
	results := map[string]*ssm.MaintenanceWindow{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.MaintenanceWindow:
			results[name] = resource
		}
	}
	return results
}

// GetSSMMaintenanceWindowWithName retrieves all ssm.MaintenanceWindow items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMMaintenanceWindowWithName(name string) (*ssm.MaintenanceWindow, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.MaintenanceWindow:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.MaintenanceWindow not found", name)
}

// GetAllSSMMaintenanceWindowTargetResources retrieves all ssm.MaintenanceWindowTarget items from an AWS CloudFormation template
func (t *Template) GetAllSSMMaintenanceWindowTargetResources() map[string]*ssm.MaintenanceWindowTarget {
	results := map[string]*ssm.MaintenanceWindowTarget{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.MaintenanceWindowTarget:
			results[name] = resource
		}
	}
	return results
}

// GetSSMMaintenanceWindowTargetWithName retrieves all ssm.MaintenanceWindowTarget items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMMaintenanceWindowTargetWithName(name string) (*ssm.MaintenanceWindowTarget, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.MaintenanceWindowTarget:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.MaintenanceWindowTarget not found", name)
}

// GetAllSSMMaintenanceWindowTaskResources retrieves all ssm.MaintenanceWindowTask items from an AWS CloudFormation template
func (t *Template) GetAllSSMMaintenanceWindowTaskResources() map[string]*ssm.MaintenanceWindowTask {
	results := map[string]*ssm.MaintenanceWindowTask{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.MaintenanceWindowTask:
			results[name] = resource
		}
	}
	return results
}

// GetSSMMaintenanceWindowTaskWithName retrieves all ssm.MaintenanceWindowTask items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMMaintenanceWindowTaskWithName(name string) (*ssm.MaintenanceWindowTask, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.MaintenanceWindowTask:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.MaintenanceWindowTask not found", name)
}

// GetAllSSMParameterResources retrieves all ssm.Parameter items from an AWS CloudFormation template
func (t *Template) GetAllSSMParameterResources() map[string]*ssm.Parameter {
	results := map[string]*ssm.Parameter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.Parameter:
			results[name] = resource
		}
	}
	return results
}

// GetSSMParameterWithName retrieves all ssm.Parameter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMParameterWithName(name string) (*ssm.Parameter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.Parameter:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.Parameter not found", name)
}

// GetAllSSMPatchBaselineResources retrieves all ssm.PatchBaseline items from an AWS CloudFormation template
func (t *Template) GetAllSSMPatchBaselineResources() map[string]*ssm.PatchBaseline {
	results := map[string]*ssm.PatchBaseline{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.PatchBaseline:
			results[name] = resource
		}
	}
	return results
}

// GetSSMPatchBaselineWithName retrieves all ssm.PatchBaseline items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMPatchBaselineWithName(name string) (*ssm.PatchBaseline, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.PatchBaseline:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.PatchBaseline not found", name)
}

// GetAllSSMResourceDataSyncResources retrieves all ssm.ResourceDataSync items from an AWS CloudFormation template
func (t *Template) GetAllSSMResourceDataSyncResources() map[string]*ssm.ResourceDataSync {
	results := map[string]*ssm.ResourceDataSync{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssm.ResourceDataSync:
			results[name] = resource
		}
	}
	return results
}

// GetSSMResourceDataSyncWithName retrieves all ssm.ResourceDataSync items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMResourceDataSyncWithName(name string) (*ssm.ResourceDataSync, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssm.ResourceDataSync:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssm.ResourceDataSync not found", name)
}

// GetAllSSMContactsContactResources retrieves all ssmcontacts.Contact items from an AWS CloudFormation template
func (t *Template) GetAllSSMContactsContactResources() map[string]*ssmcontacts.Contact {
	results := map[string]*ssmcontacts.Contact{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssmcontacts.Contact:
			results[name] = resource
		}
	}
	return results
}

// GetSSMContactsContactWithName retrieves all ssmcontacts.Contact items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMContactsContactWithName(name string) (*ssmcontacts.Contact, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssmcontacts.Contact:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssmcontacts.Contact not found", name)
}

// GetAllSSMContactsContactChannelResources retrieves all ssmcontacts.ContactChannel items from an AWS CloudFormation template
func (t *Template) GetAllSSMContactsContactChannelResources() map[string]*ssmcontacts.ContactChannel {
	results := map[string]*ssmcontacts.ContactChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssmcontacts.ContactChannel:
			results[name] = resource
		}
	}
	return results
}

// GetSSMContactsContactChannelWithName retrieves all ssmcontacts.ContactChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMContactsContactChannelWithName(name string) (*ssmcontacts.ContactChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssmcontacts.ContactChannel:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssmcontacts.ContactChannel not found", name)
}

// GetAllSSMIncidentsReplicationSetResources retrieves all ssmincidents.ReplicationSet items from an AWS CloudFormation template
func (t *Template) GetAllSSMIncidentsReplicationSetResources() map[string]*ssmincidents.ReplicationSet {
	results := map[string]*ssmincidents.ReplicationSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssmincidents.ReplicationSet:
			results[name] = resource
		}
	}
	return results
}

// GetSSMIncidentsReplicationSetWithName retrieves all ssmincidents.ReplicationSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMIncidentsReplicationSetWithName(name string) (*ssmincidents.ReplicationSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssmincidents.ReplicationSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssmincidents.ReplicationSet not found", name)
}

// GetAllSSMIncidentsResponsePlanResources retrieves all ssmincidents.ResponsePlan items from an AWS CloudFormation template
func (t *Template) GetAllSSMIncidentsResponsePlanResources() map[string]*ssmincidents.ResponsePlan {
	results := map[string]*ssmincidents.ResponsePlan{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ssmincidents.ResponsePlan:
			results[name] = resource
		}
	}
	return results
}

// GetSSMIncidentsResponsePlanWithName retrieves all ssmincidents.ResponsePlan items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSMIncidentsResponsePlanWithName(name string) (*ssmincidents.ResponsePlan, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ssmincidents.ResponsePlan:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ssmincidents.ResponsePlan not found", name)
}

// GetAllSSOAssignmentResources retrieves all sso.Assignment items from an AWS CloudFormation template
func (t *Template) GetAllSSOAssignmentResources() map[string]*sso.Assignment {
	results := map[string]*sso.Assignment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sso.Assignment:
			results[name] = resource
		}
	}
	return results
}

// GetSSOAssignmentWithName retrieves all sso.Assignment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSOAssignmentWithName(name string) (*sso.Assignment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sso.Assignment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sso.Assignment not found", name)
}

// GetAllSSOInstanceAccessControlAttributeConfigurationResources retrieves all sso.InstanceAccessControlAttributeConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllSSOInstanceAccessControlAttributeConfigurationResources() map[string]*sso.InstanceAccessControlAttributeConfiguration {
	results := map[string]*sso.InstanceAccessControlAttributeConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sso.InstanceAccessControlAttributeConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetSSOInstanceAccessControlAttributeConfigurationWithName retrieves all sso.InstanceAccessControlAttributeConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSOInstanceAccessControlAttributeConfigurationWithName(name string) (*sso.InstanceAccessControlAttributeConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sso.InstanceAccessControlAttributeConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sso.InstanceAccessControlAttributeConfiguration not found", name)
}

// GetAllSSOPermissionSetResources retrieves all sso.PermissionSet items from an AWS CloudFormation template
func (t *Template) GetAllSSOPermissionSetResources() map[string]*sso.PermissionSet {
	results := map[string]*sso.PermissionSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sso.PermissionSet:
			results[name] = resource
		}
	}
	return results
}

// GetSSOPermissionSetWithName retrieves all sso.PermissionSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSSOPermissionSetWithName(name string) (*sso.PermissionSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sso.PermissionSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sso.PermissionSet not found", name)
}

// GetAllSageMakerAppResources retrieves all sagemaker.App items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerAppResources() map[string]*sagemaker.App {
	results := map[string]*sagemaker.App{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.App:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerAppWithName retrieves all sagemaker.App items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerAppWithName(name string) (*sagemaker.App, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.App:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.App not found", name)
}

// GetAllSageMakerAppImageConfigResources retrieves all sagemaker.AppImageConfig items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerAppImageConfigResources() map[string]*sagemaker.AppImageConfig {
	results := map[string]*sagemaker.AppImageConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.AppImageConfig:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerAppImageConfigWithName retrieves all sagemaker.AppImageConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerAppImageConfigWithName(name string) (*sagemaker.AppImageConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.AppImageConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.AppImageConfig not found", name)
}

// GetAllSageMakerCodeRepositoryResources retrieves all sagemaker.CodeRepository items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerCodeRepositoryResources() map[string]*sagemaker.CodeRepository {
	results := map[string]*sagemaker.CodeRepository{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.CodeRepository:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerCodeRepositoryWithName retrieves all sagemaker.CodeRepository items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerCodeRepositoryWithName(name string) (*sagemaker.CodeRepository, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.CodeRepository:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.CodeRepository not found", name)
}

// GetAllSageMakerDataQualityJobDefinitionResources retrieves all sagemaker.DataQualityJobDefinition items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerDataQualityJobDefinitionResources() map[string]*sagemaker.DataQualityJobDefinition {
	results := map[string]*sagemaker.DataQualityJobDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.DataQualityJobDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerDataQualityJobDefinitionWithName retrieves all sagemaker.DataQualityJobDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerDataQualityJobDefinitionWithName(name string) (*sagemaker.DataQualityJobDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.DataQualityJobDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.DataQualityJobDefinition not found", name)
}

// GetAllSageMakerDeviceResources retrieves all sagemaker.Device items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerDeviceResources() map[string]*sagemaker.Device {
	results := map[string]*sagemaker.Device{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Device:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerDeviceWithName retrieves all sagemaker.Device items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerDeviceWithName(name string) (*sagemaker.Device, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Device:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Device not found", name)
}

// GetAllSageMakerDeviceFleetResources retrieves all sagemaker.DeviceFleet items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerDeviceFleetResources() map[string]*sagemaker.DeviceFleet {
	results := map[string]*sagemaker.DeviceFleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.DeviceFleet:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerDeviceFleetWithName retrieves all sagemaker.DeviceFleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerDeviceFleetWithName(name string) (*sagemaker.DeviceFleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.DeviceFleet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.DeviceFleet not found", name)
}

// GetAllSageMakerDomainResources retrieves all sagemaker.Domain items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerDomainResources() map[string]*sagemaker.Domain {
	results := map[string]*sagemaker.Domain{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Domain:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerDomainWithName retrieves all sagemaker.Domain items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerDomainWithName(name string) (*sagemaker.Domain, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Domain:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Domain not found", name)
}

// GetAllSageMakerEndpointResources retrieves all sagemaker.Endpoint items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerEndpointResources() map[string]*sagemaker.Endpoint {
	results := map[string]*sagemaker.Endpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Endpoint:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerEndpointWithName retrieves all sagemaker.Endpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerEndpointWithName(name string) (*sagemaker.Endpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Endpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Endpoint not found", name)
}

// GetAllSageMakerEndpointConfigResources retrieves all sagemaker.EndpointConfig items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerEndpointConfigResources() map[string]*sagemaker.EndpointConfig {
	results := map[string]*sagemaker.EndpointConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.EndpointConfig:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerEndpointConfigWithName retrieves all sagemaker.EndpointConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerEndpointConfigWithName(name string) (*sagemaker.EndpointConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.EndpointConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.EndpointConfig not found", name)
}

// GetAllSageMakerFeatureGroupResources retrieves all sagemaker.FeatureGroup items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerFeatureGroupResources() map[string]*sagemaker.FeatureGroup {
	results := map[string]*sagemaker.FeatureGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.FeatureGroup:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerFeatureGroupWithName retrieves all sagemaker.FeatureGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerFeatureGroupWithName(name string) (*sagemaker.FeatureGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.FeatureGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.FeatureGroup not found", name)
}

// GetAllSageMakerImageResources retrieves all sagemaker.Image items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerImageResources() map[string]*sagemaker.Image {
	results := map[string]*sagemaker.Image{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Image:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerImageWithName retrieves all sagemaker.Image items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerImageWithName(name string) (*sagemaker.Image, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Image:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Image not found", name)
}

// GetAllSageMakerImageVersionResources retrieves all sagemaker.ImageVersion items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerImageVersionResources() map[string]*sagemaker.ImageVersion {
	results := map[string]*sagemaker.ImageVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.ImageVersion:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerImageVersionWithName retrieves all sagemaker.ImageVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerImageVersionWithName(name string) (*sagemaker.ImageVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.ImageVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.ImageVersion not found", name)
}

// GetAllSageMakerModelResources retrieves all sagemaker.Model items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerModelResources() map[string]*sagemaker.Model {
	results := map[string]*sagemaker.Model{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Model:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerModelWithName retrieves all sagemaker.Model items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerModelWithName(name string) (*sagemaker.Model, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Model:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Model not found", name)
}

// GetAllSageMakerModelBiasJobDefinitionResources retrieves all sagemaker.ModelBiasJobDefinition items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerModelBiasJobDefinitionResources() map[string]*sagemaker.ModelBiasJobDefinition {
	results := map[string]*sagemaker.ModelBiasJobDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.ModelBiasJobDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerModelBiasJobDefinitionWithName retrieves all sagemaker.ModelBiasJobDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerModelBiasJobDefinitionWithName(name string) (*sagemaker.ModelBiasJobDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.ModelBiasJobDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.ModelBiasJobDefinition not found", name)
}

// GetAllSageMakerModelExplainabilityJobDefinitionResources retrieves all sagemaker.ModelExplainabilityJobDefinition items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerModelExplainabilityJobDefinitionResources() map[string]*sagemaker.ModelExplainabilityJobDefinition {
	results := map[string]*sagemaker.ModelExplainabilityJobDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.ModelExplainabilityJobDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerModelExplainabilityJobDefinitionWithName retrieves all sagemaker.ModelExplainabilityJobDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerModelExplainabilityJobDefinitionWithName(name string) (*sagemaker.ModelExplainabilityJobDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.ModelExplainabilityJobDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.ModelExplainabilityJobDefinition not found", name)
}

// GetAllSageMakerModelPackageGroupResources retrieves all sagemaker.ModelPackageGroup items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerModelPackageGroupResources() map[string]*sagemaker.ModelPackageGroup {
	results := map[string]*sagemaker.ModelPackageGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.ModelPackageGroup:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerModelPackageGroupWithName retrieves all sagemaker.ModelPackageGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerModelPackageGroupWithName(name string) (*sagemaker.ModelPackageGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.ModelPackageGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.ModelPackageGroup not found", name)
}

// GetAllSageMakerModelQualityJobDefinitionResources retrieves all sagemaker.ModelQualityJobDefinition items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerModelQualityJobDefinitionResources() map[string]*sagemaker.ModelQualityJobDefinition {
	results := map[string]*sagemaker.ModelQualityJobDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.ModelQualityJobDefinition:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerModelQualityJobDefinitionWithName retrieves all sagemaker.ModelQualityJobDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerModelQualityJobDefinitionWithName(name string) (*sagemaker.ModelQualityJobDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.ModelQualityJobDefinition:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.ModelQualityJobDefinition not found", name)
}

// GetAllSageMakerMonitoringScheduleResources retrieves all sagemaker.MonitoringSchedule items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerMonitoringScheduleResources() map[string]*sagemaker.MonitoringSchedule {
	results := map[string]*sagemaker.MonitoringSchedule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.MonitoringSchedule:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerMonitoringScheduleWithName retrieves all sagemaker.MonitoringSchedule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerMonitoringScheduleWithName(name string) (*sagemaker.MonitoringSchedule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.MonitoringSchedule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.MonitoringSchedule not found", name)
}

// GetAllSageMakerNotebookInstanceResources retrieves all sagemaker.NotebookInstance items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerNotebookInstanceResources() map[string]*sagemaker.NotebookInstance {
	results := map[string]*sagemaker.NotebookInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.NotebookInstance:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerNotebookInstanceWithName retrieves all sagemaker.NotebookInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerNotebookInstanceWithName(name string) (*sagemaker.NotebookInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.NotebookInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.NotebookInstance not found", name)
}

// GetAllSageMakerNotebookInstanceLifecycleConfigResources retrieves all sagemaker.NotebookInstanceLifecycleConfig items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerNotebookInstanceLifecycleConfigResources() map[string]*sagemaker.NotebookInstanceLifecycleConfig {
	results := map[string]*sagemaker.NotebookInstanceLifecycleConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.NotebookInstanceLifecycleConfig:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerNotebookInstanceLifecycleConfigWithName retrieves all sagemaker.NotebookInstanceLifecycleConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerNotebookInstanceLifecycleConfigWithName(name string) (*sagemaker.NotebookInstanceLifecycleConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.NotebookInstanceLifecycleConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.NotebookInstanceLifecycleConfig not found", name)
}

// GetAllSageMakerPipelineResources retrieves all sagemaker.Pipeline items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerPipelineResources() map[string]*sagemaker.Pipeline {
	results := map[string]*sagemaker.Pipeline{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Pipeline:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerPipelineWithName retrieves all sagemaker.Pipeline items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerPipelineWithName(name string) (*sagemaker.Pipeline, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Pipeline:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Pipeline not found", name)
}

// GetAllSageMakerProjectResources retrieves all sagemaker.Project items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerProjectResources() map[string]*sagemaker.Project {
	results := map[string]*sagemaker.Project{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Project:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerProjectWithName retrieves all sagemaker.Project items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerProjectWithName(name string) (*sagemaker.Project, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Project:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Project not found", name)
}

// GetAllSageMakerUserProfileResources retrieves all sagemaker.UserProfile items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerUserProfileResources() map[string]*sagemaker.UserProfile {
	results := map[string]*sagemaker.UserProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.UserProfile:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerUserProfileWithName retrieves all sagemaker.UserProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerUserProfileWithName(name string) (*sagemaker.UserProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.UserProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.UserProfile not found", name)
}

// GetAllSageMakerWorkteamResources retrieves all sagemaker.Workteam items from an AWS CloudFormation template
func (t *Template) GetAllSageMakerWorkteamResources() map[string]*sagemaker.Workteam {
	results := map[string]*sagemaker.Workteam{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sagemaker.Workteam:
			results[name] = resource
		}
	}
	return results
}

// GetSageMakerWorkteamWithName retrieves all sagemaker.Workteam items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSageMakerWorkteamWithName(name string) (*sagemaker.Workteam, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sagemaker.Workteam:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sagemaker.Workteam not found", name)
}

// GetAllSecretsManagerResourcePolicyResources retrieves all secretsmanager.ResourcePolicy items from an AWS CloudFormation template
func (t *Template) GetAllSecretsManagerResourcePolicyResources() map[string]*secretsmanager.ResourcePolicy {
	results := map[string]*secretsmanager.ResourcePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *secretsmanager.ResourcePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetSecretsManagerResourcePolicyWithName retrieves all secretsmanager.ResourcePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSecretsManagerResourcePolicyWithName(name string) (*secretsmanager.ResourcePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *secretsmanager.ResourcePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type secretsmanager.ResourcePolicy not found", name)
}

// GetAllSecretsManagerRotationScheduleResources retrieves all secretsmanager.RotationSchedule items from an AWS CloudFormation template
func (t *Template) GetAllSecretsManagerRotationScheduleResources() map[string]*secretsmanager.RotationSchedule {
	results := map[string]*secretsmanager.RotationSchedule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *secretsmanager.RotationSchedule:
			results[name] = resource
		}
	}
	return results
}

// GetSecretsManagerRotationScheduleWithName retrieves all secretsmanager.RotationSchedule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSecretsManagerRotationScheduleWithName(name string) (*secretsmanager.RotationSchedule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *secretsmanager.RotationSchedule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type secretsmanager.RotationSchedule not found", name)
}

// GetAllSecretsManagerSecretResources retrieves all secretsmanager.Secret items from an AWS CloudFormation template
func (t *Template) GetAllSecretsManagerSecretResources() map[string]*secretsmanager.Secret {
	results := map[string]*secretsmanager.Secret{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *secretsmanager.Secret:
			results[name] = resource
		}
	}
	return results
}

// GetSecretsManagerSecretWithName retrieves all secretsmanager.Secret items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSecretsManagerSecretWithName(name string) (*secretsmanager.Secret, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *secretsmanager.Secret:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type secretsmanager.Secret not found", name)
}

// GetAllSecretsManagerSecretTargetAttachmentResources retrieves all secretsmanager.SecretTargetAttachment items from an AWS CloudFormation template
func (t *Template) GetAllSecretsManagerSecretTargetAttachmentResources() map[string]*secretsmanager.SecretTargetAttachment {
	results := map[string]*secretsmanager.SecretTargetAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *secretsmanager.SecretTargetAttachment:
			results[name] = resource
		}
	}
	return results
}

// GetSecretsManagerSecretTargetAttachmentWithName retrieves all secretsmanager.SecretTargetAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSecretsManagerSecretTargetAttachmentWithName(name string) (*secretsmanager.SecretTargetAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *secretsmanager.SecretTargetAttachment:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type secretsmanager.SecretTargetAttachment not found", name)
}

// GetAllSecurityHubHubResources retrieves all securityhub.Hub items from an AWS CloudFormation template
func (t *Template) GetAllSecurityHubHubResources() map[string]*securityhub.Hub {
	results := map[string]*securityhub.Hub{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *securityhub.Hub:
			results[name] = resource
		}
	}
	return results
}

// GetSecurityHubHubWithName retrieves all securityhub.Hub items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSecurityHubHubWithName(name string) (*securityhub.Hub, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *securityhub.Hub:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type securityhub.Hub not found", name)
}

// GetAllServerlessApiResources retrieves all serverless.Api items from an AWS CloudFormation template
func (t *Template) GetAllServerlessApiResources() map[string]*serverless.Api {
	results := map[string]*serverless.Api{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *serverless.Api:
			results[name] = resource
		}
	}
	return results
}

// GetServerlessApiWithName retrieves all serverless.Api items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServerlessApiWithName(name string) (*serverless.Api, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *serverless.Api:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type serverless.Api not found", name)
}

// GetAllServerlessApplicationResources retrieves all serverless.Application items from an AWS CloudFormation template
func (t *Template) GetAllServerlessApplicationResources() map[string]*serverless.Application {
	results := map[string]*serverless.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *serverless.Application:
			results[name] = resource
		}
	}
	return results
}

// GetServerlessApplicationWithName retrieves all serverless.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServerlessApplicationWithName(name string) (*serverless.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *serverless.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type serverless.Application not found", name)
}

// GetAllServerlessFunctionResources retrieves all serverless.Function items from an AWS CloudFormation template
func (t *Template) GetAllServerlessFunctionResources() map[string]*serverless.Function {
	results := map[string]*serverless.Function{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *serverless.Function:
			results[name] = resource
		}
	}
	return results
}

// GetServerlessFunctionWithName retrieves all serverless.Function items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServerlessFunctionWithName(name string) (*serverless.Function, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *serverless.Function:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type serverless.Function not found", name)
}

// GetAllServerlessLayerVersionResources retrieves all serverless.LayerVersion items from an AWS CloudFormation template
func (t *Template) GetAllServerlessLayerVersionResources() map[string]*serverless.LayerVersion {
	results := map[string]*serverless.LayerVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *serverless.LayerVersion:
			results[name] = resource
		}
	}
	return results
}

// GetServerlessLayerVersionWithName retrieves all serverless.LayerVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServerlessLayerVersionWithName(name string) (*serverless.LayerVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *serverless.LayerVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type serverless.LayerVersion not found", name)
}

// GetAllServerlessSimpleTableResources retrieves all serverless.SimpleTable items from an AWS CloudFormation template
func (t *Template) GetAllServerlessSimpleTableResources() map[string]*serverless.SimpleTable {
	results := map[string]*serverless.SimpleTable{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *serverless.SimpleTable:
			results[name] = resource
		}
	}
	return results
}

// GetServerlessSimpleTableWithName retrieves all serverless.SimpleTable items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServerlessSimpleTableWithName(name string) (*serverless.SimpleTable, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *serverless.SimpleTable:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type serverless.SimpleTable not found", name)
}

// GetAllServerlessStateMachineResources retrieves all serverless.StateMachine items from an AWS CloudFormation template
func (t *Template) GetAllServerlessStateMachineResources() map[string]*serverless.StateMachine {
	results := map[string]*serverless.StateMachine{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *serverless.StateMachine:
			results[name] = resource
		}
	}
	return results
}

// GetServerlessStateMachineWithName retrieves all serverless.StateMachine items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServerlessStateMachineWithName(name string) (*serverless.StateMachine, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *serverless.StateMachine:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type serverless.StateMachine not found", name)
}

// GetAllServiceCatalogAcceptedPortfolioShareResources retrieves all servicecatalog.AcceptedPortfolioShare items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogAcceptedPortfolioShareResources() map[string]*servicecatalog.AcceptedPortfolioShare {
	results := map[string]*servicecatalog.AcceptedPortfolioShare{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.AcceptedPortfolioShare:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogAcceptedPortfolioShareWithName retrieves all servicecatalog.AcceptedPortfolioShare items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogAcceptedPortfolioShareWithName(name string) (*servicecatalog.AcceptedPortfolioShare, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.AcceptedPortfolioShare:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.AcceptedPortfolioShare not found", name)
}

// GetAllServiceCatalogCloudFormationProductResources retrieves all servicecatalog.CloudFormationProduct items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogCloudFormationProductResources() map[string]*servicecatalog.CloudFormationProduct {
	results := map[string]*servicecatalog.CloudFormationProduct{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.CloudFormationProduct:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogCloudFormationProductWithName retrieves all servicecatalog.CloudFormationProduct items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogCloudFormationProductWithName(name string) (*servicecatalog.CloudFormationProduct, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.CloudFormationProduct:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.CloudFormationProduct not found", name)
}

// GetAllServiceCatalogCloudFormationProvisionedProductResources retrieves all servicecatalog.CloudFormationProvisionedProduct items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogCloudFormationProvisionedProductResources() map[string]*servicecatalog.CloudFormationProvisionedProduct {
	results := map[string]*servicecatalog.CloudFormationProvisionedProduct{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.CloudFormationProvisionedProduct:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogCloudFormationProvisionedProductWithName retrieves all servicecatalog.CloudFormationProvisionedProduct items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogCloudFormationProvisionedProductWithName(name string) (*servicecatalog.CloudFormationProvisionedProduct, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.CloudFormationProvisionedProduct:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.CloudFormationProvisionedProduct not found", name)
}

// GetAllServiceCatalogLaunchNotificationConstraintResources retrieves all servicecatalog.LaunchNotificationConstraint items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogLaunchNotificationConstraintResources() map[string]*servicecatalog.LaunchNotificationConstraint {
	results := map[string]*servicecatalog.LaunchNotificationConstraint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.LaunchNotificationConstraint:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogLaunchNotificationConstraintWithName retrieves all servicecatalog.LaunchNotificationConstraint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogLaunchNotificationConstraintWithName(name string) (*servicecatalog.LaunchNotificationConstraint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.LaunchNotificationConstraint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.LaunchNotificationConstraint not found", name)
}

// GetAllServiceCatalogLaunchRoleConstraintResources retrieves all servicecatalog.LaunchRoleConstraint items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogLaunchRoleConstraintResources() map[string]*servicecatalog.LaunchRoleConstraint {
	results := map[string]*servicecatalog.LaunchRoleConstraint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.LaunchRoleConstraint:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogLaunchRoleConstraintWithName retrieves all servicecatalog.LaunchRoleConstraint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogLaunchRoleConstraintWithName(name string) (*servicecatalog.LaunchRoleConstraint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.LaunchRoleConstraint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.LaunchRoleConstraint not found", name)
}

// GetAllServiceCatalogLaunchTemplateConstraintResources retrieves all servicecatalog.LaunchTemplateConstraint items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogLaunchTemplateConstraintResources() map[string]*servicecatalog.LaunchTemplateConstraint {
	results := map[string]*servicecatalog.LaunchTemplateConstraint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.LaunchTemplateConstraint:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogLaunchTemplateConstraintWithName retrieves all servicecatalog.LaunchTemplateConstraint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogLaunchTemplateConstraintWithName(name string) (*servicecatalog.LaunchTemplateConstraint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.LaunchTemplateConstraint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.LaunchTemplateConstraint not found", name)
}

// GetAllServiceCatalogPortfolioResources retrieves all servicecatalog.Portfolio items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogPortfolioResources() map[string]*servicecatalog.Portfolio {
	results := map[string]*servicecatalog.Portfolio{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.Portfolio:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogPortfolioWithName retrieves all servicecatalog.Portfolio items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogPortfolioWithName(name string) (*servicecatalog.Portfolio, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.Portfolio:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.Portfolio not found", name)
}

// GetAllServiceCatalogPortfolioPrincipalAssociationResources retrieves all servicecatalog.PortfolioPrincipalAssociation items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogPortfolioPrincipalAssociationResources() map[string]*servicecatalog.PortfolioPrincipalAssociation {
	results := map[string]*servicecatalog.PortfolioPrincipalAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.PortfolioPrincipalAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogPortfolioPrincipalAssociationWithName retrieves all servicecatalog.PortfolioPrincipalAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogPortfolioPrincipalAssociationWithName(name string) (*servicecatalog.PortfolioPrincipalAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.PortfolioPrincipalAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.PortfolioPrincipalAssociation not found", name)
}

// GetAllServiceCatalogPortfolioProductAssociationResources retrieves all servicecatalog.PortfolioProductAssociation items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogPortfolioProductAssociationResources() map[string]*servicecatalog.PortfolioProductAssociation {
	results := map[string]*servicecatalog.PortfolioProductAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.PortfolioProductAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogPortfolioProductAssociationWithName retrieves all servicecatalog.PortfolioProductAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogPortfolioProductAssociationWithName(name string) (*servicecatalog.PortfolioProductAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.PortfolioProductAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.PortfolioProductAssociation not found", name)
}

// GetAllServiceCatalogPortfolioShareResources retrieves all servicecatalog.PortfolioShare items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogPortfolioShareResources() map[string]*servicecatalog.PortfolioShare {
	results := map[string]*servicecatalog.PortfolioShare{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.PortfolioShare:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogPortfolioShareWithName retrieves all servicecatalog.PortfolioShare items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogPortfolioShareWithName(name string) (*servicecatalog.PortfolioShare, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.PortfolioShare:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.PortfolioShare not found", name)
}

// GetAllServiceCatalogResourceUpdateConstraintResources retrieves all servicecatalog.ResourceUpdateConstraint items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogResourceUpdateConstraintResources() map[string]*servicecatalog.ResourceUpdateConstraint {
	results := map[string]*servicecatalog.ResourceUpdateConstraint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.ResourceUpdateConstraint:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogResourceUpdateConstraintWithName retrieves all servicecatalog.ResourceUpdateConstraint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogResourceUpdateConstraintWithName(name string) (*servicecatalog.ResourceUpdateConstraint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.ResourceUpdateConstraint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.ResourceUpdateConstraint not found", name)
}

// GetAllServiceCatalogServiceActionResources retrieves all servicecatalog.ServiceAction items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogServiceActionResources() map[string]*servicecatalog.ServiceAction {
	results := map[string]*servicecatalog.ServiceAction{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.ServiceAction:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogServiceActionWithName retrieves all servicecatalog.ServiceAction items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogServiceActionWithName(name string) (*servicecatalog.ServiceAction, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.ServiceAction:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.ServiceAction not found", name)
}

// GetAllServiceCatalogServiceActionAssociationResources retrieves all servicecatalog.ServiceActionAssociation items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogServiceActionAssociationResources() map[string]*servicecatalog.ServiceActionAssociation {
	results := map[string]*servicecatalog.ServiceActionAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.ServiceActionAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogServiceActionAssociationWithName retrieves all servicecatalog.ServiceActionAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogServiceActionAssociationWithName(name string) (*servicecatalog.ServiceActionAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.ServiceActionAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.ServiceActionAssociation not found", name)
}

// GetAllServiceCatalogStackSetConstraintResources retrieves all servicecatalog.StackSetConstraint items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogStackSetConstraintResources() map[string]*servicecatalog.StackSetConstraint {
	results := map[string]*servicecatalog.StackSetConstraint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.StackSetConstraint:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogStackSetConstraintWithName retrieves all servicecatalog.StackSetConstraint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogStackSetConstraintWithName(name string) (*servicecatalog.StackSetConstraint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.StackSetConstraint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.StackSetConstraint not found", name)
}

// GetAllServiceCatalogTagOptionResources retrieves all servicecatalog.TagOption items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogTagOptionResources() map[string]*servicecatalog.TagOption {
	results := map[string]*servicecatalog.TagOption{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.TagOption:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogTagOptionWithName retrieves all servicecatalog.TagOption items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogTagOptionWithName(name string) (*servicecatalog.TagOption, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.TagOption:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.TagOption not found", name)
}

// GetAllServiceCatalogTagOptionAssociationResources retrieves all servicecatalog.TagOptionAssociation items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogTagOptionAssociationResources() map[string]*servicecatalog.TagOptionAssociation {
	results := map[string]*servicecatalog.TagOptionAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalog.TagOptionAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogTagOptionAssociationWithName retrieves all servicecatalog.TagOptionAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogTagOptionAssociationWithName(name string) (*servicecatalog.TagOptionAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalog.TagOptionAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalog.TagOptionAssociation not found", name)
}

// GetAllServiceCatalogAppRegistryApplicationResources retrieves all servicecatalogappregistry.Application items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogAppRegistryApplicationResources() map[string]*servicecatalogappregistry.Application {
	results := map[string]*servicecatalogappregistry.Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.Application:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogAppRegistryApplicationWithName retrieves all servicecatalogappregistry.Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogAppRegistryApplicationWithName(name string) (*servicecatalogappregistry.Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.Application:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalogappregistry.Application not found", name)
}

// GetAllServiceCatalogAppRegistryAttributeGroupResources retrieves all servicecatalogappregistry.AttributeGroup items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogAppRegistryAttributeGroupResources() map[string]*servicecatalogappregistry.AttributeGroup {
	results := map[string]*servicecatalogappregistry.AttributeGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.AttributeGroup:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogAppRegistryAttributeGroupWithName retrieves all servicecatalogappregistry.AttributeGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogAppRegistryAttributeGroupWithName(name string) (*servicecatalogappregistry.AttributeGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.AttributeGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalogappregistry.AttributeGroup not found", name)
}

// GetAllServiceCatalogAppRegistryAttributeGroupAssociationResources retrieves all servicecatalogappregistry.AttributeGroupAssociation items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogAppRegistryAttributeGroupAssociationResources() map[string]*servicecatalogappregistry.AttributeGroupAssociation {
	results := map[string]*servicecatalogappregistry.AttributeGroupAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.AttributeGroupAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogAppRegistryAttributeGroupAssociationWithName retrieves all servicecatalogappregistry.AttributeGroupAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogAppRegistryAttributeGroupAssociationWithName(name string) (*servicecatalogappregistry.AttributeGroupAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.AttributeGroupAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalogappregistry.AttributeGroupAssociation not found", name)
}

// GetAllServiceCatalogAppRegistryResourceAssociationResources retrieves all servicecatalogappregistry.ResourceAssociation items from an AWS CloudFormation template
func (t *Template) GetAllServiceCatalogAppRegistryResourceAssociationResources() map[string]*servicecatalogappregistry.ResourceAssociation {
	results := map[string]*servicecatalogappregistry.ResourceAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.ResourceAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetServiceCatalogAppRegistryResourceAssociationWithName retrieves all servicecatalogappregistry.ResourceAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceCatalogAppRegistryResourceAssociationWithName(name string) (*servicecatalogappregistry.ResourceAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicecatalogappregistry.ResourceAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicecatalogappregistry.ResourceAssociation not found", name)
}

// GetAllServiceDiscoveryHttpNamespaceResources retrieves all servicediscovery.HttpNamespace items from an AWS CloudFormation template
func (t *Template) GetAllServiceDiscoveryHttpNamespaceResources() map[string]*servicediscovery.HttpNamespace {
	results := map[string]*servicediscovery.HttpNamespace{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicediscovery.HttpNamespace:
			results[name] = resource
		}
	}
	return results
}

// GetServiceDiscoveryHttpNamespaceWithName retrieves all servicediscovery.HttpNamespace items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceDiscoveryHttpNamespaceWithName(name string) (*servicediscovery.HttpNamespace, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicediscovery.HttpNamespace:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicediscovery.HttpNamespace not found", name)
}

// GetAllServiceDiscoveryInstanceResources retrieves all servicediscovery.Instance items from an AWS CloudFormation template
func (t *Template) GetAllServiceDiscoveryInstanceResources() map[string]*servicediscovery.Instance {
	results := map[string]*servicediscovery.Instance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicediscovery.Instance:
			results[name] = resource
		}
	}
	return results
}

// GetServiceDiscoveryInstanceWithName retrieves all servicediscovery.Instance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceDiscoveryInstanceWithName(name string) (*servicediscovery.Instance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicediscovery.Instance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicediscovery.Instance not found", name)
}

// GetAllServiceDiscoveryPrivateDnsNamespaceResources retrieves all servicediscovery.PrivateDnsNamespace items from an AWS CloudFormation template
func (t *Template) GetAllServiceDiscoveryPrivateDnsNamespaceResources() map[string]*servicediscovery.PrivateDnsNamespace {
	results := map[string]*servicediscovery.PrivateDnsNamespace{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicediscovery.PrivateDnsNamespace:
			results[name] = resource
		}
	}
	return results
}

// GetServiceDiscoveryPrivateDnsNamespaceWithName retrieves all servicediscovery.PrivateDnsNamespace items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceDiscoveryPrivateDnsNamespaceWithName(name string) (*servicediscovery.PrivateDnsNamespace, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicediscovery.PrivateDnsNamespace:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicediscovery.PrivateDnsNamespace not found", name)
}

// GetAllServiceDiscoveryPublicDnsNamespaceResources retrieves all servicediscovery.PublicDnsNamespace items from an AWS CloudFormation template
func (t *Template) GetAllServiceDiscoveryPublicDnsNamespaceResources() map[string]*servicediscovery.PublicDnsNamespace {
	results := map[string]*servicediscovery.PublicDnsNamespace{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicediscovery.PublicDnsNamespace:
			results[name] = resource
		}
	}
	return results
}

// GetServiceDiscoveryPublicDnsNamespaceWithName retrieves all servicediscovery.PublicDnsNamespace items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceDiscoveryPublicDnsNamespaceWithName(name string) (*servicediscovery.PublicDnsNamespace, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicediscovery.PublicDnsNamespace:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicediscovery.PublicDnsNamespace not found", name)
}

// GetAllServiceDiscoveryServiceResources retrieves all servicediscovery.Service items from an AWS CloudFormation template
func (t *Template) GetAllServiceDiscoveryServiceResources() map[string]*servicediscovery.Service {
	results := map[string]*servicediscovery.Service{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *servicediscovery.Service:
			results[name] = resource
		}
	}
	return results
}

// GetServiceDiscoveryServiceWithName retrieves all servicediscovery.Service items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetServiceDiscoveryServiceWithName(name string) (*servicediscovery.Service, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *servicediscovery.Service:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type servicediscovery.Service not found", name)
}

// GetAllSignerProfilePermissionResources retrieves all signer.ProfilePermission items from an AWS CloudFormation template
func (t *Template) GetAllSignerProfilePermissionResources() map[string]*signer.ProfilePermission {
	results := map[string]*signer.ProfilePermission{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *signer.ProfilePermission:
			results[name] = resource
		}
	}
	return results
}

// GetSignerProfilePermissionWithName retrieves all signer.ProfilePermission items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSignerProfilePermissionWithName(name string) (*signer.ProfilePermission, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *signer.ProfilePermission:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type signer.ProfilePermission not found", name)
}

// GetAllSignerSigningProfileResources retrieves all signer.SigningProfile items from an AWS CloudFormation template
func (t *Template) GetAllSignerSigningProfileResources() map[string]*signer.SigningProfile {
	results := map[string]*signer.SigningProfile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *signer.SigningProfile:
			results[name] = resource
		}
	}
	return results
}

// GetSignerSigningProfileWithName retrieves all signer.SigningProfile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSignerSigningProfileWithName(name string) (*signer.SigningProfile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *signer.SigningProfile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type signer.SigningProfile not found", name)
}

// GetAllStepFunctionsActivityResources retrieves all stepfunctions.Activity items from an AWS CloudFormation template
func (t *Template) GetAllStepFunctionsActivityResources() map[string]*stepfunctions.Activity {
	results := map[string]*stepfunctions.Activity{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *stepfunctions.Activity:
			results[name] = resource
		}
	}
	return results
}

// GetStepFunctionsActivityWithName retrieves all stepfunctions.Activity items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetStepFunctionsActivityWithName(name string) (*stepfunctions.Activity, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *stepfunctions.Activity:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type stepfunctions.Activity not found", name)
}

// GetAllStepFunctionsStateMachineResources retrieves all stepfunctions.StateMachine items from an AWS CloudFormation template
func (t *Template) GetAllStepFunctionsStateMachineResources() map[string]*stepfunctions.StateMachine {
	results := map[string]*stepfunctions.StateMachine{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *stepfunctions.StateMachine:
			results[name] = resource
		}
	}
	return results
}

// GetStepFunctionsStateMachineWithName retrieves all stepfunctions.StateMachine items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetStepFunctionsStateMachineWithName(name string) (*stepfunctions.StateMachine, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *stepfunctions.StateMachine:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type stepfunctions.StateMachine not found", name)
}

// GetAllSyntheticsCanaryResources retrieves all synthetics.Canary items from an AWS CloudFormation template
func (t *Template) GetAllSyntheticsCanaryResources() map[string]*synthetics.Canary {
	results := map[string]*synthetics.Canary{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *synthetics.Canary:
			results[name] = resource
		}
	}
	return results
}

// GetSyntheticsCanaryWithName retrieves all synthetics.Canary items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSyntheticsCanaryWithName(name string) (*synthetics.Canary, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *synthetics.Canary:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type synthetics.Canary not found", name)
}

// GetAllTimestreamDatabaseResources retrieves all timestream.Database items from an AWS CloudFormation template
func (t *Template) GetAllTimestreamDatabaseResources() map[string]*timestream.Database {
	results := map[string]*timestream.Database{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *timestream.Database:
			results[name] = resource
		}
	}
	return results
}

// GetTimestreamDatabaseWithName retrieves all timestream.Database items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetTimestreamDatabaseWithName(name string) (*timestream.Database, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *timestream.Database:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type timestream.Database not found", name)
}

// GetAllTimestreamScheduledQueryResources retrieves all timestream.ScheduledQuery items from an AWS CloudFormation template
func (t *Template) GetAllTimestreamScheduledQueryResources() map[string]*timestream.ScheduledQuery {
	results := map[string]*timestream.ScheduledQuery{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *timestream.ScheduledQuery:
			results[name] = resource
		}
	}
	return results
}

// GetTimestreamScheduledQueryWithName retrieves all timestream.ScheduledQuery items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetTimestreamScheduledQueryWithName(name string) (*timestream.ScheduledQuery, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *timestream.ScheduledQuery:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type timestream.ScheduledQuery not found", name)
}

// GetAllTimestreamTableResources retrieves all timestream.Table items from an AWS CloudFormation template
func (t *Template) GetAllTimestreamTableResources() map[string]*timestream.Table {
	results := map[string]*timestream.Table{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *timestream.Table:
			results[name] = resource
		}
	}
	return results
}

// GetTimestreamTableWithName retrieves all timestream.Table items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetTimestreamTableWithName(name string) (*timestream.Table, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *timestream.Table:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type timestream.Table not found", name)
}

// GetAllTransferServerResources retrieves all transfer.Server items from an AWS CloudFormation template
func (t *Template) GetAllTransferServerResources() map[string]*transfer.Server {
	results := map[string]*transfer.Server{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *transfer.Server:
			results[name] = resource
		}
	}
	return results
}

// GetTransferServerWithName retrieves all transfer.Server items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetTransferServerWithName(name string) (*transfer.Server, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *transfer.Server:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type transfer.Server not found", name)
}

// GetAllTransferUserResources retrieves all transfer.User items from an AWS CloudFormation template
func (t *Template) GetAllTransferUserResources() map[string]*transfer.User {
	results := map[string]*transfer.User{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *transfer.User:
			results[name] = resource
		}
	}
	return results
}

// GetTransferUserWithName retrieves all transfer.User items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetTransferUserWithName(name string) (*transfer.User, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *transfer.User:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type transfer.User not found", name)
}

// GetAllTransferWorkflowResources retrieves all transfer.Workflow items from an AWS CloudFormation template
func (t *Template) GetAllTransferWorkflowResources() map[string]*transfer.Workflow {
	results := map[string]*transfer.Workflow{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *transfer.Workflow:
			results[name] = resource
		}
	}
	return results
}

// GetTransferWorkflowWithName retrieves all transfer.Workflow items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetTransferWorkflowWithName(name string) (*transfer.Workflow, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *transfer.Workflow:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type transfer.Workflow not found", name)
}

// GetAllWAFByteMatchSetResources retrieves all waf.ByteMatchSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFByteMatchSetResources() map[string]*waf.ByteMatchSet {
	results := map[string]*waf.ByteMatchSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *waf.ByteMatchSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFByteMatchSetWithName retrieves all waf.ByteMatchSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFByteMatchSetWithName(name string) (*waf.ByteMatchSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *waf.ByteMatchSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type waf.ByteMatchSet not found", name)
}

// GetAllWAFIPSetResources retrieves all waf.IPSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFIPSetResources() map[string]*waf.IPSet {
	results := map[string]*waf.IPSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *waf.IPSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFIPSetWithName retrieves all waf.IPSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFIPSetWithName(name string) (*waf.IPSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *waf.IPSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type waf.IPSet not found", name)
}

// GetAllWAFRuleResources retrieves all waf.Rule items from an AWS CloudFormation template
func (t *Template) GetAllWAFRuleResources() map[string]*waf.Rule {
	results := map[string]*waf.Rule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *waf.Rule:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRuleWithName retrieves all waf.Rule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRuleWithName(name string) (*waf.Rule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *waf.Rule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type waf.Rule not found", name)
}

// GetAllWAFSizeConstraintSetResources retrieves all waf.SizeConstraintSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFSizeConstraintSetResources() map[string]*waf.SizeConstraintSet {
	results := map[string]*waf.SizeConstraintSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *waf.SizeConstraintSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFSizeConstraintSetWithName retrieves all waf.SizeConstraintSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFSizeConstraintSetWithName(name string) (*waf.SizeConstraintSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *waf.SizeConstraintSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type waf.SizeConstraintSet not found", name)
}

// GetAllWAFSqlInjectionMatchSetResources retrieves all waf.SqlInjectionMatchSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFSqlInjectionMatchSetResources() map[string]*waf.SqlInjectionMatchSet {
	results := map[string]*waf.SqlInjectionMatchSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *waf.SqlInjectionMatchSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFSqlInjectionMatchSetWithName retrieves all waf.SqlInjectionMatchSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFSqlInjectionMatchSetWithName(name string) (*waf.SqlInjectionMatchSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *waf.SqlInjectionMatchSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type waf.SqlInjectionMatchSet not found", name)
}

// GetAllWAFWebACLResources retrieves all waf.WebACL items from an AWS CloudFormation template
func (t *Template) GetAllWAFWebACLResources() map[string]*waf.WebACL {
	results := map[string]*waf.WebACL{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *waf.WebACL:
			results[name] = resource
		}
	}
	return results
}

// GetWAFWebACLWithName retrieves all waf.WebACL items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFWebACLWithName(name string) (*waf.WebACL, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *waf.WebACL:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type waf.WebACL not found", name)
}

// GetAllWAFXssMatchSetResources retrieves all waf.XssMatchSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFXssMatchSetResources() map[string]*waf.XssMatchSet {
	results := map[string]*waf.XssMatchSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *waf.XssMatchSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFXssMatchSetWithName retrieves all waf.XssMatchSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFXssMatchSetWithName(name string) (*waf.XssMatchSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *waf.XssMatchSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type waf.XssMatchSet not found", name)
}

// GetAllWAFRegionalByteMatchSetResources retrieves all wafregional.ByteMatchSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalByteMatchSetResources() map[string]*wafregional.ByteMatchSet {
	results := map[string]*wafregional.ByteMatchSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.ByteMatchSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalByteMatchSetWithName retrieves all wafregional.ByteMatchSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalByteMatchSetWithName(name string) (*wafregional.ByteMatchSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.ByteMatchSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.ByteMatchSet not found", name)
}

// GetAllWAFRegionalGeoMatchSetResources retrieves all wafregional.GeoMatchSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalGeoMatchSetResources() map[string]*wafregional.GeoMatchSet {
	results := map[string]*wafregional.GeoMatchSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.GeoMatchSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalGeoMatchSetWithName retrieves all wafregional.GeoMatchSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalGeoMatchSetWithName(name string) (*wafregional.GeoMatchSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.GeoMatchSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.GeoMatchSet not found", name)
}

// GetAllWAFRegionalIPSetResources retrieves all wafregional.IPSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalIPSetResources() map[string]*wafregional.IPSet {
	results := map[string]*wafregional.IPSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.IPSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalIPSetWithName retrieves all wafregional.IPSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalIPSetWithName(name string) (*wafregional.IPSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.IPSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.IPSet not found", name)
}

// GetAllWAFRegionalRateBasedRuleResources retrieves all wafregional.RateBasedRule items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalRateBasedRuleResources() map[string]*wafregional.RateBasedRule {
	results := map[string]*wafregional.RateBasedRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.RateBasedRule:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalRateBasedRuleWithName retrieves all wafregional.RateBasedRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalRateBasedRuleWithName(name string) (*wafregional.RateBasedRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.RateBasedRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.RateBasedRule not found", name)
}

// GetAllWAFRegionalRegexPatternSetResources retrieves all wafregional.RegexPatternSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalRegexPatternSetResources() map[string]*wafregional.RegexPatternSet {
	results := map[string]*wafregional.RegexPatternSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.RegexPatternSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalRegexPatternSetWithName retrieves all wafregional.RegexPatternSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalRegexPatternSetWithName(name string) (*wafregional.RegexPatternSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.RegexPatternSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.RegexPatternSet not found", name)
}

// GetAllWAFRegionalRuleResources retrieves all wafregional.Rule items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalRuleResources() map[string]*wafregional.Rule {
	results := map[string]*wafregional.Rule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.Rule:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalRuleWithName retrieves all wafregional.Rule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalRuleWithName(name string) (*wafregional.Rule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.Rule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.Rule not found", name)
}

// GetAllWAFRegionalSizeConstraintSetResources retrieves all wafregional.SizeConstraintSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalSizeConstraintSetResources() map[string]*wafregional.SizeConstraintSet {
	results := map[string]*wafregional.SizeConstraintSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.SizeConstraintSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalSizeConstraintSetWithName retrieves all wafregional.SizeConstraintSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalSizeConstraintSetWithName(name string) (*wafregional.SizeConstraintSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.SizeConstraintSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.SizeConstraintSet not found", name)
}

// GetAllWAFRegionalSqlInjectionMatchSetResources retrieves all wafregional.SqlInjectionMatchSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalSqlInjectionMatchSetResources() map[string]*wafregional.SqlInjectionMatchSet {
	results := map[string]*wafregional.SqlInjectionMatchSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.SqlInjectionMatchSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalSqlInjectionMatchSetWithName retrieves all wafregional.SqlInjectionMatchSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalSqlInjectionMatchSetWithName(name string) (*wafregional.SqlInjectionMatchSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.SqlInjectionMatchSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.SqlInjectionMatchSet not found", name)
}

// GetAllWAFRegionalWebACLResources retrieves all wafregional.WebACL items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalWebACLResources() map[string]*wafregional.WebACL {
	results := map[string]*wafregional.WebACL{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.WebACL:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalWebACLWithName retrieves all wafregional.WebACL items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalWebACLWithName(name string) (*wafregional.WebACL, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.WebACL:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.WebACL not found", name)
}

// GetAllWAFRegionalWebACLAssociationResources retrieves all wafregional.WebACLAssociation items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalWebACLAssociationResources() map[string]*wafregional.WebACLAssociation {
	results := map[string]*wafregional.WebACLAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.WebACLAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalWebACLAssociationWithName retrieves all wafregional.WebACLAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalWebACLAssociationWithName(name string) (*wafregional.WebACLAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.WebACLAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.WebACLAssociation not found", name)
}

// GetAllWAFRegionalXssMatchSetResources retrieves all wafregional.XssMatchSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFRegionalXssMatchSetResources() map[string]*wafregional.XssMatchSet {
	results := map[string]*wafregional.XssMatchSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafregional.XssMatchSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFRegionalXssMatchSetWithName retrieves all wafregional.XssMatchSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFRegionalXssMatchSetWithName(name string) (*wafregional.XssMatchSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafregional.XssMatchSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafregional.XssMatchSet not found", name)
}

// GetAllWAFv2IPSetResources retrieves all wafv2.IPSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFv2IPSetResources() map[string]*wafv2.IPSet {
	results := map[string]*wafv2.IPSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafv2.IPSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFv2IPSetWithName retrieves all wafv2.IPSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFv2IPSetWithName(name string) (*wafv2.IPSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafv2.IPSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafv2.IPSet not found", name)
}

// GetAllWAFv2LoggingConfigurationResources retrieves all wafv2.LoggingConfiguration items from an AWS CloudFormation template
func (t *Template) GetAllWAFv2LoggingConfigurationResources() map[string]*wafv2.LoggingConfiguration {
	results := map[string]*wafv2.LoggingConfiguration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafv2.LoggingConfiguration:
			results[name] = resource
		}
	}
	return results
}

// GetWAFv2LoggingConfigurationWithName retrieves all wafv2.LoggingConfiguration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFv2LoggingConfigurationWithName(name string) (*wafv2.LoggingConfiguration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafv2.LoggingConfiguration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafv2.LoggingConfiguration not found", name)
}

// GetAllWAFv2RegexPatternSetResources retrieves all wafv2.RegexPatternSet items from an AWS CloudFormation template
func (t *Template) GetAllWAFv2RegexPatternSetResources() map[string]*wafv2.RegexPatternSet {
	results := map[string]*wafv2.RegexPatternSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafv2.RegexPatternSet:
			results[name] = resource
		}
	}
	return results
}

// GetWAFv2RegexPatternSetWithName retrieves all wafv2.RegexPatternSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFv2RegexPatternSetWithName(name string) (*wafv2.RegexPatternSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafv2.RegexPatternSet:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafv2.RegexPatternSet not found", name)
}

// GetAllWAFv2RuleGroupResources retrieves all wafv2.RuleGroup items from an AWS CloudFormation template
func (t *Template) GetAllWAFv2RuleGroupResources() map[string]*wafv2.RuleGroup {
	results := map[string]*wafv2.RuleGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafv2.RuleGroup:
			results[name] = resource
		}
	}
	return results
}

// GetWAFv2RuleGroupWithName retrieves all wafv2.RuleGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFv2RuleGroupWithName(name string) (*wafv2.RuleGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafv2.RuleGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafv2.RuleGroup not found", name)
}

// GetAllWAFv2WebACLResources retrieves all wafv2.WebACL items from an AWS CloudFormation template
func (t *Template) GetAllWAFv2WebACLResources() map[string]*wafv2.WebACL {
	results := map[string]*wafv2.WebACL{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafv2.WebACL:
			results[name] = resource
		}
	}
	return results
}

// GetWAFv2WebACLWithName retrieves all wafv2.WebACL items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFv2WebACLWithName(name string) (*wafv2.WebACL, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafv2.WebACL:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafv2.WebACL not found", name)
}

// GetAllWAFv2WebACLAssociationResources retrieves all wafv2.WebACLAssociation items from an AWS CloudFormation template
func (t *Template) GetAllWAFv2WebACLAssociationResources() map[string]*wafv2.WebACLAssociation {
	results := map[string]*wafv2.WebACLAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wafv2.WebACLAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetWAFv2WebACLAssociationWithName retrieves all wafv2.WebACLAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWAFv2WebACLAssociationWithName(name string) (*wafv2.WebACLAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wafv2.WebACLAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wafv2.WebACLAssociation not found", name)
}

// GetAllWisdomAssistantResources retrieves all wisdom.Assistant items from an AWS CloudFormation template
func (t *Template) GetAllWisdomAssistantResources() map[string]*wisdom.Assistant {
	results := map[string]*wisdom.Assistant{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wisdom.Assistant:
			results[name] = resource
		}
	}
	return results
}

// GetWisdomAssistantWithName retrieves all wisdom.Assistant items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWisdomAssistantWithName(name string) (*wisdom.Assistant, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wisdom.Assistant:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wisdom.Assistant not found", name)
}

// GetAllWisdomAssistantAssociationResources retrieves all wisdom.AssistantAssociation items from an AWS CloudFormation template
func (t *Template) GetAllWisdomAssistantAssociationResources() map[string]*wisdom.AssistantAssociation {
	results := map[string]*wisdom.AssistantAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wisdom.AssistantAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetWisdomAssistantAssociationWithName retrieves all wisdom.AssistantAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWisdomAssistantAssociationWithName(name string) (*wisdom.AssistantAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wisdom.AssistantAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wisdom.AssistantAssociation not found", name)
}

// GetAllWisdomKnowledgeBaseResources retrieves all wisdom.KnowledgeBase items from an AWS CloudFormation template
func (t *Template) GetAllWisdomKnowledgeBaseResources() map[string]*wisdom.KnowledgeBase {
	results := map[string]*wisdom.KnowledgeBase{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *wisdom.KnowledgeBase:
			results[name] = resource
		}
	}
	return results
}

// GetWisdomKnowledgeBaseWithName retrieves all wisdom.KnowledgeBase items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWisdomKnowledgeBaseWithName(name string) (*wisdom.KnowledgeBase, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *wisdom.KnowledgeBase:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type wisdom.KnowledgeBase not found", name)
}

// GetAllWorkSpacesConnectionAliasResources retrieves all workspaces.ConnectionAlias items from an AWS CloudFormation template
func (t *Template) GetAllWorkSpacesConnectionAliasResources() map[string]*workspaces.ConnectionAlias {
	results := map[string]*workspaces.ConnectionAlias{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *workspaces.ConnectionAlias:
			results[name] = resource
		}
	}
	return results
}

// GetWorkSpacesConnectionAliasWithName retrieves all workspaces.ConnectionAlias items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWorkSpacesConnectionAliasWithName(name string) (*workspaces.ConnectionAlias, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *workspaces.ConnectionAlias:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type workspaces.ConnectionAlias not found", name)
}

// GetAllWorkSpacesWorkspaceResources retrieves all workspaces.Workspace items from an AWS CloudFormation template
func (t *Template) GetAllWorkSpacesWorkspaceResources() map[string]*workspaces.Workspace {
	results := map[string]*workspaces.Workspace{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *workspaces.Workspace:
			results[name] = resource
		}
	}
	return results
}

// GetWorkSpacesWorkspaceWithName retrieves all workspaces.Workspace items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetWorkSpacesWorkspaceWithName(name string) (*workspaces.Workspace, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *workspaces.Workspace:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type workspaces.Workspace not found", name)
}

// GetAllXRayGroupResources retrieves all xray.Group items from an AWS CloudFormation template
func (t *Template) GetAllXRayGroupResources() map[string]*xray.Group {
	results := map[string]*xray.Group{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *xray.Group:
			results[name] = resource
		}
	}
	return results
}

// GetXRayGroupWithName retrieves all xray.Group items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetXRayGroupWithName(name string) (*xray.Group, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *xray.Group:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type xray.Group not found", name)
}

// GetAllXRaySamplingRuleResources retrieves all xray.SamplingRule items from an AWS CloudFormation template
func (t *Template) GetAllXRaySamplingRuleResources() map[string]*xray.SamplingRule {
	results := map[string]*xray.SamplingRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *xray.SamplingRule:
			results[name] = resource
		}
	}
	return results
}

// GetXRaySamplingRuleWithName retrieves all xray.SamplingRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetXRaySamplingRuleWithName(name string) (*xray.SamplingRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *xray.SamplingRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type xray.SamplingRule not found", name)
}

// GetAllASKSkillResources retrieves all ask.Skill items from an AWS CloudFormation template
func (t *Template) GetAllASKSkillResources() map[string]*ask.Skill {
	results := map[string]*ask.Skill{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ask.Skill:
			results[name] = resource
		}
	}
	return results
}

// GetASKSkillWithName retrieves all ask.Skill items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetASKSkillWithName(name string) (*ask.Skill, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ask.Skill:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ask.Skill not found", name)
}
