package delete

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func deleteNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	updateAuthConfigMap := true

	cmd := &cobra.Command{
		Use:     "nodegroup",
		Short:   "Delete a nodegroup",
		Aliases: []string{"ng"},
		Run: func(_ *cobra.Command, args []string) {
			if err := doDeleteNodeGroup(p, cfg, ng, cmdutils.GetNameArg(args), updateAuthConfigMap); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name (required)")
		cmdutils.AddRegionFlag(fs, p)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete (required)")
		cmdutils.AddWaitFlag(&wait, fs)
		fs.BoolVar(&updateAuthConfigMap, "update-auth-config-map", true, "Remove nodegroup IAM role from aws-auth config map")
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)

	return cmd
}

func doDeleteNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup, nameArg string, updateAuthConfigMap bool) error {
	ctl := eks.New(p, cfg)

	if err := api.Register(); err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return errors.New("--cluster must be set")
	}

	if ng.Name != "" && nameArg != "" {
		return cmdutils.ErrNameFlagAndArg(ng.Name, nameArg)
	}

	if nameArg != "" {
		ng.Name = nameArg
	}

	if ng.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	stackManager := ctl.NewStackManager(cfg)

	if ng.IAM.InstanceRoleARN == "" {
		if err := ctl.GetNodeGroupIAM(cfg, ng); err != nil {
			logger.Warning("%s getting instance role ARN for node group %q", err.Error(), ng.Name)
		}
	}

	logger.Info("deleting nodegroup %q in cluster %q", ng.Name, cfg.Metadata.Name)

	{
		var (
			err  error
			verb string
		)
		if wait {
			err = stackManager.BlockingWaitDeleteNodeGroup(ng.Name, false)
			verb = "was"
		} else {
			err = stackManager.DeleteNodeGroup(ng.Name)
			verb = "will be"
		}
		if err != nil {
			return errors.Wrapf(err, "failed to delete nodegroup %q", ng.Name)
		}
		logger.Success("nodegroup %q %s deleted", ng.Name, verb)
	}

	// post-deletion action
	if updateAuthConfigMap {
		clientSet, err := ctl.NewStdClientSet(cfg)
		if err != nil {
			return err
		}

		// remove node group from config map
		if err := authconfigmap.RemoveNodeGroup(clientSet, ng); err != nil {
			logger.Warning(err.Error())
		}
	}

	return nil
}
