package connector_test

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/smithy-go"
	"github.com/pkg/errors"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/connector"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type connectorCase struct {
	cluster             connector.ExternalCluster
	getManifestTemplate func() (connector.ManifestTemplate, error)

	expectedErr string
}

var _ = Describe("EKS Connector", func() {
	readManifest := func(filename string) (connector.ManifestFile, error) {
		data, err := os.ReadFile(path.Join("testdata", filename))
		if err != nil {
			return connector.ManifestFile{}, nil
		}
		return connector.ManifestFile{
			Data:     data,
			Filename: filename,
		}, nil
	}

	readResources := func() (connector.ManifestTemplate, error) {
		connectorResources, err := readManifest("eks-connector.yaml")
		if err != nil {
			return connector.ManifestTemplate{}, err
		}
		clusterRole, err := readManifest("eks-connector-clusterrole.yaml")
		if err != nil {
			return connector.ManifestTemplate{}, nil
		}
		consoleAccessResources, err := readManifest("eks-connector-console-dashboard-full-access-group.yaml")
		if err != nil {
			return connector.ManifestTemplate{}, nil
		}

		return connector.ManifestTemplate{
			Connector:     connectorResources,
			ClusterRole:   clusterRole,
			ConsoleAccess: consoleAccessResources,
		}, nil
	}

	DescribeTable("Register cluster", func(cc connectorCase) {
		mockProvider := mockprovider.NewMockProvider()

		mockDescribeCluster(mockProvider, cc.cluster.Name)
		mockProvider.MockEKS().On("RegisterCluster", mock.Anything, mock.MatchedBy(func(input *eks.RegisterClusterInput) bool {
			return *input.Name == cc.cluster.Name && string(input.ConnectorConfig.Provider) == strings.ToUpper(cc.cluster.Provider)
		})).Return(&eks.RegisterClusterOutput{
			Cluster: &ekstypes.Cluster{
				ConnectorConfig: &ekstypes.ConnectorConfigResponse{
					ActivationId:     aws.String("activation-id-123"),
					ActivationCode:   aws.String("activation-code-123"),
					ActivationExpiry: aws.Time(time.Now()),
				},
			},
		}, nil)

		mockProvider.MockSTS().On("GetCallerIdentity", mock.Anything, mock.Anything).Return(&sts.GetCallerIdentityOutput{
			Arn: aws.String("arn:aws:iam::12356:user/eksctl"),
		}, nil)

		if cc.cluster.ConnectorRoleARN == "" {
			mockIAM(mockProvider, cc.cluster.Name)
		}

		manifestTemplate, err := cc.getManifestTemplate()
		Expect(err).NotTo(HaveOccurred())

		c := connector.EKSConnector{
			Provider:         mockProvider,
			ManifestTemplate: manifestTemplate,
		}

		resourceList, err := c.RegisterCluster(context.Background(), cc.cluster)
		if cc.expectedErr != "" {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(cc.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())

		assertManifestEquals := func(m connector.ManifestFile, expectedFile string) {
			expected, err := os.ReadFile(path.Join("testdata", expectedFile))
			Expect(err).NotTo(HaveOccurred())
			Expect(m.Data).To(Equal(expected), m.Filename)
		}

		assertManifestEquals(resourceList.ConnectorResources, "eks-connector-expected.yaml")
		assertManifestEquals(resourceList.ClusterRoleResources, "eks-connector-clusterrole-expected.yaml")
		assertManifestEquals(resourceList.ConsoleAccessResources, "eks-connector-console-dashboard-full-access-group-expected.yaml")

	},
		Entry("valid name and provider", connectorCase{
			cluster: connector.ExternalCluster{
				Name:     "web",
				Provider: "gke",
			},
			getManifestTemplate: readResources,
		}),

		Entry("pre-existing IAM role", connectorCase{
			cluster: connector.ExternalCluster{
				Name:             "web",
				Provider:         "gke",
				ConnectorRoleARN: "arn:aws:iam::1234567890:role/custom-connector-role",
			},
			getManifestTemplate: readResources,
		}),

		Entry("malformed manifests", connectorCase{
			cluster: connector.ExternalCluster{
				Name:             "web",
				Provider:         "gke",
				ConnectorRoleARN: "arn:aws:iam::1234567890:role/custom-connector-role",
			},
			getManifestTemplate: func() (connector.ManifestTemplate, error) {
				return connector.ManifestTemplate{
					Connector: connector.ManifestFile{
						Data: []byte("malformed"),
					},
				}, nil
			},

			expectedErr: "unexpected error parsing manifests for EKS Connector",
		}),
	)

	Describe("Register cluster failure", func() {

		It("should suggest creating SLR if it does not exist", func() {
			cluster := connector.ExternalCluster{
				Name:             "external",
				Provider:         "gke",
				ConnectorRoleARN: "arn:aws:iam::1234567890:role/custom-connector-role",
			}

			mockProvider := mockprovider.NewMockProvider()

			mockDescribeCluster(mockProvider, cluster.Name)
			mockProvider.MockEKS().On("RegisterCluster", mock.Anything, mock.MatchedBy(func(input *eks.RegisterClusterInput) bool {
				return *input.Name == cluster.Name
			})).Return(nil, &smithy.OperationError{
				Err: errors.New("Cluster Management role arn:aws:iam::12345:role/aws-service-role/eks-connector.amazonaws.com/AWSServiceRoleForAmazonEKSConnector is not available"),
			})
			mockProvider.MockIAM().On("DeleteRole").Return(&iam.DeleteRoleOutput{}, nil)

			c := &connector.EKSConnector{
				Provider: mockProvider,
			}
			_, err := c.RegisterCluster(context.Background(), cluster)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("SLR for EKS Connector does not exist; please run `aws iam create-service-linked-role --aws-service-name eks-connector.amazonaws.com` first")))
		})

		It("should clean up IAM role if RegisterCluster fails", func() {
			cluster := connector.ExternalCluster{
				Name:     "external",
				Provider: "gke",
			}

			mockProvider := mockprovider.NewMockProvider()

			mockDescribeCluster(mockProvider, cluster.Name)
			mockProvider.MockEKS().On("RegisterCluster", mock.Anything, mock.MatchedBy(func(input *eks.RegisterClusterInput) bool {
				return *input.Name == cluster.Name
			})).Return(nil, &ekstypes.InvalidRequestException{
				Message: aws.String("test"),
			})

			mockIAM(mockProvider, cluster.Name)

			mockProvider.MockIAM().
				On("DeleteRolePolicy", mock.Anything, mock.MatchedBy(func(input *iam.DeleteRolePolicyInput) bool {
					return matchesRole(*input.RoleName) && *input.PolicyName == "eks-connector-agent"
				})).Return(&iam.DeleteRolePolicyOutput{}, nil).
				On("DeleteRole", mock.Anything, mock.MatchedBy(func(input *iam.DeleteRoleInput) bool {
					return matchesRole(*input.RoleName)
				})).Return(&iam.DeleteRoleOutput{}, nil)

			c := &connector.EKSConnector{
				Provider: mockProvider,
			}
			_, err := c.RegisterCluster(context.Background(), cluster)
			Expect(err).To(HaveOccurred())
		})
	})
})

