package addons

import (
	"bytes"
	"html/template"

	"github.com/kris-nova/logger"
)

const templateName = "cw-template"

// Deployer interface
type Deployer interface {
	CreateOrReplace(manifest []byte, plan bool) error
}

// NewCloudwatchAgent creates a new cloudwatch agent
func NewCloudwatchAgent(client Deployer, clusterName, region string, planMode bool) *CloudwatchAgent {
	cw := &CloudwatchAgent{
		client:      client,
		clusterName: clusterName,
		region:      region,
		planMode:    planMode,
	}
	return cw
}

// A CloudwatchAgent deploys a new cloudwatch agent
type CloudwatchAgent struct {
	client      Deployer
	clusterName string
	region      string
	planMode    bool
}

// Deploy deploys CW agent into given cluster
func (cw *CloudwatchAgent) Deploy() (err error) {
	bArr, err := cwAgentPrometheusEksYamlBytes()
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(templateName).Parse(string(bArr)))
	var out bytes.Buffer
	if err := tmpl.Execute(&out, struct {
		ClusterName string
		Region      string
	}{
		ClusterName: cw.clusterName,
		Region:      cw.region,
	}); err != nil {
		return err
	}
	logger.Debug("install cw agent with template: %s", out.String())
	return cw.client.CreateOrReplace(out.Bytes(), cw.planMode)
}
