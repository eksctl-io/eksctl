package waiter

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

// WaitForNodegroupUpdate waits for an update to finish. Once it's done it returns the final status.
// If the status was not Successful it will also return an error.
func WaitForNodegroupUpdate(ctx context.Context, updateID string, api awsapi.EKS, timeout time.Duration, nextDelay NextDelay) (string, error) {
	var lastStatus string
	waiter := &Waiter{
		NextDelay: nextDelay,
		Operation: func() (bool, error) {
			var (
				err     error
				success bool
			)
			lastStatus, success, err = describeUpdateStatus(ctx, updateID, api)
			return success, err
		},
	}

	if err := waiter.WaitWithTimeout(timeout); err != nil {
		return lastStatus, err
	}
	return lastStatus, nil
}

func describeUpdateStatus(ctx context.Context, updateID string, api awsapi.EKS) (string, bool, error) {
	logger.Info("waiting for nodegroup update %s to finish", updateID)
	update, err := api.DescribeUpdate(ctx, &eks.DescribeUpdateInput{
		UpdateId: aws.String(updateID),
	})
	if err != nil {
		return "", false, fmt.Errorf("failed to describe update for update id %s: %w", updateID, err)
	}
	if update.Update == nil {
		return "", false, fmt.Errorf("update field of output is empty")
	}
	switch update.Update.Status {
	case ekstypes.UpdateStatusInProgress:
		return string(ekstypes.UpdateStatusInProgress), false, nil
	case ekstypes.UpdateStatusCancelled, ekstypes.UpdateStatusFailed:
		logger.Info("attempting to display any errors that could have been received...")
		if len(update.Update.Errors) == 0 {
			logger.Info("no errors received from update")
		} else {
			l := len(update.Update.Errors)
			for i, err := range update.Update.Errors {
				if err.ErrorMessage != nil {
					logger.Info("received following error message(s) (%d/%d): %s", i, l, err.ErrorMessage)
				}
			}
		}
		return string(update.Update.Status), false, fmt.Errorf("update failed or was cancelled")
	case ekstypes.UpdateStatusSuccessful:
		return string(ekstypes.UpdateStatusSuccessful), true, nil
	}
	return string(update.Update.Status), false, nil
}
