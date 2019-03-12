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
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, false, ct)
				_, err := rc.CreateOrReplace()
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "GET", "PUT"}))
			}

			Expect(ct.Updated()).ToNot(BeEmpty())
			Expect(ct.UpdatedItems()).To(HaveLen(6))
			Expect(ct.Created()).To(BeEmpty())
			Expect(ct.CreatedItems()).To(BeEmpty())
		})

		It("can create objects that don't exist yet", func() {
			sampleAddons := testutils.LoadSamples("../addons/default/testdata/sample-1.10.json")
			ct := testutils.NewCollectionTracker()

			for _, item := range sampleAddons {
				rc, track := testutils.NewFakeRawResource(item, true, ct)
				_, err := rc.CreateOrReplace()
				Expect(err).ToNot(HaveOccurred())
				Expect(track).ToNot(BeNil())
				Expect(track.Methods()).To(Equal([]string{"GET", "POST"}))
			}

			Expect(ct.Created()).ToNot(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(6))
			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.UpdatedItems()).To(BeEmpty())
		})
	})
})
