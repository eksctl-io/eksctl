package flux_test

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/executor/fakes"
	"github.com/weaveworks/eksctl/pkg/flux"
)

var _ = Describe("Flux", func() {
	var (
		fakeExecutor *fakes.FakeExecutor
		fluxClient   *flux.Client
		opts         *api.Flux

		binDir string
	)

	BeforeEach(func() {
		opts = &api.Flux{
			GitProvider: "github",
		}

		fakeExecutor = new(fakes.FakeExecutor)
		var err error
		fluxClient, err = flux.NewClient(opts)
		Expect(err).NotTo(HaveOccurred())
		fluxClient.SetExecutor(fakeExecutor)
		fakeExecutor.ExecWithOutReturns([]byte("flux version 2.0.0-rc.1\n"), nil)

		binDir, err := os.MkdirTemp("", "bin")
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

		When("a kubeconfig is provided in flags", func() {
			BeforeEach(func() {
				opts.Flags = api.FluxFlags{"kubeconfig": "some-path"}
			})

			It("sets the kubeconfig flag on the command", func() {
				Expect(fluxClient.PreFlight()).To(Succeed())
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
				_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
				Expect(receivedArgs).To(Equal([]string{"check", "--pre", "--kubeconfig", "some-path"}))
			})
		})

		When("a context is provided in flags", func() {
			BeforeEach(func() {
				opts.Flags = api.FluxFlags{"context": "foo"}
			})

			It("sets the kubeconfig flag on the command", func() {
				Expect(fluxClient.PreFlight()).To(Succeed())
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
				_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
				Expect(receivedArgs).To(Equal([]string{"check", "--pre", "--context", "foo"}))
			})
		})

		When("the flux binary is not found on the path", func() {
			BeforeEach(func() {
				Expect(os.Unsetenv("PATH")).To(Succeed())
			})

			It("returns the error", func() {
				Expect(fluxClient.PreFlight()).To(MatchError("flux not found, required"))
			})
		})

		Context("checking the flux version", func() {
			When("the flux version is < 0.32.0", func() {
				BeforeEach(func() {
					fakeExecutor.ExecWithOutReturns([]byte("flux version 0.31.5\n"), nil)
				})

				It("returns an error saying older versions are not supported", func() {
					Expect(fluxClient.PreFlight()).To(MatchError(ContainSubstring("found flux version 0.31.5, eksctl requires >= 0.32.0")))
				})
			})

			When("the flux version command returns unexpected output", func() {
				BeforeEach(func() {
					fakeExecutor.ExecWithOutReturns([]byte("hmmm"), nil)
				})

				It("returns an error", func() {
					Expect(fluxClient.PreFlight()).To(MatchError("unexpected format returned from 'flux --version': [hmmm]"))
				})
			})

			When("the flux version is not valid semver", func() {
				BeforeEach(func() {
					fakeExecutor.ExecWithOutReturns([]byte("flux version a.b.c"), nil)
				})

				It("returns an error", func() {
					Expect(fluxClient.PreFlight()).To(MatchError(ContainSubstring("failed to parse Flux version")))
				})
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
			standardArgs []string
		)

		BeforeEach(func() {
			standardArgs = []string{"bootstrap", opts.GitProvider}
		})

		It("executes the Flux binary with the correct subcommands", func() {
			Expect(fluxClient.Bootstrap()).To(Succeed())
			Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
			_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
			Expect(receivedArgs).To(Equal(standardArgs))
		})

		When("opts.Flags are set", func() {
			BeforeEach(func() {
				opts.Flags = api.FluxFlags{"foo": "bar"}
			})

			It("parses the flags and appends them to the command args", func() {
				Expect(fluxClient.Bootstrap()).To(Succeed())
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
				_, receivedArgs := fakeExecutor.ExecArgsForCall(0)
				Expect(receivedArgs).To(Equal(append(standardArgs, "--foo", "bar")))
			})
		})

		When("execution fails", func() {
			BeforeEach(func() {
				fakeExecutor.ExecReturns(errors.New("omg"))
			})

			It("returns the error", func() {
				Expect(fluxClient.Bootstrap()).To(MatchError("omg"))
				Expect(fakeExecutor.ExecCallCount()).To(Equal(1))
			})
		})
	})
})
