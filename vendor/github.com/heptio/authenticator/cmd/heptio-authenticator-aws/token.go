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

	"github.com/heptio/authenticator/pkg/token"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Authenticate using AWS IAM and get token for Kubernetes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		roleARN := viper.GetString("role")
		clusterID := viper.GetString("clusterID")

		if clusterID == "" {
			fmt.Fprintf(os.Stderr, "Error: cluster ID not specified\n")
			cmd.Usage()
			os.Exit(1)
		}

		var tok string
		var err error
		gen, err := token.NewGenerator()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get token: %v\n", err)
			os.Exit(1)
		}
		if roleARN != "" {
			// if a role was provided, assume that role for the token
			tok, err = gen.GetWithRole(clusterID, roleARN)
		} else {
			// otherwise sign the token with immediately available credentials
			tok, err = gen.Get(clusterID)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get token: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(gen.FormatJSON(tok))
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().StringP("role", "r", "", "Assume an IAM Role ARN before signing this token")
	viper.BindPFlag("role", tokenCmd.Flags().Lookup("role"))
	viper.BindEnv("role", "DEFAULT_ROLE")
}
