package main

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/version"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Output the version of eksctl",
		Run: func(_ *cobra.Command, _ []string) {

			logger.Info("%#v", version.Get())
		},
	}
}
