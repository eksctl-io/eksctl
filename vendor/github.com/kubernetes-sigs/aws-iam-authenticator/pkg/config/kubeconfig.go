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
	"os"
	"text/template"
)

var webhookKubeconfigTemplate = template.Must(
	template.New("kubeconfig").Option("missingkey=error").Parse(`
# clusters refers to the remote service.
clusters:
  - name: aws-iam-authenticator
    cluster:
      certificate-authority-data: {{.CertificateAuthorityBase64}}
      server: {{.ServerURL}}
# users refers to the API Server's webhook configuration
# (we don't need to authenticate the API server).
users:
  - name: apiserver
# kubeconfig files require a context. Provide one for the API Server.
current-context: webhook
contexts:
- name: webhook
  context:
    cluster: aws-iam-authenticator
    user: apiserver
`))

type kubeconfigParams struct {
	ServerURL                  string
	CertificateAuthorityBase64 string
}

func (p kubeconfigParams) writeTo(outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return webhookKubeconfigTemplate.Execute(f, p)
}
