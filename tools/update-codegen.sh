#!/usr/bin/env bash

set -o errexit
set -o nounset
#set -o pipefail


SCRIPT_ROOT=$(git rev-parse --show-toplevel)

# Grab code-generator version from go.sum.
CODEGEN_VERSION=$(grep 'k8s.io/code-generator' go.sum | awk '{print $2}' | head -1)
CODEGEN_PKG=$(echo `go env GOPATH`"/pkg/mod/k8s.io/code-generator@${CODEGEN_VERSION}")

echo ">> Using ${CODEGEN_PKG}"

# code-generator does work with go.mod but makes assumptions about
# the project living in `$GOPATH/src`. To work around this and support
# any location; create a temporary directory, use this as an output
# base, and copy everything back once generated.
TEMP_DIR=$(mktemp -d)
cleanup() {
    echo ">> Removing ${TEMP_DIR}"
    rm -rf ${TEMP_DIR}
}
trap "cleanup" EXIT SIGINT

echo ">> Temporary output directory ${TEMP_DIR}"

# Ensure we can execute.
chmod +x ${CODEGEN_PKG}/generate-groups.sh

GOPATH=$(go env GOPATH) ${CODEGEN_PKG}/generate-groups.sh deepcopy,defaulter \
    _ github.com/weaveworks/eksctl/pkg/apis \
    eksctl.io:v1alpha5 \
    --go-header-file .license-header \
    --output-base "${TEMP_DIR}"

# Copy everything back.
cp -r "${TEMP_DIR}/github.com/weaveworks/eksctl/." "${SCRIPT_ROOT}/"
