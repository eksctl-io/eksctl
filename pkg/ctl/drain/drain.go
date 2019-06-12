package drain

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `drain` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("drain", "Drain resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, drainNodeGroupCmd)

	return verbCmd
}
