/*
Copyright 2017 by the contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func (c *Config) certPath() string {
	return filepath.Join(c.StateDir, certFilename)
}

func (c *Config) keyPath() string {
	return filepath.Join(c.StateDir, keyFilename)
}

// GetOrCreateCertificate will create a certificate if it cannot find one based on the config
func (c *Config) GetOrCreateCertificate() (*tls.Certificate, error) {
	// first try to load the existing keypair
	cert, err := c.LoadExistingCertificate()
	if err != nil {
		return nil, err
	}
	// if we found it,
	if cert != nil {
		return cert, nil
	}

	// generate a self-signed certificate and write out the certificate and private key
	certBytes, keyBytes, err := c.selfSignCertificate()
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"certPath": c.certPath(),
		"keyPath":  c.keyPath(),
	}).Info("saving new key and certificate")
	err = dumpPEM(c.certPath(), 0666, "CERTIFICATE", certBytes)
	if err != nil {
		return nil, err
	}

	err = dumpPEM(c.keyPath(), 0600, "RSA PRIVATE KEY", keyBytes)
	if err != nil {
		return nil, err
	}

	newCert, err := tls.LoadX509KeyPair(c.certPath(), c.keyPath())
	return &newCert, err
}

// LoadExistingCertificate will load certificates from a local path
func (c *Config) LoadExistingCertificate() (*tls.Certificate, error) {

	// if either file does not exist, we'll consider that not an error but
	// return a nil
	if _, err := os.Stat(c.certPath()); os.IsNotExist(err) {
		return nil, nil
	}
	if _, err := os.Stat(c.keyPath()); os.IsNotExist(err) {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(c.certPath(), c.keyPath())
	if err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"certPath": c.certPath(),
		"keyPath":  c.keyPath(),
	}).Info("loaded existing keypair")
	return &cert, nil
}

func dumpPEM(filename string, mode os.FileMode, blockType string, bytes []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{Type: blockType, Bytes: bytes})
}

func (c *Config) selfSignCertificate() ([]byte, []byte, error) {

	// generate a new RSA-2048 keypair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// choose a beginning and end for the cert's lifetime (currently ~infinite)
	notBefore := time.Now()
	notAfter := notBefore.Add(certLifetime)

	// choose a random 128 bit serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "aws-iam-authenticator",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:        true,
		DNSNames:    []string{c.Hostname},
		IPAddresses: []net.IP{net.ParseIP(c.Address)},
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	logrus.WithFields(logrus.Fields{
		"certBytes": len(certBytes),
		"keyBytes":  len(keyBytes),
	}).Info("generated a new private key and certificate")
	return certBytes, keyBytes, nil
}

// certToPEMBase64 returns the Base64 encoded PEM block for a given DER
// certificate (i.e., it returns "Base64(PEM(asn1))").
func certToPEMBase64(der []byte) string {
	return base64.StdEncoding.EncodeToString(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	}))
}
