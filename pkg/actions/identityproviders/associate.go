package identityproviders

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type AssociateIdentityProvidersOptions struct {
	Providers []api.IdentityProvider
}

func (ipm *IdentityProviderManager) associateOIDC(idP api.OIDCIdentityProvider) error {
	oidc := &eks.OidcIdentityProviderConfigRequest{
		ClientId:                   aws.String(idP.ClientID),
		IssuerUrl:                  aws.String(idP.IssuerURL),
		IdentityProviderConfigName: aws.String(idP.Name),
	}
	if idP.GroupsClaim != "" {
		oidc.GroupsClaim = aws.String(idP.GroupsClaim)
	}
	if idP.GroupsPrefix != "" {
		oidc.GroupsPrefix = aws.String(idP.GroupsPrefix)
	}
	if len(idP.RequiredClaims) > 0 {
		oidc.RequiredClaims = aws.StringMap(idP.RequiredClaims)
	}
	if idP.UsernameClaim != "" {
		oidc.UsernameClaim = aws.String(idP.UsernameClaim)
	}
	if idP.UsernamePrefix != "" {
		oidc.UsernamePrefix = aws.String(idP.UsernamePrefix)
	}
	input := eks.AssociateIdentityProviderConfigInput{
		ClusterName: aws.String(ipm.metadata.Name),
		Oidc:        oidc,
	}
	if len(idP.Tags) > 0 {
		input.Tags = aws.StringMap(idP.Tags)
	}
	_, err := ipm.eksAPI.AssociateIdentityProviderConfig(&input)
	if err != nil {
		return err
	}
	return nil
}

func (ipm *IdentityProviderManager) Associate(options AssociateIdentityProvidersOptions) error {
	for _, generalIdP := range options.Providers {
		switch idP := generalIdP.Inner().(type) {
		case *api.OIDCIdentityProvider:
			if err := ipm.associateOIDC(*idP); err != nil {
				return err
			}
		default:
			continue
		}
	}
	return nil
}
