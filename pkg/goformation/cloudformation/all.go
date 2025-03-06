package cloudformation

import (
	"fmt"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/autoscaling"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/cloudformation"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/cloudwatch"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/ec2"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/ecr"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/eks"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/elasticloadbalancing"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/elasticloadbalancingv2"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/events"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/iam"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/kinesis"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/kms"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/lambda"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/rds"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/rolesanywhere"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/route53"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/s3"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/serverless"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/sns"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/sqs"
)

// AllResources fetches an iterable map all CloudFormation and SAM resources
func AllResources() map[string]Resource {
	return map[string]Resource{
		"AWS::AutoScaling::AutoScalingGroup":                               &autoscaling.AutoScalingGroup{},
		"AWS::AutoScaling::LaunchConfiguration":                            &autoscaling.LaunchConfiguration{},
		"AWS::AutoScaling::LifecycleHook":                                  &autoscaling.LifecycleHook{},
		"AWS::AutoScaling::ScalingPolicy":                                  &autoscaling.ScalingPolicy{},
		"AWS::AutoScaling::ScheduledAction":                                &autoscaling.ScheduledAction{},
		"AWS::AutoScaling::WarmPool":                                       &autoscaling.WarmPool{},
		"AWS::CloudFormation::CustomResource":                              &cloudformation.CustomResource{},
		"AWS::CloudFormation::GuardHook":                                   &cloudformation.GuardHook{},
		"AWS::CloudFormation::HookDefaultVersion":                          &cloudformation.HookDefaultVersion{},
		"AWS::CloudFormation::HookTypeConfig":                              &cloudformation.HookTypeConfig{},
		"AWS::CloudFormation::HookVersion":                                 &cloudformation.HookVersion{},
		"AWS::CloudFormation::LambdaHook":                                  &cloudformation.LambdaHook{},
		"AWS::CloudFormation::Macro":                                       &cloudformation.Macro{},
		"AWS::CloudFormation::ModuleDefaultVersion":                        &cloudformation.ModuleDefaultVersion{},
		"AWS::CloudFormation::ModuleVersion":                               &cloudformation.ModuleVersion{},
		"AWS::CloudFormation::PublicTypeVersion":                           &cloudformation.PublicTypeVersion{},
		"AWS::CloudFormation::Publisher":                                   &cloudformation.Publisher{},
		"AWS::CloudFormation::ResourceDefaultVersion":                      &cloudformation.ResourceDefaultVersion{},
		"AWS::CloudFormation::ResourceVersion":                             &cloudformation.ResourceVersion{},
		"AWS::CloudFormation::Stack":                                       &cloudformation.Stack{},
		"AWS::CloudFormation::StackSet":                                    &cloudformation.StackSet{},
		"AWS::CloudFormation::TypeActivation":                              &cloudformation.TypeActivation{},
		"AWS::CloudFormation::WaitCondition":                               &cloudformation.WaitCondition{},
		"AWS::CloudFormation::WaitConditionHandle":                         &cloudformation.WaitConditionHandle{},
		"AWS::CloudWatch::Alarm":                                           &cloudwatch.Alarm{},
		"AWS::CloudWatch::AnomalyDetector":                                 &cloudwatch.AnomalyDetector{},
		"AWS::CloudWatch::CompositeAlarm":                                  &cloudwatch.CompositeAlarm{},
		"AWS::CloudWatch::Dashboard":                                       &cloudwatch.Dashboard{},
		"AWS::CloudWatch::InsightRule":                                     &cloudwatch.InsightRule{},
		"AWS::CloudWatch::MetricStream":                                    &cloudwatch.MetricStream{},
		"AWS::EC2::CapacityReservation":                                    &ec2.CapacityReservation{},
		"AWS::EC2::CapacityReservationFleet":                               &ec2.CapacityReservationFleet{},
		"AWS::EC2::CarrierGateway":                                         &ec2.CarrierGateway{},
		"AWS::EC2::ClientVpnAuthorizationRule":                             &ec2.ClientVpnAuthorizationRule{},
		"AWS::EC2::ClientVpnEndpoint":                                      &ec2.ClientVpnEndpoint{},
		"AWS::EC2::ClientVpnRoute":                                         &ec2.ClientVpnRoute{},
		"AWS::EC2::ClientVpnTargetNetworkAssociation":                      &ec2.ClientVpnTargetNetworkAssociation{},
		"AWS::EC2::CustomerGateway":                                        &ec2.CustomerGateway{},
		"AWS::EC2::DHCPOptions":                                            &ec2.DHCPOptions{},
		"AWS::EC2::EC2Fleet":                                               &ec2.EC2Fleet{},
		"AWS::EC2::EIP":                                                    &ec2.EIP{},
		"AWS::EC2::EIPAssociation":                                         &ec2.EIPAssociation{},
		"AWS::EC2::EgressOnlyInternetGateway":                              &ec2.EgressOnlyInternetGateway{},
		"AWS::EC2::EnclaveCertificateIamRoleAssociation":                   &ec2.EnclaveCertificateIamRoleAssociation{},
		"AWS::EC2::FlowLog":                                                &ec2.FlowLog{},
		"AWS::EC2::GatewayRouteTableAssociation":                           &ec2.GatewayRouteTableAssociation{},
		"AWS::EC2::Host":                                                   &ec2.Host{},
		"AWS::EC2::IPAM":                                                   &ec2.IPAM{},
		"AWS::EC2::IPAMAllocation":                                         &ec2.IPAMAllocation{},
		"AWS::EC2::IPAMPool":                                               &ec2.IPAMPool{},
		"AWS::EC2::IPAMPoolCidr":                                           &ec2.IPAMPoolCidr{},
		"AWS::EC2::IPAMResourceDiscovery":                                  &ec2.IPAMResourceDiscovery{},
		"AWS::EC2::IPAMResourceDiscoveryAssociation":                       &ec2.IPAMResourceDiscoveryAssociation{},
		"AWS::EC2::IPAMScope":                                              &ec2.IPAMScope{},
		"AWS::EC2::Instance":                                               &ec2.Instance{},
		"AWS::EC2::InstanceConnectEndpoint":                                &ec2.InstanceConnectEndpoint{},
		"AWS::EC2::InternetGateway":                                        &ec2.InternetGateway{},
		"AWS::EC2::KeyPair":                                                &ec2.KeyPair{},
		"AWS::EC2::LaunchTemplate":                                         &ec2.LaunchTemplate{},
		"AWS::EC2::LocalGatewayRoute":                                      &ec2.LocalGatewayRoute{},
		"AWS::EC2::LocalGatewayRouteTable":                                 &ec2.LocalGatewayRouteTable{},
		"AWS::EC2::LocalGatewayRouteTableVPCAssociation":                   &ec2.LocalGatewayRouteTableVPCAssociation{},
		"AWS::EC2::LocalGatewayRouteTableVirtualInterfaceGroupAssociation": &ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation{},
		"AWS::EC2::NatGateway":                                             &ec2.NatGateway{},
		"AWS::EC2::NetworkAcl":                                             &ec2.NetworkAcl{},
		"AWS::EC2::NetworkAclEntry":                                        &ec2.NetworkAclEntry{},
		"AWS::EC2::NetworkInsightsAccessScope":                             &ec2.NetworkInsightsAccessScope{},
		"AWS::EC2::NetworkInsightsAccessScopeAnalysis":                     &ec2.NetworkInsightsAccessScopeAnalysis{},
		"AWS::EC2::NetworkInsightsAnalysis":                                &ec2.NetworkInsightsAnalysis{},
		"AWS::EC2::NetworkInsightsPath":                                    &ec2.NetworkInsightsPath{},
		"AWS::EC2::NetworkInterface":                                       &ec2.NetworkInterface{},
		"AWS::EC2::NetworkInterfaceAttachment":                             &ec2.NetworkInterfaceAttachment{},
		"AWS::EC2::NetworkInterfacePermission":                             &ec2.NetworkInterfacePermission{},
		"AWS::EC2::NetworkPerformanceMetricSubscription":                   &ec2.NetworkPerformanceMetricSubscription{},
		"AWS::EC2::PlacementGroup":                                         &ec2.PlacementGroup{},
		"AWS::EC2::PrefixList":                                             &ec2.PrefixList{},
		"AWS::EC2::Route":                                                  &ec2.Route{},
		"AWS::EC2::RouteTable":                                             &ec2.RouteTable{},
		"AWS::EC2::SecurityGroup":                                          &ec2.SecurityGroup{},
		"AWS::EC2::SecurityGroupEgress":                                    &ec2.SecurityGroupEgress{},
		"AWS::EC2::SecurityGroupIngress":                                   &ec2.SecurityGroupIngress{},
		"AWS::EC2::SecurityGroupVpcAssociation":                            &ec2.SecurityGroupVpcAssociation{},
		"AWS::EC2::SnapshotBlockPublicAccess":                              &ec2.SnapshotBlockPublicAccess{},
		"AWS::EC2::SpotFleet":                                              &ec2.SpotFleet{},
		"AWS::EC2::Subnet":                                                 &ec2.Subnet{},
		"AWS::EC2::SubnetCidrBlock":                                        &ec2.SubnetCidrBlock{},
		"AWS::EC2::SubnetNetworkAclAssociation":                            &ec2.SubnetNetworkAclAssociation{},
		"AWS::EC2::SubnetRouteTableAssociation":                            &ec2.SubnetRouteTableAssociation{},
		"AWS::EC2::TrafficMirrorFilter":                                    &ec2.TrafficMirrorFilter{},
		"AWS::EC2::TrafficMirrorFilterRule":                                &ec2.TrafficMirrorFilterRule{},
		"AWS::EC2::TrafficMirrorSession":                                   &ec2.TrafficMirrorSession{},
		"AWS::EC2::TrafficMirrorTarget":                                    &ec2.TrafficMirrorTarget{},
		"AWS::EC2::TransitGateway":                                         &ec2.TransitGateway{},
		"AWS::EC2::TransitGatewayAttachment":                               &ec2.TransitGatewayAttachment{},
		"AWS::EC2::TransitGatewayConnect":                                  &ec2.TransitGatewayConnect{},
		"AWS::EC2::TransitGatewayMulticastDomain":                          &ec2.TransitGatewayMulticastDomain{},
		"AWS::EC2::TransitGatewayMulticastDomainAssociation":               &ec2.TransitGatewayMulticastDomainAssociation{},
		"AWS::EC2::TransitGatewayMulticastGroupMember":                     &ec2.TransitGatewayMulticastGroupMember{},
		"AWS::EC2::TransitGatewayMulticastGroupSource":                     &ec2.TransitGatewayMulticastGroupSource{},
		"AWS::EC2::TransitGatewayPeeringAttachment":                        &ec2.TransitGatewayPeeringAttachment{},
		"AWS::EC2::TransitGatewayRoute":                                    &ec2.TransitGatewayRoute{},
		"AWS::EC2::TransitGatewayRouteTable":                               &ec2.TransitGatewayRouteTable{},
		"AWS::EC2::TransitGatewayRouteTableAssociation":                    &ec2.TransitGatewayRouteTableAssociation{},
		"AWS::EC2::TransitGatewayRouteTablePropagation":                    &ec2.TransitGatewayRouteTablePropagation{},
		"AWS::EC2::TransitGatewayVpcAttachment":                            &ec2.TransitGatewayVpcAttachment{},
		"AWS::EC2::VPC":                                                    &ec2.VPC{},
		"AWS::EC2::VPCBlockPublicAccessExclusion":                          &ec2.VPCBlockPublicAccessExclusion{},
		"AWS::EC2::VPCBlockPublicAccessOptions":                            &ec2.VPCBlockPublicAccessOptions{},
		"AWS::EC2::VPCCidrBlock":                                           &ec2.VPCCidrBlock{},
		"AWS::EC2::VPCDHCPOptionsAssociation":                              &ec2.VPCDHCPOptionsAssociation{},
		"AWS::EC2::VPCEndpoint":                                            &ec2.VPCEndpoint{},
		"AWS::EC2::VPCEndpointConnectionNotification":                      &ec2.VPCEndpointConnectionNotification{},
		"AWS::EC2::VPCEndpointService":                                     &ec2.VPCEndpointService{},
		"AWS::EC2::VPCEndpointServicePermissions":                          &ec2.VPCEndpointServicePermissions{},
		"AWS::EC2::VPCGatewayAttachment":                                   &ec2.VPCGatewayAttachment{},
		"AWS::EC2::VPCPeeringConnection":                                   &ec2.VPCPeeringConnection{},
		"AWS::EC2::VPNConnection":                                          &ec2.VPNConnection{},
		"AWS::EC2::VPNConnectionRoute":                                     &ec2.VPNConnectionRoute{},
		"AWS::EC2::VPNGateway":                                             &ec2.VPNGateway{},
		"AWS::EC2::VPNGatewayRoutePropagation":                             &ec2.VPNGatewayRoutePropagation{},
		"AWS::EC2::VerifiedAccessEndpoint":                                 &ec2.VerifiedAccessEndpoint{},
		"AWS::EC2::VerifiedAccessGroup":                                    &ec2.VerifiedAccessGroup{},
		"AWS::EC2::VerifiedAccessInstance":                                 &ec2.VerifiedAccessInstance{},
		"AWS::EC2::VerifiedAccessTrustProvider":                            &ec2.VerifiedAccessTrustProvider{},
		"AWS::EC2::Volume":                                                 &ec2.Volume{},
		"AWS::EC2::VolumeAttachment":                                       &ec2.VolumeAttachment{},
		"AWS::ECR::PublicRepository":                                       &ecr.PublicRepository{},
		"AWS::ECR::PullThroughCacheRule":                                   &ecr.PullThroughCacheRule{},
		"AWS::ECR::RegistryPolicy":                                         &ecr.RegistryPolicy{},
		"AWS::ECR::ReplicationConfiguration":                               &ecr.ReplicationConfiguration{},
		"AWS::ECR::Repository":                                             &ecr.Repository{},
		"AWS::ECR::RepositoryCreationTemplate":                             &ecr.RepositoryCreationTemplate{},
		"AWS::EKS::AccessEntry":                                            &eks.AccessEntry{},
		"AWS::EKS::Addon":                                                  &eks.Addon{},
		"AWS::EKS::Cluster":                                                &eks.Cluster{},
		"AWS::EKS::FargateProfile":                                         &eks.FargateProfile{},
		"AWS::EKS::IdentityProviderConfig":                                 &eks.IdentityProviderConfig{},
		"AWS::EKS::Nodegroup":                                              &eks.Nodegroup{},
		"AWS::EKS::PodIdentityAssociation":                                 &eks.PodIdentityAssociation{},
		"AWS::ElasticLoadBalancing::LoadBalancer":                          &elasticloadbalancing.LoadBalancer{},
		"AWS::ElasticLoadBalancingV2::Listener":                            &elasticloadbalancingv2.Listener{},
		"AWS::ElasticLoadBalancingV2::ListenerCertificate":                 &elasticloadbalancingv2.ListenerCertificate{},
		"AWS::ElasticLoadBalancingV2::ListenerRule":                        &elasticloadbalancingv2.ListenerRule{},
		"AWS::ElasticLoadBalancingV2::LoadBalancer":                        &elasticloadbalancingv2.LoadBalancer{},
		"AWS::ElasticLoadBalancingV2::TargetGroup":                         &elasticloadbalancingv2.TargetGroup{},
		"AWS::ElasticLoadBalancingV2::TrustStore":                          &elasticloadbalancingv2.TrustStore{},
		"AWS::ElasticLoadBalancingV2::TrustStoreRevocation":                &elasticloadbalancingv2.TrustStoreRevocation{},
		"AWS::Events::ApiDestination":                                      &events.ApiDestination{},
		"AWS::Events::Archive":                                             &events.Archive{},
		"AWS::Events::Connection":                                          &events.Connection{},
		"AWS::Events::Endpoint":                                            &events.Endpoint{},
		"AWS::Events::EventBus":                                            &events.EventBus{},
		"AWS::Events::EventBusPolicy":                                      &events.EventBusPolicy{},
		"AWS::Events::Rule":                                                &events.Rule{},
		"AWS::IAM::AccessKey":                                              &iam.AccessKey{},
		"AWS::IAM::Group":                                                  &iam.Group{},
		"AWS::IAM::GroupPolicy":                                            &iam.GroupPolicy{},
		"AWS::IAM::InstanceProfile":                                        &iam.InstanceProfile{},
		"AWS::IAM::ManagedPolicy":                                          &iam.ManagedPolicy{},
		"AWS::IAM::OIDCProvider":                                           &iam.OIDCProvider{},
		"AWS::IAM::Policy":                                                 &iam.Policy{},
		"AWS::IAM::Role":                                                   &iam.Role{},
		"AWS::IAM::RolePolicy":                                             &iam.RolePolicy{},
		"AWS::IAM::SAMLProvider":                                           &iam.SAMLProvider{},
		"AWS::IAM::ServerCertificate":                                      &iam.ServerCertificate{},
		"AWS::IAM::ServiceLinkedRole":                                      &iam.ServiceLinkedRole{},
		"AWS::IAM::User":                                                   &iam.User{},
		"AWS::IAM::UserPolicy":                                             &iam.UserPolicy{},
		"AWS::IAM::UserToGroupAddition":                                    &iam.UserToGroupAddition{},
		"AWS::IAM::VirtualMFADevice":                                       &iam.VirtualMFADevice{},
		"AWS::KMS::Alias":                                                  &kms.Alias{},
		"AWS::KMS::Key":                                                    &kms.Key{},
		"AWS::KMS::ReplicaKey":                                             &kms.ReplicaKey{},
		"AWS::Kinesis::ResourcePolicy":                                     &kinesis.ResourcePolicy{},
		"AWS::Kinesis::Stream":                                             &kinesis.Stream{},
		"AWS::Kinesis::StreamConsumer":                                     &kinesis.StreamConsumer{},
		"AWS::Lambda::Alias":                                               &lambda.Alias{},
		"AWS::Lambda::CodeSigningConfig":                                   &lambda.CodeSigningConfig{},
		"AWS::Lambda::EventInvokeConfig":                                   &lambda.EventInvokeConfig{},
		"AWS::Lambda::EventSourceMapping":                                  &lambda.EventSourceMapping{},
		"AWS::Lambda::Function":                                            &lambda.Function{},
		"AWS::Lambda::LayerVersion":                                        &lambda.LayerVersion{},
		"AWS::Lambda::LayerVersionPermission":                              &lambda.LayerVersionPermission{},
		"AWS::Lambda::Permission":                                          &lambda.Permission{},
		"AWS::Lambda::Url":                                                 &lambda.Url{},
		"AWS::Lambda::Version":                                             &lambda.Version{},
		"AWS::RDS::CustomDBEngineVersion":                                  &rds.CustomDBEngineVersion{},
		"AWS::RDS::DBCluster":                                              &rds.DBCluster{},
		"AWS::RDS::DBClusterParameterGroup":                                &rds.DBClusterParameterGroup{},
		"AWS::RDS::DBInstance":                                             &rds.DBInstance{},
		"AWS::RDS::DBParameterGroup":                                       &rds.DBParameterGroup{},
		"AWS::RDS::DBProxy":                                                &rds.DBProxy{},
		"AWS::RDS::DBProxyEndpoint":                                        &rds.DBProxyEndpoint{},
		"AWS::RDS::DBProxyTargetGroup":                                     &rds.DBProxyTargetGroup{},
		"AWS::RDS::DBSecurityGroup":                                        &rds.DBSecurityGroup{},
		"AWS::RDS::DBSecurityGroupIngress":                                 &rds.DBSecurityGroupIngress{},
		"AWS::RDS::DBShardGroup":                                           &rds.DBShardGroup{},
		"AWS::RDS::DBSubnetGroup":                                          &rds.DBSubnetGroup{},
		"AWS::RDS::EventSubscription":                                      &rds.EventSubscription{},
		"AWS::RDS::GlobalCluster":                                          &rds.GlobalCluster{},
		"AWS::RDS::Integration":                                            &rds.Integration{},
		"AWS::RDS::OptionGroup":                                            &rds.OptionGroup{},
		"AWS::RolesAnywhere::CRL":                                          &rolesanywhere.CRL{},
		"AWS::RolesAnywhere::Profile":                                      &rolesanywhere.Profile{},
		"AWS::RolesAnywhere::TrustAnchor":                                  &rolesanywhere.TrustAnchor{},
		"AWS::Route53::CidrCollection":                                     &route53.CidrCollection{},
		"AWS::Route53::DNSSEC":                                             &route53.DNSSEC{},
		"AWS::Route53::HealthCheck":                                        &route53.HealthCheck{},
		"AWS::Route53::HostedZone":                                         &route53.HostedZone{},
		"AWS::Route53::KeySigningKey":                                      &route53.KeySigningKey{},
		"AWS::Route53::RecordSet":                                          &route53.RecordSet{},
		"AWS::Route53::RecordSetGroup":                                     &route53.RecordSetGroup{},
		"AWS::S3::AccessGrant":                                             &s3.AccessGrant{},
		"AWS::S3::AccessGrantsInstance":                                    &s3.AccessGrantsInstance{},
		"AWS::S3::AccessGrantsLocation":                                    &s3.AccessGrantsLocation{},
		"AWS::S3::AccessPoint":                                             &s3.AccessPoint{},
		"AWS::S3::Bucket":                                                  &s3.Bucket{},
		"AWS::S3::BucketPolicy":                                            &s3.BucketPolicy{},
		"AWS::S3::MultiRegionAccessPoint":                                  &s3.MultiRegionAccessPoint{},
		"AWS::S3::MultiRegionAccessPointPolicy":                            &s3.MultiRegionAccessPointPolicy{},
		"AWS::S3::StorageLens":                                             &s3.StorageLens{},
		"AWS::S3::StorageLensGroup":                                        &s3.StorageLensGroup{},
		"AWS::SNS::Subscription":                                           &sns.Subscription{},
		"AWS::SNS::Topic":                                                  &sns.Topic{},
		"AWS::SNS::TopicInlinePolicy":                                      &sns.TopicInlinePolicy{},
		"AWS::SNS::TopicPolicy":                                            &sns.TopicPolicy{},
		"AWS::SQS::Queue":                                                  &sqs.Queue{},
		"AWS::SQS::QueueInlinePolicy":                                      &sqs.QueueInlinePolicy{},
		"AWS::SQS::QueuePolicy":                                            &sqs.QueuePolicy{},
		"AWS::Serverless::Api":                                             &serverless.Api{},
		"AWS::Serverless::Application":                                     &serverless.Application{},
		"AWS::Serverless::Function":                                        &serverless.Function{},
		"AWS::Serverless::LayerVersion":                                    &serverless.LayerVersion{},
		"AWS::Serverless::SimpleTable":                                     &serverless.SimpleTable{},
		"AWS::Serverless::StateMachine":                                    &serverless.StateMachine{},
	}
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

// GetAllCloudFormationGuardHookResources retrieves all cloudformation.GuardHook items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationGuardHookResources() map[string]*cloudformation.GuardHook {
	results := map[string]*cloudformation.GuardHook{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.GuardHook:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationGuardHookWithName retrieves all cloudformation.GuardHook items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationGuardHookWithName(name string) (*cloudformation.GuardHook, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.GuardHook:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.GuardHook not found", name)
}

// GetAllCloudFormationHookDefaultVersionResources retrieves all cloudformation.HookDefaultVersion items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationHookDefaultVersionResources() map[string]*cloudformation.HookDefaultVersion {
	results := map[string]*cloudformation.HookDefaultVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.HookDefaultVersion:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationHookDefaultVersionWithName retrieves all cloudformation.HookDefaultVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationHookDefaultVersionWithName(name string) (*cloudformation.HookDefaultVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.HookDefaultVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.HookDefaultVersion not found", name)
}

// GetAllCloudFormationHookTypeConfigResources retrieves all cloudformation.HookTypeConfig items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationHookTypeConfigResources() map[string]*cloudformation.HookTypeConfig {
	results := map[string]*cloudformation.HookTypeConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.HookTypeConfig:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationHookTypeConfigWithName retrieves all cloudformation.HookTypeConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationHookTypeConfigWithName(name string) (*cloudformation.HookTypeConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.HookTypeConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.HookTypeConfig not found", name)
}

// GetAllCloudFormationHookVersionResources retrieves all cloudformation.HookVersion items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationHookVersionResources() map[string]*cloudformation.HookVersion {
	results := map[string]*cloudformation.HookVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.HookVersion:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationHookVersionWithName retrieves all cloudformation.HookVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationHookVersionWithName(name string) (*cloudformation.HookVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.HookVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.HookVersion not found", name)
}

// GetAllCloudFormationLambdaHookResources retrieves all cloudformation.LambdaHook items from an AWS CloudFormation template
func (t *Template) GetAllCloudFormationLambdaHookResources() map[string]*cloudformation.LambdaHook {
	results := map[string]*cloudformation.LambdaHook{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *cloudformation.LambdaHook:
			results[name] = resource
		}
	}
	return results
}

// GetCloudFormationLambdaHookWithName retrieves all cloudformation.LambdaHook items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetCloudFormationLambdaHookWithName(name string) (*cloudformation.LambdaHook, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *cloudformation.LambdaHook:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type cloudformation.LambdaHook not found", name)
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

// GetAllEC2IPAMPoolCidrResources retrieves all ec2.IPAMPoolCidr items from an AWS CloudFormation template
func (t *Template) GetAllEC2IPAMPoolCidrResources() map[string]*ec2.IPAMPoolCidr {
	results := map[string]*ec2.IPAMPoolCidr{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.IPAMPoolCidr:
			results[name] = resource
		}
	}
	return results
}

