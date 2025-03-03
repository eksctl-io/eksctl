#!/usr/bin/env bash

set -o errexit
set -o nounset


PROJECT_ROOT=$(git rev-parse --show-toplevel)

# Grab code-generator pkg
CODEGEN_PKG=$(go list -m -f '{{.Dir}}' 'k8s.io/code-generator')
echo ">> Using ${CODEGEN_PKG}"


source "${CODEGEN_PKG}/kube_codegen.sh"

kube::codegen::gen_helpers \
    --boilerplate <(printf "/*\n%s\n*/\n" "$(cat LICENSE)") \
    "${PROJECT_ROOT}/pkg/apis"
