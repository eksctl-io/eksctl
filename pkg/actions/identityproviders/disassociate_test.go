package identityproviders_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

var _ = Describe("Disassociate", func() {
	var eksAPI mocksv2.EKS
	BeforeEach(func() {
		eksAPI = mocksv2.EKS{}
		eksAPI.On("DescribeIdentityProviderConfig", mock.Anything, &eks.DescribeIdentityProviderConfigInput{
			ClusterName: aws.String("idp-test"),
			IdentityProviderConfig: &ekstypes.IdentityProviderConfig{
				Name: aws.String("pool-1"),
				Type: aws.String("oidc"),
			},
		}).Return(&eks.DescribeIdentityProviderConfigOutput{
			IdentityProviderConfig: &ekstypes.IdentityProviderConfigResponse{
				Oidc: &ekstypes.OidcIdentityProviderConfig{
					IdentityProviderConfigName: aws.String("pool-1"),
					Status:                     ekstypes.ConfigStatusActive,
				},
			},
		}, nil)
		eksAPI.On("DisassociateIdentityProviderConfig", mock.Anything, &eks.DisassociateIdentityProviderConfigInput{
			ClusterName: aws.String("idp-test"),
			IdentityProviderConfig: &ekstypes.IdentityProviderConfig{
				Name: aws.String("pool-1"),
				Type: aws.String("oidc"),
			},
		}).Return(&eks.DisassociateIdentityProviderConfigOutput{
			Update: &ekstypes.Update{
				Id:   aws.String("1"),
				Type: ekstypes.UpdateTypeDisassociateIdentityProviderConfig,
			},
		}, nil)
	})
	It("disassociates from all providers", func() {
		manager := identityproviders.NewManager(api.ClusterMeta{
			Name: "idp-test",
		}, &eksAPI)
		err := manager.Disassociate(context.Background(), identityproviders.DisassociateIdentityProvidersOptions{
			Providers: []identityproviders.DisassociateIdentityProvider{
				{
					Name: "pool-1",
					Type: "oidc",
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		eksAPI.AssertExpectations(GinkgoT())
	})
	It("disassociates from all providers and waits", func() {
		manager := identityproviders.NewManager(api.ClusterMeta{
			Name: "idp-test",
		}, &eksAPI)
		updateInput := &eks.DescribeUpdateInput{
			UpdateId: aws.String("1"),
			Name:     aws.String("idp-test"),
		}
		updateOutput := &eks.DescribeUpdateOutput{
			Update: &ekstypes.Update{
				Status: ekstypes.UpdateStatusSuccessful,
				Type:   ekstypes.UpdateTypeDisassociateIdentityProviderConfig,
			},
		}
		eksAPI.On("DescribeUpdate", mock.Anything, updateInput, mock.Anything).Return(updateOutput, nil)
		err := manager.Disassociate(context.Background(), identityproviders.DisassociateIdentityProvidersOptions{
			WaitTimeout: 1 * time.Minute,
			Providers: []identityproviders.DisassociateIdentityProvider{
				{
					Name: "pool-1",
					Type: "oidc",
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		eksAPI.AssertExpectations(GinkgoT())
	})
})
