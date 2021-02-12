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

var _ = Describe("Associate", func() {
	var eksAPI mocks.EKSAPI
	BeforeEach(func() {
		eksAPI = mocks.EKSAPI{}
		eksAPI.On("AssociateIdentityProviderConfig", &eks.AssociateIdentityProviderConfigInput{
			ClusterName: aws.String(""),
			Oidc: &eks.OidcIdentityProviderConfigRequest{
				IdentityProviderConfigName: aws.String("pool-1"),
				IssuerUrl:                  aws.String("url"),
				ClientId:                   aws.String("id"),
				UsernameClaim:              aws.String("usernameClaim"),
				UsernamePrefix:             aws.String("usernamePrefix"),
				GroupsClaim:                aws.String("groupsClaim"),
				GroupsPrefix:               aws.String("groupsPrefix"),
				RequiredClaims:             aws.StringMap(map[string]string{"permission": "true"}),
			},
			Tags: aws.StringMap(map[string]string{"department": "a"}),
		}).Return(&eks.AssociateIdentityProviderConfigOutput{
			Update: &eks.Update{
				Id:   aws.String("1"),
				Type: aws.String("AssociateIdentityProviderConfig"),
			},
		}, nil)
	})
	It("associates with all providers", func() {
		manager := identityproviders.NewManager(api.ClusterMeta{}, &eksAPI)
		err := manager.Associate(identityproviders.AssociateIdentityProvidersOptions{
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
		manager := identityproviders.NewManager(api.ClusterMeta{}, &eksAPI)
		client := mockprovider.NewMockAWSClient()
		updateInput := eks.DescribeUpdateInput{
			UpdateId: aws.String("1"),
			Name:     aws.String(""),
		}
		updateOutput := eks.DescribeUpdateOutput{
			Update: &eks.Update{
				Status: aws.String(eks.UpdateStatusSuccessful),
				Type:   aws.String("AssociateIdentityProviderConfig"),
			},
		}
		eksAPI.On("DescribeUpdateRequest", &updateInput).Return(
			client.MockRequestForGivenOutput(&updateInput, &updateOutput), &updateOutput,
		)
		wait := 1 * time.Minute
		err := manager.Associate(identityproviders.AssociateIdentityProvidersOptions{
			WaitTimeout: &wait,
			Providers: []api.IdentityProvider{
				api.IdentityProvider{Inner: &api.OIDCIdentityProvider{
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
