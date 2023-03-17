package identityproviders

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/waiter"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type DisassociateIdentityProvidersOptions struct {
	Providers   []DisassociateIdentityProvider
	WaitTimeout time.Duration
}

type DisassociateIdentityProvider struct {
	Name string
	Type api.IdentityProviderType
}

func (m *Manager) Disassociate(ctx context.Context, options DisassociateIdentityProvidersOptions) error {
	taskTree := tasks.TaskTree{
		Parallel: true,
	}

	for _, idP := range options.Providers {
		taskTree.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("disassociate %s", idP.Name),
			Doer: func() error {
				idPConfig := ekstypes.IdentityProviderConfig{
					Name: aws.String(idP.Name),
					Type: aws.String(string(idP.Type)),
				}
				describeInput := eks.DescribeIdentityProviderConfigInput{
					ClusterName:            aws.String(m.metadata.Name),
					IdentityProviderConfig: &idPConfig,
				}
				idPDescription, err := m.eksAPI.DescribeIdentityProviderConfig(ctx, &describeInput)
				if err != nil {
					return err
				}
				if idPDescription.IdentityProviderConfig.Oidc.Status == ekstypes.ConfigStatusDeleting {
					logger.Warning("provider already deleting")
					return nil
				}

				disassociateInput := eks.DisassociateIdentityProviderConfigInput{
					ClusterName:            aws.String(m.metadata.Name),
					IdentityProviderConfig: &idPConfig,
				}

				update, err := m.eksAPI.DisassociateIdentityProviderConfig(ctx, &disassociateInput)
				if err != nil {
					return err
				}
				logger.Debug("identity provider disassociate update: %v", *update.Update)

				logger.Info("started disassociating identity provider %s", idP.Name)

				if options.WaitTimeout > 0 {
					updateWaiter := waiter.NewUpdateWaiter(m.eksAPI, func(options *waiter.UpdateWaiterOptions) {
						options.RetryAttemptLogMessage = fmt.Sprintf("waiting for update %q in cluster %q to succeed", *update.Update.Id, m.metadata.Name)
					})
					return updateWaiter.Wait(ctx, &eks.DescribeUpdateInput{
						Name:     aws.String(m.metadata.Name),
						UpdateId: update.Update.Id,
					}, options.WaitTimeout)
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
