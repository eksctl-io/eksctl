package set

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command creates the `set` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("set", "Set values", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, setLabelsCmd)

	return verbCmd
}
