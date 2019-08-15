package flux

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PKI", func() {
	pki, err := newPKI("foo.com", 5*365*24*time.Hour, 2048) // shorter key, to save time
	It("should not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})
	for _, content := range [][]byte{
		pki.caCertificate, pki.caKey, pki.serverCertificate,
		pki.serverCertificate, pki.serverKey,
		pki.clientCertificate, pki.clientKey,
	} {

		It("should not be empty", func() {
			Expect(len(content)).NotTo(Equal(0))
		})
	}
})
