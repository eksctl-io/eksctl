package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("update control plane component config", func() {

	type updateControlPlaneComponentConfigEntry struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e updateControlPlaneComponentConfigEntry) {
		cmd := newMockCmd(append([]string{"update-control-plane-component-config"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing --config-file", updateControlPlaneComponentConfigEntry{
			expectedErr: "--config-file/-f is required to update control plane component config",
		}),
		Entry("missing --config-file when --cluster is set", updateControlPlaneComponentConfigEntry{
			args:        []string{"--cluster", "test"},
			expectedErr: "--config-file/-f is required to update control plane component config",
		}),
		Entry("setting --cluster and --config-file at the same time", updateControlPlaneComponentConfigEntry{
			args:        []string{"--cluster", "test", "--config-file", "../../../examples/01-simple-cluster.yaml"},
			expectedErr: "cannot use --cluster when --config-file/-f is set",
		}),
		Entry("no component config set in the config file", updateControlPlaneComponentConfigEntry{
			args:        []string{"--config-file", "../../../examples/01-simple-cluster.yaml"},
			expectedErr: "at least one of kubeAPIServerConfig, kubeSchedulerConfig or kubeControllerManagerConfig is required",
		}),
	)

	// The loader only requires that one of the three components is present, so a component
	// holding no values at all gets past it. The command itself has to reject those, rather
	// than sending a request whose components are all empty.
	DescribeTable("component config with no values set", func(cfg *api.ClusterConfig) {
		cfg.TypeMeta = api.ClusterConfigTypeMeta()
		cfg.Metadata = &api.ClusterMeta{
			Name:   "cluster-1",
			Region: "us-west-2",
		}
		cmd := newMockCmd("update-control-plane-component-config", "--config-file", ctltest.CreateConfigFile(cfg))
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring("no control plane component config values are set; nothing to update")))
	},
		Entry("empty kubeAPIServerConfig", &api.ClusterConfig{
			KubeAPIServerConfig: &api.KubeAPIServerConfig{},
		}),
		Entry("empty serviceNodePortRange", &api.ClusterConfig{
			KubeAPIServerConfig: &api.KubeAPIServerConfig{
				ServiceNodePortRange: &api.ServiceNodePortRange{},
			},
		}),
		Entry("empty kubeSchedulerConfig", &api.ClusterConfig{
			KubeSchedulerConfig: &api.KubeSchedulerConfig{},
		}),
		Entry("empty scoringStrategy", &api.ClusterConfig{
			KubeSchedulerConfig: &api.KubeSchedulerConfig{
				NodeResourcesFit: &api.NodeResourcesFitConfig{
					ScoringStrategy: &api.ScoringStrategy{},
				},
			},
		}),
		Entry("empty kubeControllerManagerConfig", &api.ClusterConfig{
			KubeControllerManagerConfig: &api.KubeControllerManagerConfig{},
		}),
		Entry("all three components present but empty", &api.ClusterConfig{
			KubeAPIServerConfig:         &api.KubeAPIServerConfig{},
			KubeSchedulerConfig:         &api.KubeSchedulerConfig{},
			KubeControllerManagerConfig: &api.KubeControllerManagerConfig{},
		}),
	)
})
