package cmdutils

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var (
	accessEntryFlagsIncompatibleWithoutConfigFile = []string{}
	accessEntryFlagsIncompatibleWithConfigFile    = []string{"principal-arn"}
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

// NewGetAccessEntriesLoader loads config file and validates command for `eksctl get accessentry`.
func NewGetAccessEntryLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(accessEntryFlagsIncompatibleWithConfigFile...)
	l.flagsIncompatibleWithoutConfigFile.Insert(accessEntryFlagsIncompatibleWithoutConfigFile...)

	l.validateWithoutConfigFile = func() error {
		meta := cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if cmd.NameArg != "" {
			return ErrUnsupportedNameArg()
		}
		return nil
	}

	return l
}

// NewDeleteAccessEntryLoader loads config file and validates command for `eksctl delete accessentry`.
func NewDeleteAccessEntryLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(accessEntryFlagsIncompatibleWithConfigFile...)
	l.flagsIncompatibleWithoutConfigFile.Insert(accessEntryFlagsIncompatibleWithoutConfigFile...)

	l.validateWithConfigFile = func() error {
		if len(cmd.ClusterConfig.AccessConfig.AccessEntries) == 0 {
			return fmt.Errorf("no access entries specified")
		}
		for _, ae := range cmd.ClusterConfig.AccessConfig.AccessEntries {
			if ae.PrincipalARN.IsZero() {
				return fmt.Errorf("must specify access entry principalArn")
			}
			for _, policy := range ae.AccessPolicies {
				if policy.PolicyARN.IsZero() {
					return fmt.Errorf("must specify access policy arn")
				}
			}
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if cmd.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if len(cmd.ClusterConfig.AccessConfig.AccessEntries) == 0 {
			return fmt.Errorf("no access entries specified")
		}
		if cmd.ClusterConfig.AccessConfig.AccessEntries[0].PrincipalARN.IsZero() {
			return fmt.Errorf("must specify access entry principalArn")
		}
		return nil
	}

	return l
}
