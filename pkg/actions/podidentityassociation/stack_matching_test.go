package podidentityassociation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIAMResourcesStack(t *testing.T) {
	tests := []struct {
		name          string
		stackNames    []string
		identifier    Identifier
		expectedStack string
		expectedFound bool
		description   string
	}{
		{
			name: "exact match found",
			stackNames: []string{
				"eksctl-cluster-podidentityrole-service-service-account",
			},
			identifier: Identifier{
				Namespace:          "service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "eksctl-cluster-podidentityrole-service-service-account",
			expectedFound: true,
			description:   "Should find exact match",
		},
		{
			name: "no substring false positive",
			stackNames: []string{
				"eksctl-cluster-podidentityrole-other-service-service-account",
			},
			identifier: Identifier{
				Namespace:          "service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "",
			expectedFound: false,
			description:   "Should not match 'service-service-account' as substring of 'other-service-service-account'",
		},
		{
			name: "multiple stacks with correct match",
			stackNames: []string{
				"eksctl-cluster-podidentityrole-other-service-service-account",
				"eksctl-cluster-podidentityrole-service-service-account",
				"eksctl-cluster-podidentityrole-another-service-different-account",
			},
			identifier: Identifier{
				Namespace:          "service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "eksctl-cluster-podidentityrole-service-service-account",
			expectedFound: true,
			description:   "Should find correct match among multiple stacks",
		},
		{
			name: "no match found",
			stackNames: []string{
				"eksctl-cluster-podidentityrole-other-service-service-account",
				"eksctl-cluster-podidentityrole-another-service-different-account",
			},
			identifier: Identifier{
				Namespace:          "service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "",
			expectedFound: false,
			description:   "Should not find any match",
		},
		{
			name: "customer reported case",
			stackNames: []string{
				"eksctl-cluster-podidentityrole-service-service-account",
				"eksctl-cluster-podidentityrole-other-service-service-account",
			},
			identifier: Identifier{
				Namespace:          "service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "eksctl-cluster-podidentityrole-service-service-account",
			expectedFound: true,
			description:   "Customer reported case: should match 'service' not 'other-service'",
		},
		{
			name: "customer reported case reverse",
			stackNames: []string{
				"eksctl-cluster-podidentityrole-service-service-account",
				"eksctl-cluster-podidentityrole-other-service-service-account",
			},
			identifier: Identifier{
				Namespace:          "other-service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "eksctl-cluster-podidentityrole-other-service-service-account",
			expectedFound: true,
			description:   "Customer reported case reverse: should match 'other-service' not 'service'",
		},
		{
			name: "IRSAv1 pattern match",
			stackNames: []string{
				"eksctl-cluster-addon-iamserviceaccount-service-service-account",
			},
			identifier: Identifier{
				Namespace:          "service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "eksctl-cluster-addon-iamserviceaccount-service-service-account",
			expectedFound: true,
			description:   "Should match IRSAv1 pattern",
		},
		{
			name: "IRSAv1 no substring false positive",
			stackNames: []string{
				"eksctl-cluster-addon-iamserviceaccount-other-service-service-account",
			},
			identifier: Identifier{
				Namespace:          "service",
				ServiceAccountName: "service-account",
			},
			expectedStack: "",
			expectedFound: false,
			description:   "Should not match IRSAv1 'service-service-account' as substring of 'other-service-service-account'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualStack, actualFound := getIAMResourcesStack(tt.stackNames, tt.identifier)
			assert.Equal(t, tt.expectedFound, actualFound, tt.description)
			assert.Equal(t, tt.expectedStack, actualStack, tt.description)
		})
	}
}
