package update

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("update addon", func() {
	Describe("configuration values has invalid format", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   "cluster-1",
				Region: "us-west-2",
			},
			Addons: []*api.Addon{
				{
					Name:                "core-dns",
					ConfigurationValues: "@not a json not a yaml",
				},
			},
		}
		It("should return an error", func() {
			cmd := newMockCmd("addon", "--config-file", ctltest.CreateConfigFile(cfg))
			_, err := cmd.execute()
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("\"%s\" is not valid, supported format(s) are: JSON and YAML", cfg.Addons[0].ConfigurationValues))))
		})
	})
})
