#!/bin/sh -eu

# Simple wrapper for authenticator, so we can read AWS_EKS_CLUSTER_NAME from `metadata.env`
# instead of having to update `kubeconfig.yaml` (as `kubectl config` doesn't support plugin fields)

. /etc/eksctl/metadata.env

heptio-authenticator-aws token -i "${AWS_EKS_CLUSTER_NAME}"