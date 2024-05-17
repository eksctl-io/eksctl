//go:build integration
// +build integration

package accessentries

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kubicorn/kubicorn/pkg/namer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/pkg/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

const (
	viewPolicyARN  = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"
	editPolicyARN  = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSEditPolicy"
	adminPolicyARN = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"

	defaultRoleName   = "eksctl-default-role"
	clusterRoleName   = "eksctl-cluster-role"
	namespaceRoleName = "eksctl-namespace-role"
)

var (
	params *tests.Params
	ctl    *eks.ClusterProvider

	defaultRoleARN   string
	clusterRoleARN   string
	namespaceRoleARN string
	err              error

	apiEnabledCluster  = "accessentries-api-enabled"
	apiDisabledCluster = "accessentries-api-disabled"
)

//go:embed testdata/node-role.yaml
var nodeRoleJSON []byte

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("")
}

func TestAccessEntries(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = SynchronizedBeforeSuite(func() []byte {
	var (
		err              error
		alreadyExistsErr *iamtypes.EntityAlreadyExistsException
	)
	maybeCreateRoleAndGetARN := func(name string) (string, error) {
		createOut, err := ctl.AWSProvider.IAM().CreateRole(context.Background(), &iam.CreateRoleInput{
			RoleName:                 aws.String(name),
			AssumeRolePolicyDocument: trustPolicy,
		})
		if err == nil {
			return *createOut.Role.Arn, nil
		}
		if !errors.As(err, &alreadyExistsErr) {
			return "", fmt.Errorf("creating role %q: %w", name, err)
		}
		getOut, err := ctl.AWSProvider.IAM().GetRole(context.Background(), &iam.GetRoleInput{
			RoleName: aws.String(name),
		})
		if err != nil {
			return "", fmt.Errorf("fetching role %q: %w", name, err)
		}
		return *getOut.Role.Arn, nil
	}
	ctl, err = eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, nil)
	Expect(err).NotTo(HaveOccurred())

	defaultRoleARN, err = maybeCreateRoleAndGetARN(defaultRoleName)
	Expect(err).NotTo(HaveOccurred())

	clusterRoleARN, err = maybeCreateRoleAndGetARN(clusterRoleName)
	Expect(err).NotTo(HaveOccurred())

	namespaceRoleARN, err = maybeCreateRoleAndGetARN(namespaceRoleName)
	Expect(err).NotTo(HaveOccurred())

	return []byte(defaultRoleARN + "," + clusterRoleARN + "," + namespaceRoleARN)
}, func(arns []byte) {
	iamARNs := strings.Split(string(arns), ",")
	defaultRoleARN, clusterRoleARN, namespaceRoleARN = iamARNs[0], iamARNs[1], iamARNs[2]

	var err error
	ctl, err = eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, nil)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("(Integration) [AccessEntries Test]", func() {

	Context("Cluster with access entries API disabled", Ordered, func() {
		var (
			cfg *api.ClusterConfig
		)

		BeforeAll(func() {
			cfg = makeClusterConfig(apiDisabledCluster)
		})

		It("should create a cluster with authenticationMode set to CONFIG_MAP and allow self-managed nodes to join via aws-auth", func() {
			cfg.AccessConfig.AuthenticationMode = ekstypes.AuthenticationModeConfigMap
			cfg.NodeGroups = append(cfg.NodeGroups, &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "aws-auth-ng",
					ScalingConfig: &api.ScalingConfig{
						DesiredCapacity: aws.Int(1),
					},
				},
			})
			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))).To(RunSuccessfully())

			Expect(ctl.RefreshClusterStatus(context.Background(), cfg)).NotTo(HaveOccurred())
			Expect(ctl.IsAccessEntryEnabled()).To(BeFalse())

			Expect(params.EksctlGetCmd.WithArgs(
				"nodegroup",
				"--cluster", apiDisabledCluster,
				"--name", "aws-auth-ng",
				"-o", "yaml",
			)).To(runner.RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("Status: CREATE_COMPLETE")),
			))
		})

		It("should fail early when trying to create access entries", func() {
			session := params.EksctlCreateCmd.
				WithArgs(
					"accessentry",
					"--cluster", apiDisabledCluster,
					"--principal-arn", defaultRoleARN,
				).Run()
			Expect(session.ExitCode()).To(Equal(1))
			Expect(session.Err.Contents()).To(ContainSubstring(accessentry.ErrDisabledAccessEntryAPI.Error()))
		})

		It("should fail early when trying to fetch access entries", func() {
			session := params.EksctlGetCmd.
				WithArgs(
					"accessentry",
					"--cluster", apiDisabledCluster,
				).Run()
			Expect(session.ExitCode()).To(Equal(1))
			Expect(session.Err.Contents()).To(ContainSubstring(accessentry.ErrDisabledAccessEntryAPI.Error()))
		})

		It("should fail early when trying to delete access entries", func() {
			session := params.EksctlDeleteCmd.
				WithArgs(
					"accessentry",
					"--cluster", apiDisabledCluster,
					"--principal-arn", defaultRoleARN,
				).Run()
			Expect(session.ExitCode()).To(Equal(1))
			Expect(session.Err.Contents()).To(ContainSubstring(accessentry.ErrDisabledAccessEntryAPI.Error()))
		})

		It("should change cluster authenticationMode to API_AND_CONFIG_MAP", func() {
			cfg.AccessConfig.AuthenticationMode = ekstypes.AuthenticationModeApiAndConfigMap
			Expect(params.EksctlUtilsCmd.
				WithArgs(
					"update-authentication-mode",
					"--cluster", apiDisabledCluster,
					"--authentication-mode", string(ekstypes.AuthenticationModeApiAndConfigMap),
				)).To(RunSuccessfully())

			Expect(ctl.RefreshClusterStatus(context.Background(), cfg)).NotTo(HaveOccurred())
			Expect(ctl.IsAccessEntryEnabled()).To(BeTrue())
		})

		It("should create access entries", func() {
			cfg.AccessConfig.AccessEntries = getAccessEntries()

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(params.EksctlCreateCmd.
				WithArgs(
					"accessentry",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))).To(RunSuccessfully())
		})

		It("should fetch all expected access entries", func() {
			var output []api.AccessEntry
			session := params.EksctlGetCmd.
				WithArgs(
					"accessentry",
					"--cluster", apiDisabledCluster,
					"--output", "json",
				).Run()
			Expect(session.ExitCode()).To(Equal(0))
			Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
			// taking into account the cluster creator admin permission access entry
			Expect(output).To(HaveLen(4))
			Expect(output).To(ContainElements(cfg.AccessConfig.AccessEntries))
		})

		It("should delete an access entry via CLI flags", func() {
			Expect(params.EksctlDeleteCmd.
				WithArgs(
					"accessentry",
					"--cluster", apiDisabledCluster,
					"--principal-arn", defaultRoleARN,
				)).To(RunSuccessfully())
		})

		It("should delete multiple access entries via config file", func() {
			cfg.AccessConfig.AccessEntries = []api.AccessEntry{
				{
					PrincipalARN: api.MustParseARN(clusterRoleARN),
				},
				{
					PrincipalARN: api.MustParseARN(namespaceRoleARN),
				},
			}

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			cmd := params.EksctlDeleteCmd.
				WithArgs(
					"accessentry",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())
		})

		It("should have removed all expected access entries", func() {
			var output []api.AccessEntry
			session := params.EksctlGetCmd.
				WithArgs(
					"accessentry",
					"--cluster", apiDisabledCluster,
					"--output", "json",
				).Run()
			Expect(session.ExitCode()).To(Equal(0))
			Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
			// taking into account the cluster creator admin permission access entry
			Expect(output).To(HaveLen(1))
		})
	})

	Context("Cluster with access entries API enabled", Ordered, func() {
		var (
			cfg *api.ClusterConfig
		)

		BeforeAll(func() {
			cfg = makeClusterConfig(apiEnabledCluster)
			cfg.AccessConfig.AccessEntries = getAccessEntries()[:2]
		})

		It("should create a cluster with access entries", func() {
			cfg.AccessConfig.BootstrapClusterCreatorAdminPermissions = aws.Bool(false)

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			// cluster creation tasks that require access to K8s API will fail
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--without-nodegroup",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).NotTo(RunSuccessfully())
		})

		It("should fetch all expected access entries", func() {
			var output []api.AccessEntry
			session := params.EksctlGetCmd.
				WithArgs(
					"accessentry",
					"--cluster", apiEnabledCluster,
					"--output", "json",
				).Run()
			Expect(session.ExitCode()).To(Equal(0))
			Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
			Expect(output).To(HaveLen(2))
			Expect(output).To(ContainElements(cfg.AccessConfig.AccessEntries))
		})

		It("should fail to delete the cluster without access to K8s API server", func() {
			session := params.EksctlDeleteCmd.
				WithArgs(
					"cluster",
					"--name", apiEnabledCluster,
					"--wait",
				).Run()
			Expect(session.ExitCode()).To(Equal(1))
			Expect(session.Err.Contents()).To(ContainSubstring("Unauthorized"))
		})

		It("should create an access entry to give admin permissions to cluster creator", func() {
			cfg.AccessConfig.AccessEntries = []api.AccessEntry{
				{
					PrincipalARN: api.MustParseARN(extractIAMRoleARN(ctl.Status.IAMRoleARN)),
					AccessPolicies: []api.AccessPolicy{
						{
							PolicyARN: api.MustParseARN(adminPolicyARN),
							AccessScope: api.AccessScope{
								Type: "cluster",
							},
						},
					},
				},
			}
			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(params.EksctlCreateCmd.
				WithArgs(
					"accessentry",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))).To(RunSuccessfully())
		})

		Context("Unowned access entries", func() {
			var (
				principalARN string
			)

			BeforeAll(func() {
				principalARN = getAccessEntries()[2].PrincipalARN.String()
				_, err := ctl.AWSProvider.EKS().CreateAccessEntry(context.Background(), &awseks.CreateAccessEntryInput{
					ClusterName:  &apiEnabledCluster,
					PrincipalArn: &principalARN,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fetch the unowned access entry", func() {
				Eventually(func() runner.Cmd {
					return params.EksctlGetCmd.
						WithArgs(
							"accessentry",
							"--cluster", apiEnabledCluster,
							"--principal-arn", principalARN,
						)
				}, "5m", "30s").Should(RunSuccessfully())
			})

			It("should delete the unowned access entry", func() {
				Expect(params.EksctlDeleteCmd.
					WithArgs(
						"accessentry",
						"--cluster", apiEnabledCluster,
						"--principal-arn", principalARN,
					)).To(RunSuccessfully())

				Eventually(func() runner.Cmd {
					return params.EksctlGetCmd.
						WithArgs(
							"accessentry",
							"--cluster", apiEnabledCluster,
							"--principal-arn", principalARN,
						)
				}, "5m", "30s").ShouldNot(RunSuccessfully())
			})
		})

		Context("self-managed nodegroup with a pre-existing instanceRoleARN", func() {
			waitAndGetRoleARN := func(ctx context.Context, stackName *string) (string, error) {
				waiter := cloudformation.NewStackCreateCompleteWaiter(ctl.AWSProvider.CloudFormation())
				output, err := waiter.WaitForOutput(ctx, &cloudformation.DescribeStacksInput{
					StackName: stackName,
				}, 5*time.Minute)
				if err != nil {
					return "", err
				}
				if len(output.Stacks) != 1 {
					return "", fmt.Errorf("expected to find 1 stack; got %d", len(output.Stacks))
				}
				var roleARN string
				if err := outputs.Collect(output.Stacks[0], map[string]outputs.Collector{
					"NodeInstanceRoleARN": func(v string) error {
						roleARN = v
						return nil
					},
				}, nil); err != nil {
					return "", err
				}
				return roleARN, nil
			}

			var roleARN string

			BeforeAll(func() {
				By("creating a CloudFormation stack for NodeInstanceRoleARN")
				stackName := aws.String(fmt.Sprintf("%s-%s-role", apiEnabledCluster, namer.RandomName()))
				ctx := context.Background()
				_, err := ctl.AWSProvider.CloudFormation().CreateStack(ctx, &cloudformation.CreateStackInput{
					StackName:    stackName,
					TemplateBody: aws.String(string(nodeRoleJSON)),
					Capabilities: []cfntypes.Capability{cfntypes.CapabilityCapabilityIam},
				})
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(func() {
					_, err := ctl.AWSProvider.CloudFormation().DeleteStack(ctx, &cloudformation.DeleteStackInput{
						StackName: stackName,
					})
					Expect(err).NotTo(HaveOccurred())
				})
				By("waiting for the stack to be created successfully")
				roleARN, err = waitAndGetRoleARN(ctx, stackName)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should create an access entry but not delete it", func() {
				clusterConfig := makeClusterConfig(apiEnabledCluster)
				ng := api.NewNodeGroup()
				ng.Name = "with-instance-role-arn"
				ng.IAM.InstanceRoleARN = roleARN
				ng.ScalingConfig = &api.ScalingConfig{
					DesiredCapacity: aws.Int(1),
				}
				clusterConfig.NodeGroups = []*api.NodeGroup{ng}

				cmd := params.EksctlCreateNodegroupCmd.
					WithArgs("--config-file", "-").
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.Reader(clusterConfig))
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(
					fmt.Sprintf("nodegroup %s: created access entry for principal ARN %q", ng.Name, roleARN),
				)))

				assertAccessEntryExists := func() {
					_, err = ctl.AWSProvider.EKS().DescribeAccessEntry(context.Background(), &awseks.DescribeAccessEntryInput{
						ClusterName:  aws.String(apiEnabledCluster),
						PrincipalArn: aws.String(roleARN),
					})
					ExpectWithOffset(1, err).NotTo(HaveOccurred())
				}
				By("asserting an access entry was created for the nodegroup")
				assertAccessEntryExists()

				By("deleting nodegroup with a pre-existing instance role ARN")
				cmd = params.EksctlDeleteCmd.
					WithArgs(
						"nodegroup",
						"--config-file", "-",
						"--wait",
					).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.Reader(clusterConfig))
				Expect(cmd).To(RunSuccessfully())
				By("asserting the associated access entry is not deleted")
				assertAccessEntryExists()
			})
		})
	})
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	if ctl == nil {
		return
	}

	Expect(params.EksctlDeleteCmd.
		WithArgs(
			"cluster",
			"--name", apiDisabledCluster,
			"--disable-nodegroup-eviction",
			"--wait",
		)).To(RunSuccessfully())

	Expect(params.EksctlDeleteCmd.
		WithArgs(
			"cluster",
			"--name", apiEnabledCluster,
			"--disable-nodegroup-eviction",
			"--wait",
		)).To(RunSuccessfully())

	_, err = ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(defaultRoleName),
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(clusterRoleName),
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(namespaceRoleName),
	})
	Expect(err).NotTo(HaveOccurred())
})

