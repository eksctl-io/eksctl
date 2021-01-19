package identitymapping

import (
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
)

func (m *Manager) Delete(identityMappings []*api.IAMIdentityMapping, all bool) error {
	acm, err := authconfigmap.NewFromClientSet(m.clientSet)
	if err != nil {
		return err
	}

	for _, identityMapping := range identityMappings {
		if err := acm.RemoveIdentity(identityMapping.ARN, all); err != nil {
			return err
		}
		if err := acm.Save(); err != nil {
			return err
		}

		// Check whether we have more roles that match
		identities, err := acm.Identities()
		if err != nil {
			return err
		}

		duplicates := 0
		for _, identity := range identities {

			if identityMapping.ARN == identity.ARN() {
				duplicates++
			}
		}

		if duplicates > 0 {
			logger.Warning("there are %d mappings left with same arn %q (use --all to delete them at once)", duplicates, identityMapping.ARN)
		}
	}
	return nil
}
