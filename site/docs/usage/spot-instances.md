## Spot instances

`eksctl` has support for spot instances through the MixedInstancesPolicy for Auto Scaling Groups.

Here is an example of a nodegroup that uses 50% spot instances and 50% on demand instances:

```yaml
nodeGroups:
  - name: ng-1
    minSize: 2
    maxSize: 5
    instancesDistribution:
      maxPrice: 0.017
      instanceTypes: ["t3.small", "t3.medium"] # At least two instance types should be specified
      onDemandBaseCapacity: 0
      onDemandPercentageAboveBaseCapacity: 50
      spotInstancePools: 2
```

Note that the `nodeGroups.X.instanceType` field shouldn't be set when using the `instancesDistribution` field.

This example uses GPU instances:

```yaml
nodeGroups:
  - name: ng-gpu
    instanceType: mixed
    desiredCapacity: 1
    instancesDistribution:
      instanceTypes:
        - p2.xlarge
        - p2.8xlarge
        - p2.16xlarge
      maxPrice: 0.50
```

Here is a minimal example:

```yaml
nodeGroups:
  - name: ng-1
    instancesDistribution:
      instanceTypes: ["t3.small", "t3.medium"] # At least two instance types should be specified
```

### Parameters in instancesDistribution

|                                     | type        | required | default value   |
| ----------------------------------- | ----------- | -------- | --------------- |
| instanceTypes                       | []string    | required | -               |
| maxPrice                            | float       | optional | on demand price |
| onDemandBaseCapacity                | int         | optional | 0               |
| onDemandPercentageAboveBaseCapacity | int [1-100] | optional | 100             |
| spotInstancePools                   | int [1-20]  | optional | 2               |
