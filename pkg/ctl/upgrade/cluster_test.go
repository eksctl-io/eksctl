package upgrade

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("upgrade cluster", func() {

	newMockUpgradeClusterCmd := func(args ...string) *ctltest.MockCmd {
		return ctltest.NewMockCmd(upgradeClusterWithRunFunc, "upgrade", args...)
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

		It("should accept the --version flag", func() {
			cmd := newMockUpgradeClusterCmd("cluster", "--name", "clus-1", "--region", "eu-north-1", "--version", "1.16")
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

			// The version should be empty when not specified as a flag
			Expect(cfg.Metadata.Version).To(Equal(""))

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
			configFile = ctltest.CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("accepts --approve flag with the config file", func() {
			configFile = ctltest.CreateConfigFile(cfg)
			cmd := newMockUpgradeClusterCmd("cluster", "--config-file", configFile, "--approve")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails without a cluster name", func() {
			cfg.Metadata.Name = ""
			configFile = ctltest.CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.name must be set"))
		})

		It("fails without a region", func() {
			cfg.Metadata.Region = ""
			configFile = ctltest.CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.region must be set"))
		})

		It("fails if the config file and the --name are specified", func() {
			configFile = ctltest.CreateConfigFile(cfg)
			cmd := newMockUpgradeClusterCmd("cluster", "--name", "clus-1", "--config-file", configFile)
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("cannot use --name when --config-file/-f is set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("loads the config file correctly", func() {
			cfg.Metadata.Version = "1.16"
			configFile = ctltest.CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "--config-file", configFile)
			_, err := cmd.Execute()
			Expect(err).To(Not(HaveOccurred()))

			loadedCfg := cmd.Cmd.ClusterConfig.Metadata
			Expect(loadedCfg.Name).To(Equal("cluster-1"))
			Expect(loadedCfg.Region).To(Equal("us-west-2"))
			Expect(loadedCfg.Version).To(Equal("1.16"))
		})

		It("when not specified in the config file the version is empty", func() {
			cfg.Metadata.Version = ""
			configFile = ctltest.CreateConfigFile(cfg)

			cmd := newMockUpgradeClusterCmd("cluster", "--config-file", configFile)
			_, err := cmd.Execute()
			Expect(err).To(Not(HaveOccurred()))

			loadedCfg := cmd.Cmd.ClusterConfig.Metadata
			Expect(loadedCfg.Version).To(Equal(""))
		})
	})
})
