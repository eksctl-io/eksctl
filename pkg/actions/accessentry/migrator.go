package accessentry

import (
	"context"
	"time"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type AddonCreator interface {
	Create(ctx context.Context, addon *api.Addon, waitTimeout time.Duration) error
}

type AccessEntryMigrationOptions struct {
	RemoveOIDCProviderTrustRelationship bool
	TargetAuthMode                      string
	Approve                             bool
	Timeout                             time.Duration
}
