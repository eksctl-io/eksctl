package builder

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (m *ManagedNodeGroupResourceSet) makeLaunchTemplateData(ctx context.Context) (*gfnec2.LaunchTemplate_LaunchTemplateData, error) {
	mng := m.nodeGroup
	launchTemplateData := &gfnec2.LaunchTemplate_LaunchTemplateData{
		TagSpecifications: makeTags(mng.NodeGroupBase, m.clusterConfig.Metadata),
		MetadataOptions:   makeMetadataOptions(mng.NodeGroupBase),
	}

	userData, err := m.bootstrapper.UserData()
	if err != nil {
		return nil, err
	}
	if userData != "" {
		launchTemplateData.UserData = gfnt.NewString(userData)
	}

	securityGroupIDs := m.vpcImporter.SecurityGroups()
	for _, sgID := range mng.SecurityGroups.AttachIDs {
		securityGroupIDs = append(securityGroupIDs, gfnt.NewString(sgID))
	}

	if mng.AMI != "" {
		launchTemplateData.ImageId = gfnt.NewString(mng.AMI)
	}

	if mng.SSH != nil && api.IsSetAndNonEmptyString(mng.SSH.PublicKeyName) {
		launchTemplateData.KeyName = gfnt.NewString(*mng.SSH.PublicKeyName)

		if *mng.SSH.Allow {
			vpcID := m.vpcImporter.VPC()
			sshRef := m.newResource("SSH", &gfnec2.SecurityGroup{
				GroupName:            gfnt.MakeFnSubString(fmt.Sprintf("${%s}-remoteAccess", gfnt.StackName)),
				VpcId:                vpcID,
				SecurityGroupIngress: makeSSHIngressRules(mng.NodeGroupBase, m.clusterConfig.VPC.CIDR.String(), fmt.Sprintf("managed worker nodes in group %s", mng.Name)),
				GroupDescription:     gfnt.NewString("Allow SSH access"),
			})
			securityGroupIDs = append(securityGroupIDs, sshRef)
		}
	}

	if api.IsEnabled(mng.EFAEnabled) {
		// we don't want to touch the network interfaces at all if we have a
		// managed nodegroup, unless EFA is enabled
		desc := "worker nodes in group " + m.nodeGroup.Name
		efaSG := m.addEFASecurityGroup(m.vpcImporter.VPC(), m.clusterConfig.Metadata.Name, desc)
		securityGroupIDs = append(securityGroupIDs, efaSG)
		if err := buildNetworkInterfaces(ctx, launchTemplateData, mng.InstanceTypeList(), true, securityGroupIDs, m.ec2API); err != nil {
			return nil, errors.Wrap(err, "couldn't build network interfaces for launch template data")
		}
		if mng.Placement == nil {
			groupName := m.newResource("NodeGroupPlacementGroup", &gfnec2.PlacementGroup{
				Strategy: gfnt.NewString("cluster"),
			})
			launchTemplateData.Placement = &gfnec2.LaunchTemplate_Placement{
				GroupName: groupName,
			}
		}
	} else {
		launchTemplateData.SecurityGroupIds = gfnt.NewSlice(securityGroupIDs...)
	}

	if mng.EBSOptimized != nil {
		launchTemplateData.EbsOptimized = gfnt.NewBoolean(*mng.EBSOptimized)
	}

	if mng.Placement != nil {
		launchTemplateData.Placement = &gfnec2.LaunchTemplate_Placement{
			GroupName: gfnt.NewString(mng.Placement.GroupName),
		}
	}

	if mng.EnableDetailedMonitoring != nil {
		launchTemplateData.Monitoring = &gfnec2.LaunchTemplate_Monitoring{
			Enabled: gfnt.NewBoolean(*mng.EnableDetailedMonitoring),
		}
	}

	launchTemplateData.BlockDeviceMappings = makeBlockDeviceMappings(mng.NodeGroupBase)

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

	var launchTemplateTagSpecs []gfnec2.LaunchTemplate_TagSpecification

	launchTemplateTagSpecs = append(launchTemplateTagSpecs,
		gfnec2.LaunchTemplate_TagSpecification{
			ResourceType: gfnt.NewString("instance"),
			Tags:         cfnTags,
		}, gfnec2.LaunchTemplate_TagSpecification{
			ResourceType: gfnt.NewString("volume"),
			Tags:         cfnTags,
		}, gfnec2.LaunchTemplate_TagSpecification{
			ResourceType: gfnt.NewString("network-interface"),
			Tags:         cfnTags,
		})

	return launchTemplateTagSpecs
}
