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

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a token for debugging purpose",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tok := viper.GetString("token")
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

		fmt.Printf("%+v\n", id)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringP("token", "t", "", "Verify this token")
	viper.BindPFlag("token", verifyCmd.Flags().Lookup("token"))
}
