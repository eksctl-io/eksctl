package install

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `install` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("install", "Install components in a cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, installFluxCmd)

	return verbCmd
}
