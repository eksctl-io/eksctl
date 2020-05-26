package enable

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("enable profile", func() {

	newMockEnableProfileCmd := func(args ...string) *MockCmd {
		return NewMockCmd(enableProfileWithRunFunc, "enable", args...)
	}

	Describe("without a config file", func() {

		It("should accept a name argument", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "app-dev")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.Source).To(Equal("app-dev"))
		})

		It("should accept a --profile-source flag", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "--profile-source", "app-dev")
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.Source).To(Equal("app-dev"))
		})

		It("should reject name argument and --profile-source flag", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "--profile-source", "app-dev", "app-dev")
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--profile-source=app-dev and argument app-dev cannot be used at the same time"))
		})

		It("requires the --profile-source flag", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--profile-source must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("requires the --git-email flag", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--git-url", "git@example.com:repo.git", "app-dev")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--git-email must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("requires the --git-url flag", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--region", "eu-north-1", "--git-email", "user@example.com", "app-dev")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--git-url must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("requires the --cluster flag", func() {
			cmd := newMockEnableProfileCmd("profile", "--region", "eu-north-1", "--git-email", "user@example.com", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "app-dev")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})
		It("requires the --region flag", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--git-email", "user@example.com", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "app-dev")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--region must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})
		It("requires the git-private-ssh-key file to exist", func() {
			cmd := newMockEnableProfileCmd("profile", "--cluster", "clus-1", "--git-email", "user@example.com", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "--git-private-ssh-key-path", "./inexistent-file", "app-dev")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--region must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("loads all flags correctly", func() {
			cmd := newMockEnableProfileCmd("profile",
				"--cluster", "clus-1",
				"--region", "us-west-2",
				"--git-url", "git@example.com:repo.git",
				"--git-email", "user@example.com",
				"--git-branch", "master",
				"--git-user", "user1",
				"--git-private-ssh-key-path", "./profile_test.go",
				"--profile-source", "app-dev",
				"--profile-revision", "branch-2",
			)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			cfg := cmd.Cmd.ClusterConfig
			Expect(cfg.Metadata.Name).To(Equal("clus-1"))
			Expect(cfg.Metadata.Region).To(Equal("us-west-2"))
			Expect(cfg.Git).ToNot(BeNil())
			Expect(cfg.Git.Repo).ToNot(BeNil())
			Expect(cfg.Git.Repo.URL).To(Equal("git@example.com:repo.git"))
			Expect(cfg.Git.Repo.Email).To(Equal("user@example.com"))
			Expect(cfg.Git.Repo.Branch).To(Equal("master"))
			Expect(cfg.Git.Repo.User).To(Equal("user1"))
			Expect(cfg.Git.Repo.PrivateSSHKeyPath).To(Equal("./profile_test.go"))

			Expect(cfg.Git.BootstrapProfile).ToNot(BeNil())
			Expect(cfg.Git.BootstrapProfile.Source).To(Equal("app-dev"))
			Expect(cfg.Git.BootstrapProfile.Revision).To(Equal("branch-2"))
		})
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
					Repo: &api.Repo{
						URL:   "git@github.com:org/repo1",
						Email: "user@example.com",
					},
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

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("loads the correct defaults", func() {
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())

			gitCfg := cmd.Cmd.ClusterConfig.Git
			Expect(gitCfg).ToNot(BeNil())
			Expect(gitCfg.Repo.Branch).To(Equal("master"))
			Expect(gitCfg.Repo.User).To(Equal("Flux"))
			Expect(gitCfg.Repo.FluxPath).To(Equal("flux/"))
			Expect(gitCfg.Repo.Paths).To(BeEmpty())
			Expect(gitCfg.Repo.PrivateSSHKeyPath).To(Equal(""))

			Expect(gitCfg.Operator.Namespace).To(Equal("flux"))
			Expect(gitCfg.Operator.Label).To(Equal("flux"))
			Expect(gitCfg.Operator.WithHelm).ToNot(BeNil())
			Expect(*gitCfg.Operator.WithHelm).To(BeTrue())

			Expect(gitCfg.BootstrapProfile.Revision).To(Equal(""))
			Expect(gitCfg.BootstrapProfile.OutputPath).To(Equal("./app-dev"))
		})

		It("fails without a cluster name", func() {
			cfg.Metadata.Name = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.name must be set"))
		})

		It("fails without a region", func() {
			cfg.Metadata.Region = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.region must be set"))
		})

		It("fails without a nil repo", func() {
			cfg.Git.Repo = nil
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("git.repo.URL must be set"))
		})

		It("fails without a git url", func() {
			cfg.Git.Repo.URL = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("git.repo.URL must be set"))
		})

		It("fails without a user email", func() {
			cfg.Git.Repo.Email = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("git.repo.email must be set"))
		})

		It("fails when the private ssh key file does not exist", func() {
			cfg.Git.Repo.PrivateSSHKeyPath = "non-existent-file"
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("please supply a valid file for git.repo.privateSSHKeyPath: invalid path to private SSH key: non-existent-file"))
		})

		It("fails without bootstrap profiles", func() {
			cfg.Git.BootstrapProfile = nil
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("git.bootstrapProfile.Source must be set"))
		})

		It("fails with empty bootstrap profile", func() {
			cfg.Git.BootstrapProfile = &api.Profile{}
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableProfileCmd("profile", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("git.bootstrapProfile.Source must be set"))
		})
	})

})
