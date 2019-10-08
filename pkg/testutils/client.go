package testutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/kubernetes"

	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"

	restfake "k8s.io/client-go/rest/fake"
)

func LoadSamples(manifest string) []runtime.Object {
	var (
		samples []runtime.Object
	)

	samplesData, err := ioutil.ReadFile(manifest)
	Expect(err).ToNot(HaveOccurred())
	samplesList, err := kubernetes.NewList(samplesData)
	Expect(err).ToNot(HaveOccurred())
	Expect(samplesList).ToNot(BeNil())

	for _, item := range samplesList.Items {
		kind := item.Object.GetObjectKind().GroupVersionKind().Kind
		if kind == "CustomResourceDefinition" {
			continue // fake client doesn't support CRDs, save it from a panic
		}
		samples = append(samples, item.Object)
	}

	return samples
}

func NewFakeClientSetWithSamples(manifest string) (*fake.Clientset, []runtime.Object) {
	samples := LoadSamples(manifest)
	return fake.NewSimpleClientset(samples...), samples
}

var mapper = testrestmapper.TestOnlyStaticRESTMapper(scheme.Scheme)

type CollectionTracker struct {
	created map[string]runtime.Object
	updated map[string]runtime.Object
	deleted map[string]runtime.Object
	objects map[string]runtime.Object
}

func NewCollectionTracker() *CollectionTracker {
	return &CollectionTracker{
		created: make(map[string]runtime.Object),
		updated: make(map[string]runtime.Object),
		deleted: make(map[string]runtime.Object),
		objects: make(map[string]runtime.Object),
	}
}

type requestTracker struct {
	requests   *[]*http.Request
	missing    *bool
	unionised  *bool
	collection *CollectionTracker
}

func objectReqKey(req *http.Request, item runtime.Object) string {
	return fmt.Sprintf("%s [%s] (%s)",
		req.Method, req.URL.Path, item.(metav1.Object).GetName())
}

func objectKey(req *http.Request, item runtime.Object) string {
	switch req.Method {
	case http.MethodPost:
		return fmt.Sprintf("%s/%s", req.URL.Path, item.(metav1.Object).GetName())
	case http.MethodGet, http.MethodPut, http.MethodDelete:
		return fmt.Sprintf("%s", req.URL.Path)
	}
	return fmt.Sprintf("%s [%s] (%s)",
		req.Method, req.URL.Path, item.(metav1.Object).GetName())
}

func (t *requestTracker) Append(req *http.Request) { *t.requests = append(*t.requests, req) }

func (t *requestTracker) Methods() (m []string) {
	for _, r := range *t.requests {
		m = append(m, r.Method)
	}
	return
}

func (t *requestTracker) IsMissing(req *http.Request, item runtime.Object) bool {
	if *t.unionised && t.collection != nil {
		k := objectKey(req, item)
		_, ok := t.collection.objects[k]
		return !ok
	}
	return *t.missing
}

func (t *requestTracker) Create(req *http.Request, item runtime.Object) bool {
	*t.missing = false
	if t.collection != nil {
		t.collection.created[objectReqKey(req, item)] = item
		if *t.unionised {
			k := objectKey(req, item)
			if _, ok := t.collection.objects[k]; ok {
				return false
			}
			t.collection.objects[k] = item
		}
	}
	return true
}

func (c *CollectionTracker) Created() map[string]runtime.Object { return c.created }

func (c *CollectionTracker) CreatedItems() (items []runtime.Object) {
	for _, item := range c.Created() {
		items = append(items, item)
	}
	return
}

func (t *requestTracker) Update(req *http.Request, item runtime.Object) {
	if t.collection != nil {
		t.collection.updated[objectReqKey(req, item)] = item
		if *t.unionised {
			k := objectKey(req, item)
			t.collection.objects[k] = item
		}
	}
}

func (c *CollectionTracker) Updated() map[string]runtime.Object { return c.updated }

func (c *CollectionTracker) UpdatedItems() (items []runtime.Object) {
	for _, item := range c.Updated() {
		items = append(items, item)
	}
	return
}

