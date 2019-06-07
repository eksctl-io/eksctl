---
title: Cluster with only spot instances
weight: 10
---

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
    name: spot-cluster
    region: eu-north-1

nodeGroups:
    - name: ng-1
      minSize: 2
      maxSize: 5
      instancesDistribution:
        maxPrice: 0.018
        instanceTypes: ["m5.xlarge", "m5.large"] # At least two instance types should be specified
        onDemandBaseCapacity: 0
        onDemandPercentageAboveBaseCapacity: 0
        spotInstancePools: 2
```
