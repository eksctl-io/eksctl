---
title: "Enabling CloudWatch logging"
weight: 120
url: usage/cloudwatch-cluster-logging
---

[eksdocs]: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html

## Enabling CloudWatch logging

[CloudWatch logging][eksdocs] for EKS control plane is not enabled by default due to data
ingestion and storage costs.

To enable control plane logging when cluster is created, you will need to define **`cloudWatch.clusterLogging.enableTypes`** setting in your `ClusterConfig` (see below for examples).

So if you have a config file with correct **`cloudWatch.clusterLogging.enableTypes`**
setting, you can create a cluster with `eksctl create cluster --config-file=<path>`.

If you have created a cluster already, you can use `eksctl utils enable-logging`.

> **NOTE**: this command runs in plan mode by default, you will need to specify `--approve` flag to
> apply the changes to your cluster.

If you are using a config file, run:
```
eksctl utils enable-logging --config-file=<path>
```

Alternatively, you can use CLI flags.

To enable all types of logs, run:
```
eksctl utils enable-logging
```

To enable `audit` logs, run:
```
eksctl utils enable-logging --audit
```

To enable all but `controllerManager` logs, run:
```
eksctl utils enable-logging --controllerManager=false
```

To disable all types of logs, run:
```
eksctl utils enable-logging --all=false
```

### `ClusterConfig` Examples

There 5 types of logs that you may wish to enable (see [EKS documentation](eksdocs) for more details):

- `api`
- `audit`
- `authenticator`
- `controllerManager`
- `scheduler`

You can enable all types with `"*"` or `"all"`, i.e.:

```YAML
cloudWatch:
  clusterLogging:
    enableTypes: ["*"]
```

To disable all types, use `[]` or remove `cloudWatch` section completely.

You can enable a subset of types by listing the types you want to enable:

```YAML
cloudWatch:
  clusterLogging:
    enableTypes:
      - "audit"
      - "authenticator"
```

Full example:
```YAML
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-11
  region: eu-west-2

nodeGroups:
  - name: ng-1
    instanceType: m5.large
    desiredCapacity: 1

cloudWatch:
  clusterLogging:
    enableTypes: ["audit", "authenticator"]
```