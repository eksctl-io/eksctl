//go:build integration
// +build integration

package accessentries

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

const (
	viewPolicyARN  = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"
	editPolicyARN  = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSEditPolicy"
	adminPolicyARN = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"

	userName          = "test-user"
	clusterRoleName   = "test-cluster-role"
	namespaceRoleName = "test-namespace-role"
)

var trustPolicy = aws.String(`{
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

var (
	params           *tests.Params
	ctl              *eks.ClusterProvider
	userARN          string
	clusterRoleARN   string
	namespaceRoleARN string
	err              error
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("accessentries-api-disabled")
	ctl, err = eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, getInitialClusterConfig())
	if err != nil {
		panic(err)
	}
}

func TestAccessEntries(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = SynchronizedBeforeSuite(func() []byte {
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
		RoleName:                 aws.String(clusterRoleName),
		AssumeRolePolicyDocument: trustPolicy,
		Tags: []iamtypes.Tag{
			{
				Key:   aws.String(api.ClusterNameTag),
				Value: aws.String(params.ClusterName),
			},
		},
	})
	Expect(err).NotTo(HaveOccurred())
	clusterRoleARN = *roleOutput.Role.Arn

	roleOutput, err = ctl.AWSProvider.IAM().CreateRole(context.Background(), &iam.CreateRoleInput{
		RoleName:                 aws.String(namespaceRoleName),
		AssumeRolePolicyDocument: trustPolicy,
		Tags: []iamtypes.Tag{
			{
				Key:   aws.String(api.ClusterNameTag),
				Value: aws.String(params.ClusterName),
			},
		},
	})
	Expect(err).NotTo(HaveOccurred())
	namespaceRoleARN = *roleOutput.Role.Arn

	return []byte(userARN + "," + clusterRoleARN + "," + namespaceRoleARN)
}, func(arns []byte) {
	iamARNs := strings.Split(string(arns), ",")
	userARN, clusterRoleARN, namespaceRoleARN = iamARNs[0], iamARNs[1], iamARNs[2]
})

var _ = Describe("(Integration) [AccessEntries Test]", func() {

	Context("Cluster without access entries", Ordered, func() {
		var (
			cfg *api.ClusterConfig
		)

		BeforeAll(func() {
			cfg = getInitialClusterConfig()
		})

		It("should create a cluster with default authenticationMode set to CONFIG_MAP", func() {
			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--without-nodegroup",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))).To(RunSuccessfully())

			Expect(ctl.RefreshClusterStatus(context.Background(), cfg)).NotTo(HaveOccurred())
			Expect(ctl.IsAccessEntryEnabled()).To(BeFalse())
		})

		It("should fail early when trying to create access entries", func() {
			Expect(params.EksctlCreateCmd.
				WithArgs(
					"accessentry",
					"--cluster", params.ClusterName,
					"--principal-arn", userARN,
				)).NotTo(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(accessentry.ErrDisabledAccessEntryAPI.Error())),
			))
		})

		It("should fail early when trying to fetch access entries", func() {
			Expect(params.EksctlGetCmd.
				WithArgs(
					"accessentry",
					"--cluster", params.ClusterName,
				)).NotTo(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(accessentry.ErrDisabledAccessEntryAPI.Error())),
			))
		})

		It("should fail early when trying to delete access entries", func() {
			Expect(params.EksctlDeleteCmd.
				WithArgs(
					"accessentry",
					"--cluster", params.ClusterName,
					"--principal-arn", userARN,
				)).NotTo(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(accessentry.ErrDisabledAccessEntryAPI.Error())),
			))
		})

		It("should change cluster authenticationMode to API_AND_CONFIG_MAP", func() {
			Expect(params.EksctlUtilsCmd.
				WithArgs(
					"--cluster", params.ClusterName,
					"--authentication-mode", string(ekstypes.AuthenticationModeApiAndConfigMap),
				)).To(RunSuccessfully())

			Expect(ctl.RefreshClusterStatus(context.Background(), cfg)).NotTo(HaveOccurred())
			Expect(ctl.IsAccessEntryEnabled()).To(BeTrue())
		})

		It("should create access entries", func() {
			addAccessEntriesToConfig(cfg)

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
			Eventually(func() runner.Cmd {
				return params.EksctlGetCmd.
					WithArgs(
						"accessentry",
						"--cluster", params.ClusterName,
					)
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(SatisfyAll(
				ContainElement(ContainSubstring(viewPolicyARN)),
				ContainElement(ContainSubstring(adminPolicyARN)),
				ContainElement(ContainSubstring("default")),
				ContainElement(ContainSubstring("dev")),
			)))
		})

		It("should delete an access entry via CLI flags", func() {
			Expect(params.EksctlDeleteCmd.
				WithArgs(
					"accessentry",
					"--cluster", params.ClusterName,
					"--principal-arn", userARN,
				)).To(RunSuccessfully())
		})

		It("should delete multiple access entries via config file", func() {
			clusterConfig := getInitialClusterConfig()
			clusterConfig.AccessConfig.AccessEntries = append(clusterConfig.AccessConfig.AccessEntries,
				api.AccessEntry{
					PrincipalARN: api.MustParseARN(clusterRoleARN),
				},
				api.AccessEntry{
					PrincipalARN: api.MustParseARN(namespaceRoleARN),
				},
			)

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

		It("should have removed all access entries", func() {

		})
	})

	Context("Cluster with access entries", Ordered, func() {
		var (
			cfg *api.ClusterConfig
		)

		BeforeAll(func() {
			cfg = getInitialClusterConfig()
			cfg.Metadata.Name = params.NewClusterName("accessentries-api-enabled")
		})

		It("should create a cluster with access entries", func() {
			addAccessEntriesToConfig(cfg)
			cfg.AccessConfig.AuthenticationMode = ekstypes.AuthenticationModeApiAndConfigMap

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--without-nodegroup",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())
		})

		It("should fetch all expected access entries", func() {
			Eventually(func() runner.Cmd {
				return params.EksctlGetCmd.
					WithArgs(
						"accessentries",
						"--cluster", cfg.Metadata.Name,
						"--verbose", "2",
						"--output", "yaml",
					)
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(SatisfyAll(
				ContainElement(ContainSubstring(userARN)),
				ContainElement(ContainSubstring(viewPolicyARN)),
				ContainElement(ContainSubstring(editPolicyARN)),
			)))
		})
	})
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	_, err := ctl.AWSProvider.IAM().DeleteUser(context.Background(), &iam.DeleteUserInput{
		UserName: aws.String(userName),
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

	params.DeleteClusters()
})

func getInitialClusterConfig() *api.ClusterConfig {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
	clusterConfig.Metadata.Version = params.Version
	clusterConfig.Metadata.Region = params.Region
	return clusterConfig
}

func addAccessEntriesToConfig(cfg *api.ClusterConfig) {
	cfg.AccessConfig.AccessEntries = append(cfg.AccessConfig.AccessEntries,
		api.AccessEntry{
			PrincipalARN:     api.MustParseARN(userARN),
			KubernetesGroups: []string{"group1", "group2"},
		},
		api.AccessEntry{
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
		api.AccessEntry{
			PrincipalARN: api.MustParseARN(namespaceRoleARN),
			AccessPolicies: []api.AccessPolicy{
				{
					PolicyARN: api.MustParseARN(viewPolicyARN),
					AccessScope: api.AccessScope{
						Type:       "namespace",
						Namespaces: []string{"default", "dev"},
					},
				},
			},
		},
	)
}
