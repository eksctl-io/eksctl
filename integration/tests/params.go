// +build integration

package tests

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/onsi/ginkgo"
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
	clusterNamePrefix string
	ClusterName       string
	SkipCreate        bool
	SkipDelete        bool
	KubeconfigPath    string
	KubeconfigTemp    bool
	TestDirectory     string
	// privateSSHKeyPath is the SSH key to use for Git operations.
	PrivateSSHKeyPath       string
	EksctlCmd               runner.Cmd
	EksctlCreateCmd         runner.Cmd
	EksctlUpdateCmd         runner.Cmd
	EksctlUpgradeCmd        runner.Cmd
	EksctlGetCmd            runner.Cmd
	EksctlDeleteCmd         runner.Cmd
	EksctlDeleteClusterCmd  runner.Cmd
	EksctlScaleNodeGroupCmd runner.Cmd
	EksctlUtilsCmd          runner.Cmd
	EksctlExperimentalCmd   runner.Cmd
	// Keep track of created clusters, for post-tests clean-up.
	clustersToDelete []string
}

// SetRegion sets the provided region and re-generates other fields depending
// on the region field.
func (p *Params) SetRegion(region string) {
	p.Region = region
	p.GenerateCommands()
}

// GenerateKubeconfigPath generates a path in ${TEMP} based on these params' cluster name.
func (p *Params) GenerateKubeconfigPath() {
	p.KubeconfigPath = path.Join(os.TempDir(), fmt.Sprintf("%s.yaml", p.ClusterName))
}

// GenerateCommands generates eksctl commands with the various options & values
// provided to this `Params` object.
func (p *Params) GenerateCommands() {
	p.EksctlCmd = runner.NewCmd(p.EksctlPath).
		WithArgs("--region", p.Region).
		WithTimeout(30 * time.Minute)

	p.EksctlCreateCmd = p.EksctlCmd.
		WithArgs("create").
		WithTimeout(25 * time.Minute)

	p.EksctlUpdateCmd = p.EksctlCmd.
		WithArgs("update").
		WithTimeout(55 * time.Minute)

	p.EksctlUpgradeCmd = p.EksctlCmd.
		WithArgs("upgrade")

	p.EksctlGetCmd = p.EksctlCmd.
		WithArgs("get").
		WithTimeout(1 * time.Minute)

	p.EksctlDeleteCmd = p.EksctlCmd.
		WithArgs("delete").
		WithTimeout(15 * time.Minute)

	p.EksctlDeleteClusterCmd = p.EksctlDeleteCmd.
		WithArgs("cluster", "--verbose", "4")

	p.EksctlScaleNodeGroupCmd = p.EksctlCmd.
		WithArgs("scale", "nodegroup", "--verbose", "4").
		WithTimeout(5 * time.Minute)

	p.EksctlUtilsCmd = p.EksctlCmd.
		WithArgs("utils").
		WithTimeout(5 * time.Minute)

	p.EksctlExperimentalCmd = p.EksctlCmd.
		WithEnv("EKSCTL_EXPERIMENTAL=true")
}

// NewClusterName generates a new cluster name using the provided prefix, and
// adds the cluster to the list of clusters to eventually delete, once the test
// suite has run.
func (p *Params) NewClusterName(prefix string) string {
	return p.formatClusterName(prefix, names.ForCluster("", ""))
}

func (p *Params) formatClusterName(prefix string, name string) string {
	clusterName := fmt.Sprintf("it-%s-%s", prefix, name)
	p.addToDeleteList(clusterName)
	return clusterName
}

// addToDeleteList adds the provided cluster name to the list of clusters to eventually delete.
func (p *Params) addToDeleteList(clusterName string) {
	p.clustersToDelete = append(p.clustersToDelete, clusterName)
}

func (p Params) DeleteClusters() {
	if p.SkipDelete {
		return
	}
	for _, clusterName := range p.clustersToDelete {
		cmd := p.EksctlDeleteClusterCmd.WithArgs(
			"--name", clusterName,
		)
		session := cmd.Start()
		if session.ExitCode() != 1 {
			fmt.Fprintf(GinkgoWriter, "Warning: cluster %s's deletion failed", clusterName)
		}
	}
}

const (
	defaultTestDirectory     = "test_profile"
	defaultPrivateSSHKeyPath = "~/.ssh/eksctl-bot_id_rsa"
)

// NewParams creates a new Test instance from CLI args, grouping all test parameters.
func NewParams(clusterNamePrefix string) *Params {
	params := Params{clusterNamePrefix: clusterNamePrefix}

	flag.StringVar(&params.EksctlPath, "eksctl.path", "../../../eksctl", "Path to eksctl")
	flag.StringVar(&params.Region, "eksctl.region", api.DefaultRegion, "Region to use for the tests")
	flag.StringVar(&params.Version, "eksctl.version", api.DefaultVersion, "Version of Kubernetes to test")
	flag.StringVar(&params.TestDirectory, "eksctl.test.dir", defaultTestDirectory, "Test directory. Defaulted to: "+defaultTestDirectory)
	// Flags to help with the development of the integration tests
	flag.StringVar(&params.ClusterName, "eksctl.cluster", "", "Cluster name (default: generate one)")
	flag.BoolVar(&params.SkipCreate, "eksctl.skip.create", false, "Skip the creation tests. Useful for debugging the tests")
	flag.BoolVar(&params.SkipDelete, "eksctl.skip.delete", false, "Skip the cleanup after the tests have run")
	flag.StringVar(&params.KubeconfigPath, "eksctl.kubeconfig", "", "Path to kubeconfig (default: create a temporary file)")
	flag.StringVar(&params.PrivateSSHKeyPath, "eksctl.git.sshkeypath", defaultPrivateSSHKeyPath, fmt.Sprintf("Path to the SSH key to use for Git operations (default: %s)", defaultPrivateSSHKeyPath))

	// go1.13+ testing flags regression fix: https://github.com/golang/go/issues/31859
	flag.Parse()
	if params.ClusterName == "" {
		params.ClusterName = params.NewClusterName(clusterNamePrefix)
	} else {
		params.addToDeleteList(params.ClusterName)
	}
	if params.KubeconfigPath == "" {
		params.GenerateKubeconfigPath()
	}
	params.GenerateCommands()
	return &params
}
