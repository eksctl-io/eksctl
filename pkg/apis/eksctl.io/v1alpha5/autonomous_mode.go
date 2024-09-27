package v1alpha5

import (
	"errors"
	"fmt"
	"slices"
)

// Values for `AutonomousModeNodePool`.
const (
	AutonomousModeNodePoolGeneralPurpose = "general-purpose"
	AutonomousModeNodePoolSystem         = "system"
)

// AutonomousModeKnownNodePools is a slice of known node pools for Autonomous Mode.
var AutonomousModeKnownNodePools = []string{AutonomousModeNodePoolGeneralPurpose, AutonomousModeNodePoolSystem}

type AutonomousModeConfig struct {
	// Enabled enables or disables Autonomous Mode.
	Enabled *bool `json:"enabled,omitempty"`
	// NodeRoleARN is the node role to use for nodes launched by Autonomous Mode.
	NodeRoleARN ARN `json:"nodeRoleARN,omitempty"`
	// NodePools is a list of node pools to create.
	NodePools *[]string `json:"nodePools,omitempty"`
}

// HasNodePools reports whether any node pools are specified.
func (a *AutonomousModeConfig) HasNodePools() bool {
	return a.NodePools != nil && len(*a.NodePools) > 0
}

// ValidateAutonomousModeConfig validates the autonomous mode config.
func ValidateAutonomousModeConfig(clusterConfig *ClusterConfig) error {
	autonomousModeConfig := clusterConfig.AutonomousModeConfig
	if autonomousModeConfig == nil {
		return nil
	}
	if IsEnabled(autonomousModeConfig.Enabled) {
		if clusterConfig.IsControlPlaneOnOutposts() {
			return errors.New("Autonomous Mode is not supported on Outposts")
		}
		if autonomousModeConfig.NodePools != nil {
			if len(*autonomousModeConfig.NodePools) == 0 && !autonomousModeConfig.NodeRoleARN.IsZero() {
				return errors.New("cannot specify autonomousModeConfig.nodeRoleARN when autonomousModeConfig.nodePools is empty")
			}
			for _, np := range *autonomousModeConfig.NodePools {
				if !slices.Contains(AutonomousModeKnownNodePools, np) {
					return fmt.Errorf("invalid NodePool %q", np)
				}
			}
		}
	} else if !autonomousModeConfig.NodeRoleARN.IsZero() || autonomousModeConfig.HasNodePools() {
		return errors.New("cannot set autonomousModeConfig.nodeRoleARN or autonomousModeConfig.nodePools when Autonomous Mode is disabled")
	}
	return nil
}
