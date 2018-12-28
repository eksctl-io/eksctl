package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("StackCollection Template", func() {
	var (
		nodeLabels NodeLabels
	)

	BeforeEach(func() {
		nodeLabels = NodeLabels{}
	})

	Describe("Type", func() {
		It("should have the fixed, expected value", func() {
			Expect(nodeLabels.Type()).To(Equal("NodeLabels"))
		})
	})

	Describe("Set", func() {
		BeforeEach(func() {
			nodeLabels.Set("k1=v1,k2=v2")
			nodeLabels.Set("k2=v2_,k3=v3")
		})

		It("should parse and merge key-value pairs", func() {
			Expect(strings.Split(nodeLabels.String(), ",")).To(ConsistOf("k3=v3", "k2=v2_", "k1=v1"))
		})
	})
})
