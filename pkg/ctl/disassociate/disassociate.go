package disassociate

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `disassociate` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("disassociate", "Disassociate resources from a cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, disassociateIdentityProvider)

	return verbCmd
}
