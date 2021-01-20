package identitymapping_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/identitymapping"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("Delete", func() {
	It("delete the identity mapping", func() {
		rawClient := testutils.NewFakeRawClient()
		fakeACM := new(fakes.FakeManager)
		manager := identitymapping.New(rawClient, fakeACM)
		err := manager.Delete([]*api.IAMIdentityMapping{
			{
				ARN:      "arn123",
				Username: "user",
				Groups:   []string{"system:masters"},
			},
		}, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeACM.RemoveIdentityCallCount()).To(Equal(1))
		arn, all := fakeACM.RemoveIdentityArgsForCall(0)
		Expect(arn).To(Equal("arn123"))
		Expect(all).To(BeFalse())
		Expect(fakeACM.SaveCallCount()).To(Equal(1))
	})
})
