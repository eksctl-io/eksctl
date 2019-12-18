package utils

import (
	"encoding/csv"
	"net"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func setPublicAccessCIDRsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("set-public-access-cidrs", "Update public access CIDRs", "CIDR blocks that EKS uses to create a security group on the public endpoint")

	cmd.CobraCommand.Args = cobra.ExactArgs(1)
	cmd.CobraCommand.RunE = func(c *cobra.Command, args []string) error {
		cidrs, err := parseCIDRs(args[0])
		if err != nil {
			return err
		}
		validCIDRs, err := validateCIDRs(cidrs)
		if err != nil {
			return err
		}
		cmd.ClusterConfig.VPC.PublicAccessCIDRs = validCIDRs
		return doUpdatePublicAccessCIDRs(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func parseCIDRs(arg string) ([]string, error) {
	reader := strings.NewReader(arg)
	csvReader := csv.NewReader(reader)
	return csvReader.Read()
}

func validateCIDRs(cidrs []string) ([]string, error) {
	var validCIDRs []string
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		validCIDRs = append(validCIDRs, ipNet.String())
	}
	return validCIDRs, nil
}

func doUpdatePublicAccessCIDRs(cmd *cmdutils.Cmd) error {
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

	clusterVPCConfig, err := ctl.GetCurrentClusterVPCConfig(cfg)
	if err != nil {
		return err
	}

	logger.Info("current public access CIDRs: %v", clusterVPCConfig.PublicAccessCIDRs)

	if cidrsEqual(clusterVPCConfig.PublicAccessCIDRs, cfg.VPC.PublicAccessCIDRs) {
		logger.Success("Public Endpoint Restrictions for cluster %q in %q is already up to date",
			meta.Name, meta.Region)
		return nil
	}

	cmdutils.LogIntendedAction(
		cmd.Plan, "update Public Endpoint Restrictions for cluster %q in %q to: %v",
		meta.Name, meta.Region, cfg.VPC.PublicAccessCIDRs)

	if !cmd.Plan {
		if err := ctl.UpdatePublicAccessCIDRs(cfg); err != nil {
			return errors.Wrap(err, "error updating CIDRs for public access")
		}
		cmdutils.LogCompletedAction(
			false,
			"Public Endpoint Restrictions for cluster %q in %q have been updated to: %v",
			meta.Name, meta.Region, cfg.VPC.PublicAccessCIDRs)
	}
	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}

func cidrsEqual(currentValues, newValues []string) bool {
	return sets.NewString(currentValues...).Equal(sets.NewString(newValues...))
}
