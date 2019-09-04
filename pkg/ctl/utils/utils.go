package utils

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `utils` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("utils", "Various utils", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, waitNodesCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, writeKubeconfigCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, describeStacksCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateClusterStackCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateKubeProxyCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateAWSNodeCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateCoreDNSCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, enableLoggingCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, associateIAMOIDCProviderCmd)

	return verbCmd
}
