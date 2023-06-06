# CloudWatch logging

## Enabling CloudWatch logging

[CloudWatch logging][eksdocs] for EKS control plane is not enabled by default due to data
ingestion and storage costs.

To enable control plane logging when cluster is created, you will need to define **`cloudWatch.clusterLogging.enableTypes`** setting in your `ClusterConfig` (see below for examples).

So if you have a config file with correct **`cloudWatch.clusterLogging.enableTypes`**
setting, you can create a cluster with `eksctl create cluster --config-file=<path>`.

If you have created a cluster already, you can use `eksctl utils update-cluster-logging`.

???+ note
    this command runs in plan mode by default, you will need to specify `--approve` flag to
    apply the changes to your cluster.

If you are using a config file, run:

```
eksctl utils update-cluster-logging --config-file=<path>
```

Alternatively, you can use CLI flags.

To enable all types of logs, run:

```
eksctl utils update-cluster-logging --enable-types all
```

To enable `audit` logs, run:
```
eksctl utils update-cluster-logging --enable-types audit
```

To enable all but `controllerManager` logs, run:
```
eksctl utils update-cluster-logging --enable-types=all --disable-types=controllerManager
```

If the `api` and `scheduler` log types were already enabled, to disable `scheduler` and enable `controllerManager` at
the same time, run:

```
eksctl utils update-cluster-logging --enable-types=controllerManager --disable-types=scheduler
```

This will leave `api` and `controllerManager` as the only log types enabled.

To disable all types of logs, run:
```
eksctl utils update-cluster-logging --disable-types all
```

## `ClusterConfig` Examples

In an EKS cluster, the `enableTypes` field under `clusterLogging` can take a list of possible values to enable the different types of logs for the control plane components.

The following are the possible values:

- `api`: Enables the Kubernetes API server logs.
- `audit`: Enables the Kubernetes audit logs.
- `authenticator`: Enables the authenticator logs.
- `controllerManager`: Enables the Kubernetes controller manager logs.
- `scheduler`: Enables the Kubernetes scheduler logs.

To learn more, see [EKS documentation][eksdocs]. 

### Disable all logs
To disable all types, use `[]` or remove the `cloudWatch` section completely.

### Enable all logs
You can enable all types with `"*"` or `"all"`. For example:

```yaml
cloudWatch:
  clusterLogging:
    enableTypes: ["*"]
```

### Enable one or more logs
You can enable a subset of types by listing the types you want to enable. For example:

```yaml
cloudWatch:
  clusterLogging:
    enableTypes:
      - "audit"
      - "authenticator"
```

### Log retention period
By default, logs are stored in CloudWatch Logs, indefinitely. You can specify the number of days for which the control plane logs should be retained in CloudWatch Logs. The following example retains logs for 7 days:

```yaml
cloudWatch:
  clusterLogging:
    logRetentionInDays: 7
```

### Complete example

```yaml
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
    logRetentionInDays: 7
```

[eksdocs]: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
