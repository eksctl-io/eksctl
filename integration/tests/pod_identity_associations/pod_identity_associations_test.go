//go:build integration
// +build integration

package podidentityassociations

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"

	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

const (
	nsInitial    = "initial"
	nsCLI        = "cli"
	nsConfigFile = "config-file"
	nsUnowned    = "unowned"

	sa1 = "service-account-1"
	sa2 = "service-account-2"
	sa3 = "service-account-3"

	initialRole1 = "pod-identity-role-11"
	initialRole2 = "pod-identity-role-12"
)

var (
	params             *tests.Params
	ctl                *eks.ClusterProvider
	role1ARN, role2ARN string
	err                error
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParamsWithGivenClusterName("pod-identity-associations", "test")
	ctl, err = eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, nil)
	if err != nil {
		panic(err)
	}
}

func TestPodIdentityAssociations(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	roleOutput, err := ctl.AWSProvider.IAM().CreateRole(context.Background(), &iam.CreateRoleInput{
		RoleName:                 aws.String(initialRole1),
		AssumeRolePolicyDocument: trustPolicy,
	})
	Expect(err).NotTo(HaveOccurred())
	role1ARN = *roleOutput.Role.Arn

	roleOutput, err = ctl.AWSProvider.IAM().CreateRole(context.Background(), &iam.CreateRoleInput{
		RoleName:                 aws.String(initialRole2),
		AssumeRolePolicyDocument: trustPolicy,
	})
	Expect(err).NotTo(HaveOccurred())
	role2ARN = *roleOutput.Role.Arn
})

