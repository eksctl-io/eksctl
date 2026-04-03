package v1alpha5

import "slices"

var KnownAddons = map[string]struct {
	IsDefault             bool
	CreateBeforeNodeGroup bool
	IsDefaultAutoMode     bool
	DontRequireWait       bool
	ExcludedRegions       []string
}{
	VPCCNIAddon: {
		IsDefault:             true,
		CreateBeforeNodeGroup: true,
	},
	KubeProxyAddon: {
		IsDefault:             true,
		CreateBeforeNodeGroup: true,
	},
	CoreDNSAddon: {
		IsDefault:             true,
		CreateBeforeNodeGroup: true,
	},
	PodIdentityAgentAddon: {
		CreateBeforeNodeGroup: true,
	},
	AWSEBSCSIDriverAddon: {},
	AWSEFSCSIDriverAddon: {},
	MetricsServerAddon: {
		CreateBeforeNodeGroup: false, // Create after nodegroup so we get scheduled on Fargate profiles.
		IsDefaultAutoMode:     true,
		ExcludedRegions: []string{
			RegionCNNorthwest1,
			RegionCNNorth1,
			RegionUSISOEast1,
			RegionUSISOWest1,
			RegionUSISOBEast1,
			RegionUSGovWest1,
			RegionUSGovEast1,
			RegionUSISOFEast1,
			RegionUSISOFSouth1,
			RegionEUISOEWest1,
		},
		// Don't require waiting for metrics-server to be up if it's the only add-on to wait for.
		// This is because this add-on is installed by default and we don't
		// want to make a large change blocking cluster creating on this add-on coming up successfully in setups.
		DontRequireWait: true,
	},
}

// HasDefaultNonAutoAddon reports whether addons contains at least one non-auto mode default addon
func HasDefaultNonAutoAddon(addons []*Addon) bool {
	for _, addon := range addons {
		addonConfig, ok := KnownAddons[addon.Name]
		if ok && addonConfig.IsDefault && !addonConfig.IsDefaultAutoMode {
			return true
		}
	}
	return false
}

// HasAllDefaultAddons reports whether addonNames contains all default addons.
func HasAllDefaultAddons(addonNames []string) bool {
	for addonName, addon := range KnownAddons {
		if addon.IsDefault && !slices.Contains(addonNames, addonName) {
			return false
		}
	}
	return true
}
