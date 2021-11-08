//go:build integration
// +build integration

package cluster_dns

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/utilities/kube"

	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("cluster-dns")
}

func TestClusterDNS(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [Cluster DNS test]", func() {

	Context("Cluster with non-default ServiceIPv4CIDR", func() {
		BeforeSuite(func() {
			if params.SkipCreate {
				return
			}
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = params.ClusterName
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.Metadata.Version = params.Version
			clusterConfig.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
				ServiceIPv4CIDR: "172.16.0.0/12",
			}
			clusterConfig.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "dns",
						ScalingConfig: &api.ScalingConfig{
							DesiredCapacity: aws.Int(1),
						},
					},
				},
			}

			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--config-file", "-",
				"--kubeconfig", params.KubeconfigPath,
				"--verbose", "4",
			).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))

			Expect(cmd).To(RunSuccessfully())
		})

		It("cluster DNS should work", func() {
			test, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-dns.yaml")
			test.WaitForDaemonSetReady(d, 2*time.Minute)
			ds, err := test.GetDaemonSet(test.Namespace, d.Name)
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
		})

		AfterSuite(func() {
			if params.SkipDelete {
				return
			}
			cmd := params.EksctlDeleteCmd.WithArgs(
				"cluster", params.ClusterName,
				"--verbose", "2",
			)
			Expect(cmd).To(RunSuccessfully())
		})
	})

})
