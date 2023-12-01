package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// STS provides an interface to the AWS STS service.
type STS interface {
	// GetCallerIdentity returns details about the IAM user or role whose credentials are used to call
	// the operation. No permissions are required to perform this operation. If an
	// administrator attaches a policy to your identity that explicitly denies access
	// to the sts:GetCallerIdentity action, you can still perform this operation.
	// Permissions are not required because the same information is returned when
	// access is denied. To view an example response, see I Am Not Authorized to
	// Perform: iam:DeleteVirtualMFADevice (https://docs.aws.amazon.com/IAM/latest/UserGuide/troubleshoot_general.html#troubleshoot_general_access-denied-delete-mfa)
	// in the IAM User Guide.
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}
