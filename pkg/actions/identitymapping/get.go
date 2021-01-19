package identitymapping

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/iam"
)

func (m *Manager) Get(arn string) ([]iam.Identity, error) {
	acm, err := authconfigmap.NewFromClientSet(m.clientSet)
	if err != nil {
		return nil, err
	}
	identities, err := acm.Identities()
	if err != nil {
		return nil, err
	}

	if arn != "" {
		selectedIdentities := []iam.Identity{}

		for _, identity := range identities {
			if identity.ARN() == arn {
				selectedIdentities = append(selectedIdentities, identity)
			}
		}

		identities = selectedIdentities
		// If a filter was given, we error if none was found
		if len(identities) == 0 {
			return nil, fmt.Errorf("no iamidentitymapping with arn %q found", arn)
		}
	}

	return identities, nil
}
