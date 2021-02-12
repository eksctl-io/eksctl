package get

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

type getCmdParams struct {
	chunkSize int
	output    printers.Type
}

// Command will create the `get` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("get", "Get resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getClusterCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getNodeGroupCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getIdentityProvider)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getIAMServiceAccountCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getIAMIdentityMappingCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getLabelsCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getFargateProfile)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, getAddonCmd)

	return verbCmd
}
