package delete

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks"
)

var (
	updateAuthConfigMap  bool
	deleteNodeGroupDrain bool
)

func deleteNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:     "nodegroup",
		Short:   "Delete a nodegroup",
		Aliases: []string{"ng"},
		Run: func(_ *cobra.Command, args []string) {
			if err := doDeleteNodeGroup(p, cfg, ng, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, p)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete (required)")
		cmdutils.AddWaitFlag(&wait, fs, "deletion of all resources")
		cmdutils.AddUpdateAuthConfigMap(&updateAuthConfigMap, fs, "Remove nodegroup IAM role from aws-auth configmap")
		fs.BoolVar(&deleteNodeGroupDrain, "drain", true, "Drain and cordon all nodes in the nodegroup before deletion")
		cmdutils.AddCommonFlagsForDeleteCmd(fs, &output)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)

	return cmd
}

func doDeleteNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup, nameArg string) error {
	ctl := eks.New(p, cfg)

	if err := api.Register(); err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--cluster")
	}

	if ng.Name != "" && nameArg != "" {
		return cmdutils.ErrNameFlagAndArg(ng.Name, nameArg)
	}

	if nameArg != "" {
		ng.Name = nameArg
	}

	if ng.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	if updateAuthConfigMap {
		// remove node group from config map
		if err := authconfigmap.RemoveNodeGroup(clientSet, ng); err != nil {
			logger.Warning(err.Error())
		}
	}

	if deleteNodeGroupDrain {
		if err := drain.NodeGroup(clientSet, ng, ctl.Provider.WaitTimeout(), false); err != nil {
			return err
		}
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

	return nil
}
