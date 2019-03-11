#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

skaffold build --profile=local
make integration-test-container-pre-built EKSCTL_IMAGE=weaveworks/eksctl:local TEST_V=1
