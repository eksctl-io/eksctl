package identityproviders

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

type waitForUpdateTask struct {
	ipm *IdentityProviderManager
	update eks.Update
	timeout time.Duration
}

func (ipm *IdentityProviderManager) waitForUpdate(
	update eks.Update, timeout time.Duration,
) error {
	clusterName := ipm.metadata.Name
	newRequest := func() *request.Request {
		input := &eks.DescribeUpdateInput{
			Name:     aws.String(ipm.metadata.Name),
			UpdateId: update.Id,
		}
		req, _ := ipm.eksAPI.DescribeUpdateRequest(input)
		return req
	}

	acceptors := waiters.MakeAcceptors(
		"Update.Status",
		eks.UpdateStatusSuccessful,
		[]string{
			eks.UpdateStatusCancelled,
			eks.UpdateStatusFailed,
		},
	)

	msg := fmt.Sprintf(
		"waiting for update %q in cluster %q to succeed",
		*update.Type,
		clusterName,
	)

	if err := waiters.Wait(clusterName, msg, acceptors, newRequest, timeout, nil); err != nil {
		return err
	}
	return nil
}
