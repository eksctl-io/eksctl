package podidentityassociation

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// ValidatePodIdentityAssociation validates a pod identity association configuration.
func ValidatePodIdentityAssociation(pia *api.PodIdentityAssociation) error {
	if pia.Namespace == "" {
		return errors.New("namespace cannot be empty")
	}

	if pia.ServiceAccountName == "" {
		return errors.New("serviceAccountName cannot be empty")
	}

	// If targetRoleARN is specified, validate it
	if pia.TargetRoleARN != "" {
		if err := validateRoleARN(pia.TargetRoleARN); err != nil {
			return fmt.Errorf("invalid targetRoleARN: %w", err)
		}

		// If roleARN is empty but targetRoleARN is specified, we need to create a source role
		if pia.RoleARN == "" && pia.RoleName == "" {
			// This is fine, we'll create a source role
		} else if pia.RoleARN != "" {
			// Validate that the source role is in the same account as the cluster
			sourceAccountID, err := getAccountIDFromARN(pia.RoleARN)
			if err != nil {
				return fmt.Errorf("invalid roleARN: %w", err)
			}

			targetAccountID, err := getAccountIDFromARN(pia.TargetRoleARN)
			if err != nil {
				return fmt.Errorf("invalid targetRoleARN: %w", err)
			}

			if sourceAccountID == targetAccountID {
				return errors.New("targetRoleARN must be in a different account than roleARN for cross-account access")
			}
		}
	}

	return nil
}

// validateRoleARN validates that the provided string is a valid IAM role ARN.
func validateRoleARN(arn string) error {
	if !strings.HasPrefix(arn, "arn:aws:iam::") {
		return errors.New("ARN must start with 'arn:aws:iam::'")
	}

	parts := strings.Split(arn, ":")
	if len(parts) != 6 {
		return errors.New("ARN must have 6 parts separated by colons")
	}

	if parts[5] == "" || !strings.HasPrefix(parts[5], "role/") {
		return errors.New("ARN must end with 'role/ROLE_NAME'")
	}

	accountID := parts[4]
	if len(accountID) != 12 || !isNumeric(accountID) {
		return errors.New("account ID in ARN must be a 12-digit number")
	}

	return nil
}

// getAccountIDFromARN extracts the account ID from an IAM role ARN.
func getAccountIDFromARN(arn string) (string, error) {
	parts := strings.Split(arn, ":")
	if len(parts) != 6 {
		return "", errors.New("ARN must have 6 parts separated by colons")
	}
	return parts[4], nil
}

// isNumeric checks if a string contains only digits.
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
