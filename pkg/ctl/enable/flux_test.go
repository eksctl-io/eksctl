package enable

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("enable flux", func() {
	var mockEnableFluxCmd func(args ...string) *MockCmd

	BeforeEach(func() {
		mockEnableFluxCmd = func(args ...string) *MockCmd {
			return NewMockCmd(configureAndRun, "enable", args...)
		}
	})

	When("--config-file is not provided", func() {
		It("should fail", func() {
			cmd := mockEnableFluxCmd("flux")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("--config-file/-f <file> must be set"))
			Expect(out).To(ContainSubstring("Usage"))
		})
	})

	When("name arg is provided", func() {
		It("should fail", func() {
			cmd := mockEnableFluxCmd("flux", "foo")
			out, err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("name argument is not supported"))
			Expect(out).To(ContainSubstring("Usage"))
		})
	})

	When("--config-file is provided", func() {
		var (
			configFile string
			cfg        *api.ClusterConfig

			cmd *MockCmd
			err error
		)

		BeforeEach(func() {
			// Minimal valid cluster config for the command to work
			cfg = &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Name:   "cluster-1",
					Region: "us-west-2",
				},
				GitOps: &api.GitOps{
					Flux: &api.Flux{
						GitProvider: "github",
						Repository:  "repo1",
						Owner:       "username",
					},
				},
			}
		})

		JustBeforeEach(func() {
			configFile = CreateConfigFile(cfg)
			cmd = mockEnableFluxCmd("flux", "-f", configFile)
			_, err = cmd.Execute()
		})

		AfterEach(func() {
			Expect(os.Remove(configFile)).To(Succeed())
		})

		It("succeeds with the basic configuration", func() {
			Expect(err).ToNot(HaveOccurred())

			fluxCfg := cmd.Cmd.ClusterConfig.GitOps.Flux
			Expect(fluxCfg).ToNot(BeNil())
			Expect(fluxCfg.Repository).To(Equal("repo1"))
			Expect(fluxCfg.GitProvider).To(Equal("github"))
			Expect(fluxCfg.Owner).To(Equal("username"))
		})

		It("loads the correct default", func() {
			Expect(err).ToNot(HaveOccurred())

			fluxCfg := cmd.Cmd.ClusterConfig.GitOps.Flux
			Expect(fluxCfg).ToNot(BeNil())
			Expect(fluxCfg.Namespace).To(Equal("flux-system"))
		})

		When("metadata.cluster is not provided", func() {
			BeforeEach(func() {
				cfg.Metadata.Name = ""
			})

			It("fails", func() {
				Expect(err).To(MatchError("metadata.name must be set"))
			})
		})

		When("metadata.region is not provided", func() {
			BeforeEach(func() {
				cfg.Metadata.Region = ""
			})

			It("fails", func() {
				Expect(err).To(MatchError("metadata.region must be set"))
			})
		})

		When("gitops.flux is not provided", func() {
			BeforeEach(func() {
				cfg.GitOps.Flux = nil
			})

			It("fails", func() {
				Expect(err).To(MatchError("no configuration found for enable flux"))
			})
		})

		When("gitops.flux.gitProvider is not provided", func() {
			BeforeEach(func() {
				cfg.GitOps.Flux.GitProvider = ""
			})

			It("fails", func() {
				Expect(err).To(MatchError("gitops.flux.gitProvider must be set"))
			})
		})

		When("gitops.flux.repository is not provided", func() {
			BeforeEach(func() {
				cfg.GitOps.Flux.Repository = ""
			})

			It("fails", func() {
				Expect(err).To(MatchError("gitops.flux.repository must be set"))
			})
		})

		When("gitops.flux.owner is not provided", func() {
			BeforeEach(func() {
				cfg.GitOps.Flux.Owner = ""
			})

			It("fails", func() {
				Expect(err).To(MatchError("gitops.flux.owner must be set"))
			})
		})

		When("deprecated git configuration is provided", func() {
			BeforeEach(func() {
				cfg.Git = &api.Git{Repo: &api.Repo{}}
			})

			It("fails", func() {
				Expect(err).To(MatchError("config cannot be provided for git.repo, git.bootstrapProfile or git.operator alongside gitops.*"))
			})
		})
	})
})
