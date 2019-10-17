---
title: "Manage IAM users and roles"
weight: 100
url: usage/iam-identity-mappings
---

## Managing IAM users and roles

EKS clusters use IAM users and roles to control access to the cluster. The rules are implemented in a config map
called `aws-auth`. `eksctl` provides commands to read and edit this config map.

Get all identity mappings:

```bash
eksctl get iamidentitymapping --cluster my-cluster-1
```

Get all identity mappings matching an arn:

```bash
eksctl get iamidentitymapping --cluster my-cluster-1 --arn arn:aws:iam::123456:role/testing-role
```

Create an identity mapping:

```bash
 eksctl create iamidentitymapping --cluster  my-cluster-1 --arn arn:aws:iam::123456:role/testing --group system:masters --username admin
```

Delete a mapping:

```bash
eksctl delete iamidentitymapping --cluster  my-cluster-1 --arn arn:aws:iam::123456:role/testing
```

_Note_: this deletes a single mapping FIFO unless `--all`is given in which case it removes all matching. Will warn if
more mappings matching this role are found.
