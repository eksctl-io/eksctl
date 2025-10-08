# Cluster Upgrade Policy

This document describes how to configure the upgrade policy for your EKS cluster using eksctl.

## Overview

The `upgradePolicy` field allows you to specify the support type for your EKS cluster. This determines the level of support AWS provides for your cluster version.

## Support Types

- **STANDARD**: The default support type that provides standard AWS support for the cluster
- **EXTENDED**: Provides extended support for older Kubernetes versions beyond the standard support period

## Configuration

You can specify the upgrade policy in your cluster configuration file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-cluster
  region: us-west-2

upgradePolicy:
  supportType: "EXTENDED"  # or "STANDARD"
```

## Command Line Usage

When creating a cluster with a specific upgrade policy:

```bash
eksctl create cluster --config-file=cluster-config.yaml
```

## Notes

- If no `upgradePolicy` is specified, AWS will use its default behavior
- The upgrade policy can only be set during cluster creation
- Extended support may incur additional costs - refer to AWS EKS pricing documentation
