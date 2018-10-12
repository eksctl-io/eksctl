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
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func (c *Config) ListenURL() string {
	return fmt.Sprintf("https://%s/authenticate", c.ListenAddr())
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.HostPort)
}

func (c *Config) GenerateFiles() error {
	// load or generate a certificate+private key
	_, err := c.GetOrCreateCertificate()
	if err != nil {
		return fmt.Errorf("could not load/generate a certificate")
	}
	err = c.CreateKubeconfig()
	if err != nil {
		return fmt.Errorf("could not generate a webhook kubeconfig")
	}
	return nil
}

func (c *Config) CreateKubeconfig() error {
	cert, err := c.LoadExistingCertificate()
	if err != nil {
		return fmt.Errorf("failed to load an existing certificate: %v", err)
	}

	// write a kubeconfig suitable for the API server to call us
	logrus.WithField("kubeconfigPath", c.GenerateKubeconfigPath).Info("writing webhook kubeconfig file")
	err = kubeconfigParams{
		ServerURL:                  c.ListenURL(),
		CertificateAuthorityBase64: certToPEMBase64(cert.Certificate[0]),
	}.writeTo(c.GenerateKubeconfigPath)
	if err != nil {
		logrus.WithField("kubeconfigPath", c.GenerateKubeconfigPath).WithError(err).Fatal("could not write kubeconfig")
	}
	return nil
}

// CertPath returns the path to the pem file containing the certificate
func (c *Config) CertPath() string {
	return filepath.Join(c.StateDir, "cert.pem")
}

// KeyPath returns the path to the pem file containing the private key
func (c *Config) KeyPath() string {
	return filepath.Join(c.StateDir, "key.pem")
}
