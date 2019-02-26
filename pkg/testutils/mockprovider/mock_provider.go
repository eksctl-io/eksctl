package mockprovider

import (
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
)

// MockProvider stores the mocked APIs
type MockProvider struct {
	cfnRoleARN string
	cfn        *mocks.CloudFormationAPI
	eks        *mocks.EKSAPI
	ec2        *mocks.EC2API
	sts        *mocks.STSAPI
	iam        *mocks.IAMAPI
	cloudtrail *mocks.CloudTrailAPI
}

// NewMockProvider returns a new MockProvider
func NewMockProvider() *MockProvider {
	return &MockProvider{
		cfn:        &mocks.CloudFormationAPI{},
		eks:        &mocks.EKSAPI{},
		ec2:        &mocks.EC2API{},
		sts:        &mocks.STSAPI{},
		iam:        &mocks.IAMAPI{},
		cloudtrail: &mocks.CloudTrailAPI{},
	}
}

// ProviderConfig holds current global config
var ProviderConfig = &api.ProviderConfig{
	Region:      api.DefaultRegion,
	Profile:     "default",
	WaitTimeout: 1200000000000,
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

// MockEC2 returns a mocked EC2 API
func (m MockProvider) MockEC2() *mocks.EC2API { return m.EC2().(*mocks.EC2API) }

// STS returns a representation of the STS API
func (m MockProvider) STS() stsiface.STSAPI { return m.sts }

// MockSTS returns a mocked STS API
func (m MockProvider) MockSTS() *mocks.STSAPI { return m.STS().(*mocks.STSAPI) }

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
