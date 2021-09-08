package connector_test

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/connector"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type connectorCase struct {
	cluster             connector.ExternalCluster
	getManifestTemplate func() (connector.ManifestTemplate, error)

	expectedErr string
}

var _ = Describe("EKS Connector", func() {
	readResources := func() (connector.ManifestTemplate, error) {
		connectorResources, err := ioutil.ReadFile("testdata/eks-connector.yaml")
		if err != nil {
			return connector.ManifestTemplate{}, err
		}
		bindingResources, err := ioutil.ReadFile("testdata/eks-connector-binding.yaml")
		if err != nil {
			return connector.ManifestTemplate{}, nil
		}
		return connector.ManifestTemplate{
			Connector:   connectorResources,
			RoleBinding: bindingResources,
		}, nil
	}

	DescribeTable("Register cluster", func(cc connectorCase) {
		mockProvider := mockprovider.NewMockProvider()

		mockProvider.MockEKS().On("RegisterCluster", mock.MatchedBy(func(input *eks.RegisterClusterInput) bool {
			return *input.Name == cc.cluster.Name && *input.ConnectorConfig.Provider == strings.ToUpper(cc.cluster.Provider)
		})).Return(&eks.RegisterClusterOutput{
			Cluster: &eks.Cluster{
				ConnectorConfig: &eks.ConnectorConfigResponse{
					ActivationId:     aws.String("activation-id-123"),
					ActivationCode:   aws.String("activation-code-123"),
					ActivationExpiry: aws.Time(time.Now()),
				},
			},
		}, nil).On("DescribeCluster", mock.MatchedBy(func(input *eks.DescribeClusterInput) bool {
			return *input.Name == cc.cluster.Name
		})).Return(nil, &eks.ResourceNotFoundException{
			ClusterName: aws.String(cc.cluster.Name),
		})

		mockProvider.MockSTS().On("GetCallerIdentity", mock.Anything).Return(&sts.GetCallerIdentityOutput{
			Arn: aws.String("arn:aws:iam::12356:user/eksctl"),
		}, nil)

		if cc.cluster.ConnectorRole == "" {
			matchesRole := func(roleName string) bool {
				return strings.HasPrefix(roleName, "eksctl-")
			}
			mockProvider.MockIAM().On("CreateRole", mock.MatchedBy(func(input *iam.CreateRoleInput) bool {
				return matchesRole(*input.RoleName)
			})).Return(&iam.CreateRoleOutput{
				Role: &iam.Role{
					Arn: aws.String("arn:aws:iam::1234567890:role/eksctl-12345"),
				},
			}, nil).On("PutRolePolicy", mock.MatchedBy(func(input *iam.PutRolePolicyInput) bool {
				return matchesRole(*input.RoleName)
			})).Return(&iam.PutRolePolicyOutput{}, nil).On("WaitUntilRoleExists", mock.MatchedBy(func(input *iam.GetRoleInput) bool {
				return matchesRole(*input.RoleName)
			})).Return(nil)
		}

		manifestTemplate, err := cc.getManifestTemplate()
		Expect(err).ToNot(HaveOccurred())

		c := connector.EKSConnector{
			Provider:         mockProvider,
			ManifestTemplate: manifestTemplate,
		}

		resourceList, err := c.RegisterCluster(cc.cluster)
		if cc.expectedErr != "" {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(cc.expectedErr)))
			return
		}

		Expect(err).ToNot(HaveOccurred())

		assertFileEquals := func(filename string, actual []byte) {
			expected, err := ioutil.ReadFile(filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected), filename)
		}

		assertFileEquals("testdata/eks-connector-expected.yaml", resourceList.ConnectorResources)

		assertFileEquals("testdata/eks-connector-binding-expected.yaml", resourceList.ClusterRoleResources)
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
				Name:          "web",
				Provider:      "gke",
				ConnectorRole: "arn:aws:iam::1234567890:role/custom-connector-role",
			},
			getManifestTemplate: readResources,
		}),

		Entry("malformed manifests", connectorCase{
			cluster: connector.ExternalCluster{
				Name:          "web",
				Provider:      "gke",
				ConnectorRole: "arn:aws:iam::1234567890:role/custom-connector-role",
			},
			getManifestTemplate: func() (connector.ManifestTemplate, error) {
				return connector.ManifestTemplate{
					Connector: []byte("malformed"),
				}, nil
			},

			expectedErr: "unexpected error parsing manifests for EKS Connector",
		}),
	)
})
