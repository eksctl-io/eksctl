package builder_test

import (
	"fmt"

	"github.com/onsi/gomega/types"

	. "github.com/weaveworks/eksctl/pkg/cfn/builder"
)

type ResourceSetRenderingMatcher struct {
	templateBody *[]byte
	err          error
}

func RenderWithoutErrors(templateBody *[]byte) types.GomegaMatcher {
	return &ResourceSetRenderingMatcher{
		templateBody: templateBody,
	}
}

func (m *ResourceSetRenderingMatcher) Match(actualResourceSet interface{}) (bool, error) {
	if actualResourceSet == nil {
		return false, fmt.Errorf("resourceset is nil")
	}

	if _, ok := actualResourceSet.(ResourceSet); !ok {
		return false, fmt.Errorf("not a resourceset")
	}

	m.err = actualResourceSet.(ResourceSet).AddAllResources()
	if m.err != nil {
		return false, nil
	}

	*m.templateBody, m.err = actualResourceSet.(ResourceSet).RenderJSON()
	if m.err != nil {
		return false, nil
	}
	return true, nil
}

func (m *ResourceSetRenderingMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected to add all resouerces and render JSON without errors, got: %s", m.err.Error())
}

func (m *ResourceSetRenderingMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected to NOT load template from JSON without errors")
}
