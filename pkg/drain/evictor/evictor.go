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

package evictor

import (
	"context"
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
	// EvictionSubresource represents the kind of evictions object as Pod's subresource
	EvictionSubresource = "pods/eviction"
)

// Evictor contains the parameters to control the behaviour of the evictor
type Evictor struct {
	podSelector string

	client kubernetes.Interface

	force  bool
	DryRun bool

	maxGracePeriodSeconds int

	ignoreAllDaemonSets bool
	ignoreDaemonSets    []metav1.ObjectMeta
	deleteLocalData     bool
	disableEviction     bool

	policyAPIGroupVersion string
	UseEvictions          bool
}

func New(clientSet kubernetes.Interface, maxGracePeriod time.Duration, ignoreDaemonSets []metav1.ObjectMeta, disableEviction bool) *Evictor {
	return &Evictor{
		client: clientSet,
		// TODO: force, DeleteLocalData & IgnoreAllDaemonSets shouldn't
		// be enabled by default, we need flags to control these, but that
		// requires more improvements in the underlying drain package,
		// as it currently produces errors and warnings with references
		// to kubectl flags
		force:                 true,
		deleteLocalData:       true,
		ignoreAllDaemonSets:   true,
		maxGracePeriodSeconds: int(maxGracePeriod.Seconds()),
		ignoreDaemonSets:      ignoreDaemonSets,
		disableEviction:       disableEviction,
	}
}

// CanUseEvictions uses Discovery API to find out if evictions are supported
func (d *Evictor) CanUseEvictions() error {
	if d.disableEviction {
		d.UseEvictions = false
		return nil
	}
	discoveryClient := d.client.Discovery()
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

func (d *Evictor) makeDeleteOptions(pod corev1.Pod) metav1.DeleteOptions {
	deleteOptions := metav1.DeleteOptions{}

	gracePeriodSeconds := int64(corev1.DefaultTerminationGracePeriodSeconds)
	if pod.Spec.TerminationGracePeriodSeconds != nil {
		if *pod.Spec.TerminationGracePeriodSeconds < int64(d.maxGracePeriodSeconds) {
			gracePeriodSeconds = *pod.Spec.TerminationGracePeriodSeconds
		} else {
			gracePeriodSeconds = int64(d.maxGracePeriodSeconds)
		}
	}

	deleteOptions.GracePeriodSeconds = &gracePeriodSeconds
	return deleteOptions
}

// EvictOrDeletePod will evict Pod if policy API is available, otherwise deletes it. If disableEviction is true, we skip straight to the delete step
// NOTE: CanUseEvictions must be called prior to this
func (d *Evictor) EvictOrDeletePod(pod corev1.Pod) error {
	if d.UseEvictions {
		return d.evictPod(pod)
	}
	return d.deletePod(pod)
}

// evictPod will evict the give Pod, or return an error if it couldn't
// NOTE: CanUseEvictions must be called prior to this
func (d *Evictor) evictPod(pod corev1.Pod) error {
	deleteOptions := d.makeDeleteOptions(pod)
	eviction := &policyv1beta1.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: d.policyAPIGroupVersion,
			Kind:       EvictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: &deleteOptions,
	}
	return d.client.PolicyV1beta1().Evictions(eviction.Namespace).Evict(context.TODO(), eviction)
}

// deletePod will Delete the given Pod, or return an error if it couldn't
func (d *Evictor) deletePod(pod corev1.Pod) error {
	return d.client.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, d.makeDeleteOptions(pod))
}

// GetPodsForEviction lists all pods on a given node, filters those using the default
// filters, and returns PodDeleteList along with any errors. All pods that are ready
// to be deleted can be obtained with .Pods(), and string with all warning can be obtained
// with .Warnings()
func (d *Evictor) GetPodsForEviction(nodeName string) (*PodDeleteList, []error) {
	labelSelector, err := labels.Parse(d.podSelector)
	if err != nil {
		return nil, []error{err}
	}

	podList, err := d.client.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector.String(),
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": nodeName}).String()})
	if err != nil {
		return nil, []error{err}
	}

	pods := []PodDelete{}

	for _, pod := range podList.Items {
		var status PodDeleteStatus
		for _, filter := range d.makeFilters() {
			status = filter(pod)
			if !status.Delete {
				// short-circuit as soon as Pod is filtered out
				// at that point, there is no Reason to run Pod
				// through any additional filters
				break
			}
		}
		pods = append(pods, PodDelete{
			Pod:    pod,
			Status: status,
		})
	}

	list := &PodDeleteList{Items: pods}

	if errs := list.errors(); len(errs) > 0 {
		return list, errs
	}

	return list, nil
}
