---
title: Config file schema
weight: 200
---

```yaml
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
