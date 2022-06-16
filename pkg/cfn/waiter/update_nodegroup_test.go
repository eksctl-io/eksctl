package waiter

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

var _ = Describe("WaitForNodegroupUpdate", func() {
	It("can wait for a successful outcome", func() {
		eksAPI := &mocksv2.EKS{}
		eksAPI.On("DescribeUpdate", mock.Anything, &eks.DescribeUpdateInput{
			UpdateId: aws.String("update-1"),
		}).Return(&eks.DescribeUpdateOutput{
			Update: &ekstypes.Update{
				Status: ekstypes.UpdateStatusSuccessful,
			},
		}, nil)
		status, err := WaitForNodegroupUpdate(context.TODO(), "update-1", eksAPI, 35*time.Second, func(attempts int) time.Duration {
			return 1 * time.Nanosecond
		})
		Expect(status).To(Equal(string(ekstypes.UpdateStatusSuccessful)))
		Expect(err).NotTo(HaveOccurred())
	})
	When("describe fails", func() {
		It("errors", func() {
			eksAPI := &mocksv2.EKS{}
			eksAPI.On("DescribeUpdate", mock.Anything, &eks.DescribeUpdateInput{
				UpdateId: aws.String("update-1"),
			}).Return(&eks.DescribeUpdateOutput{
				Update: &ekstypes.Update{
					Status: ekstypes.UpdateStatusFailed,
				},
			}, errors.New("nope"))
			status, err := WaitForNodegroupUpdate(context.TODO(), "update-1", eksAPI, 35*time.Second, func(attempts int) time.Duration {
				return 1 * time.Nanosecond
			})
			Expect(status).To(BeEmpty())
			Expect(err).To(MatchError(ContainSubstring("failed to describe update for update id update-1: nope")))
		})
	})
	When("status is failed", func() {
		It("errors", func() {
			eksAPI := &mocksv2.EKS{}
			eksAPI.On("DescribeUpdate", mock.Anything, &eks.DescribeUpdateInput{
				UpdateId: aws.String("update-1"),
			}).Return(&eks.DescribeUpdateOutput{
				Update: &ekstypes.Update{
					Status: ekstypes.UpdateStatusFailed,
				},
			}, nil)
			status, err := WaitForNodegroupUpdate(context.TODO(), "update-1", eksAPI, 35*time.Second, func(attempts int) time.Duration {
				return 1 * time.Nanosecond
			})
			Expect(status).To(Equal(string(ekstypes.UpdateStatusFailed)))
			Expect(err).To(MatchError(ContainSubstring("update failed or was cancelled")))
		})
	})
})
