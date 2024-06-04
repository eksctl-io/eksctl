//go:build integration
// +build integration

package podidentityassociations

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

const (
	clusterIRSAv1 = "iam-service-accounts"
	clusterIRSAv2 = "pod-identity-associations"

	nsDefault    = "default"
	nsInitial    = "initial"
	nsCLI        = "cli"
	nsConfigFile = "config-file"
	nsUnowned    = "unowned"

	sa1 = "service-account-1"
	sa2 = "service-account-2"
	sa3 = "service-account-3"

	rolePrefix   = "eksctl-"
	initialRole1 = rolePrefix + "pod-identity-role-1"
	initialRole2 = rolePrefix + "pod-identity-role-2"
)

var (
	params             *tests.Params
	ctl                *eks.ClusterProvider
	role1ARN, role2ARN string
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParamsWithGivenClusterName("", "test")
}

func TestPodIdentityAssociations(t *testing.T) {
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

	role1ARN, err := maybeCreateRoleAndGetARN(initialRole1)
	Expect(err).NotTo(HaveOccurred())

	role2ARN, err := maybeCreateRoleAndGetARN(initialRole2)
	Expect(err).NotTo(HaveOccurred())

	return []byte(role1ARN + "," + role2ARN)
}, func(arns []byte) {
	roleARNs := strings.Split(string(arns), ",")
	role1ARN, role2ARN = roleARNs[0], roleARNs[1]

	var err error
	ctl, err = eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, nil)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("(Integration) [PodIdentityAssociations Test]", func() {

	Context("Cluster with iam service accounts", Ordered, func() {
		var (
			cfg *api.ClusterConfig
		)

		BeforeAll(func() {
			cfg = makeClusterConfig(clusterIRSAv1)
		})

		It("should create a cluster with iam service accounts", func() {
			cfg.IAM = &api.ClusterIAM{
				WithOIDC: aws.Bool(true),
				ServiceAccounts: []*api.ClusterIAMServiceAccount{
					{
						ClusterIAMMeta: api.ClusterIAMMeta{
							Name: sa1,
						},
						AttachPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"},
					},
					{
						ClusterIAMMeta: api.ClusterIAMMeta{
							Name: sa2,
						},
						AttachRoleARN: role1ARN,
					},
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

			awsConfig := NewConfig(params.Region)
			stackNamePrefix := fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", clusterIRSAv1)
			Expect(awsConfig).To(HaveExistingStack(stackNamePrefix + "default-service-account-1"))
		})

		It("should migrate to pod identity associations", func() {
			Expect(params.EksctlUtilsCmd.
				WithArgs(
					"migrate-to-pod-identity",
					"--cluster", clusterIRSAv1,
					"--remove-oidc-provider-trust-relationship",
					"--approve",
				)).To(RunSuccessfully())
		})

		It("should fetch all expected associations", func() {
			var output []podidentityassociation.Summary
			session := params.EksctlGetCmd.
				WithArgs(
					"podidentityassociation",
					"--cluster", clusterIRSAv1,
					"--output", "json",
				).Run()
			Expect(session.ExitCode()).To(Equal(0))
			Expect(json.Unmarshal(session.Out.Contents(), &output)).To(Succeed())
			Expect(output).To(HaveLen(3))
		})

		It("should not return any iam service accounts", func() {
			Expect(params.EksctlGetCmd.
				WithArgs(
					"iamserviceaccount",
					"--cluster", clusterIRSAv1,
				)).To(RunSuccessfullyWithOutputStringLines(ContainElement("No iamserviceaccounts found")))
		})

		It("should fail to update an owned migrated role", func() {
			session := params.EksctlUpdateCmd.
				WithArgs(
					"podidentityassociation",
					"--cluster", clusterIRSAv1,
					"--namespace", nsDefault,
					"--service-account-name", sa1,
					"--role-arn", role1ARN,
				).Run()
			Expect(session.ExitCode()).To(Equal(1))
			Expect(session.Err.Contents()).To(ContainSubstring("cannot change podIdentityAssociation.roleARN since the role was created by eksctl"))
		})

		It("should update an unowned migrated role", func() {
			Expect(params.EksctlUpdateCmd.
				WithArgs(
					"podidentityassociation",
					"--cluster", clusterIRSAv1,
					"--namespace", nsDefault,
					"--service-account-name", sa2,
					"--role-arn", role1ARN,
				),
			).To(RunSuccessfully())
		})

		It("should delete an owned migrated role", func() {
			Expect(params.EksctlDeleteCmd.
				WithArgs(
					"podidentityassociation",
					"--cluster", clusterIRSAv1,
					"--namespace", nsDefault,
					"--service-account-name", sa1,
				)).To(RunSuccessfully())
		})
	})

	Context("Cluster with pod identity associations", Ordered, func() {
		var (
			cfg *api.ClusterConfig
		)

		BeforeAll(func() {
			cfg = makeClusterConfig(clusterIRSAv2)
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
					"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
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
				session := params.EksctlUpdateCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", clusterIRSAv2,
						"--namespace", nsCLI,
						"--service-account-name", sa1,
						"--role-arn", role1ARN,
					).Run()
				Expect(session.ExitCode()).To(Equal(1))
				Expect(session.Err.Contents()).To(ContainSubstring("cannot change podIdentityAssociation.roleARN since the role was created by eksctl"))
			})

			It("should update an association via CLI", func() {
				Expect(params.EksctlUpdateCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
						"--namespace", nsInitial,
					)).To(RunSuccessfullyWithOutputStringLines(ContainElement("No podidentityassociations found")))
			})
		})

		Context("Unowned pod identity association", func() {
			BeforeAll(func() {
				_, err := ctl.AWSProvider.EKS().CreatePodIdentityAssociation(context.Background(),
					&awseks.CreatePodIdentityAssociationInput{
						ClusterName:    aws.String(clusterIRSAv2),
						Namespace:      aws.String(nsUnowned),
						ServiceAccount: aws.String(sa1),
						RoleArn:        &role1ARN,
					})
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fetch an unowned association", func() {
				Expect(params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", clusterIRSAv2,
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
						"--cluster", clusterIRSAv2,
						"--namespace", nsUnowned,
						"--service-account-name", sa1,
					)).To(RunSuccessfully())

				Expect(params.EksctlGetCmd.
					WithArgs(
						"podidentityassociation",
						"--cluster", clusterIRSAv2,
						"--namespace", nsUnowned,
						"--service-account-name", sa1,
					)).To(RunSuccessfullyWithOutputStringLines(ContainElement("No podidentityassociations found")))
			})
		})
	})
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	if ctl == nil {
		return
	}

	Expect(params.EksctlDeleteCmd.WithArgs(
		"cluster", clusterIRSAv1,
	)).To(RunSuccessfully())

	Expect(params.EksctlDeleteCmd.WithArgs(
		"cluster", clusterIRSAv2,
	)).To(RunSuccessfully())

	_, err := ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(initialRole1),
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ctl.AWSProvider.IAM().DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(initialRole2),
	})
	Expect(err).NotTo(HaveOccurred())
})

var (
	makeClusterConfig = func(clusterName string) *api.ClusterConfig {
		cfg := api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		cfg.Metadata.Version = params.Version
		cfg.Metadata.Region = params.Region
		return cfg
	}

	trustPolicy = aws.String(fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Effect": "Allow",
			"Principal": {
			  "Service": [
				"%s"
			  ]
			},
			"Action": [
			  "sts:AssumeRole",
			  "sts:TagSession"
			]
		  }
		]
	  }`, api.EKSServicePrincipal))

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
