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
}

func NewCollectionTracker() *CollectionTracker {
	return &CollectionTracker{
		created: make(map[string]runtime.Object),
		updated: make(map[string]runtime.Object),
	}
}

type requestTracker struct {
	requests   *[]*http.Request
	missing    *bool
	collection *CollectionTracker
}

func objectKey(req *http.Request, item runtime.Object) string {
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

func (t *requestTracker) IsMissing() bool { return *t.missing }

func (t *requestTracker) Create(req *http.Request, item runtime.Object) {
	*t.missing = false
	if t.collection != nil {
		t.collection.created[objectKey(req, item)] = item
	}
}

func (c *CollectionTracker) Created() map[string]runtime.Object { return c.created }

func (c *CollectionTracker) CreatedItems() (items []runtime.Object) {
	for _, item := range c.created {
		items = append(items, item)
	}
	return
}

func (t *requestTracker) Update(req *http.Request, item runtime.Object) {
	if t.collection != nil {
		t.collection.updated[objectKey(req, item)] = item
	}
}

func (c *CollectionTracker) Updated() map[string]runtime.Object { return c.updated }

func (c *CollectionTracker) UpdatedItems() (items []runtime.Object) {
	for _, item := range c.updated {
		items = append(items, item)
	}
	return
}

func NewFakeRawResource(item runtime.Object, missing bool, ct *CollectionTracker) (*kubernetes.RawResource, requestTracker) {
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
		collection: ct,
	}

	notFound := http.Response{StatusCode: 404, Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}

	echo := func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: req.Body}, nil
	}

	client := &restfake.RESTClient{
		GroupVersion:         gvk.GroupVersion(),
		NegotiatedSerializer: scheme.Codecs,
		Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			rt.Append(req)
			switch req.Method {
			case http.MethodGet:
				if rt.IsMissing() {
					return &notFound, nil
				}
				data, err := runtime.Encode(unstructured.UnstructuredJSONScheme, item)
				Expect(err).To(Not(HaveOccurred()))
				res := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(data))}
				return res, nil
			case http.MethodPost:
				rt.Create(req, item)
				return echo(req)
			case http.MethodPut:
				rt.Update(req, item)
				return echo(req)
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
