package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Output the version of eksctl",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))
			return nil
		},
	}
}
