package identityproviders

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/weaveworks/eksctl/pkg/utils/waiters"
	"github.com/weaveworks/logger"
)

type DisassociateIdentityProvidersOptions struct {
	Providers   []DisassociateIdentityProvider
	WaitTimeout *time.Duration
}

type DisassociateIdentityProvider struct {
	Name string
	Type string
}

func (ipm *IdentityProviderManager) Disassociate(options DisassociateIdentityProvidersOptions) error {
	clusterName := ipm.metadata.Name
	for _, idP := range options.Providers {
		idPConfig := eks.IdentityProviderConfig{
			Name: aws.String(idP.Name),
			Type: aws.String(idP.Type),
		}
		describeInput := eks.DescribeIdentityProviderConfigInput{
			ClusterName:            aws.String(ipm.metadata.Name),
			IdentityProviderConfig: &idPConfig,
		}
		idPDescription, err := ipm.eksAPI.DescribeIdentityProviderConfig(&describeInput)
		if err != nil {
			return err
		}
		if aws.StringValue(idPDescription.IdentityProviderConfig.Oidc.Status) == "DELETING" {
			logger.Warning("provider already deleting")
			return nil
		}

		disassociateInput := eks.DisassociateIdentityProviderConfigInput{
			ClusterName:            aws.String(clusterName),
			IdentityProviderConfig: &idPConfig,
		}

		disassociation, err := ipm.eksAPI.DisassociateIdentityProviderConfig(&disassociateInput)
		if err != nil {
			return err
		}

		if options.WaitTimeout != nil {
			newRequest := func() *request.Request {
				input := &eks.DescribeUpdateInput{
					Name:     aws.String(ipm.metadata.Name),
					UpdateId: disassociation.Update.Id,
				}
				req, _ := ipm.eksAPI.DescribeUpdateRequest(input)
				return req
			}

			acceptors := waiters.MakeAcceptors(
				"Update.Status",
				eks.UpdateStatusSuccessful,
				[]string{
					eks.UpdateStatusCancelled,
					eks.UpdateStatusFailed,
				},
			)

			msg := fmt.Sprintf(
				"waiting for requested identity provider %q in cluster %q to succeed",
				*disassociation.Update.Type,
				clusterName,
			)

			return waiters.Wait(clusterName, msg, acceptors, newRequest, *options.WaitTimeout, nil)
		}

		logger.Info("started disassociating identity provider %s", idP.Name)
	}
	return nil
}
