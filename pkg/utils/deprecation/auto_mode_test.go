package deprecation

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func TestShouldShowAutoModeDeprecationWarning_NilConfig(t *testing.T) {
	cfg := &v1alpha5.ClusterConfig{
		AutoModeConfig: nil,
	}

	result := shouldShowAutoModeDeprecationWarning(cfg)
	if !result {
		t.Error("Expected warning to be shown when AutoModeConfig is nil")
	}
}

func TestShouldShowAutoModeDeprecationWarning_NilEnabled(t *testing.T) {
	cfg := &v1alpha5.ClusterConfig{
		AutoModeConfig: &v1alpha5.AutoModeConfig{
			Enabled: nil,
		},
	}

	result := shouldShowAutoModeDeprecationWarning(cfg)
	if !result {
		t.Error("Expected warning to be shown when AutoModeConfig.Enabled is nil")
	}
}

func TestShouldShowAutoModeDeprecationWarning_ExplicitlyFalse(t *testing.T) {
	enabled := false
	cfg := &v1alpha5.ClusterConfig{
		AutoModeConfig: &v1alpha5.AutoModeConfig{
			Enabled: &enabled,
		},
	}

	result := shouldShowAutoModeDeprecationWarning(cfg)
	if result {
		t.Error("Expected no warning when AutoModeConfig.Enabled is explicitly false")
	}
}

func TestShouldShowAutoModeDeprecationWarning_ExplicitlyTrue(t *testing.T) {
	enabled := true
	cfg := &v1alpha5.ClusterConfig{
		AutoModeConfig: &v1alpha5.AutoModeConfig{
			Enabled: &enabled,
		},
	}

	result := shouldShowAutoModeDeprecationWarning(cfg)
	if result {
		t.Error("Expected no warning when AutoModeConfig.Enabled is explicitly true")
	}
}

func TestDefaultAutoModeDeprecationConfig(t *testing.T) {
	config := DefaultAutoModeDeprecationConfig()

	if config.Message == "" {
		t.Error("Expected default message to be non-empty")
	}

	if config.DocumentationURL == "" {
		t.Error("Expected default documentation URL to be non-empty")
	}

	if config.LogLevel != "warning" {
		t.Errorf("Expected default log level to be 'warning', got '%s'", config.LogLevel)
	}

	expectedMessage := "Auto Mode will be enabled by default in an upcoming release of eksctl. " +
		"This means managed node groups and managed networking add-ons will no longer be created by default. " +
		"To maintain current behavior, explicitly set 'autoModeConfig.enabled: false' in your cluster configuration."

	if config.Message != expectedMessage {
		t.Errorf("Expected default message to match specification")
	}
}

func TestCheckAutoModeDeprecation_WithNilConfig(t *testing.T) {
	cfg := &v1alpha5.ClusterConfig{
		AutoModeConfig: nil,
	}

	// This test verifies the function doesn't panic and executes the warning logic
	// We can't easily test the actual logging without mocking the logger
	CheckAutoModeDeprecation(cfg)
}

func TestCheckAutoModeDeprecation_WithExplicitlySetConfig(t *testing.T) {
	enabled := false
	cfg := &v1alpha5.ClusterConfig{
		AutoModeConfig: &v1alpha5.AutoModeConfig{
			Enabled: &enabled,
		},
	}

	// This test verifies the function doesn't show warning when explicitly set
	// Since we can't easily mock the logger, we just ensure it doesn't panic
	CheckAutoModeDeprecation(cfg)
}

func TestCheckAutoModeDeprecationWithConfig_CustomConfig(t *testing.T) {
	cfg := &v1alpha5.ClusterConfig{
		AutoModeConfig: nil,
	}

	customConfig := &AutoModeDeprecationConfig{
		Message:          "Custom test message",
		DocumentationURL: "https://example.com",
		LogLevel:         "info",
	}

	// This test verifies the function accepts custom configuration
	CheckAutoModeDeprecationWithConfig(cfg, customConfig)
}
