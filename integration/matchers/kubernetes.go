package matchers

import (
	"fmt"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

//BeNotFoundError succeeds if actual is a non-nil error
//which represents a missing kubernetes resource
func BeNotFoundError() types.GomegaMatcher {
	return &notFoundMatcher{}
}

type notFoundMatcher struct {
}

func (matcher *notFoundMatcher) Match(actual interface{}) (success bool, err error) {
	isErr, err := gomega.HaveOccurred().Match(actual)
	if !isErr || err != nil {
		return isErr, err
	}
	return apierrors.IsNotFound(actual.(error)), nil
}

func (matcher *notFoundMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected a NotFound API error to have occurred.  Got:\n%s", format.Object(actual, 1))
}

func (matcher *notFoundMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Unexpected NotFound API error:\n%s\n%s\n%s", format.Object(actual, 1), format.IndentString(actual.(error).Error(), 1), "occurred")
}
