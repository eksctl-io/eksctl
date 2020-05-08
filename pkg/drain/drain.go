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
	"time"

	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	// EvictionKind represents the kind of evictions object
	EvictionKind = "Eviction"
	// EvictionSubresource represents the kind of evictions object as pod's subresource
	EvictionSubresource = "pods/eviction"
	// retryDelay is how long is slept before retry after an error occurs during drainage
	retryDelay = 5 * time.Second
)

// Helper contains the parameters to control the behaviour of drainer
type Helper struct {
	Selector    string
	PodSelector string

	Client kubernetes.Interface

	Force  bool
	DryRun bool

	GracePeriodSeconds int
	Timeout            time.Duration

	IgnoreAllDaemonSets bool
	IgnoreDaemonSets    []metav1.ObjectMeta
	DeleteLocalData     bool

	policyAPIGroupVersion string
	UseEvictions          bool
}

// CanUseEvictions uses Discovery API to find out if evictions are supported
func (d *Helper) CanUseEvictions() error {
	discoveryClient := d.Client.Discovery()
	groupList, err := discoveryClient.ServerGroups()
	if err != nil {
		return err
	}
	foundPolicyGroup := false
	for _, group := range groupList.Groups {
		if group.Name == "policy" {
			foundPolicyGroup = true
			d.policyAPIGroupVersion = group.PreferredVersion.GroupVersion
			break
		}
	}
	if !foundPolicyGroup {
		return nil
	}
	resourceList, err := discoveryClient.ServerResourcesForGroupVersion("v1")
	if err != nil {
		return err
	}
	for _, resource := range resourceList.APIResources {
		if resource.Name == EvictionSubresource && resource.Kind == EvictionKind {
			d.UseEvictions = true
			return nil
		}
	}
	return nil
}

func (d *Helper) makeDeleteOptions() *metav1.DeleteOptions {
	deleteOptions := &metav1.DeleteOptions{}
	if d.GracePeriodSeconds >= 0 {
		gracePeriodSeconds := int64(d.GracePeriodSeconds)
		deleteOptions.GracePeriodSeconds = &gracePeriodSeconds
	}
	return deleteOptions
}

// EvictOrDeletePod will evict pod if policy API is available, otherwise deletes it
// NOTE: CanUseEvictions must be called prior to this
func (d *Helper) EvictOrDeletePod(pod corev1.Pod) error {
	if d.UseEvictions {
		return d.EvictPod(pod)
	}
	return d.DeletePod(pod)
}

// EvictPod will evict the give pod, or return an error if it couldn't
// NOTE: CanUseEvictions must be called prior to this
func (d *Helper) EvictPod(pod corev1.Pod) error {
	eviction := &policyv1beta1.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: d.policyAPIGroupVersion,
			Kind:       EvictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: d.makeDeleteOptions(),
	}
	return d.Client.PolicyV1beta1().Evictions(eviction.Namespace).Evict(eviction)
}

// DeletePod will delete the given pod, or return an error if it couldn't
func (d *Helper) DeletePod(pod corev1.Pod) error {
	return d.Client.CoreV1().Pods(pod.Namespace).Delete(pod.Name, d.makeDeleteOptions())
}

// getPodsForDeletion lists all pods on a given node, filters those using the default
// filters, and returns podDeleteList along with any errors. All pods that are ready
// to be deleted can be obtained with .Pods(), and string with all warning can be obtained
// with .Warnings()
func (d *Helper) getPodsForDeletion(nodeName string) (*podDeleteList, []error) {
	labelSelector, err := labels.Parse(d.PodSelector)
	if err != nil {
		return nil, []error{err}
	}

	podList, err := d.Client.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labelSelector.String(),
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": nodeName}).String()})
	if err != nil {
		return nil, []error{err}
	}

	pods := []podDelete{}

	for _, pod := range podList.Items {
		var status podDeleteStatus
		for _, filter := range d.makeFilters() {
			status = filter(pod)
			if !status.delete {
				// short-circuit as soon as pod is filtered out
				// at that point, there is no reason to run pod
				// through any additional filters
				break
			}
		}
		pods = append(pods, podDelete{
			pod:    pod,
			status: status,
		})
	}

	list := &podDeleteList{items: pods}

	if errs := list.errors(); len(errs) > 0 {
		return list, errs
	}

	return list, nil
}
