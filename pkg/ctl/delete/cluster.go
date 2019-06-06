package delete

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/ssh"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

func deleteClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doDeleteCluster(cp); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
		cmdutils.AddWaitFlag(fs, &cp.Wait, "deletion of all resources")
		cmdutils.AddConfigFileFlag(fs, &cp.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, true)

	group.AddTo(cp.Command)
	return cp.Command
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

func doDeleteCluster(cp *cmdutils.CommonParams) error {
	if err := cmdutils.NewMetadataLoader(cp).Load(); err != nil {
		return err
	}

	cfg := cp.ClusterConfig
	meta := cp.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()
	ctl := eks.New(cp.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(cp.ProviderConfig)
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

	ssh.DeleteKeys(meta.Name, ctl.Provider)

	kubeconfig.MaybeDeleteConfig(meta)

	if hasDeprectatedStacks, err := deleteDeprecatedStacks(stackManager); hasDeprectatedStacks {
		if err != nil {
			return err
		}
		return nil
	}

	{
		tasks, err := stackManager.NewTasksToDeleteClusterWithNodeGroups(cp.Wait, func(errs chan error, _ string) error {
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
