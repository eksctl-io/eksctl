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
})
