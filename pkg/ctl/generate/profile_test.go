package generate

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("generate profile", func() {

	newMockGenerateProfileCmd := func(args ...string) *MockCmd {
		return NewMockCmd(generateProfileWithRunFunc, "generate", args...)
	}

	Describe("without a config file", func() {

		It("should accept a name argument", func() {
			cmd := newMockGenerateProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "app-dev")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.Source).To(Equal("app-dev"))
		})

		It("should accept --profile-source flag", func() {
			cmd := newMockGenerateProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--profile-source", "app-dev")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.Source).To(Equal("app-dev"))
		})

		It("should reject name argument and --profile-source flag", func() {
			cmd := newMockGenerateProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--profile-source", "app-dev", "app-dev")
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--profile-source=app-dev and argument app-dev cannot be used at the same time"))
		})

		It("requires the --profile-source flag", func() {
			cmd := newMockGenerateProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1")
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--profile-source must be set"))
		})

		It("requires the --cluster flag", func() {
			cmd := newMockGenerateProfileCmd("profile", "--region", "eu-north-1", "--profile-source", "app-dev")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})
		It("requires the --region flag", func() {
			cmd := newMockGenerateProfileCmd("profile", "--cluster", "clus-1", "--profile-source", "app-dev")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--region must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("loads all flags correctly", func() {
			cmd := newMockGenerateProfileCmd("profile",
				"--profile-source", "app-dev",
				"--profile-revision", "branch-2",
				"--profile-path", "test-output-dir/dir2",
				"--cluster", "clus-1",
				"--region", "us-west-2",
			)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig
			Expect(cfg.Metadata.Name).To(Equal("clus-1"))
			Expect(cfg.Metadata.Region).To(Equal("us-west-2"))
			Expect(cfg.Git).ToNot(BeNil())

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.Source).To(Equal("app-dev"))
			Expect(cfg.Git.BootstrapProfile.Revision).To(Equal("branch-2"))
			Expect(cfg.Git.BootstrapProfile.OutputPath).To(Equal("test-output-dir/dir2"))
		})

		It("defaults the --output-path to the profile name", func() {
			cmd := newMockGenerateProfileCmd("profile",
				"--profile-source", "app-dev",
				"--profile-revision", "branch-2",
				"--cluster", "clus-1",
				"--region", "us-west-2",
			)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.OutputPath).To(Equal("./app-dev"))
		})

		It("defaults the --output-path to the repo name", func() {
			cmd := newMockGenerateProfileCmd("profile",
				"--profile-source", "git@github.com:weaveworks/eks-quickstart-app-dev.git",
				"--profile-revision", "branch-2",
				"--cluster", "clus-1",
				"--region", "us-west-2",
			)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.OutputPath).To(Equal("./eks-quickstart-app-dev"))
		})

		Describe("with a config file", func() {
			var configFile string
			var cfg *api.ClusterConfig

			BeforeEach(func() {
				// Minimal cluster config for the command to work
				cfg = &api.ClusterConfig{
					TypeMeta: api.ClusterConfigTypeMeta(),
					Metadata: &api.ClusterMeta{
						Name:   "cluster-1",
						Region: "us-west-2",
					},
					Git: &api.Git{
						BootstrapProfile: &api.Profile{Source: "app-dev"},
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

				cmd := newMockGenerateProfileCmd("profile", "-f", configFile)
				_, err := cmd.Execute()
				Expect(err).ToNot(HaveOccurred())
			})

			It("fails without a cluster name", func() {
				cfg.Metadata.Name = ""
				configFile = CreateConfigFile(cfg)

				cmd := newMockGenerateProfileCmd("profile", "-f", configFile)
				_, err := cmd.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("metadata.name must be set"))
			})

			It("fails without a region", func() {
				cfg.Metadata.Region = ""
				configFile = CreateConfigFile(cfg)

				cmd := newMockGenerateProfileCmd("profile", "-f", configFile)
				_, err := cmd.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("metadata.region must be set"))
			})

			It("fails without bootstrap profiles", func() {
				cfg.Git.BootstrapProfile = nil
				configFile = CreateConfigFile(cfg)

				cmd := newMockGenerateProfileCmd("profile", "-f", configFile)
				_, err := cmd.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("git.bootstrapProfile.Source must be set"))
			})

			It("fails with empty bootstrap profile", func() {
				cfg.Git.BootstrapProfile = &api.Profile{}
				configFile = CreateConfigFile(cfg)

				cmd := newMockGenerateProfileCmd("profile", "-f", configFile)
				_, err := cmd.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("git.bootstrapProfile.Source must be set"))
			})

			It("defaults the outputPath to the profile name", func() {
				configFile = CreateConfigFile(cfg)

				cmd := newMockGenerateProfileCmd("profile", "-f", configFile)
				_, err := cmd.Execute()
				Expect(err).ToNot(HaveOccurred())

				cfg := cmd.Cmd.ClusterConfig

				Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
				Expect(cfg.Git.BootstrapProfile.OutputPath).To(Equal("./app-dev"))
			})

			It("defaults the outputPath to the repo name", func() {
				cfg.Git.BootstrapProfile.Source = "git@github.com:some-org/some-repo.git"
				configFile = CreateConfigFile(cfg)

				cmd := newMockGenerateProfileCmd("profile", "-f", configFile)
				_, err := cmd.Execute()
				Expect(err).ToNot(HaveOccurred())

				cfg := cmd.Cmd.ClusterConfig

				Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
				Expect(cfg.Git.BootstrapProfile.OutputPath).To(Equal("./some-repo"))
			})
		})
	})
})
