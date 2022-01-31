#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

yum install -y amazon-ssm-agent
systemctl enable amazon-ssm-agent
systemctl start amazon-ssm-agent
