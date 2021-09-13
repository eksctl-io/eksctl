package deregister

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command creates the `set` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("deregister", "Deregister a non-EKS cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, deregisterClusterCmd)

	return verbCmd
}
