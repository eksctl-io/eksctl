package main

import (
	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

//go:generate go run version_generate.go

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Output the version of eksctl",
		Run: func(_ *cobra.Command, _ []string) {
			versionInfo := map[string]string{
				"builtAt":   builtAt,
				"gitCommit": gitCommit,
			}

			if gitTag != "" {
				versionInfo["gitTag"] = gitTag
			}

			logger.Info("versionInfo = %#v", versionInfo)
		},
	}
}
