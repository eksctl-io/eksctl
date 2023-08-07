package cloudwatch

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"k8s.io/apimachinery/pkg/util/sets"
)

// LogEnabledFeatures logs enabled features
func LogEnabledFeatures(clusterConfig *api.ClusterConfig) {
	if clusterConfig.HasClusterEndpointAccess() && api.EndpointsEqual(*clusterConfig.VPC.ClusterEndpoints, *api.ClusterEndpointAccessDefaults()) {
		logger.Info(clusterConfig.DefaultEndpointsMsg())
	} else {
		logger.Info(clusterConfig.CustomEndpointsMsg())
	}

	if !clusterConfig.HasClusterCloudWatchLogging() {
		logger.Info("CloudWatch logging will not be enabled for cluster %q in %q", clusterConfig.Metadata.Name, clusterConfig.Metadata.Region)
		logger.Info("you can enable it with 'eksctl utils update-cluster-logging --enable-types={SPECIFY-YOUR-LOG-TYPES-HERE (e.g. all)} --region=%s --cluster=%s'", clusterConfig.Metadata.Region, clusterConfig.Metadata.Name)
		return
	}

	all := sets.NewString(api.SupportedCloudWatchClusterLogTypes()...)

	enabled := sets.NewString()
	if clusterConfig.HasClusterCloudWatchLogging() {
		enabled.Insert(clusterConfig.CloudWatch.ClusterLogging.EnableTypes...)
	}

	disabled := all.Difference(enabled)

	describeEnabledTypes := "no types enabled"
	if enabled.Len() > 0 {
		describeEnabledTypes = fmt.Sprintf("enabled types: %s", strings.Join(enabled.List(), ", "))
	}

	describeDisabledTypes := "no types disabled"
	if disabled.Len() > 0 {
		describeDisabledTypes = fmt.Sprintf("disabled types: %s", strings.Join(disabled.List(), ", "))
	}

	logger.Info("configuring CloudWatch logging for cluster %q in %q (%s & %s)",
		clusterConfig.Metadata.Name, clusterConfig.Metadata.Region, describeEnabledTypes, describeDisabledTypes,
	)
}
