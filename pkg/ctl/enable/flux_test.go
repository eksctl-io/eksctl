package enable

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("enable flux", func() {
	var mockEnableFluxCmd func(args ...string) *ctltest.MockCmd

	BeforeEach(func() {
		mockEnableFluxCmd = func(args ...string) *ctltest.MockCmd {
			return ctltest.NewMockCmd(configureAndRun, "enable", args...)
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

			cmd *ctltest.MockCmd
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
						Flags: api.FluxFlags{
							"repository": "repo1",
							"owner":      "username",
						},
					},
				},
			}
		})

		JustBeforeEach(func() {
			configFile = ctltest.CreateConfigFile(cfg)
			cmd = mockEnableFluxCmd("flux", "-f", configFile)
			_, err = cmd.Execute()
		})

		AfterEach(func() {
			Expect(os.Remove(configFile)).To(Succeed())
		})

		It("succeeds with the basic configuration", func() {
			Expect(err).NotTo(HaveOccurred())

			fluxCfg := cmd.Cmd.ClusterConfig.GitOps.Flux
			Expect(fluxCfg).NotTo(BeNil())
			Expect(fluxCfg.GitProvider).To(Equal("github"))
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

		When("gitops is not provided", func() {
			BeforeEach(func() {
				cfg.GitOps = nil
			})

			It("fails", func() {
				Expect(err).To(MatchError("gitops.flux must be set"))
			})
		})

		When("gitops.flux is not provided", func() {
			BeforeEach(func() {
				cfg.GitOps.Flux = nil
			})

			It("fails", func() {
				Expect(err).To(MatchError("gitops.flux must be set"))
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

		When("gitops.flux.flags are not provided", func() {
			BeforeEach(func() {
				cfg.GitOps.Flux.Flags = api.FluxFlags{}
			})

			It("fails", func() {
				Expect(err).To(MatchError("gitops.flux.flags must be set"))
			})
		})
	})
})
