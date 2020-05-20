package enable

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("enable repo", func() {

	newMockEnableRepoCmd := func(args ...string) *MockCmd {
		return NewMockCmd(enableRepoWithRunFunc, "enable", args...)
	}

	Describe("without a config file", func() {

		It("with name argument should fail", func() {
			cmd := newMockEnableRepoCmd("repo", "--cluster", "clus-1", "--region", "eu-north-1", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "name-argument")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name argument is not supported"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("requires the --git-email flag", func() {
			cmd := newMockEnableRepoCmd("repo", "--cluster", "clus-1", "--region", "eu-north-1", "--git-url", "git@example.com:repo.git")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--git-email must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("requires the --git-url flag", func() {
			cmd := newMockEnableRepoCmd("repo", "--cluster", "clus-1", "--region", "eu-north-1", "--git-email", "user@example.com")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--git-url must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("requires the --cluster flag", func() {
			cmd := newMockEnableRepoCmd("repo", "--region", "eu-north-1", "--git-email", "user@example.com", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})
		It("requires the --region flag", func() {
			cmd := newMockEnableRepoCmd("repo", "--cluster", "clus-1", "--git-email", "user@example.com", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--region must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})
		It("requires the git-private-ssh-key file to exist", func() {
			cmd := newMockEnableRepoCmd("repo", "--cluster", "clus-1", "--git-email", "user@example.com", "--git-url", "git@example.com:repo.git", "--git-email", "user@example.com", "--git-private-ssh-key-path", "./inexistent-file")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--region must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})

		It("loads all flags correctly", func() {
			cmd := newMockEnableRepoCmd("repo",
				"--cluster", "clus-1",
				"--region", "us-west-2",
				"--git-url", "git@example.com:repo.git",
				"--git-email", "user@example.com",
				"--git-branch", "master",
				"--git-user", "user1",
				"--git-private-ssh-key-path", "./repo_test.go",
				"--git-paths", "base,flux,upgrades",
				"--git-label", "flux2",
				"--git-flux-subdir", "flux-dir/",
				"--namespace", "gitops",
				"--with-helm",
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
			Expect(cfg.Git.Repo.PrivateSSHKeyPath).To(Equal("./repo_test.go"))
			Expect(cfg.Git.Repo.Paths).To(ConsistOf("base", "flux", "upgrades"))
			Expect(cfg.Git.Repo.FluxPath).To(Equal("flux-dir/"))

			Expect(cfg.Git.Operator).ToNot(BeNil())
			Expect(cfg.Git.Operator.Label).To(Equal("flux2"))
			Expect(cfg.Git.Operator.Namespace).To(Equal("gitops"))
			Expect(*cfg.Git.Operator.WithHelm).To(BeTrue())
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

			cmd := newMockEnableRepoCmd("repo", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails without a cluster name", func() {
			cfg.Metadata.Name = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableRepoCmd("repo", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.name must be set"))
		})

		It("fails without a region", func() {
			cfg.Metadata.Region = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableRepoCmd("repo", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata.region must be set"))
		})

		It("fails without a git url", func() {
			cfg.Git.Repo.URL = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableRepoCmd("repo", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("git.repo.URL must be set"))
		})

		It("fails without a user email", func() {
			cfg.Git.Repo.Email = ""
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableRepoCmd("repo", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("git.repo.email must be set"))
		})

		It("fails when the private ssh key file does not exist", func() {
			cfg.Git.Repo.PrivateSSHKeyPath = "non-existent-file"
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableRepoCmd("repo", "-f", configFile)
			_, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("please supply a valid file for git.repo.privateSSHKeyPath: invalid path to private SSH key: non-existent-file"))
		})

		It("loads the correct defaults", func() {
			configFile = CreateConfigFile(cfg)

			cmd := newMockEnableRepoCmd("repo", "-f", configFile)
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
		})

	})
})
