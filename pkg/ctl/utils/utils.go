package utils

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var errUnsupportedLocalCluster = errors.New("this operation is not supported on local clusters")

// Command will create the `utils` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("utils", "Various utils", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, writeKubeconfigCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, describeStacksCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateKubeProxyCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateAWSNodeCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateCoreDNSCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateLegacySubnetSettings)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, enableLoggingCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, associateIAMOIDCProviderCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, installWindowsVPCController)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateClusterEndpointsCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, publicAccessCIDRsCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, enableSecretsEncryptionCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, schemaCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, nodeGroupHealthCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, describeAddonVersionsCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, describeAddonConfigurationCmd)

	return verbCmd
}
