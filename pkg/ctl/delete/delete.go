package delete

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `delete` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("delete", "Delete resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, deleteClusterCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, deleteNodeGroupCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, deleteIAMIdentityMappingCmd)

	return verbCmd
}
