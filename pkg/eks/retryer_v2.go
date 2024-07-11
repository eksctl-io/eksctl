package eks

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/smithy-go"
)

const (
	maxRetries = 13
)

// RetryerV2 implements aws.Retryer
type RetryerV2 struct {
	aws.Retryer
}

// NewRetryerV2 returns a new *RetryerV2
func NewRetryerV2() *RetryerV2 {
	standard := retry.AddWithMaxAttempts(retry.NewStandard(func(o *retry.StandardOptions) {
		o.MaxAttempts = maxRetries
	}), maxRetries)

	return &RetryerV2{
		Retryer: standard,
	}
}

// IsErrorRetryable implements aws.Retryer
func (r *RetryerV2) IsErrorRetryable(err error) bool {
	if !r.Retryer.IsErrorRetryable(err) {
		return false
	}

	var oe *smithy.OperationError
	if !errors.As(err, &oe) {
		return true
	}
	return oe.Err != nil && isErrorRetryable(oe.Err)
}

func isErrorRetryable(err error) bool {
	if aerr, ok := err.(awserr.RequestFailure); ok && aerr.Code() == "EC2MetadataError" {
		switch aerr.StatusCode() {
		case http.StatusForbidden, http.StatusNotFound, http.StatusMethodNotAllowed:
			return false
		}
	}
	return true
}
