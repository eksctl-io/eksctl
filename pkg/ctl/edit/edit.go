package edit

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `edit` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("edit", "edit resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, editClusterCmd)

	return verbCmd
}
