package powershell

import (
	"fmt"
	"strings"
)

// A KeyValue holds a key-value pair.
type KeyValue struct {
	Key   string
	Value string
}

// FormatStringVariables formats params as PowerShell string variables.
// The caller is responsible for ensuring the keys are valid PowerShell variable names
// and the values can be surrounded by single-quoted strings.
func FormatStringVariables(params []KeyValue) string {
	variables := make([]string, 0, len(params))
	for _, param := range params {
		variables = append(variables, fmt.Sprintf("[string]$%s = '%s'", param.Key, param.Value))
	}
	return strings.Join(variables, "\n")
}

// JoinVariables joins the specified variables.
func JoinVariables(variables ...string) string {
	return strings.Join(variables, "\n")
}

// FormatHashTable formats table as a PowerShell hashtable.
// The caller is responsible for ensuring variableName is a valid PowerShell variable name
// and the key-value pairs can be surrounded by single-quoted strings.
func FormatHashTable(table []KeyValue, variableName string) string {
	values := make([]string, 0, len(table))
	for _, kv := range table {
		values = append(values, fmt.Sprintf(" '%s' = '%s'", kv.Key, kv.Value))
	}
	return fmt.Sprintf("$%s = @{%s}", variableName, strings.Join(values, ";"))
}

// FormatParams formats params into `-key "value"`, ignoring keys with empty values.
func FormatParams(params []KeyValue) string {
	var args []string
	for _, param := range params {
		if param.Value != "" {
			args = append(args, fmt.Sprintf("-%s %q", param.Key, param.Value))
		}
	}
	return strings.Join(args, " ")
}

// ToCLIArgs formats params as CLI arguments.
func ToCLIArgs(params []KeyValue) string {
	var args []string
	for _, param := range params {
		args = append(args, fmt.Sprintf("--%s=%s", param.Key, param.Value))
	}
	return strings.Join(args, " ")
}
