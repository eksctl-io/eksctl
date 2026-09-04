# Control Plane Component Configuration

EKS allows you to customize the configuration of managed control plane components. Instead of accepting
the defaults EKS applies, you can tune selected settings on the `kube-apiserver`, `kube-scheduler`, and
`kube-controller-manager` to suit the workloads running on your cluster.

Each of the three components is configured through its own top-level field, and every field is optional.
If you omit a field, EKS applies its default.

## Creating a cluster with control plane component configuration

```yaml
# control-plane-component-config.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-cluster
  region: us-west-2

kubeAPIServerConfig:
  eventTTL: 30m
  serviceNodePortRange:
    minPort: 30000
    maxPort: 32767

kubeSchedulerConfig:
  nodeResourcesFit:
    scoringStrategy:
      type: MostAllocated
      resources:
        - name: cpu
          weight: 2
        - name: memory
          weight: 1

kubeControllerManagerConfig:
  horizontalPodAutoscalerControllerConfig:
    horizontalPodAutoscalerSyncPeriod: 15s
  podGCControllerConfig:
    terminatedPodGCThreshold: 12000
```

```shell
$ eksctl create cluster -f control-plane-component-config.yaml
```

## Updating an existing cluster

To change the configuration of cluster that already exists, edit the same fields in your config file and run

```shell
$ eksctl utils update-control-plane-component-config -f control-plane-component-config.yaml
```

A config file is required, and at least one of `kubeAPIServerConfig`, `kubeSchedulerConfig` or
`kubeControllerManagerConfig` must be set. Only the components present in the file are updated; any component
you omit is left unchanged.

## kube-apiserver

```yaml
kubeAPIServerConfig:
  eventTTL: 30m
  serviceNodePortRange:
    minPort: 30000
    maxPort: 32767
```

- `eventTTL` — how long Kubernetes events are retained, as a duration string such as `30m` or `1h`.
- `serviceNodePortRange` — the port range from which `NodePort` services are allocated. Both `minPort`
  and `maxPort` are port numbers.

## kube-scheduler

```yaml
kubeSchedulerConfig:
  nodeResourcesFit:
    scoringStrategy:
      type: MostAllocated
      resources:
        - name: cpu
          weight: 2
        - name: memory
          weight: 1
```

Configures the `NodeResourcesFit` scheduler plugin, which scores nodes by how their resource usage would
look after a pod is placed on them.

- `scoringStrategy.type` — either `LeastAllocated` (spread pods across nodes) or `MostAllocated` (pack pods
  onto fewer nodes).
- `scoringStrategy.resources` — the resources considered when scoring, each with a relative `weight`. In the
  example above `cpu` counts twice as much as `memory`.

## kube-controller-manager

```yaml
kubeControllerManagerConfig:
  horizontalPodAutoscalerControllerConfig:
    horizontalPodAutoscalerSyncPeriod: 15s
  podGCControllerConfig:
    terminatedPodGCThreshold: 12000
```

- `horizontalPodAutoscalerSyncPeriod` — the interval between horizontal pod autoscaler syncs, as a duration
  string such as `15s`. Shorter intervals make autoscaling more responsive at the cost of more frequent
  metric queries. Only configurable on provisioned control plane (PCP) scaling tiers (for example `tier-xl`);
  on the Standard tier it is fixed at the default (15s).
- `terminatedPodGCThreshold` — the number of terminated pods that can accumulate before the pod garbage
  collector starts deleting them. Lower values reclaim resources sooner; higher values keep more terminated
  pods around for inspection. Only configurable on provisioned control plane (PCP) scaling tiers (for example
  `tier-xl`); on the Standard tier it is fixed at the default (12500).

## Notes

- The accepted values depend on the scaling tier of the cluster and on its Kubernetes version, so a value
  accepted for one cluster may be rejected for another. Values are validated by EKS rather than by `eksctl`,
  on both create and update. Consult the EKS documentation for the limits that apply to your cluster.
