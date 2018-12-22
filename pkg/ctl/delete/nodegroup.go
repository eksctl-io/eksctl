package delete

import (
	"fmt"
	"os"

	"errors"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

func deleteNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "nodegroup NAME",
		Short: "Delete a nodegroup",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := cmdutils.GetNameArg(args)
			if name != "" {
				ng.Name = name
			}
			if err := doDeleteNodeGroup(p, cfg, ng.Name); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	group := &cmdutils.NamedFlagSetGroup{}

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name (required)")
		cmdutils.AddRegionFlag(fs, p)
		fs.BoolVarP(&waitDelete, "wait", "w", false, "Wait for deletion of all resources before exiting")
	})

	group.InFlagSet("Nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup. Generated if unset, e.g. \"ng-a345f4\"")
	})

	cmdutils.AddCommonFlagsForAWS(group, p)

	group.AddTo(cmd)

	return cmd
}

func doDeleteNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, name string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return errors.New("--cluster must be set")
	}

	logger.Info("deleting EKS nodegroup %q-nodegroup-%s", cfg.Metadata.Name, name)

	var deletedResources []string

	handleIfError := func(err error, name string) bool {
		if err != nil {
			logger.Debug("continue despite error: %v", err)
			return true
		}
		logger.Debug("deleted %q", name)
		deletedResources = append(deletedResources, name)
		return false
	}

	// We can remove all 'DeprecatedDelete*' calls in 0.2.0

	stackManager := ctl.NewStackManager(cfg)

	{
		err := stackManager.WaitDeleteNodeGroup(nil, name)
		errs := []error{err}
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred while deleting nodegroup(s)", len(errs))
			for _, err := range errs {
				if err != nil {
					logger.Critical("%s\n", err.Error())
				}
			}
			handleIfError(fmt.Errorf("failed to delete nodegroup(s)"), "nodegroup(s)")
		}
		logger.Debug("all nodegroups were deleted")
	}

	return nil
}
