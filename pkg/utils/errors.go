package utils

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
)

func HasAwsErrorCode(err error, awsErrCode string) bool {
	for err != nil {
		awserr, ok := err.(awserr.Error)

		if !ok {
			cause := errors.Cause(err)
			if cause == err {
				break
			}
			err = cause
		} else {
			return awserr.Code() == awsErrCode
		}
	}

	return false
}
