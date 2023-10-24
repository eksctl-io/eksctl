//go:build integration
// +build integration

package tests

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"

	"github.com/weaveworks/eksctl/integration/runner"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/names"
)

const (
	owner = "eksctl-bot"
)

// Params groups all test parameters.
type Params struct {
	userID     string
	EksctlPath string
	Region     string
	Version    string
	// Flags to help with the development of the integration tests
	clusterNamePrefix        string
	ClusterName              string
	SkipCreate               bool
	SkipDelete               bool
	KubeconfigPath           string
	GitopsOwner              string
	IsGitopsOwnerPersonal    bool
	KubeconfigTemp           bool
	TestDirectory            string
	EksctlCmd                runner.Cmd
	EksctlCreateCmd          runner.Cmd
	EksctlCreateNodegroupCmd runner.Cmd
	EksctlUpgradeCmd         runner.Cmd
	EksctlUpdateCmd          runner.Cmd
	EksctlGetCmd             runner.Cmd
	EksctlSetLabelsCmd       runner.Cmd
	EksctlUnsetLabelsCmd     runner.Cmd
	EksctlDeleteCmd          runner.Cmd
	EksctlDeleteClusterCmd   runner.Cmd
	EksctlDrainNodeGroupCmd  runner.Cmd
	EksctlScaleNodeGroupCmd  runner.Cmd
	EksctlUtilsCmd           runner.Cmd
	EksctlEnableCmd          runner.Cmd
	EksctlAnywhereCmd        runner.Cmd
	EksctlHelpCmd            runner.Cmd
	EksctlRegisterCmd        runner.Cmd
	EksctlDeregisterCmd      runner.Cmd
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

	p.EksctlHelpCmd = runner.NewCmd(p.EksctlPath).
		WithArgs("--help").
		WithTimeout(30 * time.Minute)

	p.EksctlCreateCmd = p.EksctlCmd.
		WithArgs("create").
		WithTimeout(90 * time.Minute)

	p.EksctlUpgradeCmd = p.EksctlCmd.
		WithArgs("upgrade").
		WithTimeout(90 * time.Minute)

	p.EksctlUpdateCmd = p.EksctlCmd.
		WithArgs("update").
		WithTimeout(90 * time.Minute)

	p.EksctlGetCmd = p.EksctlCmd.
		WithArgs("get").
		WithTimeout(2 * time.Minute)

	p.EksctlSetLabelsCmd = p.EksctlCmd.
		WithArgs("set", "labels").
		// increased timeout reason: set label updates the cloudformation stack to propagate labels
		// down to the nodes. That can take a while...
		WithTimeout(10 * time.Minute)

	p.EksctlUnsetLabelsCmd = p.EksctlCmd.
		WithArgs("unset", "labels").
		// increased timeout reason: unset label updates the cloudformation stack to propagate labels
		// down to the nodes. That can take a while...
		WithTimeout(10 * time.Minute)

	p.EksctlDeleteCmd = p.EksctlCmd.
		WithArgs("delete").
		WithTimeout(30 * time.Minute)

	p.EksctlDeleteClusterCmd = p.EksctlDeleteCmd.
		WithArgs("cluster", "--verbose", "4").
		WithTimeout(40 * time.Minute)

	p.EksctlDrainNodeGroupCmd = p.EksctlCmd.
		WithArgs("drain", "nodegroup", "--verbose", "4").
		WithTimeout(10 * time.Minute)

	p.EksctlScaleNodeGroupCmd = p.EksctlCmd.
		WithArgs("scale", "nodegroup", "--verbose", "4").
		WithTimeout(5 * time.Minute)

	p.EksctlUtilsCmd = p.EksctlCmd.
		WithArgs("utils").
		WithTimeout(5 * time.Minute)

	p.EksctlEnableCmd = runner.NewCmd(p.EksctlPath).
		WithArgs("enable").
		WithTimeout(10 * time.Minute)

	p.EksctlRegisterCmd = runner.NewCmd(p.EksctlPath).
		WithArgs("register").
		WithTimeout(2 * time.Minute)

	p.EksctlDeregisterCmd = runner.NewCmd(p.EksctlPath).
		WithArgs("deregister").
		WithTimeout(1 * time.Minute)

	p.EksctlCreateNodegroupCmd = runner.NewCmd(p.EksctlPath).
		WithArgs("create", "nodegroup").
		WithTimeout(40 * time.Minute)

	p.EksctlAnywhereCmd = runner.NewCmd(p.EksctlPath).
		WithArgs("anywhere").
		WithTimeout(1 * time.Minute)
}

