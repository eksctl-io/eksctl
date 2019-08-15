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

var _ = Describe("(Integration) Create & Delete before Active", func() {
	const (
		initNG = "ng-0"
		testNG = "ng-1"
	)

	Describe("when creating a cluster with 1 node", func() {
		var clName string
		Context("for deleting a creating cluster", func() {
			It("should run eksctl and not wait for it to finish", func() {

				fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", kubeconfigPath)

				if clName == "" {
					clName = cmdutils.ClusterName("", "") + "-delb4active"
				}

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
		})

		Context("when deleting the (creating) cluster", func() {

			It("should not return an error", func() {

				eksctlSuccess("delete", "cluster",
					"--verbose", "4",
					"--name", clName,
					"--region", region,
					"--wait",
				)
			})

			It("and should have deleted the EKS cluster and both CloudFormation stacks", func() {

				awsSession := aws.NewSession(region)

				Expect(awsSession).ToNot(HaveExistingCluster(clName, awseks.ClusterStatusActive, version))

				Expect(awsSession).ToNot(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", clName)))
				Expect(awsSession).ToNot(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-ng-%d", clName, 0)))
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
})
