package experimental

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `experimental` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("experimental", "Various commands that are currently experimental",
		"WARNING: In any release these commands may get removed, renamed or their behaviour may change significantly")

	// cmdutils.AddResourceCmd(flagGrouping, verbCmd, fn)

	return verbCmd
}
