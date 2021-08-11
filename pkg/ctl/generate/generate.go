package generate

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command creates `generate` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("generate",
		"DEPRECATED: https://github.com/weaveworks/eksctl/issues/2963\nGenerate gitops manifests",
		"",
	)
	verbCmd.Hidden = true
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, generateProfile)
	return verbCmd
}
