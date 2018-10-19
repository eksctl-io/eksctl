#!/bin/sh -ex

if [ -n "${TEST_OUTPUT}" ]; then

BUILDID="$(docker images --filter "label=eksctl.builder=true" --format '{{.ID}}')"
mkdir -p ${PWD}/test-results/ginkgo
docker run -i --rm -v ${PWD}/test-results/ginkgo:/mnt/test-results  ${BUILDID}  sh -s <<EOF
cp /go/src/github.com/weaveworks/eksctl/test-results/ginkgo/*.xml /mnt/test-results
EOF

else
    echo "not retrieving test results"
fi
