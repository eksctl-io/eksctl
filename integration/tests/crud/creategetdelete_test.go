//go:build integration
// +build integration

package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/weaveworks/eksctl/integration/matchers"
	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/iam"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

var (
	params        *tests.Params
	extraSubnetID string
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	if err := api.Register(); err != nil {
		panic(errors.Wrap(err, "unexpected error registering API scheme"))
	}
	params = tests.NewParamsWithGivenClusterName("crud", "test")
}

func TestCRUD(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const (
	deleteNg        = "ng-delete"
	taintsNg1       = "ng-taints-1"
	taintsNg2       = "ng-taints-2"
	scaleSingleNg   = "ng-scale-single"
	scaleMultipleNg = "ng-scale-multiple"

	scaleMultipleMng       = "mng-scale-multiple"
	GPUMng                 = "mng-gpu"
	drainMng               = "mng-drain"
	newSubnetCLIMng        = "mng-new-subnet-CLI"
	newSubnetConfigFileMng = "mng-new-subnet-config-file"
)

func makeClusterConfig() *api.ClusterConfig {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
	clusterConfig.Metadata.Region = params.Region
	clusterConfig.Metadata.Version = params.Version
	return clusterConfig
}

var _ = SynchronizedBeforeSuite(func() []byte {
	params.KubeconfigTemp = false
	if params.KubeconfigPath == "" {
		wd, _ := os.Getwd()
		f, _ := os.CreateTemp(wd, "kubeconfig-")
		params.KubeconfigPath = f.Name()
		params.KubeconfigTemp = true
	}

	if params.SkipCreate {
		fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", params.ClusterName)
		if !file.Exists(params.KubeconfigPath) {
			// Generate the Kubernetes configuration that eksctl create
			// would have generated otherwise:
			cmd := params.EksctlUtilsCmd.WithArgs(
				"write-kubeconfig",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		}
	}
	fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

	cfg := makeClusterConfig()
	cfg.NodeGroups = []*api.NodeGroup{
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name: deleteNg,
			},
		},
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name: scaleSingleNg,
			},
		},
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name: scaleMultipleNg,
			},
		},
	}
	cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name: drainMng,
			},
		},
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name: scaleMultipleMng,
			},
		},
	}
	cfg.AvailabilityZones = []string{"us-west-2b", "us-west-2c"}
	cfg.Metadata.Tags = map[string]string{
		"alpha.eksctl.io/description": "eksctl integration test",
	}

	Expect(params.EksctlCreateCmd.
		WithArgs(
			"cluster",
			"--config-file", "-",
			"--verbose", "4",
			"--kubeconfig", params.KubeconfigPath,
		).
		WithoutArg("--region", params.Region).
		WithStdin(clusterutils.Reader(cfg))).To(RunSuccessfully())

	// create an additional subnet to test nodegroup creation within it later on
	extraSubnetID = createAdditionalSubnet(cfg)
	return []byte(extraSubnetID)
}, func(subnetID []byte) {
	extraSubnetID = string(subnetID)
})

