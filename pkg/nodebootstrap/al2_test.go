package nodebootstrap_test

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("AmazonLinux2 User Data", func() {
	var (
		clusterConfig *api.ClusterConfig
		ng            *api.NodeGroup
		bootstrapper  nodebootstrap.Bootstrapper
	)

	BeforeEach(func() {
		clusterConfig = api.NewClusterConfig()
		clusterConfig.Metadata.Name = "something-awesome"
		clusterConfig.Status = &api.ClusterStatus{}
		ng = &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				AMIFamily: "AmazonLinux2",
				SSH:       &api.NodeGroupSSH{},
			},
		}
	})

	When("SSM is enabled", func() {
		BeforeEach(func() {
			ng.SSH.EnableSSM = api.Enabled()
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("does not add the SSM install script to the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)

			var paths []string
			for _, f := range cloudCfg.WriteFiles {
				paths = append(paths, f.Path)
			}
			Expect(paths).NotTo(ContainElement("/var/lib/cloud/scripts/eksctl/install-ssm.al2.sh"))
		})
	})

	When("EFA is enabled", func() {
		BeforeEach(func() {
			enabled := true
			ng.EFAEnabled = &enabled
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("adds the ssm install script to the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/var/lib/cloud/scripts/eksctl/efa.al2.sh"))
			Expect(cloudCfg.WriteFiles[2].Permissions).To(Equal("0755"))
		})
	})

	type bootScriptEntry struct {
		clusterConfig    *api.ClusterConfig
		ng               *api.NodeGroup
		expectedUserData string
	}

	DescribeTable("Boot script environment variable in userdata", func(be bootScriptEntry) {
		if be.clusterConfig == nil {
			be.clusterConfig = api.NewClusterConfig()
			be.clusterConfig.Metadata.Name = "userdata-test"
			be.clusterConfig.Status = &api.ClusterStatus{}
		}
		be.ng.AMIFamily = "AmazonLinux2"
		bootstrapper := newBootstrapper(be.clusterConfig, be.ng)
		userData, err := bootstrapper.UserData()
		Expect(err).NotTo(HaveOccurred())
		cloudCfg := decode(userData)
		file := cloudCfg.WriteFiles[1]
		Expect(file.Path).To(Equal("/etc/eksctl/kubelet.env"))

		actualLines := strings.Split(file.Content, "\n")
		expectedLines := strings.Split(be.expectedUserData, "\n")
		Expect(actualLines).To(ConsistOf(expectedLines))
	},
		Entry("no fields set", bootScriptEntry{
			ng: api.NewNodeGroup(),
			expectedUserData: `CLUSTER_NAME=userdata-test
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=
NODE_TAINTS=
CONTAINER_RUNTIME=`,
		}),
		Entry("maxPods set", bootScriptEntry{
			ng: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					MaxPodsPerNode: 123,
					SSH:            &api.NodeGroupSSH{},
				},
			},
			expectedUserData: `CLUSTER_NAME=userdata-test
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=
NODE_TAINTS=
MAX_PODS=123
CONTAINER_RUNTIME=`,
		}),
		Entry("labels and taints set", bootScriptEntry{
			ng: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Labels: map[string]string{
						"role": "worker",
					},
					SSH: &api.NodeGroupSSH{},
				},
				Taints: []api.NodeGroupTaint{
					{
						Key:    "key1",
						Value:  "value1",
						Effect: "NoSchedule",
					},
				},
			},
			expectedUserData: `CLUSTER_NAME=userdata-test
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=role=worker
NODE_TAINTS=key1=value1:NoSchedule
CONTAINER_RUNTIME=`,
		}),
		Entry("container runtime set", bootScriptEntry{
			ng: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					SSH: &api.NodeGroupSSH{},
				},
				ContainerRuntime: aws.String(api.ContainerRuntimeContainerD),
			},
			expectedUserData: `CLUSTER_NAME=userdata-test
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=
NODE_TAINTS=
CONTAINER_RUNTIME=containerd`,
		}),

		Entry("non-default ServiceIPv4CIDR", bootScriptEntry{
			clusterConfig: func() *api.ClusterConfig {
				clusterConfig := api.NewClusterConfig()
				clusterConfig.Metadata.Name = "custom-dns"
				clusterConfig.Status = &api.ClusterStatus{
					KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
						ServiceIPv4CIDR: "172.16.0.0/12",
					},
				}
				return clusterConfig
			}(),
			ng: api.NewNodeGroup(),
			expectedUserData: `CLUSTER_NAME=custom-dns
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=
NODE_TAINTS=
CLUSTER_DNS=172.16.0.10
CONTAINER_RUNTIME=`,
		}),
	)

	Context("standard userdata", func() {
		var (
			err      error
			userData string
		)

		BeforeEach(func() {
			bootstrapper = newBootstrapper(clusterConfig, ng)
			userData, err = bootstrapper.UserData()
		})

		It("adds the boot script environment variable file to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[1].Path).To(Equal("/etc/eksctl/kubelet.env"))
			contentLines := strings.Split(cloudCfg.WriteFiles[1].Content, "\n")
			Expect(contentLines).To(ConsistOf(strings.Split(`CLUSTER_NAME=something-awesome
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=
NODE_TAINTS=
CONTAINER_RUNTIME=`, "\n")))
			Expect(cloudCfg.WriteFiles[1].Permissions).To(Equal("0644"))
		})

		It("adds the common linux boot script to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[3].Path).To(Equal("/var/lib/cloud/scripts/eksctl/bootstrap.helper.sh"))
			Expect(cloudCfg.WriteFiles[3].Permissions).To(Equal("0755"))
		})

		It("adds the al2 boot script to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/var/lib/cloud/scripts/eksctl/bootstrap.al2.sh"))
			Expect(cloudCfg.WriteFiles[2].Permissions).To(Equal("0755"))
		})
	})

	When("KubeletExtraConfig is provided by the user", func() {
		BeforeEach(func() {
			ng.KubeletExtraConfig = &api.InlineDocument{"foo": "bar"}
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("adds the settings to the kubelet extra args file in the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet-extra.json"))
			Expect(cloudCfg.WriteFiles[0].Content).To(Equal("{\"foo\":\"bar\"}"))
			Expect(cloudCfg.WriteFiles[0].Permissions).To(Equal("0644"))
		})
	})

	When("labels are set on the node config", func() {
		BeforeEach(func() {
			ng.Labels = map[string]string{"foo": "bar"}
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("adds the labels to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[1].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[1].Content).To(ContainSubstring("NODE_LABELS=foo=bar"))
			Expect(cloudCfg.WriteFiles[1].Permissions).To(Equal("0644"))
		})
	})

	When("taints are set on the node config", func() {
		BeforeEach(func() {
			ng.Taints = []api.NodeGroupTaint{
				{
					Key:    "foo",
					Effect: "NoSchedule",
				},
			}
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("adds the taints to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[1].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[1].Content).To(ContainSubstring("NODE_TAINTS=foo=:NoSchedule"))
			Expect(cloudCfg.WriteFiles[1].Permissions).To(Equal("0644"))
		})
	})

	When("clusterDNS is set on the node config", func() {
		BeforeEach(func() {
			ng.ClusterDNS = "1.2.3.4"
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("adds the taints to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[1].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[1].Content).To(ContainSubstring("CLUSTER_DNS=1.2.3.4"))
			Expect(cloudCfg.WriteFiles[1].Permissions).To(Equal("0644"))
		})
	})

	When("PreBootstrapCommands are set", func() {
		BeforeEach(func() {
			ng.PreBootstrapCommands = []string{"echo 'rubarb'"}
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("adds them to the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.Commands[0]).To(ContainElement("echo 'rubarb'"))
		})
	})

	When("OverrideBootstrapCommand is set", func() {
		var (
			err      error
			userData string
		)

		BeforeEach(func() {
			override := "echo 'crashoverride'"
			ng.OverrideBootstrapCommand = &override
			bootstrapper = newBootstrapper(clusterConfig, ng)

			userData, err = bootstrapper.UserData()
		})

		It("adds it to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.Commands[0]).To(ContainElement("echo 'crashoverride'"))
		})

		It("does not add the standard scripts to the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.Commands).To(HaveLen(2))
			Expect(cloudCfg.WriteFiles).To(HaveLen(3))
		})
	})
})
