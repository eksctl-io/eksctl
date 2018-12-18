// +build integration

package integration_test

import (
	"flag"
	"testing"
	"time"

	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

const (
	createTimeout = 25 * time.Minute
	deleteTimeout = 15 * time.Minute
	getTimeout    = 1 * time.Minute
	scaleTimeout  = 5 * time.Minute
	region        = api.DefaultRegion
)

var (
	eksctlPath string

	// Flags to help with the development of the integration tests
	clusterName    string
	doCreate       bool
	doDelete       bool
	kubeconfigPath string

	kubeconfigTemp bool
)

func init() {
	flag.StringVar(&eksctlPath, "eksctl.path", "../eksctl", "Path to eksctl")

	// Flags to help with the development of the integration tests
	flag.StringVar(&clusterName, "eksctl.cluster", "", "Cluster name (default: generate one)")
	flag.BoolVar(&doCreate, "eksctl.create", true, "Skip the creation tests. Useful for debugging the tests")
	flag.BoolVar(&doDelete, "eksctl.delete", true, "Skip the cleanup after the tests have run")
	flag.StringVar(&kubeconfigPath, "eksctl.kubeconfig", "", "Path to kubeconfig (default: create it a temporary file)")
}

func TestCreateIntegration(t *testing.T) {
	testutils.RegisterAndRun(t, "(Integration) Create, Get, Scale & Delete")
}
