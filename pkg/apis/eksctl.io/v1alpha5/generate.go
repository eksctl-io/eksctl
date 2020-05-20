package v1alpha5

//go:generate go run ../../../../cmd/schema/generate.go assets/schema.json
//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o schema.go assets
