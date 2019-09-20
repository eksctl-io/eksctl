package cmdutils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type gitopsConfigLoader struct {
	cmd                                *Cmd
	flagsIncompatibleWithConfigFile    sets.String
	flagsIncompatibleWithoutConfigFile sets.String
	validateWithConfigFile             func() error
	validateWithoutConfigFile          func() error
}

// Load ClusterConfig or use flags
func (l *gitopsConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.cmd.ClusterConfigFile == "" {
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.validateWithoutConfigFile()
	}

	var err error

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
	if l.cmd.ClusterConfig, err = eks.LoadConfigFromFile(l.cmd.ClusterConfigFile); err != nil {
		return err
	}
	meta := l.cmd.ClusterConfig.Metadata

	if meta == nil {
		return ErrMustBeSet("metadata")
	}

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	if meta.Name == "" {
		return ErrMustBeSet("metadata.name")
	}

	if meta.Region == "" {
		return ErrMustBeSet("metadata.region")
	}
	l.cmd.ProviderConfig.Region = meta.Region

	return l.validateWithConfigFile()
}

// NewEnableProfileLoader handles loading of clusterConfigFile vs using flags for gitops commands
func NewGitopsConfigLoader(cmd *Cmd) ClusterConfigLoader {
	l := newGitopsConfigLoader(cmd)

	l.validateWithoutConfigFile = func() error {
		meta := l.cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet("--cluster")
		}
		return nil
	}

	return l
}

func newGitopsConfigLoader(cmd *Cmd) *gitopsConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &gitopsConfigLoader{
		cmd: cmd,

		validateWithConfigFile: nilValidatorFunc,
		flagsIncompatibleWithConfigFile: sets.NewString(
			"region",
			"version",
			"cluster",
		),
		validateWithoutConfigFile:          nilValidatorFunc,
		flagsIncompatibleWithoutConfigFile: sets.NewString(),
	}
}
