package identityproviders

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

func (m *Manager) waitForUpdate(
	update eks.Update, timeout time.Duration,
) error {
	clusterName := m.metadata.Name
	newRequest := func() *request.Request {
		input := &eks.DescribeUpdateInput{
			Name:     aws.String(clusterName),
			UpdateId: update.Id,
		}
		req, _ := m.eksAPI.DescribeUpdateRequest(input)
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

	return waiters.Wait(clusterName, msg, acceptors, newRequest, timeout, nil)
}
