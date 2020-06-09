package upgrade

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("enable repo", func() {

	newMockUpgradeClusterCmd := func(args ...string) *MockCmd {
		return NewMockCmd(upgradeClusterWithRunFunc, "upgrade", args...)
	}

	Describe("without a config file", func() {

		It("should accept a name argument", func() {
			cmd := newMockUpgradeClusterCmd("cluster", "clus-1")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})
		It("should accept a --name flag", func() {
			cmd := newMockUpgradeClusterCmd("cluster", "--name", "clus-1")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept the --region flag", func() {
			cmd := newMockUpgradeClusterCmd("cluster", "--name", "clus-1", "--region", "eu-north-1")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("accepts --approve flag", func() {
			cmd := newMockUpgradeClusterCmd("cluster", "--name", "clus-1", "--approve")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("loads all flags correctly", func() {
			cmd := newMockUpgradeClusterCmd("cluster",
				"--name", "clus-1",
				"--region", "us-west-2",
				"--timeout", "123m",
				"--approve",
			)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig
			Expect(cfg.Metadata.Name).To(Equal("clus-1"))
			// I cannot test the region here because this flag is loaded into the cmd.ProviderConfig.Region
			//Expect(cfg.Metadata.Region).To(Equal("us-west-2"))
			Expect(cmd.Cmd.ProviderConfig.Region).To(Equal("us-west-2"))
			Expect(cmd.Cmd.Plan).To(BeFalse())
			Expect(cmd.Cmd.ProviderConfig.WaitTimeout).To(Equal(123 * time.Minute))
		})
	})

	Describe("with a config file", func() {
		var configFile string
		var cfg *api.ClusterConfig

		BeforeEach(func() {
			// Minimal valid cluster config for the command to work
			cfg = &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Name:   "cluster-1",
					Region: "us-west-2",
				},
			}
		})

		AfterEach(func() {
			if configFile != "" {
				os.Remove(configFile)
			}
		})

		It("succeeds with the basic configuration", func() {
			configFile = CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("accepts --approve flag with the config file", func() {
			configFile = CreateConfigFile(cfg)
			cmd := newMockUpgradeClusterCmd("cluster", "--config-file", configFile, "--approve")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails without a cluster name", func() {
			cfg.Metadata.Name = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.name must be set"))
		})

		It("fails without a region", func() {
			cfg.Metadata.Region = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.region must be set"))
		})

		It("fails if the config file and the --name are specified", func() {
			configFile = CreateConfigFile(cfg)
			cmd := newMockUpgradeClusterCmd("cluster", "--name", "clus-1", "--config-file", configFile)
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("cannot use --name when --config-file/-f is set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

	})
})
