//go:build integration
// +build integration

package kube

import (
	"fmt"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"
	"github.com/onsi/ginkgo/v2"
)

type tHelper struct{ ginkgo.FullGinkgoTInterface }

func (t *tHelper) Helper()      {}
func (t *tHelper) Name() string { return "eksctl-test" }

// NewTest creates a new test harness to more easily run integration tests against the provided Kubernetes cluster.
func NewTest(kubeconfigPath string) (*harness.Test, error) {
	t := &tHelper{ginkgo.GinkgoT()}
	l := &logger.TestLogger{}
	h := harness.New(harness.Options{Logger: l.ForTest(t)})
	if err := h.Setup(); err != nil {
		return nil, err
	}
	if err := h.SetKubeconfig(kubeconfigPath); err != nil {
		return nil, err
	}
	test := h.NewTest(t)
	test.Namespace += fmt.Sprintf("-%d", t.ParallelProcess())
	test.Setup()
	return test, nil
}
