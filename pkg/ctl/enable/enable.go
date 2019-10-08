package enable

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `enable` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("enable", "Enable features in a cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, enableProfileCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, enableRepo)

	return verbCmd
}
