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

// ClusterConfigLoader holds common parameters required for loading
// ClusterConfig objects from files vs using flags fallback
type ClusterConfigLoader struct {
	provider *api.ProviderConfig
	path     string
	cmd      *cobra.Command

	Spec    *api.ClusterConfig
	NameArg string

	FlagsIncompatibleWithConfigFile, FlagsIncompatibleWithoutConfigFile sets.String

	ValidateWithConfigFile, ValidateWithoutConfigFile func() error
}

// NewClusterConfigLoader constructs a standard loader
func NewClusterConfigLoader(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile string, cmd *cobra.Command) *ClusterConfigLoader {
	nilValidator := func() error { return nil }

	return &ClusterConfigLoader{
		provider: provider,
		path:     clusterConfigFile,
		cmd:      cmd,

		Spec: spec,

		FlagsIncompatibleWithConfigFile:    sets.NewString("name", "region", "version"),
		ValidateWithConfigFile:             nilValidator,
		FlagsIncompatibleWithoutConfigFile: sets.NewString(),
		ValidateWithoutConfigFile:          nilValidator,
	}
}

// Load ClusterConfig or use flags
func (l *ClusterConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.path == "" {
		for f := range l.FlagsIncompatibleWithoutConfigFile {
			if l.cmd.Flag(f).Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.ValidateWithoutConfigFile()
	}

	if err := eks.LoadConfigFromFile(l.path, l.Spec); err != nil {
		return err
	}
	meta := l.Spec.Metadata

	for f := range l.FlagsIncompatibleWithConfigFile {
		if flag := l.cmd.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	if l.NameArg != "" {
		return ErrCannotUseWithConfigFile(fmt.Sprintf("name argument %q", l.NameArg))
	}

	if meta.Name == "" {
		return ErrMustBeSet("metadata.name")
	}

	if meta.Region == "" {
		return ErrMustBeSet("metadata.region")
	}
	l.provider.Region = meta.Region

	return l.ValidateWithConfigFile()
}

// LoadMetadata handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fileds, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func LoadMetadata(provider *api.ProviderConfig, spec *api.ClusterConfig, clusterConfigFile, nameArg string, cmd *cobra.Command) error {
	l := NewClusterConfigLoader(provider, spec, clusterConfigFile, cmd)

	l.NameArg = nameArg

	l.ValidateWithoutConfigFile = func() error {
		meta := l.Spec.Metadata

		if meta.Name != "" && l.NameArg != "" {
			return ErrNameFlagAndArg(meta.Name, l.NameArg)
		}

		if l.NameArg != "" {
			meta.Name = l.NameArg
		}

		if meta.Name == "" {
			return ErrMustBeSet("--name")
		}

		return nil
	}

	return l.Load()
}
