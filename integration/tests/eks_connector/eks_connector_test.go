//go:build integration
// +build integration

//revive:disable Not changing package name
package eks_connector_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
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
			if !params.SkipCreate {
				By("creating an EKS cluster")
				cmd := params.EksctlCreateCmd.
					WithArgs(
						"cluster",
						"--name",
						params.ClusterName,
					)

				Expect(cmd).To(RunSuccessfully())
			}

			By("registering the new cluster")
			connectedClusterName := makeConnectedClusterName()

			cmd := params.EksctlRegisterCmd.WithArgs("cluster").
				WithArgs(
					"--name", connectedClusterName,
					"--provider", "OTHER",
				)

			wd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(fmt.Sprintf("registered cluster %q successfully", connectedClusterName))),
				ContainElement(ContainSubstring(fmt.Sprintf("wrote file eks-connector.yaml to %s", wd))),
				ContainElement(ContainSubstring(fmt.Sprintf("wrote file eks-connector-clusterrole.yaml to %s", wd))),
				ContainElement(ContainSubstring(fmt.Sprintf("wrote file eks-connector-console-dashboard-full-access-group.yaml to %s", wd))),
			))

			resourceFilenames := []string{"eks-connector.yaml", "eks-connector-clusterrole.yaml", "eks-connector-console-dashboard-full-access-group.yaml"}
			var resourcePaths []string
			for _, f := range resourceFilenames {
				resourcePaths = append(resourcePaths, path.Join(wd, f))
			}

			defer func() {
				for _, f := range resourcePaths {
					if err := os.Remove(f); err != nil {
						fmt.Fprintf(GinkgoWriter, "unexpected error removing file %q", f)
					}
				}
			}()

			By("applying the generated EKS Connector manifests to the EKS cluster")

			ctx := context.Background()
			rawClient := getRawClient(ctx, params.ClusterName, params.Region)
			for _, f := range resourcePaths {
				bytes, err := os.ReadFile(f)
				Expect(err).NotTo(HaveOccurred())
				Expect(rawClient.CreateOrReplace(bytes, false)).To(Succeed())
			}

			provider, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Name:   params.ClusterName,
					Region: params.Region,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			By("ensuring the registered cluster is active and visible")
			describeClusterInput := &awseks.DescribeClusterInput{
				Name: aws.String(connectedClusterName),
			}
			Eventually(func() string {
				connectedCluster, err := provider.AWSProvider.EKS().DescribeCluster(ctx, describeClusterInput)
				Expect(err).NotTo(HaveOccurred())
				return string(connectedCluster.Cluster.Status)
			}, "5m", "8s").Should(Equal("ACTIVE"))

			cmd = params.EksctlGetCmd.WithArgs("clusters", "-n", connectedClusterName)
			Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("OTHER")))

			By("ensuring `get nodegroup` fails early with a user-friendly error")
			cmd = params.EksctlGetCmd.
				WithArgs(
					"nodegroup",
					"--cluster", connectedClusterName,
				)

			session := cmd.Run()
			Expect(session.ExitCode()).NotTo(Equal(0))
			output := string(session.Err.Contents())
			Expect(output).To(ContainSubstring(fmt.Sprintf("cannot perform this operation on a non-EKS cluster; please follow the documentation for "+
				"cluster %s's Kubernetes provider", connectedClusterName)))

			By("deregistering the cluster")
			cmd = params.EksctlDeregisterCmd.WithArgs("cluster").
				WithArgs("--name", connectedClusterName)
			Expect(cmd).To(RunSuccessfully())

			_, err = provider.AWSProvider.EKS().DescribeCluster(ctx, describeClusterInput)
			Expect(err).To(HaveOccurred())
			var notFoundErr *ekstypes.ResourceNotFoundException
			Expect(errors.As(err, &notFoundErr)).To(BeTrue())
		})
	})
})

func makeConnectedClusterName() string {
	return fmt.Sprintf("connected-%s", params.ClusterName)
}

func deregisterCluster() {
	connectedClusterName := makeConnectedClusterName()
	cmd := params.EksctlDeregisterCmd.WithArgs("cluster").
		WithArgs("--name", makeConnectedClusterName())

	session := cmd.Run()
	if session.ExitCode() == 0 {
		fmt.Fprintf(GinkgoWriter, "cleaned up registered cluster %q successfully", connectedClusterName)
	} else {
		fmt.Fprintf(GinkgoWriter, "warning: failed to deregister cluster %q; this can be ignored if the test ran successfully to completion", connectedClusterName)
	}
}

var _ = AfterSuite(func() {
	if !params.SkipCreate && !params.SkipDelete {
		params.DeleteClusters()
	}
	deregisterCluster()
})

func getRawClient(ctx context.Context, clusterName, region string) *kubewrapper.RawClient {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   clusterName,
			Region: region,
		},
	}
	ctl, err := eks.New(context.TODO(), &api.ProviderConfig{Region: region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(ctx, cfg)
	Expect(err).ShouldNot(HaveOccurred())
	rawClient, err := ctl.NewRawClient(cfg)
	Expect(err).NotTo(HaveOccurred())
	return rawClient
}
