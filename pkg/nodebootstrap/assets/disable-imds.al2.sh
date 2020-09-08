#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

yum install -y iptables-services
iptables --insert FORWARD 1 --in-interface eni+ --destination 169.254.169.254/32 --jump DROP
iptables-save | tee /etc/sysconfig/iptables
systemctl enable --now iptables