// GetEC2IPAMPoolCidrWithName retrieves all ec2.IPAMPoolCidr items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2IPAMPoolCidrWithName(name string) (*ec2.IPAMPoolCidr, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.IPAMPoolCidr:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.IPAMPoolCidr not found", name)
}

// GetAllEC2IPAMResourceDiscoveryResources retrieves all ec2.IPAMResourceDiscovery items from an AWS CloudFormation template
func (t *Template) GetAllEC2IPAMResourceDiscoveryResources() map[string]*ec2.IPAMResourceDiscovery {
	results := map[string]*ec2.IPAMResourceDiscovery{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.IPAMResourceDiscovery:
			results[name] = resource
		}
	}
	return results
}

// GetEC2IPAMResourceDiscoveryWithName retrieves all ec2.IPAMResourceDiscovery items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2IPAMResourceDiscoveryWithName(name string) (*ec2.IPAMResourceDiscovery, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.IPAMResourceDiscovery:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.IPAMResourceDiscovery not found", name)
}

// GetAllEC2IPAMResourceDiscoveryAssociationResources retrieves all ec2.IPAMResourceDiscoveryAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2IPAMResourceDiscoveryAssociationResources() map[string]*ec2.IPAMResourceDiscoveryAssociation {
	results := map[string]*ec2.IPAMResourceDiscoveryAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.IPAMResourceDiscoveryAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2IPAMResourceDiscoveryAssociationWithName retrieves all ec2.IPAMResourceDiscoveryAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2IPAMResourceDiscoveryAssociationWithName(name string) (*ec2.IPAMResourceDiscoveryAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.IPAMResourceDiscoveryAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.IPAMResourceDiscoveryAssociation not found", name)
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

// GetAllEC2InstanceConnectEndpointResources retrieves all ec2.InstanceConnectEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllEC2InstanceConnectEndpointResources() map[string]*ec2.InstanceConnectEndpoint {
	results := map[string]*ec2.InstanceConnectEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.InstanceConnectEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetEC2InstanceConnectEndpointWithName retrieves all ec2.InstanceConnectEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2InstanceConnectEndpointWithName(name string) (*ec2.InstanceConnectEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.InstanceConnectEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.InstanceConnectEndpoint not found", name)
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

// GetAllEC2KeyPairResources retrieves all ec2.KeyPair items from an AWS CloudFormation template
func (t *Template) GetAllEC2KeyPairResources() map[string]*ec2.KeyPair {
	results := map[string]*ec2.KeyPair{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.KeyPair:
			results[name] = resource
		}
	}
	return results
}

// GetEC2KeyPairWithName retrieves all ec2.KeyPair items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2KeyPairWithName(name string) (*ec2.KeyPair, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.KeyPair:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.KeyPair not found", name)
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

// GetAllEC2LocalGatewayRouteTableResources retrieves all ec2.LocalGatewayRouteTable items from an AWS CloudFormation template
func (t *Template) GetAllEC2LocalGatewayRouteTableResources() map[string]*ec2.LocalGatewayRouteTable {
	results := map[string]*ec2.LocalGatewayRouteTable{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRouteTable:
			results[name] = resource
		}
	}
	return results
}

