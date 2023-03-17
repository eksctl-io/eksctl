package iamidentitymapping

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/iam"
)

func (m *Manager) Create(ctx context.Context, mapping *api.IAMIdentityMapping) error {

	acm, err := authconfigmap.NewFromClientSet(m.clientSet)
	if err != nil {
		return err
	}

	if mapping.ServiceName != "" {
		parsedARN, err := arn.Parse(m.clusterConfig.Status.ARN)
		if err != nil {
			return errors.Wrap(err, "error parsing cluster ARN")
		}
		sa := authconfigmap.NewServiceAccess(m.rawClient, acm, parsedARN.AccountID)
		return sa.Grant(mapping.ServiceName, mapping.Namespace, api.Partition(m.region))
	}

	if mapping.Account == "" {
		id, err := iam.NewIdentity(mapping.ARN, mapping.Username, mapping.Groups)
		if err != nil {
			return err
		}

		// Check whether role already exists.
		identities, err := acm.GetIdentities()
		if err != nil {
			return err
		}

		createdArn := id.ARN() // The call to Valid above makes sure this cannot error
		logger.Info("checking arn %s against entries in the auth ConfigMap", id.ARN())
		for _, identity := range identities {
			arn := identity.ARN()
			if mapping.NoDuplicateARNs && iam.CompareIdentity(id, identity) {
				logger.Warning("found existing mapping that matches the one being created, skipping.")
				return nil
			}

			if createdArn == arn && mapping.NoDuplicateARNs {
				return fmt.Errorf("found existing mapping with the same arn %q and shadowing is disabled", createdArn)
			}

			if createdArn == arn {
				logger.Warning("found existing mappings with same arn %q (which will be shadowed by your new mapping)", createdArn)
				break
			}
		}

		if err := acm.AddIdentity(id); err != nil {
			return err
		}
	} else {
		if err := acm.AddAccount(mapping.Account); err != nil {
			return err
		}
	}
	return acm.Save()

}
