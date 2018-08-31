package harness

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func operatorToSelectionOperator(op metav1.LabelSelectorOperator) selection.Operator {
	switch op {
	case metav1.LabelSelectorOpIn:
		return selection.In
	case metav1.LabelSelectorOpNotIn:
		return selection.NotIn
	case metav1.LabelSelectorOpExists:
		return selection.Exists
	case metav1.LabelSelectorOpDoesNotExist:
		// XXX: the selection package doesn't define the DoesNotExist operator.
		fallthrough
	default:
		panic(fmt.Sprintf("selector: unexpected operator '%s'", op))
	}
}

func selectorToString(selector *metav1.LabelSelector) (string, error) {
	set := labels.Set(selector.MatchLabels)
	s := set.AsSelector()
	for _, req := range selector.MatchExpressions {
		newReq, err := labels.NewRequirement(req.Key, operatorToSelectionOperator(req.Operator), req.Values)
		if err != nil {
			return "", err
		}
		s.Add(*newReq)
	}
	return s.String(), nil
}
