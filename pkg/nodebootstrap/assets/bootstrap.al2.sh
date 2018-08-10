#!/bin/sh -eu

echo "NODE_IP=$(hostname -i)" > /etc/eksctl/kubelet.local.env

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
