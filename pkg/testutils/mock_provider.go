package testutils

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
)

type MockProvider struct {
	cfn *mocks.CloudFormationAPI
	eks *mocks.EKSAPI
	ec2 *mocks.EC2API
	sts *mocks.STSAPI
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		cfn: &mocks.CloudFormationAPI{},
		eks: &mocks.EKSAPI{},
		ec2: &mocks.EC2API{},
		sts: &mocks.STSAPI{},
	}
}

func (m MockProvider) CloudFormation() cloudformationiface.CloudFormationAPI { return m.cfn }
func (m MockProvider) MockCloudFormation() *mocks.CloudFormationAPI {
	return m.CloudFormation().(*mocks.CloudFormationAPI)
}

func (m MockProvider) EKS() eksiface.EKSAPI   { return m.eks }
func (m MockProvider) MockEKS() *mocks.EKSAPI { return m.EKS().(*mocks.EKSAPI) }
func (m MockProvider) EC2() ec2iface.EC2API   { return m.ec2 }
func (m MockProvider) MockEC2() *mocks.EC2API { return m.EC2().(*mocks.EC2API) }
func (m MockProvider) STS() stsiface.STSAPI   { return m.sts }
func (m MockProvider) MockSTS() *mocks.STSAPI { return m.STS().(*mocks.STSAPI) }
