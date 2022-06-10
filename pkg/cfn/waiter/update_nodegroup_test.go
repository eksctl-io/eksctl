package waiter

import (
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
		status, err := WaitForNodegroupUpdate("update-1", eksAPI, 35*time.Second, func(attempts int) time.Duration {
			return 1 * time.Second
		})
		Expect(status).To(Equal(string(ekstypes.UpdateStatusSuccessful)))
		Expect(err).NotTo(HaveOccurred())
	})
})
