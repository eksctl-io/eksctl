#!/bin/sh -ex

docker run -i -v "${PWD}":/src --workdir /src "${EKSCTL_DEPENDENCIES_IMAGE}" sh -s << EOF
go list -tags 'integration tools' -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n" }}{{"\n"}}{{join .XTestImports "\n" }}' ./...  |
  sort | uniq | grep -v eksctl | xargs go list -f '{{ if not .Standard }}{{.ImportPath}}{{end}}'
EOF
