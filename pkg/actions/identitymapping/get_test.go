package identitymapping_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/identitymapping"
	"github.com/weaveworks/eksctl/pkg/authconfigmap/fakes"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("Get", func() {
	It("returns the identity mappings matching the arn", func() {
		rawClient := testutils.NewFakeRawClient()
		fakeACM := new(fakes.FakeManager)
		manager := identitymapping.New(rawClient, fakeACM)
		fakeACM.IdentitiesReturns([]iam.Identity{
			iam.UserIdentity{UserARN: "foo"},
			iam.UserIdentity{UserARN: "bar"},
		}, nil)
		identities, err := manager.Get("foo")

		Expect(err).NotTo(HaveOccurred())
		Expect(fakeACM.IdentitiesCallCount()).To(Equal(1))
		Expect(identities).To(HaveLen(1))
		Expect(identities[0].ARN()).To(Equal("foo"))
	})
})
