#!/bin/sh -eu

echo "NODE_IP=$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)" > /etc/eksctl/kubelet.local.env

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
