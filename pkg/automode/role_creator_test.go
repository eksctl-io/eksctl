package automode_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/automode"
	"github.com/weaveworks/eksctl/pkg/automode/mocks"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

type roleCreatorTest struct {
	updateMock func(*mocks.StackCreator)

	expectedNodeRoleARN string
	expectedErr         string
}

var _ = DescribeTable("Role Creator", func(t roleCreatorTest) {
	var stackCreator mocks.StackCreator
	roleCreator := &automode.RoleCreator{
		StackCreator: &stackCreator,
	}
	if t.updateMock != nil {
		t.updateMock(&stackCreator)
	}
	nodeRoleARN, err := roleCreator.CreateOrImport(context.Background(), "cluster")
	if t.expectedErr != "" {
		Expect(err).To(MatchError(t.expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
	Expect(nodeRoleARN).To(Equal(t.expectedNodeRoleARN))
	stackCreator.AssertExpectations(GinkgoT())
},
	Entry("Auto Mode role exists in cluster stack", roleCreatorTest{
		updateMock: func(s *mocks.StackCreator) {
			s.EXPECT().GetClusterStackIfExists(mock.Anything).Return(&cfntypes.Stack{
				Outputs: []cfntypes.Output{
					{
						OutputKey:   aws.String("AutoModeNodeRoleARN"),
						OutputValue: aws.String("arn:aws:iam::000:role/AutoModeNodeRole"),
					},
				},
			}, nil).Once()
		},
		expectedNodeRoleARN: "arn:aws:iam::000:role/AutoModeNodeRole",
	}),
	Entry("Auto Mode role is missing in cluster stack", roleCreatorTest{
		updateMock: func(s *mocks.StackCreator) {
			s.EXPECT().GetClusterStackIfExists(mock.Anything).Return(&cfntypes.Stack{}, nil).Once()
			s.EXPECT().CreateStack(mock.Anything, "eksctl-cluster-auto-mode-role", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, _ string, resourceSet builder.ResourceSetReader, _, _ map[string]string, errCh chan error) error {
					if err := resourceSet.GetAllOutputs(cfntypes.Stack{
						Outputs: []cfntypes.Output{
							{
								OutputKey:   aws.String("AutoModeNodeRoleARN"),
								OutputValue: aws.String("arn:aws:iam::000:role/AutoModeNodeRole"),
							},
						},
					}); err != nil {
						return err
					}
					close(errCh)
					return nil
				}).Once()
		},
		expectedNodeRoleARN: "arn:aws:iam::000:role/AutoModeNodeRole",
	}),
)
