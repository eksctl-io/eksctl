package cmdutils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// AddConfigFileFlag adds common --config-file flag
func AddConfigFileFlag(path *string, fs *pflag.FlagSet) {
	fs.StringVarP(path, "config-file", "f", "", "load configuration from a file")
}

// ClusterConfigLoader is an inteface that loaders should implement
type ClusterConfigLoader interface {
	Load() error
}

type commonClusterConfigLoader struct {
	provider *api.ProviderConfig
	path     string
	cmd      *cobra.Command
	spec     *api.ClusterConfig
	nameArg  string

	flagsIncompatibleWithConfigFile, flagsIncompatibleWithoutConfigFile sets.String

	validateWithConfigFile, validateWithoutConfigFile func() error
}

func newCommonClusterConfigLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile string, cmd *cobra.Command) *commonClusterConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &commonClusterConfigLoader{
		provider: provider,
		path:     clusterConfigFile,
		cmd:      cmd,
		spec:     spec,

		validateWithConfigFile:             nilValidatorFunc,
		flagsIncompatibleWithConfigFile:    sets.NewString("name", "region", "version"),
		validateWithoutConfigFile:          nilValidatorFunc,
		flagsIncompatibleWithoutConfigFile: sets.NewString(),
	}
}

// Load ClusterConfig or use flags
func (l *commonClusterConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.path == "" {
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if l.cmd.Flag(f).Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.validateWithoutConfigFile()
	}

	if err := eks.LoadConfigFromFile(l.path, l.spec); err != nil {
		return err
	}
	meta := l.spec.Metadata

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.cmd.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	if l.nameArg != "" {
		return ErrCannotUseWithConfigFile(fmt.Sprintf("name argument %q", l.nameArg))
	}

	if meta.Name == "" {
		return ErrMustBeSet("metadata.name")
	}

	if meta.Region == "" {
		return ErrMustBeSet("metadata.region")
	}
	l.provider.Region = meta.Region

	return l.validateWithConfigFile()
}

// NewMetadataLoader handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fileds, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func NewMetadataLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile, nameArg string, cmd *cobra.Command) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(provider, spec, clusterConfigFile, cmd)

	l.nameArg = nameArg

	l.validateWithoutConfigFile = func() error {
		meta := l.spec.Metadata

		if meta.Name != "" && l.nameArg != "" {
			return ErrNameFlagAndArg(meta.Name, l.nameArg)
		}

		if l.nameArg != "" {
			meta.Name = l.nameArg
		}

		if meta.Name == "" {
			return ErrMustBeSet("--name")
		}

		return nil
	}

	return l
}
