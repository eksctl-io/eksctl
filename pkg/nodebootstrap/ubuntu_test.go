package nodebootstrap_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Ubuntu User Data", func() {
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
				AMIFamily: "Ubuntu2004",
			},
		}
	})

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
NODE_TAINTS=`, "\n")))
			Expect(cloudCfg.WriteFiles[1].Permissions).To(Equal("0644"))
		})

		It("adds the common linux boot script to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[3].Path).To(Equal("/var/lib/cloud/scripts/eksctl/bootstrap.helper.sh"))
			Expect(cloudCfg.WriteFiles[3].Permissions).To(Equal("0755"))
		})

		It("adds the ubuntu boot script to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/var/lib/cloud/scripts/eksctl/bootstrap.ubuntu.sh"))
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
					Effect: "NoExecute",
				},
				{
					Key:    "one",
					Value:  "two",
					Effect: "NoSchedule",
				},
			}
			bootstrapper = newBootstrapper(clusterConfig, ng)
		})

		It("adds the taints to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			file := cloudCfg.WriteFiles[1]
			Expect(file.Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(file.Content).To(ContainSubstring("NODE_TAINTS"))
			Expect(file.Content).To(ContainSubstring("foo=:NoExecute"))
			Expect(file.Content).To(ContainSubstring("one=two:NoSchedule"))
			Expect(file.Permissions).To(Equal("0644"))
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
