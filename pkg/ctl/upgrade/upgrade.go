package upgrade

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `upgrade` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("upgrade", "Upgrade resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, upgradeNodeGroupCmd)

	return verbCmd
}
