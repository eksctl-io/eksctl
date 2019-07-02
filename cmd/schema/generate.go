package main

import (
	"github.com/alecthomas/jsonschema"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		panic("expected one argument with the output file")
	}
	outputFile := os.Args[1]

	var document strings.Builder
	document.WriteString(`---
title: Config file schema
weight: 200
url: usage/schema
---

`)
	document.WriteString("```yaml\n")

	schema := jsonschema.Reflect(&api.ClusterConfig{})
	yamlSchema, err := yaml.Marshal(schema.Definitions)
	if err != nil {
		panic(err)
	}
	document.Write(yamlSchema)
	document.WriteString("```\n")

	err = ioutil.WriteFile(outputFile, []byte(document.String()), 0755)

	if err != nil {
		panic(err)
	}

}
