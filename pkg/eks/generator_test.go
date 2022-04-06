package eks_test

import (
	"context"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/credentials/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("TokenGenerator", func() {
	var (
		provider *mockprovider.MockProvider
		clock    *fakes.FakeClock
	)

	BeforeEach(func() {
		provider = mockprovider.NewMockProvider()
		clock = &fakes.FakeClock{}
	})
	Context("GetWithSTS", func() {
		It("can generate a token", func() {
			fakeGenerator := provider.MockSTSPresigner()
			fakeGenerator.PresignGetCallerIdentityReturns(&v4.PresignedHTTPRequest{
				URL: "https://example.com",
			}, nil)
			clock.NowReturns(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC))
			generator := eks.NewGenerator(provider.MockSTSPresigner(), clock)
			token, err := generator.GetWithSTS(context.TODO(), "cluster-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(token.Token).To(Equal("k8s-aws-v1.aHR0cHM6Ly9leGFtcGxlLmNvbQ"))
		})
		When("PresignGetCaller returns an error", func() {
			It("errors", func() {
				fakeGenerator := provider.MockSTSPresigner()
				fakeGenerator.PresignGetCallerIdentityReturns(nil, errors.New("nope"))
				clock.NowReturns(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC))
				generator := eks.NewGenerator(provider.MockSTSPresigner(), clock)
				_, err := generator.GetWithSTS(context.TODO(), "cluster-id")
				Expect(err).To(MatchError("failed to presign caller identity: nope"))
			})
		})
	})
})
