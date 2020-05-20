package utils

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func schemaCmd(cmd *cmdutils.Cmd) {
	cmd.SetDescription("schema", "Output the ClusterConfig JSON Schema", "")

	cmd.CobraCommand.Run = func(_ *cobra.Command, args []string) {
		doSchemaCmd(cmd)
	}
}

func doSchemaCmd(cmd *cmdutils.Cmd) {
	schema, err := api.Asset("schema.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s", schema)
}
