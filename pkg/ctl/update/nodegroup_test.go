package update

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
)

var _ = Describe("update nodegroup", func() {
	It("returns error if cluster is not set", func() {
		cmd := newMockCmd("nodegroup")
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("cluster name must be set")))
	})

	It("returns error if nodegroup name is not set as flag", func() {
		cmd := newMockCmd("nodegroup", "--cluster", "cluster-name")
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("nodegroup name must be set")))
	})

	It("returns error if nodegroup is not set in config", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   "cluster-1",
				Region: "us-west-2",
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("managedNodeGroups must be set")))
	})

	It("returns error if cluster is set in cfg and as flag", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name: "cluster-1",
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--cluster", "cluster-name", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("cannot use --cluster when --config-file/-f is set")))
	})

	It("returns error if ng name is set in cfg and as flag", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name: "cluster-1",
			},
			ManagedNodeGroups: []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "ngName",
					},
				},
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--name", "ng-name", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("cannot use --name when --config-file/-f is set")))
	})

	It("returns error if region is set in cfg and as flag", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   "cluster-1",
				Region: "us-west-2",
			},
			ManagedNodeGroups: []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:           "ngName",
						MaxPodsPerNode: 5,
					},
				},
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--region", "region-name", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("cannot use --region when --config-file/-f is set")))
	})

	It("returns error if cfg contains unsupported field in NodeGroupBase fields", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   "cluster-1",
				Region: "us-west-2",
			},
			ManagedNodeGroups: []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:           "ngName",
						MaxPodsPerNode: 5,
					},
				},
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("unsupported field: MaxPodsPerNode cannot be used with `eksctl update nodegroup`")))
	})

	It("returns error if cfg contains unsupported field in NodeGroup fields", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   "cluster-1",
				Region: "us-west-2",
			},
			ManagedNodeGroups: []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "ngName",
					},
					Spot: true,
				},
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("unsupported field: Spot cannot be used with `eksctl update nodegroup`")))
	})

	It("returns error if cfg contains multiple nodegroups", func() {
		cfg := &api.ClusterConfig{
			TypeMeta: api.ClusterConfigTypeMeta(),
			Metadata: &api.ClusterMeta{
				Name:   "cluster-1",
				Region: "us-west-2",
			},
			ManagedNodeGroups: []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "ngName",
					},
				},
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "ngName",
					},
				},
			},
		}
		config := ctltest.CreateConfigFile(cfg)
		cmd := newMockCmd("nodegroup", "--config-file", config)
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("please update one NodeGroup at a time")))
	})
})
