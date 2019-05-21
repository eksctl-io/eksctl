package nodebootstrap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"strings"
)

var _ = Describe("User data", func() {
	Describe("generating max pods", func() {
		It("max pods mapping has the correct format", func() {
			maxPods := makeMaxPodsMapping()
			lines := strings.Split(strings.TrimSpace(maxPods), "\n")
			for _, line := range lines {
				parts := strings.Split(line, " ")
				Expect(parts[0]).To(MatchRegexp(`[a-zA-Z][a-zA-Z0-9]*\.[a-zA-Z0-9]+`))
				_, err := strconv.Atoi(parts[1])
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})
})
