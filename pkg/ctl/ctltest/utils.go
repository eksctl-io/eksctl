package ctltest

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type CommandFunction func(*cmdutils.Cmd, func(*cmdutils.Cmd) error)

// NewMockCmd returns a mock Cmd for a parent command (such as enable, get, utils) to be used for testing
// cli related features like flag and config file loading
func NewMockCmd(cmdFunc CommandFunction, parentCommand string, args ...string) *MockCmd {
	mockCmd := &MockCmd{}
	grouping := cmdutils.NewGrouping()
	parentCmd := cmdutils.NewVerbCmd(parentCommand, "", "")
	cmdutils.AddResourceCmd(grouping, parentCmd, func(cmd *cmdutils.Cmd) {
		noOpRunFunc := func(cmd *cmdutils.Cmd) error {
			mockCmd.Cmd = cmd
			return nil // no-op, to only test input aggregation & validation.
		}
		cmdFunc(cmd, noOpRunFunc)
	})
	parentCmd.SetArgs(args)
	mockCmd.parentCmd = parentCmd
	return mockCmd
}

type MockCmd struct {
	parentCmd *cobra.Command
	Cmd       *cmdutils.Cmd
}

func (c MockCmd) Execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOut(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}

// CreateConfigFile creates a temporary configuration file for testing by marshalling the given object in yaml. It
// returns the path to the file
func CreateConfigFile(cfg *api.ClusterConfig) string {
	contents, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("unable to marshal object: %s", err.Error())
	}
	file, err := ioutil.TempFile("", "enable-test-config-file")
	if err != nil {
		log.Fatalf("unable to create temp config file: %s", err.Error())
	}

	defer file.Close()

	_, err = file.Write(contents)
	if err != nil {
		log.Fatalf("unable to write to temp config file: %s", err.Error())
	}
	return file.Name()
}
