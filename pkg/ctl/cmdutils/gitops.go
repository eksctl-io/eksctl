package cmdutils

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// GitOpsConfigLoader handles loading of ClusterConfigFile v.s. using CLI
// flags for GitOps-related commands.
type GitOpsConfigLoader struct {
	cmd                    *Cmd
	validateWithConfigFile func() error
}

// NewGitOpsConfigLoader creates a new ClusterConfigLoader which handles
// loading of ClusterConfigFile GitOps-related commands.
func NewGitOpsConfigLoader(cmd *Cmd) *GitOpsConfigLoader {
	l := &GitOpsConfigLoader{
		cmd: cmd,
	}

	l.validateWithConfigFile = func() error {
		meta := l.cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet("metadata.name")
		}

		if meta.Region == "" {
			return ErrMustBeSet("metadata.region")
		}

		if l.cmd.ClusterConfig.GitOps == nil || l.cmd.ClusterConfig.GitOps.Flux == nil {
			return ErrMustBeSet("gitops.flux")
		}

		fluxCfg := l.cmd.ClusterConfig.GitOps.Flux
		if fluxCfg.GitProvider == "" {
			return ErrMustBeSet("gitops.flux.gitProvider")
		}

		if len(fluxCfg.Flags) == 0 {
			return ErrMustBeSet("gitops.flux.flags")
		}

		return nil
	}

	return l
}

// Load ClusterConfig or use CLI flags.
func (l *GitOpsConfigLoader) Load() error {
	if err := api.Register(); err != nil {
		return err
	}

	if l.cmd.ClusterConfigFile == "" {
		return ErrMustBeSet("--config-file/-f <file>")
	}

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
	var err error
	if l.cmd.ClusterConfig, err = eks.LoadConfigFromFile(l.cmd.ClusterConfigFile); err != nil {
		return err
	}

	meta := l.cmd.ClusterConfig.Metadata
	if meta == nil {
		return ErrMustBeSet("metadata")
	}

	if meta.Region != "" {
		l.cmd.ProviderConfig.Region = meta.Region
	}

	return l.validateWithConfigFile()
}
