#!/bin/sh -ex

BUILDID="$(docker images --filter "label=eksctl.builder=true" --format '{{.ID}}')"
mkdir -p ${PWD}/test-results/ginkgo
docker run -i --rm -v ${PWD}/test-results/ginkgo:/mnt/test-results ${BUILDID}  sh -s <<EOF
cp /src/test-results/ginkgo/*.xml /mnt/test-results
EOF