// NewClusterName generates a new cluster name using the provided prefix, and
// adds the cluster to the list of clusters to eventually delete, once the test
// suite has run.
func (p *Params) NewClusterName(prefix string) string {
	return p.newClusterName(prefix, "")
}

func (p *Params) newClusterName(prefix string, givenName string) string {
	return p.formatClusterName(prefix, names.ForCluster(givenName, ""))
}

func (p *Params) formatClusterName(prefix string, name string) string {
	clusterName := fmt.Sprintf("it-%s-%s", prefix, name)
	if p.userID != "" {
		clusterName = p.userID + "-" + clusterName
	}
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
		session := cmd.Run()
		if session.ExitCode() != 0 {
			fmt.Fprintf(GinkgoWriter, "Warning: cluster %s's deletion failed", clusterName)
		}
	}
}

const (
	defaultTestDirectory = "test_profile"
)

// NewParams creates a new Test instance from CLI args, grouping all test parameters.
func NewParams(clusterNamePrefix string) *Params {
	return NewParamsWithGivenClusterName(clusterNamePrefix, "")
}

func NewParamsWithGivenClusterName(clusterNamePrefix string, givenName string) *Params {
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
	flag.StringVar(&params.GitopsOwner, "eksctl.owner", "", "User or org name to create gitops repo under")
	flag.BoolVar(&params.IsGitopsOwnerPersonal, "eksctl.personal", false, "Whether the gitops repo is associated to an org or not")

	// go1.13+ testing flags regression fix: https://github.com/golang/go/issues/31859
	flag.Parse()

	params.attemptSettingUserID()
	if params.ClusterName == "" {
		params.ClusterName = params.newClusterName(clusterNamePrefix, givenName)
	} else {
		params.addToDeleteList(params.ClusterName)
	}
	if params.KubeconfigPath == "" {
		params.GenerateKubeconfigPath()
	}
	if params.GitopsOwner == "" {
		params.GitopsOwner = owner
	}
	params.GenerateCommands()
	return &params
}

// LogStacksEventsOnFailure lists stacks and their events in case the current suite has failed
func (p *Params) LogStacksEventsOnFailure() bool {
	return p.LogStacksEventsOnFailureForCluster(p.ClusterName)
}

func (p *Params) LogStacksEventsOnFailureForCluster(clusterName string) bool {
	return AfterEach(func() {
		if CurrentSpecReport().Failed() {
			Expect(p.EksctlUtilsCmd.WithArgs(
				"describe-stacks",
				"--cluster", clusterName,
				"--events",
				"--resource-status", "CREATE_FAILED,DELETE_FAILED,UPDATE_FAILED,UPDATE_ROLLBACK_FAILED,ROLLBACK_FAILED",
			)).To(RunSuccessfully())
		}
	})
}

// attemptSettingUserID will attempt to fetch the first 4 characters of the userID running the
// integration tests to be used as a prefix for cluster names. It fails silently if an error occurs
func (p *Params) attemptSettingUserID() {
	clusterProvider, err := eks.New(context.TODO(), &api.ProviderConfig{Region: p.Region}, &api.ClusterConfig{
		TypeMeta: api.ClusterConfigTypeMeta(),
		Metadata: &api.ClusterMeta{
			Version: p.Version,
			Name:    p.ClusterName,
			Region:  p.Region,
		},
	})

	msg := "Warning: failed to get UserID, proceeding without configuring cluster name prefix. Error: %v"
	if err != nil {
		fmt.Fprintf(GinkgoWriter, msg, err)
		return
	}

	out, err := clusterProvider.AWSProvider.STS().GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Fprintf(GinkgoWriter, msg, err)
		return
	}

	if out.UserId == nil {
		fmt.Fprintf(GinkgoWriter, msg, fmt.Errorf("UserId value not set in AWS response"))
		return
	}

	//example value: ABCD1234:my-name@email-provider
	parts := strings.Split(*out.UserId, ":")
	if len(parts) == 2 {
		parts = strings.Split(parts[1], "@")
		if len(parts) == 2 {
			p.userID = parts[0][:2]
			return
		}
	}
	fmt.Fprintf(GinkgoWriter, msg, fmt.Errorf("failed to split UserID %q into parts", *out.UserId))
}
