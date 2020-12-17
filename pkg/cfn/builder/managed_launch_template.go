package builder

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/weaveworks/eksctl/pkg/nodebootstrap"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
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

	userData, err := makeUserData(mng.NodeGroupBase, m.UserDataMimeBoundary)
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

		var sgIngressRules []gfnec2.SecurityGroup_Ingress
		if *mng.SSH.Allow {
			if len(mng.SSH.SourceSecurityGroupIDs) > 0 {
				for _, sgID := range mng.SSH.SourceSecurityGroupIDs {
					sgIngressRules = append(sgIngressRules, gfnec2.SecurityGroup_Ingress{
						FromPort:              sgPortSSH,
						ToPort:                sgPortSSH,
						IpProtocol:            sgProtoTCP,
						SourceSecurityGroupId: gfnt.NewString(sgID),
					})
				}
			} else {
				makeSSHIngress := func(cidrIP *gfnt.Value) gfnec2.SecurityGroup_Ingress {
					return gfnec2.SecurityGroup_Ingress{
						FromPort:   sgPortSSH,
						ToPort:     sgPortSSH,
						IpProtocol: sgProtoTCP,
						CidrIp:     cidrIP,
					}
				}

				if mng.PrivateNetworking {
					allInternalIPv4 := gfnt.NewString(m.clusterConfig.VPC.CIDR.String())
					sgIngressRules = []gfnec2.SecurityGroup_Ingress{makeSSHIngress(allInternalIPv4)}
				} else {
					sgIngressRules = []gfnec2.SecurityGroup_Ingress{
						makeSSHIngress(sgSourceAnywhereIPv4),
						{
							FromPort:   sgPortSSH,
							ToPort:     sgPortSSH,
							IpProtocol: sgProtoTCP,
							CidrIpv6:   sgSourceAnywhereIPv6,
						},
					}
				}
			}

			sshRef := m.newResource("SSH", &gfnec2.SecurityGroup{
				GroupName:            gfnt.MakeFnSubString(fmt.Sprintf("${%s}-remoteAccess", gfnt.StackName)),
				VpcId:                makeImportValue(m.clusterStackName, outputs.ClusterVPC),
				SecurityGroupIngress: sgIngressRules,
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
		if *mng.VolumeType == api.NodeVolumeTypeIO1 {
			mapping.Ebs.Iops = gfnt.NewInteger(*mng.VolumeIOPS)
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

func createMimeMessage(writer io.Writer, scripts []string, mimeBoundary string) error {
	mw := multipart.NewWriter(writer)
	if mimeBoundary != "" {
		if err := mw.SetBoundary(mimeBoundary); err != nil {
			return errors.Wrap(err, "unexpected error setting MIME boundary")
		}
	}
	fmt.Fprint(writer, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(writer, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", mw.Boundary())

	for _, script := range scripts {
		part, err := mw.CreatePart(map[string][]string{
			"Content-Type": {"text/x-shellscript", `charset="us-ascii"`},
		})

		if err != nil {
			return err
		}

		if _, err = part.Write([]byte(script)); err != nil {
			return err
		}
	}
	return mw.Close()
}

func makeUserData(ng *api.NodeGroupBase, mimeBoundary string) (string, error) {
	var (
		buf     bytes.Buffer
		scripts []string
	)

	if ng.SSH.EnableSSM != nil && *ng.SSH.EnableSSM {
		installSSMScript, err := nodebootstrap.Asset("install-ssm.al2.sh")
		if err != nil {
			return "", err
		}

		scripts = append(scripts, string(installSSMScript))
	}

	if len(ng.PreBootstrapCommands) > 0 {
		scripts = append(scripts, ng.PreBootstrapCommands...)
	}

	if ng.OverrideBootstrapCommand != nil {
		scripts = append(scripts, *ng.OverrideBootstrapCommand)
	} else if ng.MaxPodsPerNode != 0 {
		scripts = append(scripts, makeMaxPodsScript(ng.MaxPodsPerNode))
	}

	if len(scripts) == 0 {
		return "", nil
	}

	if err := createMimeMessage(&buf, scripts, mimeBoundary); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func makeMaxPodsScript(maxPods int) string {
	script := `#!/bin/sh
set -ex
sed -i -E "s/^USE_MAX_PODS=\"\\$\{USE_MAX_PODS:-true}\"/USE_MAX_PODS=false/" /etc/eks/bootstrap.sh
KUBELET_CONFIG=/etc/kubernetes/kubelet/kubelet-config.json
`
	script += fmt.Sprintf(`echo "$(jq ".maxPods=%v" $KUBELET_CONFIG)" > $KUBELET_CONFIG`, maxPods)
	return script
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
