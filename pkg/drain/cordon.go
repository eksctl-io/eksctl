/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package drain

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
)

// CordonHelper wraps functionality to cordon/uncordon nodes
type CordonHelper struct {
	node   *corev1.Node
	status DesiredCordonStatus
}

type DesiredCordonStatus string

const (
	CordonNode   DesiredCordonStatus = "cordon"
	UncordonNode DesiredCordonStatus = "uncordon"
)

func (n DesiredCordonStatus) String() string {
	return string(n)
}

// NewCordonHelper returns a new CordonHelper
func NewCordonHelper(node *corev1.Node, desired DesiredCordonStatus) *CordonHelper {
	return &CordonHelper{
		node:   node,
		status: desired,
	}
}

// NewCordonHelperFromRuntimeObject returns a new CordonHelper, or an error if given object is not a
// node or cannot be encoded as JSON
func NewCordonHelperFromRuntimeObject(nodeObject runtime.Object, scheme *runtime.Scheme, gvk schema.GroupVersionKind, desired DesiredCordonStatus) (*CordonHelper, error) {
	nodeObject, err := scheme.ConvertToVersion(nodeObject, gvk.GroupVersion())
	if err != nil {
		return nil, err
	}

	node, ok := nodeObject.(*corev1.Node)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", nodeObject)
	}

	return NewCordonHelper(node, desired), nil
}

// IsUpdateRequired returns true if c.node.Spec.Unschedulable matches desired state,
// or false when it is
func (c *CordonHelper) IsUpdateRequired() bool {
	mustCordon := !c.node.Spec.Unschedulable && c.status == CordonNode

	mustUncordon := c.node.Spec.Unschedulable && c.status == UncordonNode

	return mustCordon || mustUncordon
}

// PatchOrReplace uses given clientset to update the node status, either by patching or
// updating the given node object; it may return error if the object cannot be encoded as
// JSON, or if either patch or update calls fail; it will also return a second error
// whenever creating a patch has failed
func (c *CordonHelper) PatchOrReplace(clientset kubernetes.Interface) (error, error) {
	client := clientset.Core().Nodes()

	oldData, err := json.Marshal(c.node)
	if err != nil {
		return err, nil
	}

	switch c.status {
	case CordonNode:
		c.node.Spec.Unschedulable = true
	case UncordonNode:
		c.node.Spec.Unschedulable = false
	}

	newData, err := json.Marshal(c.node)
	if err != nil {
		return err, nil
	}

	patchBytes, patchErr := strategicpatch.CreateTwoWayMergePatch(oldData, newData, c.node)
	if patchErr == nil {
		_, err = client.Patch(c.node.Name, types.StrategicMergePatchType, patchBytes)
	} else {
		_, err = client.Update(c.node)
	}
	return err, patchErr
}
