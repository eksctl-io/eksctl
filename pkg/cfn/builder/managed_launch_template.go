package builder

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

func (m *ManagedNodeGroupResourceSet) makeLaunchTemplateData() (*gfnec2.LaunchTemplate_LaunchTemplateData, error) {
	mng := m.nodeGroup

	launchTemplateData := &gfnec2.LaunchTemplate_LaunchTemplateData{
		TagSpecifications: makeTags(mng.NodeGroupBase, m.clusterConfig.Metadata),
		MetadataOptions:   makeMetadataOptions(mng.NodeGroupBase),
	}

	userData, err := nodebootstrap.MakeManagedUserData(mng.NodeGroupBase, m.UserDataMimeBoundary)
	if err != nil {
		return nil, err
	}
	if userData != "" {
		launchTemplateData.UserData = gfnt.NewString(userData)
	}

	securityGroupIDs := gfnt.Slice{makeImportValue(m.clusterStackName, outputs.ClusterDefaultSecurityGroup)}
	for _, sgID := range mng.SecurityGroups.AttachIDs {
		securityGroupIDs = append(securityGroupIDs, gfnt.NewString(sgID))
	}

	if mng.AMI != "" {
		launchTemplateData.ImageId = gfnt.NewString(mng.AMI)
	}

	if mng.SSH != nil && api.IsSetAndNonEmptyString(mng.SSH.PublicKeyName) {
		launchTemplateData.KeyName = gfnt.NewString(*mng.SSH.PublicKeyName)

		if *mng.SSH.Allow {
			sshRef := m.newResource("SSH", &gfnec2.SecurityGroup{
				GroupName:            gfnt.MakeFnSubString(fmt.Sprintf("${%s}-remoteAccess", gfnt.StackName)),
				VpcId:                makeImportValue(m.clusterStackName, outputs.ClusterVPC),
				SecurityGroupIngress: makeSSHIngressRules(mng.NodeGroupBase, m.clusterConfig.VPC.CIDR.String(), fmt.Sprintf("managed worker nodes in group %s", mng.Name)),
				GroupDescription:     gfnt.NewString("Allow SSH access"),
			})
			securityGroupIDs = append(securityGroupIDs, sshRef)
		}
	}

	launchTemplateData.SecurityGroupIds = gfnt.NewValue(securityGroupIDs)

	if mng.EBSOptimized != nil {
		launchTemplateData.EbsOptimized = gfnt.NewBoolean(*mng.EBSOptimized)
	}

	if volumeSize := mng.VolumeSize; volumeSize != nil && *volumeSize > 0 {
		mapping := gfnec2.LaunchTemplate_BlockDeviceMapping{
			Ebs: &gfnec2.LaunchTemplate_Ebs{
				VolumeSize: gfnt.NewInteger(*volumeSize),
				VolumeType: gfnt.NewString(*mng.VolumeType),
			},
		}
		if mng.VolumeEncrypted != nil {
			mapping.Ebs.Encrypted = gfnt.NewBoolean(*mng.VolumeEncrypted)
		}
		if api.IsSetAndNonEmptyString(mng.VolumeKmsKeyID) {
			mapping.Ebs.KmsKeyId = gfnt.NewString(*mng.VolumeKmsKeyID)
		}

		if *mng.VolumeType == api.NodeVolumeTypeIO1 || *mng.VolumeType == api.NodeVolumeTypeGP3 {
			if mng.VolumeIOPS != nil {
				mapping.Ebs.Iops = gfnt.NewInteger(*mng.VolumeIOPS)
			}
		}

		if *mng.VolumeType == api.NodeVolumeTypeGP3 && mng.VolumeThroughput != nil {
			mapping.Ebs.Throughput = gfnt.NewInteger(*mng.VolumeThroughput)
		}

		if mng.VolumeName != nil {
			mapping.DeviceName = gfnt.NewString(*mng.VolumeName)
		} else {
			mapping.DeviceName = gfnt.NewString("/dev/xvda")
		}

		launchTemplateData.BlockDeviceMappings = []gfnec2.LaunchTemplate_BlockDeviceMapping{mapping}
	}

	if mng.Placement != nil {
		launchTemplateData.Placement = &gfnec2.LaunchTemplate_Placement{
			GroupName: gfnt.NewString(mng.Placement.GroupName),
		}
	}

	return launchTemplateData, nil
}

func makeSSHIngressRules(n *api.NodeGroupBase, vpcCIDR, description string) []gfnec2.SecurityGroup_Ingress {
	var sgIngressRules []gfnec2.SecurityGroup_Ingress
	if *n.SSH.Allow {
		if len(n.SSH.SourceSecurityGroupIDs) > 0 {
			for _, sgID := range n.SSH.SourceSecurityGroupIDs {
				sgIngressRules = append(sgIngressRules, gfnec2.SecurityGroup_Ingress{
					FromPort:              sgPortSSH,
					ToPort:                sgPortSSH,
					IpProtocol:            sgProtoTCP,
					SourceSecurityGroupId: gfnt.NewString(sgID),
				})
			}
		} else {
			makeSSHIngress := func(cidrIP *gfnt.Value, sshDesc string) gfnec2.SecurityGroup_Ingress {
				return gfnec2.SecurityGroup_Ingress{
					FromPort:    sgPortSSH,
					ToPort:      sgPortSSH,
					IpProtocol:  sgProtoTCP,
					CidrIp:      cidrIP,
					Description: gfnt.NewString(sshDesc),
				}
			}

			sshDesc := "Allow SSH access to " + description

			if n.PrivateNetworking {
				allInternalIPv4 := gfnt.NewString(vpcCIDR)
				sgIngressRules = []gfnec2.SecurityGroup_Ingress{makeSSHIngress(allInternalIPv4, sshDesc+" (private, only inside VPC)")}
			} else {
				sgIngressRules = append(sgIngressRules,
					makeSSHIngress(sgSourceAnywhereIPv4, sshDesc),
					gfnec2.SecurityGroup_Ingress{
						CidrIpv6:    sgSourceAnywhereIPv6,
						Description: gfnt.NewString(sshDesc),
						IpProtocol:  sgProtoTCP,
						FromPort:    sgPortSSH,
						ToPort:      sgPortSSH,
					})
			}
		}
	}
	return sgIngressRules
}

func makeTags(ng *api.NodeGroupBase, meta *api.ClusterMeta) []gfnec2.LaunchTemplate_TagSpecification {
	cfnTags := []cloudformation.Tag{
		{
			Key:   gfnt.NewString("Name"),
			Value: gfnt.NewString(generateNodeName(ng, meta)),
		},
	}
	for k, v := range ng.Tags {
		cfnTags = append(cfnTags, cloudformation.Tag{
			Key:   gfnt.NewString(k),
			Value: gfnt.NewString(v),
		})
	}
	return []gfnec2.LaunchTemplate_TagSpecification{
		{
			ResourceType: gfnt.NewString("instance"),
			Tags:         cfnTags,
		},
	}
}
