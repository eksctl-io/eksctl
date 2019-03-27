package create_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"

	. "github.com/weaveworks/eksctl/pkg/ctl/create"
)

var _ = Describe("create nodegroup utils", func() {

	Context("apply nodegroup defaults", func() {

		It("should iterate over unique nodegroups and apply defaults with NewNodeGroupChecker", func() {
			cfg := newClusterConfig()
			addGroupA(cfg)
			addGroupB(cfg)

			filter := cmdutils.NewNodeGroupFilter()
			printer := printers.NewJSONPrinter()
			names := []string{}

			filter.CheckEachNodeGroup(cfg.NodeGroups, SetNodeGroupDefaults)

			filter.CheckEachNodeGroup(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				Expect(nodeGroup).To(Equal(cfg.NodeGroups[i]))
				names = append(names, nodeGroup.Name)
				return nil
			})
			Expect(names).To(Equal([]string{"test-ng1a", "test-ng2a", "test-ng3a", "test-ng1b", "test-ng2b", "test-ng3b"}))

			w := &bytes.Buffer{}

			printer.PrintObj(cfg, w)

			Expect(w.Bytes()).To(MatchJSON(expected))
		})
	})
})

func newClusterConfig() *api.ClusterConfig {
	cfg := api.NewClusterConfig()

	cfg.Metadata.Name = "test-3x3-ngs"
	cfg.Metadata.Region = "eu-central-1"

	return cfg
}

func addGroupA(cfg *api.ClusterConfig) {
	var ng *api.NodeGroup

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng1a"
	ng.VolumeSize = 768
	ng.VolumeType = "io1"
	ng.IAM.AttachPolicyARNs = []string{"foo"}
	ng.Labels = map[string]string{"group": "a", "seq": "1"}

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng2a"
	ng.IAM.AttachPolicyARNs = []string{"bar"}
	ng.Labels = map[string]string{"group": "a", "seq": "2"}

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng3a"
	ng.ClusterDNS = "1.2.3.4"
	ng.InstanceType = "m3.large"
	ng.AllowSSH = true
	ng.Labels = map[string]string{"group": "a", "seq": "3"}
}

func addGroupB(cfg *api.ClusterConfig) {
	var ng *api.NodeGroup

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng1b"
	ng.AllowSSH = true
	ng.Labels = map[string]string{"group": "b", "seq": "1"}

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng2b"
	ng.ClusterDNS = "4.2.8.14"
	ng.InstanceType = "m5.xlarge"
	ng.SecurityGroups.AttachIDs = []string{"sg-1", "sg-2"}
	ng.SecurityGroups.WithLocal = api.NewBoolFalse()
	ng.Labels = map[string]string{"group": "b", "seq": "1"}

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng3b"
	ng.VolumeSize = 192
	ng.SecurityGroups.AttachIDs = []string{"sg-1", "sg-2"}
	ng.SecurityGroups.WithLocal = api.NewBoolFalse()
	ng.Labels = map[string]string{"group": "b", "seq": "1"}
}

const expected = `
  {
		"kind": "ClusterConfig",
		"apiVersion": "eksctl.io/v1alpha4",
		"metadata": {
		  "name": "test-3x3-ngs",
		  "region": "eu-central-1",
		  "version": "1.11"
		},
		"iam": {},
		"vpc": {
		  "cidr": "192.168.0.0/16"
		},
		"nodeGroups": [
		  {
			  "name": "test-ng1a",
			  "ami": "static",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 768,
			  "volumeType": "io1",
			  "labels": {
			    "group": "a",
			    "seq": "1"
			  },
			  "allowSSH": false,
			  "iam": {
			    "attachPolicyARNs": [
			  	"foo"
			    ],
			    "withAddonPolicies": {
			  	"imageBuilder": false,
			  	"autoScaler": false,
			  	"externalDNS": false,
			  	"appMesh": false,
			  	"ebs": false
			    }
			  }
		  },
		  {
			  "name": "test-ng2a",
			  "ami": "static",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 0,
			  "volumeType": "gp2",
			  "labels": {
			    "group": "a",
			    "seq": "2"
			  },
			  "allowSSH": false,
			  "iam": {
			    "attachPolicyARNs": [
			  	"bar"
			    ],
			    "withAddonPolicies": {
			  	"imageBuilder": false,
			  	"autoScaler": false,
			  	"externalDNS": false,
			  	"appMesh": false,
			  	"ebs": false
			    }
			  }
		  },
		  {
			  "name": "test-ng3a",
			  "ami": "static",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m3.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 0,
			  "volumeType": "gp2",
			  "labels": {
			    "group": "a",
			    "seq": "3"
			  },
			  "allowSSH": true,
			  "sshPublicKeyPath": "~/.ssh/id_rsa.pub",
			  "iam": {
			    "withAddonPolicies": {
			  	"imageBuilder": false,
			  	"autoScaler": false,
			  	"externalDNS": false,
			  	"appMesh": false,
			  	"ebs": false
			    }
			  },
			  "clusterDNS": "1.2.3.4"
		  },
		  {
			  "name": "test-ng1b",
			  "ami": "static",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 0,
			  "volumeType": "gp2",
			  "labels": {
			    "group": "b",
			    "seq": "1"
			  },
			  "allowSSH": true,
			  "sshPublicKeyPath": "~/.ssh/id_rsa.pub",
			  "iam": {
			    "withAddonPolicies": {
			  	"imageBuilder": false,
			  	"autoScaler": false,
			  	"externalDNS": false,
			  	"appMesh": false,
			  	"ebs": false
			    }
			  }
		  },
		  {
			  "name": "test-ng2b",
			  "ami": "static",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.xlarge",
			  "privateNetworking": false,
			  "securityGroups": {
			    "attachIDs": [
			  	"sg-1",
			  	"sg-2"
			    ],
			    "withShared": true,
			    "withLocal": false
			  },
			  "volumeSize": 0,
			  "volumeType": "gp2",
			  "labels": {
			    "group": "b",
			    "seq": "1"
			  },
			  "allowSSH": false,
			  "iam": {
			    "withAddonPolicies": {
			  	"imageBuilder": false,
			  	"autoScaler": false,
			  	"externalDNS": false,
			  	"appMesh": false,
			  	"ebs": false
			    }
			  },
			  "clusterDNS": "4.2.8.14"
		  },
		  {
			  "name": "test-ng3b",
			  "ami": "static",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "attachIDs": [
			  	"sg-1",
			  	"sg-2"
			    ],
			    "withShared": true,
			    "withLocal": false
			  },
			  "volumeSize": 192,
			  "volumeType": "gp2",
			  "labels": {
			    "group": "b",
			    "seq": "1"
			  },
			  "allowSSH": false,
			  "iam": {
			    "withAddonPolicies": {
			  	"imageBuilder": false,
			  	"autoScaler": false,
			  	"externalDNS": false,
			  	"appMesh": false,
			  	"ebs": false
			    }
			  }
		  }
		]
  }
`
