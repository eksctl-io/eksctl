package v1alpha5

import (
	"errors"
	"fmt"
	"slices"
)

// Values for `AutoModeNodePool`.
const (
	AutoModeNodePoolGeneralPurpose = "general-purpose"
	AutoModeNodePoolSystem         = "system"
)

// AutoModeKnownNodePools is a slice of known node pools for Auto Mode.
var AutoModeKnownNodePools = []string{AutoModeNodePoolGeneralPurpose, AutoModeNodePoolSystem}

type AutoModeConfig struct {
	// Enabled enables or disables Auto Mode.
	Enabled *bool `json:"enabled,omitempty"`
	// NodeRoleARN is the node role to use for nodes launched by Auto Mode.
	NodeRoleARN ARN `json:"nodeRoleARN,omitempty"`
	// NodePools is a list of node pools to create.
	NodePools *[]string `json:"nodePools,omitempty"`
}

// HasNodePools reports whether any node pools are specified.
func (a *AutoModeConfig) HasNodePools() bool {
	return a.NodePools != nil && len(*a.NodePools) > 0
}

// ValidateAutoModeConfig validates the Auto Mode config.
func ValidateAutoModeConfig(clusterConfig *ClusterConfig) error {
	autoModeConfig := clusterConfig.AutoModeConfig
	if autoModeConfig == nil {
		return nil
	}
	if IsEnabled(autoModeConfig.Enabled) {
		if clusterConfig.IsControlPlaneOnOutposts() {
			return errors.New("Auto Mode is not supported on Outposts")
		}
		if autoModeConfig.NodePools != nil {
			if len(*autoModeConfig.NodePools) == 0 && !autoModeConfig.NodeRoleARN.IsZero() {
				return errors.New("cannot specify autoModeConfig.nodeRoleARN when autoModeConfig.nodePools is empty")
			}
			seenNodePools := map[string]struct{}{}
			for _, np := range *autoModeConfig.NodePools {
				if _, ok := seenNodePools[np]; ok {
					return fmt.Errorf("found duplicate node pool: %q", np)
				}
				if !slices.Contains(AutoModeKnownNodePools, np) {
					return fmt.Errorf("invalid NodePool %q", np)
				}
				seenNodePools[np] = struct{}{}
			}
		}
	} else if !autoModeConfig.NodeRoleARN.IsZero() || autoModeConfig.HasNodePools() {
		return errors.New("cannot set autoModeConfig.nodeRoleARN or autoModeConfig.nodePools when Auto Mode is disabled")
	}
	return nil
}
