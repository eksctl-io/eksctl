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

func updateEndpointAccessCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-cluster-api-access", "Update cluster API endpoint configuration", "")

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	var private, public bool
	cmd.FlagSetGroup.InFlagSet("Allow/Disallow Cluster API access flags",
		func(fs *pflag.FlagSet) {
			fs.BoolVar(&private, "private-access", false, "access for private (VPC) clients")
			fs.BoolVar(&public, "public-access", true, "access for public clients")
		})
	cmd.SetRunFuncWithNameArg(func() error {
		return doConfigureEndpointAccess(cmd, private, public)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doConfigureEndpointAccess(cmd *cmdutils.Cmd, newPrivate bool, newPublic bool) error {
	if err := cmdutils.NewUtilsEnableEndpointAccessLoader(cmd).Load(); err != nil {
		return err
	}

	if cmd.ClusterConfig.HasClusterEndpointAccess() {
		if err := validateEndpointConfig(newPrivate, newPublic); err != nil {
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
	if err!= nil {
		return err
	}
	logger.Info("Current endpoint access: private access:%v, public access: %v",
		*curPrivate, *curPublic)

	if newPrivate == false && newPublic == false {
		logger.Critical(api.NoAccessMsg(
			&api.ClusterEndpoints{PublicAccess: &newPublic, PrivateAccess: &newPrivate}))
		os.Exit(2)
	} else if newPrivate == true && newPublic == false {
		logger.Warning(api.PrivateOnlyAwsChangesNeededMsg())
	}

	describeAccessToUpdate := []string{"no updates to make"}

	needsUpdate := false
	if newPrivate != *curPrivate || newPublic != *curPublic {
		needsUpdate = true
	}

	if needsUpdate {
		cfg.VPC.ClusterEndpoints.PrivateAccess = &newPrivate
		describeAccessToUpdate[0] = fmt.Sprintf("private access:%v", newPrivate)

		cfg.VPC.ClusterEndpoints.PublicAccess = &newPublic
		msg := fmt.Sprintf("public access:%v", newPublic)
		describeAccessToUpdate = append(describeAccessToUpdate, msg)

		describeAccessUpdate := strings.Join(describeAccessToUpdate, ", ")

		cmdutils.LogIntendedAction(
			cmd.Plan, "update Cluster API Endpoint Access for cluster %q in %q to: (%s)",
			meta.Name, meta.Region, describeAccessUpdate)
		if !cmd.Plan {
			if err := ctl.UpdateClusterConfigForEndpoints(cfg); err != nil {
				return err
			}
			cmdutils.LogCompletedAction(
				false,
				"Cluster API Endpoint Access for cluster %q in %q has been upddated to: {%s}",
				meta.Name, meta.Region, describeAccessUpdate)
		}
	} else {
		logger.Success("Cluster API Endpoint Access for cluster %q in %q is already up to date", meta.Name, meta.Region)
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && needsUpdate)

	return nil
}

func validateEndpointConfig(private, public bool) error {
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
