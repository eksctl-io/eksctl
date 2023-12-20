package accessentry

import (
	"fmt"

	"github.com/kris-nova/logger"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var (
	ErrDisabledAccessEntryAPI = fmt.Errorf("access entries API is not currently enabled; please enable it using `eksctl utils update-authentication-mode --cluster <> --authentication-mode=API_AND_CONFIG_MAP`")
)

// Service is a service for access entries.
type Service struct {
	// ClusterStateGetter returns the cluster state.
	ClusterStateGetter
}

// ClusterStateGetter returns the cluster state.
type ClusterStateGetter interface {
	GetClusterState() *ekstypes.Cluster
}

// IsEnabled reports whether the cluster has access entries enabled.
func (s *Service) IsEnabled() bool {
	cluster := s.GetClusterState()
	return cluster.AccessConfig != nil && IsEnabled(cluster.AccessConfig.AuthenticationMode)
}

// IsAWSAuthDisabled reports whether the cluster has authentication mode set to API.
func (s *Service) IsAWSAuthDisabled() bool {
	accessConfig := s.GetClusterState().AccessConfig
	return accessConfig == nil || accessConfig.AuthenticationMode == ekstypes.AuthenticationModeApi
}

// IsEnabled reports whether the authenticationMode indicates that the cluster has access entries enabled.
func IsEnabled(authenticationMode ekstypes.AuthenticationMode) bool {
	return authenticationMode != ekstypes.AuthenticationModeConfigMap
}

// ValidateAPIServerAccess validates whether the API server is accessible for clusterConfig, and logs warning messages
// for operations that might fail later.
func ValidateAPIServerAccess(clusterConfig *api.ClusterConfig) error {
	if !api.IsDisabled(clusterConfig.AccessConfig.BootstrapClusterCreatorAdminPermissions) {
		return nil
	}

	const (
		apiServerConnectivityMsg = "eksctl features that require connectivity to the Kubernetes API server will fail"
		bootstrapFalseMsg        = "bootstrapClusterCreatorAdminPermissions is false"
	)
	switch clusterConfig.AccessConfig.AuthenticationMode {
	case ekstypes.AuthenticationModeConfigMap:
		if len(clusterConfig.NodeGroups) > 0 {
			return fmt.Errorf("cannot create self-managed nodegroups when authenticationMode is %s and %s", ekstypes.AuthenticationModeConfigMap, bootstrapFalseMsg)
		}
		logger.Warning("%s; %s", bootstrapFalseMsg, apiServerConnectivityMsg)
	default:
		if len(clusterConfig.AccessConfig.AccessEntries) == 0 {
			if len(clusterConfig.NodeGroups) > 0 {
				return fmt.Errorf("cannot create self-managed nodegroups when %s and no access entries are configured", bootstrapFalseMsg)
			}
			logger.Warning("%s and no access entries are configured; %s", bootstrapFalseMsg, apiServerConnectivityMsg)
			return nil
		}
		logger.Warning("%s; if no configured access entries allow access to the Kubernetes API server, %s", bootstrapFalseMsg, apiServerConnectivityMsg)
	}
	return nil
}
