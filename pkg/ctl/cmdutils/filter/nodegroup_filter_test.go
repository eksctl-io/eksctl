package filter

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/printers"
)

var _ = Describe("nodegroup filter", func() {

	getNodeGroupNames := func(clusterConfig *api.ClusterConfig) []string {
		var ngNames []string
		for _, ng := range clusterConfig.NodeGroups {
			ngNames = append(ngNames, ng.NameString())
		}
		return ngNames
	}

	Context("Match", func() {
		var (
			filter *NodeGroupFilter
			cfg    *api.ClusterConfig
		)

		BeforeEach(func() {
			cfg = newClusterConfig()
			addGroupA(cfg)
			addGroupB(cfg)

			filter = NewNodeGroupFilter()
		})

		It("regression: should only match the ones included in the filter when non existing ngs are present in the config file", func() {
			nonExistentNg := cfg.NewNodeGroup()
			nonExistentNg.Name = "non-existing-in-cluster"
			cfg.NodeGroups = append(cfg.NodeGroups, nonExistentNg)

			err := filter.AppendIncludeGlobs(getNodeGroupNames(cfg), "test-ng1?")
			Expect(err).ToNot(HaveOccurred())

			Expect(filter.Match("test-ng3x")).To(BeFalse())
			Expect(filter.Match("test-ng3b")).To(BeFalse())
			Expect(filter.Match("xyz1")).To(BeFalse())
			Expect(filter.Match("yz1")).To(BeFalse())
			Expect(filter.Match("test-ng1")).To(BeFalse())
			Expect(filter.Match("test-ng1a")).To(BeTrue())
			Expect(filter.Match("test-ng1n")).To(BeTrue())

			included, excluded := filter.matchAll(filter.collectNames(cfg.NodeGroups))
			Expect(included).To(HaveLen(2))
			Expect(included.HasAll("test-ng1a", "test-ng1b")).To(BeTrue())
			Expect(excluded).To(HaveLen(5))
			Expect(excluded.HasAll("test-ng2a", "test-ng2b", "test-ng3a", "test-ng3b", "non-existing-in-cluster")).To(BeTrue())
		})

		It("only-missing (only-remote) works correctly", func() {
			mockLister := newMockStackLister(
				"test-ng1a",
				"test-ng2a",
				"test-ng3a",
				"test-ng1b",
				"test-ng2b",
				"test-ng3b",
				"non-existing-in-cfg-1",
				"non-existing-in-cfg-2",
			)
			err := filter.SetOnlyRemote(mockLister, cfg)
			Expect(err).ToNot(HaveOccurred())

			included, excluded := filter.matchAll(filter.collectNames(cfg.NodeGroups))
			Expect(included).To(HaveLen(2))
			Expect(included.HasAll("non-existing-in-cfg-1", "non-existing-in-cfg-2")).To(BeTrue())
			Expect(excluded).To(HaveLen(6))
			Expect(excluded.HasAll(
				"test-ng1a",
				"test-ng2a",
				"test-ng3a",
				"test-ng1b",
				"test-ng2b",
				"test-ng3b",
			)).To(BeTrue())
		})

		It("should match only local nodegroups", func() {
			err := filter.AppendIncludeGlobs(getNodeGroupNames(cfg), "test-ng1?")
			Expect(err).ToNot(HaveOccurred())

			mockLister := newMockStackLister(
				"test-ng1a",
				"test-ng2a",
				"test-ng3a",
			)
			err = filter.SetOnlyLocal(mockLister, cfg)
			Expect(err).ToNot(HaveOccurred())

			included, excluded := filter.matchAll(filter.collectNames(cfg.NodeGroups))
			Expect(included).To(HaveLen(1))
			Expect(included.HasAll("test-ng1b")).To(BeTrue())
			Expect(excluded).To(HaveLen(5))
			Expect(excluded.HasAll("test-ng1a", "test-ng2a", "test-ng3a", "test-ng2b", "test-ng3b")).To(BeTrue())
		})

		It("should match only local nodegroups with exclude and include rules", func() {
			err := filter.AppendIncludeGlobs(getNodeGroupNames(cfg), "test-ng?a", "test-ng?b")
			Expect(err).ToNot(HaveOccurred())

			mockLister := newMockStackLister(
				"test-ng2a",
				"test-ng1b",
				"test-ng2b",
			)
			err = filter.SetOnlyLocal(mockLister, cfg)
			Expect(err).ToNot(HaveOccurred())

			err = filter.AppendExcludeGlobs("test-ng1a", "test-ng2?")
			Expect(err).ToNot(HaveOccurred())

			included, excluded := filter.matchAll(filter.collectNames(cfg.NodeGroups))
			Expect(included).To(HaveLen(2))
			Expect(included.HasAll("test-ng3a", "test-ng3b")).To(BeTrue())
			Expect(excluded).To(HaveLen(4))
			Expect(excluded.HasAll("test-ng1a", "test-ng1b", "test-ng2a", "test-ng2b")).To(BeTrue())
		})
	})

	Context("ForEach", func() {

		It("should iterate over unique nodegroups, apply defaults and validate", func() {
			cfg := newClusterConfig()
			addGroupA(cfg)
			addGroupB(cfg)

			filter := NewNodeGroupFilter()
			printer := printers.NewJSONPrinter()
			names := []string{}

			err := filter.ForEach(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				api.SetNodeGroupDefaults(nodeGroup, cfg.Metadata)
				err := api.ValidateNodeGroup(i, nodeGroup)
				Expect(err).ToNot(HaveOccurred())
				return nil
			})
			Expect(err).ToNot(HaveOccurred())

			err = filter.ForEach(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				Expect(nodeGroup).To(Equal(cfg.NodeGroups[i]))
				names = append(names, nodeGroup.Name)
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(names).To(Equal([]string{"test-ng1a", "test-ng2a", "test-ng3a", "test-ng1b", "test-ng2b", "test-ng3b"}))

			w := &bytes.Buffer{}

			_ = printer.PrintObj(cfg, w)

			Expect(w.Bytes()).To(MatchJSON(expected))
		})

		It("should be able to skip all nodegroups", func() {
			cfg := newClusterConfig()
			addGroupA(cfg)
			addGroupB(cfg)

			filter := NewNodeGroupFilter()
			filter.delegate.ExcludeAll = true

			err := filter.ForEach(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				api.SetNodeGroupDefaults(nodeGroup, cfg.Metadata)
				err := api.ValidateNodeGroup(i, nodeGroup)
				Expect(err).ToNot(HaveOccurred())
				return nil
			})
			Expect(err).ToNot(HaveOccurred())

			callback := false
			err = filter.ForEach(cfg.NodeGroups, func(_ int, _ *api.NodeGroup) error {
				callback = true
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(callback).To(BeFalse())
		})

		It("should iterate over unique nodegroups", func() {
			cfg := newClusterConfig()
			addGroupA(cfg)

			filter := NewNodeGroupFilter()
			names := []string{}

			err := filter.ForEach(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				Expect(nodeGroup).To(Equal(cfg.NodeGroups[i]))
				names = append(names, nodeGroup.Name)
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(names).To(Equal([]string{"test-ng1a", "test-ng2a", "test-ng3a"}))

			names = []string{}
			cfg.NodeGroups[0].Name = "ng-x0"
			cfg.NodeGroups[1].Name = "ng-x1"
			cfg.NodeGroups[2].Name = "ng-x2"

			err = filter.ForEach(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				Expect(nodeGroup).To(Equal(cfg.NodeGroups[i]))
				names = append(names, nodeGroup.Name)
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(names).To(Equal([]string{"ng-x0", "ng-x1", "ng-x2"}))
		})

		It("should iterate over unique nodegroups and filter some out", func() {
			cfg := newClusterConfig()
			addGroupA(cfg)
			addGroupB(cfg)

			filter := NewNodeGroupFilter()
			names := []string{}

			err := filter.ForEach(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				Expect(nodeGroup).To(Equal(cfg.NodeGroups[i]))
				names = append(names, nodeGroup.Name)
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(names).To(Equal([]string{"test-ng1a", "test-ng2a", "test-ng3a", "test-ng1b", "test-ng2b", "test-ng3b"}))

			names = []string{}

			err = filter.AppendIncludeGlobs(getNodeGroupNames(cfg), "t?xyz?", "ab*z123?")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`no nodegroups match include glob filter specification: "t?xyz?,ab*z123?"`))

			err = filter.AppendIncludeGlobs(getNodeGroupNames(cfg), "test-ng1?", "te*-ng3?")
			Expect(err).ToNot(HaveOccurred())
			err = filter.ForEach(cfg.NodeGroups, func(i int, nodeGroup *api.NodeGroup) error {
				Expect(nodeGroup).To(Equal(cfg.NodeGroups[i]))
				names = append(names, nodeGroup.Name)
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(names).To(Equal([]string{"test-ng1a", "test-ng3a", "test-ng1b", "test-ng3b"}))
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

	var (
		ng1aVolSize = 768
		ng1aVolType = api.NodeVolumeTypeIO1
		ng1aVolIOPS = 200
	)

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng1a"
	ng.VolumeSize = &ng1aVolSize
	ng.VolumeType = &ng1aVolType
	ng.VolumeIOPS = &ng1aVolIOPS
	ng.IAM.AttachPolicyARNs = []string{"arn:aws:iam::aws:policy/Foo"}
	ng.Labels = map[string]string{"group": "a", "seq": "1"}
	ng.SSH = nil

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng2a"
	ng.IAM.AttachPolicyARNs = []string{"arn:aws:iam::aws:policy/Bar"}
	ng.Labels = map[string]string{"group": "a", "seq": "2"}
	ng.SSH = nil

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng3a"
	ng.ClusterDNS = "1.2.3.4"
	ng.InstanceType = "m3.large"
	ng.SSH.Allow = api.Enabled()
	ng.SSH.PublicKeyPath = nil
	ng.Labels = map[string]string{"group": "a", "seq": "3"}
}

func addGroupB(cfg *api.ClusterConfig) {
	var ng *api.NodeGroup

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng1b"
	ng.SSH.Allow = api.Enabled()
	ng.SSH.PublicKeyPath = nil
	ng.Labels = map[string]string{"group": "b", "seq": "1"}

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng2b"
	ng.ClusterDNS = "4.2.8.14"
	ng.InstanceType = "m5.xlarge"
	ng.SecurityGroups.AttachIDs = []string{"sg-1", "sg-2"}
	ng.SecurityGroups.WithLocal = api.Disabled()
	ng.Labels = map[string]string{"group": "b", "seq": "1"}
	ng.SSH = nil

	ng3bVolSize := 192

	ng = cfg.NewNodeGroup()
	ng.Name = "test-ng3b"
	ng.VolumeSize = &ng3bVolSize
	ng.SecurityGroups.AttachIDs = []string{"sg-1", "sg-2"}
	ng.SecurityGroups.WithLocal = api.Disabled()
	ng.Labels = map[string]string{"group": "b", "seq": "1"}
	ng.SSH = nil
}

const expected = `
  {
		"kind": "ClusterConfig",
		"apiVersion": "eksctl.io/v1alpha5",
		"metadata": {
		  "name": "test-3x3-ngs",
		  "region": "eu-central-1",
		  "version": "1.18"
		},
		"iam": {
		  "withOIDC": false
		},
		"vpc": {
			"cidr": "192.168.0.0/16",
			"autoAllocateIPv6": false,
			"nat": {
				"gateway": "Single"
			  },
             "clusterEndpoints": {}
		},
		"cloudWatch": {
		  "clusterLogging": {}
		},
		"privateCluster": {
			"enabled": false
		},
		"nodeGroups": [
		  {
			  "name": "test-ng1a",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 768,
			  "volumeType": "io1",
              "OSVolumeName": null,
			  "volumeIOPS": 200,
			  "labels": {
				"alpha.eksctl.io/cluster-name": "test-3x3-ngs",
				"alpha.eksctl.io/nodegroup-name": "test-ng1a",
			    "group": "a",
			    "seq": "1"
			  },
			  "ssh": {
				"allow": false
			  },
			  "iam": {
			    "attachPolicyARNs": [
				  "arn:aws:iam::aws:policy/Foo"
			    ],
			    "withAddonPolicies": {
				  "imageBuilder": false,
				  "autoScaler": false,
				  "externalDNS": false,
				  "certManager": false,
				  "appMesh": false,
				  "appMeshPreview": false,
				  "ebs": false,
				  "fsx": false,
				  "efs": false,
				  "albIngress": false,
				  "xRay": false,
				  "cloudWatch": false
			    }
			  },
			  "disableIMDSv1": false,
			  "disablePodIMDS": false
		  },
		  {
			  "name": "test-ng2a",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 80,
			  "volumeType": "gp2",
              "OSVolumeName": null,
			  "labels": {
				"alpha.eksctl.io/cluster-name": "test-3x3-ngs",
				"alpha.eksctl.io/nodegroup-name": "test-ng2a",
			    "group": "a",
			    "seq": "2"
			  },
			  "ssh": {
			    "allow": false
			  },
			  "iam": {
			    "attachPolicyARNs": [
				  "arn:aws:iam::aws:policy/Bar"
			    ],
			    "withAddonPolicies": {
				  "imageBuilder": false,
				  "autoScaler": false,
				  "externalDNS": false,
				  "certManager": false,
				  "appMesh": false,
				  "appMeshPreview": false,
				  "ebs": false,
				  "fsx": false,
				  "efs": false,
				  "albIngress": false,
				  "xRay": false,
				  "cloudWatch": false
			    }
			  },
			  "disableIMDSv1": false,
			  "disablePodIMDS": false
		  },
		  {
			  "name": "test-ng3a",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m3.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 80,
			  "volumeType": "gp2",
              "OSVolumeName": null,
			  "labels": {
				"alpha.eksctl.io/cluster-name": "test-3x3-ngs",
				"alpha.eksctl.io/nodegroup-name": "test-ng3a",
			    "group": "a",
			    "seq": "3"
			  },
			  "ssh": {
			    "allow": true,
			    "publicKeyPath": "~/.ssh/id_rsa.pub"
			  },
			  "iam": {
				"withAddonPolicies": {
				  "imageBuilder": false,
				  "autoScaler": false,
				  "externalDNS": false,
				  "certManager": false,
				  "appMesh": false,
				  "appMeshPreview": false,
				  "ebs": false,
				  "fsx": false,
				  "efs": false,
				  "albIngress": false,
				  "xRay": false,
				  "cloudWatch": false
			    }
			  },
			  "clusterDNS": "1.2.3.4",
			  "disableIMDSv1": false,
			  "disablePodIMDS": false
		  },
		  {
			  "name": "test-ng1b",
			  "amiFamily": "AmazonLinux2",
			  "instanceType": "m5.large",
			  "privateNetworking": false,
			  "securityGroups": {
			    "withShared": true,
			    "withLocal": true
			  },
			  "volumeSize": 80,
			  "volumeType": "gp2",
              "OSVolumeName": null,
			  "labels": {
				"alpha.eksctl.io/cluster-name": "test-3x3-ngs",
				"alpha.eksctl.io/nodegroup-name": "test-ng1b",
			    "group": "b",
			    "seq": "1"
			  },
			  "ssh": {
			    "allow": true,
			    "publicKeyPath": "~/.ssh/id_rsa.pub"
              },
			  "iam": {
			    "withAddonPolicies": {
			  	  "imageBuilder": false,
			  	  "autoScaler": false,
			  	  "externalDNS": false,
			  	  "certManager": false,
			  	  "appMesh": false,
				  "appMeshPreview": false,
			  	  "ebs": false,
			  	  "fsx": false,
			  	  "efs": false,
				  "albIngress": false,
				  "xRay": false,
				  "cloudWatch": false
			    }
			  },
			  "disableIMDSv1": false,
			  "disablePodIMDS": false
		  },
		  {
			  "name": "test-ng2b",
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
			  "volumeSize": 80,
			  "volumeType": "gp2",
              "OSVolumeName": null,
			  "labels": {
				"alpha.eksctl.io/cluster-name": "test-3x3-ngs",
				"alpha.eksctl.io/nodegroup-name": "test-ng2b",
			    "group": "b",
			    "seq": "1"
			  },
			  "ssh": {
			    "allow": false
              },
			  "iam": {
			    "withAddonPolicies": {
			  	  "imageBuilder": false,
			  	  "autoScaler": false,
			  	  "externalDNS": false,
			  	  "certManager": false,
			  	  "appMesh": false,
				  "appMeshPreview": false,
			  	  "ebs": false,
			  	  "fsx": false,
			  	  "efs": false,
				  "albIngress": false,
				  "xRay": false,
				  "cloudWatch": false
			    }
			  },
			  "clusterDNS": "4.2.8.14",
			  "disableIMDSv1": false,
			  "disablePodIMDS": false
		  },
		  {
			  "name": "test-ng3b",
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
              "OSVolumeName": null,
			  "labels": {
				"alpha.eksctl.io/cluster-name": "test-3x3-ngs",
				"alpha.eksctl.io/nodegroup-name": "test-ng3b",
			    "group": "b",
			    "seq": "1"
			  },
			  "ssh": {
			    "allow": false
			  },
			  "iam": {
			    "withAddonPolicies": {
			  	  "imageBuilder": false,
			  	  "autoScaler": false,
			  	  "externalDNS": false,
			  	  "certManager": false,
			  	  "appMesh": false,
				  "appMeshPreview": false,
			  	  "ebs": false,
			  	  "fsx": false,
			  	  "efs": false,
				  "albIngress": false,
				  "xRay": false,
				  "cloudWatch": false
			    }
			  },
			  "disableIMDSv1": false,
			  "disablePodIMDS": false
		  }
		]
  }
`

type mockStackLister struct {
	nodesResult []manager.NodeGroupStack
}

func (s *mockStackLister) ListNodeGroupStacks() ([]manager.NodeGroupStack, error) {
	return s.nodesResult, nil
}

func newMockStackLister(ngs ...string) *mockStackLister {
	stacks := make([]manager.NodeGroupStack, 0)
	for _, ng := range ngs {
		stacks = append(stacks, manager.NodeGroupStack{
			NodeGroupName: ng,
		})
	}
	return &mockStackLister{
		nodesResult: stacks,
	}
}
