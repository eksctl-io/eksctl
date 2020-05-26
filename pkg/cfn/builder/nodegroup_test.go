package builder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func TestGenerateNodeName(t *testing.T) {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = "foo"

	nodeGroup := api.NewNodeGroup()
	nodeGroup.Name = "bar"

	nodeNameTests := []struct {
		prefix      string
		name        string
		expected    string
		description string
	}{
		{
			prefix:      "",
			name:        "",
			expected:    "foo-bar-Node",
			description: "Default naming convention",
		},
		{
			prefix:      "hello",
			name:        "",
			expected:    "hello-foo-bar-Node",
			description: "Each node has a specific prefix",
		},
		{
			prefix:      "",
			name:        "i-am-the-master",
			expected:    "i-am-the-master",
			description: "Each node has a full override of the name",
		},
		{
			prefix:      "prefix",
			name:        "specific-name",
			expected:    "prefix-specific-name",
			description: "Each node has a prefix with a specific name",
		},
	}

	for i, tt := range nodeNameTests {
		t.Run(fmt.Sprintf("%d: %s", i, tt.description), func(t *testing.T) {
			nodeGroup.InstancePrefix = tt.prefix
			nodeGroup.InstanceName = tt.name

			n := NewNodeGroupResourceSet(nil, clusterConfig, "cluster", nodeGroup, false)
			nodeName := n.generateNodeName()

			assert.Equal(t, tt.expected, nodeName)
		})
	}
}
