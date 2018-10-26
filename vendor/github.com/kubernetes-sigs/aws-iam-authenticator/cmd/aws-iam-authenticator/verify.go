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
	"encoding/json"
	"fmt"
	"os"

	"github.com/kubernetes-sigs/aws-iam-authenticator/pkg/token"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a token for debugging purpose",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tok := viper.GetString("token")
		output := viper.GetString("output")
		clusterID := viper.GetString("clusterID")

		if tok == "" {
			fmt.Fprintf(os.Stderr, "error: token not specified\n")
			cmd.Usage()
			os.Exit(1)
		}

		if clusterID == "" {
			fmt.Fprintf(os.Stderr, "error: cluster ID not specified\n")
			cmd.Usage()
			os.Exit(1)
		}

		id, err := token.NewVerifier(clusterID).Verify(tok)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not verify token: %v\n", err)
			os.Exit(1)
		}

		if output == "json" {
			value, err := json.MarshalIndent(id, "", "    ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not unmarshal token: %v\n", err)
			}
			fmt.Printf("%s\n", value)
		} else {
			fmt.Printf("%+v\n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringP("token", "t", "", "Token to verify")
	verifyCmd.Flags().StringP("output", "o", "", "Output format. Only `json` is supported currently.")
	viper.BindPFlag("token", verifyCmd.Flags().Lookup("token"))
	viper.BindPFlag("output", verifyCmd.Flags().Lookup("output"))
}
