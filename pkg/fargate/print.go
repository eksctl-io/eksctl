package fargate

import (
	"io"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/printers"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	kindFargateProfiles = "fargateprofiles"
)

// PrintProfiles formats the provided profiles in the provided printer type
// ("table", "json", "yaml") and prints them to the provided writer.
func PrintProfiles(profiles []*api.FargateProfile, writer io.Writer, printerType printers.Type) error {
	printer, err := printers.NewPrinter(printerType)
	if err != nil {
		return err
	}
	switch printerType {
	case printers.TableType:
		addFargateProfileColumns(printer.(*printers.TablePrinter))
		return printer.PrintObjWithKind(kindFargateProfiles, toTable(profiles), writer)
	default:
		return printer.PrintObjWithKind(kindFargateProfiles, profiles, writer)
	}
}

type row struct {
	Name                string
	PodExecutionRoleARN string
	Subnets             []string
	Selector            api.FargateProfileSelector
	Tags                map[string]string
}

func toTable(profiles []*api.FargateProfile) []*row {
	table := []*row{}
	for _, profile := range profiles {
		for _, selector := range profile.Selectors {
			table = append(table, &row{
				Name:                profile.Name,
				PodExecutionRoleARN: profile.PodExecutionRoleARN,
				Subnets:             profile.Subnets,
				Selector:            selector,
				Tags:                profile.Tags,
			})
		}
	}
	return table
}

func addFargateProfileColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(r *row) string {
		return r.Name
	})
	printer.AddColumn("SELECTOR_NAMESPACE", func(r *row) string {
		return r.Selector.Namespace
	})
	printer.AddColumn("SELECTOR_LABELS", func(r *row) string {
		return labels.FormatLabels(r.Selector.Labels)
	})
	printer.AddColumn("POD_EXECUTION_ROLE_ARN", func(r *row) string {
		return r.PodExecutionRoleARN
	})
	printer.AddColumn("SUBNETS", func(r *row) string {
		if len(r.Subnets) == 0 {
			return "<none>"
		}
		return strings.Join(r.Subnets, ",")
	})
	printer.AddColumn("TAGS", func(r *row) string {
		return labels.FormatLabels(r.Tags)
	})
}
