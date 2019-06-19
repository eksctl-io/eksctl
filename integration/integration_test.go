// +build integration

package integration_test

import (
	"flag"
	"testing"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var (
	eksctlPath string

	region  string
	version string

	// Flags to help with the development of the integration tests
	clusterName    string
	doCreate       bool
	doDelete       bool
	kubeconfigPath string

	kubeconfigTemp bool
)

func init() {
	flag.StringVar(&eksctlPath, "eksctl.path", "../eksctl", "Path to eksctl")

	flag.StringVar(&region, "eksctl.region", api.DefaultRegion, "Region to use for the tests")
	flag.StringVar(&version, "eksctl.version", api.DefaultVersion, "Version of Kubernetes to test")

	// Flags to help with the development of the integration tests
	flag.StringVar(&clusterName, "eksctl.cluster", "", "Cluster name (default: generate one)")
	flag.BoolVar(&doCreate, "eksctl.create", true, "Skip the creation tests. Useful for debugging the tests")
	flag.BoolVar(&doDelete, "eksctl.delete", true, "Skip the cleanup after the tests have run")
	flag.StringVar(&kubeconfigPath, "eksctl.kubeconfig", "", "Path to kubeconfig (default: create it a temporary file)")
}

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}