// GetEC2LocalGatewayRouteTableWithName retrieves all ec2.LocalGatewayRouteTable items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2LocalGatewayRouteTableWithName(name string) (*ec2.LocalGatewayRouteTable, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRouteTable:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.LocalGatewayRouteTable not found", name)
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

// GetAllEC2LocalGatewayRouteTableVirtualInterfaceGroupAssociationResources retrieves all ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2LocalGatewayRouteTableVirtualInterfaceGroupAssociationResources() map[string]*ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation {
	results := map[string]*ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2LocalGatewayRouteTableVirtualInterfaceGroupAssociationWithName retrieves all ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2LocalGatewayRouteTableVirtualInterfaceGroupAssociationWithName(name string) (*ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.LocalGatewayRouteTableVirtualInterfaceGroupAssociation not found", name)
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

// GetAllEC2NetworkInsightsAccessScopeResources retrieves all ec2.NetworkInsightsAccessScope items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkInsightsAccessScopeResources() map[string]*ec2.NetworkInsightsAccessScope {
	results := map[string]*ec2.NetworkInsightsAccessScope{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsAccessScope:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkInsightsAccessScopeWithName retrieves all ec2.NetworkInsightsAccessScope items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkInsightsAccessScopeWithName(name string) (*ec2.NetworkInsightsAccessScope, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsAccessScope:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkInsightsAccessScope not found", name)
}

// GetAllEC2NetworkInsightsAccessScopeAnalysisResources retrieves all ec2.NetworkInsightsAccessScopeAnalysis items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkInsightsAccessScopeAnalysisResources() map[string]*ec2.NetworkInsightsAccessScopeAnalysis {
	results := map[string]*ec2.NetworkInsightsAccessScopeAnalysis{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsAccessScopeAnalysis:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkInsightsAccessScopeAnalysisWithName retrieves all ec2.NetworkInsightsAccessScopeAnalysis items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkInsightsAccessScopeAnalysisWithName(name string) (*ec2.NetworkInsightsAccessScopeAnalysis, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkInsightsAccessScopeAnalysis:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkInsightsAccessScopeAnalysis not found", name)
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

// GetAllEC2NetworkPerformanceMetricSubscriptionResources retrieves all ec2.NetworkPerformanceMetricSubscription items from an AWS CloudFormation template
func (t *Template) GetAllEC2NetworkPerformanceMetricSubscriptionResources() map[string]*ec2.NetworkPerformanceMetricSubscription {
	results := map[string]*ec2.NetworkPerformanceMetricSubscription{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.NetworkPerformanceMetricSubscription:
			results[name] = resource
		}
	}
	return results
}

// GetEC2NetworkPerformanceMetricSubscriptionWithName retrieves all ec2.NetworkPerformanceMetricSubscription items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2NetworkPerformanceMetricSubscriptionWithName(name string) (*ec2.NetworkPerformanceMetricSubscription, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.NetworkPerformanceMetricSubscription:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.NetworkPerformanceMetricSubscription not found", name)
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

// GetAllEC2SecurityGroupVpcAssociationResources retrieves all ec2.SecurityGroupVpcAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEC2SecurityGroupVpcAssociationResources() map[string]*ec2.SecurityGroupVpcAssociation {
	results := map[string]*ec2.SecurityGroupVpcAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroupVpcAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SecurityGroupVpcAssociationWithName retrieves all ec2.SecurityGroupVpcAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SecurityGroupVpcAssociationWithName(name string) (*ec2.SecurityGroupVpcAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SecurityGroupVpcAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SecurityGroupVpcAssociation not found", name)
}

// GetAllEC2SnapshotBlockPublicAccessResources retrieves all ec2.SnapshotBlockPublicAccess items from an AWS CloudFormation template
func (t *Template) GetAllEC2SnapshotBlockPublicAccessResources() map[string]*ec2.SnapshotBlockPublicAccess {
	results := map[string]*ec2.SnapshotBlockPublicAccess{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.SnapshotBlockPublicAccess:
			results[name] = resource
		}
	}
	return results
}

// GetEC2SnapshotBlockPublicAccessWithName retrieves all ec2.SnapshotBlockPublicAccess items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2SnapshotBlockPublicAccessWithName(name string) (*ec2.SnapshotBlockPublicAccess, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.SnapshotBlockPublicAccess:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.SnapshotBlockPublicAccess not found", name)
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

// GetAllEC2VPCBlockPublicAccessExclusionResources retrieves all ec2.VPCBlockPublicAccessExclusion items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCBlockPublicAccessExclusionResources() map[string]*ec2.VPCBlockPublicAccessExclusion {
	results := map[string]*ec2.VPCBlockPublicAccessExclusion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCBlockPublicAccessExclusion:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCBlockPublicAccessExclusionWithName retrieves all ec2.VPCBlockPublicAccessExclusion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCBlockPublicAccessExclusionWithName(name string) (*ec2.VPCBlockPublicAccessExclusion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCBlockPublicAccessExclusion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCBlockPublicAccessExclusion not found", name)
}

// GetAllEC2VPCBlockPublicAccessOptionsResources retrieves all ec2.VPCBlockPublicAccessOptions items from an AWS CloudFormation template
func (t *Template) GetAllEC2VPCBlockPublicAccessOptionsResources() map[string]*ec2.VPCBlockPublicAccessOptions {
	results := map[string]*ec2.VPCBlockPublicAccessOptions{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VPCBlockPublicAccessOptions:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VPCBlockPublicAccessOptionsWithName retrieves all ec2.VPCBlockPublicAccessOptions items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VPCBlockPublicAccessOptionsWithName(name string) (*ec2.VPCBlockPublicAccessOptions, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VPCBlockPublicAccessOptions:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VPCBlockPublicAccessOptions not found", name)
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

// GetAllEC2VerifiedAccessEndpointResources retrieves all ec2.VerifiedAccessEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllEC2VerifiedAccessEndpointResources() map[string]*ec2.VerifiedAccessEndpoint {
	results := map[string]*ec2.VerifiedAccessEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessEndpoint:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VerifiedAccessEndpointWithName retrieves all ec2.VerifiedAccessEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VerifiedAccessEndpointWithName(name string) (*ec2.VerifiedAccessEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessEndpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VerifiedAccessEndpoint not found", name)
}

// GetAllEC2VerifiedAccessGroupResources retrieves all ec2.VerifiedAccessGroup items from an AWS CloudFormation template
func (t *Template) GetAllEC2VerifiedAccessGroupResources() map[string]*ec2.VerifiedAccessGroup {
	results := map[string]*ec2.VerifiedAccessGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessGroup:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VerifiedAccessGroupWithName retrieves all ec2.VerifiedAccessGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VerifiedAccessGroupWithName(name string) (*ec2.VerifiedAccessGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VerifiedAccessGroup not found", name)
}

// GetAllEC2VerifiedAccessInstanceResources retrieves all ec2.VerifiedAccessInstance items from an AWS CloudFormation template
func (t *Template) GetAllEC2VerifiedAccessInstanceResources() map[string]*ec2.VerifiedAccessInstance {
	results := map[string]*ec2.VerifiedAccessInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessInstance:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VerifiedAccessInstanceWithName retrieves all ec2.VerifiedAccessInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VerifiedAccessInstanceWithName(name string) (*ec2.VerifiedAccessInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VerifiedAccessInstance not found", name)
}

// GetAllEC2VerifiedAccessTrustProviderResources retrieves all ec2.VerifiedAccessTrustProvider items from an AWS CloudFormation template
func (t *Template) GetAllEC2VerifiedAccessTrustProviderResources() map[string]*ec2.VerifiedAccessTrustProvider {
	results := map[string]*ec2.VerifiedAccessTrustProvider{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessTrustProvider:
			results[name] = resource
		}
	}
	return results
}

// GetEC2VerifiedAccessTrustProviderWithName retrieves all ec2.VerifiedAccessTrustProvider items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEC2VerifiedAccessTrustProviderWithName(name string) (*ec2.VerifiedAccessTrustProvider, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ec2.VerifiedAccessTrustProvider:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ec2.VerifiedAccessTrustProvider not found", name)
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

// GetAllECRPullThroughCacheRuleResources retrieves all ecr.PullThroughCacheRule items from an AWS CloudFormation template
func (t *Template) GetAllECRPullThroughCacheRuleResources() map[string]*ecr.PullThroughCacheRule {
	results := map[string]*ecr.PullThroughCacheRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecr.PullThroughCacheRule:
			results[name] = resource
		}
	}
	return results
}

// GetECRPullThroughCacheRuleWithName retrieves all ecr.PullThroughCacheRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECRPullThroughCacheRuleWithName(name string) (*ecr.PullThroughCacheRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecr.PullThroughCacheRule:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecr.PullThroughCacheRule not found", name)
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

// GetAllECRRepositoryCreationTemplateResources retrieves all ecr.RepositoryCreationTemplate items from an AWS CloudFormation template
func (t *Template) GetAllECRRepositoryCreationTemplateResources() map[string]*ecr.RepositoryCreationTemplate {
	results := map[string]*ecr.RepositoryCreationTemplate{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *ecr.RepositoryCreationTemplate:
			results[name] = resource
		}
	}
	return results
}

// GetECRRepositoryCreationTemplateWithName retrieves all ecr.RepositoryCreationTemplate items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetECRRepositoryCreationTemplateWithName(name string) (*ecr.RepositoryCreationTemplate, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *ecr.RepositoryCreationTemplate:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type ecr.RepositoryCreationTemplate not found", name)
}

// GetAllEKSAccessEntryResources retrieves all eks.AccessEntry items from an AWS CloudFormation template
func (t *Template) GetAllEKSAccessEntryResources() map[string]*eks.AccessEntry {
	results := map[string]*eks.AccessEntry{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eks.AccessEntry:
			results[name] = resource
		}
	}
	return results
}

// GetEKSAccessEntryWithName retrieves all eks.AccessEntry items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEKSAccessEntryWithName(name string) (*eks.AccessEntry, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eks.AccessEntry:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eks.AccessEntry not found", name)
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

// GetAllEKSIdentityProviderConfigResources retrieves all eks.IdentityProviderConfig items from an AWS CloudFormation template
func (t *Template) GetAllEKSIdentityProviderConfigResources() map[string]*eks.IdentityProviderConfig {
	results := map[string]*eks.IdentityProviderConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eks.IdentityProviderConfig:
			results[name] = resource
		}
	}
	return results
}

