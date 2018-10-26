/*
Copyright 2017 by the contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Pre-generate certificate, private key, and kubeconfig files for the server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := getConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get config: %v\n", err)
			os.Exit(1)
		}

		localCfg := cfg
		localCfg.GenerateKubeconfigPath = "aws-iam-authenticator.kubeconfig"
		localCfg.StateDir = "./"

		err = localCfg.GenerateFiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not initialize: %v\n", err)
			os.Exit(1)
		}

		logrus.Infof("copy %s to %s on kubernetes master node(s)", localCfg.CertPath(), cfg.CertPath())
		logrus.Infof("copy %s to %s on kubernetes master node(s)", localCfg.KeyPath(), cfg.KeyPath())
		logrus.Infof("copy %s to %s on kubernetes master node(s)", localCfg.GenerateKubeconfigPath, cfg.GenerateKubeconfigPath)
		logrus.Infof("configure your apiserver with `--authentication-token-webhook-config-file=%s` to enable authentication with aws-iam-authenticator", cfg.GenerateKubeconfigPath)
	},
}

func init() {
	initCmd.Flags().String(
		"hostname",
		"localhost",
		"Hostname that should be used for writing the self-signed certificates")
	viper.BindPFlag("server.hostname", initCmd.Flags().Lookup("hostname"))

	initCmd.Flags().String(
		"address",
		"127.0.0.1",
		"IP Address to bind the server to listen to. (should be a 127.0.0.1 or 0.0.0.0)")
	viper.BindPFlag("server.address", initCmd.Flags().Lookup("address"))

	rootCmd.AddCommand(initCmd)
}
