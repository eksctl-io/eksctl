package karpenter

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/karpenter/providers/fakes"
)

var _ = Describe("InstallKarpenter", func() {
	Context("Install", func() {
		var (
			fakeHelmInstaller *fakes.FakeHelmInstaller
		)
		BeforeEach(func() {
			fakeHelmInstaller = &fakes.FakeHelmInstaller{}
		})
		AfterEach(func() {})
		It("installs karpenter into an existing cluster", func() {
			fakeHelmInstaller.AddRepoReturns(nil)
			fakeHelmInstaller.InstallChartReturns(nil)
			ki := &Installer{
				Options: Options{
					HelmInstaller:         fakeHelmInstaller,
					Namespace:             "karpenter",
					ClusterName:           "test-cluster",
					AddDefaultProvisioner: true,
					CreateServiceAccount:  true,
					ClusterEndpoint:       "https://endpoint.com",
					Version:               "0.4.3",
				},
			}
			err := ki.InstallKarpenter(context.Background())
			Expect(err).NotTo(HaveOccurred())
		})
		When("add repo fails", func() {
			It("returns an error", func() {
				fakeHelmInstaller.AddRepoReturns(errors.New("nope"))
				ki := &Installer{
					Options: Options{
						HelmInstaller:         fakeHelmInstaller,
						Namespace:             "karpenter",
						ClusterName:           "test-cluster",
						AddDefaultProvisioner: true,
						CreateServiceAccount:  true,
						ClusterEndpoint:       "https://endpoint.com",
						Version:               "0.4.3",
					},
				}
				err := ki.InstallKarpenter(context.Background())
				Expect(err).To(MatchError(ContainSubstring("failed to karpenter repo: nope")))
			})
		})
		When("install chart fails", func() {
			It("returns an error", func() {
				fakeHelmInstaller.AddRepoReturns(nil)
				fakeHelmInstaller.InstallChartReturns(errors.New("nope"))
				ki := &Installer{
					Options: Options{
						HelmInstaller:         fakeHelmInstaller,
						Namespace:             "karpenter",
						ClusterName:           "test-cluster",
						AddDefaultProvisioner: true,
						CreateServiceAccount:  true,
						ClusterEndpoint:       "https://endpoint.com",
						Version:               "0.4.3",
					},
				}
				err := ki.InstallKarpenter(context.Background())
				Expect(err).To(MatchError(ContainSubstring("failed to install karpenter chart: nope")))
			})
		})
	})
})
