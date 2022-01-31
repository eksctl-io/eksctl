package register

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command creates the `set` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("register", "Register a non-EKS cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, registerClusterCmd)

	return verbCmd
}
