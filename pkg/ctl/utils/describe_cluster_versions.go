package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func describeClusterVersionsCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription("describe-cluster-versions", "Describe Supported Kubernetes Versions", "")

	var clusterVersions []string
	var defaultOnly, include bool
	var clusterTypes, status string

	cmd.FlagSetGroup.InFlagSet("Versions", func(fs *pflag.FlagSet) {
		fs.StringVar(&clusterTypes, "cluster-types", "", "(optional) Specify the cluster type(s) to filter results. Valid options: eks, eks-local-outposts. Example: --cluster-types \"eks\"")
		fs.StringSliceVar(&clusterVersions, "cluster-versions", []string{}, "(optional) Filter results by specific Kubernetes versions. Accepts multiple comma-separated values. Example: --cluster-versions \"1.31,1.30\"")
		fs.BoolVar(&include, "include-all", false, "(optional) When set, includes unsupported versions in the results. Default: false")
		fs.StringVar(&status, "status", "", "(optional) Filter results by support plan status. Valid options: EXTENDED_SUPPORT, STANDARD_SUPPORT")
		fs.BoolVar(&defaultOnly, "default-only", false, "(optional) When set, returns only the default version for the specified cluster type(s). Default: false")
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return describeClusterVersions(cmd, clusterTypes, clusterVersions, include, defaultOnly, status)
	}

}

func describeClusterVersions(cmd *cmdutils.Cmd, clusterTypes string, clusterVersions []string, include, defaultOnly bool, status string) error {
	clusterProvider, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	ctx := context.TODO()

	versions, err := clusterProvider.AWSProvider.EKS().DescribeClusterVersions(
		ctx,
		&eks.DescribeClusterVersionsInput{
			ClusterType:     &clusterTypes,
			ClusterVersions: clusterVersions,
			IncludeAll:      &include,
			DefaultOnly:     &defaultOnly,
			Status:          types.ClusterVersionStatus(status),
		},
	)

	if err != nil {
		return err
	}

	jsonPrinter := printers.NewJSONPrinter()
	printerErr := jsonPrinter.PrintObj(struct {
		ClusterVersions []types.ClusterVersionInformation `json:"clusterVersions"`
	}{
		ClusterVersions: versions.ClusterVersions,
	}, os.Stdout)
	if printerErr != nil {
		fmt.Printf("Error printing JSON: %v\n", printerErr)
	}

	return nil
}
