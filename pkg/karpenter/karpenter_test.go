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
			fakeHelmInstaller  *fakes.FakeHelmInstaller
			installerUnderTest *Installer
		)

		BeforeEach(func() {
			fakeHelmInstaller = &fakes.FakeHelmInstaller{}
			installerUnderTest = &Installer{
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
		})

		It("installs karpenter into an existing cluster", func() {
			Expect(installerUnderTest.InstallKarpenter(context.Background())).To(Succeed())
		})
		When("add repo fails", func() {

			BeforeEach(func() {
				fakeHelmInstaller.AddRepoReturns(errors.New("nope"))
			})

			It("errors", func() {
				Expect(installerUnderTest.InstallKarpenter(context.Background())).
					To(MatchError(ContainSubstring("failed to karpenter repo: nope")))
			})
		})
		When("install chart fails", func() {

			BeforeEach(func() {
				fakeHelmInstaller.AddRepoReturns(nil)
				fakeHelmInstaller.InstallChartReturns(errors.New("nope"))
			})

			It("errors", func() {
				Expect(installerUnderTest.InstallKarpenter(context.Background())).
					To(MatchError(ContainSubstring("failed to install karpenter chart: nope")))
			})
		})
	})
})
