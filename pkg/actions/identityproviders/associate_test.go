package identityproviders_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Associate", func() {
	var eksAPI *mocksv2.EKS
	BeforeEach(func() {
		eksAPI = &mocksv2.EKS{}
		eksAPI.On("AssociateIdentityProviderConfig", mock.Anything, &eks.AssociateIdentityProviderConfigInput{
			ClusterName: aws.String("idp-test"),
			Oidc: &ekstypes.OidcIdentityProviderConfigRequest{
				IdentityProviderConfigName: aws.String("pool-1"),
				IssuerUrl:                  aws.String("url"),
				ClientId:                   aws.String("id"),
				UsernameClaim:              aws.String("usernameClaim"),
				UsernamePrefix:             aws.String("usernamePrefix"),
				GroupsClaim:                aws.String("groupsClaim"),
				GroupsPrefix:               aws.String("groupsPrefix"),
				RequiredClaims:             map[string]string{"permission": "true"},
			},
			Tags: map[string]string{"department": "a"},
		}).Return(&eks.AssociateIdentityProviderConfigOutput{
			Update: &ekstypes.Update{
				Id:   aws.String("1"),
				Type: ekstypes.UpdateTypeAssociateIdentityProviderConfig,
			},
		}, nil)
	})
	It("associates with all providers", func() {
		manager := identityproviders.NewManager(api.ClusterMeta{
			Name: "idp-test",
		}, eksAPI)
		err := manager.Associate(context.Background(), identityproviders.AssociateIdentityProvidersOptions{
			Providers: []api.IdentityProvider{
				{Inner: &api.OIDCIdentityProvider{
					Name:           "pool-1",
					IssuerURL:      "url",
					ClientID:       "id",
					UsernameClaim:  "usernameClaim",
					UsernamePrefix: "usernamePrefix",
					GroupsClaim:    "groupsClaim",
					GroupsPrefix:   "groupsPrefix",
					RequiredClaims: map[string]string{"permission": "true"},
					Tags:           map[string]string{"department": "a"},
				}},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		eksAPI.AssertExpectations(GinkgoT())
	})
	It("associates with all providers and waits", func() {
		manager := identityproviders.NewManager(api.ClusterMeta{
			Name: "idp-test",
		}, eksAPI)
		updateInput := &eks.DescribeUpdateInput{
			UpdateId: aws.String("1"),
			Name:     aws.String("idp-test"),
		}
		updateOutput := &eks.DescribeUpdateOutput{
			Update: &ekstypes.Update{
				Status: ekstypes.UpdateStatusSuccessful,
				Type:   ekstypes.UpdateTypeAssociateIdentityProviderConfig,
			},
		}
		eksAPI.On("DescribeUpdate", mock.Anything, updateInput, mock.Anything).Return(updateOutput, nil)
		err := manager.Associate(context.Background(), identityproviders.AssociateIdentityProvidersOptions{
			WaitTimeout: 1 * time.Minute,
			Providers: []api.IdentityProvider{
				{Inner: &api.OIDCIdentityProvider{
					Name:           "pool-1",
					IssuerURL:      "url",
					ClientID:       "id",
					UsernameClaim:  "usernameClaim",
					UsernamePrefix: "usernamePrefix",
					GroupsClaim:    "groupsClaim",
					GroupsPrefix:   "groupsPrefix",
					RequiredClaims: map[string]string{"permission": "true"},
					Tags:           map[string]string{"department": "a"},
				}},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		eksAPI.AssertExpectations(GinkgoT())
	})
})
