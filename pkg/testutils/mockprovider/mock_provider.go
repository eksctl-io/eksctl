package mockprovider

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/awsapi"

	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting"

	//"github.com/aws/aws-sdk-go/awstesting/unit"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
)

// ProviderConfig holds current global config
var ProviderConfig = &api.ProviderConfig{
	Region:      api.DefaultRegion,
	Profile:     "default",
	WaitTimeout: 1200000000000,
}

var _ api.ClusterProvider = &MockProvider{}

type MockAWSClient struct {
	*client.Client
}

type MockInput struct{}
type MockOutput struct {
	States []*MockState
}
type MockState struct {
	State *string
}

// MockProvider stores the mocked APIs
type MockProvider struct {
	Client *MockAWSClient

	region         string
	cfnRoleARN     string
	asg            *mocks.AutoScalingAPI
	cfn            *mocks.CloudFormationAPI
	eks            *mocks.EKSAPI
	ec2            *mocks.EC2API
	sts            *mocks.STSAPI
	ssm            *mocksv2.SSM
	iam            *mocks.IAMAPI
	cloudtrail     *mocks.CloudTrailAPI
	cloudwatchlogs *mocks.CloudWatchLogsAPI
	configProvider *mocks.ConfigProvider

	stsV2            *mocksv2.STS
	cloudformationV2 *mocksv2.CloudFormation
	elb              *mocksv2.ELB
	elbV2            *mocksv2.ELBV2
}

// NewMockProvider returns a new MockProvider
func NewMockProvider() *MockProvider {
	return &MockProvider{
		Client: NewMockAWSClient(),

		asg:            &mocks.AutoScalingAPI{},
		cfn:            &mocks.CloudFormationAPI{},
		eks:            &mocks.EKSAPI{},
		ec2:            &mocks.EC2API{},
		sts:            &mocks.STSAPI{},
		ssm:            &mocksv2.SSM{},
		iam:            &mocks.IAMAPI{},
		cloudtrail:     &mocks.CloudTrailAPI{},
		cloudwatchlogs: &mocks.CloudWatchLogsAPI{},
		configProvider: &mocks.ConfigProvider{},

		stsV2:            &mocksv2.STS{},
		cloudformationV2: &mocksv2.CloudFormation{},
		elb:              &mocksv2.ELB{},
		elbV2:            &mocksv2.ELBV2{},
	}
}

// STSV2 returns a representation of the STS v2 API
func (m MockProvider) STSV2() awsapi.STS {
	return m.stsV2
}

// MockSTSV2 returns a mocked STS v2 API
func (m MockProvider) MockSTSV2() *mocksv2.STS {
	return m.stsV2
}

// CloudFormationV2 returns a representation of the CloudFormation v2 API
func (m MockProvider) CloudFormationV2() awsapi.CloudFormation {
	return m.cloudformationV2
}

