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

// Parse parses taints but does not validate the fields for correctness
func Parse(taints map[string]string) []corev1.Taint {
	var parsedTaints []corev1.Taint
	for k, v := range taints {
		parsedTaints = append(parsedTaints, parseTaint(k, v))
	}
	return parsedTaints
}

// parseTaint parses a taint from valueEffect, whose form must be either
// '<value>:<effect>' or ':<effect>'.
func parseTaint(key, valueEffect string) corev1.Taint {
	var (
		value  string
		effect corev1.TaintEffect
	)
	parts := strings.Split(valueEffect, ":")
	switch len(parts) {
	case 1:
		effect = corev1.TaintEffect(parts[0])
	case 2:
		value, effect = parts[0], corev1.TaintEffect(parts[1])
	}

	return corev1.Taint{
		Key:    key,
		Value:  value,
		Effect: effect,
	}
}

func validateTaintEffect(effect corev1.TaintEffect) error {
	switch effect {
	case corev1.TaintEffectNoSchedule, corev1.TaintEffectPreferNoSchedule, corev1.TaintEffectNoExecute:
		return nil
	default:
		return fmt.Errorf(
			"invalid taint effect: %v, unsupported taint effect. Valid taint effects are: %s, %s and %s",
			effect,
			corev1.TaintEffectNoSchedule,
			corev1.TaintEffectNoExecute,
			corev1.TaintEffectPreferNoSchedule,
		)
	}
}
