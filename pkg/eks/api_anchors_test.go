package eks_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

var _ = Describe("ParseConfig with YAML anchors and aliases", func() {
	BeforeEach(func() {
		err := api.Register()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should parse ClusterConfig with YAML anchors and aliases", func() {
		yamlWithAnchors := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

aliases:
  genericAttachPolicyARNs: &genericAttachPolicyARNs
  - arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly
  - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy

  genericNodeGroupSettings: &genericNodeGroupSettings
    minSize: 2
    maxSize: 5
    desiredCapacity: 2
    volumeSize: 30
    volumeType: gp3
    instanceTypes:
    - "t3a.small"
    - "t3.small"

metadata:
  name: eks-demo
  region: us-east-1

managedNodeGroups:
  - name: generic-1
    <<: *genericNodeGroupSettings
    iam:
      attachPolicyARNs: *genericAttachPolicyARNs

  - name: generic-2
    <<: *genericNodeGroupSettings
    iam:
      attachPolicyARNs: *genericAttachPolicyARNs`

		cfg, err := eks.ParseConfig([]byte(yamlWithAnchors))
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())
		Expect(cfg.Metadata.Name).To(Equal("eks-demo"))
		Expect(cfg.Metadata.Region).To(Equal("us-east-1"))
		Expect(cfg.ManagedNodeGroups).To(HaveLen(2))

		// Verify first node group
		ng1 := cfg.ManagedNodeGroups[0]
		Expect(ng1.Name).To(Equal("generic-1"))
		Expect(*ng1.MinSize).To(Equal(2))
		Expect(*ng1.MaxSize).To(Equal(5))
		Expect(*ng1.DesiredCapacity).To(Equal(2))
		Expect(*ng1.VolumeSize).To(Equal(30))
		Expect(*ng1.VolumeType).To(Equal("gp3"))
		Expect(ng1.InstanceTypes).To(Equal([]string{"t3a.small", "t3.small"}))
		Expect(ng1.IAM.AttachPolicyARNs).To(Equal([]string{
			"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
		}))

		// Verify second node group has same settings
		ng2 := cfg.ManagedNodeGroups[1]
		Expect(ng2.Name).To(Equal("generic-2"))
		Expect(*ng2.MinSize).To(Equal(2))
		Expect(*ng2.MaxSize).To(Equal(5))
		Expect(*ng2.DesiredCapacity).To(Equal(2))
		Expect(*ng2.VolumeSize).To(Equal(30))
		Expect(*ng2.VolumeType).To(Equal("gp3"))
		Expect(ng2.InstanceTypes).To(Equal([]string{"t3a.small", "t3.small"}))
		Expect(ng2.IAM.AttachPolicyARNs).To(Equal([]string{
			"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
		}))
	})

	It("should parse ClusterConfig with inline anchors", func() {
		yamlWithInlineAnchors := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: eks-demo
  region: us-east-1

managedNodeGroups:
  - name: generic-al2
    instanceTypes: &instance-types
    - "t3a.small"
    - "t3.small"
    minSize: 2
    iam:
      attachPolicyARNs: &policy-arns
      - arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly

  - name: generic-al2023
    instanceTypes: *instance-types
    minSize: 3
    iam:
      attachPolicyARNs: *policy-arns`

		cfg, err := eks.ParseConfig([]byte(yamlWithInlineAnchors))
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.ManagedNodeGroups).To(HaveLen(2))
		Expect(cfg.ManagedNodeGroups[0].InstanceTypes).To(Equal(cfg.ManagedNodeGroups[1].InstanceTypes))
		Expect(cfg.ManagedNodeGroups[0].IAM.AttachPolicyARNs).To(Equal(cfg.ManagedNodeGroups[1].IAM.AttachPolicyARNs))
	})

	It("should parse ClusterConfig without anchors normally", func() {
		yamlWithoutAnchors := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: simple-cluster
  region: us-west-2

managedNodeGroups:
  - name: workers
    instanceTypes:
    - t3.medium`

		cfg, err := eks.ParseConfig([]byte(yamlWithoutAnchors))
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.Metadata.Name).To(Equal("simple-cluster"))
		Expect(cfg.ManagedNodeGroups).To(HaveLen(1))
	})

	It("should reject invalid YAML with malformed anchors", func() {
		malformedYAML := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test
  region: us-east-1

managedNodeGroups:
  - name: workers
    instanceTypes: *nonexistent-anchor`

		_, err := eks.ParseConfig([]byte(malformedYAML))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown anchor"))
	})

	It("should reject YAML input that is too large", func() {
		// Create a YAML that exceeds the 1MB limit
		largeYAML := "apiVersion: eksctl.io/v1alpha5\nkind: ClusterConfig\nmetadata:\n  name: " + strings.Repeat("a", 1025*1024)
		_, err := eks.ParseConfig([]byte(largeYAML))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("YAML input too large"))
	})

	It("should reject YAML with excessive nesting depth", func() {
		// Create deeply nested YAML that exceeds the 10 level limit
		deepYAML := "apiVersion: eksctl.io/v1alpha5\nkind: ClusterConfig\nmetadata:\n  name: test\n"
		for i := 0; i < 12; i++ {
			deepYAML += strings.Repeat("  ", i) + "level" + fmt.Sprintf("%d:\n", i)
		}
		deepYAML += strings.Repeat("  ", 12) + "value: deep"

		_, err := eks.ParseConfig([]byte(deepYAML))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("YAML nesting too deep"))
	})

	It("should reject unknown top-level fields while allowing aliases", func() {
		yamlWithUnknownField := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test
  region: us-east-1

unknownField: should-be-rejected

managedNodeGroups:
  - name: workers
    instanceTypes:
    - t3.medium`

		_, err := eks.ParseConfig([]byte(yamlWithUnknownField))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown field"))
	})

	It("should handle empty YAML gracefully", func() {
		emptyYAML := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig`

		cfg, err := eks.ParseConfig([]byte(emptyYAML))
		Expect(err).NotTo(HaveOccurred()) // Basic parsing should succeed
		Expect(cfg).NotTo(BeNil())
		Expect(cfg.TypeMeta.APIVersion).To(Equal("eksctl.io/v1alpha5"))
		Expect(cfg.TypeMeta.Kind).To(Equal("ClusterConfig"))
	})

	It("should handle YAML with only aliases section", func() {
		aliasOnlyYAML := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

aliases:
  unused: &unused
    value: test

metadata:
  name: test
  region: us-east-1`

		cfg, err := eks.ParseConfig([]byte(aliasOnlyYAML))
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.Metadata.Name).To(Equal("test"))
	})

	It("should handle nested anchors", func() {
		nestedAnchorsYAML := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

aliases:
  base: &base
    minSize: 1
  extended: &extended
    <<: *base
    maxSize: 5

metadata:
  name: test
  region: us-east-1

managedNodeGroups:
  - name: workers
    <<: *extended`

		cfg, err := eks.ParseConfig([]byte(nestedAnchorsYAML))
		Expect(err).NotTo(HaveOccurred())
		Expect(*cfg.ManagedNodeGroups[0].MinSize).To(Equal(1))
		Expect(*cfg.ManagedNodeGroups[0].MaxSize).To(Equal(5))
	})

	It("should reject invalid YAML syntax", func() {
		invalidYAML := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: test
  region: us-east-1
  invalid: [unclosed array`

		_, err := eks.ParseConfig([]byte(invalidYAML))
		Expect(err).To(HaveOccurred())
	})

	It("should handle mixed anchor types in same config", func() {
		mixedAnchorsYAML := `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test
  region: us-east-1

managedNodeGroups:
  - name: ng1
    instanceTypes: &instances
    - t3.medium
    minSize: &min-size 2
    maxSize: 5
  - name: ng2
    instanceTypes: *instances
    minSize: *min-size
    maxSize: 10`

		cfg, err := eks.ParseConfig([]byte(mixedAnchorsYAML))
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.ManagedNodeGroups[0].InstanceTypes).To(Equal(cfg.ManagedNodeGroups[1].InstanceTypes))
		Expect(*cfg.ManagedNodeGroups[0].MinSize).To(Equal(*cfg.ManagedNodeGroups[1].MinSize))
	})
})
