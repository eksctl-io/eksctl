package apierrors

import (
	"errors"

	"github.com/aws/smithy-go"
)

func IsRetriableError(err error) bool {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return true
	}
	switch {
	case isBadRequestErrorCode(apiErr),
		isNotFoundErrorCode(apiErr),
		isAccessDeniedErrorCode(apiErr):
		return false
	default:
		return true
	}
}

func IsServerError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && isServerErrorCode(apiErr)
}

func isServerErrorCode(apiErr smithy.APIError) bool {
	return apiErr.Error() == "ServerException"
}

func IsServiceUnavailableError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && isServiceUnavailableErrorCode(apiErr)
}

func isServiceUnavailableErrorCode(apiErr smithy.APIError) bool {
	return apiErr.Error() == "ServiceUnavailableException"
}

func IsInvalidRequestError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && isInvalidRequestErrorCode(apiErr)
}

func isInvalidRequestErrorCode(apiErr smithy.APIError) bool {
	return apiErr.Error() == "InvalidRequestException"
}

func IsBadRequestError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && isBadRequestErrorCode(apiErr)
}

func isBadRequestErrorCode(apiErr smithy.APIError) bool {
	return apiErr.ErrorCode() == "BadRequestException"
}

func IsNotFoundError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && isNotFoundErrorCode(apiErr)
}

func isNotFoundErrorCode(apiErr smithy.APIError) bool {
	return apiErr.ErrorCode() == "NotFoundException" || apiErr.ErrorCode() == "ResourceNotFoundException"
}

func IsAccessDeniedError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && isAccessDeniedErrorCode(apiErr)
}

func isAccessDeniedErrorCode(apiErr smithy.APIError) bool {
	return apiErr.ErrorCode() == "AccessDenied" || apiErr.ErrorCode() == "AccessDeniedException"
}
