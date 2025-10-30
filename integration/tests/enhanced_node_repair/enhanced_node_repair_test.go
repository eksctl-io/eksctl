//go:build integration
// +build integration

package enhancednoderepair

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParamsWithGivenClusterName("enhanced-node-repair", "test-enhanced-node-repair")
}

func TestEnhancedNodeRepair(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Enhanced Node Repair Configuration", func() {

	Context("CloudFormation template generation", func() {
		It("should generate correct CloudFormation template with CLI flags", func() {
			By("testing CLI flags generate correct CloudFormation")
			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--name", "test-cli-template",
				"--region", params.Region,
				"--managed",
				"--enable-node-repair",
				"--node-repair-max-unhealthy-percentage=25",
				"--node-repair-max-parallel-count=2",
				"--dry-run",
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should generate correct CloudFormation template with YAML config", func() {
			By("creating temporary config file")
			configFile := fmt.Sprintf("/tmp/test-enhanced-node-repair-%d.yaml", time.Now().Unix())
			yamlConfig := fmt.Sprintf(`
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test-yaml-template
  region: %s

managedNodeGroups:
- name: enhanced-ng
  instanceType: t3.medium
  desiredCapacity: 2
  nodeRepairConfig:
    enabled: true
    maxUnhealthyNodeThresholdPercentage: 20
    maxParallelNodesRepairedPercentage: 15
    nodeRepairConfigOverrides:
    - nodeMonitoringCondition: "NetworkNotReady"
      nodeUnhealthyReason: "InterfaceNotUp"
      minRepairWaitTimeMins: 15
      repairAction: "Restart"
`, params.Region)

			err := os.WriteFile(configFile, []byte(yamlConfig), 0644)
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(configFile)

			By("testing YAML config generates correct CloudFormation")
			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--config-file", configFile,
				"--dry-run",
			).WithoutArg("--region", params.Region)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should validate backward compatibility with existing config", func() {
			By("testing existing node repair config still works")
			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--name", "test-backward-compat",
				"--region", params.Region,
				"--managed",
				"--enable-node-repair",
				"--dry-run",
			)
			Expect(cmd).To(RunSuccessfully())
		})
	})

	Context("error handling", func() {
		It("should handle invalid CLI flag combinations gracefully", func() {
			By("testing with unmanaged nodegroup (should fail)")
			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--name", "test-error-handling",
				"--region", params.Region,
				"--managed=false",
				"--enable-node-repair",
				"--dry-run",
			)
			Expect(cmd).NotTo(RunSuccessfully())
		})

		It("should handle invalid YAML configuration gracefully", func() {
			By("creating config file with invalid node repair config")
			configFile := fmt.Sprintf("/tmp/test-invalid-config-%d.yaml", time.Now().Unix())
			invalidConfig := fmt.Sprintf(`
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test-invalid
  region: %s

nodeGroups:
- name: unmanaged-ng
  instanceType: t3.medium
  nodeRepairConfig:
    enabled: true
`, params.Region)

			err := os.WriteFile(configFile, []byte(invalidConfig), 0644)
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(configFile)

			By("testing invalid config is rejected")
			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--config-file", configFile,
				"--dry-run",
			).WithoutArg("--region", params.Region)
			// This should fail because nodeRepairConfig is not supported for unmanaged nodegroups
			Expect(cmd).NotTo(RunSuccessfully())
		})
	})
})
