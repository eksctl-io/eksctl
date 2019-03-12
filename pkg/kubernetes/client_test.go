package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("default addons", func() {
	Describe("can create or replace missing objects", func() {
		It("can update objects that already exist", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.10.json")

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, false)
				_, err := rc.CreateOrReplace()
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "GET", "PUT"}))
			}
		})

		It("can create objects that don't exist yet", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.10.json")

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, true)
				_, err := rc.CreateOrReplace()
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "POST"}))
			}
		})
	})
})
