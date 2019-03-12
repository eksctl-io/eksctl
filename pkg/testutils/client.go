package testutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kris-nova/logger"
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

type requestTracker struct {
	requests *[]*http.Request
	missing  *bool
}

func (t *requestTracker) Append(req *http.Request) {
	*t.requests = append(*t.requests, req)
	logger.Critical("requests = %v", t.Methods())
}

func (t *requestTracker) Methods() (m []string) {
	for _, r := range *t.requests {
		m = append(m, r.Method)
	}
	return
}

func (t *requestTracker) IsMissing() bool {
	return *t.missing
}

func NewFakeRawResource(item runtime.Object, missing bool) (*kubernetes.RawResource, requestTracker) {
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
		requests: &[]*http.Request{},
		missing:  &missing,
	}

	client := &restfake.RESTClient{
		GroupVersion:         gvk.GroupVersion(),
		NegotiatedSerializer: scheme.Codecs,
		Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			rt.Append(req)
			switch req.Method {
			case http.MethodGet:
				logger.Critical("missing = %v", rt.IsMissing())
				if rt.IsMissing() {
					res := &http.Response{StatusCode: 404, Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}
					return res, nil
				}
				data, err := runtime.Encode(unstructured.UnstructuredJSONScheme, item)
				Expect(err).To(Not(HaveOccurred()))
				res := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(data))}
				return res, nil
			case http.MethodPut, http.MethodPost:
				*rt.missing = false
				res := &http.Response{StatusCode: 200, Body: req.Body}
				return res, nil
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
