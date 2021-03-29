package nodebootstrap_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/cloudconfig"
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
