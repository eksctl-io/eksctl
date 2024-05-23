package nodebootstrap

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	nodeadm "github.com/awslabs/amazon-eks-ami/nodeadm/api/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
)

// ManagedAL2 is a bootstrapper for managed Amazon Linux 2 nodegroups
type ManagedAL2 struct {
	ng *api.ManagedNodeGroup
	// UserDataMimeBoundary sets the MIME boundary for user data
	UserDataMimeBoundary string
}

// NewManagedAL2Bootstrapper creates a new ManagedAL2 bootstrapper
func NewManagedAL2Bootstrapper(ng *api.ManagedNodeGroup) *ManagedAL2 {
	return &ManagedAL2{
		ng: ng,
	}
}

// UserData returns user data for AL2 managed nodegroups
func (m *ManagedAL2) UserData() (string, error) {
	ng := m.ng

	if strings.HasPrefix(ng.AMI, "ami-") {
		return makeCustomAMIUserData(ng.NodeGroupBase, m.UserDataMimeBoundary)
	}

	var (
		buf       bytes.Buffer
		scripts   []string
		cloudboot []string
	)

	if len(ng.PreBootstrapCommands) > 0 {
		scripts = append(scripts, ng.PreBootstrapCommands...)
	}

	if ng.OverrideBootstrapCommand != nil {
		scripts = append(scripts, *ng.OverrideBootstrapCommand)
	} else if ng.MaxPodsPerNode != 0 {
		scripts = append(scripts, makeMaxPodsScript(ng.MaxPodsPerNode))
	}

	if api.IsEnabled(ng.EFAEnabled) {
		cloudboot = append(cloudboot, assets.EfaManagedBoothook)
	}

	if len(scripts) == 0 && len(cloudboot) == 0 {
		return "", nil
	}

	if err := createMimeMessage(&buf, scripts, cloudboot, nil, m.UserDataMimeBoundary); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func makeCustomAMIUserData(ng *api.NodeGroupBase, mimeBoundary string) (string, error) {
	var (
		buf     bytes.Buffer
		scripts []string
	)

	if len(ng.PreBootstrapCommands) > 0 {
		scripts = append(scripts, ng.PreBootstrapCommands...)
	}

	if ng.OverrideBootstrapCommand != nil {
		scripts = append(scripts, *ng.OverrideBootstrapCommand)
	}

	if len(scripts) == 0 {
		return "", nil
	}

	if err := createMimeMessage(&buf, scripts, nil, nil, mimeBoundary); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func makeMaxPodsScript(maxPods int) string {
	script := `#!/bin/sh
set -ex
`
	script += fmt.Sprintf(`sed -i 's/KUBELET_EXTRA_ARGS=$2/KUBELET_EXTRA_ARGS="$2 --max-pods=%d"/' /etc/eks/bootstrap.sh`, maxPods)
	return script
}

func createMimeMessage(writer io.Writer, scripts, cloudboots []string, nodeConfig *nodeadm.NodeConfig, mimeBoundary string) error {
	mw := multipart.NewWriter(writer)
	if mimeBoundary != "" {
		if err := mw.SetBoundary(mimeBoundary); err != nil {
			return errors.Wrap(err, "unexpected error setting MIME boundary")
		}
	}
	fmt.Fprint(writer, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(writer, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", mw.Boundary())

	for _, script := range scripts {
		part, err := mw.CreatePart(map[string][]string{
			"Content-Type": {"text/x-shellscript", `charset="us-ascii"`},
		})

		if err != nil {
			return err
		}

		if _, err = part.Write([]byte(script)); err != nil {
			return err
		}
	}
	for _, cloudboot := range cloudboots {
		part, err := mw.CreatePart(map[string][]string{
			"Content-Type": {"text/cloud-boothook", `charset="us-ascii"`},
		})

		if err != nil {
			return err
		}

		if _, err = part.Write([]byte(cloudboot)); err != nil {
			return err
		}
	}

	if nodeConfig != nil {
		yamlData, err := yaml.Marshal(nodeConfig)
		if err != nil {
			return fmt.Errorf("error marshalling node configuration: %w", err)
		}
		part, err := mw.CreatePart(map[string][]string{
			"Content-Type": {"application/node.eks.aws"},
		})
		if err != nil {
			return err
		}
		if _, err = part.Write(yamlData); err != nil {
			return err
		}
	}

	return mw.Close()
}
