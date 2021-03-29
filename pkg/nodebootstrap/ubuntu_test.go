package nodebootstrap_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Ubuntu User Data", func() {
	var (
		clusterName  string
		ng           *api.NodeGroup
		bootstrapper *nodebootstrap.Ubuntu
	)

	BeforeEach(func() {
		clusterName = "something-awesome"
		ng = &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{},
		}
	})

	Context("standard userdata", func() {
		var (
			err      error
			userData string
		)

		BeforeEach(func() {
			bootstrapper = nodebootstrap.NewUbuntuBootstrapper(clusterName, ng)
			userData, err = bootstrapper.UserData()
		})

		It("adds the boot script environment variable file to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[0].Content).To(Equal("NODE_LABELS=\nNODE_TAINTS=\nCLUSTER_NAME=something-awesome"))
		})

		It("adds the common linux boot script to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[1].Path).To(Equal("/var/lib/cloud/scripts/eksctl/bootstrap.linux.sh"))
		})

		It("adds the ubuntu boot script to the userdata", func() {
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[2].Path).To(Equal("/var/lib/cloud/scripts/eksctl/bootstrap.ubuntu.sh"))
		})
	})

	When("KubeletExtraConfig is set", func() {
		BeforeEach(func() {
			ng.KubeletExtraConfig = &api.InlineDocument{"foo": "bar"}
			bootstrapper = nodebootstrap.NewUbuntuBootstrapper(clusterName, ng)
		})

		It("adds the kubelet extra args file to the userdata", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet-extra.json"))
			Expect(cloudCfg.WriteFiles[0].Content).To(Equal("{\"foo\":\"bar\"}"))
		})
	})

	When("labels are set on the node config", func() {
		BeforeEach(func() {
			ng.Labels = map[string]string{"foo": "bar"}
			bootstrapper = nodebootstrap.NewUbuntuBootstrapper(clusterName, ng)
		})

		It("adds the labels to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[0].Content).To(ContainSubstring("NODE_LABELS=foo=bar"))
		})
	})

	When("taints are set on the node config", func() {
		BeforeEach(func() {
			ng.Taints = map[string]string{"foo": "bar"}
			bootstrapper = nodebootstrap.NewUbuntuBootstrapper(clusterName, ng)
		})

		It("adds the taints to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[0].Content).To(ContainSubstring("NODE_TAINTS=foo=:bar"))
		})
	})

	When("clusterDNS is set on the node config", func() {
		BeforeEach(func() {
			ng.ClusterDNS = "1.2.3.4"
			bootstrapper = nodebootstrap.NewUbuntuBootstrapper(clusterName, ng)
		})

		It("adds the taints to the env file", func() {
			userData, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())

			cloudCfg := decode(userData)
			Expect(cloudCfg.WriteFiles[0].Path).To(Equal("/etc/eksctl/kubelet.env"))
			Expect(cloudCfg.WriteFiles[0].Content).To(ContainSubstring("CLUSTER_DNS=1.2.3.4"))
		})
	})

	When("PreBootstrapCommands are set", func() {
		BeforeEach(func() {
			ng.PreBootstrapCommands = []string{"echo 'rubarb'"}
			bootstrapper = nodebootstrap.NewUbuntuBootstrapper(clusterName, ng)
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
			bootstrapper = nodebootstrap.NewUbuntuBootstrapper(clusterName, ng)

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
