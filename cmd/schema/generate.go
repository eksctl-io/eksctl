package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alecthomas/jsonschema"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io"
	v1alpha5 "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func main() {

	if len(os.Args) != 2 {
		panic("expected one argument with the output file")
	}
	outputFile := os.Args[1]

	var document strings.Builder
	schema := jsonschema.Reflect(&v1alpha5.ClusterConfig{})
	// We have to manually add examples here, because we can't tag `TypeMeta`
	// from the k8s package
	if kind, ok := schema.Definitions["ClusterConfig"].Properties.Get("kind"); ok {
		t := kind.(*jsonschema.Type)
		t.Examples = []interface{}{"ClusterConfig"}
	}
	if kind, ok := schema.Definitions["ClusterConfig"].Properties.Get("apiVersion"); ok {
		t := kind.(*jsonschema.Type)
		t.Examples = []interface{}{fmt.Sprintf("%s/%s", api.GroupName, v1alpha5.CurrentGroupVersion)}
	}
	jsonSchema, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err)
	}
	document.Write(jsonSchema)

	err = ioutil.WriteFile(outputFile, []byte(document.String()), 0755)

	if err != nil {
		panic(err)
	}

}
