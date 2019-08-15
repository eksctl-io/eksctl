---
title: "Configuration file YAML reference"
weight: 30
---


# Configuration file YAML reference

## apiVersion <a name="apiVersion"></a>

The `eksctl` API version that supports this YAML configuration file.

### valid values

`eksctl.io/v1alpha5`

### example

`apiVersion: eksctl.io/v1alpha5`

### parent

Top-level YAML field.

## availabilityZones <a name="availabilityZones"></a>

The AWS availability zones for the managed cluster or nodegroup represented as an array of strings.

### valid values

All supported AWS availability zones.

### example

For a single zone:
`availabilityZones: ["eu-west-2a"]`

For multiple zones:
`availabilityZones: ["eu-west-2a", "eu-west-2b"]`

### parent

- Top-level YAML field for the ClusterConfig type
- [nodegroups](#nodegroups)
  - Inherited from the cluster if unspecified.

### flag

- [zones](../02-flags#zones)
- [node-zones](../02-flags#node-zones)

## cloudWatch <a name="cloudWatch"></a>

CloudWatch configurations for the managed EKS cluster.

**Also used as** a boolean attribute of nodegroup IAM add-on policies.

### parent

- Top-level YAML field for the ClusterConfig type
- [withAddonPolicies](#withAddonPolicies)

### flag

- [enable-types](../02-flags#enable-types)

### command

- [utils update-cluster-logging](../01-commands#utils-update-cluster-logging)

## iam <a name="iam"></a>

### parent

Top-level YAML field for the ClusterConfig type.

## kind <a name="kind"></a>

The kind of `eksctl` object defined in this YAML configuration file.

***Note***: The `ClusterConfig` kind is currently the only supported top-level object.

### valid values

`ClusterConfig`

### example

`kind: ClusterConfig`

### parent

Top-level YAML field.

### command

- [create cluster](../01-commands#create-common-flags)

## metadata <a name="metadata"></a>

Metadata describing the EKS cluster.

### example

```
metadata:
  name: basic-cluster
  region: eu-north-1
```

### parent

Top-level YAML field for the ClusterConfig type.

## nodeGroups <a name="nodeGroups"></a>

### parent

Top-level YAML field for the ClusterConfig type.

## status <a name="status"></a>

### parent

Top-level YAML field for the ClusterConfig type.

## vpc <a name="vpc"></a>

### parent

Top-level YAML field for the ClusterConfig type.


# YAML Schema

```yaml
ClusterCloudWatch:
  additionalProperties: false
  properties:
    clusterLogging:
      $ref: '#/definitions/ClusterCloudWatchLogging'
      $schema: http://json-schema.org/draft-04/schema#
  type: object
ClusterCloudWatchLogging:
  additionalProperties: false
  properties:
    enableTypes:
      items:
        type: string
      type: array
  type: object
ClusterConfig:
  additionalProperties: false
  properties:
    TypeMeta:
      $ref: '#/definitions/TypeMeta'
      $schema: http://json-schema.org/draft-04/schema#
    availabilityZones:
      items:
        type: string
      type: array
    cloudWatch:
      $ref: '#/definitions/ClusterCloudWatch'
      $schema: http://json-schema.org/draft-04/schema#
    iam:
      $ref: '#/definitions/ClusterIAM'
      $schema: http://json-schema.org/draft-04/schema#
    metadata:
      $ref: '#/definitions/ClusterMeta'
      $schema: http://json-schema.org/draft-04/schema#
    nodeGroups:
      items:
        $ref: '#/definitions/NodeGroup'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    status:
      $ref: '#/definitions/ClusterStatus'
      $schema: http://json-schema.org/draft-04/schema#
    vpc:
      $ref: '#/definitions/ClusterVPC'
      $schema: http://json-schema.org/draft-04/schema#
  required:
  - TypeMeta
  - metadata
  - iam
  type: object
ClusterIAM:
  additionalProperties: false
  properties:
    serviceRoleARN:
      type: string
  type: object
ClusterMeta:
  additionalProperties: false
  properties:
    name:
      type: string
    region:
      type: string
    tags:
      patternProperties:
        .*:
          type: string
      type: object
    version:
      type: string
  required:
  - name
  - region
  type: object
ClusterNAT:
  additionalProperties: false
  properties:
    gateway:
      type: string
  type: object
ClusterStatus:
  additionalProperties: false
  properties:
    arn:
      type: string
    certificateAuthorityData:
      media:
        binaryEncoding: base64
      type: string
    endpoint:
      type: string
    stackName:
      type: string
  type: object
ClusterSubnets:
  additionalProperties: false
  properties:
    private:
      patternProperties:
        .*:
          $ref: '#/definitions/Network'
      type: object
    public:
      patternProperties:
        .*:
          $ref: '#/definitions/Network'
      type: object
  type: object
ClusterVPC:
  additionalProperties: false
  properties:
    Network:
      $ref: '#/definitions/Network'
      $schema: http://json-schema.org/draft-04/schema#
    autoAllocateIPv6:
      type: boolean
    extraCIDRs:
      items:
        $ref: '#/definitions/IPNet'
      type: array
    nat:
      $ref: '#/definitions/ClusterNAT'
      $schema: http://json-schema.org/draft-04/schema#
    securityGroup:
      type: string
    sharedNodeSecurityGroup:
      type: string
    subnets:
      $ref: '#/definitions/ClusterSubnets'
      $schema: http://json-schema.org/draft-04/schema#
  required:
  - Network
  type: object
IPNet:
  additionalProperties: false
  properties:
    IP:
      format: ipv4
      type: string
    Mask:
      items:
        type: integer
      type: array
  required:
  - IP
  - Mask
  type: object
Network:
  additionalProperties: false
  properties:
    cidr:
      $ref: '#/definitions/IPNet'
      $schema: http://json-schema.org/draft-04/schema#
    id:
      type: string
  type: object
NodeGroup:
  additionalProperties: false
  properties:
    ami:
      type: string
    amiFamily:
      type: string
    availabilityZones:
      items:
        type: string
      type: array
    clusterDNS:
      type: string
    desiredCapacity:
      type: integer
    iam:
      $ref: '#/definitions/NodeGroupIAM'
      $schema: http://json-schema.org/draft-04/schema#
    instanceType:
      type: string
    instancesDistribution:
      $ref: '#/definitions/NodeGroupInstancesDistribution'
      $schema: http://json-schema.org/draft-04/schema#
    kubeletExtraConfig:
      patternProperties:
        .*:
          additionalProperties: true
          type: object
      type: object
    labels:
      patternProperties:
        .*:
          type: string
      type: object
    maxPodsPerNode:
      type: integer
    maxSize:
      type: integer
    minSize:
      type: integer
    name:
      type: string
    overrideBootstrapCommand:
      type: string
    preBootstrapCommands:
      items:
        type: string
      type: array
    privateNetworking:
      type: boolean
    securityGroups:
      $ref: '#/definitions/NodeGroupSGs'
      $schema: http://json-schema.org/draft-04/schema#
    ssh:
      $ref: '#/definitions/NodeGroupSSH'
      $schema: http://json-schema.org/draft-04/schema#
    tags:
      patternProperties:
        .*:
          type: string
      type: object
    taints:
      patternProperties:
        .*:
          type: string
      type: object
    targetGroupARNs:
      items:
        type: string
      type: array
    volumeEncrypted:
      type: boolean
    volumeIOPS:
      type: integer
    volumeKmsKeyID:
      type: string
    volumeName:
      type: string
    volumeSize:
      type: integer
    volumeType:
      type: string
  required:
  - name
  - privateNetworking
  - volumeSize
  - volumeType
  - volumeIOPS
  - ssh
  - iam
  type: object
NodeGroupIAM:
  additionalProperties: false
  properties:
    attachPolicyARNs:
      items:
        type: string
      type: array
    instanceProfileARN:
      type: string
    instanceRoleARN:
      type: string
    instanceRoleName:
      type: string
    withAddonPolicies:
      $ref: '#/definitions/NodeGroupIAMAddonPolicies'
      $schema: http://json-schema.org/draft-04/schema#
  type: object
NodeGroupIAMAddonPolicies:
  additionalProperties: false
  properties:
    albIngress:
      type: boolean
    appMesh:
      type: boolean
    autoScaler:
      type: boolean
    certManager:
      type: boolean
    cloudWatch:
      type: boolean
    ebs:
      type: boolean
    efs:
      type: boolean
    externalDNS:
      type: boolean
    fsx:
      type: boolean
    imageBuilder:
      type: boolean
    xRay:
      type: boolean
  required:
  - imageBuilder
  - autoScaler
  - externalDNS
  - certManager
  - appMesh
  - ebs
  - fsx
  - efs
  - albIngress
  - xRay
  - cloudWatch
  type: object
NodeGroupInstancesDistribution:
  additionalProperties: false
  properties:
    instanceTypes:
      items:
        type: string
      type: array
    maxPrice:
      type: number
    onDemandBaseCapacity:
      type: integer
    onDemandPercentageAboveBaseCapacity:
      type: integer
    spotInstancePools:
      type: integer
  required:
  - instanceTypes
  - onDemandBaseCapacity
  - onDemandPercentageAboveBaseCapacity
  - spotInstancePools
  type: object
NodeGroupSGs:
  additionalProperties: false
  properties:
    attachIDs:
      items:
        type: string
      type: array
    withLocal:
      type: boolean
    withShared:
      type: boolean
  required:
  - withShared
  - withLocal
  type: object
NodeGroupSSH:
  additionalProperties: false
  properties:
    allow:
      type: boolean
    publicKey:
      type: string
    publicKeyName:
      type: string
    publicKeyPath:
      type: string
  required:
  - allow
  type: object
TypeMeta:
  additionalProperties: false
  properties:
    apiVersion:
      type: string
    kind:
      type: string
  type: object
```

# Examples

## Create Cluster

You can create a cluster using a config file instead of flags.

First, create `cluster.yaml` file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: basic-cluster
  region: eu-north-1

nodeGroups:
  - name: ng-1
    instanceType: m5.large
    desiredCapacity: 10
    ssh:
      allow: true # will use ~/.ssh/id_rsa.pub as the default ssh key
  - name: ng-2
    instanceType: m5.xlarge
    desiredCapacity: 2
    ssh:
      publicKeyPath: ~/.ssh/ec2_id_rsa.pub
```

If you needed to use an existing VPC, you can use a config file like this:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-in-existing-vpc
  region: eu-north-1

vpc:
  subnets:
    private:
      eu-north-1a: { id: subnet-0ff156e0c4a6d300c }
      eu-north-1b: { id: subnet-0549cdab573695c03 }
      eu-north-1c: { id: subnet-0426fb4a607393184 }

nodeGroups:
  - name: ng-1-workers
    labels: { role: workers }
    instanceType: m5.xlarge
    desiredCapacity: 10
    privateNetworking: true
  - name: ng-2-builders
    labels: { role: builders }
    instanceType: m5.2xlarge
    desiredCapacity: 2
    privateNetworking: true
    iam:
      withAddonPolicies:
        imageBuilder: true
```

## Managing Nodegroups

Nodegroups can also be created through a cluster definition or config file. Given the following example config file
and an existing cluster called ``dev-cluster:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev-cluster
  region: eu-north-1

nodeGroups:
  - name: ng-1-workers
    labels: { role: workers }
    instanceType: m5.xlarge
    desiredCapacity: 10
    privateNetworking: true
  - name: ng-2-builders
    labels: { role: builders }
    instanceType: m5.2xlarge
    desiredCapacity: 2
    privateNetworking: true
```

## Autoscaling

BEFORE:

```yaml
nodeGroups:
  - name: ng1-public
    instanceType: m5.xlarge
    # availabilityZones: ["eu-west-2a", "eu-west-2b"]
```

AFTER:

```yaml
nodeGroups:
  - name: ng1-public-2a
    instanceType: m5.xlarge
    availabilityZones: ["eu-west-2a"]
  - name: ng1-public-2b
    instanceType: m5.xlarge
    availabilityZones: ["eu-west-2b"]
```

## VPC Networking

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-north-1

nodeGroups:
- name: ng-1
  clusterDNS: 169.254.20.10
```

Note that this configuration only accepts one IP address. To specify more than one address, use the
[`extraKubeletConfig` parameter](../customizing-the-kubelet):

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-north-1

nodeGroups:
  - name: ng-1
    kubeletExtraConfig:
        clusterDNS: ["169.254.20.10","172.20.0.10"]
```

The NAT Gateway for a cluster can be configured to be `Disabled`, `Single` (default) or `HighlyAvailable`. It can be
specified through the `--vpc-nat-mode` CLI flag or in the cluster config file like the example below:


```yaml
vpc:
  nat:
    gateway: HighlyAvailable # other options: Disable, Single (default)
```

## Spot Instances

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

## IAM policies

Example of all supported add-on policies:

```yaml
nodeGroups:
  - name: ng-1
    instanceType: m5.xlarge
    desiredCapacity: 1
    iam:
      withAddonPolicies:
        imageBuilder: true
        autoScaler: true
        externalDNS: true
        certManager: true
        appMesh: true
        ebs: true
        fsx: true
        efs: true
        albIngress: true
        xRay: true
        cloudWatch: true
```
### Adding a custom instance role

This example creates a nodegroup that reuses an existing IAM Instance Role from another cluster:

```yaml
apiVersion: eksctl.io/v1alpha4
kind: ClusterConfig
metadata:
  name: test-cluster-c-1
  region: eu-north-1

nodeGroups:
  - name: ng2-private
    instanceType: m5.large
    desiredCapacity: 1
    iam:
      instanceProfileARN: "arn:aws:iam::123:instance-profile/eksctl-test-cluster-a-3-nodegroup-ng2-private-NodeInstanceProfile-Y4YKHLNINMXC"
      instanceRoleARN: "arn:aws:iam::123:role/eksctl-test-cluster-a-3-nodegroup-NodeInstanceRole-DNGMQTQHQHBJ"
```

## Attaching policies by ARN

```yaml
nodeGroups:
  - name: my-special-nodegroup
    iam:
      attachPolicyARNs:
        - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy
        - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
        - arn:aws:iam::aws:policy/ElasticLoadBalancingFullAccess
        - arn:aws:iam::1111111111:policy/kube2iam
      withAddonPolicies:
        autoScaler: true
        imageBuilder: true
```

## Customizing the Kubelet

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev-cluster-1
  region: eu-north-1

nodeGroups:
  - name: ng-1
    instanceType: m5a.xlarge
    desiredCapacity: 1
    kubeletExtraConfig:
        kubeReserved:
            cpu: "300m"
            memory: "300Mi"
            ephemeral-storage: "1Gi"
        kubeReservedCgroup: "/kube-reserved"
        systemReserved:
            cpu: "300m"
            memory: "300Mi"
            ephemeral-storage: "1Gi"
        evictionHard:
            memory.available:  "200Mi"
            nodefs.available: "10%"
        featureGates:
            DynamicKubeletConfig: true
            RotateKubeletServerCertificate: true # has to be enabled, otherwise it will be disabled
```

## Cloudwatch Cluster Logging

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

## Troubleshooting

### subnet ID "subnet-11111111" is not the same as "subnet-22222222"

Given a config file specifying subnets for a VPC like the following:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test
  region: us-east-1

vpc:
  subnets:
    public:
      us-east-1a: {id: subnet-11111111}
      us-east-1b: {id: subnet-22222222}
    private:
      us-east-1a: {id: subnet-33333333}
      us-east-1b: {id: subnet-44444444}

nodeGroups: []
```

An error `subnet ID "subnet-11111111" is not the same as "subnet-22222222"` means that the subnets specified are not 
placed in the right Availability zone. Check in the AWS console which is the right subnet ID for each Availability Zone.

In this example, the correct configuration for the VPC would be:

```yaml
vpc:
  subnets:
    public:
      us-east-1a: {id: subnet-22222222}
      us-east-1b: {id: subnet-11111111}
    private:
      us-east-1a: {id: subnet-33333333}
      us-east-1b: {id: subnet-44444444}
```