// GetEKSIdentityProviderConfigWithName retrieves all eks.IdentityProviderConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEKSIdentityProviderConfigWithName(name string) (*eks.IdentityProviderConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eks.IdentityProviderConfig:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eks.IdentityProviderConfig not found", name)
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

// GetAllEKSPodIdentityAssociationResources retrieves all eks.PodIdentityAssociation items from an AWS CloudFormation template
func (t *Template) GetAllEKSPodIdentityAssociationResources() map[string]*eks.PodIdentityAssociation {
	results := map[string]*eks.PodIdentityAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *eks.PodIdentityAssociation:
			results[name] = resource
		}
	}
	return results
}

// GetEKSPodIdentityAssociationWithName retrieves all eks.PodIdentityAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEKSPodIdentityAssociationWithName(name string) (*eks.PodIdentityAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *eks.PodIdentityAssociation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type eks.PodIdentityAssociation not found", name)
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

// GetAllElasticLoadBalancingV2TrustStoreResources retrieves all elasticloadbalancingv2.TrustStore items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingV2TrustStoreResources() map[string]*elasticloadbalancingv2.TrustStore {
	results := map[string]*elasticloadbalancingv2.TrustStore{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.TrustStore:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingV2TrustStoreWithName retrieves all elasticloadbalancingv2.TrustStore items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingV2TrustStoreWithName(name string) (*elasticloadbalancingv2.TrustStore, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.TrustStore:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancingv2.TrustStore not found", name)
}

// GetAllElasticLoadBalancingV2TrustStoreRevocationResources retrieves all elasticloadbalancingv2.TrustStoreRevocation items from an AWS CloudFormation template
func (t *Template) GetAllElasticLoadBalancingV2TrustStoreRevocationResources() map[string]*elasticloadbalancingv2.TrustStoreRevocation {
	results := map[string]*elasticloadbalancingv2.TrustStoreRevocation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.TrustStoreRevocation:
			results[name] = resource
		}
	}
	return results
}

