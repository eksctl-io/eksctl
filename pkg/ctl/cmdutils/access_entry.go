package cmdutils

import (
	"errors"
	"fmt"

	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/sets"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const principalARNFlag = "principal-arn"

var (
	accessEntryFlagsIncompatibleWithoutConfigFile = []string{}
	accessEntryFlagsIncompatibleWithConfigFile    = []string{"principal-arn"}
)

// NewCreateAccessEntryLoader creates a new loader for access entries.
func NewCreateAccessEntryLoader(cmd *Cmd, accessEntry *api.AccessEntry) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

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
		l.ClusterConfig.AccessConfig.AccessEntries = []api.AccessEntry{*accessEntry}
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
func NewDeleteAccessEntryLoader(cmd *Cmd, accessEntry api.AccessEntry) ClusterConfigLoader {
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
		if accessEntry.PrincipalARN.IsZero() {
			return ErrMustBeSet(fmt.Sprintf("--%s", principalARNFlag))
		}
		l.ClusterConfig.AccessConfig.AccessEntries = []api.AccessEntry{accessEntry}
		return nil
	}

	return l
}

// NewUtilsUpdateAuthenticationModeLoader loads config or uses flags for `eksctl utils update-autentication-mode`
func NewUtilsUpdateAuthenticationModeLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"cluster",
		"authentication-mode",
	)

	validateAuthenticationMode := func(authenticationMode ekstypes.AuthenticationMode) error {
		if authenticationMode == "" {
			return ErrMustBeSet("--authentication-mode")
		}
		if !slices.Contains(authenticationMode.Values(), cmd.ClusterConfig.AccessConfig.AuthenticationMode) {
			return fmt.Errorf("invalid value %q provided for authenticationMode, choose one of: %q, %q, %q",
				authenticationMode, ekstypes.AuthenticationModeConfigMap, ekstypes.AuthenticationModeApiAndConfigMap, ekstypes.AuthenticationModeApi)
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		meta := cmd.ClusterConfig.Metadata
		authenticationMode := cmd.ClusterConfig.AccessConfig.AuthenticationMode
		if meta.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if cmd.NameArg != "" {
			return ErrUnsupportedNameArg()
		}
		return validateAuthenticationMode(authenticationMode)
	}

	l.validateWithConfigFile = func() error {
		meta := cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		return validateAuthenticationMode(cmd.ClusterConfig.AccessConfig.AuthenticationMode)
	}

	return l
}
