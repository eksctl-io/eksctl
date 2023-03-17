package addon_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

var _ = Describe("Addon", func() {
	When("the version is supported", func() {
		It("does not error", func() {
			_, err := addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{Version: "1.18"}}, &mocksv2.EKS{}, nil, false, nil, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
