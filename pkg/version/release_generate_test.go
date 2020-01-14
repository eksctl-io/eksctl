// +build release
package main

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/version"
)

func TestRelease(t *testing.T) {
	version.Version = "0.5.0"
	version.PreReleaseId = "dev"
	v, p := prepareRelease()

	Expect(v).To(Equal("0.5.0"))
	Expect(p).To(BeEmpty())
}
