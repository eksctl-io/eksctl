package powershell_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/nodebootstrap/powershell"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestPowerShellUtils(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("PowerShell", func() {

	type powerShellEntry struct {
		params         []powershell.KeyValue
		expectedOutput string
	}

	DescribeTable("FormatStringVariables", func(e powerShellEntry) {
		variables := powershell.FormatStringVariables(e.params)
		Expect(variables).To(Equal(e.expectedOutput))
	},
		Entry("standard params", powerShellEntry{
			params: []powershell.KeyValue{
				{
					Key:   "EKSClusterName",
					Value: "cluster",
				},
				{
					Key:   "APIServerEndpoint",
					Value: "https://test.com",
				},
				{
					Key:   "Var",
					Value: "Value",
				},
				{
					Key:   "Key1",
					Value: `"Value1"`,
				},
			},
			expectedOutput: `[string]$EKSClusterName = 'cluster'
[string]$APIServerEndpoint = 'https://test.com'
[string]$Var = 'Value'
[string]$Key1 = '"Value1"'`,
		}),

		Entry("empty params", powerShellEntry{
			expectedOutput: "",
		}),

		Entry("single param", powerShellEntry{
			params: []powershell.KeyValue{
				{
					Key:   "EKSClusterName",
					Value: "cluster",
				},
			},
			expectedOutput: `[string]$EKSClusterName = 'cluster'`,
		}),
	)

	DescribeTable("FormatHashTable", func(e powerShellEntry, variableName string) {
		hashtable := powershell.FormatHashTable(e.params, variableName)
		Expect(hashtable).To(Equal(e.expectedOutput))
	},

		Entry("standard params", powerShellEntry{
			params: []powershell.KeyValue{
				{
					Key:   "node-labels",
					Value: "",
				},
				{
					Key:   "key1",
					Value: "val1",
				},
				{
					Key:   "key2",
					Value: "val2",
				},
			},
			expectedOutput: `$KubeletExtraArgs = @{ 'node-labels' = ''; 'key1' = 'val1'; 'key2' = 'val2'}`,
		},
			"KubeletExtraArgs",
		),

		Entry(
			"empty params",
			powerShellEntry{
				expectedOutput: `$Map = @{}`,
			},
			"Map",
		),
	)

	DescribeTable("FormatParams", func(e powerShellEntry) {
		formatted := powershell.FormatParams(e.params)
		Expect(formatted).To(Equal(e.expectedOutput))
	},
		Entry("standard params", powerShellEntry{
			params: []powershell.KeyValue{
				{
					Key:   "EKSClusterName",
					Value: "cluster",
				},
				{
					Key:   "APIServerEndpoint",
					Value: "https://test.com",
				},
				{
					Key:   "Var",
					Value: "Value",
				},
			},
			expectedOutput: `-EKSClusterName "cluster" -APIServerEndpoint "https://test.com" -Var "Value"`,
		}),

		Entry("empty params", powerShellEntry{
			expectedOutput: "",
		}),
	)

	DescribeTable("ToCLIArgs", func(e powerShellEntry) {
		cliArgs := powershell.ToCLIArgs(e.params)
		Expect(cliArgs).To(Equal(e.expectedOutput))
	},
		Entry("standard params", powerShellEntry{
			params: []powershell.KeyValue{
				{
					Key:   "node-labels",
					Value: "",
				},
				{
					Key:   "key1",
					Value: "val1",
				},
				{
					Key:   "key2",
					Value: "val2",
				},
			},
			expectedOutput: "--node-labels= --key1=val1 --key2=val2",
		}),

		Entry("empty params", powerShellEntry{
			expectedOutput: "",
		}),
	)
})
