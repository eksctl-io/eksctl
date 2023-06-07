package mockprovider

import (
	"context"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5/fakes"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

// ProviderConfig holds current global config
var ProviderConfig = &api.ProviderConfig{
	Region: api.DefaultRegion,
	Profile: api.Profile{
		Name: "default",
	},
	WaitTimeout: 1200000000000,
}

var _ api.ClusterProvider = &MockProvider{}

type MockInput struct{}
type MockOutput struct {
	States []*MockState
}
type MockState struct {
	State *string
}

// MockProvider stores the mocked APIs
type MockProvider struct {
	waitTimeout         *time.Duration
	region              string
	cfnRoleARN          string
	asg                 *mocksv2.ASG
	eks                 *mocksv2.EKS
	cloudtrail          *mocksv2.CloudTrail
	cloudwatchlogs      *mocksv2.CloudWatchLogs
	credentialsProvider *mocksv2.CredentialsProvider

	cfn          *mocksv2.CloudFormation
	sts          *mocksv2.STS
	stsPresigner api.STSPresigner
	elb          *mocksv2.ELB
	elbV2        *mocksv2.ELBV2
	ssm          *mocksv2.SSM
	iam          *mocksv2.IAM
	ec2          *mocksv2.EC2
	outposts     *mocksv2.Outposts
}

// NewMockProvider returns a new MockProvider
func NewMockProvider() *MockProvider {
	return &MockProvider{

		asg:            &mocksv2.ASG{},
		eks:            &mocksv2.EKS{},
		cloudtrail:     &mocksv2.CloudTrail{},
		cloudwatchlogs: &mocksv2.CloudWatchLogs{},

		sts:                 &mocksv2.STS{},
		stsPresigner:        &fakes.FakeSTSPresigner{},
		cfn:                 &mocksv2.CloudFormation{},
		elb:                 &mocksv2.ELB{},
		elbV2:               &mocksv2.ELBV2{},
		ssm:                 &mocksv2.SSM{},
		iam:                 &mocksv2.IAM{},
		ec2:                 &mocksv2.EC2{},
		outposts:            &mocksv2.Outposts{},
		credentialsProvider: &mocksv2.CredentialsProvider{},
	}
}

// STS returns a representation of the STS v2 API
func (m MockProvider) STS() awsapi.STS {
	return m.sts
}

func (m MockProvider) STSPresigner() api.STSPresigner {
	return m.stsPresigner
}

// MockSTS returns a mocked STS v2 API
func (m MockProvider) MockSTS() *mocksv2.STS {
	return m.sts
}

// MockSTSPresigner returns a mocked STS v2 API
func (m MockProvider) MockSTSPresigner() *fakes.FakeSTSPresigner {
	return m.stsPresigner.(*fakes.FakeSTSPresigner)
}

// CloudFormation returns a representation of the CloudFormation v2 API
func (m MockProvider) CloudFormation() awsapi.CloudFormation {
	return m.cfn
}

// MockCloudFormation returns a mocked CloudFormation v2 API
func (m MockProvider) MockCloudFormation() *mocksv2.CloudFormation {
	return m.cfn
}

func (m *MockProvider) ELB() awsapi.ELB {
	return m.elb
}

func (m *MockProvider) MockELB() *mocksv2.ELB {
	return m.elb
}

func (m *MockProvider) ELBV2() awsapi.ELBV2 {
	return m.elbV2
}

func (m *MockProvider) MockELBV2() *mocksv2.ELBV2 {
	return m.elbV2
}

// CloudFormationRoleARN returns, if any, a service role used by CloudFormation to call AWS API on your behalf
func (m MockProvider) CloudFormationRoleARN() string { return m.cfnRoleARN }

// CloudFormationDisableRollback returns whether stacks should not rollback on failure
func (m MockProvider) CloudFormationDisableRollback() bool {
	return false
}

// ASG returns a representation of the ASG API
func (m MockProvider) ASG() awsapi.ASG { return m.asg }

// MockASG returns a mocked ASG API
func (m MockProvider) MockASG() *mocksv2.ASG { return m.ASG().(*mocksv2.ASG) }

// EKS returns a representation of the EKS API
func (m MockProvider) EKS() awsapi.EKS { return m.eks }

// MockEKS returns a mocked EKS API
func (m MockProvider) MockEKS() *mocksv2.EKS { return m.eks }

// EC2 returns a representation of the EC2 API
func (m MockProvider) EC2() awsapi.EC2 { return m.ec2 }

// MockEC2 returns a mocked EC2 API
func (m MockProvider) MockEC2() *mocksv2.EC2 { return m.ec2 }

// SSM returns a representation of the SSM API
func (m MockProvider) SSM() awsapi.SSM { return m.ssm }

// MockSSM returns a mocked SSM API
func (m MockProvider) MockSSM() *mocksv2.SSM { return m.ssm }

// IAM returns a representation of the IAM API
func (m MockProvider) IAM() awsapi.IAM { return m.iam }

// MockIAM returns a mocked IAM API
func (m MockProvider) MockIAM() *mocksv2.IAM { return m.iam }

// CloudTrail returns a representation of the CloudTrail API
func (m MockProvider) CloudTrail() awsapi.CloudTrail { return m.cloudtrail }

// MockCloudTrail returns a mocked CloudTrail API
func (m MockProvider) MockCloudTrail() *mocksv2.CloudTrail {
	return m.CloudTrail().(*mocksv2.CloudTrail)
}

// CloudWatchLogs returns a representation of the CloudWatchLogs API
func (m MockProvider) CloudWatchLogs() awsapi.CloudWatchLogs { return m.cloudwatchlogs }

// MockCloudWatchLogs returns a mocked CloudWatchLogs API
func (m MockProvider) MockCloudWatchLogs() *mocksv2.CloudWatchLogs {
	return m.CloudWatchLogs().(*mocksv2.CloudWatchLogs)
}

func (m MockProvider) CredentialsProvider() awsv2.CredentialsProvider {
	return m.credentialsProvider
}

func (m MockProvider) MockCredentialsProvider() *mocksv2.CredentialsProvider {
	return m.CredentialsProvider().(*mocksv2.CredentialsProvider)
}

// Outposts returns a representation of the Outposts API
func (m MockProvider) Outposts() awsapi.Outposts { return m.outposts }

// MockOutposts returns a mocked Outposts API
func (m MockProvider) MockOutposts() *mocksv2.Outposts {
	return m.outposts
}

// Profile returns current profile setting
func (m MockProvider) Profile() api.Profile { return ProviderConfig.Profile }

// Region returns current region setting
func (m MockProvider) Region() string {
	if m.region != "" {
		return m.region
	}
	return ProviderConfig.Region
}

// SetRegion can be used to set the region of the provider
func (m *MockProvider) SetRegion(r string) {
	m.region = r
}

// WaitTimeout returns current timeout setting
func (m MockProvider) WaitTimeout() time.Duration {
	if m.waitTimeout != nil {
		return *m.waitTimeout
	}
	return ProviderConfig.WaitTimeout
}

func (m *MockProvider) SetWaitTimeout(t time.Duration) {
	m.waitTimeout = &t
}

func (m MockProvider) AWSConfig() awsv2.Config {
	options := []func(options *config.LoadOptions) error{
		config.WithRegion(ProviderConfig.Region),
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), options...)
	if err != nil {
		panic(err)
	}
	return cfg
}
