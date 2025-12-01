package create

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func TestConfigureCreateCapabilityCmd(t *testing.T) {
	flagGrouping := cmdutils.NewGrouping()
	cmd := &cmdutils.Cmd{
		CobraCommand: &cobra.Command{},
	}
	cmd.FlagSetGroup = flagGrouping.New(cmd.CobraCommand)

	createCapabilityCmd(cmd)

	if cmd.CobraCommand.Use != "capability" {
		t.Errorf("Expected command use 'capability', got %s", cmd.CobraCommand.Use)
	}

	if cmd.CobraCommand.Short == "" {
		t.Error("Expected command to have a short description")
	}
}
