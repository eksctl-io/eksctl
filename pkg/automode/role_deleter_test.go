package automode_test

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/smithy-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/automode"
	"github.com/weaveworks/eksctl/pkg/automode/mocks"
)

type roleDeleterTest struct {
	updateMock func(*mocks.StackDeleter)
	cluster    *ekstypes.Cluster

	expectedErr string
}

var _ = DescribeTable("Role Creator", func(t roleDeleterTest) {
	var stackDeleter mocks.StackDeleter
	roleDeleter := &automode.RoleDeleter{
		Cluster:      t.cluster,
		StackDeleter: &stackDeleter,
	}
	if t.updateMock != nil {
		t.updateMock(&stackDeleter)
	}
	err := roleDeleter.DeleteIfRequired(context.Background())
	if t.expectedErr != "" {
		Expect(err).To(MatchError(t.expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
	stackDeleter.AssertExpectations(GinkgoT())
},
	Entry("Auto Mode is disabled", roleDeleterTest{
		cluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
	}),
	Entry("role stack does not exist", roleDeleterTest{
		cluster: &ekstypes.Cluster{
			Name: aws.String("cluster"),
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(true),
			},
		},
		updateMock: func(d *mocks.StackDeleter) {
			d.EXPECT().DescribeStack(mock.Anything, &cfntypes.Stack{StackName: aws.String("eksctl-cluster-auto-mode-role")}).
				Return(nil, &smithy.OperationError{
					Err: errors.New("ValidationError"),
				}).Once()
		},
	}),
	Entry("role stack exists", roleDeleterTest{
		cluster: &ekstypes.Cluster{
			Name: aws.String("cluster"),
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(true),
			},
		},
		updateMock: func(d *mocks.StackDeleter) {
			var stack cfntypes.Stack
			d.EXPECT().DescribeStack(mock.Anything, &cfntypes.Stack{StackName: aws.String("eksctl-cluster-auto-mode-role")}).
				Return(&stack, nil).Once()
			d.EXPECT().DeleteStackSync(mock.Anything, &stack).Return(nil).Once()
		},
	}),
)
