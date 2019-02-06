// +build integration

package integration_test

import (
	"fmt"
	"os/exec"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type tInterface interface {
	GinkgoTInterface
	Helper()
}

type tHelper struct{ GinkgoTInterface }

func (t *tHelper) Helper()      { return }
func (t *tHelper) Name() string { return "eksctl-test" }

func newKubeTest() (*harness.Test, error) {
	t := &tHelper{GinkgoT()}
	l := &logger.TestLogger{}
	h := harness.New(harness.Options{Logger: l.ForTest(t)})
	if err := h.Setup(); err != nil {
		return nil, err
	}
	if err := h.SetKubeconfig(kubeconfigPath); err != nil {
		return nil, err
	}
	return h.NewTest(t), nil
}

func eksctl(args ...string) *gexec.Session {
	command := exec.Command(eksctlPath, args...)
	fmt.Fprintf(GinkgoWriter, "calling %q with %v\n", eksctlPath, args)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	if err != nil {
		Fail(fmt.Sprintf("error starting process: %v\n", err), 1)
	}

	t := time.Minute
	switch args[0] {
	case "create":
		t *= 25
	case "delete":
		t *= 15
	case "get":
		t *= 1
	case "scale":
		t *= 5
	default:
		t *= 30
	}

	session.Wait(t)
	Expect(session.ExitCode()).To(Equal(0))
	return session
}
