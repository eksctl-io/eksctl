#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

echo "eksctl: starting bootstrap"
/var/lib/cloud/scripts/eksctl/bootstrap.linux.sh

echo "eksctl: restarting kubelet-eks"
snap restart kubelet-eks
