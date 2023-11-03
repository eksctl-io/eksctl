package cmdutils

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewCreateAccessEntryLoader creates a new loader for access entries.
func NewCreateAccessEntryLoader(cmd *Cmd, accessEntry api.AccessEntry) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	const principalARNFlag = "principal-arn"
	l.flagsIncompatibleWithConfigFile = sets.NewString(
		principalARNFlag,
		"kubernetes-groups",
		"kubernetes-username",
	)

	l.validateWithConfigFile = func() error {
		if len(cmd.ClusterConfig.AccessConfig.AccessEntries) == 0 {
			return errors.New("at least one access entry is required")
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if accessEntry.PrincipalARN.Partition == "" {
			return fmt.Errorf("--%s is required", principalARNFlag)
		}
		l.ClusterConfig.AccessConfig.AccessEntries = []api.AccessEntry{accessEntry}
		return nil
	}

	return l
}
