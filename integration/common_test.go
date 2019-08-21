// +build integration

package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

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
	test := h.NewTest(t)
	test.Setup()
	return test, nil
}

func eksctl(args ...string) *gexec.Session {
	command := exec.Command(eksctlPath, args...)
	command.Env = os.Environ()
	command.Env = append(command.Env, "EKSCTL_EXPERIMENTAL=true")
	fmt.Fprintf(GinkgoWriter, "calling %q with %s and %v\n", eksctlPath, "EKSCTL_EXPERIMENTAL=true", args)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).To(BeNil())

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
	return session
}

func eksctlSuccess(args ...string) *gexec.Session {
	session := eksctl(args...)
	Expect(session.ExitCode()).To(BeZero())
	return session
}

func eksctlFail(args ...string) *gexec.Session {
	session := eksctl(args...)
	Expect(session.ExitCode()).ToNot(BeZero())
	return session
}

//eksctlStart starts running an eksctl command but doesn't wait for it to finish the command
//This is primarily so that we can run eksctl create and then subsequently call eksctl delete
//on the same cluster, but might be useful for other test scenarios as well.
func eksctlStart(args ...string) {
	fmt.Fprintf(GinkgoWriter, "calling %q with %v\n", eksctlPath, args)
	cmd := exec.Command(eksctlPath, args...)
	err := cmd.Start()
	Expect(err).To(BeNil())
}
