package addon_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
)

var _ = Describe("Addon", func() {
	When("the version is supported", func() {
		It("does not error", func() {
			_, err := addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{Version: "1.18"}}, &mocks.EKSAPI{}, nil, false, nil, nil, 5*time.Minute)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("the version is not supported", func() {
		It("errors", func() {
			_, err := addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{Version: "1.17"}}, &mocks.EKSAPI{}, nil, false, nil, nil, 5*time.Minute)
			Expect(err).To(MatchError("addons not supported on 1.17. Must be using 1.18 or newer"))
		})
	})
})