var (
	makeClusterConfig = func(clusterName string) *api.ClusterConfig {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = clusterName
		clusterConfig.Metadata.Version = params.Version
		clusterConfig.Metadata.Region = params.Region
		return clusterConfig
	}

	trustPolicy = aws.String(`{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Effect": "Allow",
			"Principal": {
			  "Service": [
				"eks.amazonaws.com"
			  ]
			},
			"Action": "sts:AssumeRole"
		  }
		]
	  }`)

	extractIAMRoleARN = func(assumedRoleARN string) string {
		roleARN := strings.Replace(assumedRoleARN, "assumed-role", "role", 1)
		roleARN = strings.Replace(roleARN, "sts", "iam", 1)
		parts := strings.Split(roleARN, "/")
		if len(parts) > 2 {
			return strings.Join(parts[:2], "/")
		}
		return roleARN
	}

	getAccessEntries = func() []api.AccessEntry {
		return []api.AccessEntry{
			{
				PrincipalARN:     api.MustParseARN(defaultRoleARN),
				KubernetesGroups: []string{"test-group"},
			},
			{
				PrincipalARN: api.MustParseARN(clusterRoleARN),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN(adminPolicyARN),
						AccessScope: api.AccessScope{
							Type: "cluster",
						},
					},
				},
			},
			{
				PrincipalARN: api.MustParseARN(namespaceRoleARN),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN(viewPolicyARN),
						AccessScope: api.AccessScope{
							Type:       "namespace",
							Namespaces: []string{"test-namespace"},
						},
					},
				},
			},
		}
	}
)
