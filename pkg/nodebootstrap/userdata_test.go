package nodebootstrap

import (
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	kubeletapi "k8s.io/kubelet/config/v1beta1"
	"sigs.k8s.io/yaml"
)

var _ = Describe("User data", func() {
	Describe("generating max pods", func() {
		It("max pods mapping has the correct format", func() {
			maxPods := makeMaxPodsMapping()
			lines := strings.Split(strings.TrimSpace(maxPods), "\n")
			for _, line := range lines {
				parts := strings.Split(line, " ")
				Expect(parts[0]).To(MatchRegexp(`[a-zA-Z][a-zA-Z0-9]*\.[a-zA-Z0-9]+`))
				_, err := strconv.Atoi(parts[1])
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})

	Describe("creating kubelet config", func() {
		var (
			clusterConfig *api.ClusterConfig
			ng            *api.NodeGroup
		)
		BeforeEach(func() {
			clusterConfig = api.NewClusterConfig()
			ng = &api.NodeGroup{}
		})

		It("the kubelet is serialized with the correct format", func() {
			data, err := makeKubeletConfigYAML(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())

			kubelet := &kubeletapi.KubeletConfiguration{}

			errUnmarshal := yaml.UnmarshalStrict(data, kubelet)
			Expect(errUnmarshal).ToNot(HaveOccurred())
		})

		It("the kubelet config contains the overwritten values", func() {

			ng.KubeletExtraConfig = &api.InlineDocument{
				"kubeReserved": &map[string]string{
					"cpu":               "300m",
					"memory":            "300Mi",
					"ephemeral-storage": "1Gi",
				},
				"featureGates": map[string]bool{
					"HugePages":            false,
					"DynamicKubeletConfig": true,
				},
			}
			data, err := makeKubeletConfigYAML(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())

			kubelet := &kubeletapi.KubeletConfiguration{}

			errUnmarshal := yaml.UnmarshalStrict(data, kubelet)
			Expect(errUnmarshal).ToNot(HaveOccurred())

			Expect(kubelet.KubeReserved).ToNot(BeNil())
			Expect(kubelet.KubeReserved["cpu"]).To(Equal("300m"))
			Expect(kubelet.KubeReserved["memory"]).To(Equal("300Mi"))
			Expect(kubelet.KubeReserved["ephemeral-storage"]).To(Equal("1Gi"))
			Expect(kubelet.FeatureGates["HugePages"]).To(Equal(false))
			Expect(kubelet.FeatureGates["DynamicKubeletConfig"]).To(Equal(true))
			Expect(kubelet.FeatureGates["RotateKubeletServerCertificate"]).To(Equal(false))
		})
	})
})
