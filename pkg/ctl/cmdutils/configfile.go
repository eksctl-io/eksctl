package cmdutils

import (
	"fmt"

	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// LoadMetadata handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fileds, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func LoadMetadata(p *api.ProviderConfig, cfg *api.ClusterConfig, clusterConfigFile, nameArg string, cmd *cobra.Command) error {
	meta := cfg.Metadata

	if clusterConfigFile != "" {
		if err := eks.LoadConfigFromFile(clusterConfigFile, cfg); err != nil {
			return err
		}
		meta = cfg.Metadata

		incompatibleFlags := []string{
			"name",
			"region",
			"version",
		}

		for _, f := range incompatibleFlags {
			if flag := cmd.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s when --config-file/-f is set", f)
			}
		}

		if nameArg != "" {
			return fmt.Errorf("cannot use name argument %q when --config-file/-f is set", nameArg)
		}

		if meta.Name == "" {
			return fmt.Errorf("metadata.name must be set")
		}

		// region is always required when config file is used
		if meta.Region == "" {
			return fmt.Errorf("metadata.region must be set")
		}

		p.Region = meta.Region

		// version has different default values in some command
		// we don't check it here
	} else {
		if meta.Name != "" && nameArg != "" {
			return ErrNameFlagAndArg(meta.Name, nameArg)
		}

		if nameArg != "" {
			meta.Name = nameArg
		}

		if meta.Name == "" {
			return fmt.Errorf("--name must be set")
		}

		// default region will get picked by eks.New, and
		// version validation gets handled separately
	}

	return nil
}
