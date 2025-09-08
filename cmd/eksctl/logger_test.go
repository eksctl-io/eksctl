package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/fatih/color"
)

func TestInitLoggerDisablesColors(t *testing.T) {
	tests := []struct {
		name        string
		colorValue  string
		noColorEnv  string
		expectColor bool
	}{
		{
			name:        "color disabled with false",
			colorValue:  "false",
			expectColor: false,
		},
		{
			name:        "color disabled with NO_COLOR env",
			colorValue:  "true",
			noColorEnv:  "1",
			expectColor: false,
		},
		{
			name:        "color enabled with true",
			colorValue:  "true",
			expectColor: true,
		},
		{
			name:        "color enabled with fabulous",
			colorValue:  "fabulous",
			expectColor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset color state
			color.NoColor = false

			// Set NO_COLOR environment variable if specified
			if tt.noColorEnv != "" {
				os.Setenv("NO_COLOR", tt.noColorEnv)
				defer os.Unsetenv("NO_COLOR")
			}

			// Initialize logger
			logBuffer := new(bytes.Buffer)
			initLogger(3, tt.colorValue, logBuffer, false)

			// Check if colors are disabled as expected
			if tt.expectColor && color.NoColor {
				t.Errorf("Expected colors to be enabled, but color.NoColor is true")
			}
			if !tt.expectColor && !color.NoColor {
				t.Errorf("Expected colors to be disabled, but color.NoColor is false")
			}
		})
	}
}
