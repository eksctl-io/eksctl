# Spot Ocean Nodegroups

[Ocean](https://spot.io/products/ocean/) by [Spot](https://spot.io/) simplifies infrastructure management for Kubernetes.  With robust, container-driven infrastructure auto-scaling and intelligent right-sizing for container resource requirements, operations can literally "set and forget" the underlying cluster.

Ocean seamlessly integrates with your existing nodegroups, as a drop-in replacement for AWS Auto Scaling groups, and allows you to streamline and optimize the entire workflow, from initially creating your cluster to managing and optimizing it on an ongoing basis.

## Features

- **Simplify Cluster Management** —
Ocean's Virtual Node Groups make it easy to run different infrastructure in a single cluster, which can span multiple AWS VPC availability zones and subnets for high-availability.

- **Container-Driven Autoscaling and Vertical Rightsizing** —
Auto-detect your container infrastructure requirements so the appropriate instance size or type will always be available. Measure real-time CPU/Memory consumption of your Pods for ongoing resource optimization.

- **Cloud-Native Showback** —
Gain a granular view of your cluster's cost breakdown (compute and storage) for each and every one of the cluster's resources such as Namespaces, Deployments, Daemon Sets, Jobs, and Pods.

- **Optimized Pricing and Utilization** —
Ocean not only intelligently leverages Spot Instances and reserved capacity to reduce costs, but also eliminates underutilized instances with container-driven autoscaling and advanced bin-packing.

## Prerequisites

Make sure you have [installed eksctl](https://eksctl.io/introduction/#installation) and [installed kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

## Setup your environment

### Spot
Generate your credentials [here](https://console.spotinst.com/spt/settings/tokens/permanent). If you are not a Spot Ocean user, sign up for free [here](https://console.spotinst.com/spt/auth/signUp). For further information, please checkout our [Spot API](https://help.spot.io/spotinst-api/) guide, available on the [Spot Help Center](https://help.spot.io/) website.

To use environment variables, run:
```bash
export SPOTINST_TOKEN=<spotinst_token>
export SPOTINST_ACCOUNT=<spotinst_account>
```

To use credentials file, run the [spotctl configure](https://github.com/spotinst/spotctl#getting-started) command:
```bash
spotctl configure
? Enter your access token [? for help] **********************************
? Select your default account  [Use arrows to move, ? for more help]
> act-01234567 (prod)
  act-0abcdefg (dev)
```

Or, manually create an INI formatted file like this:
```ini
[default]
token   = <spotinst_token>
account = <spotinst_account>
```

and place it in:

- Unix/Linux/macOS:
```bash
~/.spotinst/credentials
```
- Windows:
```bash
%UserProfile%\.spotinst\credentials
```

### AWS

Make sure to set up your AWS credentials. Please refer to [Configuration and Credential File Settings](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) for further details.

## Creating a Cluster

You can add an Ocean nodegroup to new or existing clusters. To create a new cluster with a Ocean nodegroup, run:

```bash
eksctl create cluster --spot-ocean
```

To create multiple Ocean nodegroups and have more control over the configuration, a config file can be used.

!!!note
    Ocean nodegroups are integrated with unmanaged nodegroups and managed by Ocean.

```yaml
# cluster.yaml
# A cluster with two Ocean nodegroups.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: example
  region: us-west-2

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
  spotOcean: {} # enable Ocean integration

- name: ocean-ng-2
  instanceType: t2.large
  minSize: 2
  maxSize: 3
  spotOcean: {} # enable Ocean integration
```

!!!note
    It's possible to have a cluster with both Ocean-managed and unmanaged nodegroups.

```yaml
# cluster.yaml
# A cluster with an unmanaged nodegroup and an Ocean-managed nodegroup.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: example
  region: us-west-2

nodeGroups:

# Use Ocean.
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
  spotOcean: {} # enable Ocean integration

# Use AWS Auto Scaling group.
- name: ng-2
  instanceType: t2.large
  privateNetworking: true
  minSize: 2
  maxSize: 3
```

## Creating a Nodegroup

To create a new nodegroup, run:

```bash
eksctl create nodegroup \
  --cluster <cluster-name> \
  --nodegroup-name <nodegroup-name> \
  --spot-ocean
```

To create multiple nodegroups and have more control over the configuration, a config file can be used.

```yaml
# cluster.yaml
# A cluster with two Ocean nodegroups.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: example
  region: us-west-2

nodeGroups:
- name: ocean-ng-1
  [... nodegroup standard fields; ssh, tags, etc.]

  # Enable Ocean integration and use all defaults.
  spotOcean: {}

- name: ocean-ng-2
  [... nodegroup standard fields; ssh, tags, etc.]

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
      # Enable the Ocean autoscaler.
      enabled: true

      # Cooldown period between scaling actions.
      cooldown: 300

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

## Nodegroups Immutability
By design, nodegroups are immutable. This means that if you need to change something like the AMI or the instance type of a nodegroup, you would need to create a new nodegroup with the desired changes, move the load and delete the old one. Check [Deleting and draining](managing-nodegroups.md#deleting-and-draining).

## Documentation

If you're new to [Spot](https://spot.io/) and want to get started, please checkout our [Getting Started](https://help.spot.io/getting-started-with-spotinst/) guide, available on the [Spot Help Center](https://help.spot.io/) website.

## Getting Help

Please use these community resources for getting help:

- Join our [Spot](https://spot.io/) community on [Slack](http://slack.spot.io/).
- Open a GitHub [issue](https://github.com/spotinst/weaveworks-eksctl/issues/new/choose/).
- Ask a question on [Stack Overflow](https://stackoverflow.com/) and tag it with `spot-ocean`.
