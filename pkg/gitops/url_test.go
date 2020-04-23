package gitops_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/gitops"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("profile", func() {
	Describe("RepositoryURL", func() {
		It("returns Git URLs as-is", func() {
			url, err := gitops.RepositoryURL("https://github.com/eksctl-bot/my-gitops-repo")
			Expect(err).To(Not(HaveOccurred()))
			Expect(url).To(Equal("https://github.com/eksctl-bot/my-gitops-repo"))
		})

		It("returns full Git URLs for supported mnemonics", func() {
			mnemonicToURLs := []struct {
				mnemonic string
				url      string
			}{
				{mnemonic: "app-dev", url: "https://github.com/weaveworks/eks-quickstart-app-dev"},
				{mnemonic: "appmesh", url: "https://github.com/weaveworks/eks-appmesh-profile"},
			}
			for _, mnemonicToURL := range mnemonicToURLs {
				url, err := gitops.RepositoryURL(mnemonicToURL.mnemonic)
				Expect(err).To(Not(HaveOccurred()))
				Expect(url).To(Equal(mnemonicToURL.url))
			}
		})

		It("returns an error otherwise", func() {
			url, err := gitops.RepositoryURL("foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid URL or unknown Quick Start profile: foo"))
			Expect(url).To(Equal(""))
		})
	})
})