// GetElasticLoadBalancingV2TrustStoreRevocationWithName retrieves all elasticloadbalancingv2.TrustStoreRevocation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetElasticLoadBalancingV2TrustStoreRevocationWithName(name string) (*elasticloadbalancingv2.TrustStoreRevocation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *elasticloadbalancingv2.TrustStoreRevocation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type elasticloadbalancingv2.TrustStoreRevocation not found", name)
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

// GetAllEventsEndpointResources retrieves all events.Endpoint items from an AWS CloudFormation template
func (t *Template) GetAllEventsEndpointResources() map[string]*events.Endpoint {
	results := map[string]*events.Endpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *events.Endpoint:
			results[name] = resource
		}
	}
	return results
}

// GetEventsEndpointWithName retrieves all events.Endpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetEventsEndpointWithName(name string) (*events.Endpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *events.Endpoint:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type events.Endpoint not found", name)
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

// GetAllIAMGroupPolicyResources retrieves all iam.GroupPolicy items from an AWS CloudFormation template
func (t *Template) GetAllIAMGroupPolicyResources() map[string]*iam.GroupPolicy {
	results := map[string]*iam.GroupPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.GroupPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetIAMGroupPolicyWithName retrieves all iam.GroupPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMGroupPolicyWithName(name string) (*iam.GroupPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.GroupPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.GroupPolicy not found", name)
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

// GetAllIAMRolePolicyResources retrieves all iam.RolePolicy items from an AWS CloudFormation template
func (t *Template) GetAllIAMRolePolicyResources() map[string]*iam.RolePolicy {
	results := map[string]*iam.RolePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.RolePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetIAMRolePolicyWithName retrieves all iam.RolePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMRolePolicyWithName(name string) (*iam.RolePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.RolePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.RolePolicy not found", name)
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

// GetAllIAMUserPolicyResources retrieves all iam.UserPolicy items from an AWS CloudFormation template
func (t *Template) GetAllIAMUserPolicyResources() map[string]*iam.UserPolicy {
	results := map[string]*iam.UserPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *iam.UserPolicy:
			results[name] = resource
		}
	}
	return results
}

// GetIAMUserPolicyWithName retrieves all iam.UserPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetIAMUserPolicyWithName(name string) (*iam.UserPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *iam.UserPolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type iam.UserPolicy not found", name)
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

// GetAllKinesisResourcePolicyResources retrieves all kinesis.ResourcePolicy items from an AWS CloudFormation template
func (t *Template) GetAllKinesisResourcePolicyResources() map[string]*kinesis.ResourcePolicy {
	results := map[string]*kinesis.ResourcePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *kinesis.ResourcePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetKinesisResourcePolicyWithName retrieves all kinesis.ResourcePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetKinesisResourcePolicyWithName(name string) (*kinesis.ResourcePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *kinesis.ResourcePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type kinesis.ResourcePolicy not found", name)
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

// GetAllLambdaUrlResources retrieves all lambda.Url items from an AWS CloudFormation template
func (t *Template) GetAllLambdaUrlResources() map[string]*lambda.Url {
	results := map[string]*lambda.Url{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *lambda.Url:
			results[name] = resource
		}
	}
	return results
}

// GetLambdaUrlWithName retrieves all lambda.Url items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetLambdaUrlWithName(name string) (*lambda.Url, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *lambda.Url:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type lambda.Url not found", name)
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

// GetAllRDSCustomDBEngineVersionResources retrieves all rds.CustomDBEngineVersion items from an AWS CloudFormation template
func (t *Template) GetAllRDSCustomDBEngineVersionResources() map[string]*rds.CustomDBEngineVersion {
	results := map[string]*rds.CustomDBEngineVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.CustomDBEngineVersion:
			results[name] = resource
		}
	}
	return results
}

// GetRDSCustomDBEngineVersionWithName retrieves all rds.CustomDBEngineVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSCustomDBEngineVersionWithName(name string) (*rds.CustomDBEngineVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.CustomDBEngineVersion:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.CustomDBEngineVersion not found", name)
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

// GetAllRDSDBShardGroupResources retrieves all rds.DBShardGroup items from an AWS CloudFormation template
func (t *Template) GetAllRDSDBShardGroupResources() map[string]*rds.DBShardGroup {
	results := map[string]*rds.DBShardGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.DBShardGroup:
			results[name] = resource
		}
	}
	return results
}

