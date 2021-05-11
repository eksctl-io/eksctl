package taints

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

// Validate validates taints
func Validate(t corev1.Taint) error {
	if t.Key == "" {
		return errors.New("taint key must be non-empty")
	}

	if errs := validation.IsQualifiedName(t.Key); len(errs) > 0 {
		return errors.Errorf("invalid taint key: %v, %s", t.Key, strings.Join(errs, "; "))
	}

	if t.Value != "" {
		if errs := validation.IsValidLabelValue(t.Value); len(errs) > 0 {
			return errors.Errorf("invalid taint value: %v, %s", t.Value, strings.Join(errs, "; "))
		}
	}

	return validateTaintEffect(t.Effect)
}

func validateTaintEffect(effect corev1.TaintEffect) error {
	switch effect {
	case corev1.TaintEffectNoSchedule, corev1.TaintEffectPreferNoSchedule, corev1.TaintEffectNoExecute:
		return nil
	default:
		return fmt.Errorf("invalid taint effect: %v, unsupported taint effect", effect)
	}
}
