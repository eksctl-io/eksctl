package cmdutils

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Auto Mode loader", func() {
	type loaderTest struct {
		config      string
		expectedErr string
	}

	DescribeTable("validates kubernetesNetworkConfig.elasticLoadBalancing", func(t loaderTest) {
		configFile := filepath.Join(GinkgoT().TempDir(), "config.yaml")
		Expect(os.WriteFile(configFile, []byte(`apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: cluster
  region: us-west-2
`+t.config), 0o600)).To(Succeed())

		cobraCmd := &cobra.Command{Use: "test", Run: func(_ *cobra.Command, _ []string) {}}
		cobraCmd.Flags().Int("drain-parallel", 1, "")
		cobraCmd.Flags().Bool("drain-all-nodegroups", false, "")
		cmd := &Cmd{
			CobraCommand:      cobraCmd,
			ClusterConfigFile: configFile,
			ClusterConfig:     api.NewClusterConfig(),
			ProviderConfig:    api.ProviderConfig{},
		}
		err := NewAutoModeLoader(cmd).Load()
		if t.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(t.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
	},
		Entry("enabling Auto Mode without elasticLoadBalancing", loaderTest{
			config: `autoModeConfig:
  enabled: true
`,
		}),
		Entry("enabling Auto Mode with gatewayApi.crds.managed set", loaderTest{
			config: `autoModeConfig:
  enabled: true
kubernetesNetworkConfig:
  elasticLoadBalancing:
    gatewayApi:
      crds:
        managed: false
`,
		}),
		Entry("gatewayApi.crds without managed", loaderTest{
			config: `autoModeConfig:
  enabled: true
kubernetesNetworkConfig:
  elasticLoadBalancing:
    gatewayApi:
      crds: {}
`,
			expectedErr: "kubernetesNetworkConfig.elasticLoadBalancing.gatewayApi.crds.managed must be set to either true or false",
		}),
		Entry("gatewayApi without crds", loaderTest{
			config: `autoModeConfig:
  enabled: true
kubernetesNetworkConfig:
  elasticLoadBalancing:
    gatewayApi: {}
`,
			expectedErr: "kubernetesNetworkConfig.elasticLoadBalancing.gatewayApi.crds.managed must be set to either true or false",
		}),
		Entry("gatewayApi when disabling Auto Mode", loaderTest{
			config: `autoModeConfig:
  enabled: false
kubernetesNetworkConfig:
  elasticLoadBalancing:
    gatewayApi:
      crds:
        managed: true
`,
			expectedErr: "cannot set kubernetesNetworkConfig.elasticLoadBalancing.gatewayApi when the load balancing capability is disabled",
		}),
		Entry("elasticLoadBalancing.enabled does not match Auto Mode", loaderTest{
			config: `autoModeConfig:
  enabled: true
kubernetesNetworkConfig:
  elasticLoadBalancing:
    enabled: false
`,
			expectedErr: "kubernetesNetworkConfig.elasticLoadBalancing.enabled must match autoModeConfig.enabled",
		}),
	)
})
