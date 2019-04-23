package v1alpha4

// SetNodeGroupDefaults will set defaults for a given nodegroup
func SetNodeGroupDefaults(_ int, ng *NodeGroup) error {
	if ng.InstanceType == "" {
		ng.InstanceType = DefaultNodeType
	}
	if ng.AMIFamily == "" {
		ng.AMIFamily = DefaultNodeImageFamily
	}
	if ng.AMI == "" {
		ng.AMI = "static"
	}

	if ng.SecurityGroups == nil {
		ng.SecurityGroups = &NodeGroupSGs{
			AttachIDs: []string{},
		}
	}
	if ng.SecurityGroups.WithLocal == nil {
		ng.SecurityGroups.WithLocal = Enabled()
	}
	if ng.SecurityGroups.WithShared == nil {
		ng.SecurityGroups.WithShared = Enabled()
	}

	if ng.SSH == nil {
		ng.SSH = &NodeGroupSSH{
			Allow: Disabled(),
		}
	}

	numSSHFlagsEnabled := countEnabledFields(
		ng.SSH.PublicKeyName,
		ng.SSH.PublicKeyPath,
		ng.SSH.PublicKey)

	if numSSHFlagsEnabled > 0 {
		ng.SSH.Allow = Enabled()
	} else {
		if ng.SSH.Allow != nil && *ng.SSH.Allow {
			ng.SSH.PublicKeyPath = &DefaultNodeSSHPublicKeyPath
		} else {
			ng.SSH.Allow = Disabled()
		}
	}

	if ng.VolumeSize > 0 {
		if ng.VolumeType == "" {
			ng.VolumeType = DefaultNodeVolumeType
		}
	}

	if ng.IAM == nil {
		ng.IAM = &NodeGroupIAM{}
	}
	if ng.IAM.WithAddonPolicies.ImageBuilder == nil {
		ng.IAM.WithAddonPolicies.ImageBuilder = Disabled()
	}
	if ng.IAM.WithAddonPolicies.AutoScaler == nil {
		ng.IAM.WithAddonPolicies.AutoScaler = Disabled()
	}
	if ng.IAM.WithAddonPolicies.ExternalDNS == nil {
		ng.IAM.WithAddonPolicies.ExternalDNS = Disabled()
	}
	if ng.IAM.WithAddonPolicies.ALBIngress == nil {
		ng.IAM.WithAddonPolicies.ALBIngress = Disabled()
	}
	if ng.IAM.WithAddonPolicies.EBS == nil {
		ng.IAM.WithAddonPolicies.EBS = Disabled()
	}
	if ng.IAM.WithAddonPolicies.FSX == nil {
		ng.IAM.WithAddonPolicies.FSX = Disabled()
	}
	if ng.IAM.WithAddonPolicies.EFS == nil {
		ng.IAM.WithAddonPolicies.EFS = Disabled()
	}

	return nil
}
