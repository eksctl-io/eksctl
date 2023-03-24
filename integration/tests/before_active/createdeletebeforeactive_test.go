//go:build integration
// +build integration

//revive:disable Not changing package name
package before_active

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params
var deleteAfterSuite bool

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("b4active")
	deleteAfterSuite = true
}

func TestBeforeActive(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const (
	pollInterval   = 15   //seconds
	timeOutSeconds = 1200 // 20 minutes
)

var _ = BeforeSuite(func() {
	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--verbose", "2",
		"--name", params.ClusterName,
		"--tags", "alpha.eksctl.io/description=eksctl delete before active test",
		"--without-nodegroup",
		"--version", params.Version,
	)
	cmd.Start()
	cfg := NewConfig(params.Region)
	Eventually(cfg, timeOutSeconds, pollInterval).Should(
		HaveExistingCluster(params.ClusterName, string(types.ClusterStatusCreating), params.Version))
})

var _ = Describe("(Integration) Create & Delete before Active", func() {
	const initNG = "ng-0"
	params.LogStacksEventsOnFailure()

	Context("when deleting the cluster in process of being created", func() {
		It("deleting cluster should have a zero exitcode", func() {
			cmd := params.EksctlDeleteClusterCmd.WithArgs(
				"--name", params.ClusterName,
			)
			Expect(cmd).To(RunSuccessfully())
			deleteAfterSuite = false
		})
	})

	Context("after the delete of the cluster in progress has been initiated", func() {
		It("should eventually delete the EKS cluster and both CloudFormation stacks", func() {
			config := NewConfig(params.Region)
			Eventually(config, timeOutSeconds, pollInterval).ShouldNot(
				HaveExistingCluster(params.ClusterName, string(types.ClusterStatusActive), params.Version))
			Eventually(config, timeOutSeconds, pollInterval).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Eventually(config, timeOutSeconds, pollInterval).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initNG)))
		})
	})

	Context("when trying to delete the cluster again", func() {
		It("should return an a non-zero exit code", func() {
			cmd := params.EksctlDeleteClusterCmd.WithArgs(
				"--name", params.ClusterName,
			)
			Expect(cmd).NotTo(RunSuccessfully())
		})
	})
})

var _ = AfterSuite(func() {
	if deleteAfterSuite {
		params.DeleteClusters()
	}
})