// GetRDSDBShardGroupWithName retrieves all rds.DBShardGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSDBShardGroupWithName(name string) (*rds.DBShardGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.DBShardGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.DBShardGroup not found", name)
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

// GetAllRDSIntegrationResources retrieves all rds.Integration items from an AWS CloudFormation template
func (t *Template) GetAllRDSIntegrationResources() map[string]*rds.Integration {
	results := map[string]*rds.Integration{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rds.Integration:
			results[name] = resource
		}
	}
	return results
}

// GetRDSIntegrationWithName retrieves all rds.Integration items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRDSIntegrationWithName(name string) (*rds.Integration, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rds.Integration:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rds.Integration not found", name)
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

// GetAllRolesAnywhereCRLResources retrieves all rolesanywhere.CRL items from an AWS CloudFormation template
func (t *Template) GetAllRolesAnywhereCRLResources() map[string]*rolesanywhere.CRL {
	results := map[string]*rolesanywhere.CRL{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rolesanywhere.CRL:
			results[name] = resource
		}
	}
	return results
}

// GetRolesAnywhereCRLWithName retrieves all rolesanywhere.CRL items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRolesAnywhereCRLWithName(name string) (*rolesanywhere.CRL, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rolesanywhere.CRL:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rolesanywhere.CRL not found", name)
}

