package create

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

var _ = Describe("create nodegroup", func() {
	Describe("un-managed node group", func() {
		DescribeTable("create nodegroup successfully",
			func(args ...string) {
				commandArgs := append([]string{"nodegroup", "--cluster", "clusterName"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createNodeGroupCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options nodegroupOptions) error {
						Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
						Expect(ng.Name).NotTo(BeNil())
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).To(Not(HaveOccurred()))
				Expect(count).To(Equal(1))
			},
			Entry("with nodegroup name as flag", "--name", "nodegroupName"),
			Entry("with nodegroup name with a hyphen as flag", "--name", "nodegroup-name"),
			Entry("with nodegroup name as argument", "nodegroupName"),
			Entry("with nodegroup name with a hyphen as argument", "nodegroup-name"),
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with node-volume-type flag", "--node-volume-type", "gp2"),
			Entry("with max-pods-per-node flag", "--max-pods-per-node", "20"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with enable-ssm flag", "--enable-ssm"),
			Entry("with node-ami flag", "--node-ami", "ami-dummy-123"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-security-groups flag", "--node-security-groups", "sg-123"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3"),
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
			Entry("with subnet-ids flag", "--subnet-ids", "id1,id2,id3"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"nodegroup", "--managed=false"}, c.args...)
				cmd := newDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(c.error)))
			},
			Entry("without cluster name", invalidParamsCase{
				args:  []string{"--name", "nodegroupName"},
				error: "--cluster must be set",
			}),
			Entry("with nodegroup name as argument and flag", invalidParamsCase{
				args:  []string{"--cluster", "clusterName", "--name", "nodegroupName", "nodegroupName"},
				error: "--name=nodegroupName and argument nodegroupName cannot be used at the same time",
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"--invalid", "dummy"},
				error: "unknown flag: --invalid",
			}),
			Entry("with spot flag", invalidParamsCase{
				args:  []string{"--cluster", "foo", "--spot"},
				error: "--spot is only valid with managed nodegroups (--managed)",
			}),
			Entry("with instance-types flag", invalidParamsCase{
				args:  []string{"--cluster", "foo", "--instance-types", "some-type"},
				error: "--instance-types is only valid with managed nodegroups (--managed)",
			}),
			Entry("with nodegroup name as flag with invalid characters", invalidParamsCase{
				args:  []string{"--cluster", "clusterName", "--name", "eksctl-ng_k8s_nodegroup1"},
				error: "validation for eksctl-ng_k8s_nodegroup1 failed, name must satisfy regular expression pattern: [a-zA-Z][-a-zA-Z0-9]*",
			}),
		)
	})

	Describe("managed node group", func() {
		DescribeTable("create nodegroup successfully",
			func(args ...string) {
				commandArgs := append([]string{"nodegroup", "--managed", "--cluster", "clusterName"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createNodeGroupCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options nodegroupOptions) error {
						Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
						Expect(ng.Name).NotTo(BeNil())
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).To(Not(HaveOccurred()))
				Expect(count).To(Equal(1))
			},
			Entry("without nodegroup name", ""),
			Entry("with nodegroup name as flag", "--name", "nodegroupName"),
			Entry("with nodegroup name with a hyphen as flag", "--name", "nodegroup-name"),
			Entry("with nodegroup name as argument", "nodegroupName"),
			Entry("with nodegroup name with a hyphen as argument", "nodegroup-name"),
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with enable-ssm flag", "--enable-ssm"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3"),
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
			Entry("with Ubuntu AMI", "--node-ami-family", "Ubuntu2004"),
			Entry("with Bottlerocket AMI", "--node-ami-family", "Bottlerocket"),
			Entry("with subnet-ids flag", "--subnet-ids", "id1,id2,id3"),
			Entry("with Windows AMI", "--node-ami-family", "WindowsServer2019FullContainer"),
			Entry("with Windows AMI", "--node-ami-family", "WindowsServer2019CoreContainer"),
			Entry("with Windows AMI", "--node-ami-family", "WindowsServer2022FullContainer"),
			Entry("with Windows AMI", "--node-ami-family", "WindowsServer2022CoreContainer"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"nodegroup", "--managed", "--cluster", "clusterName"}, c.args...)
				cmd := newDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(c.error)))
			},
			Entry("with nodegroup name as argument and flag", invalidParamsCase{
				args:  []string{"--name", "nodegroupName", "nodegroupName"},
				error: "--name=nodegroupName and argument nodegroupName cannot be used at the same time",
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"--invalid", "dummy"},
				error: "unknown flag: --invalid",
			}),
			Entry("with nodegroup name as flag with invalid characters", invalidParamsCase{
				args:  []string{"--name", "eksctl-ng_k8s_nodegroup1"},
				error: "validation for eksctl-ng_k8s_nodegroup1 failed, name must satisfy regular expression pattern: [a-zA-Z][-a-zA-Z0-9]*",
			}),
			Entry("with version flag", invalidParamsCase{
				args:  []string{"--version", "1.18"},
				error: "--version is only valid with unmanaged nodegroups",
			}),
		)
	})

	type checkNodeGroupVersionInput struct {
		ctl             *eks.ClusterProvider
		meta            *api.ClusterMeta
		expectedVersion string
		expectedErr     string
	}

	providerVersion1_23 := &eks.ClusterProvider{
		Status: &eks.ProviderStatus{
			ClusterInfo: &eks.ClusterInfo{
				Cluster: &types.Cluster{
					Version: aws.String(api.Version1_23),
				},
			},
		},
	}

	DescribeTable("checkNodeGroupVersion",
		func(input checkNodeGroupVersionInput) {
			err := checkNodeGroupVersion(input.ctl, input.meta)
			if input.expectedErr != "" {
				Expect(err).To(MatchError(ContainSubstring(input.expectedErr)))
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(input.meta.Version).To(Equal(input.expectedVersion))
		},
		Entry("version is left empty", checkNodeGroupVersionInput{
			ctl:             providerVersion1_23,
			meta:            &api.ClusterMeta{},
			expectedVersion: api.Version1_23,
		}),
		Entry("version is set to auto", checkNodeGroupVersionInput{
			ctl: providerVersion1_23,
			meta: &api.ClusterMeta{
				Version: "auto",
			},
			expectedVersion: api.Version1_23,
		}),
		Entry("version is set to latest", checkNodeGroupVersionInput{
			ctl: providerVersion1_23,
			meta: &api.ClusterMeta{
				Version: "latest",
			},
			expectedVersion: api.LatestVersion,
		}),
		Entry("version is set to deprecated version", checkNodeGroupVersionInput{
			meta: &api.ClusterMeta{
				Version: api.Version1_15,
			},
			expectedErr: fmt.Sprintf("invalid version, %s is no longer supported", api.Version1_15),
		}),
		Entry("version is set to unsupported version", checkNodeGroupVersionInput{
			meta: &api.ClusterMeta{
				Version: "100",
			},
			expectedErr: fmt.Sprintf("invalid version 100, supported values: auto, default, latest, %s", strings.Join(api.SupportedVersions(), ", ")),
		}),
		Entry("fails to retrieve control plane version", checkNodeGroupVersionInput{
			ctl: &eks.ClusterProvider{
				Status: &eks.ProviderStatus{},
			},
			meta: &api.ClusterMeta{
				Version: "auto",
			},
			expectedErr: "unable to get control plane version",
		}),
	)
})
