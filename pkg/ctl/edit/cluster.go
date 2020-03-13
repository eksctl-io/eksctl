package edit

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func editClusterCmd(cmd *cmdutils.Cmd) {
	cmd.SetDescription("cluster", "Edit cluster", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doEditClusterCmd(cmd)
	}
}

// doEditClusterCmd edits a cluster by calling eksctl get, open data as yml in editor then update the cluster
func doEditClusterCmd(cmd *cmdutils.Cmd) error {
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		editor := GetEditorEnvironmentVariable("EDITOR")

		resourceData, err := exec.Command("eksctl", "get", args[0]).Output()
		if err != nil {
			log.Fatal(err)
		}

		tmpfile, err := ioutil.TempFile("", "resource-data.*.yml")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(resourceData); err != nil {
			tmpfile.Close()
			log.Fatal(err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}

		openEditorCmd := exec.Command(editor, tmpfile.Name())
		openEditorCmd.Stdin = os.Stdin
		openEditorCmd.Stdout = os.Stdout
		openEditorCmd.Stderr = os.Stderr

		if err := openEditorCmd.Run(); err != nil {
			log.Fatal(err)
		}

		exec.Command("eksctl", "edit", args[0], "-f", tmpfile.Name())

		return nil
	}
	return nil
}

// GetEditorEnvironmentVariable gets $EDITOR environment variable to use in opening file
func GetEditorEnvironmentVariable(key string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return "vim"
	}

	return value
}
