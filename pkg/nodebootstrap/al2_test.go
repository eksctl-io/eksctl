package nodebootstrap_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("AmazonLinux2 User Data", func() {
	var (
		clusterName  string
		ng           *api.NodeGroup
		bootstrapper *nodebootstrap.AmazonLinux2
	)

	BeforeEach(func() {
		clusterName = "something-awesome"
		ng = &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				SSH: &api.NodeGroupSSH{},
			},
		}
	})

	When("SSM is enabled", func() {
		BeforeEach(func() {
			enabled := true
			ng.SSH.EnableSSM = &enabled
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
		})

		It("adds the ssm install script to the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[3].Path).To(Equal("/var/lib/cloud/scripts/eksctl/install-ssm.al2.sh"))
			Expect(cloudCfg.WriteFiles[3].Permissions).To(Equal("0755"))
		})
	})

	When("EFA is enabled", func() {
		BeforeEach(func() {
			enabled := true
			ng.EFAEnabled = &enabled
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
		})

		It("adds the ssm install script to the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[3].Path).To(Equal("/var/lib/cloud/scripts/eksctl/efa.al2.sh"))
			Expect(cloudCfg.WriteFiles[3].Permissions).To(Equal("0755"))
		})
	})

	Context("standard userdata", func() {
		var (
			err      error
			userData string
		)

		BeforeEach(func() {
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
			userData, err = bootstrapper.UserData()
		})

		It("adds the kubelet extra args and docker daemon extra args files to the userdata (sets cgroupDriver to systemd", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet-extra.json"))
			Expect(cloudCfg.WriteFiles[0].Content).To(Equal("{\"cgroupDriver\":\"systemd\"}"))
			Expect(cloudCfg.WriteFiles[0].Permissions).To(Equal("0644"))
			Expect(cloudCfg.WriteFiles[1].Path).To(Equal("/etc/eksctl/docker-extra.json"))
			Expect(cloudCfg.WriteFiles[1].Content).To(Equal("{\"exec-opts\":[\"native.cgroupdriver=systemd\"]}"))
			Expect(cloudCfg.WriteFiles[1].Permissions).To(Equal("0644"))
		})

		It("adds the boot script environment variable file to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[2].Content).To(Equal("NODE_LABELS=\nNODE_TAINTS=\nCLUSTER_NAME=something-awesome"))
			Expect(cloudCfg.WriteFiles[2].Permissions).To(Equal("0644"))
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
			Expect(cloudCfg.WriteFiles[4].Path).To(Equal("/var/lib/cloud/scripts/eksctl/bootstrap.al2.sh"))
			Expect(cloudCfg.WriteFiles[4].Permissions).To(Equal("0755"))
		})
	})

	When("KubeletExtraConfig is provided by the user", func() {
		BeforeEach(func() {
			ng.KubeletExtraConfig = &api.InlineDocument{"foo": "bar"}
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
		})

		It("adds the settings to the kubelet extra args file in the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet-extra.json"))
			Expect(cloudCfg.WriteFiles[0].Content).To(Equal("{\"cgroupDriver\":\"systemd\",\"foo\":\"bar\"}"))
			Expect(cloudCfg.WriteFiles[0].Permissions).To(Equal("0644"))
		})
	})

	When("labels are set on the node config", func() {
		BeforeEach(func() {
			ng.Labels = map[string]string{"foo": "bar"}
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
		})

		It("adds the labels to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[2].Content).To(ContainSubstring("NODE_LABELS=foo=bar"))
			Expect(cloudCfg.WriteFiles[2].Permissions).To(Equal("0644"))
		})
	})

	When("taints are set on the node config", func() {
		BeforeEach(func() {
			ng.Taints = map[string]string{"foo": "bar"}
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
		})

		It("adds the taints to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[2].Content).To(ContainSubstring("NODE_TAINTS=foo=:bar"))
			Expect(cloudCfg.WriteFiles[2].Permissions).To(Equal("0644"))
		})
	})

	When("clusterDNS is set on the node config", func() {
		BeforeEach(func() {
			ng.ClusterDNS = "1.2.3.4"
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
		})

		It("adds the taints to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[2].Content).To(ContainSubstring("CLUSTER_DNS=1.2.3.4"))
			Expect(cloudCfg.WriteFiles[2].Permissions).To(Equal("0644"))
		})
	})

	When("PreBootstrapCommands are set", func() {
		BeforeEach(func() {
			ng.PreBootstrapCommands = []string{"echo 'rubarb'"}
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)
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
			bootstrapper = nodebootstrap.NewAL2Bootstrapper(clusterName, ng)

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
			Expect(cloudCfg.Commands).To(HaveLen(1))
			Expect(cloudCfg.WriteFiles).To(HaveLen(0))
		})
	})
})
