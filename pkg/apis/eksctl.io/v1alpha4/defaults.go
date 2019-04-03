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
		ng.SecurityGroups.WithLocal = NewBoolTrue()
	}
	if ng.SecurityGroups.WithShared == nil {
		ng.SecurityGroups.WithShared = NewBoolTrue()
	}

	if ng.SSH == nil {
		ng.SSH = &SSHConfig{
			Allow: NewBoolFalse(),
		}
	}

	// Enable SSH when a key is provided
	if ng.SSH.PublicKeyPath != nil {
		ng.SSH.Allow = NewBoolTrue()
	}

	if *ng.SSH.Allow && ng.SSH.PublicKeyPath == nil {
		ng.SSH.PublicKeyPath = &DefaultNodeSSHPublicKeyPath
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
		ng.IAM.WithAddonPolicies.ImageBuilder = NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.AutoScaler == nil {
		ng.IAM.WithAddonPolicies.AutoScaler = NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.ExternalDNS == nil {
		ng.IAM.WithAddonPolicies.ExternalDNS = NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.ALBIngress == nil {
		ng.IAM.WithAddonPolicies.ALBIngress = NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.EBS == nil {
		ng.IAM.WithAddonPolicies.EBS = NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.FSX == nil {
		ng.IAM.WithAddonPolicies.FSX = NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.EFS == nil {
		ng.IAM.WithAddonPolicies.EFS = NewBoolFalse()
	}

	return nil
}
