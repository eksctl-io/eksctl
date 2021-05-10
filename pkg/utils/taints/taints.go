package taints

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

// Parse parses taint values
func Parse(taints map[string]string) ([]corev1.Taint, error) {
	var parsedTaints []corev1.Taint
	for k, v := range taints {
		taint, err := parseTaint(k, v)
		if err != nil {
			return nil, err
		}
		parsedTaints = append(parsedTaints, taint)
	}
	return parsedTaints, nil
}

// parseTaint parses a taint from valueEffect, whose form must be either
// '<value>:<effect>' or ':<effect>'.
func parseTaint(key, valueEffect string) (corev1.Taint, error) {
	if errs := validation.IsQualifiedName(key); len(errs) > 0 {
		return corev1.Taint{}, fmt.Errorf("invalid taint key: %v, %s", key, strings.Join(errs, "; "))
	}

	parts := strings.Split(valueEffect, ":")
	var (
		value  string
		effect corev1.TaintEffect
	)
	switch len(parts) {
	case 1:
		effect = corev1.TaintEffect(parts[0])
	case 2:
		value, effect = parts[0], corev1.TaintEffect(parts[1])
		if errs := validation.IsValidLabelValue(value); len(errs) > 0 {
			return corev1.Taint{}, fmt.Errorf("invalid taint value: %v, %s", value, strings.Join(errs, "; "))
		}
	}

	if err := validateTaintEffect(effect); err != nil {
		return corev1.Taint{}, err
	}

	return corev1.Taint{
		Key:    key,
		Value:  value,
		Effect: effect,
	}, nil
}

func validateTaintEffect(effect corev1.TaintEffect) error {
	switch effect {
	case corev1.TaintEffectNoSchedule, corev1.TaintEffectPreferNoSchedule, corev1.TaintEffectNoExecute:
		return nil
	default:
		return fmt.Errorf("invalid taint effect: %v, unsupported taint effect", effect)
	}
}
