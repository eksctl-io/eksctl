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
				Version:                "0.4.3",
				CreateServiceAccount:   api.Disabled(),
				DefaultInstanceProfile: nil,
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
			Expect(installerUnderTest.Install(context.Background(), "role-arn", "role/profile")).To(Succeed())
			_, args := fakeHelmInstaller.InstallChartArgsForCall(0)
			values := map[string]interface{}{
				clusterName:     cfg.Metadata.Name,
				clusterEndpoint: cfg.Status.Endpoint,
				serviceAccount: map[string]interface{}{
					create: api.IsEnabled(cfg.Karpenter.CreateServiceAccount),
					serviceAccountAnnotation: map[string]interface{}{
						api.AnnotationEKSRoleARN: "role-arn",
					},
					serviceAccountName: DefaultServiceAccountName,
				},
				aws: map[string]interface{}{
					defaultInstanceProfile: "role/profile",
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
				Expect(installerUnderTest.Install(context.Background(), "", "role/profile")).
					To(MatchError(ContainSubstring("failed to add Karpenter repository: nope")))
			})
		})
		When("install chart fails", func() {

			BeforeEach(func() {
				fakeHelmInstaller.AddRepoReturns(nil)
				fakeHelmInstaller.InstallChartReturns(errors.New("nope"))
			})

			It("errors", func() {
				Expect(installerUnderTest.Install(context.Background(), "", "role/profile")).
					To(MatchError(ContainSubstring("failed to install Karpenter chart: nope")))
			})
		})

		When("service account is defined", func() {
			It("add role to the values for the helm chart", func() {
				Expect(installerUnderTest.Install(context.Background(), "role/account", "role/profile")).To(Succeed())
				_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
				values := map[string]interface{}{
					clusterName:     cfg.Metadata.Name,
					clusterEndpoint: cfg.Status.Endpoint,
					serviceAccount: map[string]interface{}{
						create: api.IsEnabled(cfg.Karpenter.CreateServiceAccount),
						serviceAccountAnnotation: map[string]interface{}{
							api.AnnotationEKSRoleARN: "role/account",
						},
						serviceAccountName: DefaultServiceAccountName,
					},
					aws: map[string]interface{}{
						defaultInstanceProfile: "role/profile",
					},
				}
				Expect(opts.Values).To(Equal(values))
			})
		})
	})
})
