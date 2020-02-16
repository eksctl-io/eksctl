// +build integration

package a_test

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
	fmt.Println("BEFORE: a")
	fmt.Println(time.Now())
})

var _ = Describe("a", func() {
	It("1=1", func() {
		fmt.Println("a: 1=1")
		fmt.Printf("%v\n", params)
		time.Sleep(5 * time.Second)
		Expect(1).To(Equal(1))
		fmt.Fprintln(GinkgoWriter, "Almost there...")
	})

	It("2=2", func() {
		fmt.Println("a: 2=2")
		time.Sleep(10 * time.Second)
		Expect(2).To(Equal(2))
	})
})

var _ = AfterSuite(func() {
	fmt.Println("AFTER: a")
})
