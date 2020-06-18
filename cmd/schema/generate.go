package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	schemapkg "github.com/weaveworks/eksctl/pkg/schema"
)

func main() {
	if len(os.Args) != 2 {
		panic("expected one argument with the output file")
	}
	outputFile := os.Args[1]

	input := filepath.Join("../../../..", "pkg", "apis", "eksctl.io")
	schema, err := schemapkg.GenerateSchema(input, "v1alpha5", "ClusterConfig", false)
	if err != nil {
		panic(err)
	}

	// We add some examples and blacklist some descriptions
	if t, ok := schema.Definitions["ClusterConfig"].Properties["kind"]; ok {
		t.Examples = []string{"ClusterConfig"}
		t.Description = ""
		t.HTMLDescription = ""
	}
	if t, ok := schema.Definitions["ClusterConfig"].Properties["apiVersion"]; ok {
		t.Examples = []string{fmt.Sprintf("%s/%s", api.GroupName, v1alpha5.CurrentGroupVersion)}
		t.Description = ""
		t.HTMLDescription = ""
	}

	bytes, err := schemapkg.ToJSON(schema)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(outputFile, bytes, 0755)
	if err != nil {
		panic(err)
	}
}
