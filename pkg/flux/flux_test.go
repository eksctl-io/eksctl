package flux_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/executor/fakes"
	"github.com/weaveworks/eksctl/pkg/flux"
)

var _ = Describe("Flux", func() {
	var (
		fakeExecutor *fakes.FakeExecutor
		fluxClient   *flux.Client

		binDir string
	)

	BeforeEach(func() {
		fakeExecutor = new(fakes.FakeExecutor)
		fluxClient = flux.NewClient("", "github")
		fluxClient.SetExecutor(fakeExecutor)

		binDir, err := ioutil.TempDir("", "bin")
		Expect(err).NotTo(HaveOccurred())
		f, err := os.Create(filepath.Join(binDir, "flux"))
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Chmod(f.Name(), 0777)).To(Succeed())
		Expect(os.Setenv("PATH", binDir)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(binDir)).To(Succeed())
	})

	Context("PreFlight", func() {
		It("executes the Flux binary with the correct args", func() {
			Expect(fluxClient.PreFlight()).To(Succeed())
			Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
			_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
			Expect(receivedArgs).To(Equal([]string{"check", "--pre"}))
		})

		When("the flux binary is not found on the path", func() {
			BeforeEach(func() {
				Expect(os.Unsetenv("PATH")).To(Succeed())
			})

			It("returns the error", func() {
				Expect(fluxClient.PreFlight()).To(MatchError("flux not found, required"))
			})
		})

		When("execution fails", func() {
			BeforeEach(func() {
				fakeExecutor.ExecReturns(errors.New("omg"))
			})

			It("returns the error", func() {
				Expect(fluxClient.PreFlight()).To(MatchError("omg"))
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
			})
		})
	})

	Context("Bootstrap", func() {
		var (
			opts         *api.Flux
			standardArgs []string
		)

		BeforeEach(func() {
			opts = &api.Flux{
				GitProvider: "github",
				Repository:  "some-repo",
				Owner:       "theadversary_destroyerofkings_angelofthebottomlesspit_princeofthisworld_and_lordofdarkness",
			}

			standardArgs = []string{"bootstrap", opts.GitProvider, "--repository", opts.Repository, "--owner", opts.Owner}
		})

		It("executes the Flux binary with the correct default args", func() {
			Expect(fluxClient.Bootstrap(opts)).To(Succeed())
			Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
			_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
			Expect(receivedArgs).To(Equal(standardArgs))
		})

		When("opts.Personal is true", func() {
			BeforeEach(func() {
				opts.Personal = true
			})

			It("executes the Flux binary with the --personal flag", func() {
				Expect(fluxClient.Bootstrap(opts)).To(Succeed())
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
				_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
				Expect(receivedArgs).To(Equal(append(standardArgs, "--personal")))
			})
		})

		When("opts.Path is set", func() {
			BeforeEach(func() {
				opts.Path = "road_to_somewhere"
			})

			It("executes the Flux binary with the --path flag", func() {
				Expect(fluxClient.Bootstrap(opts)).To(Succeed())
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
				_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
				Expect(receivedArgs).To(Equal(append(standardArgs, "--path", opts.Path)))
			})
		})

		When("opts.Branch is set", func() {
			BeforeEach(func() {
				opts.Branch = "more-of-twig-really"
			})

			It("executes the Flux binary with the --path flag", func() {
				Expect(fluxClient.Bootstrap(opts)).To(Succeed())
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
				_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
				Expect(receivedArgs).To(Equal(append(standardArgs, "--branch", opts.Branch)))
			})
		})

		When("opts.Namespace is set", func() {
			BeforeEach(func() {
				opts.Namespace = "socially-distanced-space"
			})

			It("executes the Flux binary with the --path flag", func() {
				Expect(fluxClient.Bootstrap(opts)).To(Succeed())
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
				_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
				Expect(receivedArgs).To(Equal(append(standardArgs, "--namespace", opts.Namespace)))
			})
		})

		When("execution fails", func() {
			BeforeEach(func() {
				fakeExecutor.ExecReturns(errors.New("omg"))
			})

			It("returns the error", func() {
				Expect(fluxClient.Bootstrap(opts)).To(MatchError("omg"))
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
			})
		})
	})
})
