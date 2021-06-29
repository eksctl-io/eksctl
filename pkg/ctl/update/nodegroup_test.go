package update

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("update nodegroup", func() {
	It("returns error if config file is not set", func() {
		cmd := newMockCmd("nodegroup")
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("--config-file must be set")))
	})

	It("returns error if nodegroup is not set in config", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   "cluster-1",
				Region: "us-west-2",
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("managedNodeGroups field must be set")))
	})
})