var _ = Describe("(Integration) Create, Get, Scale & Delete", func() {

	makeClientset := func() *kubernetes.Clientset {
		config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		clientset, err := kubernetes.NewForConfig(config)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		return clientset
	}

	Context("validating cluster setup", func() {
		It("should have created an EKS cluster and 6 CloudFormation stacks", func() {
			awsConfig := NewConfig(params.Region)
			Expect(awsConfig).To(HaveExistingCluster(params.ClusterName, string(ekstypes.ClusterStatusActive), params.Version))
			Expect(awsConfig).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(awsConfig).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, drainMng)))
			Expect(awsConfig).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, scaleMultipleMng)))
			Expect(awsConfig).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, deleteNg)))
			Expect(awsConfig).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, scaleSingleNg)))
			Expect(awsConfig).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, scaleMultipleNg)))
		})

		It("should have created a valid kubectl config file", func() {
			kubeConfig, err := clientcmd.LoadFromFile(params.KubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*kubeConfig, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(kubeConfig.CurrentContext).To(ContainSubstring("eksctl"))
			Expect(kubeConfig.CurrentContext).To(ContainSubstring(params.ClusterName))
			Expect(kubeConfig.CurrentContext).To(ContainSubstring(params.Region))
		})

		It("should successfully fetch the cluster", func() {
			AssertContainsCluster(
				params.EksctlGetCmd.WithArgs("clusters", "--all-regions"),
				GetClusterOutput{
					ClusterName:   params.ClusterName,
					Region:        params.Region,
					EksctlCreated: "True",
				},
			)
		})

		It("should successfully describe cluster's CFN stacks", func() {
			session := params.EksctlUtilsCmd.WithArgs("describe-stacks", "--cluster", params.ClusterName, "-o", "yaml").Run()
			Expect(session.ExitCode()).To(BeZero())
			var stacks []*cfntypes.Stack
			Expect(yaml.Unmarshal(session.Out.Contents(), &stacks)).To(Succeed())
			Expect(stacks).To(HaveLen(6))

			var (
				names, descriptions []string
				ngPrefix            = params.ClusterName + "-nodegroup-"
			)
			for _, s := range stacks {
				names = append(names, *s.StackName)
				descriptions = append(descriptions, *s.Description)
			}

			Expect(names).To(ContainElements(
				ContainSubstring(params.ClusterName+"-cluster"),
				ContainSubstring(ngPrefix+deleteNg),
				ContainSubstring(ngPrefix+scaleSingleNg),
				ContainSubstring(ngPrefix+scaleMultipleNg),
				ContainSubstring(ngPrefix+scaleMultipleMng),
				ContainSubstring(ngPrefix+drainMng),
			))
			Expect(descriptions).To(ContainElements(
				"EKS cluster (dedicated VPC: true, dedicated IAM: true) [created and managed by eksctl]",
				"EKS Managed Nodes (SSH access: false) [created by eksctl]",
				"EKS nodes (AMI family: AmazonLinux2, SSH access: false, private networking: false) [created and managed by eksctl]",
			))
		})
	})

	Context("creating cluster workloads", func() {
		var (
			err           error
			test          *harness.Test
			commonTimeout = 10 * time.Minute
		)

		BeforeEach(func() {
			test, err = kube.NewTest(params.KubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			test.Close()
			Eventually(func() int {
				return len(test.ListPods(test.Namespace, metav1.ListOptions{}).Items)
			}, "3m", "1s").Should(BeZero())
		})

		It("should deploy podinfo service to the cluster and access it via proxy", func() {
			d := test.CreateDeploymentFromFile(test.Namespace, "../../data/podinfo.yaml")
			test.WaitForDeploymentReady(d, commonTimeout)

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

		It("should have functional DNS", func() {
			d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-dns.yaml")
			test.WaitForDaemonSetReady(d, commonTimeout)
			ds, err := test.GetDaemonSet(test.Namespace, d.Name)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
		})

		It("should have access to HTTP(S) sites", func() {
			d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-http.yaml")
			test.WaitForDaemonSetReady(d, commonTimeout)
			ds, err := test.GetDaemonSet(test.Namespace, d.Name)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
		})
	})

	Context("configuring IAM service accounts", Ordered, func() {
		var (
			clientSet       kubernetes.Interface
			test            *harness.Test
			awsConfig       aws.Config
			oidc            *iamoidc.OpenIDConnectManager
			stackNamePrefix string
		)

		BeforeAll(func() {
			cfg := makeClusterConfig()

			ctl, err := eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, cfg)
			Expect(err).NotTo(HaveOccurred())

			err = ctl.RefreshClusterStatus(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())

			clientSet, err = ctl.NewStdClientSet(cfg)
			Expect(err).ShouldNot(HaveOccurred())

			test, err = kube.NewTest(params.KubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			stackNamePrefix = fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", params.ClusterName)
			awsConfig = NewConfig(params.Region)

			oidc, err = ctl.NewOpenIDConnectManager(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterAll(func() {
			test.Close()
			Eventually(func() int {
				return len(test.ListPods(test.Namespace, metav1.ListOptions{}).Items)
			}, "3m", "1s").Should(BeZero())
		})

		It("should have OIDC disabled by default", func() {
			exists, err := oidc.CheckProviderExists(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should successfully enable OIDC", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"associate-iam-oidc-provider",
				"--cluster", params.ClusterName,
				"--approve",
			)).To(RunSuccessfully())
		})

		It("should successfully create two iamserviceaccounts", func() {
			Expect([]Cmd{
				params.EksctlCreateCmd.WithArgs(
					"iamserviceaccount",
					"--cluster", params.ClusterName,
					"--name", "app-cache-access",
					"--namespace", "app1",
					"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
					"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
					"--approve",
				),
				params.EksctlCreateCmd.WithArgs(
					"iamserviceaccount",
					"--cluster", params.ClusterName,
					"--name", "s3-reader",
					"--namespace", test.Namespace,
					"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
					"--approve",
				),
			}).To(RunSuccessfully())
			Expect(awsConfig).To(HaveExistingStack(stackNamePrefix + test.Namespace + "-s3-reader"))
			Expect(awsConfig).To(HaveExistingStack(stackNamePrefix + "app1-app-cache-access"))

			sa, err := clientSet.CoreV1().ServiceAccounts(test.Namespace).Get(context.TODO(), "s3-reader", metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sa.Annotations).To(HaveLen(1))
			Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
			Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + params.ClusterName + ".*$"))

			sa, err = clientSet.CoreV1().ServiceAccounts("app1").Get(context.TODO(), "app-cache-access", metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sa.Annotations).To(HaveLen(1))
			Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
			Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + params.ClusterName + ".*$"))
		})

		It("should successfully update service account policy", func() {
			Expect(params.EksctlUpdateCmd.WithArgs(
				"iamserviceaccount",
				"--cluster", params.ClusterName,
				"--name", "app-cache-access",
				"--namespace", "app1",
				"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
				"--approve",
			)).To(RunSuccessfully())
		})

		It("should successfully list both iamserviceaccounts", func() {
			Expect(params.EksctlGetCmd.WithArgs(
				"iamserviceaccount",
				"--cluster", params.ClusterName,
			)).To(RunSuccessfullyWithOutputString(MatchRegexp(
				`(?m:^NAMESPACE\s+NAME\s+ROLE\sARN$)` +
					`|(?m:^app1\s+app-cache-access\s+arn:aws:iam::.*$)` +
					fmt.Sprintf(`|(?m:^%s\s+s3-reader\s+arn:aws:iam::.*$)`, test.Namespace),
			)))
		})

		It("should successfully run pods with an iamserviceaccount", func() {
			d := test.CreateDeploymentFromFile(test.Namespace, "../../data/iamserviceaccount-checker.yaml")
			test.WaitForDeploymentReady(d, 10*time.Minute)

			pods := test.ListPodsFromDeployment(d)
			Expect(len(pods.Items)).To(Equal(2))

			// For each pod of the Deployment, check we get expected environment variables
			// via a GET request on /env.
			type sessionObject struct {
				AssumedRoleUser struct {
					AssumedRoleID, Arn string
				}
				Audience, Provider, SubjectFromWebIdentityToken string
				Credentials                                     struct {
					SecretAccessKey, SessionToken, Expiration, AccessKeyID string
				}
			}

			for _, pod := range pods.Items {
				Expect(pod.Namespace).To(Equal(test.Namespace))

				so := sessionObject{}

				var js []string
				test.PodProxyGetJSON(&pod, "", "/env", &js)

				Expect(js).To(ContainElement(HavePrefix("AWS_ROLE_ARN=arn:aws:iam::")))
				Expect(js).To(ContainElement("AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/eks.amazonaws.com/serviceaccount/token"))
				Expect(js).To(ContainElement(HavePrefix("AWS_SESSION_OBJECT=")))

				for _, envVar := range js {
					if strings.HasPrefix(envVar, "AWS_SESSION_OBJECT=") {
						err := json.Unmarshal([]byte(strings.TrimPrefix(envVar, "AWS_SESSION_OBJECT=")), &so)
						Expect(err).ShouldNot(HaveOccurred())
					}
				}

				Expect(so.AssumedRoleUser.AssumedRoleID).To(HaveSuffix(":integration-test"))

				Expect(so.AssumedRoleUser.Arn).To(MatchRegexp("^arn:aws:sts::.*:assumed-role/eksctl-" + params.ClusterName + "-.*/integration-test$"))

				Expect(so.Audience).To(Equal("sts.amazonaws.com"))

				Expect(so.Provider).To(MatchRegexp("^arn:aws:iam::.*:oidc-provider/oidc.eks." + params.Region + ".amazonaws.com/id/.*$"))

				Expect(so.SubjectFromWebIdentityToken).To(Equal("system:serviceaccount:" + test.Namespace + ":s3-reader"))

				Expect(so.Credentials.SecretAccessKey).NotTo(BeEmpty())
				Expect(so.Credentials.SessionToken).NotTo(BeEmpty())
				Expect(so.Credentials.Expiration).NotTo(BeEmpty())
				Expect(so.Credentials.AccessKeyID).NotTo(BeEmpty())
			}
		})

		It("should successfully delete both iamserviceaccounts", func() {
			Expect([]Cmd{
				params.EksctlDeleteCmd.WithArgs(
					"iamserviceaccount",
					"--cluster", params.ClusterName,
					"--name", "s3-reader",
					"--namespace", test.Namespace,
					"--wait",
				),
				params.EksctlDeleteCmd.WithArgs(
					"iamserviceaccount",
					"--cluster", params.ClusterName,
					"--name", "app-cache-access",
					"--namespace", "app1",
					"--wait",
				),
			}).To(RunSuccessfully())
			Expect(awsConfig).NotTo(HaveExistingStack(stackNamePrefix + test.Namespace + "-s3-reader"))
			Expect(awsConfig).NotTo(HaveExistingStack(stackNamePrefix + "app1-app-cache-access"))
		})
	})

	Context("configuring K8s API", Serial, Ordered, func() {
		var (
			k8sAPICall func() error
		)

		BeforeAll(func() {
			cfg := makeClusterConfig()

			ctl, err := eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, cfg)
			Expect(err).NotTo(HaveOccurred())

			err = ctl.RefreshClusterStatus(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())

			clientSet, err := ctl.NewStdClientSet(cfg)
			Expect(err).ShouldNot(HaveOccurred())

			k8sAPICall = func() error {
				_, err = clientSet.CoreV1().ServiceAccounts(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
				return err
			}
		})

		It("should have public access by default", func() {
			Expect(k8sAPICall()).ShouldNot(HaveOccurred())
		})

		It("should disable public access", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"set-public-access-cidrs",
				"--cluster", params.ClusterName,
				"1.1.1.1/32,2.2.2.0/24",
				"--approve",
			)).To(RunSuccessfully())
			Expect(k8sAPICall()).Should(HaveOccurred())
		})

		It("should re-enable public access", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"set-public-access-cidrs",
				"--cluster", params.ClusterName,
				"0.0.0.0/0",
				"--approve",
			)).To(RunSuccessfully())
			Expect(k8sAPICall()).ShouldNot(HaveOccurred())
		})
	})

	Context("configuring Cloudwatch logging", Serial, Ordered, func() {
		var (
			cfg *api.ClusterConfig
			ctl *eks.ClusterProvider
			err error
		)

		BeforeEach(func() {
			cfg = makeClusterConfig()

			ctl, err = eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, cfg)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have all types disabled by default", func() {
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sets.List(enabled)).To(HaveLen(0))
			Expect(sets.List(disabled)).To(HaveLen(5))
		})

		It("should plan to enable two of the types using flags", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"update-cluster-logging",
				"--cluster", params.ClusterName,
				"--enable-types", "api,controllerManager",
			)).To(RunSuccessfully())
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sets.List(enabled)).To(HaveLen(0))
			Expect(sets.List(disabled)).To(HaveLen(5))
		})

		It("should enable two of the types using flags", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"update-cluster-logging",
				"--cluster", params.ClusterName,
				"--approve",
				"--enable-types", "api,controllerManager",
			)).To(RunSuccessfully())
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sets.List(enabled)).To(HaveLen(2))
			Expect(sets.List(disabled)).To(HaveLen(3))
			Expect(sets.List(enabled)).To(ConsistOf("api", "controllerManager"))
			Expect(sets.List(disabled)).To(ConsistOf("audit", "authenticator", "scheduler"))
		})

		It("should enable all of the types using --enable-types=all", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"update-cluster-logging",
				"--cluster", params.ClusterName,
				"--approve",
				"--enable-types", "all",
			)).To(RunSuccessfully())
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sets.List(enabled)).To(HaveLen(5))
			Expect(sets.List(disabled)).To(HaveLen(0))
		})

		It("should enable all but one type", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"update-cluster-logging",
				"--cluster", params.ClusterName,
				"--approve",
				"--enable-types", "all",
				"--disable-types", "controllerManager",
			)).To(RunSuccessfully())
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sets.List(enabled)).To(HaveLen(4))
			Expect(sets.List(disabled)).To(HaveLen(1))
			Expect(sets.List(enabled)).To(ConsistOf("api", "audit", "authenticator", "scheduler"))
			Expect(sets.List(disabled)).To(ConsistOf("controllerManager"))
		})

		It("should disable all but one type", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"update-cluster-logging",
				"--cluster", params.ClusterName,
				"--approve",
				"--disable-types", "all",
				"--enable-types", "controllerManager",
			)).To(RunSuccessfully())
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sets.List(disabled)).To(HaveLen(4))
			Expect(sets.List(enabled)).To(HaveLen(1))
			Expect(sets.List(disabled)).To(ConsistOf("api", "audit", "authenticator", "scheduler"))
			Expect(sets.List(enabled)).To(ConsistOf("controllerManager"))
		})

		It("should disable all of the types using --disable-types=all", func() {
			Expect(params.EksctlUtilsCmd.WithArgs(
				"update-cluster-logging",
				"--cluster", params.ClusterName,
				"--approve",
				"--disable-types", "all",
			)).To(RunSuccessfully())
			enabled, disabled, err := ctl.GetCurrentClusterConfigForLogging(context.Background(), cfg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sets.List(enabled)).To(HaveLen(0))
			Expect(sets.List(disabled)).To(HaveLen(5))
			Expect(disabled.HasAll(api.SupportedCloudWatchClusterLogTypes()...)).To(BeTrue())
		})
	})

	Context("configuring iam identity mappings", Serial, Ordered, func() {
		var (
			expR0, expR1, expU0 string
			role0, role1        iam.Identity
			user0               iam.Identity
			admin               = "admin"
			alice               = "alice"
		)

		BeforeAll(func() {
			roleCanonicalArn := "arn:aws:iam::123456:role/eksctl-testing-XYZ"
			role0 = iam.RoleIdentity{
				RoleARN: roleCanonicalArn,
				KubernetesIdentity: iam.KubernetesIdentity{
					KubernetesUsername: admin,
					KubernetesGroups:   []string{"system:masters", "system:nodes"},
				},
			}
			role1 = iam.RoleIdentity{
				RoleARN: roleCanonicalArn,
				KubernetesIdentity: iam.KubernetesIdentity{
					KubernetesGroups: []string{"system:something"},
				},
			}

			userCanonicalArn := "arn:aws:iam::123456:user/alice"

			user0 = iam.UserIdentity{
				UserARN: userCanonicalArn,
				KubernetesIdentity: iam.KubernetesIdentity{
					KubernetesUsername: alice,
					KubernetesGroups:   []string{"system:masters", "cryptographers"},
				},
			}

			bs, err := yaml.Marshal([]iam.Identity{role0})
			Expect(err).ShouldNot(HaveOccurred())
			expR0 = string(bs)

			bs, err = yaml.Marshal([]iam.Identity{role1})
			Expect(err).ShouldNot(HaveOccurred())
			expR1 = string(bs)

			bs, err = yaml.Marshal([]iam.Identity{user0})
			Expect(err).ShouldNot(HaveOccurred())
			expU0 = string(bs)
		})

		It("should fail to get unknown role mapping", func() {
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", "arn:aws:iam::123456:role/idontexist",
				"-o", "yaml",
			)).NotTo(RunSuccessfully())
		})

		It("should fail to get unknown user mapping", func() {
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", "arn:aws:iam::123456:user/bob",
				"-o", "yaml",
			)).NotTo(RunSuccessfully())
		})

		It("should create role mappings", func() {
			Expect(params.EksctlCreateCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role0.ARN(),
				"--username", role0.Username(),
				"--group", role0.Groups()[0],
				"--group", role0.Groups()[1],
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role0.ARN(),
				"-o", "yaml",
			)).To(RunSuccessfullyWithOutputString(MatchYAML(expR0)))
		})

		It("should create user mappings", func() {
			Expect(params.EksctlCreateCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", user0.ARN(),
				"--username", user0.Username(),
				"--group", user0.Groups()[0],
				"--group", user0.Groups()[1],
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", user0.ARN(),
				"-o", "yaml",
			)).To(RunSuccessfullyWithOutputString(MatchYAML(expU0)))
		})

		It("should create a duplicate role mapping", func() {
			Expect(params.EksctlCreateCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role0.ARN(),
				"--username", role0.Username(),
				"--group", role0.Groups()[0],
				"--group", role0.Groups()[1],
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role0.ARN(),
				"-o", "yaml",
			)).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0)))
		})

		It("should create a duplicate user mapping", func() {
			Expect(params.EksctlCreateCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", user0.ARN(),
				"--username", user0.Username(),
				"--group", user0.Groups()[0],
				"--group", user0.Groups()[1],
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", user0.ARN(),
				"-o", "yaml",
			)).To(RunSuccessfullyWithOutputString(MatchYAML(expU0 + expU0)))
		})

		It("should create a duplicate role mapping with different identity", func() {
			Expect(params.EksctlCreateCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role1.ARN(),
				"--group", role1.Groups()[0],
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role1.ARN(),
				"-o", "yaml",
			)).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0 + expR1)))
		})

		It("should delete a single role mapping (fifo)", func() {
			Expect(params.EksctlDeleteCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role1.ARN(),
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role1.ARN(),
				"-o", "yaml",
			)).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR1)))
		})

		It("should fail to delete unknown mapping", func() {
			Expect(params.EksctlDeleteCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", "arn:aws:iam::123456:role/idontexist",
			)).NotTo(RunSuccessfully())
		})

		It("should delete duplicate role mappings with --all", func() {
			Expect(params.EksctlDeleteCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role1.ARN(),
				"--all",
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", role1.ARN(),
				"-o", "yaml",
			)).NotTo(RunSuccessfully())
		})

		It("should delete duplicate user mappings with --all", func() {
			Expect(params.EksctlDeleteCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", user0.ARN(),
				"--all",
			)).To(RunSuccessfully())
			Expect(params.EksctlGetCmd.WithArgs(
				"iamidentitymapping",
				"--cluster", params.ClusterName,
				"--arn", user0.ARN(),
				"-o", "yaml",
			)).NotTo(RunSuccessfully())
		})
	})

	Context("creating nodegroups", func() {
		It("should be able to create two nodegroups with taints and maxPods", func() {
			By("creating them")
			Expect(params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/taints-max-pods.yaml"))).To(RunSuccessfully())

			By("asserting that both formats for taints are supported")
			clientset := makeClientset()
			nodeListN1 := tests.ListNodes(clientset, taintsNg1)
			nodeListN2 := tests.ListNodes(clientset, taintsNg2)

			tests.AssertNodeTaints(nodeListN1, []corev1.Taint{
				{
					Key:    "key1",
					Value:  "val1",
					Effect: "NoSchedule",
				},
				{
					Key:    "key2",
					Effect: "NoExecute",
				},
			})
			tests.AssertNodeTaints(nodeListN2, []corev1.Taint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: "NoSchedule",
				},
				{
					Key:    "key2",
					Effect: "NoExecute",
				},
			})

			By("asserting that maxPods is set correctly")
			expectedMaxPods := 123
			for _, node := range nodeListN1.Items {
				maxPods, _ := node.Status.Allocatable.Pods().AsInt64()
				Expect(maxPods).To(Equal(int64(expectedMaxPods)))
			}
		})

		It("should be able to create a new GPU nodegroup", func() {
			Expect(params.EksctlCreateCmd.WithArgs(
				"nodegroup",
				"--timeout=45m",
				"--cluster", params.ClusterName,
				"--nodes", "1",
				"--instance-types", "p2.xlarge,p3.2xlarge,p3.8xlarge,g3s.xlarge,g4ad.xlarge,g4ad.2xlarge",
				"--node-private-networking",
				"--node-zones", "us-west-2b,us-west-2c",
				GPUMng,
			)).To(RunSuccessfully())
		})

		Context("creating nodegroups within a new subnet", func() {
			var (
				vpcID      string
				subnetName string
			)
			BeforeEach(func() {
				ec2 := awsec2.NewFromConfig(NewConfig(params.Region))
				output, err := ec2.DescribeSubnets(context.Background(), &awsec2.DescribeSubnetsInput{
					SubnetIds: []string{extraSubnetID},
				})
				Expect(err).NotTo(HaveOccurred())
				vpcID = *output.Subnets[0].VpcId
				subnetName = "new-subnet"
			})

			It("should be able to create a nodegroup in a new subnet via config file", func() {
				clusterConfig := makeClusterConfig()
				clusterConfig.VPC = &api.ClusterVPC{
					Network: api.Network{
						ID: vpcID,
					},
					Subnets: &api.ClusterSubnets{
						Public: api.AZSubnetMapping{
							subnetName: api.AZSubnetSpec{
								ID: extraSubnetID,
							},
						},
					},
				}
				clusterConfig.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name: newSubnetConfigFileMng,
							ScalingConfig: &api.ScalingConfig{
								DesiredCapacity: aws.Int(1),
							},
							Subnets: []string{subnetName},
						},
					},
				}

				Expect(params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
					"--timeout", time.Hour.String(),
				).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.Reader(clusterConfig))).To(RunSuccessfully())

				Expect(params.EksctlDeleteCmd.WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"--wait",
					newSubnetConfigFileMng,
				)).To(RunSuccessfully())
			})

			It("should be able to create a nodegroup in a new subnet via CLI", func() {
				Expect(params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--timeout", time.Hour.String(),
					"--cluster", params.ClusterName,
					"--nodes", "1",
					"--node-type", "p2.xlarge",
					"--subnet-ids", extraSubnetID,
					newSubnetCLIMng,
				)).To(RunSuccessfully())

				Expect(params.EksctlDeleteCmd.WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"--wait",
					newSubnetCLIMng,
				)).To(RunSuccessfully())
			})
		})
	})

	Context("scaling nodegroup(s)", func() {
		scaleNgCmd := func(desiredCapacity string) runner.Cmd {
			return params.EksctlScaleNodeGroupCmd.WithArgs(
				"--cluster", params.ClusterName,
				"--nodes-min", desiredCapacity,
				"--nodes", desiredCapacity,
				"--nodes-max", desiredCapacity,
				"--name", scaleSingleNg,
			)
		}
		getNgCmd := func(ngName string) runner.Cmd {
			return params.EksctlGetCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--name", ngName,
				"-o", "yaml",
			)
		}

		It("should be able to scale a single nodegroup", func() {
			By("upscaling a nodegroup without --wait flag")
			Expect(scaleNgCmd("3")).To(RunSuccessfully())
			Eventually(getNgCmd(scaleSingleNg), "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("Type: unmanaged")),
				ContainElement(ContainSubstring("MaxSize: 3")),
				ContainElement(ContainSubstring("MinSize: 3")),
				ContainElement(ContainSubstring("DesiredCapacity: 3")),
				ContainElement(ContainSubstring("Status: CREATE_COMPLETE")),
			))

			By("upscaling a nodegroup with --wait flag")
			Expect(scaleNgCmd("4")).To(RunSuccessfully())
			Eventually(getNgCmd(scaleSingleNg), "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("Type: unmanaged")),
				ContainElement(ContainSubstring("MaxSize: 4")),
				ContainElement(ContainSubstring("MinSize: 4")),
				ContainElement(ContainSubstring("DesiredCapacity: 4")),
				ContainElement(ContainSubstring("Status: CREATE_COMPLETE")),
			))

			By("downscaling a nodegroup")
			Expect(scaleNgCmd("1")).To(RunSuccessfully())
			Eventually(getNgCmd(scaleSingleNg), "5m", "30s").Should(runner.RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("Type: unmanaged")),
				ContainElement(ContainSubstring("MaxSize: 1")),
				ContainElement(ContainSubstring("MinSize: 1")),
				ContainElement(ContainSubstring("DesiredCapacity: 1")),
				ContainElement(ContainSubstring("Status: CREATE_COMPLETE")),
			))
		})

		It("should be able to scale multiple nodegroups", func() {
			By("passing a config file")
			Expect(params.EksctlScaleNodeGroupCmd.WithArgs(
				"--config-file", "-",
			).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/scale-nodegroups.yaml")),
			).To(RunSuccessfully())

			Eventually(getNgCmd(scaleMultipleNg), "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("Type: unmanaged")),
				ContainElement(ContainSubstring("MaxSize: 5")),
				ContainElement(ContainSubstring("MinSize: 5")),
				ContainElement(ContainSubstring("DesiredCapacity: 5")),
				ContainElement(ContainSubstring("Status: CREATE_COMPLETE")),
			))

			Eventually(getNgCmd(scaleMultipleMng), "5m", "30s").Should(runner.RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("Type: managed")),
				ContainElement(ContainSubstring("MaxSize: 5")),
				ContainElement(ContainSubstring("MinSize: 5")),
				ContainElement(ContainSubstring("DesiredCapacity: 5")),
				ContainElement(ContainSubstring("Status: ACTIVE")),
			))
		})
	})

	Context("access entries", Serial, func() {
		hasAuthIdentity := func(nodeRoleARN string) bool {
			auth, err := authconfigmap.NewFromClientSet(makeClientset())
			Expect(err).NotTo(HaveOccurred())
			identities, err := auth.GetIdentities()
			Expect(err).NotTo(HaveOccurred())
			for _, id := range identities {
				if roleID, ok := id.(iam.RoleIdentity); ok && roleID.RoleARN == nodeRoleARN {
					return true
				}
			}
			return false
		}

		type authMode int
		const (
			authModeAccessEntry authMode = iota + 1
			authModeAWSAuthConfigMap
		)

		type ngAuthTest struct {
			ngName       string
			createNgArgs []string

			expectedAuthMode       authMode
			expectedCreateNgOutput string
		}

		DescribeTable("authorising self-managed nodegroups", Serial, func(nt ngAuthTest) {
			cmd := params.EksctlCreateNodegroupCmd.
				WithArgs(
					"--cluster", params.ClusterName,
					"--name", nt.ngName,
					"--nodes", "1",
					"--managed=false",
				).WithArgs(nt.createNgArgs...)
			session := cmd.Run()
			Expect(session.ExitCode()).To(BeZero())
			if nt.expectedCreateNgOutput != "" {
				Expect(session.Out.Contents()).To(ContainSubstring(nt.expectedCreateNgOutput))
			}

			DeferCleanup(func() {
				cmd := params.EksctlDeleteCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--name", nt.ngName,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			cmd = params.EksctlGetCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--name", nt.ngName,
				"-o", "json",
			)
			session = cmd.Run()
			Expect(session.ExitCode()).To(BeZero())
			var ngSummaries []nodegroup.Summary
			Expect(json.Unmarshal(session.Out.Contents(), &ngSummaries)).To(Succeed())
			Expect(ngSummaries).To(HaveLen(1))
			ngSummary := ngSummaries[0]
			Expect(ngSummary.NodeInstanceRoleARN).NotTo(BeEmpty())

			By("checking access entries")
			cmd = params.EksctlGetCmd.
				WithArgs(
					"accessentry",
					"--cluster", params.ClusterName,
				)
			session = cmd.Run()
			Expect(session.ExitCode()).To(BeZero())
			if nt.expectedAuthMode == authModeAccessEntry {
				Expect(session.Out.Contents()).To(ContainSubstring(ngSummary.NodeInstanceRoleARN), "failed to find access entry for nodegroup")
			} else {
				Expect(session.Out.Contents()).NotTo(ContainSubstring(ngSummary.NodeInstanceRoleARN), "found access entry for nodegroup")
			}

			By("checking aws-auth ConfigMap")
			Expect(hasAuthIdentity(ngSummary.NodeInstanceRoleARN)).To(Equal(nt.expectedAuthMode == authModeAWSAuthConfigMap))
		},
			Entry("with access entry", ngAuthTest{
				ngName:           "ng-access-entry",
				expectedAuthMode: authModeAccessEntry,
			}),

			Entry("with --update-auth-configmap", ngAuthTest{
				ngName:       "ng-update-auth",
				createNgArgs: []string{"--update-auth-configmap"},

				expectedCreateNgOutput: "--update-auth-configmap is deprecated and will be removed soon; the recommended way " +
					"to authorize nodes is by creating EKS access entries",
				expectedAuthMode: authModeAWSAuthConfigMap,
			}),

			Entry("with --update-auth-configmap=false", ngAuthTest{
				ngName:       "ng-update-auth-false",
				createNgArgs: []string{"--update-auth-configmap=false"},

				expectedCreateNgOutput: "--update-auth-configmap is deprecated and will be removed soon; eksctl now uses " +
					"EKS Access Entries to authorize nodes if it is enabled on the cluster",
			}),
		)
	})

	Context("draining nodegroup(s)", func() {
		It("should be able to drain a nodegroup", func() {
			Expect(params.EksctlDrainNodeGroupCmd.WithArgs(
				"--cluster", params.ClusterName,
				"--name", drainMng,
			)).To(RunSuccessfully())
		})
	})

	Context("deleting nodegroup(s)", func() {
		It("should be able to delete an unmanaged nodegroup", func() {
			Expect(params.EksctlDeleteCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--name", deleteNg,
				"--wait",
			)).To(RunSuccessfully())
		})
	})

	Context("fetching nodegroup(s)", Serial, func() {
		It("should be able to get all expected nodegroups", func() {
			Expect(params.EksctlGetCmd.WithArgs(
				"nodegroup",
				"-o", "json",
				"--cluster", params.ClusterName,
			)).To(RunSuccessfullyWithOutputString(BeNodeGroupsWithNamesWhich(
				ContainElement(taintsNg1),
				ContainElement(taintsNg2),
				ContainElement(scaleSingleNg),
				ContainElement(scaleMultipleNg),
				ContainElement(scaleMultipleMng),
				ContainElement(GPUMng),
				ContainElement(drainMng),
			)))
		})
	})
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	// before deleting the cluster, first delete the additional subnet
	ec2 := awsec2.NewFromConfig(NewConfig(params.Region))
	_, err := ec2.DeleteSubnet(context.Background(), &awsec2.DeleteSubnetInput{SubnetId: &extraSubnetID})
	Expect(err).NotTo(HaveOccurred())

	Expect(params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
		"--wait",
	)).To(RunSuccessfully())

	gexec.KillAndWait()
	if params.KubeconfigTemp {
		os.Remove(params.KubeconfigPath)
	}
	os.RemoveAll(params.TestDirectory)
})

