package utils

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLoggingFlags(t *testing.T) {
	var (
		errBothEmpty        = "at least one flag has to be provided"
		errAllBoth          = "cannot use `all` for both"
		errOverlappingTypes = "log types cannot be part of --enable-types and --disable-types"
		errUnknownType      = "unknown log type"
	)

	flagTests := []struct {
		toEnable   []string
		toDisable  []string
		errPattern string
	}{
		{
			toEnable:   []string{"all"},
			toDisable:  []string{"audit", "scheduler", "api"},
			errPattern: "",
		},
		{
			toEnable:   []string{"audit", "scheduler", "api"},
			toDisable:  []string{"controllerManager", "authenticator"},
			errPattern: "",
		},
		{
			toEnable:   []string{"all"},
			toDisable:  []string{"controllerManager", "authenticator"},
			errPattern: "",
		},
		{
			toEnable:   []string{"all"},
			toDisable:  []string{"all"},
			errPattern: errAllBoth,
		},
		{
			toEnable:  []string{"all", "api"},
			toDisable: []string{"all"},
			// TODO improve error reporting for {"all", "api", ...}
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{""},
			toDisable:  []string{"all"},
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{"all", "invalid"},
			errPattern: errUnknownType,
		},
		{
			toDisable:  []string{"all", "api", "invalid", "scheduler"},
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{""},
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{""},
			toDisable:  []string{""},
			errPattern: errUnknownType,
		},
		{
			errPattern: errBothEmpty,
		},
		{
			toEnable:   []string{"api", "audit"},
			toDisable:  []string{"audit", "authenticator", "scheduler"},
			errPattern: errOverlappingTypes,
		},
		{
			toEnable:   []string{"audit", "authenticator", "scheduler"},
			toDisable:  []string{"api", "scheduler", "controllerManager"},
			errPattern: errOverlappingTypes,
		},
	}

	for i, tt := range flagTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := validateLoggingFlags(tt.toEnable, tt.toDisable)
			if err != nil {
				if tt.errPattern == "" {
					t.Errorf("unexpected error: %v", err)
				} else if !strings.Contains(err.Error(), tt.errPattern) {
					t.Errorf("expected error %q to match %q", err, tt.errPattern)
				}
			} else if tt.errPattern != "" {
				t.Errorf("expected error %q; got nil", tt.errPattern)
			}
		})
	}

}

func TestParseCIDRs(t *testing.T) {
	cidrTests := []struct {
		input       string
		expected    []string
		shouldError bool
	}{
		{
			input:    "192.168.0.1/32,192.168.0.2/32,192.168.0.3/32",
			expected: []string{"192.168.0.1/32", "192.168.0.2/32", "192.168.0.3/32"},
		},
		{
			input:    "192.168.0.1/32",
			expected: []string{"192.168.0.1/32"},
		},
		{
			input:       "fail",
			shouldError: true,
		},
		{
			input:       "256.0.0.0/32",
			shouldError: true,
		},
		{
			input:       "192.168.0.1/32,192.168.0.2/33",
			shouldError: true,
		},
		{
			input:    "192.168.0.0/24,192.168.32.0/19,192.168.64.0/19",
			expected: []string{"192.168.0.0/24", "192.168.32.0/19", "192.168.64.0/19"},
		},
		{
			input:       "192.168.0.1/32, 192.168.0.2/32, 192.168.0.3/32",
			shouldError: true, // space after comma
		},
	}

	for i, tt := range cidrTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cidrs, err := parseCIDRs(tt.input)
			if err != nil {
				if tt.shouldError {
					return
				}
				t.Errorf("expected %v; got %v", tt.expected, err)
				return
			}
			assert.Equal(t, cidrs, tt.expected)

		})
	}
}
