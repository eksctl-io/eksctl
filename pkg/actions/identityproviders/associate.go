package identityproviders

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/weaveworks/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type AssociateIdentityProvidersOptions struct {
	Providers   []api.IdentityProvider
	WaitTimeout *time.Duration
}

func (m *Manager) Associate(options AssociateIdentityProvidersOptions) error {
	taskTree := tasks.TaskTree{
		Parallel: true,
	}

	for _, generalIDP := range options.Providers {
		switch idP := (generalIDP.Inner).(type) {
		case *api.OIDCIdentityProvider:
			taskTree.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("associate %s", idP.Name),
				Doer: func() error {
					update, err := m.associateOIDC(*idP)
					if err != nil {
						return err
					}

					logger.Info("started associating identity provider %s", idP.Name)
					if options.WaitTimeout != nil {
						if err := m.waitForUpdate(update, *options.WaitTimeout); err != nil {
							return err
						}
					}
					return nil
				},
			})

		default:
			panic("unsupported identity provider")
		}
	}

	errs := taskTree.DoAllSync()
	for _, err := range errs {
		logger.Critical(err.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf("one or more providers failed to associate")
	}
	return nil
}

func (m *Manager) associateOIDC(idP api.OIDCIdentityProvider) (eks.Update, error) {
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
		ClusterName: aws.String(m.metadata.Name),
		Oidc:        oidc,
	}
	if len(idP.Tags) > 0 {
		input.Tags = aws.StringMap(idP.Tags)
	}

	update, err := m.eksAPI.AssociateIdentityProviderConfig(&input)
	if err != nil {
		return eks.Update{}, err
	}
	logger.Debug("identity provider associate update: %v", *update.Update)

	return *update.Update, nil
}
