# Support Node Repair Configuration for EKS Managed Nodegroups

EKS Managed Nodegroups supports Node Repair, where the health of managed nodes are monitored,
and unhealthy worker nodes are replaced or rebooted in response. eksctl now provides comprehensive
configuration options for fine-grained control over node repair behavior.

## Basic Node Repair Configuration

### Using CLI flags

To create a cluster with a managed nodegroup using basic node repair:

```shell
$ eksctl create cluster --enable-node-repair
```

To create a managed nodegroup with node repair on an existing cluster:

```shell
$ eksctl create nodegroup --cluster=<clusterName> --enable-node-repair
```

### Using configuration files

```yaml
# basic-node-repair.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: basic-node-repair-cluster
  region: us-west-2

managedNodeGroups:
- name: ng-1
  nodeRepairConfig:
    enabled: true
```

```shell
$ eksctl create cluster -f basic-node-repair.yaml
```

## Enhanced Node Repair Configuration

### Threshold Configuration

You can configure when node repair actions will stop using either percentage or count-based thresholds. **Note: You cannot use both percentage and count thresholds at the same time.**

#### CLI flags for thresholds

```shell
# Percentage-based threshold - repair stops when 20% of nodes are unhealthy
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-unhealthy-percentage=20

# Count-based threshold - repair stops when 5 nodes are unhealthy
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-unhealthy-count=5
```

#### Configuration file for thresholds

```yaml
managedNodeGroups:
- name: threshold-ng
  nodeRepairConfig:
    enabled: true
    # Stop repair actions when 20% of nodes are unhealthy
    maxUnhealthyNodeThresholdPercentage: 20
    # Alternative: stop repair actions when 3 nodes are unhealthy
    # maxUnhealthyNodeThresholdCount: 3
    # Note: Cannot use both percentage and count thresholds simultaneously
```

### Parallel Repair Limits

Control the maximum number of nodes that can be repaired concurrently or in parallel. This gives you finer-grained control over the pace of node replacements. **Note: You cannot use both percentage and count limits at the same time.**

#### CLI flags for parallel limits

```shell
# Percentage-based parallel limits - repair at most 15% of unhealthy nodes in parallel
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-parallel-percentage=15

# Count-based parallel limits - repair at most 2 unhealthy nodes in parallel
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-parallel-count=2
```

#### Configuration file for parallel limits

```yaml
managedNodeGroups:
- name: parallel-ng
  nodeRepairConfig:
    enabled: true
    # Repair at most 15% of unhealthy nodes in parallel
    maxParallelNodesRepairedPercentage: 15
    # Alternative: repair at most 2 unhealthy nodes in parallel
    # maxParallelNodesRepairedCount: 2
    # Note: Cannot use both percentage and count limits simultaneously
```

### Custom Repair Overrides

Specify granular overrides for specific repair actions. These overrides control the repair action and the repair delay time before a node is considered eligible for repair. **If you use this, you must specify all the values for each override.**

```yaml
managedNodeGroups:
- name: custom-repair-ng
  instanceType: g4dn.xlarge  # GPU instances
  nodeRepairConfig:
    enabled: true
    maxUnhealthyNodeThresholdPercentage: 25
    maxParallelNodesRepairedCount: 1
    nodeRepairConfigOverrides:
      # Handle GPU-related failures with immediate termination
      - nodeMonitoringCondition: "AcceleratedInstanceNotReady"
        nodeUnhealthyReason: "NvidiaXID13Error"
        minRepairWaitTimeMins: 10
        repairAction: "Terminate"
      # Handle network issues with restart after waiting
      - nodeMonitoringCondition: "NetworkNotReady"
        nodeUnhealthyReason: "InterfaceNotUp"
        minRepairWaitTimeMins: 20
        repairAction: "Restart"
```

## Complete Configuration Examples

