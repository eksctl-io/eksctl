# Creating a Cluster

You can add an Ocean nodegroup to new or existing clusters.\
To create a new cluster managed by Spot Ocean run:

```bash
eksctl create cluster \
 --name example \
 --spot-ocean
 --managed=false
```

The command will create an eks cluster "example" managed by Spot Ocean.
### Note:
Ocean nodegroups are integrated with eks unmanaged nodegroups and managed by Ocean.
nodegroups that managed by ocean are called Virtual Node Groups (VNGs).

spotOcean object in the configuration file on cluster level will define the default VNG properties, spotOcean object in the configuration file on a specific nodegroup level will define a custom VNG properties for this nodegroup.\
Ocean nodegroups will inherit their default parameters from the cluster’s default VNG unless
those parameters are set for the nodegroup’s custom VNG, you can find more info about it [here](https://docs.spot.io/ocean/features/vngs/?id=default-vng).

## Creating a Cluster With Config File
To have more control over the cluster configuration, like creating multiple Ocean nodegroups, a configuration file can be used.

Please see below the cluster configuration examples:

```yaml
# cluster.yaml
# A cluster with two Ocean nodegroups.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: example
  region: us-west-2

    # enable Ocean integration on the entire cluster
spotOcean:
    strategy:
        utilizeReservedInstances: true
        fallbackToOnDemand: true

    scheduling:
        shutdownHours:
            isEnabled: true
            timeWindows:
                - Mon:22:00-Tue:06:00
                - Tue:22:00-Wed:06:00
                - Wed:22:00-Thu:06:00
                - Thu:22:00-Fri:06:00
                - Fri:22:00-Mon:06:00

        tasks:
            - isEnabled: true
              taskType: manualHeadroomUpdate
              cronExpression: 0 1 * * *
              config:
                  headrooms:
                      - cpuPerUnit: 2000
                        memoryPerUnit: 4000
                        gpuPerUnit: 0
                        numOfUnits: 1
            - isEnabled: true
              taskType: manualHeadroomUpdate
              cronExpression: 0 2 * * *
              config:
                  headrooms:
                      - cpuPerUnit: 0
                        memoryPerUnit: 0
                        gpuPerUnit: 0
                        numOfUnits: 0

    autoScaler:
        enabled: true
        cooldown: 300
        autoConfig: false
        headrooms:
            cpuPerUnit: 2
            gpuPerUnit: 0
            memoryPerUnit: 64
            numOfUnits: 1

    compute:
        instanceTypes:
            whitelist:
                - t3a.large
                - t3a.xlarge
                - t3a.2xlarge
                - m5a.large
                - m5a.xlarge
                - m5a.2xlarge
                - m5a.4xlarge
                - c5.large
                - c5.xlarge
                - c5.2xlarge
                - c5.4xlarge

nodeGroups:
- name: ocean-ng-1
  minSize: 2
  maxSize: 4
  desiredCapacity: 3
  volumeSize: 20
  ssh:
    allow: true
    publicKeyPath: ~/.ssh/ec2_id_rsa.pub
    sourceSecurityGroupIds: ["sg-00241fbb12c607007"]
  labels: {role: worker}
  tags:
    nodegroup-role: worker
  iam:
    withAddonPolicies:
      externalDNS: true
      certManager: true
  # enable Ocean integration on nodegroup
  spotOcean:
    strategy:
    # Percentage of Spot instances that would spin up from the desired capacity.
      spotPercentage: 100
    compute:
      instanceTypes:
        - c4.large
        - t2.large

- name: ocean-ng-2
  instanceType: t2.large
  minSize: 2
  maxSize: 3
  # enable Ocean integration on nodegroup
  spotOcean:
    strategy:
    # Percentage of Spot instances that would spin up from the desired capacity.
      spotPercentage: 100
    compute:
      instanceTypes:
        - c4.large
        - t2.large
```

### Note:
The following configuration file example enables us to create cluster with both Ocean-managed and unmanaged nodegroups.

```yaml
# cluster.yaml
# A cluster with an unmanaged nodegroup and an Ocean-managed nodegroup.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: example
  region: us-west-2

# enable Ocean integration on the entire cluster
spotOcean:
    strategy:
        utilizeReservedInstances: true
        fallbackToOnDemand: true

    scheduling:
        shutdownHours:
            isEnabled: true
            timeWindows:
                - Mon:22:00-Tue:06:00
                - Tue:22:00-Wed:06:00
                - Wed:22:00-Thu:06:00
                - Thu:22:00-Fri:06:00
                - Fri:22:00-Mon:06:00

        tasks:
            - isEnabled: true
              taskType: manualHeadroomUpdate
              cronExpression: 0 1 * * *
              config:
                  headrooms:
                      - cpuPerUnit: 2000
                        memoryPerUnit: 4000
                        gpuPerUnit: 0
                        numOfUnits: 1
            - isEnabled: true
              taskType: manualHeadroomUpdate
              cronExpression: 0 2 * * *
              config:
                  headrooms:
                      - cpuPerUnit: 0
                        memoryPerUnit: 0
                        gpuPerUnit: 0
                        numOfUnits: 0

    autoScaler:
        enabled: true
        cooldown: 300
        autoConfig: false
        headrooms:
            cpuPerUnit: 2
            gpuPerUnit: 0
            memoryPerUnit: 64
            numOfUnits: 1

    compute:
        instanceTypes:
            whitelist:
                - t3a.large
                - t3a.xlarge
                - t3a.2xlarge
                - m5a.large
                - m5a.xlarge
                - m5a.2xlarge
                - m5a.4xlarge
                - c5.large
                - c5.xlarge
                - c5.2xlarge
                - c5.4xlarge

nodeGroups:
# Use Ocean Managed Group
- name: ocean-ng-1
  minSize: 2
  maxSize: 4
  desiredCapacity: 3
  volumeSize: 20
  ssh:
    allow: true
    publicKeyPath: ~/.ssh/ec2_id_rsa.pub # provide correct public access key name
    sourceSecurityGroupIds: ["sg-00241fbb12c607007"] # update correct security group id
  labels: {role: worker}
  tags:
    nodegroup-role: worker
  iam:
    withAddonPolicies:
      externalDNS: true
      certManager: true
  # enable Ocean integration on nodegroup
  spotOcean:
    strategy:
      # Percentage of Spot instances that would spin up from the desired capacity.
      spotPercentage: 100
    compute:
      instanceTypes:
        - c4.large
        - t2.large


# Use AWS Auto Scaling group.
- name: ng-2
  instanceType: t2.large
  privateNetworking: true
  minSize: 2
  maxSize: 3
```

To create a new cluster managed by Spot Ocean with the definition file run:

```bash
eksctl create cluster -f cluster.yaml
```

It's also possible to create the cluster without the nodegroups by running:

```bash
eksctl create cluster -f cluster.yaml --without-nodegroup
```

- Also see [Spot Ocean Nodegroups](./nodegroup/spot-ocean-nodegroups.md/#nodegroups)
