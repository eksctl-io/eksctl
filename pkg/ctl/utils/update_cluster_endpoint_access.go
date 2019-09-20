package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateClusterEndpointsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-cluster-endpoints", "Update Kubernetes API endpoint access configuration", "")

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	var private, public bool
	cmd.FlagSetGroup.InFlagSet("Update private/public Kubernetes API endpoint access configuration",
		func(fs *pflag.FlagSet) {
			fs.BoolVar(&private, "private-access", false, "access for private (VPC) clients")
			fs.BoolVar(&public, "public-access", true, "access for public clients")
		})
	cmd.SetRunFunc(func() error {
		return doUpdateClusterEndpoints(cmd, private, public)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doUpdateClusterEndpoints(cmd *cmdutils.Cmd, newPrivate bool, newPublic bool) error {
	if err := cmdutils.NewUtilsEnableEndpointAccessLoader(cmd).Load(); err != nil {
		return err
	}

	if cmd.ClusterConfig.HasClusterEndpointAccess() {
		if err := validateClusterEndpointConfig(newPrivate, newPublic); err != nil {
			return err
		}
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

	if newPrivate == false && newPublic == false {
		logger.Critical(api.NoAccessMsg(
			&api.ClusterEndpoints{PublicAccess: &newPublic, PrivateAccess: &newPrivate},
		))
		os.Exit(2)
	} else if newPrivate == true && newPublic == false {
		logger.Warning(api.PrivateOnlyAwsChangesNeededMsg())
	}

	needsUpdate := false
	if newPrivate != curPrivate || newPublic != curPublic {
		needsUpdate = true
	}

	if needsUpdate {
		describeAccessToUpdate := []string{}

		cfg.VPC.ClusterEndpoints.PrivateAccess = &newPrivate
		if newPrivate != curPrivate {
			describeAccessToUpdate =
				append(describeAccessToUpdate, fmt.Sprintf("privateAccess=%v", newPrivate))
		}

		cfg.VPC.ClusterEndpoints.PublicAccess = &newPublic
		if newPublic != curPublic {
			describeAccessToUpdate =
				append(describeAccessToUpdate, fmt.Sprintf("publicAccess=%v", newPublic))
		}
		describeAccessUpdate := strings.Join(describeAccessToUpdate, ", ")

		cmdutils.LogIntendedAction(
			cmd.Plan, "update Kubernetes API Endpoint Access for cluster %q in %q to: %s",
			meta.Name, meta.Region, describeAccessUpdate)
		if !cmd.Plan {
			if err := ctl.UpdateClusterConfigForEndpoints(cfg); err != nil {
				return err
			}
			cmdutils.LogCompletedAction(
				false,
				"the Kubernetes API Endpoint Access for cluster %q in %q has been upddated to: privateAccess=%s, publicAccess=%s",
				meta.Name, meta.Region, newPrivate, newPublic)
		}
	} else {
		logger.Success("Kubernetes API Endpoint Access for cluster %q in %q is already up to date", meta.Name, meta.Region)
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
