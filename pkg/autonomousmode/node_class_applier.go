package autonomousmode

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/kris-nova/logger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

type EKSNodeClassApplier struct {
	DynamicClient *dynamic.DynamicClient
}

func (a *EKSNodeClassApplier) PatchSubnets(ctx context.Context, subnetIDs []string) error {
	const eksNodeClassName = "default"
	if err := a.patchSubnets(ctx, eksNodeClassName, subnetIDs); err != nil {
		return fmt.Errorf("patching %q EKSNodeClass to use subnets %v: %w", eksNodeClassName, subnetIDs, err)
	}
	return nil
}

func (a *EKSNodeClassApplier) patchSubnets(ctx context.Context, eksNodeClassName string, subnetIDs []string) error {
	const eksNodeClassResourceName = "nodeclasses"
	eksNodeClasses := a.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    "eks.amazonaws.com",
		Version:  "v1",
		Resource: eksNodeClassResourceName,
	})

	eksNodeClass, err := backoff.RetryWithData(func() (*unstructured.Unstructured, error) {
		eksNodeClass, err := eksNodeClasses.Get(ctx, eksNodeClassName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return nil, fmt.Errorf("%s %q does not exist", eksNodeClassResourceName, eksNodeClassName)
			}
			return nil, fmt.Errorf("getting %s/%s: %w", eksNodeClassResourceName, eksNodeClassName, err)
		}
		return eksNodeClass, nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(3*time.Second), 20))
	// TODO: timeouts and retries.
	if err != nil {
		return err
	}
	subnetSelectorTerms, found, err := unstructured.NestedSlice(eksNodeClass.UnstructuredContent(), "spec", "subnetSelectorTerms")
	if err != nil {
		return err
	}
	if !found {
		logger.Warning("expected to find spec.subnetSelectorTerms in %s", eksNodeClassResourceName)
	} else {
		logger.Debug("existing subnetSelectorTerms: %v", subnetSelectorTerms)
	}

	type subnetSelector struct {
		ID string `json:"id,omitempty"`
	}
	var subnetSelectors []subnetSelector
	for _, subnetID := range subnetIDs {
		subnetSelectors = append(subnetSelectors, subnetSelector{
			ID: subnetID,
		})
	}
	patch := []struct {
		Op    string           `json:"op"`
		Path  string           `json:"path"`
		Value []subnetSelector `json:"value"`
	}{
		{
			Op:    "replace",
			Path:  "/spec/subnetSelectorTerms",
			Value: subnetSelectors,
		},
	}
	patchData, err := json.Marshal(patch)
	if err != nil {
		return fmt.Errorf("error creating patch: %w", err)
	}
	if _, err := eksNodeClasses.Patch(ctx, eksNodeClassName, types.JSONPatchType, patchData, metav1.PatchOptions{
		FieldManager: "eksctl",
	}); err != nil {
		return fmt.Errorf("patching %s %q to use subnets %v: %w", eksNodeClassResourceName, eksNodeClassName, subnetIDs, err)
	}
	return nil
}