func createAdditionalSubnet(cfg *api.ClusterConfig) string {
	ctl, err := eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())
	cl, err := ctl.GetCluster(context.Background(), params.ClusterName)
	Expect(err).NotTo(HaveOccurred())

	ec2 := awsec2.NewFromConfig(NewConfig(params.Region))
	existingSubnets, err := ec2.DescribeSubnets(context.Background(), &awsec2.DescribeSubnetsInput{
		SubnetIds: cl.ResourcesVpcConfig.SubnetIds,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(len(existingSubnets.Subnets) > 0).To(BeTrue())
	existingSubnet := existingSubnets.Subnets[0]

	cidr := *existingSubnet.CidrBlock
	var (
		i1, i2, i3, i4, ic int
	)
	fmt.Sscanf(cidr, "%d.%d.%d.%d/%d", &i1, &i2, &i3, &i4, &ic)
	cidr = fmt.Sprintf("%d.%d.%s.%d/%d", i1, i2, "255", i4, ic)

	var tags []ec2types.Tag

	// filter aws: tags
	for _, t := range existingSubnet.Tags {
		if !strings.HasPrefix(*t.Key, "aws:") {
			tags = append(tags, t)
		}
	}

	// create a new subnet in that given vpc and zone.
	output, err := ec2.CreateSubnet(context.Background(), &awsec2.CreateSubnetInput{
		AvailabilityZone: aws.String("us-west-2a"),
		CidrBlock:        aws.String(cidr),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
				Tags:         tags,
			},
		},
		VpcId: existingSubnet.VpcId,
	})
	Expect(err).NotTo(HaveOccurred())

	moutput, err := ec2.ModifySubnetAttribute(context.Background(), &awsec2.ModifySubnetAttributeInput{
		MapPublicIpOnLaunch: &types.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		SubnetId: output.Subnet.SubnetId,
	})
	Expect(err).NotTo(HaveOccurred(), moutput.ResultMetadata)

	subnet := output.Subnet
	routeTables, err := ec2.DescribeRouteTables(context.Background(), &awsec2.DescribeRouteTablesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []string{*existingSubnet.SubnetId},
			},
		},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(len(routeTables.RouteTables) > 0).To(BeTrue(), fmt.Sprintf("route table ended up being empty: %+v", routeTables))

	routput, err := ec2.AssociateRouteTable(context.Background(), &awsec2.AssociateRouteTableInput{
		RouteTableId: routeTables.RouteTables[0].RouteTableId,
		SubnetId:     subnet.SubnetId,
	})
	Expect(err).NotTo(HaveOccurred(), routput)

	return *subnet.SubnetId
}