// MockCloudFormationV2 returns a mocked CloudFormation v2 API
func (m MockProvider) MockCloudFormationV2() *mocksv2.CloudFormation {
	return m.cloudformationV2
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

// CloudFormation returns a representation of the CloudFormation API
func (m MockProvider) CloudFormation() cloudformationiface.CloudFormationAPI { return m.cfn }

// CloudFormationRoleARN returns, if any, a service role used by CloudFormation to call AWS API on your behalf
func (m MockProvider) CloudFormationRoleARN() string { return m.cfnRoleARN }

// CloudFormationDisableRollback returns whether stacks should not rollback on failure
func (m MockProvider) CloudFormationDisableRollback() bool {
	return false
}

// MockCloudFormation returns a mocked CloudFormation API
func (m MockProvider) MockCloudFormation() *mocks.CloudFormationAPI {
	return m.CloudFormation().(*mocks.CloudFormationAPI)
}

// ASG returns a representation of the ASG API
func (m MockProvider) ASG() autoscalingiface.AutoScalingAPI { return m.asg }

// MockASG returns a mocked ASG API
func (m MockProvider) MockASG() *mocks.AutoScalingAPI { return m.ASG().(*mocks.AutoScalingAPI) }

// EKS returns a representation of the EKS API
func (m MockProvider) EKS() eksiface.EKSAPI { return m.eks }

// MockEKS returns a mocked EKS API
func (m MockProvider) MockEKS() *mocks.EKSAPI { return m.EKS().(*mocks.EKSAPI) }

// EC2 returns a representation of the EC2 API
func (m MockProvider) EC2() ec2iface.EC2API { return m.ec2 }

// MockEC2 returns a mocked EC2 API
func (m MockProvider) MockEC2() *mocks.EC2API { return m.EC2().(*mocks.EC2API) }

// STS returns a representation of the STS API
func (m MockProvider) STS() stsiface.STSAPI { return m.sts }

// MockSTS returns a mocked STS API
func (m MockProvider) MockSTS() *mocks.STSAPI { return m.STS().(*mocks.STSAPI) }

// SSM returns a representation of the SSM API
func (m MockProvider) SSM() awsapi.SSM { return m.ssm }

// MockSSM returns a mocked SSM API
func (m MockProvider) MockSSM() *mocksv2.SSM { return m.ssm }

// IAM returns a representation of the IAM API
func (m MockProvider) IAM() iamiface.IAMAPI { return m.iam }

// MockIAM returns a mocked IAM API
func (m MockProvider) MockIAM() *mocks.IAMAPI { return m.IAM().(*mocks.IAMAPI) }

// CloudTrail returns a representation of the CloudTrail API
func (m MockProvider) CloudTrail() cloudtrailiface.CloudTrailAPI { return m.cloudtrail }

// MockCloudTrail returns a mocked CloudTrail API
func (m MockProvider) MockCloudTrail() *mocks.CloudTrailAPI {
	return m.CloudTrail().(*mocks.CloudTrailAPI)
}

// CloudWatchLogs returns a representation of the CloudWatchLogs API
func (m MockProvider) CloudWatchLogs() cloudwatchlogsiface.CloudWatchLogsAPI { return m.cloudwatchlogs }

// MockCloudWatchLogs returns a mocked CloudWatchLogs API
func (m MockProvider) MockCloudWatchLogs() *mocks.CloudWatchLogsAPI {
	return m.CloudWatchLogs().(*mocks.CloudWatchLogsAPI)
}

// Profile returns current profile setting
func (m MockProvider) Profile() string { return ProviderConfig.Profile }

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
func (m MockProvider) WaitTimeout() time.Duration { return ProviderConfig.WaitTimeout }

// ConfigProvider returns a representation of the ConfigProvider
func (m MockProvider) ConfigProvider() client.ConfigProvider {
	return m.configProvider
}

// MockConfigProvider returns a mocked ConfigProvider
func (m MockProvider) MockConfigProvider() client.ConfigProvider {
	return m.configProvider
}

func (m MockProvider) Session() *session.Session {
	panic("not implemented")
}

func NewMockAWSClient() *MockAWSClient {
	m := &MockAWSClient{
		Client: awstesting.NewClient(&aws.Config{
			Region: &ProviderConfig.Region,
		}),
	}

	m.Handlers.Send.Clear()
	m.Handlers.Unmarshal.Clear()
	m.Handlers.UnmarshalMeta.Clear()
	m.Handlers.ValidateResponse.Clear()

	return m
}

func (m *MockAWSClient) MockRequestForMockOutput(input *MockInput) (*request.Request, *MockOutput) {
	op := &request.Operation{
		Name:       "Mock",
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &MockInput{}
	}

	output := &MockOutput{}
	req := m.NewRequest(op, input, output)
	req.Data = output
	return req, output
}

func (m *MockAWSClient) MockRequestForGivenOutput(input, output interface{}) *request.Request {
	op := &request.Operation{
		Name:       "Mock",
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &MockInput{}
	}

	req := m.NewRequest(op, input, output)
	req.Data = output
	return req
}
