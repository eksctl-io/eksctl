package utils_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/utils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("update authentication mode", func() {

	type updateAuthenticationModeEntry struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e updateAuthenticationModeEntry) {
		cmd := newMockCmd(append([]string{"update-authentication-mode"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing required flag --cluster", updateAuthenticationModeEntry{
			expectedErr: "Error: --cluster must be set",
		}),
		Entry("missing required flag --authentication-mode", updateAuthenticationModeEntry{
			args:        []string{"--cluster", "test"},
			expectedErr: "Error: --authentication-mode must be set",
		}),
		Entry("unsupported name argument", updateAuthenticationModeEntry{
			expectedErr: "Error: name argument is not supported",
			args:        []string{"--cluster", "test", "CONFIG_MAP"},
		}),
		Entry("setting --cluster and --config-file at the same time", updateAuthenticationModeEntry{
			expectedErr: "Error: cannot use --cluster when --config-file/-f is set",
			args:        []string{"--cluster", "test", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
		Entry("setting --authentication-mode and --config-file at the same time", updateAuthenticationModeEntry{
			expectedErr: "Error: cannot use --authentication-mode when --config-file/-f is set",
			args:        []string{"--authentication-mode", "CONFIG_MAP", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
	)

	Describe("updating the cluster", func() {
		const clusterName = "test"

		makeClusterProvider := func(currentMode ekstypes.AuthenticationMode) (*mockprovider.MockProvider, *eks.ClusterProvider) {
			provider := mockprovider.NewMockProvider()
			return provider, &eks.ClusterProvider{
				AWSProvider: provider,
				Status: &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: &ekstypes.Cluster{
							Name: aws.String(clusterName),
							AccessConfig: &ekstypes.AccessConfigResponse{
								AuthenticationMode: currentMode,
							},
						},
					},
				},
			}
		}

		makeClusterConfig := func(desiredMode ekstypes.AuthenticationMode) *api.ClusterConfig {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = clusterName
			clusterConfig.AccessConfig.AuthenticationMode = desiredMode
			return clusterConfig
		}

		When("the cluster already has the desired authentication mode", func() {
			It("skips the update", func() {
				provider, clusterProvider := makeClusterProvider(ekstypes.AuthenticationModeApiAndConfigMap)

				err := utils.UpdateAuthenticationModeForCluster(context.Background(), clusterProvider, makeClusterConfig(ekstypes.AuthenticationModeApiAndConfigMap))
				Expect(err).NotTo(HaveOccurred())
				provider.MockEKS().AssertNotCalled(GinkgoT(), "UpdateClusterConfig", mock.Anything, mock.Anything)
			})
		})

		When("the cluster has a different authentication mode", func() {
			It("updates the cluster with the desired mode", func() {
				provider, clusterProvider := makeClusterProvider(ekstypes.AuthenticationModeConfigMap)

				updateID := aws.String("update-1")
				provider.MockEKS().On("UpdateClusterConfig", mock.Anything, mock.MatchedBy(func(input *awseks.UpdateClusterConfigInput) bool {
					return *input.Name == clusterName && input.AccessConfig.AuthenticationMode == ekstypes.AuthenticationModeApiAndConfigMap
				})).Return(&awseks.UpdateClusterConfigOutput{
					Update: &ekstypes.Update{
						Id: updateID,
					},
				}, nil).Once()
				provider.MockEKS().On("DescribeUpdate", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeUpdateInput) bool {
					return *input.UpdateId == *updateID
				}), mock.Anything).Return(&awseks.DescribeUpdateOutput{
					Update: &ekstypes.Update{
						Id:     updateID,
						Type:   ekstypes.UpdateTypeAccessConfigUpdate,
						Status: ekstypes.UpdateStatusSuccessful,
					},
				}, nil).Once()

				err := utils.UpdateAuthenticationModeForCluster(context.Background(), clusterProvider, makeClusterConfig(ekstypes.AuthenticationModeApiAndConfigMap))
				Expect(err).NotTo(HaveOccurred())
				provider.MockEKS().AssertExpectations(GinkgoT())
			})
		})
	})
})
