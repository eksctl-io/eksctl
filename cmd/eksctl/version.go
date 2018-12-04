package main

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/version"

	"github.com/kubicorn/kubicorn/pkg/logger"
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
