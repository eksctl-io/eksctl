// +build integration

package integration_test

import (
	"fmt"
	"time"

	"github.com/dlespiau/kube-test-harness"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/pkg/utils/names"
)

var _ = Describe("(Integration) Fargate", func() {

	// Fargate is not supported in us-west-2 yet
	const region = "ap-northeast-1"

	deleteCluster := func(clusterName string) {
		cmd := eksctlDeleteCmd.WithArgs(
			"cluster", clusterName,
			"--verbose", "4",
			"--region", region,
		)
		Expect(cmd).To(RunSuccessfully())
	}

	type fargateTest struct {
		clusterName string
		kubeTest    *harness.Test
	}

	setup := func(ft *fargateTest, createArgs ...string) {
		ft.clusterName = "fargate-" + names.ForCluster("", "")
		args := []string{
			"cluster",
			"--name", ft.clusterName,
			"--verbose", "4",
			"--region", region,
			"--kubeconfig", kubeconfigPath,
		}

		args = append(args, createArgs...)
		cmd := eksctlCreateCmd.WithArgs(args...)
		Expect(cmd).To(RunSuccessfully())

		var err error
		ft.kubeTest, err = newKubeTest()
		Expect(err).ToNot(HaveOccurred())
	}

	assertFargateDefaultProfile := func(clusterName string, kubeTest *harness.Test) {
		By("having a default Fargate profile")
		cmd := eksctlGetCmd.WithArgs(
			"fargateprofile",
			"--cluster", clusterName,
			"--verbose", "4",
			"--region", region,
		)
		Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("fp-default")))

		By("scheduling pods matching the default profile onto Fargate")
		d := kubeTest.CreateDeploymentFromFile("default", "podinfo.yaml")
		kubeTest.WaitForDeploymentReady(d, 5*time.Minute)

		pods := kubeTest.ListPodsFromDeployment(d)
		Expect(len(pods.Items)).To(Equal(2))
		for _, pod := range pods.Items {
			Expect(pod.Spec.NodeName).To(HavePrefix("fargate-"))
		}
		cmd = eksctlDeleteCmd.WithArgs(
			"fargateprofile",
			"--cluster", clusterName,
			"--name", "fp-default",
			"--region", region,
			"--wait",
			"--verbose", "4",
		)
		Expect(cmd).To(RunSuccessfully())
	}

	assertFargateNewProfile := func(clusterName string, kubeTest *harness.Test) {
		By("creating a new Fargate profile")
		profileName := "profile-1"
		cmd := eksctlCreateCmd.WithArgs(
			"fargateprofile",
			"--cluster", clusterName,
			"--name", profileName,
			"--namespace", kubeTest.Namespace,
			"--labels", "run-on=fargate",
			"--verbose", "4",
			"--region", region,
		)
		Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(profileName)))

		By("scheduling pods matching the selector onto Fargate")
		d := kubeTest.LoadDeployment("podinfo.yaml")
		d.Spec.Template.Labels["run-on"] = "fargate"

		kubeTest.CreateDeployment(kubeTest.Namespace, d)
		kubeTest.WaitForDeploymentReady(d, 5*time.Minute)
		pods := kubeTest.ListPodsFromDeployment(d)
		Expect(len(pods.Items)).To(Equal(2))
		for _, pod := range pods.Items {
			Expect(pod.Spec.NodeName).To(HavePrefix("fargate-"))
		}

		By(fmt.Sprintf("deleting Fargate profile: %q", profileName))
		cmd = eksctlDeleteCmd.WithArgs(
			"fargateprofile",
			"--cluster", clusterName,
			"--name", profileName,
			"--wait",
			"--region", region,
			"--verbose", "4",
		)
		Expect(cmd).To(RunSuccessfully())
	}

	Context("Creating a cluster with --fargate", func() {
		ft := &fargateTest{}

		BeforeEach(func() {
			setup(ft, "--fargate")
		})

		It("should support Fargate", func() {
			assertFargateDefaultProfile(ft.clusterName, ft.kubeTest)
			assertFargateNewProfile(ft.clusterName, ft.kubeTest)
		})

		AfterEach(func() {
			deleteCluster(ft.clusterName)
		})
	})

	Context("Creating a cluster with --fargate and --managed", func() {
		ft := &fargateTest{}

		BeforeEach(func() {
			setup(ft, "--fargate", "--managed")
		})

		It("should support Fargate", func() {
			assertFargateDefaultProfile(ft.clusterName, ft.kubeTest)
			assertFargateNewProfile(ft.clusterName, ft.kubeTest)
		})

		AfterEach(func() {
			deleteCluster(ft.clusterName)
		})
	})

	Context("Creating a cluster without --fargate", func() {
		ft := &fargateTest{}

		BeforeEach(func() {
			setup(ft)
		})

		It("should allow creation of new Fargate profiles", func() {
			assertFargateNewProfile(ft.clusterName, ft.kubeTest)
		})

		AfterEach(func() {
			deleteCluster(ft.clusterName)
		})
	})

})
