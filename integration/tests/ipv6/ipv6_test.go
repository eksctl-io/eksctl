//go:build integration
// +build integration

package ipv6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/xgfone/netaddr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("ipv6")
}

func TestIPv6(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var clusterConfig *api.ClusterConfig

var _ = BeforeSuite(func() {
	f, err := os.CreateTemp("", "kubeconfig-")
	Expect(err).NotTo(HaveOccurred())
	params.KubeconfigPath = f.Name()
	params.KubeconfigTemp = true

	clusterConfig = api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
	clusterConfig.Metadata.Version = api.LatestVersion
	clusterConfig.Metadata.Region = params.Region
	clusterConfig.KubernetesNetworkConfig.IPFamily = "iPv6"
	clusterConfig.VPC.NAT = nil
	clusterConfig.IAM.WithOIDC = api.Enabled()
	clusterConfig.Addons = []*api.Addon{
		{
			Name: api.VPCCNIAddon,
		},
		{
			Name: api.KubeProxyAddon,
		},
		{
			Name: api.CoreDNSAddon,
		},
	}
	clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name: "mng-1",
			},
		},
	}

	data, err := json.Marshal(clusterConfig)
	Expect(err).ToNot(HaveOccurred())

	cmd := params.EksctlCreateCmd.
		WithArgs(
			"cluster",
			"--config-file", "-",
			"--verbose", "4",
			"--kubeconfig", params.KubeconfigPath,
		).
		WithoutArg("--region", params.Region).
		WithStdin(bytes.NewReader(data))
	Expect(cmd).To(RunSuccessfully())
})

var _ = Describe("(Integration) [EKS IPv6 test]", func() {
	params.LogStacksEventsOnFailure()

	Context("Creating a cluster with IPv6", func() {
		clusterName := params.ClusterName

		It("should support ipv6", func() {
			By("creating a VPC that has an IPv6 CIDR")
			config := NewConfig(params.Region)
			cfnSession := cfn.NewFromConfig(config)

			var describeStackOut *cfn.DescribeStacksOutput
			describeStackOut, err := cfnSession.DescribeStacks(context.Background(), &cfn.DescribeStacksInput{
				StackName: aws.String(fmt.Sprintf("eksctl-%s-cluster", clusterName)),
			})
			Expect(err).NotTo(HaveOccurred())

			var vpcID string
			for _, output := range describeStackOut.Stacks[0].Outputs {
				if *output.OutputKey == builder.VPCResourceKey {
					vpcID = *output.OutputValue
				}
			}

			ec2 := awsec2.NewFromConfig(config)
			vpcOutput, err := ec2.DescribeVpcs(context.Background(), &awsec2.DescribeVpcsInput{
				VpcIds: []string{vpcID},
			})
			Expect(err).NotTo(HaveOccurred(), vpcOutput.ResultMetadata)
			Expect(vpcOutput.Vpcs[0].Ipv6CidrBlockAssociationSet).To(HaveLen(1))

			ctx := context.Background()
			By("ensuring the K8s cluster has IPv6 enabled")
			var clientSet kubernetes.Interface
			ctl, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, clusterConfig)
			Expect(err).NotTo(HaveOccurred())
			err = ctl.RefreshClusterStatus(ctx, clusterConfig)
			Expect(err).ShouldNot(HaveOccurred())
			clientSet, err = ctl.NewStdClientSet(clusterConfig)
			Expect(err).ShouldNot(HaveOccurred())

			svcName := "ipv6-service"
			_, err = clientSet.CoreV1().Services("default").Create(context.TODO(), &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: svcName,
				},
				Spec: corev1.ServiceSpec{
					IPFamilies: []corev1.IPFamily{corev1.IPv6Protocol},
					Selector:   map[string]string{"app": "ipv6"},
					Ports: []corev1.ServicePort{
						{
							Protocol: corev1.ProtocolTCP,
							Port:     80,
						},
					},
				},
			}, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func() int {
				svc, err := clientSet.CoreV1().Services("default").Get(context.TODO(), svcName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				svcIP, err := netaddr.NewIPAddress(svc.Spec.ClusterIP)
				if err != nil {
					return 0
				}
				return svcIP.Version()
			}, 5*time.Second, time.Minute).Should(Equal(6))

			By("ensuring workloads can run successfully")
			test, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())
			d := test.CreateDeploymentFromFile(test.Namespace, "../../data/podinfo.yaml")
			test.WaitForDeploymentReady(d, 10*time.Minute)

			pods := test.ListPodsFromDeployment(d)
			Expect(len(pods.Items)).To(Equal(2))

			// For each pod of the Deployment, check we receive a sensible response to a
			// GET request on /version.
			for _, pod := range pods.Items {
				Expect(pod.Namespace).To(Equal(test.Namespace))

				req := test.PodProxyGet(&pod, "", "/version")
				fmt.Fprintf(GinkgoWriter, "url = %#v", req.URL())

				var js interface{}
				test.PodProxyGetJSON(&pod, "", "/version", &js)

				Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.5.1"))
			}
		})
	})
})

var _ = AfterSuite(func() {
	cmd := params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
		"--verbose", "2",
	)
	Expect(cmd).To(RunSuccessfully())

	if params.KubeconfigTemp {
		os.Remove(params.KubeconfigPath)
	}
})
