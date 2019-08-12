package clusterautoscaler

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	defaultNamespace    = "kube-system"
	defaultImageVersion = "v1.12.3"
)

// TemplateParameters groups the parameters required to configure cluster
// autoscaling.
// N.B.: these parameters are injected into the cluster autoscaler's YAML
// template, hence their names need to be consistent both here and in the
// template.
type TemplateParameters struct {
	Namespace    string // Namespace under which cluster autoscaler resources will be created.
	ClusterName  string // Name of the EKS cluster, as configured in eksctl's manifest or via CLI flags.
	ImageVersion string // Version of the k8s.gcr.io/cluster-autoscaler image.
}

// Validate validates the values currently set and defaults the ones missing.
func (p *TemplateParameters) Validate() error {
	if p.ClusterName == "" {
		return errors.New("blank cluster name")
	}
	if p.Namespace == "" {
		p.Namespace = defaultNamespace
	}
	if p.ImageVersion == "" {
		p.ImageVersion = defaultImageVersion
	}
	return nil
}

// GenerateManifests fills the cluster autoscaler's YAML templates with the
// provided template parameters to generate actionable YAML manifests.
func GenerateManifests(params TemplateParameters) ([]byte, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	manifests := [][]byte{}
	err := vfsutil.WalkFiles(Templates, "/", func(path string, info os.FileInfo, rs io.ReadSeeker, err error) error {
		if err != nil {
			return fmt.Errorf("cannot walk embedded files: %s", err)
		}
		if info.IsDir() {
			return nil
		}
		manifestTemplateBytes, err := ioutil.ReadAll(rs)
		if err != nil {
			return fmt.Errorf("cannot read embedded file %q: %s", info.Name(), err)
		}
		manifestTemplate, err := template.New(info.Name()).
			Funcs(template.FuncMap{"StringsJoin": strings.Join}).
			Parse(string(manifestTemplateBytes))
		if err != nil {
			return fmt.Errorf("cannot parse embedded file %q: %s", info.Name(), err)
		}
		out := bytes.NewBuffer(nil)
		if err := manifestTemplate.Execute(out, params); err != nil {
			return fmt.Errorf("cannot execute template for embedded file %q: %s", info.Name(), err)
		}
		manifests = append(manifests, out.Bytes())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("internal error filling embedded installation templates: %s", err)
	}
	return kubernetes.ConcatManifests(manifests...), nil
}
