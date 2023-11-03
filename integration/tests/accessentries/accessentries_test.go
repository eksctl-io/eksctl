//go:build integration
// +build integration

package accessentries

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

const (
	nodePolicyARN  = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSNodePolicy"
	viewPolicyARN  = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"
	editPolicyARN  = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSEditPolicy"
	adminPolicyARN = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
	userName       = "test-user"
	roleName       = "test-role"
)

var (
	params  *tests.Params
	ctl     *eks.ClusterProvider
	userARN string
	roleARN string
	err     error
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("accessentries")
	ctl, err = eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, getInitialClusterConfig())
	if err != nil {
		panic(err)
	}
}

func TestAccessEntries(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	userOutput, err := ctl.AWSProvider.IAM().CreateUser(context.Background(), &iam.CreateUserInput{
		UserName: aws.String(userName),
		Tags: []iamtypes.Tag{
			{
				Key:   aws.String(api.ClusterNameTag),
				Value: aws.String(params.ClusterName),
			},
		},
	})
	Expect(err).NotTo(HaveOccurred())
	userARN = *userOutput.User.Arn

	roleOutput, err := ctl.AWSProvider.IAM().CreateRole(context.Background(), &iam.CreateRoleInput{
		RoleName: aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(`{
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
		  }`),
		Tags: []iamtypes.Tag{
			{
				Key:   aws.String(api.ClusterNameTag),
				Value: aws.String(params.ClusterName),
			},
		},
	})
	Expect(err).NotTo(HaveOccurred())
	roleARN = *roleOutput.Role.Arn

	cfg := getInitialClusterConfig()
	cfg.AccessEntries = []api.AccessEntry{
		{
			PrincipalARN: api.MustParseARN(userARN),
			AccessPolicies: []api.AccessPolicy{
				{
					PolicyARN: api.MustParseARN(viewPolicyARN),
					AccessScope: api.AccessScope{
						Type: "cluster",
					},
				},
				{
					PolicyARN: api.MustParseARN(editPolicyARN),
					AccessScope: api.AccessScope{
						Type:       "namespace",
						Namespaces: []string{"default"},
					},
				},
			},
		},
	}

	data, err := json.Marshal(cfg)
	Expect(err).NotTo(HaveOccurred())

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

var _ = Describe("(Integration) [AccessEntries Test]", func() {

	Describe("Cluster with access entries", func() {
		It("Should have created a cluster with access entries", func() {
			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"accessentries",
						"--cluster", params.ClusterName,
						"--verbose", "2",
						"--output", "yaml",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(SatisfyAll(
				ContainElement(ContainSubstring(userARN)),
				ContainElement(ContainSubstring(viewPolicyARN)),
				ContainElement(ContainSubstring(editPolicyARN)),
			)))
		})

		It("Should be able to create a new access entry", func() {
			clusterConfig := getInitialClusterConfig()
			clusterConfig.AccessEntries = append(clusterConfig.AccessEntries,
				api.AccessEntry{
					PrincipalARN:     api.MustParseARN(roleARN),
					KubernetesGroups: []string{"default"},
					AccessPolicies: []api.AccessPolicy{
						{
							PolicyARN: api.MustParseARN(adminPolicyARN),
							AccessScope: api.AccessScope{
								Type: "cluster",
							},
						},
					},
				})
			data, err := json.Marshal(clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"accessentry",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"accessentry",
						"--cluster", params.ClusterName,
						"--principalARN", roleARN,
						"--output", "yaml",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(SatisfyAll(
				ContainElement(ContainSubstring(adminPolicyARN)),
				ContainElement(ContainSubstring("default")),
			)))
		})

		It("Should be able to delete an access entry via CLI flags", func() {
			cmd := params.EksctlDeleteCmd.
				WithArgs(
					"accessentry",
					"--cluster", params.ClusterName,
					"--principalARN", roleARN,
				)
			Expect(cmd).To(RunSuccessfully())
		})

		It("Should be able to delete an access entry via config file", func() {
			clusterConfig := getInitialClusterConfig()
			clusterConfig.AccessEntries = append(clusterConfig.AccessEntries,
				api.AccessEntry{
					PrincipalARN: api.MustParseARN(userARN),
				})

			data, err := json.Marshal(clusterConfig)
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
	})

	Describe("Managed nodegroup authorization via access entries", func() {
		It("Should create a manage nodegroup and associated access entry", func() {
			Skip("skipping until this functionality is merged") //TODO: remove this

			clusterConfig := getInitialClusterConfig()
			clusterConfig.ManagedNodeGroups = append(clusterConfig.ManagedNodeGroups,
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "mng",
					},
				})
			data, err := json.Marshal(clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"accessentry",
						"--cluster", params.ClusterName,
						"--output", "yaml",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(SatisfyAll(
				ContainElement(ContainSubstring("NodeInstanceRole")),
				ContainElement(ContainSubstring(nodePolicyARN)),
			)))
		})

		It("Should delete the managed nodegroup and associated access entry", func() {
			Skip("skipping until this functionality is merged") //TODO: remove this

			cmd := params.EksctlDeleteCmd.
				WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--name", "mng",
				)
			Expect(cmd).To(RunSuccessfully())
		})
	})
})

var _ = AfterSuite(func() {
	_, err := ctl.AWSProvider.IAM().DeleteUser(context.Background(), &iam.DeleteUserInput{
		UserName: aws.String(userName),
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	Expect(err).NotTo(HaveOccurred())

	cmd := params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
		"--verbose", "2",
	)
	Expect(cmd).To(RunSuccessfully())
})

func getInitialClusterConfig() *api.ClusterConfig {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
	clusterConfig.Metadata.Version = api.LatestVersion
	clusterConfig.Metadata.Region = params.Region

	return clusterConfig
}
