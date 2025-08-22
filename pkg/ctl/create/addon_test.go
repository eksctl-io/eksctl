package create

import (
	"testing"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/stretchr/testify/assert"
)

func TestParseNamespaceConfig(t *testing.T) {
	tests := []struct {
		name           string
		namespaceConfig string
		expectedError  string
		expectedNamespace string
	}{
		{
			name:            "valid namespace config",
			namespaceConfig: "namespace=custom-namespace",
			expectedNamespace: "custom-namespace",
		},
		{
			name:            "valid namespace config with spaces",
			namespaceConfig: " namespace = my-namespace ",
			expectedNamespace: "my-namespace",
		},
		{
			name:            "invalid format - missing equals",
			namespaceConfig: "namespace",
			expectedError:   "expected format 'namespace=<namespace-name>', got \"namespace\"",
		},
		{
			name:            "invalid format - wrong key",
			namespaceConfig: "name=test",
			expectedError:   "unsupported key \"name\", only 'namespace' is supported",
		},
		{
			name:            "empty namespace value",
			namespaceConfig: "namespace=",
			expectedError:   "namespace value cannot be empty",
		},
		{
			name:            "empty namespace value with spaces",
			namespaceConfig: "namespace=   ",
			expectedError:   "namespace value cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addon := &api.Addon{}
			err := parseNamespaceConfig(tt.namespaceConfig, addon)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, addon.NamespaceConfig)
				assert.Equal(t, tt.expectedNamespace, addon.NamespaceConfig.Namespace)
			}
		})
	}
}

func TestParseNamespaceConfig_PreservesExistingConfig(t *testing.T) {
	addon := &api.Addon{
		NamespaceConfig: &api.AddonNamespaceConfig{
			Namespace: "existing-namespace",
		},
	}

	err := parseNamespaceConfig("namespace=new-namespace", addon)
	assert.NoError(t, err)
	assert.Equal(t, "new-namespace", addon.NamespaceConfig.Namespace)
}