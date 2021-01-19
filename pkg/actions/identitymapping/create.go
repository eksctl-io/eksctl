package identitymapping

import (
	"fmt"

	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/iam"
)

func (m *Manager) Create(identityMappings []*api.IAMIdentityMapping, clusterARN string) error {
	var failedIdentityMappings []string
	for _, identityMapping := range identityMappings {
		err := m.create(identityMapping, clusterARN)
		if err != nil {
			logger.Warning("failed to create identity mapping %s: %v", identityMapping.ARN, err)
			failedIdentityMappings = append(failedIdentityMappings, identityMapping.ARN)
		}
	}
	if len(failedIdentityMappings) > 0 {
		return fmt.Errorf("failed to create %d identity mappings: %v", len(failedIdentityMappings), failedIdentityMappings)
	}
	return nil
}

func (m *Manager) create(identityMapping *api.IAMIdentityMapping, arn string) error {
	hasARNOptions := func() bool {
		return !(identityMapping.ARN == "" && identityMapping.Username == "" && len(identityMapping.Groups) == 0)
	}

	validateNonServiceOptions := func() error {
		if identityMapping.Namespace != "" {
			return errors.New("namespace is only valid with service-name")
		}
		return nil
	}

	acm, err := authconfigmap.NewFromClientSet(m.clientSet)
	if err != nil {
		return err
	}

	if identityMapping.ServiceName != "" {
		if hasARNOptions() {
			return errors.New("cannot use arn, username, and groups with service-name")
		}

		parsedARN, err := awsarn.Parse(arn)
		if err != nil {
			return errors.Wrap(err, "error parsing cluster ARN")
		}
		sa := authconfigmap.NewServiceAccess(m.rawClient, acm, parsedARN.AccountID)
		return sa.Grant(identityMapping.ServiceName, identityMapping.Namespace)
	}

	// Check whether role already exists.
	if identityMapping.Account == "" {
		if err := validateNonServiceOptions(); err != nil {
			return err
		}
		id, err := iam.NewIdentity(identityMapping.ARN, identityMapping.Username, identityMapping.Groups)
		if err != nil {
			return err
		}
		identities, err := acm.Identities()
		if err != nil {
			return err
		}

		createdArn := id.ARN() // The call to Valid above makes sure this cannot error
		for _, identity := range identities {
			arn := identity.ARN()

			if createdArn == arn {
				logger.Warning("found existing mappings with same arn %q (which will be shadowed by your new mapping)", createdArn)
				break
			}
		}

		if err := acm.AddIdentity(id); err != nil {
			return err
		}
	} else if hasARNOptions() {
		if err := validateNonServiceOptions(); err != nil {
			return err
		}
		if err := acm.AddAccount(identityMapping.Account); err != nil {
			return err
		}
	} else {
		return errors.New("account can only be set alone")
	}

	return acm.Save()
}
