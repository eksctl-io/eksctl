#!/usr/bin/env bash

set -o errexit
set -o nounset


PROJECT_ROOT=$(git rev-parse --show-toplevel)

# Grab code-generator pkg
CODEGEN_PKG=$(go list -m -f '{{.Dir}}' 'k8s.io/code-generator')
echo ">> Using ${CODEGEN_PKG}"


source "${CODEGEN_PKG}/kube_codegen.sh"

# Create temporary boilerplate file to avoid process substitution issues on macOS
TEMP_BOILERPLATE=$(mktemp)
printf "/*\n%s\n*/\n" "$(cat LICENSE)" > "${TEMP_BOILERPLATE}"

kube::codegen::gen_helpers \
    --boilerplate "${TEMP_BOILERPLATE}" \
    "${PROJECT_ROOT}/pkg/apis"

# Clean up temporary file
rm "${TEMP_BOILERPLATE}"
