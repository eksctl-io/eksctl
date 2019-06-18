package main

import (
	"github.com/alecthomas/jsonschema"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"io/ioutil"
	"sigs.k8s.io/yaml"
	"strings"
)

const outputFile = "site/content/usage/20-schema.md"

func main() {

	var document strings.Builder
	document.WriteString(`---
title: Config file schema
weight: 200
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