func mockDescribeCluster(mockProvider *mockprovider.MockProvider, clusterName string) {
	mockProvider.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *eks.DescribeClusterInput) bool {
		return *input.Name == clusterName
	})).Return(nil, &ekstypes.ResourceNotFoundException{
		ClusterName: aws.String(clusterName),
	})
}

func matchesRole(roleName string) bool {
	return strings.HasPrefix(roleName, "eksctl-")
}

func mockIAM(mockProvider *mockprovider.MockProvider, clusterName string) {
	mockProvider.MockIAM().
		On("CreateRole", mock.Anything, mock.MatchedBy(func(input *iam.CreateRoleInput) bool {
			return matchesRole(*input.RoleName)
		})).Return(&iam.CreateRoleOutput{
		Role: &iamtypes.Role{
			Arn: aws.String("arn:aws:iam::1234567890:role/eksctl-12345"),
		},
	}, nil).
		On("PutRolePolicy", mock.Anything, mock.MatchedBy(func(input *iam.PutRolePolicyInput) bool {
			return matchesRole(*input.RoleName)
		})).Return(&iam.PutRolePolicyOutput{}, nil).
		On("GetRole", mock.Anything, mock.MatchedBy(func(input *iam.GetRoleInput) bool {
			return matchesRole(*input.RoleName)
		}), mock.Anything).
		Return(&iam.GetRoleOutput{
			Role: &iamtypes.Role{
				Tags: []iamtypes.Tag{
					{
						Key:   aws.String(v1alpha5.ClusterNameTag),
						Value: aws.String(clusterName),
					},
				},
			},
		}, nil)
}
