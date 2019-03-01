package update

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

var updateClusterDryRun = true

func updateClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Update cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doUpdateClusterCmd(p, cfg, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name (required)")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		fs.BoolVar(&updateClusterDryRun, "dry-run", updateClusterDryRun, "do not apply any change, only show what resources would be added")
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)
	return cmd
}

func doUpdateClusterCmd(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string) error {
	ctl := eks.New(p, cfg)

	printer := printers.NewJSONPrinter()

	if err := api.Register(); err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && nameArg != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.Metadata.Name, nameArg)
	}

	if nameArg != "" {
		cfg.Metadata.Name = nameArg
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	if err := ctl.GetClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n", cfg); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	updateRequired, err := stackManager.AppendNewClusterStackResource(updateClusterDryRun)
	if err != nil {
		return err
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	if updateClusterDryRun && updateRequired {
		logger.Warning("no changes were applied, run again with '--dry-run=false' to apply the changes")
	}

	return nil
}
