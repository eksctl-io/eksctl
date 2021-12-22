package karpenter

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers/fakes"
)

var _ = Describe("Install", func() {

	Context("Install", func() {

		var (
			fakeHelmInstaller  *fakes.FakeHelmInstaller
			installerUnderTest *Installer
			cfg                *api.ClusterConfig
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			cfg.Metadata.Name = "test-cluster"
			cfg.Karpenter = &api.Karpenter{
				Version:              "0.4.3",
				CreateServiceAccount: api.Disabled(),
			}
			cfg.Status = &api.ClusterStatus{
				Endpoint: "https://endpoint.com",
			}
			fakeHelmInstaller = &fakes.FakeHelmInstaller{}
			installerUnderTest = &Installer{
				Options: Options{
					HelmInstaller: fakeHelmInstaller,
					Namespace:     "karpenter",
					ClusterConfig: cfg,
				},
			}
		})

		It("installs karpenter into an existing cluster", func() {
			Expect(installerUnderTest.Install(context.Background(), "")).To(Succeed())
			_, args := fakeHelmInstaller.InstallChartArgsForCall(0)
			values := map[string]interface{}{
				controller: map[string]interface{}{
					clusterName:     cfg.Metadata.Name,
					clusterEndpoint: cfg.Status.Endpoint,
				},
				serviceAccount: map[string]interface{}{
					create: api.IsEnabled(cfg.Karpenter.CreateServiceAccount),
				},
			}
			Expect(args).To(Equal(providers.InstallChartOpts{
				ChartName:       "karpenter/karpenter",
				CreateNamespace: true,
				Namespace:       "karpenter",
				ReleaseName:     "karpenter",
				Values:          values,
				Version:         "0.4.3",
			}))
		})
		When("add repo fails", func() {

			BeforeEach(func() {
				fakeHelmInstaller.AddRepoReturns(errors.New("nope"))
			})

			It("errors", func() {
				Expect(installerUnderTest.Install(context.Background(), "")).
					To(MatchError(ContainSubstring("failed to add Karpenter repository: nope")))
			})
		})
		When("install chart fails", func() {

			BeforeEach(func() {
				fakeHelmInstaller.AddRepoReturns(nil)
				fakeHelmInstaller.InstallChartReturns(errors.New("nope"))
			})

			It("errors", func() {
				Expect(installerUnderTest.Install(context.Background(), "")).
					To(MatchError(ContainSubstring("failed to install Karpenter chart: nope")))
			})
		})

		When("the cluster configuration has fargate configured", func() {
			BeforeEach(func() {
				profile := &api.FargateProfile{
					Selectors: []api.FargateProfileSelector{
						{Namespace: "default"},
					},
				}
				cfg.FargateProfiles = []*api.FargateProfile{
					profile,
				}
				installerUnderTest = &Installer{
					Options: Options{
						HelmInstaller: fakeHelmInstaller,
						Namespace:     "karpenter",
						ClusterConfig: cfg,
					},
				}
			})
			It("will not tell helm to create a namespace", func() {
				Expect(installerUnderTest.Install(context.Background(), "")).To(Succeed())
				_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
				Expect(opts.CreateNamespace).To(BeFalse())
			})
		})

		When("there are no fargate profiles configured for the cluster", func() {
			BeforeEach(func() {
				installerUnderTest = &Installer{
					Options: Options{
						HelmInstaller: fakeHelmInstaller,
						Namespace:     "karpenter",
						ClusterConfig: cfg,
					},
				}
			})
			It("will tell helm to create the namespace", func() {
				Expect(installerUnderTest.Install(context.Background(), "")).To(Succeed())
				_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
				Expect(opts.CreateNamespace).To(BeTrue())
			})
		})

		When("service account is defined", func() {
			It("add role to the values for the helm chart", func() {
				Expect(installerUnderTest.Install(context.Background(), "role/account")).To(Succeed())
				_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
				values := map[string]interface{}{
					controller: map[string]interface{}{
						clusterName:     cfg.Metadata.Name,
						clusterEndpoint: cfg.Status.Endpoint,
					},
					serviceAccount: map[string]interface{}{
						create: api.IsEnabled(cfg.Karpenter.CreateServiceAccount),
						serviceAccountAnnotation: map[string]interface{}{
							api.AnnotationEKSRoleARN: "role/account",
						},
					},
				}
				Expect(opts.Values).To(Equal(values))
			})
		})
	})
})
