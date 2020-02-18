// +build integration

package before_active

import (
	"fmt"
	"testing"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("b4active")
}

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const (
	pollInterval   = 15   //seconds
	timeOutSeconds = 1200 // 20 minutes
)

var _ = Describe("(Integration) Create & Delete before Active", func() {
	const initNG = "ng-0"

	Context("when creating a new cluster", func() {
		It("should not return an error", func() {
			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--verbose", "2",
				"--name", params.ClusterName,
				"--tags", "alpha.eksctl.io/description=eksctl delete before active test",
				"--without-nodegroup",
				"--version", params.Version,
			)
			cmd.Start()
			awsSession := NewSession(params.Region)
			Eventually(awsSession, timeOutSeconds, pollInterval).Should(
				HaveExistingCluster(params.ClusterName, awseks.ClusterStatusCreating, params.Version))
		})
	})

	Context("when deleting the cluster in process of being created", func() {
		It("deleting cluster should have a zero exitcode", func() {
			cmd := params.EksctlDeleteClusterCmd.WithArgs(
				"--name", params.ClusterName,
			)
			Expect(cmd).To(RunSuccessfully())
		})
	})

	Context("after the delete of the cluster in progress has been initiated", func() {
		It("should eventually delete the EKS cluster and both CloudFormation stacks", func() {
			awsSession := NewSession(params.Region)
			Eventually(awsSession, timeOutSeconds, pollInterval).ShouldNot(
				HaveExistingCluster(params.ClusterName, awseks.ClusterStatusActive, params.Version))
			Eventually(awsSession, timeOutSeconds, pollInterval).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Eventually(awsSession, timeOutSeconds, pollInterval).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initNG)))
		})
	})

	Context("when trying to delete the cluster again", func() {
		It("should return an a non-zero exit code", func() {
			cmd := params.EksctlDeleteClusterCmd.WithArgs(
				"--name", params.ClusterName,
			)
			Expect(cmd).ToNot(RunSuccessfully())
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
