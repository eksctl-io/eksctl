# Creating a Nodegroup

To create a new Ocean nodegroup, run:

```bash
eksctl create nodegroup \
  --cluster example \
  --nodegroup-name nodegroup-example \
  --spot-ocean
```

The command will create a nodegroup "nodegroup-example" managed by Spot Ocean for the Ocean cluster "example".

To have more control over the nodegroups configuration, like creating multiple nodegroups, a configuration file can be used.

The following configuration file example enables us to create two nodegroups managed by Ocean.
```yaml
# cluster.yaml
# A cluster with two Ocean nodegroups.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  # existing ocean cluster
  name: cluster-name
  region: us-west-2

nodeGroups:
- name: ocean-ng-1
  #[... nodegroup standard fields; ssh, tags, etc.]

  # Enable Ocean integration and use all defaults.
  spotOcean: {}

- name: ocean-ng-2
  #[... nodegroup standard fields; ssh, tags, etc.]

  # Enable Ocean integration with custom configuration.
  spotOcean:
    strategy:
      # Percentage of Spot instances that would spin up
      # from the desired capacity.
      spotPercentage: 100

      # Allow Ocean to utilize any available reserved
      # instances first before purchasing Spot instances.
      utilizeReservedInstances: true

      # Launch On-Demand instances in case of no Spot
      # instances available.
      fallbackToOnDemand: true

    autoScaler:
      # Spare resource capacity management enabling fast
      # assignment of Pods without waiting for new resources
      # to launch.
      headrooms:

        # Number of CPUs to allocate. CPUs are denoted
        # in millicores, where 1000 millicores = 1 vCPU.
      - cpuPerUnit: 2

        # Number of GPUs to allocate.
        gpuPerUnit: 0

        # Amount of memory (MB) to allocate.
        memoryPerUnit: 64

        # Number of units to retain as headroom, where
        # each unit has the defined CPU and memory.
        numOfUnits: 1

    compute:
      instanceTypes:
        # Instance types allowed in the Ocean cluster.
        # Cannot be configured if the blacklist is configured.
        whitelist: # OR blacklist
        - t2.large
        - c5.large
```

To create a new nodgroup managed by Spot Ocean for the cluster with the definition file run:

```bash
eksctl create nodegroup -f cluster.yaml
```


## AWS Nodegroups Immutability
By design, AWS nodegroups are immutable. This means that if you need to change something like the AMI or the instance type of nodegroup, you would need to create a new nodegroup with the desired changes, move the load and delete the old one.\
Please refer to [Deleting and draining](../../../managing-nodegroups.md#deleting-and-draining) documentation for further details.
## Ocean VNGs Advantage
By using Ocean VNGs, those changes on AWS nodegroups are made for you automatically by only modifying the configuration of the spotOcean object of the nodegroup.
