package install

import (
	"testing"
	"time"
)

func TestNewPKI(t *testing.T) {
	pki, err := newPKI("foo.com", 5*365*24*time.Hour, 2048) // shorter key, to save time
	if err != nil {
		t.Fatalf("PKI creation failed: %s", err)
	}
	for _, content := range [][]byte{
		pki.caCertificate, pki.caKey, pki.serverCertificate,
		pki.serverCertificate, pki.serverKey,
		pki.clientCertificate, pki.clientKey,
	} {
		if len(content) == 0 {
			t.Fatalf("empty content found")
		}
	}
}
