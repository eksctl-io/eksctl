package utils

import (
	"fmt"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// FormatTaints constructs a taint string of the form
// key1=value1:NoEffect,key2=value2:NoSchedule
func FormatTaints(taints []api.NodeGroupTaint) string {
	var params []string
	for _, t := range taints {
		params = append(params, fmt.Sprintf("%s=%s:%s", t.Key, t.Value, t.Effect))
	}
	return strings.Join(params, ",")
}
