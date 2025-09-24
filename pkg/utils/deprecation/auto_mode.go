package deprecation

import (
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// AutoModeDeprecationConfig holds deprecation warning settings
type AutoModeDeprecationConfig struct {
	Message          string
	DocumentationURL string
	LogLevel         string
}

// DefaultAutoModeDeprecationConfig returns default deprecation warning config
func DefaultAutoModeDeprecationConfig() *AutoModeDeprecationConfig {
	return &AutoModeDeprecationConfig{
		Message: "Auto Mode will be enabled by default in an upcoming release of eksctl. " +
			"This means managed node groups and managed networking add-ons will no longer be created by default. " +
			"To maintain current behavior, explicitly set 'autoModeConfig.enabled: false' in your cluster configuration.",
		DocumentationURL: "https://eksctl.io/usage/auto-mode/",
		LogLevel:         "warning",
	}
}

// CheckAutoModeDeprecation shows deprecation warning if needed
func CheckAutoModeDeprecation(cfg *v1alpha5.ClusterConfig) {
	if shouldShowAutoModeDeprecationWarning(cfg) {
		logAutoModeDeprecationWarning(DefaultAutoModeDeprecationConfig())
	}
}

// CheckAutoModeDeprecationWithConfig shows deprecation warning with custom config
func CheckAutoModeDeprecationWithConfig(cfg *v1alpha5.ClusterConfig, deprecationConfig *AutoModeDeprecationConfig) {
	if shouldShowAutoModeDeprecationWarning(cfg) {
		logAutoModeDeprecationWarning(deprecationConfig)
	}
}

// shouldShowAutoModeDeprecationWarning returns true if warning should be shown
func shouldShowAutoModeDeprecationWarning(cfg *v1alpha5.ClusterConfig) bool {
	// Show warning only when AutoMode is nil or Enabled is nil (not explicitly set)
	return cfg.AutoModeConfig == nil || cfg.AutoModeConfig.Enabled == nil
}

// logAutoModeDeprecationWarning logs the deprecation warning
func logAutoModeDeprecationWarning(config *AutoModeDeprecationConfig) {
	fullMessage := config.Message
	if config.DocumentationURL != "" {
		fullMessage += " Learn more: " + config.DocumentationURL
	}

	// Log at configured level
	switch config.LogLevel {
	case "info":
		logger.Info(fullMessage)
	case "warning":
		logger.Warning(fullMessage)
	default:
		logger.Warning(fullMessage)
	}
}
