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
	"errors"
	"fmt"
	"os"

	"github.com/kubernetes-sigs/aws-iam-authenticator/pkg/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "aws-iam-authenticator",
	Short: "A tool to authenticate to Kubernetes using AWS IAM credentials",
}

func main() {
	Execute()
}

// Execute the CLI entrypoint
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Load configuration from `filename`")

	rootCmd.PersistentFlags().StringP("log-format", "l", "text", "Specify log format to use when logging to stderr [text or json]")

	rootCmd.PersistentFlags().StringP(
		"cluster-id",
		"i",
		"",
		"Specify the cluster `ID`, a unique-per-cluster identifier for your aws-iam-authenticator installation.",
	)
	viper.BindPFlag("clusterID", rootCmd.PersistentFlags().Lookup("cluster-id"))
	viper.BindEnv("clusterID", "KUBERNETES_AWS_AUTHENTICATOR_CLUSTER_ID")
}

func initConfig() {
	logrus.SetFormatter(getLogFormatter())
	if cfgFile == "" {
		return
	}
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Can't read configuration file %q: %v\n", cfgFile, err)
		os.Exit(1)
	}
}

func getConfig() (config.Config, error) {
	config := config.Config{
		ClusterID:                         viper.GetString("clusterID"),
		ServerEC2DescribeInstancesRoleARN: viper.GetString("server.ec2DescribeInstancesRoleARN"),
		HostPort:               viper.GetInt("server.port"),
		Hostname:               viper.GetString("server.hostname"),
		GenerateKubeconfigPath: viper.GetString("server.generateKubeconfig"),
		KubeconfigPregenerated: viper.GetBool("server.kubeconfigPregenerated"),
		StateDir:               viper.GetString("server.stateDir"),
		Address:                viper.GetString("server.address"),
	}
	if err := viper.UnmarshalKey("server.mapRoles", &config.RoleMappings); err != nil {
		return config, fmt.Errorf("invalid server role mappings: %v", err)
	}
	if err := viper.UnmarshalKey("server.mapUsers", &config.UserMappings); err != nil {
		logrus.WithError(err).Fatal("invalid server user mappings")
	}
	if err := viper.UnmarshalKey("server.mapAccounts", &config.AutoMappedAWSAccounts); err != nil {
		logrus.WithError(err).Fatal("invalid server account mappings")
	}

	if config.ClusterID == "" {
		return config, errors.New("cluster ID cannot be empty")
	}

	return config, nil
}

func getLogFormatter() logrus.Formatter {
	format, _ := rootCmd.PersistentFlags().GetString("log-format")

	if format == "json" {
		return &logrus.JSONFormatter{}
	} else if format != "text" {
		logrus.Warnf("Unknown log format specified (%s), will use default text formatter instead.", format)
	}

	return &logrus.TextFormatter{FullTimestamp: true}
}
