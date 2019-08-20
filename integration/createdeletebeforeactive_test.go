// +build integration

package integration_test

import (
	"fmt"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/testutils/aws"
	. "github.com/weaveworks/eksctl/pkg/testutils/matchers"
)

const (
	POLLINTERVAL = 15   //seconds
	TIMEOUT      = 1200 //seconds = 20 minutes
)

var _ = Describe("(Integration) Create & Delete before Active", func() {
	const initNG = "ng-0"
	var clName string

	// initialize clName (and possibly clusterName) for this test suite
	if clusterName == "" {
		clusterName = cmdutils.ClusterName("", "")
	}
	if clName == "" {
		clName = clusterName + "-delb4active"
	}

	Context("when creating a new cluster", func() {
		It("should not return an error", func() {
			eksctlStart("create", "cluster",
				"--verbose", "4",
				"--name", clName,
				"--tags", "alpha.eksctl.io/description=eksctl delete before active test",
				"--nodegroup-name", initNG,
				"--node-labels", "ng-name="+initNG,
				"--node-type", "t2.medium",
				"--nodes", "1",
				"--region", region,
				"--version", version,
			)
		})

		It("should eventually show up as creating", func() {
			awsSession := aws.NewSession(region)
			Eventually(awsSession, TIMEOUT, POLLINTERVAL).Should(
				HaveExistingCluster(clName, awseks.ClusterStatusCreating, version))
		})
	})

	Context("when deleting the cluster in process of being created", func() {
		It("deleting cluster should have a zero exitcode", func() {
			eksctlSuccess("delete", "cluster",
				"--verbose", "4",
				"--name", clName,
				"--region", region,
			)
		})
	})

	Context("after the delete of the cluster in progress has been initiated", func() {
		It("should eventually delete the EKS cluster and both CloudFormation stacks", func() {
			awsSession := aws.NewSession(region)
			Eventually(awsSession, TIMEOUT, POLLINTERVAL).ShouldNot(
				HaveExistingCluster(clName, awseks.ClusterStatusActive, version))
			Eventually(awsSession, TIMEOUT, POLLINTERVAL).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", clName)))
			Eventually(awsSession, TIMEOUT, POLLINTERVAL).ShouldNot(
				HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", clName, initNG)))
		})
	})

	Context("when trying to delete the cluster again", func() {
		It("should return an a non-zero exit code", func() {
			eksctlFail("delete", "cluster",
				"--verbose", "4",
				"--name", clName,
				"--region", region,
			)
		})
	})
})
