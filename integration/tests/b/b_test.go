// +build integration

package b_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("")
}

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	fmt.Fprintln(GinkgoWriter, ">>> B")
	fmt.Println("BEFORE: b")
	fmt.Printf("%v\n", params)
	fmt.Println(time.Now())
})

var _ = Describe("b", func() {
	It("3=3", func() {
		fmt.Fprintln(GinkgoWriter, "b: 3=3")
		time.Sleep(5 * time.Second)
		Expect(3).To(Equal(3))
	})

	It("4=4", func() {
		fmt.Println("b: 4=4")
		time.Sleep(10 * time.Second)
		Expect(4).To(Equal(4))
	})
})

var _ = AfterSuite(func() {
	fmt.Println("AFTER: b")
})
