package generate

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command creates `generate` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("generate", "Generate gitops manifests", "")
	cmdutils.AddResourceCmd(flagGrouping, verbCmd, generateProfile)
	return verbCmd
}
