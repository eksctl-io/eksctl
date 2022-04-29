//go:build integration
// +build integration

//revive:disable Not changing package name
package instance_selector

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	// No cleanup required for dry-run clusters
	params = tests.NewParams("instance-selector")
	if err := api.Register(); err != nil {
		panic(errors.Wrap(err, "unexpected error registering API scheme"))
	}
}

func TestInstanceSelector(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [Instance Selector test]", func() {

	DescribeTable("Instance Selector", func(assertionFunc func(instanceTypes []string), instanceSelectorArgs ...string) {
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--dry-run",
				"--name", params.ClusterName,
			).WithArgs(instanceSelectorArgs...)

		session := cmd.Run()
		Expect(session.ExitCode()).To(Equal(0))

		output := session.Buffer().Contents()
		clusterConfig, err := eks.ParseConfig(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterConfig.ManagedNodeGroups).To(HaveLen(1))
		if assertionFunc != nil {
			assertionFunc(clusterConfig.ManagedNodeGroups[0].InstanceTypes)
		}
	},
		Entry("non-GPU instances", func(instanceTypes []string) {
			for _, instanceType := range instanceTypes {
				Expect(instanceType).NotTo(Satisfy(instanceutils.IsGPUInstanceType))
			}
		}, "--instance-selector-vcpus=8",
			"--instance-selector-memory=32",
			"--instance-selector-gpus=0",
		),
		Entry("with vCPUs and memory", nil,
			"--instance-selector-vcpus=8",
			"--instance-selector-memory=32",
		),
	)

})
