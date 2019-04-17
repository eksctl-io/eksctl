package delete

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/vpc"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func deleteClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doDeleteCluster(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddWaitFlag(&wait, fs, "deletion of all resources")
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)
	return cmd
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}

func deleteDeprecatedStacks(stackManager *manager.StackCollection) (bool, error) {
	tasks, err := stackManager.DeleteTasksForDeprecatedStacks()
	if err != nil {
		return true, err
	}
	if count := tasks.Len(); count > 0 {
		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return true, handleErrors(errs, "deprecated stacks")
		}
		logger.Success("deleted all %s deperecated stacks", count)
		return true, nil
	}
	return false, nil
}

func doDeleteCluster(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
	printer := printers.NewJSONPrinter()

	if err := cmdutils.NewMetadataLoader(p, cfg, clusterConfigFile, nameArg, cmd).Load(); err != nil {
		return err
	}

	ctl := eks.New(p, cfg)
	meta := cfg.Metadata

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	logger.Info("deleting EKS cluster %q", meta.Name)
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	ctl.MaybeDeletePublicSSHKey(meta.Name)

	kubeconfig.MaybeDeleteConfig(meta)

	if hasDeprectatedStacks, err := deleteDeprecatedStacks(stackManager); hasDeprectatedStacks {
		if err != nil {
			return err
		}
		return nil
	}

	{
		tasks, err := stackManager.NewTasksToDeleteClusterWithNodeGroups(wait, func(errs chan error, _ string) error {
			logger.Info("trying to cleanup dangling network interfaces")
			if err := ctl.GetClusterVPC(cfg); err != nil {
				return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
			}
			go func() {
				errs <- vpc.CleanupNetworkInterfaces(ctl.Provider, cfg)
				close(errs)
			}()
			return nil
		})

		if err != nil {
			return err
		}

		if tasks.Len() == 0 {
			logger.Warning("no cluster resources were found for %q", meta.Name)
			return nil
		}

		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return handleErrors(errs, "cluster with nodegroup(s)")
		}

		logger.Success("all cluster resources were deleted")
	}

	return nil
}
