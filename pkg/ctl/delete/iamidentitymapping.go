package delete

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func deleteIAMIdentityMappingCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	var roleFlag string
	var all bool
	cmd := &cobra.Command{
		Use:   "iamidentitymapping",
		Short: "Delete a IAM identity mapping",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doDeleteIAMIdentityMapping(p, cfg, roleFlag, all); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&roleFlag, "role", "", "ARN of the IAM role to delete")
		fs.BoolVar(&all, "all", false, "Delete all matching mappings instead of just one")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, p)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)

	return cmd
}

func doDeleteIAMIdentityMapping(p *api.ProviderConfig, cfg *api.ClusterConfig, roleFlag string, all bool) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if roleFlag == "" {
		return cmdutils.ErrMustBeSet("--role")
	}
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--cluster")
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}
	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}

	if err := acm.RemoveRole(roleFlag, all); err != nil {
		return err
	}
	if err := acm.Save(); err != nil {
		return err
	}

	// Check whether we have more roles that match
	roles, err := acm.Roles()
	if err != nil {
		return err
	}
	filtered := roles.Get(roleFlag)
	if len(filtered) > 0 {
		logger.Warning("there are %d mappings left with same role %q (use --all to delete them at once)", len(filtered), roleFlag)
	}
	return nil
}
