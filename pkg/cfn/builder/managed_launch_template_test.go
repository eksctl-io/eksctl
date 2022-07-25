package builder

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	. "github.com/benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/goformation/v4"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"
)

type mngCase struct {
	ng                *api.ManagedNodeGroup
	resourcesFilename string
	mockFetcherFn     func(*mockprovider.MockProvider)

	hasUserData bool
	errMsg      string
}

var _ = Describe("ManagedNodeGroup builder", func() {
	DescribeTable("Add resources", func(m *mngCase) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "lt"
		api.SetManagedNodeGroupDefaults(m.ng, clusterConfig.Metadata)
		Expect(api.ValidateManagedNodeGroup(0, m.ng)).To(Succeed())

		provider := mockprovider.NewMockProvider()
		if m.mockFetcherFn != nil {
			m.mockFetcherFn(provider)
		}

		fakeVPCImporter := new(vpcfakes.FakeImporter)
		fakeVPCImporter.VPCReturns(gfnt.MakeFnImportValueString("eksctl-lt::VPC"))
		fakeVPCImporter.SecurityGroupsReturns(gfnt.Slice{gfnt.MakeFnImportValueString("eksctl-lt::ClusterSecurityGroupId")})
		fakeVPCImporter.SubnetsPublicReturns(gfnt.MakeFnSplit(",", gfnt.MakeFnImportValueString("eksctl-lt::SubnetsPublic")))

		bootstrapper := &fakes.FakeBootstrapper{}
		bootstrapper.UserDataStub = func() (string, error) {
			if !m.hasUserData {
				return "", nil
			}
			userData := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`/etc/eks/bootstrap.sh %s`, clusterConfig.Metadata.Name)))
			return userData, nil
		}

		stack := NewManagedNodeGroup(provider.MockEC2(), clusterConfig, m.ng, NewLaunchTemplateFetcher(provider.MockEC2()), bootstrapper, false, fakeVPCImporter)
		err := stack.AddAllResources(context.Background())
		if m.errMsg != "" {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(m.errMsg))
			return
		}

		Expect(err).NotTo(HaveOccurred())
		bytes, err := stack.RenderJSON()
		Expect(err).NotTo(HaveOccurred())

		template, err := goformation.ParseJSON(bytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(template).NotTo(BeNil())

		actual, err := json.Marshal(template.Resources)
		Expect(err).NotTo(HaveOccurred())

		expected, err := os.ReadFile(path.Join("testdata", "launch_template", m.resourcesFilename))
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(MatchOrderedJSON(expected, WithUnorderedListKeys("Tags")))

	},
		Entry("No custom AMI", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "standard",
					InstanceType: "m5.xlarge",
				},
			},
			resourcesFilename: "standard.json",
		}),
		Entry("Custom AMI", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "custom-ami",
					InstanceType: "m5.xlarge",
					AMI:          "ami-custom",
					OverrideBootstrapCommand: aws.String(`
#!/bin/bash
set -ex
B64_CLUSTER_CA=dGVzdAo=
API_SERVER_URL=https://test.com
/etc/eks/bootstrap.sh launch-template --b64-cluster-ca $B64_CLUSTER_CA --apiserver-endpoint $API_SERVER_URL
`),
				},
			},
			hasUserData: true,

			resourcesFilename: "custom_ami.json",
		}),

		Entry("Launch Template", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "custom-template",
				},
				LaunchTemplate: &api.LaunchTemplate{
					ID: "lt-1234",
				},
			},
			mockFetcherFn: mockLaunchTemplate(func(input *ec2.DescribeLaunchTemplateVersionsInput) bool {
				return *input.LaunchTemplateId == "lt-1234" && input.Versions[0] == "$Default"
			}, &ec2types.ResponseLaunchTemplateData{
				InstanceType: ec2types.InstanceTypeT2Medium,
				KeyName:      aws.String("key-name"),
			}),

			resourcesFilename: "launch_template.json",
		}),

		Entry("Launch Template With Additional Volumes", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "extra-volumes",
					AdditionalVolumes: []*api.VolumeMapping{
						{
							VolumeSize:      aws.Int(20),
							VolumeType:      aws.String(api.NodeVolumeTypeGP3),
							VolumeName:      aws.String("/foo/bar-add-1"),
							VolumeEncrypted: aws.Bool(true),
							SnapshotID:      aws.String("snapshot-id"),
						},
					},
				},
			},
			resourcesFilename: "launch_template_additional_volumes.json",
		}),

		Entry("Launch Template with volumes missing volume size", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "extra-volumes",
					AdditionalVolumes: []*api.VolumeMapping{
						{
							VolumeType:      aws.String(api.NodeVolumeTypeGP3),
							VolumeName:      aws.String("/foo/bar-add-1"),
							VolumeEncrypted: aws.Bool(true),
						},
					},
				},
			},
			resourcesFilename: "launch_template_additional_volumes_missing_size.json",
		}),

		Entry("Launch Template with custom AMI", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "template-custom-ami",
				},
				LaunchTemplate: &api.LaunchTemplate{
					ID:      "lt-1234",
					Version: aws.String("2"),
				},
			},
			mockFetcherFn: mockLaunchTemplate(func(input *ec2.DescribeLaunchTemplateVersionsInput) bool {
				return *input.LaunchTemplateId == "lt-1234" && input.Versions[0] == "2"
			}, &ec2types.ResponseLaunchTemplateData{
				ImageId:      aws.String("ami-1234"),
				InstanceType: ec2types.InstanceTypeT2Medium,
				KeyName:      aws.String("key-name"),
				UserData:     aws.String("bootstrap.sh"),
			}),

			resourcesFilename: "launch_template_custom_ami.json",
		}),

		Entry("SSH enabled", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "ssh-enabled",
					SSH: &api.NodeGroupSSH{
						Allow:         api.Enabled(),
						PublicKeyName: aws.String("test-keypair"),
					},
				},
			},
			hasUserData: true,

			resourcesFilename: "ssh_enabled.json",
		}),

		Entry("SSH configured but allowed=false", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "ssh-disabled",
					SSH: &api.NodeGroupSSH{
						Allow:         api.Disabled(),
						PublicKeyName: aws.String("test-keypair"),
					},
				},
			},
			hasUserData: true,

			// The SG should not be created
			resourcesFilename: "ssh_disabled.json",
		}),

		Entry("With placement group", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "standard",
					InstanceType: "m5.xlarge",
					Placement: &api.Placement{
						GroupName: "test",
					},
				},
			},
			resourcesFilename: "placement.json",
		}),

		Entry("With Spot instances", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "spot",
				},
				Spot:          true,
				InstanceTypes: []string{"c3.large", "c4.large", "c5.large", "c5d.large", "c5n.large", "c5a.large"},
			},
			resourcesFilename: "spot.json",
		}),

		Entry("Without instance type set", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "template-custom-ami",
				},
				LaunchTemplate: &api.LaunchTemplate{
					ID:      "lt-1234",
					Version: aws.String("2"),
				},
			},
			mockFetcherFn: mockLaunchTemplate(func(input *ec2.DescribeLaunchTemplateVersionsInput) bool {
				return *input.LaunchTemplateId == "lt-1234" && input.Versions[0] == "2"
			}, &ec2types.ResponseLaunchTemplateData{
				ImageId:  aws.String("ami-1234"),
				KeyName:  aws.String("key-name"),
				UserData: aws.String("bootstrap.sh"),
			}),
			errMsg: "instance type must be set in the launch template",
		}),

		Entry("With instance type set", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "template-custom-ami",
				},
				InstanceTypes: []string{"t2.medium"},
				LaunchTemplate: &api.LaunchTemplate{
					ID:      "lt-1234",
					Version: aws.String("2"),
				},
			},
			mockFetcherFn: mockLaunchTemplate(func(input *ec2.DescribeLaunchTemplateVersionsInput) bool {
				return *input.LaunchTemplateId == "lt-1234" && input.Versions[0] == "2"
			}, &ec2types.ResponseLaunchTemplateData{
				ImageId:      aws.String("ami-1234"),
				InstanceType: ec2types.InstanceTypeM5Large,
				KeyName:      aws.String("key-name"),
				UserData:     aws.String("bootstrap.sh"),
			}),
			errMsg: "instance type must not be set in the launch template",
		}),

		Entry("With launch template and multiple instance types", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "template-custom-ami",
				},
				InstanceTypes: []string{"c3.large", "c4.large", "c5.large", "c5d.large", "c5n.large", "c5a.large"},
				LaunchTemplate: &api.LaunchTemplate{
					ID:      "lt-1234",
					Version: aws.String("3"),
				},
			},
			mockFetcherFn: mockLaunchTemplate(func(input *ec2.DescribeLaunchTemplateVersionsInput) bool {
				return *input.LaunchTemplateId == "lt-1234" && input.Versions[0] == "3"
			}, &ec2types.ResponseLaunchTemplateData{
				ImageId:  aws.String("ami-1234"),
				KeyName:  aws.String("key-name"),
				UserData: aws.String("bootstrap.sh"),
			}),
			resourcesFilename: "lt_instance_types.json",
		}),

		Entry("Bottlerocket AMI Family with defaults", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "bottlerocket",
					AMIFamily:    api.NodeImageFamilyBottlerocket,
					InstanceType: "m5.xlarge",
				},
			},
			resourcesFilename: "bottlerocket.json",
		}),

		Entry("Bottlerocket with volumeSize set", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "bottlerocket-volume",
					AMIFamily:    api.NodeImageFamilyBottlerocket,
					InstanceType: "m5.xlarge",
					VolumeSize:   aws.Int(142),
				},
			},
			resourcesFilename: "bottlerocket_volume.json",
		}),

		Entry("CapacityReservationID is set", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "capacity-test",
					AMIFamily:    api.NodeImageFamilyBottlerocket,
					InstanceType: "m5.xlarge",
					VolumeSize:   aws.Int(142),
					CapacityReservation: &api.CapacityReservation{
						CapacityReservationTarget: &api.CapacityReservationTarget{
							CapacityReservationID: aws.String("res-id"),
						},
					},
				},
			},
			resourcesFilename: "launch_template_with_capacity_reservation_id.json",
		}),
		Entry("CapacityReservationResourceGroupAN is set", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "capacity-test",
					AMIFamily:    api.NodeImageFamilyBottlerocket,
					InstanceType: "m5.xlarge",
					VolumeSize:   aws.Int(142),
					CapacityReservation: &api.CapacityReservation{
						CapacityReservationTarget: &api.CapacityReservationTarget{
							CapacityReservationResourceGroupARN: aws.String("group-arn"),
						},
					},
				},
			},
			resourcesFilename: "launch_template_with_capacity_reservation_resource_group_arn.json",
		}),
		Entry("CapacityReservationPreference is set", &mngCase{
			ng: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "capacity-test",
					AMIFamily:    api.NodeImageFamilyBottlerocket,
					InstanceType: "m5.xlarge",
					VolumeSize:   aws.Int(142),
					CapacityReservation: &api.CapacityReservation{
						CapacityReservationPreference: aws.String("open"),
					},
				},
			},
			resourcesFilename: "launch_template_with_capacity_reservation_preference.json",
		}),
	)
})

func mockLaunchTemplate(matcher func(*ec2.DescribeLaunchTemplateVersionsInput) bool, lt *ec2types.ResponseLaunchTemplateData) func(provider *mockprovider.MockProvider) {
	return func(provider *mockprovider.MockProvider) {
		provider.MockEC2().On("DescribeLaunchTemplateVersions", mock.Anything, mock.MatchedBy(matcher)).
			Return(&ec2.DescribeLaunchTemplateVersionsOutput{
				LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
					{
						LaunchTemplateData: lt,
					},
				},
			}, nil)
	}
}
