#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

dnf install -y wget
wget -q --timeout=20 https://s3-us-west-2.amazonaws.com/aws-efa-installer/aws-efa-installer-latest.tar.gz -O /tmp/aws-efa-installer.tar.gz
tar -xf /tmp/aws-efa-installer.tar.gz -C /tmp
rm -rf /tmp/aws-efa-installer.tar.gz
cd /tmp/aws-efa-installer
./efa_installer.sh -y -g
/opt/amazon/efa/bin/fi_info -p efa
