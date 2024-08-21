package nodebootstrap_test

import (
	"encoding/base64"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

type managedEntry struct {
	ng *api.ManagedNodeGroup

	expectedUserData string
}

var _ = DescribeTable("Managed AL2", func(e managedEntry) {
	api.SetManagedNodeGroupDefaults(e.ng, &api.ClusterMeta{Name: "cluster"}, false)
	bootstrapper := nodebootstrap.NewManagedAL2Bootstrapper(e.ng)
	bootstrapper.UserDataMimeBoundary = "//"

	userData, err := bootstrapper.UserData()
	Expect(err).NotTo(HaveOccurred())
	decoded, err := base64.StdEncoding.DecodeString(userData)
	Expect(err).NotTo(HaveOccurred())
	actual := strings.ReplaceAll(string(decoded), "\r\n", "\n")
	Expect(actual).To(Equal(e.expectedUserData))
},
	Entry("SSM enabled", managedEntry{
		ng: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name: "ssm-enabled",
				SSH: &api.NodeGroupSSH{
					Allow:         api.Enabled(),
					PublicKeyName: aws.String("test-keypair"),
					EnableSSM:     api.Enabled(),
				},
			},
		},

		expectedUserData: "",
	}),

	Entry("Custom AMI with bootstrap script", managedEntry{
		ng: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "custom-ami",
				InstanceType: "m5.xlarge",
				AMI:          "ami-custom",
				PreBootstrapCommands: []string{
					"date",
					"echo hello",
				},
				OverrideBootstrapCommand: aws.String(`#!/bin/bash
set -ex
B64_CLUSTER_CA=dGVzdAo=
API_SERVER_URL=https://test.com
/etc/eks/bootstrap.sh launch-template --b64-cluster-ca $B64_CLUSTER_CA --apiserver-endpoint $API_SERVER_URL
`),
			},
		},

		expectedUserData: `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=//

--//
Content-Type: text/x-shellscript
Content-Type: charset="us-ascii"

date
--//
Content-Type: text/x-shellscript
Content-Type: charset="us-ascii"

echo hello
--//
Content-Type: text/x-shellscript
Content-Type: charset="us-ascii"

#!/bin/bash
set -ex
B64_CLUSTER_CA=dGVzdAo=
API_SERVER_URL=https://test.com
/etc/eks/bootstrap.sh launch-template --b64-cluster-ca $B64_CLUSTER_CA --apiserver-endpoint $API_SERVER_URL

--//--
`,
	}),

	Entry("EFA enabled", managedEntry{
		ng: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:       "ng",
				EFAEnabled: api.Enabled(),
			},
		},

		expectedUserData: `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=//

--//
Content-Type: text/cloud-boothook
Content-Type: charset="us-ascii"

cloud-init-per once yum_wget yum install -y wget
cloud-init-per once wget_efa wget -q --timeout=20 https://s3-us-west-2.amazonaws.com/aws-efa-installer/aws-efa-installer-latest.tar.gz -O /tmp/aws-efa-installer-latest.tar.gz

cloud-init-per once tar_efa tar -xf /tmp/aws-efa-installer-latest.tar.gz -C /tmp
cloud-init-per once rm_efa_gz rm -rf /tmp/aws-efa-installer-latest.tar.gz
pushd /tmp/aws-efa-installer
cloud-init-per once install_efa ./efa_installer.sh -y -g
pop /tmp/aws-efa-installer

cloud-init-per once efa_info /opt/amazon/efa/bin/fi_info -p efa

--//--
`,
	}),

	Entry("EFA and SSM enabled", managedEntry{
		ng: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:       "ng",
				EFAEnabled: api.Enabled(),
				SSH: &api.NodeGroupSSH{
					Allow:     api.Enabled(),
					EnableSSM: api.Enabled(),
				},
			},
		},
		expectedUserData: `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=//

--//
Content-Type: text/cloud-boothook
Content-Type: charset="us-ascii"

cloud-init-per once yum_wget yum install -y wget
cloud-init-per once wget_efa wget -q --timeout=20 https://s3-us-west-2.amazonaws.com/aws-efa-installer/aws-efa-installer-latest.tar.gz -O /tmp/aws-efa-installer-latest.tar.gz

cloud-init-per once tar_efa tar -xf /tmp/aws-efa-installer-latest.tar.gz -C /tmp
cloud-init-per once rm_efa_gz rm -rf /tmp/aws-efa-installer-latest.tar.gz
pushd /tmp/aws-efa-installer
cloud-init-per once install_efa ./efa_installer.sh -y -g
pop /tmp/aws-efa-installer

cloud-init-per once efa_info /opt/amazon/efa/bin/fi_info -p efa

--//--
`,
	}),

	Entry("maxPodsPerNode set", managedEntry{
		ng: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:           "ng",
				MaxPodsPerNode: 142,
			},
		},
		expectedUserData: `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=//

--//
Content-Type: text/x-shellscript
Content-Type: charset="us-ascii"

#!/bin/sh
set -ex
sed -i 's/KUBELET_EXTRA_ARGS=$2/KUBELET_EXTRA_ARGS="$2 --max-pods=142"/' /etc/eks/bootstrap.sh
--//--
`,
	}),
)
