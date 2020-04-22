package completion

import (
	"github.com/kris-nova/logger"
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

If you are stuck on Bash 3 (macOS) use

source /dev/stdin <<<"$(eksctl completion bash)"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenBashCompletion(cmd.OutOrStdout())
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenZshCompletion(cmd.OutOrStdout())
		},
	}

	var fishCompletionCmd = &cobra.Command{
		Use:   "fish",
		Short: "Generates fish completion scripts",
		Long: `To configure your fish shell, run:

mkdir -p ~/.config/fish/completions
eksctl completion fish > ~/.config/fish/completions/eksctl.fish

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		},
	}

	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates shell completion scripts for bash, zsh or fish",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(bashCompletionCmd)
	cmd.AddCommand(zshCompletionCmd)
	cmd.AddCommand(fishCompletionCmd)

	return cmd
}
