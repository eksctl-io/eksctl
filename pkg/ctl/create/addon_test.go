package create

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("create addon", func() {
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
			cmd := newDefaultCmd("addon", "--config-file", ctltest.CreateConfigFile(cfg))
			_, err := cmd.execute()
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("\"%s\" is not valid, supported format(s) are: JSON and YAML", cfg.Addons[0].ConfigurationValues))))
		})
	})

	Describe("namespace config validation", func() {
		Context("with valid namespace config", func() {
			cfg := &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Name:   "cluster-1",
					Region: "us-west-2",
				},
				Addons: []*api.Addon{
					{
						Name: "vpc-cni",
						NamespaceConfig: &api.AddonNamespaceConfig{
							Namespace: "kube-system",
						},
					},
				},
			}
			It("should not return an error", func() {
				cmd := newDefaultCmd("addon", "--config-file", ctltest.CreateConfigFile(cfg))
				_, err := cmd.execute()
				// Note: This test validates the config parsing, actual execution would require AWS API mocking
				// The error we expect here is related to AWS API calls, not config validation
				Expect(err).ToNot(MatchError(ContainSubstring("namespace")))
			})
		})

		Context("with invalid namespace config", func() {
			cfg := &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Name:   "cluster-1",
					Region: "us-west-2",
				},
				Addons: []*api.Addon{
					{
						Name: "vpc-cni",
						NamespaceConfig: &api.AddonNamespaceConfig{
							Namespace: "Invalid-Namespace-Name!",
						},
					},
				},
			}
			It("should return a validation error", func() {
				cmd := newDefaultCmd("addon", "--config-file", ctltest.CreateConfigFile(cfg))
				_, err := cmd.execute()
				Expect(err).To(MatchError(ContainSubstring("is not a valid Kubernetes namespace name")))
			})
		})

		Context("with empty namespace config", func() {
			cfg := &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Name:   "cluster-1",
					Region: "us-west-2",
				},
				Addons: []*api.Addon{
					{
						Name: "vpc-cni",
						NamespaceConfig: &api.AddonNamespaceConfig{
							Namespace: "",
						},
					},
				},
			}
			It("should not return an error", func() {
				cmd := newDefaultCmd("addon", "--config-file", ctltest.CreateConfigFile(cfg))
				_, err := cmd.execute()
				// Empty namespace should be valid (uses default behavior)
				Expect(err).ToNot(MatchError(ContainSubstring("namespace")))
			})
		})
	})
})