func (t *requestTracker) Delete(req *http.Request, item runtime.Object) bool {
	*t.missing = true
	if t.collection != nil {
		t.collection.deleted[objectReqKey(req, item)] = item
		if *t.unionised {
			k := objectKey(req, item)
			if _, ok := t.collection.objects[k]; ok {
				delete(t.collection.objects, k)
				return true
			}
			return false
		}
	}
	return true
}

func (c *CollectionTracker) Deleted() map[string]runtime.Object { return c.deleted }

func (c *CollectionTracker) DeletedItems() (items []runtime.Object) {
	for _, item := range c.Deleted() {
		items = append(items, item)
	}
	return
}

func (c *CollectionTracker) AllTracked() map[string]runtime.Object { return c.objects }

func (c *CollectionTracker) AllTrackedItmes() (items []runtime.Object) {
	for _, item := range c.AllTracked() {
		items = append(items, item)
	}
	return
}

func NewFakeRawResource(item runtime.Object, missing, unionised bool, ct *CollectionTracker) (*kubernetes.RawResource, requestTracker) {
	obj, ok := item.(metav1.Object)
	Expect(ok).To(BeTrue())

	gvk := item.GetObjectKind().GroupVersionKind()

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	Expect(err).To(Not(HaveOccurred()))

	info := &resource.Info{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
		Object:    item,
	}

	rt := requestTracker{
		requests:   &[]*http.Request{},
		missing:    &missing,
		unionised:  &unionised,
		collection: ct,
	}

	emptyBody := ioutil.NopCloser(bytes.NewReader([]byte{}))
	notFound := http.Response{StatusCode: http.StatusNotFound, Body: emptyBody}
	conflict := http.Response{StatusCode: http.StatusConflict, Body: emptyBody}

	echo := func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: req.Body}, nil
	}

	asResult := func(req *http.Request) (*http.Response, error) {
		data, err := runtime.Encode(unstructured.UnstructuredJSONScheme, item)
		Expect(err).To(Not(HaveOccurred()))
		res := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader(data))}
		return res, nil
	}

	client := &restfake.RESTClient{
		GroupVersion:         gvk.GroupVersion(),
		NegotiatedSerializer: scheme.Codecs,
		Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			rt.Append(req)
			switch req.Method {
			case http.MethodGet:
				if rt.IsMissing(req, item) {
					return &notFound, nil
				}
				return asResult(req)
			case http.MethodPost:
				if !rt.Create(req, item) {
					return &conflict, nil
				}
				return echo(req)
			case http.MethodPut:
				rt.Update(req, item)
				return echo(req)
			case http.MethodDelete:
				if !rt.Delete(req, item) {
					return &notFound, nil
				}
				return asResult(req)
			default:
				return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
			}
		}),
	}

	helper := resource.NewHelper(client, mapping)

	rc := &kubernetes.RawResource{
		Helper: helper,
		Info:   info,
		GVK:    &gvk,
	}

	return rc, rt
}

type FakeRawClient struct {
	Collection                 *CollectionTracker
	AssumeObjectsMissing       bool
	ClientSetUseUpdatedObjects bool
	UseUnionTracker            bool
}

func NewFakeRawClient() *FakeRawClient {
	return &FakeRawClient{
		Collection: NewCollectionTracker(),
	}
}

func (c *FakeRawClient) ClientSet() kubeclient.Interface {
	if c.UseUnionTracker {
		// TODO: try to use clientSet.Fake.Actions, clientSet.Fake.PrependReactor
		// or any of the other hooks to connect this clientset instance with
		// udnerlying CollectionTracker, so that we get proper end-to-end behaviour
		return fake.NewSimpleClientset(c.Collection.AllTrackedItmes()...)
	}
	if c.ClientSetUseUpdatedObjects {
		return fake.NewSimpleClientset(c.Collection.UpdatedItems()...)
	}
	return fake.NewSimpleClientset(c.Collection.CreatedItems()...)
}

func (c *FakeRawClient) NewRawResource(object runtime.Object) (*kubernetes.RawResource, error) {
	r, _ := NewFakeRawResource(object, c.AssumeObjectsMissing, c.UseUnionTracker, c.Collection)
	return r, nil
}

func (c *FakeRawClient) ClearUpdated() {
	for k := range c.Collection.updated {
		delete(c.Collection.updated, k)
	}
}
