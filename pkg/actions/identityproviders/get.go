package identityproviders

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// Summary holds the known info about this provider
type Summary struct {
	Type           api.IdentityProviderType
	Name           string
	ClientID       string
	IssuerURL      string
	Status         string
	Arn            string
	UsernameClaim  *string
	UsernamePrefix *string
	GroupsClaim    *string
	GroupsPrefix   *string
	RequiredClaims map[string]string
	Tags           map[string]string
}

type GetIdentityProvidersOptions struct {
	Name string
}

func (m *Manager) Get(ctx context.Context, options GetIdentityProvidersOptions) ([]Summary, error) {
	var summaries []Summary
	var configs []ekstypes.IdentityProviderConfig

	input := eks.ListIdentityProviderConfigsInput{
		ClusterName: aws.String(m.metadata.Name),
	}
	list, err := m.eksAPI.ListIdentityProviderConfigs(ctx, &input)
	if err != nil {
		return summaries, err
	}

	if options.Name == "" {
		configs = list.IdentityProviderConfigs
	} else {
		var getCfg *ekstypes.IdentityProviderConfig
		for _, cfg := range list.IdentityProviderConfigs {
			if aws.ToString(cfg.Name) == options.Name {
				getCfg = &ekstypes.IdentityProviderConfig{
					Name: aws.String(options.Name),
					Type: cfg.Type,
				}
			}
		}
		if getCfg == nil {
			return summaries, errors.Errorf("couldn't find identity provider %s", options.Name)
		}
		configs = []ekstypes.IdentityProviderConfig{*getCfg}
	}
	for _, idp := range configs {
		input := eks.DescribeIdentityProviderConfigInput{
			ClusterName:            aws.String(m.metadata.Name),
			IdentityProviderConfig: &idp,
		}
		idP, err := m.eksAPI.DescribeIdentityProviderConfig(ctx, &input)
		if err != nil {
			return summaries, err
		}
		if cfg := idP.IdentityProviderConfig.Oidc; cfg != nil {
			summaries = append(summaries, Summary{
				Type:           api.OIDCIdentityProviderType,
				Name:           aws.ToString(cfg.IdentityProviderConfigName),
				ClientID:       aws.ToString(cfg.ClientId),
				IssuerURL:      aws.ToString(cfg.IssuerUrl),
				Status:         string(cfg.Status),
				Arn:            aws.ToString(cfg.IdentityProviderConfigArn),
				UsernameClaim:  cfg.UsernameClaim,
				UsernamePrefix: cfg.UsernamePrefix,
				GroupsClaim:    cfg.GroupsClaim,
				GroupsPrefix:   cfg.GroupsPrefix,
				RequiredClaims: cfg.RequiredClaims,
				Tags:           cfg.Tags,
			})
		}
	}
	return summaries, nil
}
