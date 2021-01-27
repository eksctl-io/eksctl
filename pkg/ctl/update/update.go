package update

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `create` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("update", "Update resource(s)", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateClusterCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateAddonCmd)
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, updateIAMServiceAccountCmd)

	return verbCmd
}
