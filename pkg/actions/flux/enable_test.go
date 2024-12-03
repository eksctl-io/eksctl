package flux_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/flux"
	"github.com/weaveworks/eksctl/pkg/actions/flux/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Gitops", func() {
	var (
		fakeClientSet  *fake.Clientset
		fakeFluxClient *fakes.FakeInstallerClient
		opts           *api.GitOps
		installer      *flux.Installer
	)

	BeforeEach(func() {
		fakeClientSet = fake.NewSimpleClientset()
		fakeFluxClient = new(fakes.FakeInstallerClient)
		opts = &api.GitOps{Flux: &api.Flux{}}
	})

	JustBeforeEach(func() {
		var err error
		installer, err = flux.New(fakeClientSet, opts)
		Expect(err).NotTo(HaveOccurred())
		installer.SetFluxClient(fakeFluxClient)
	})

	It("installs Flux v2", func() {
		Expect(installer.Run()).To(Succeed())
		Expect(fakeFluxClient.PreFlightCallCount()).To(Equal(1))
		Expect(fakeFluxClient.BootstrapCallCount()).To(Equal(1))
	})

	Context("Flux v2 pre-check execution fails", func() {
		BeforeEach(func() {
			fakeFluxClient.PreFlightReturns(errors.New("flux cli not installed"))
		})

		It("returns the error from the Flux CLI", func() {
			Expect(installer.Run()).To(MatchError(ContainSubstring("flux cli not installed")))
		})
	})

	Context("Flux v2 bootstrap execution fails", func() {
		BeforeEach(func() {
			fakeFluxClient.BootstrapReturns(errors.New("this totally failed"))
		})

		It("returns the error from the Flux CLI", func() {
			Expect(installer.Run()).To(MatchError(ContainSubstring("this totally failed")))
		})
	})

	Context("Flux v1 components are already installed", func() {
		BeforeEach(func() {
			_, err := fakeClientSet.AppsV1().Deployments("flux-system").Create(context.Background(), &v1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "flux"}},
				metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not fail, and does not install v2 components", func() {
			Expect(installer.Run()).To(Succeed())
			Expect(fakeFluxClient.PreFlightCallCount()).To(Equal(0))
			Expect(fakeFluxClient.BootstrapCallCount()).To(Equal(0))
		})
	})
})
