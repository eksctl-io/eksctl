package completion

import (
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/spf13/cobra"
)

// Command will create the `completion` commands
func Command(rootCmd *cobra.Command) *cobra.Command {
	var bashCompletionCmd = &cobra.Command{
		Use:   "bash",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(eksctl completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(eksctl completion bash)
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
	var zshCompletionCmd = &cobra.Command{
		Use:   "zsh",
		Short: "Generates zsh completion scripts",
		Long: `To configure your zsh shell, run:

mkdir -p ~/.zsh/completion/
eksctl completion zsh > ~/.zsh/completion/_eksctl

and put the following in ~/.zshrc:

fpath=($fpath ~/.zsh/completion)

`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenZshCompletion(os.Stdout)
		},
	}

	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates shell completion scripts",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(bashCompletionCmd)
	cmd.AddCommand(zshCompletionCmd)

	return cmd
}
