# Enhanced Node Repair Configuration for EKS Managed Nodegroups

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

You can configure when node repair is triggered using either percentage or count-based thresholds:

#### CLI flags for thresholds

```shell
# Percentage-based thresholds
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-unhealthy-percentage=20

# Count-based thresholds  
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-unhealthy-count=5
```

#### Configuration file for thresholds

```yaml
managedNodeGroups:
- name: threshold-ng
  nodeRepairConfig:
    enabled: true
    # Trigger repair when 20% of nodes are unhealthy
    maxUnhealthyNodeThresholdPercentage: 20
    # Alternative: trigger repair when 3 nodes are unhealthy
    # maxUnhealthyNodeThresholdCount: 3
```

### Parallel Repair Limits

Control how many nodes can be repaired simultaneously:

#### CLI flags for parallel limits

```shell
# Percentage-based parallel limits
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-parallel-percentage=15

# Count-based parallel limits
$ eksctl create cluster --enable-node-repair \
  --node-repair-max-parallel-count=2
```

#### Configuration file for parallel limits

```yaml
managedNodeGroups:
- name: parallel-ng
  nodeRepairConfig:
    enabled: true
    # Repair at most 15% of nodes in parallel
    maxParallelNodesRepairedPercentage: 15
    # Alternative: repair at most 2 nodes in parallel
    # maxParallelNodesRepairedCount: 2
```

### Custom Repair Overrides

Define specialized repair behavior for specific failure scenarios:

```yaml
managedNodeGroups:
- name: custom-repair-ng
  instanceType: g4dn.xlarge  # GPU instances
  nodeRepairConfig:
    enabled: true
    maxUnhealthyNodeThresholdPercentage: 25
    maxParallelNodesRepairedCount: 1
    nodeRepairConfigOverrides:
      # Handle GPU-related failures
      - nodeMonitoringCondition: "AcceleratedInstanceNotReady"
        nodeUnhealthyReason: "NvidiaXID13Error"
        minRepairWaitTimeMins: 10
        repairAction: "Terminate"
      # Handle network issues
      - nodeMonitoringCondition: "NetworkNotReady"
        nodeUnhealthyReason: "InterfaceNotUp"
        minRepairWaitTimeMins: 20
        repairAction: "Restart"
```

## Complete Configuration Examples

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

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `enabled` | boolean | Enable/disable node repair | `true` |
| `maxUnhealthyNodeThresholdPercentage` | integer | Percentage threshold for unhealthy nodes | `20` |
| `maxUnhealthyNodeThresholdCount` | integer | Count threshold for unhealthy nodes | `5` |
| `maxParallelNodesRepairedPercentage` | integer | Percentage limit for parallel repairs | `15` |
| `maxParallelNodesRepairedCount` | integer | Count limit for parallel repairs | `2` |
| `nodeRepairConfigOverrides` | array | Custom repair behavior overrides | See examples above |

### nodeRepairConfigOverrides

| Field | Type | Description | Valid Values |
|-------|------|-------------|--------------|
| `nodeMonitoringCondition` | string | Monitoring condition | `"AcceleratedInstanceNotReady"`, `"NetworkNotReady"` |
| `nodeUnhealthyReason` | string | Reason for node being unhealthy | `"NvidiaXID13Error"`, `"InterfaceNotUp"` |
| `minRepairWaitTimeMins` | integer | Minimum wait time before repair (minutes) | Any positive integer |
| `repairAction` | string | Action to take for repair | `"Terminate"`, `"Restart"`, `"NoAction"` |

## Best Practices

### Choosing Thresholds

- **Small nodegroups (< 10 nodes)**: Use count-based thresholds for precise control
- **Large nodegroups (≥ 10 nodes)**: Use percentage-based thresholds for scalability
- **Critical workloads**: Use conservative thresholds (10-15%)
- **Development environments**: Use higher thresholds (20-30%)

### Parallel Repair Limits

- **High availability requirements**: Limit to 1-2 nodes or 10-15%
- **Batch workloads**: Allow higher parallel repairs (20-25%)
- **GPU workloads**: Limit to 1 node at a time due to cost and setup time

### Custom Overrides

- **GPU instances**: Use immediate termination for hardware failures
- **Network issues**: Try restart first, then terminate
- **Critical workloads**: Increase wait times to avoid unnecessary disruptions

## Troubleshooting

### Common Issues

1. **Configuration validation errors**: Ensure parameter values are within valid ranges
2. **Conflicting thresholds**: Don't specify both percentage and count for the same parameter
3. **Invalid override values**: Check that monitoring conditions, reasons, and actions are valid

### Monitoring Node Repair

Enable CloudWatch logging to monitor node repair activities:

```yaml
cloudWatch:
  clusterLogging:
    enableTypes: ["api", "audit", "authenticator", "controllerManager", "scheduler"]
```

## Further Information

- [EKS Managed Nodegroup Node Health][eks-user-guide]
- [EKS Node Repair Configuration][eks-node-repair]
- [eksctl Managed Nodegroups][eksctl-managed-nodegroups]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/node-health.html
[eks-node-repair]: https://docs.aws.amazon.com/eks/latest/userguide/node-repair.html
[eksctl-managed-nodegroups]: https://eksctl.io/usage/managing-nodegroups/
