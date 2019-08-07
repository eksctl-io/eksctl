package install

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
)

type publicKeyInfrastructurePaths struct {
	caCertificate     string
	caKey             string
	serverCertificate string
	serverKey         string
	clientCertificate string
	clientKey         string
}

type publicKeyInfrastructure struct {
	caCertificate     []byte
	caKey             []byte
	serverCertificate []byte
	serverKey         []byte
	clientCertificate []byte
	clientKey         []byte
}

type devnull struct{}

func (devnull) Debug(string)   {}
func (devnull) Info(string)    {}
func (devnull) Warning(string) {}
func (devnull) Err(string)     {}
func (devnull) Crit(string)    {}
func (devnull) Emerg(string)   {}
func init() {
	// Disable logging of cfssl since we handle errors explicitly
	log.SetLogger(devnull{})
}

func newPKI(hostName string, validFor time.Duration, rsaKeyBitSize int) (*publicKeyInfrastructure, error) {
	keyReq := &csr.KeyRequest{
		A: "rsa",
		S: rsaKeyBitSize,
	}

	// Generate CA
	caReq := &csr.CertificateRequest{
		KeyRequest: keyReq,
		CN:         "CA",
	}
	caCert, _, caKey, err := initca.New(caReq)
	if err != nil {
		return nil, fmt.Errorf("cannot generate root CA: %s", err)
	}

	// Generate Server Certificate
	serverCert, serverKey, err := generateCertificate(caCert, caKey, keyReq, hostName, "server", validFor)
	if err != nil {
		return nil, fmt.Errorf("cannot generate server certificate: %s", err)
	}

	// Generate Client certificate
	generateCertificate(caCert, caKey, keyReq, hostName, "client", validFor)
	clientCert, clientKey, err := generateCertificate(caCert, caKey, keyReq, hostName, "client", validFor)
	if err != nil {
		return nil, fmt.Errorf("cannot generate client certificate: %s", err)
	}
	pki := &publicKeyInfrastructure{
		caCertificate:     caCert,
		caKey:             caKey,
		serverCertificate: serverCert,
		serverKey:         serverKey,
		clientCertificate: clientCert,
		clientKey:         clientKey,
	}
	return pki, nil

}

func generateCertificate(caCert []byte, caKey []byte, keyReq *csr.KeyRequest,
	hostName string, commonName string, validFor time.Duration) ([]byte, []byte, error) {
	policy := &config.Signing{
		Default: &config.SigningProfile{
			Expiry: validFor,
			Usage:  []string{"signing", "key encipherment", "server auth", "client auth"},
		},
	}
	parsedCa, err := helpers.ParseCertificatePEM(caCert)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed generated private CA certificate: %s", err)
	}
	priv, err := helpers.ParsePrivateKeyPEM(caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed generated private CA key: %s", err)
	}
	certSigner, err := local.NewSigner(priv, parsedCa, signer.DefaultSigAlgo(priv), policy)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create signer: %s", err)
	}
	serverReq := &csr.CertificateRequest{
		KeyRequest: keyReq,
		CN:         commonName,
		Hosts:      []string{hostName},
	}
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err := g.ProcessRequest(serverReq)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot process server-certificate signing request: %s", err)
	}
	serverSignRequest := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   []string{hostName},
	}
	cert, err := certSigner.Sign(serverSignRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot generate server certificate: %s", err)
	}
	return cert, key, nil
}

func (pki *publicKeyInfrastructure) loadFrom(paths *publicKeyInfrastructurePaths) error {
	for path, destination := range map[string]*[]byte{
		paths.caCertificate:     &pki.caCertificate,
		paths.caKey:             &pki.caKey,
		paths.serverCertificate: &pki.serverCertificate,
		paths.serverKey:         &pki.serverKey,
		paths.clientCertificate: &pki.clientCertificate,
		paths.clientKey:         &pki.clientKey,
	} {
		if len(path) == 0 {
			continue
		}
		var err error
		if *destination, err = ioutil.ReadFile(path); err != nil {
			return fmt.Errorf("cannot read file %q: %s", path, err)
		}
	}
	return nil
}

func (pki *publicKeyInfrastructure) saveTo(paths *publicKeyInfrastructurePaths) error {
	for path, content := range map[string][]byte{
		paths.caCertificate:     pki.caCertificate,
		paths.caKey:             pki.caKey,
		paths.serverCertificate: pki.serverCertificate,
		paths.serverKey:         pki.serverKey,
		paths.clientCertificate: pki.clientCertificate,
		paths.clientKey:         pki.clientKey,
	} {
		if len(path) == 0 {
			continue
		}
		if err := ioutil.WriteFile(path, content, 0400); err != nil {
			return fmt.Errorf("cannot write file %q: %s", path, err)
		}
	}
	return nil
}
