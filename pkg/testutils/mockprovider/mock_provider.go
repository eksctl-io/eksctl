package mockprovider

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
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

	cfnRoleARN string
	cfn        *mocks.CloudFormationAPI
	eks        *mocks.EKSAPI
	ec2        *mocks.EC2API
	elb        *mocks.ELBAPI
	elbv2      *mocks.ELBV2API
	sts        *mocks.STSAPI
	ssm        *mocks.SSMAPI
	iam        *mocks.IAMAPI
	cloudtrail *mocks.CloudTrailAPI
}

// NewMockProvider returns a new MockProvider
func NewMockProvider() *MockProvider {
	return &MockProvider{
		Client: NewMockAWSClient(),

		cfn:        &mocks.CloudFormationAPI{},
		eks:        &mocks.EKSAPI{},
		ec2:        &mocks.EC2API{},
		elb:        &mocks.ELBAPI{},
		elbv2:      &mocks.ELBV2API{},
		sts:        &mocks.STSAPI{},
		ssm:        &mocks.SSMAPI{},
		iam:        &mocks.IAMAPI{},
		cloudtrail: &mocks.CloudTrailAPI{},
	}
}

// CloudFormation returns a representation of the CloudFormation API
func (m MockProvider) CloudFormation() cloudformationiface.CloudFormationAPI { return m.cfn }

// CloudFormationRoleARN returns, if any,  a service role used by CloudFormation to call AWS API on your behalf
func (m MockProvider) CloudFormationRoleARN() string { return m.cfnRoleARN }

// MockCloudFormation returns a mocked CloudFormation API
func (m MockProvider) MockCloudFormation() *mocks.CloudFormationAPI {
	return m.CloudFormation().(*mocks.CloudFormationAPI)
}

// EKS returns a representation of the EKS API
func (m MockProvider) EKS() eksiface.EKSAPI { return m.eks }

// MockEKS returns a mocked EKS API
func (m MockProvider) MockEKS() *mocks.EKSAPI { return m.EKS().(*mocks.EKSAPI) }

// EC2 returns a representation of the EC2 API
func (m MockProvider) EC2() ec2iface.EC2API { return m.ec2 }

// ELB returns a representation of the ELB API
func (m MockProvider) ELB() elbiface.ELBAPI { return m.elb }

// ELBV2 returns a representation of the ELBV2 API
func (m MockProvider) ELBV2() elbv2iface.ELBV2API { return m.elbv2 }

// MockEC2 returns a mocked EC2 API
func (m MockProvider) MockEC2() *mocks.EC2API { return m.EC2().(*mocks.EC2API) }

// STS returns a representation of the STS API
func (m MockProvider) STS() stsiface.STSAPI { return m.sts }

// MockSTS returns a mocked STS API
func (m MockProvider) MockSTS() *mocks.STSAPI { return m.STS().(*mocks.STSAPI) }

// SSM returns a representation of the SSM API
func (m MockProvider) SSM() ssmiface.SSMAPI { return m.ssm }

// MockSSM returns a mocked SSM API
func (m MockProvider) MockSSM() *mocks.SSMAPI { return m.SSM().(*mocks.SSMAPI) }

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

// Profile returns current profile setting
func (m MockProvider) Profile() string { return ProviderConfig.Profile }

// Region returns current region setting
func (m MockProvider) Region() string { return ProviderConfig.Region }

// WaitTimeout returns current timeout setting
func (m MockProvider) WaitTimeout() time.Duration { return ProviderConfig.WaitTimeout }

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

func BuildNewMockRequestForMockOutput(m *MockAWSClient, in *MockInput) func([]request.Option) (*request.Request, error) {
	return func(opts []request.Option) (*request.Request, error) {
		req, _ := m.MockRequestForMockOutput(in)
		req.ApplyOptions(opts...)
		return req, nil
	}
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