var _ = Describe("(Integration) [PodIdentityAssociations Test]", Ordered, func() {

	Context("Cluster with pod identity associations", func() {
		var (
			cfg *api.ClusterConfig
		)

		BeforeAll(func() {
			cfg = makeClusterConfig()
		})

		It("should create a cluster with pod identity associations", func() {
			cfg.Addons = []*api.Addon{{Name: api.PodIdentityAgentAddon}}
			cfg.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
				{
					Namespace:          nsInitial,
					ServiceAccountName: sa1,
					RoleARN:            role1ARN,
				},
				{
					Namespace:          nsInitial,
					ServiceAccountName: sa2,
					RoleARN:            role1ARN,
				},
				{
					Namespace:          nsInitial,
					ServiceAccountName: sa3,
					RoleARN:            role1ARN,
				},
			}

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
		})

		It("should fetch all expected associations", func() {
			var output []podidentityassociation.Summary
			session := params.EksctlGetCmd.
				WithArgs(
					"podidentityassociation",
					"--cluster", params.ClusterName,
					"--output", "json",
				).Run()
			Expect(session.ExitCode()).To(Equal(0))
			Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
			Expect(output).To(HaveLen(3))
		})

		Context("Create new pod identity associations", func() {
			It("should fail to create a new association for the same namespace & service account", func() {
				Expect(params.EksctlCreateCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsInitial,
						"--service-account-name", sa1,
						"--role-arn", role1ARN,
					),
				).NotTo(RunSuccessfully())
			})

			It("should create a new association via CLI", func() {
				Expect(params.EksctlCreateCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsCLI,
						"--service-account-name", sa1,
						"--well-known-policies", "certManager",
					),
				).To(RunSuccessfully())
			})

			It("should create (multiple) associations via config file", func() {
				cfg.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
					{
						Namespace:          nsConfigFile,
						ServiceAccountName: sa1,
						WellKnownPolicies: api.WellKnownPolicies{
							AutoScaler:  true,
							ExternalDNS: true,
						},
					},
					{
						Namespace:          nsConfigFile,
						ServiceAccountName: sa2,
						PermissionPolicy:   permissionPolicy,
					},
				}

				data, err := json.Marshal(cfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(params.EksctlCreateCmd.
					WithArgs(
						"podidentityassociation",
						"--config-file", "-",
					).
					WithoutArg("--region", params.Region).
					WithStdin(bytes.NewReader(data)),
				).To(RunSuccessfully())
			})
		})

		Context("Fetching pod identity associations", func() {
			It("should fetch all associations for a cluster", func() {
				var output []podidentityassociation.Summary
				session := params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--output", "json",
					).Run()
				Expect(session.ExitCode()).To(Equal(0))
				Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
				Expect(output).To(HaveLen(6))
			})

			It("should fetch all associations for a namespace", func() {
				var output []podidentityassociation.Summary
				session := params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsConfigFile,
						"--output", "json",
					).Run()
				Expect(session.ExitCode()).To(Equal(0))
				Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
				Expect(output).To(HaveLen(2))
			})

			It("should fetch a single association defined by namespace & service account", func() {
				var output []podidentityassociation.Summary
				session := params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsConfigFile,
						"--service-account-name", sa1,
						"--output", "json",
					).Run()
				Expect(session.ExitCode()).To(Equal(0))
				Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
				Expect(output).To(HaveLen(1))
			})
		})

		Context("Updating pod identity associations", func() {
			It("should fail to update an association with role created by eksctl", func() {
				Expect(params.EksctlUpdateCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsCLI,
						"--service-account-name", sa1,
						"--role-arn", role1ARN,
					),
				).NotTo(RunSuccessfully())
			})

			It("should update an association via CLI", func() {
				Expect(params.EksctlUpdateCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsInitial,
						"--service-account-name", sa1,
						"--role-arn", role2ARN,
					),
				).To(RunSuccessfully())
			})

			It("should update (multiple) associations via config file", func() {
				cfg.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
					{
						Namespace:          nsInitial,
						ServiceAccountName: sa2,
						RoleARN:            role2ARN,
					},
					{
						Namespace:          nsInitial,
						ServiceAccountName: sa3,
						RoleARN:            role2ARN,
					},
				}

				data, err := json.Marshal(cfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(params.EksctlUpdateCmd.
					WithArgs(
						"podidentityassociation",
						"--config-file", "-",
					).
					WithoutArg("--region", params.Region).
					WithStdin(bytes.NewReader(data)),
				).To(RunSuccessfully())
			})

			It("should check all associations were updated successfully", func() {
				var output []podidentityassociation.Summary
				session := params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsInitial,
						"--output", "json",
					).Run()
				Expect(session.ExitCode()).To(Equal(0))
				Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
				Expect(output).To(HaveLen(3))
				Expect(output[0].RoleARN).To(Equal(role2ARN))
				Expect(output[1].RoleARN).To(Equal(role2ARN))
				Expect(output[2].RoleARN).To(Equal(role2ARN))
			})
		})

		Context("Deleting pod identity associations", func() {
			It("should delete an association via CLI", func() {
				Expect(params.EksctlDeleteCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsInitial,
						"--service-account-name", sa1,
					),
				).To(RunSuccessfully())
			})

			It("should delete (multiple) associations via config file", func() {
				cfg.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
					{
						Namespace:          nsInitial,
						ServiceAccountName: sa2,
					},
					{
						Namespace:          nsInitial,
						ServiceAccountName: sa3,
					},
				}

				data, err := json.Marshal(cfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(params.EksctlDeleteCmd.
					WithArgs(
						"podidentityassociation",
						"--config-file", "-",
					).
					WithoutArg("--region", params.Region).
					WithStdin(bytes.NewReader(data)),
				).To(RunSuccessfully())
			})

			It("should check that all associations were deleted successfully", func() {
				Expect(params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsInitial,
					)).To(RunSuccessfullyWithOutputStringLines(ContainElement("No podidentityassociations found")))
			})
		})

		Context("Unowned pod identity association", func() {
			BeforeAll(func() {
				ctl.AWSProvider.EKS().CreatePodIdentityAssociation(context.Background(),
					&awseks.CreatePodIdentityAssociationInput{
						ClusterName:    &params.ClusterName,
						Namespace:      aws.String(nsUnowned),
						ServiceAccount: aws.String(sa1),
						RoleArn:        &role1ARN,
					})
			})

			It("should fetch an unowned association", func() {
				Expect(params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsUnowned,
						"--service-account-name", sa1,
						"--output", "json",
					)).To(RunSuccessfullyWithOutputStringLines(ContainElements(
					ContainSubstring(nsUnowned),
					ContainSubstring(sa1),
				)))
			})

			It("should delete an unowned association", func() {
				Expect(params.EksctlDeleteCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsUnowned,
						"--service-account-name", sa1,
					)).To(RunSuccessfully())

				Expect(params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", params.ClusterName,
						"--namespace", nsUnowned,
						"--service-account-name", sa1,
					)).To(RunSuccessfullyWithOutputStringLines(ContainElement("No podidentityassociations found")))
			})
		})
	})
})

var _ = AfterSuite(func() {
	_, err = ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(initialRole1),
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(initialRole2),
	})
	Expect(err).NotTo(HaveOccurred())

	params.DeleteClusters()
})

var (
	makeClusterConfig = func() *api.ClusterConfig {
		cfg := api.NewClusterConfig()
		cfg.Metadata.Name = params.ClusterName
		cfg.Metadata.Version = params.Version
		cfg.Metadata.Region = params.Region
		return cfg
	}

	trustPolicy = aws.String(`{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Effect": "Allow",
			"Principal": {
			  "Service": [
				"beta.pods.eks.aws.internal"
			  ]
			},
			"Action": [
			  "sts:AssumeRole",
			  "sts:TagSession"
			]
		  }
		]
	  }`)

	permissionPolicy = api.InlineDocument{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"autoscaling:DescribeAutoScalingGroups",
					"autoscaling:DescribeAutoScalingInstances",
				},
				"Resource": "*",
			},
		},
	}
)
