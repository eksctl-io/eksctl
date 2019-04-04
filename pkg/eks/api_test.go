package eks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"

	. "github.com/weaveworks/eksctl/pkg/eks"
)

var _ = Describe("eksctl API", func() {

	Context("loading config files", func() {
		var (
			cfg *api.ClusterConfig
		)
		BeforeEach(func() {
			err := api.Register()
			Expect(err).ToNot(HaveOccurred())
			cfg = &api.ClusterConfig{}
		})

		It("should load a valid YAML config without error", func() {
			err := LoadConfigFromFile("../../examples/01-simple-cluster.yaml", cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.Metadata.Name).To(Equal("cluster-1"))
			Expect(cfg.NodeGroups).To(HaveLen(1))
		})

		It("should load a valid JSON config without error", func() {
			err := LoadConfigFromFile("testdata/example.json", cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.Metadata.Name).To(Equal("cluster-1"))
			Expect(cfg.NodeGroups).To(HaveLen(1))
		})

		It("should error when version is a float, not a string", func() {
			err := LoadConfigFromFile("testdata/bad-type-1.yaml", cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-type-1.yaml": v1alpha4.ClusterConfig.Metadata: v1alpha4.ClusterMeta.Version: ReadString: expects " or n, but found 1`))
		})

		It("should reject unknown field in a YAML config", func() {
			err := LoadConfigFromFile("testdata/bad-field-1.yaml", cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-field-1.yaml": error unmarshaling JSON: while decoding JSON: json: unknown field "zone"`))
		})

		It("should reject unknown field in a YAML config", func() {
			err := LoadConfigFromFile("testdata/bad-field-2.yaml", cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-field-2.yaml": error unmarshaling JSON: while decoding JSON: json: unknown field "bar"`))
		})

		It("should reject unknown field in a JSON config", func() {
			err := LoadConfigFromFile("testdata/bad-field-1.json", cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-field-1.json": error unmarshaling JSON: while decoding JSON: json: unknown field "nodes"`))
		})

		It("should reject old API version", func() {
			err := LoadConfigFromFile("testdata/old-version.json", cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/old-version.json": no kind "ClusterConfig" is registered for version "eksctl.io/v1alpha3" in scheme "k8s.io/client-go/kubernetes/scheme/register.go:60"`))
		})

		It("should error when cannot read a file", func() {
			err := LoadConfigFromFile("../../examples/nothing.xml", cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`reading config file "../../examples/nothing.xml": open ../../examples/nothing.xml: no such file or directory`))
		})
	})
})
