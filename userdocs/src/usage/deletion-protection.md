# Cluster Deletion Protection

This document describes how to configure deletion protection for your EKS cluster using eksctl.

## Overview

The `deletionProtection` field allows you to enable deletion protection for your EKS cluster. This prevents accidental cluster deletion.

## Configuration

You can specify deletion protection in your cluster configuration file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-cluster
  region: us-west-2

deletionProtection: true
```

## Command Line Usage

When creating a cluster with deletion protection:

```bash
eksctl create cluster --config-file=cluster-config.yaml
```

To update deletion protection on an existing cluster:

```bash
# Enable deletion protection
eksctl utils deletion-protection --name=my-cluster --enabled=true --approve

# Disable deletion protection
eksctl utils deletion-protection --name=my-cluster --enabled=false --approve
```

## Notes

- If no `deletionProtection` is specified, AWS will use its default behavior (false)
- Deletion protection can be set during cluster creation and updated later
- When enabled, you must disable deletion protection before you can delete the cluster
