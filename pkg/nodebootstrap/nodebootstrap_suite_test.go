package nodebootstrap_test

import (
	"testing"

	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestNodebootstrap(t *testing.T) {
	testutils.RegisterAndRun(t)
}

func decode(userData string) *cloudconfig.CloudConfig {
	cloudCfg, err := cloudconfig.DecodeCloudConfig(userData)
	Expect(err).NotTo(HaveOccurred())

	return cloudCfg
}

func newBootstrapper(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) nodebootstrap.Bootstrapper {
	bootstrapper, err := nodebootstrap.NewBootstrapper(clusterConfig, ng)
	Expect(err).NotTo(HaveOccurred())
	return bootstrapper
}
