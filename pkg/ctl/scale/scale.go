package scale

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `scale` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("scale", "Scale resources(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, scaleNodeGroupCmd)

	return verbCmd
}
