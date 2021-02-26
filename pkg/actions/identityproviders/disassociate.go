package identityproviders

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type DisassociateIdentityProvidersOptions struct {
	Providers   []DisassociateIdentityProvider
	WaitTimeout *time.Duration
}

type DisassociateIdentityProvider struct {
	Name string
	Type api.IdentityProviderType
}

func (m *Manager) Disassociate(options DisassociateIdentityProvidersOptions) error {
	taskTree := tasks.TaskTree{
		Parallel: true,
	}

	for _, idP := range options.Providers {
		taskTree.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("disassociate %s", idP.Name),
			Doer: func() error {
				idPConfig := eks.IdentityProviderConfig{
					Name: aws.String(idP.Name),
					Type: aws.String(string(idP.Type)),
				}
				describeInput := eks.DescribeIdentityProviderConfigInput{
					ClusterName:            aws.String(m.metadata.Name),
					IdentityProviderConfig: &idPConfig,
				}
				idPDescription, err := m.eksAPI.DescribeIdentityProviderConfig(&describeInput)
				if err != nil {
					return err
				}
				if aws.StringValue(idPDescription.IdentityProviderConfig.Oidc.Status) == "DELETING" {
					logger.Warning("provider already deleting")
					return nil
				}

				disassociateInput := eks.DisassociateIdentityProviderConfigInput{
					ClusterName:            aws.String(m.metadata.Name),
					IdentityProviderConfig: &idPConfig,
				}

				update, err := m.eksAPI.DisassociateIdentityProviderConfig(&disassociateInput)
				if err != nil {
					return err
				}
				logger.Debug("identity provider disassociate update: %v", *update.Update)

				logger.Info("started disassociating identity provider %s", idP.Name)

				if options.WaitTimeout != nil {
					if err := m.waitForUpdate(*update.Update, *options.WaitTimeout); err != nil {
						return err
					}
				}
				return nil
			},
		})
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
