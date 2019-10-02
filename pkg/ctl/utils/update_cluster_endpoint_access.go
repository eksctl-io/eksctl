package utils

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var (
	private bool
	public  bool
)

func updateClusterEndpointsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-cluster-endpoints", "Update Kubernetes API endpoint access configuration", "")

	cmd.SetRunFunc(func() error {
		return doUpdateClusterEndpoints(cmd, private, public)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("Update private/public Kubernetes API endpoint access configuration",
		func(fs *pflag.FlagSet) {
			fs.BoolVar(&private, "private-access", false, "access for private (VPC) clients")
			fs.BoolVar(&public, "public-access", false, "access for public clients")
		})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func accessFlagsSet(cmd *cmdutils.Cmd) (privateSet, publicSet bool) {
	cmd.FlagSetGroup.InFlagSet("Update private/public Kubernetes API endpoint access configuration",
		func(fs *pflag.FlagSet) {
			fs.VisitAll(func(f *pflag.Flag) {
				switch f.Name {
				case "private-access":
					privateSet = f.Changed
				case "public-access":
					publicSet = f.Changed
				default:
					// do nothing
				}
			})
		})
	return
}

func doUpdateClusterEndpoints(cmd *cmdutils.Cmd, newPrivate bool, newPublic bool) error {
	privateSet, publicSet := accessFlagsSet(cmd)
	if err := cmdutils.NewUtilsEnableEndpointAccessLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	curPrivate, curPublic, err := ctl.GetCurrentClusterConfigForEndpoints(cfg)
	if err != nil {
		return err
	}

	logger.Info("current Kubernetes API endpoint access: privateAccess=%v, publicAccess=%v",
		curPrivate, curPublic)

	if !privateSet {
		newPrivate = curPrivate
	}
	if !publicSet {
		newPublic = curPublic
	}

	needsUpdate := false
	if newPrivate != curPrivate || newPublic != curPublic {
		needsUpdate = true
	}

	if needsUpdate {
		if newPrivate == false && newPublic == false {
			logger.Critical(api.NoAccessMsg(
				api.ClusterEndpoints{PublicAccess: &newPublic, PrivateAccess: &newPrivate},
			))
			os.Exit(2)
		} else if newPrivate == true && newPublic == false {
			logger.Warning(api.PrivateOnlyAwsChangesNeededMsg())
		}

		cfg.VPC.ClusterEndpoints.PrivateAccess = &newPrivate
		cfg.VPC.ClusterEndpoints.PublicAccess = &newPublic

		describeAccessToUpdate :=
			fmt.Sprintf("privateAccess=%v, publicAccess=%v", newPrivate, newPublic)

		cmdutils.LogIntendedAction(
			cmd.Plan, "update Kubernetes API Endpoint Access for cluster %q in %q to: %s",
			meta.Name, meta.Region, describeAccessToUpdate)
		if !cmd.Plan {
			if err := ctl.UpdateClusterConfigForEndpoints(cfg); err != nil {
				return err
			}
			cmdutils.LogCompletedAction(
				false,
				"the Kubernetes API Endpoint Access for cluster %q in %q has been upddated to: "+
					"privateAccess=%v, publicAccess=%v",
				meta.Name, meta.Region, newPrivate, newPublic)
		}
	} else {
		logger.Success("Kubernetes API Endpoint Access for cluster %q in %q is already up to date",
			meta.Name, meta.Region)
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && needsUpdate)

	return nil
}

func validateClusterEndpointConfig(private, public bool) error {
	cfg := api.NewClusterConfig()
	cfg.VPC.ClusterEndpoints = &api.ClusterEndpoints{PrivateAccess: &private, PublicAccess: &public}
	err := cfg.ValidateClusterEndpointConfig()
	if err != nil {
		// utils can change public access to false since cluster creation has already completed
		if err.Error() == api.PrivateOnlyUseUtilsMsg() {
			return nil
		}
		return err
	}
	return nil
}
