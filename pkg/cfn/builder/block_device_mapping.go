package builder

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

const rootDevice = "/dev/xvda"

func makeBlockDeviceMappings(ng *api.NodeGroupBase) []gfnec2.LaunchTemplate_BlockDeviceMapping {
	volumeSize := ng.VolumeSize
	if volumeSize == nil || *volumeSize == 0 {
		return nil
	}

	mapping := gfnec2.LaunchTemplate_BlockDeviceMapping{
		Ebs: &gfnec2.LaunchTemplate_Ebs{
			VolumeSize: gfnt.NewInteger(*volumeSize),
			VolumeType: gfnt.NewString(*ng.VolumeType),
		},
	}
	if ng.VolumeEncrypted != nil {
		mapping.Ebs.Encrypted = gfnt.NewBoolean(*ng.VolumeEncrypted)
	}
	if api.IsSetAndNonEmptyString(ng.VolumeKmsKeyID) {
		mapping.Ebs.KmsKeyId = gfnt.NewString(*ng.VolumeKmsKeyID)
	}

	if (*ng.VolumeType == api.NodeVolumeTypeIO1 || *ng.VolumeType == api.NodeVolumeTypeGP3) && ng.VolumeIOPS != nil {
		mapping.Ebs.Iops = gfnt.NewInteger(*ng.VolumeIOPS)
	}

	if *ng.VolumeType == api.NodeVolumeTypeGP3 && ng.VolumeThroughput != nil {
		mapping.Ebs.Throughput = gfnt.NewInteger(*ng.VolumeThroughput)
	}

	if ng.VolumeName != nil {
		mapping.DeviceName = gfnt.NewString(*ng.VolumeName)
	} else {
		mapping.DeviceName = gfnt.NewString(rootDevice)
	}

	mappings := []gfnec2.LaunchTemplate_BlockDeviceMapping{mapping}

	if api.IsEnabled(ng.VolumeEncrypted) && ng.AdditionalEncryptedVolume != "" {
		mappings = append(mappings, gfnec2.LaunchTemplate_BlockDeviceMapping{
			DeviceName: gfnt.NewString(ng.AdditionalEncryptedVolume),
			Ebs: &gfnec2.LaunchTemplate_Ebs{
				Encrypted: gfnt.NewBoolean(*ng.VolumeEncrypted),
				KmsKeyId:  mapping.Ebs.KmsKeyId,
			},
		})
	}

	return mappings
}
