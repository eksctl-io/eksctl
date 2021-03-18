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
