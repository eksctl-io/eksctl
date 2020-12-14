package associate

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `associate` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("associate", "Associate resources with a cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, associateIdentityProvider)

	return verbCmd
}
