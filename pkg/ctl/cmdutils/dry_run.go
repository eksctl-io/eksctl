package cmdutils

import (
	"io"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/printers"
)

// PrintDryRunConfig prints ClusterConfig for dry-run
func PrintDryRunConfig(clusterConfig *v1alpha5.ClusterConfig, writer io.Writer) error {
	yamlPrinter := printers.NewYAMLPrinter()
	return yamlPrinter.PrintObj(clusterConfig, writer)
}

// PrintNodeGroupDryRunConfig prints the dry-run config for nodegroups, omitting any cluster-wide defaults
func PrintNodeGroupDryRunConfig(clusterConfig *v1alpha5.ClusterConfig, writer io.Writer) error {
	output := &v1alpha5.ClusterConfig{
		TypeMeta:          clusterConfig.TypeMeta,
		Metadata:          clusterConfig.Metadata,
		NodeGroups:        clusterConfig.NodeGroups,
		ManagedNodeGroups: clusterConfig.ManagedNodeGroups,
	}
	return PrintDryRunConfig(output, writer)
}
