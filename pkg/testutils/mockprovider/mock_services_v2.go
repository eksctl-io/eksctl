package mockprovider

import (
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

type MockServicesV2 struct {
	sts            *mocksv2.STS
	cloudformation *mocksv2.CloudFormation
}

func (s *MockServicesV2) STSV2() awsapi.STS {
	return s.sts
}

func (s *MockServicesV2) MockSTSV2() *mocksv2.STS {
	return s.sts
}

func (s *MockServicesV2) CloudFormationV2() awsapi.CloudFormation {
	return s.cloudformation
}

func (s *MockServicesV2) MockCloudFormationV2() *mocksv2.CloudFormation {
	return s.cloudformation
}