// GetAllRolesAnywhereProfileResources retrieves all rolesanywhere.Profile items from an AWS CloudFormation template
func (t *Template) GetAllRolesAnywhereProfileResources() map[string]*rolesanywhere.Profile {
	results := map[string]*rolesanywhere.Profile{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rolesanywhere.Profile:
			results[name] = resource
		}
	}
	return results
}

// GetRolesAnywhereProfileWithName retrieves all rolesanywhere.Profile items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRolesAnywhereProfileWithName(name string) (*rolesanywhere.Profile, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rolesanywhere.Profile:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rolesanywhere.Profile not found", name)
}

// GetAllRolesAnywhereTrustAnchorResources retrieves all rolesanywhere.TrustAnchor items from an AWS CloudFormation template
func (t *Template) GetAllRolesAnywhereTrustAnchorResources() map[string]*rolesanywhere.TrustAnchor {
	results := map[string]*rolesanywhere.TrustAnchor{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *rolesanywhere.TrustAnchor:
			results[name] = resource
		}
	}
	return results
}

// GetRolesAnywhereTrustAnchorWithName retrieves all rolesanywhere.TrustAnchor items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRolesAnywhereTrustAnchorWithName(name string) (*rolesanywhere.TrustAnchor, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *rolesanywhere.TrustAnchor:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type rolesanywhere.TrustAnchor not found", name)
}

