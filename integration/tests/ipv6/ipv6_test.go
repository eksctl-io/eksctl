//go:build integration
// +build integration

package ipv6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xgfone/netaddr"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("IPv6")
}

func TestIPv6(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [EKS IPv6 test]", func() {
	var (
		clusterConfig *api.ClusterConfig
	)

	Context("Creating a cluster with IPv6", func() {
		clusterName := params.NewClusterName("ipv6")

		BeforeSuite(func() {
			clusterConfig = api.NewClusterConfig()
			clusterConfig.Metadata.Name = clusterName
			clusterConfig.Metadata.Version = "1.21"
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.VPC.IPFamily = aws.String("IPv6")
			clusterConfig.VPC.NAT = nil
			clusterConfig.IAM.WithOIDC = api.Enabled()
			clusterConfig.Addons = []*api.Addon{
				{
					Name: "vpc-cni",
				},
				{
					Name: "kube-proxy",
				},
				{
					Name: "coredns",
				},
			}

			data, err := json.Marshal(clusterConfig)
			Expect(err).ToNot(HaveOccurred())

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())
		})

		AfterSuite(func() {
			cmd := params.EksctlDeleteCmd.WithArgs(
				"cluster", clusterName,
				"--verbose", "2",
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should support ipv6", func() {
			By("creating a VPC that has an IPv6 CIDR")
			awsSession := NewSession(params.Region)
			cfnSession := cfn.New(awsSession)

			var describeStackOut *cfn.DescribeStacksOutput
			describeStackOut, err := cfnSession.DescribeStacks(&cfn.DescribeStacksInput{
				StackName: aws.String(fmt.Sprintf("eksctl-%s-cluster", clusterName)),
			})
			Expect(err).NotTo(HaveOccurred())

			var vpcID string
			for _, output := range describeStackOut.Stacks[0].Outputs {
				if *output.OutputKey == builder.VPCResourceKey {
					vpcID = *output.OutputValue
				}
			}

			ec2 := awsec2.New(awsSession)
			vpcOutput, err := ec2.DescribeVpcs(&awsec2.DescribeVpcsInput{
				VpcIds: aws.StringSlice([]string{vpcID}),
			})
			Expect(err).NotTo(HaveOccurred(), vpcOutput.GoString())
			Expect(vpcOutput.Vpcs[0].Ipv6CidrBlockAssociationSet).To(HaveLen(1))

			// TODO: get rid of this once CF bug is fixed
			By("setting AssignIpv6AddressOnCreation to true for each public subnet")
			var publicSubnets string
			for _, output := range describeStackOut.Stacks[0].Outputs {
				if *output.OutputKey == builder.PublicSubnetsOutputKey {
					publicSubnets = *output.OutputValue
				}
			}

			subnetsOutput, err := ec2.DescribeSubnets(&awsec2.DescribeSubnetsInput{
				SubnetIds: aws.StringSlice(strings.Split(publicSubnets, ",")),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(subnetsOutput.Subnets)).To(BeNumerically(">", 0))
			for _, s := range subnetsOutput.Subnets {
				Expect(s.AssignIpv6AddressOnCreation).NotTo(BeNil())
				Expect(*s.AssignIpv6AddressOnCreation).To(BeTrue())
			}

			By("the k8s cluster's having an IP family of IPv6")
			var clientSet *kubernetes.Clientset
			ctl, err := eks.New(&api.ProviderConfig{Region: params.Region}, clusterConfig)
			Expect(err).NotTo(HaveOccurred())
			err = ctl.RefreshClusterStatus(clusterConfig)
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
					Ports: []corev1.ServicePort{corev1.ServicePort{
						Protocol: corev1.ProtocolTCP,
						Port:     80,
					}},
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
		})
	})
})
