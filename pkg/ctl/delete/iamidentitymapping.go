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
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	var (
		role string
		all  bool
	)

	cp.Command = &cobra.Command{
		Use:   "iamidentitymapping",
		Short: "Delete a IAM identity mapping",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doDeleteIAMIdentityMapping(cp, role, all); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&role, "role", "", "ARN of the IAM role to delete")
		fs.BoolVar(&all, "all", false, "Delete all matching mappings instead of just one")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, false)

	group.AddTo(cp.Command)
	return cp.Command
}

func doDeleteIAMIdentityMapping(cp *cmdutils.CommonParams, role string, all bool) error {
	cfg := cp.ClusterConfig

	ctl := eks.New(cp.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if role == "" {
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

	if err := acm.RemoveRole(role, all); err != nil {
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
	filtered := roles.Get(role)
	if len(filtered) > 0 {
		logger.Warning("there are %d mappings left with same role %q (use --all to delete them at once)", len(filtered), role)
	}
	return nil
}