For a comprehensive example with all configuration options, see [examples/44-node-repair.yaml](https://github.com/eksctl-io/eksctl/blob/main/examples/44-node-repair.yaml).

### Example 1: Basic repair with percentage thresholds

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: basic-repair-cluster
  region: us-west-2

managedNodeGroups:
- name: basic-ng
  instanceType: m5.large
  desiredCapacity: 3
  nodeRepairConfig:
    enabled: true
    maxUnhealthyNodeThresholdPercentage: 20
    maxParallelNodesRepairedPercentage: 15
```

### Example 2: Conservative repair for critical workloads

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: critical-workload-cluster
  region: us-west-2

managedNodeGroups:
- name: critical-ng
  instanceType: c5.2xlarge
  desiredCapacity: 6
  nodeRepairConfig:
    enabled: true
    # Very conservative settings
    maxUnhealthyNodeThresholdPercentage: 10
    maxParallelNodesRepairedCount: 1
    nodeRepairConfigOverrides:
      # Wait longer before taking action on critical workloads
      - nodeMonitoringCondition: "NetworkNotReady"
        nodeUnhealthyReason: "InterfaceNotUp"
        minRepairWaitTimeMins: 45
        repairAction: "Restart"
```

### Example 3: GPU workload with specialized repair

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: gpu-workload-cluster
  region: us-west-2

managedNodeGroups:
- name: gpu-ng
  instanceType: g4dn.xlarge
  desiredCapacity: 4
  nodeRepairConfig:
    enabled: true
    maxUnhealthyNodeThresholdPercentage: 25
    maxParallelNodesRepairedCount: 1
    nodeRepairConfigOverrides:
      # GPU failures require immediate termination
      - nodeMonitoringCondition: "AcceleratedInstanceNotReady"
        nodeUnhealthyReason: "NvidiaXID13Error"
        minRepairWaitTimeMins: 5
        repairAction: "Terminate"
```

## CLI Reference

### Node Repair Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--enable-node-repair` | Enable automatic node repair | `--enable-node-repair` |
| `--node-repair-max-unhealthy-percentage` | Maximum percentage of unhealthy nodes before repair | `--node-repair-max-unhealthy-percentage=20` |
| `--node-repair-max-unhealthy-count` | Maximum count of unhealthy nodes before repair | `--node-repair-max-unhealthy-count=5` |
| `--node-repair-max-parallel-percentage` | Maximum percentage of nodes to repair in parallel | `--node-repair-max-parallel-percentage=15` |
| `--node-repair-max-parallel-count` | Maximum count of nodes to repair in parallel | `--node-repair-max-parallel-count=2` |

**Note:** Node repair config overrides are only supported through YAML configuration files due to their complexity.

## Configuration Reference

### nodeRepairConfig

| Field | Type | Description | Constraints | Example |
|-------|------|-------------|-------------|---------|
| `enabled` | boolean | Enable/disable node repair | - | `true` |
| `maxUnhealthyNodeThresholdPercentage` | integer | Percentage threshold of unhealthy nodes, above which node auto repair actions will stop | Cannot be used with `maxUnhealthyNodeThresholdCount` | `20` |
| `maxUnhealthyNodeThresholdCount` | integer | Count threshold of unhealthy nodes, above which node auto repair actions will stop | Cannot be used with `maxUnhealthyNodeThresholdPercentage` | `5` |
| `maxParallelNodesRepairedPercentage` | integer | Maximum percentage of unhealthy nodes that can be repaired concurrently or in parallel | Cannot be used with `maxParallelNodesRepairedCount` | `15` |
| `maxParallelNodesRepairedCount` | integer | Maximum count of unhealthy nodes that can be repaired concurrently or in parallel | Cannot be used with `maxParallelNodesRepairedPercentage` | `2` |
| `nodeRepairConfigOverrides` | array | Granular overrides for specific repair actions controlling repair action and delay time | All values must be specified for each override | See examples above |

### nodeRepairConfigOverrides

| Field | Type | Description | Valid Values |
|-------|------|-------------|--------------|
| `nodeMonitoringCondition` | string | Unhealthy condition reported by the node monitoring agent that this override applies to | `"AcceleratedInstanceNotReady"`, `"NetworkNotReady"` |
| `nodeUnhealthyReason` | string | Reason reported by the node monitoring agent that this override applies to | `"NvidiaXID13Error"`, `"InterfaceNotUp"` |
| `minRepairWaitTimeMins` | integer | Minimum time in minutes to wait before attempting to repair a node with the specified condition and reason | Any positive integer |
| `repairAction` | string | Repair action to take for nodes when all of the specified conditions are met | `"Terminate"`, `"Restart"`, `"NoAction"` |

## Further Information

- [EKS Managed Nodegroup Node Health][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/node-health.html
