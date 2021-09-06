//go:build integration
// +build integration

package eks_connector_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("eks-connector")
}

func TestEKSConnector(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [EKS Connector test]", func() {
	Describe("EKS Connector", func() {
		It("should register and deregister EKS clusters", func() {
			By("creating an EKS cluster")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--name",
					params.ClusterName,
				)

			Expect(cmd).To(RunSuccessfully())

			By("registering the new cluster")
			connectedClusterName := fmt.Sprintf("connected-%s", params.ClusterName)
			cmd = params.EksctlRegisterCmd.WithArgs("cluster").
				WithArgs(
					"--name", connectedClusterName,
					"--provider", "OTHER",
				)

			wd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(fmt.Sprintf("registered cluster %q successfully", connectedClusterName))),
				ContainElement(Equal(fmt.Sprintf("wrote file eks-connector.yaml to %s", wd))),
				ContainElement(Equal(fmt.Sprintf("wrote file eks-connector-binding.yaml to %s", wd))),
			))

			By("applying the generated EKS Connector manifests to the EKS cluster")
			rawClient := getRawClient(params.ClusterName, params.Region)
			bytes, err := ioutil.ReadFile(path.Join(wd, "eks-connector.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(rawClient.CreateOrReplace(bytes, false)).To(Succeed())

			provider, err := eks.New(&api.ProviderConfig{Region: params.Region}, &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Name:   params.ClusterName,
					Region: params.Region,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			By("ensuring the registered cluster is active and visible")
			describeClusterInput := &awseks.DescribeClusterInput{
				Name: aws.String(connectedClusterName),
			}
			Eventually(func() string {
				connectedCluster, err := provider.Provider.EKS().DescribeCluster(describeClusterInput)
				Expect(err).ToNot(HaveOccurred())
				return *connectedCluster.Cluster.Status
			}, "5m", "8s").Should(Equal("ACTIVE"))

			cmd = params.EksctlGetCmd.WithArgs("clusters", "-n", params.ClusterName)
			Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("OTHER")))

			By("deregistering the cluster")
			cmd = params.EksctlDeregisterCmd.WithArgs("cluster").
				WithArgs("--name", connectedClusterName)
			Expect(cmd).To(RunSuccessfully())

			_, err = provider.Provider.EKS().DescribeCluster(describeClusterInput)
			Expect(err).To(HaveOccurred())
			awsErr, ok := err.(awserr.Error)
			Expect(ok && awsErr.Code() == awseks.ErrCodeResourceNotFoundException).To(BeTrue())
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})

func getRawClient(clusterName, region string) *kubewrapper.RawClient {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   clusterName,
			Region: region,
		},
	}
	ctl, err := eks.New(&api.ProviderConfig{Region: region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	rawClient, err := ctl.NewRawClient(cfg)
	Expect(err).ToNot(HaveOccurred())
	return rawClient
}
