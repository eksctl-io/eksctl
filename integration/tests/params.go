// +build integration

package tests

import (
	"flag"
	"fmt"
	"time"

	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/integration/runner"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/names"
)

// Params groups all test parameters.
type Params struct {
	EksctlPath string
	Region     string
	Version    string
	// Flags to help with the development of the integration tests
	ClusterName    string
	DoCreate       bool
	DoDelete       bool
	KubeconfigPath string
	KubeconfigTemp bool
	TestDirectory  string
	// privateSSHKeyPath is the SSH key to use for Git operations.
	PrivateSSHKeyPath       string
	EksctlCmd               runner.Cmd
	EksctlCreateCmd         runner.Cmd
	EksctlGetCmd            runner.Cmd
	EksctlDeleteCmd         runner.Cmd
	EksctlDeleteClusterCmd  runner.Cmd
	EksctlScaleNodeGroupCmd runner.Cmd
	EksctlUtilsCmd          runner.Cmd
	EksctlExperimentalCmd   runner.Cmd
	// Keep track of created clusters, for post-tests clean-up.
	clustersToDelete []string
}

// NewClusterName generates a new cluster name using the provided prefix, and
// adds the cluster to the list of clusters to eventually delete, once the test
// suite has run.
func (p *Params) NewClusterName(prefix string) string {
	clusterName := fmt.Sprintf("it-%s-%s", prefix, names.ForCluster("", ""))
	p.addToDeleteList(clusterName)
	return clusterName
}

// addToDeleteList adds the provided cluster name to the list of clusters to eventually delete.
func (p *Params) addToDeleteList(clusterName string) {
	p.clustersToDelete = append(p.clustersToDelete, clusterName)
}

func (p Params) DeleteClusters() {
	if !p.DoDelete {
		return
	}
	for _, clusterName := range p.clustersToDelete {
		cmd := p.EksctlDeleteClusterCmd.WithArgs(
			"--name", clusterName,
		)
		Expect(cmd).To(runner.RunSuccessfully())
	}
}

const (
	defaultTestDirectory     = "test_profile"
	defaultPrivateSSHKeyPath = "~/.ssh/eksctl-bot_id_rsa"
)

// NewParams creates a new Test instance from CLI args, grouping all test parameters.
func NewParams(clusterNamePrefix string) *Params {
	var params Params

	flag.StringVar(&params.EksctlPath, "eksctl.path", "../../../eksctl", "Path to eksctl")
	flag.StringVar(&params.Region, "eksctl.region", api.DefaultRegion, "Region to use for the tests")
	flag.StringVar(&params.Version, "eksctl.version", api.DefaultVersion, "Version of Kubernetes to test")
	flag.StringVar(&params.TestDirectory, "eksctl.test.dir", defaultTestDirectory, "Test directory. Defaulted to: "+defaultTestDirectory)
	// Flags to help with the development of the integration tests
	flag.StringVar(&params.ClusterName, "eksctl.cluster", "", "Cluster name (default: generate one)")
	flag.BoolVar(&params.DoCreate, "eksctl.create", true, "Skip the creation tests. Useful for debugging the tests")
	flag.BoolVar(&params.DoDelete, "eksctl.delete", true, "Skip the cleanup after the tests have run")
	flag.StringVar(&params.KubeconfigPath, "eksctl.kubeconfig", "", "Path to kubeconfig (default: create a temporary file)")
	flag.StringVar(&params.PrivateSSHKeyPath, "eksctl.git.sshkeypath", defaultPrivateSSHKeyPath, fmt.Sprintf("Path to the SSH key to use for Git operations (default: %s)", defaultPrivateSSHKeyPath))

	// go1.13+ testing flags regression fix: https://github.com/golang/go/issues/31859
	flag.Parse()

	if params.ClusterName == "" {
		params.ClusterName = params.NewClusterName(clusterNamePrefix)
	} else {
		params.addToDeleteList(params.ClusterName)
	}

	params.EksctlCmd = runner.NewCmd(params.EksctlPath).
		WithArgs("--region", params.Region).
		WithTimeout(30 * time.Minute)

	params.EksctlCreateCmd = params.EksctlCmd.
		WithArgs("create").
		WithTimeout(25 * time.Minute)

	params.EksctlGetCmd = params.EksctlCmd.
		WithArgs("get").
		WithTimeout(1 * time.Minute)

	params.EksctlDeleteCmd = params.EksctlCmd.
		WithArgs("delete").
		WithTimeout(15 * time.Minute)

	params.EksctlDeleteClusterCmd = params.EksctlDeleteCmd.
		WithArgs("cluster", "--verbose", "4")

	params.EksctlScaleNodeGroupCmd = params.EksctlCmd.
		WithArgs("scale", "nodegroup", "--verbose", "4").
		WithTimeout(5 * time.Minute)

	params.EksctlUtilsCmd = params.EksctlCmd.
		WithArgs("utils").
		WithTimeout(5 * time.Minute)

	params.EksctlExperimentalCmd = params.EksctlCmd.
		WithEnv("EKSCTL_EXPERIMENTAL=true")

	return &params
}
