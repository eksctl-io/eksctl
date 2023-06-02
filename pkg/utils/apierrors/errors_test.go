package apierrors_test

import (
	"fmt"
	"testing"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/apierrors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtilsAPIErrors(t *testing.T) {
	testutils.RegisterAndRun(t)
}

type retriableErrorsEntry struct {
	err             error
	shouldBeRetried bool
}

var _ = Describe("APIErrors", func() {
	DescribeTable("IsRetriableError", func(e retriableErrorsEntry) {
		Expect(apierrors.IsRetriableError(e.err)).To(Equal(e.shouldBeRetried))
	},
		Entry("Non API Error", retriableErrorsEntry{
			err:             fmt.Errorf("Non API Error"),
			shouldBeRetried: true,
		}),
		Entry("ServerException", retriableErrorsEntry{
			err:             &ekstypes.ServerException{},
			shouldBeRetried: true,
		}),
		Entry("ServiceUnavailableException", retriableErrorsEntry{
			err:             &ekstypes.ServiceUnavailableException{},
			shouldBeRetried: true,
		}),
		Entry("InvalidRequestException", retriableErrorsEntry{
			err:             &ekstypes.InvalidRequestException{},
			shouldBeRetried: true,
		}),
		Entry("BadRequestException", retriableErrorsEntry{
			err:             &ekstypes.BadRequestException{},
			shouldBeRetried: false,
		}),
		Entry("NotFoundException", retriableErrorsEntry{
			err:             &ekstypes.NotFoundException{},
			shouldBeRetried: false,
		}),
		Entry("ResourceNotFoundException", retriableErrorsEntry{
			err:             &ekstypes.ResourceNotFoundException{},
			shouldBeRetried: false,
		}),
		Entry("AccessDeniedException", retriableErrorsEntry{
			err:             &ekstypes.AccessDeniedException{},
			shouldBeRetried: false,
		}),
	)
})
