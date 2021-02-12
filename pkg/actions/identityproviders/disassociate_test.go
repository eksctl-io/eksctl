package identityproviders_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Disassociate", func() {
	var eksAPI mocks.EKSAPI
	BeforeEach(func() {
		eksAPI = mocks.EKSAPI{}
		eksAPI.On("DescribeIdentityProviderConfig", &eks.DescribeIdentityProviderConfigInput{
			ClusterName: aws.String(""),
			IdentityProviderConfig: &eks.IdentityProviderConfig{
				Name: aws.String("pool-1"),
				Type: aws.String("oidc"),
			},
		}).Return(&eks.DescribeIdentityProviderConfigOutput{
			IdentityProviderConfig: &eks.IdentityProviderConfigResponse{
				Oidc: &eks.OidcIdentityProviderConfig{
					IdentityProviderConfigName: aws.String("pool-1"),
					Status:                     aws.String(eks.ConfigStatusActive),
				},
			},
		}, nil)
		eksAPI.On("DisassociateIdentityProviderConfig", &eks.DisassociateIdentityProviderConfigInput{
			ClusterName: aws.String(""),
			IdentityProviderConfig: &eks.IdentityProviderConfig{
				Name: aws.String("pool-1"),
				Type: aws.String("oidc"),
			},
		}).Return(&eks.DisassociateIdentityProviderConfigOutput{
			Update: &eks.Update{
				Id:   aws.String("1"),
				Type: aws.String("DisassociateIdentityProviderConfig"),
			},
		}, nil)
	})
	It("disassociates from all providers", func() {
		manager := identityproviders.NewManager(api.ClusterMeta{}, &eksAPI)
		err := manager.Disassociate(identityproviders.DisassociateIdentityProvidersOptions{
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
		manager := identityproviders.NewManager(api.ClusterMeta{}, &eksAPI)
		client := mockprovider.NewMockAWSClient()
		updateInput := eks.DescribeUpdateInput{
			UpdateId: aws.String("1"),
			Name:     aws.String(""),
		}
		updateOutput := eks.DescribeUpdateOutput{
			Update: &eks.Update{
				Status: aws.String(eks.UpdateStatusSuccessful),
				Type:   aws.String("DisassociateIdentityProviderConfig"),
			},
		}
		eksAPI.On("DescribeUpdateRequest", &updateInput).Return(
			client.MockRequestForGivenOutput(&updateInput, &updateOutput), &updateOutput,
		)
		wait := 1 * time.Minute
		err := manager.Disassociate(identityproviders.DisassociateIdentityProvidersOptions{
			WaitTimeout: &wait,
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
