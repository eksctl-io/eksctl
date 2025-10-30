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

var _ = Describe("ShouldWaitForAddons", func() {
	It("returns false when nodes is false", func() {
		addons := []*api.Addon{{Name: api.VPCCNIAddon}}
		Expect(addon.ShouldWaitForAddons(false, addons)).To(BeFalse())
	})

	It("returns false when no addons", func() {
		Expect(addon.ShouldWaitForAddons(true, nil)).To(BeFalse())
	})

	It("returns false for metrics-server", func() {
		addons := []*api.Addon{{Name: api.MetricsServerAddon}}
		Expect(addon.ShouldWaitForAddons(true, addons)).To(BeFalse())
	})

	It("returns true for unknown addons", func() {
		addons := []*api.Addon{{Name: "unknown-addon"}}
		Expect(addon.ShouldWaitForAddons(true, addons)).To(BeTrue())
	})

	It("returns true when any addon requires wait", func() {
		addons := []*api.Addon{
			{Name: api.MetricsServerAddon},
			{Name: "unknown-addon"},
		}
		Expect(addon.ShouldWaitForAddons(true, addons)).To(BeTrue())
	})
})