// GetAllRoute53CidrCollectionResources retrieves all route53.CidrCollection items from an AWS CloudFormation template
func (t *Template) GetAllRoute53CidrCollectionResources() map[string]*route53.CidrCollection {
	results := map[string]*route53.CidrCollection{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *route53.CidrCollection:
			results[name] = resource
		}
	}
	return results
}

// GetRoute53CidrCollectionWithName retrieves all route53.CidrCollection items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetRoute53CidrCollectionWithName(name string) (*route53.CidrCollection, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *route53.CidrCollection:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type route53.CidrCollection not found", name)
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

// GetAllS3AccessGrantResources retrieves all s3.AccessGrant items from an AWS CloudFormation template
func (t *Template) GetAllS3AccessGrantResources() map[string]*s3.AccessGrant {
	results := map[string]*s3.AccessGrant{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.AccessGrant:
			results[name] = resource
		}
	}
	return results
}

// GetS3AccessGrantWithName retrieves all s3.AccessGrant items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3AccessGrantWithName(name string) (*s3.AccessGrant, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.AccessGrant:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.AccessGrant not found", name)
}

// GetAllS3AccessGrantsInstanceResources retrieves all s3.AccessGrantsInstance items from an AWS CloudFormation template
func (t *Template) GetAllS3AccessGrantsInstanceResources() map[string]*s3.AccessGrantsInstance {
	results := map[string]*s3.AccessGrantsInstance{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.AccessGrantsInstance:
			results[name] = resource
		}
	}
	return results
}

// GetS3AccessGrantsInstanceWithName retrieves all s3.AccessGrantsInstance items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3AccessGrantsInstanceWithName(name string) (*s3.AccessGrantsInstance, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.AccessGrantsInstance:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.AccessGrantsInstance not found", name)
}

// GetAllS3AccessGrantsLocationResources retrieves all s3.AccessGrantsLocation items from an AWS CloudFormation template
func (t *Template) GetAllS3AccessGrantsLocationResources() map[string]*s3.AccessGrantsLocation {
	results := map[string]*s3.AccessGrantsLocation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.AccessGrantsLocation:
			results[name] = resource
		}
	}
	return results
}

// GetS3AccessGrantsLocationWithName retrieves all s3.AccessGrantsLocation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3AccessGrantsLocationWithName(name string) (*s3.AccessGrantsLocation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.AccessGrantsLocation:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.AccessGrantsLocation not found", name)
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

// GetAllS3StorageLensGroupResources retrieves all s3.StorageLensGroup items from an AWS CloudFormation template
func (t *Template) GetAllS3StorageLensGroupResources() map[string]*s3.StorageLensGroup {
	results := map[string]*s3.StorageLensGroup{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *s3.StorageLensGroup:
			results[name] = resource
		}
	}
	return results
}

// GetS3StorageLensGroupWithName retrieves all s3.StorageLensGroup items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetS3StorageLensGroupWithName(name string) (*s3.StorageLensGroup, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *s3.StorageLensGroup:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type s3.StorageLensGroup not found", name)
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

// GetAllSNSTopicInlinePolicyResources retrieves all sns.TopicInlinePolicy items from an AWS CloudFormation template
func (t *Template) GetAllSNSTopicInlinePolicyResources() map[string]*sns.TopicInlinePolicy {
	results := map[string]*sns.TopicInlinePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sns.TopicInlinePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetSNSTopicInlinePolicyWithName retrieves all sns.TopicInlinePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSNSTopicInlinePolicyWithName(name string) (*sns.TopicInlinePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sns.TopicInlinePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sns.TopicInlinePolicy not found", name)
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

// GetAllSQSQueueInlinePolicyResources retrieves all sqs.QueueInlinePolicy items from an AWS CloudFormation template
func (t *Template) GetAllSQSQueueInlinePolicyResources() map[string]*sqs.QueueInlinePolicy {
	results := map[string]*sqs.QueueInlinePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case *sqs.QueueInlinePolicy:
			results[name] = resource
		}
	}
	return results
}

// GetSQSQueueInlinePolicyWithName retrieves all sqs.QueueInlinePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetSQSQueueInlinePolicyWithName(name string) (*sqs.QueueInlinePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case *sqs.QueueInlinePolicy:
			return resource, nil
		}
	}
	return nil, fmt.Errorf("resource %q of type sqs.QueueInlinePolicy not found", name)
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
