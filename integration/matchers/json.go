package matchers

import (
	"github.com/onsi/gomega/types"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"encoding/json"
	"fmt"
)

// BeNodeGroupsWithNamesWhich helps match JSON-formatted nodegroups by
// accepting matchers on an array of nodegroup names.
func BeNodeGroupsWithNamesWhich(matchers ...types.GomegaMatcher) types.GomegaMatcher {
	return &jsonNodeGroupMatcher{
		matchers: matchers,
	}
}

type jsonNodeGroupMatcher struct {
	matchers              []types.GomegaMatcher
	failureMessage        string
	negatedFailureMessage string
}

func (matcher *jsonNodeGroupMatcher) Match(actual interface{}) (success bool, err error) {
	rawJSON, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("BeNodeGroupsWithNamesWhich matcher expects a string")
	}
	ngSummaries := []manager.NodeGroupSummary{}
	if err := json.Unmarshal([]byte(rawJSON), &ngSummaries); err != nil {
		return false, fmt.Errorf("BeNodeGroupsWithNamesWhich matcher expects a NodeGroupSummary JSON array")
	}
	ngNames := extractNames(ngSummaries)
	for _, m := range matcher.matchers {
		if ok, err := m.Match(ngNames); !ok {
			matcher.failureMessage = m.FailureMessage(ngNames)
			matcher.negatedFailureMessage = m.NegatedFailureMessage(ngNames)
			return false, err
		}
	}
	return true, nil
}

func extractNames(ngSummaries []manager.NodeGroupSummary) []string {
	ngNames := make([]string, len(ngSummaries))
	for i, ngSummary := range ngSummaries {
		ngNames[i] = ngSummary.Name
	}
	return ngNames
}

func (matcher *jsonNodeGroupMatcher) FailureMessage(unusedActual interface{}) string {
	return matcher.failureMessage
}

func (matcher *jsonNodeGroupMatcher) NegatedFailureMessage(unusedActual interface{}) string {
	return matcher.negatedFailureMessage
}
