#!/usr/bin/env bash

set -x -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
AWSSERV=${DIR}/vendor/github.com/aws/aws-sdk-go/service
MOCKERY=$GOPATH/bin/mockery
MOCKOUT=${DIR}/pkg/testutils/mocks

echo "Generating CloudFormation mock"
${MOCKERY} -dir=${AWSSERV}/cloudformation/cloudformationiface -name=CloudFormationAPI -output=${MOCKOUT}

echo "Generating EC2 mock"
${MOCKERY} -dir=${AWSSERV}/ec2/ec2iface -name=EC2API -output=${MOCKOUT}

echo "Generating EKS mock"
${MOCKERY} -dir=${AWSSERV}/eks/eksiface -name=EKSAPI -output=${MOCKOUT}

echo "Generating STS mock"
${MOCKERY} -dir=${AWSSERV}/sts/stsiface -name=STSAPI -output=${MOCKOUT}
