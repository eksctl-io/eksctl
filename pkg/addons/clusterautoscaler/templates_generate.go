// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/weaveworks/eksctl/pkg/addons/clusterautoscaler"
)

func main() {
	err := vfsgen.Generate(clusterautoscaler.Templates, vfsgen.Options{
		PackageName:  "clusterautoscaler",
		VariableName: "Templates",
		BuildTags:    "!dev",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
