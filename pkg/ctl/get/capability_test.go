package get

import (
	"testing"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func TestGetCapabilityCmd(t *testing.T) {
	cmd := &cmdutils.Cmd{
		ClusterConfig: api.NewClusterConfig(),
	}
	
	getCapabilityCmd(cmd)
	
	if cmd.CobraCommand.Use != "capability" {
		t.Errorf("Expected command use 'capability', got %s", cmd.CobraCommand.Use)
	}
	
	if cmd.CobraCommand.Short == "" {
		t.Error("Expected command to have a short description")
	}
	
	if cmd.CobraCommand.RunE == nil {
		t.Error("Expected command to have a RunE function")
	}
}