package unset

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command creates the `unset` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("unset", "Unset values", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, unsetLabelsCmd)

	return verbCmd
}
