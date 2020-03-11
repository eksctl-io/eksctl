---
title: Config file schema
weight: 200
url: usage/schema
---

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
    fargateProfiles:
      items:
        $ref: '#/definitions/FargateProfile'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    iam:
      $ref: '#/definitions/ClusterIAM'
      $schema: http://json-schema.org/draft-04/schema#
    managedNodeGroups:
      items:
        $ref: '#/definitions/ManagedNodeGroup'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    metadata:
      $ref: '#/definitions/ClusterMeta'
      $schema: http://json-schema.org/draft-04/schema#
    nodeGroups:
      items:
        $ref: '#/definitions/NodeGroup'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    secretsEncryption:
      $ref: '#/definitions/SecretsEncryption'
      $schema: http://json-schema.org/draft-04/schema#
    status:
      $ref: '#/definitions/ClusterStatus'
      $schema: http://json-schema.org/draft-04/schema#
    vpc:
      $ref: '#/definitions/ClusterVPC'
      $schema: http://json-schema.org/draft-04/schema#
  required:
  - TypeMeta
  - metadata
  type: object
ClusterEndpoints:
  additionalProperties: false
  properties:
    privateAccess:
      type: boolean
    publicAccess:
      type: boolean
  type: object
ClusterIAM:
  additionalProperties: false
  properties:
    fargatePodExecutionRoleARN:
      type: string
    fargatePodExecutionRolePermissionsBoundary:
      type: string
    serviceAccounts:
      items:
        $ref: '#/definitions/ClusterIAMServiceAccount'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    serviceRoleARN:
      type: string
    serviceRolePermissionsBoundary:
      type: string
    withOIDC:
      type: boolean
  type: object
ClusterIAMServiceAccount:
  additionalProperties: false
  properties:
    attachPolicy:
      patternProperties:
        .*:
          additionalProperties: true
          type: object
      type: object
    attachPolicyARNs:
      items:
        type: string
      type: array
    metadata:
      $ref: '#/definitions/ObjectMeta'
      $schema: http://json-schema.org/draft-04/schema#
    permissionsBoundary:
      type: string
    status:
      $ref: '#/definitions/ClusterIAMServiceAccountStatus'
      $schema: http://json-schema.org/draft-04/schema#
  type: object
ClusterIAMServiceAccountStatus:
  additionalProperties: false
  properties:
    roleARN:
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
    clusterEndpoints:
      $ref: '#/definitions/ClusterEndpoints'
      $schema: http://json-schema.org/draft-04/schema#
    extraCIDRs:
      items:
        $ref: '#/definitions/IPNet'
      type: array
    nat:
      $ref: '#/definitions/ClusterNAT'
      $schema: http://json-schema.org/draft-04/schema#
    publicAccessCIDRs:
      items:
        type: string
      type: array
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
FargateProfile:
  additionalProperties: false
  properties:
    name:
      type: string
    podExecutionRoleARN:
      type: string
    selectors:
      items:
        $ref: '#/definitions/FargateProfileSelector'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    subnets:
      items:
        type: string
      type: array
  required:
  - name
  - selectors
  type: object
FargateProfileSelector:
  additionalProperties: false
  properties:
    labels:
      patternProperties:
        .*:
          type: string
      type: object
    namespace:
      type: string
  required:
  - namespace
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
Initializer:
  additionalProperties: false
  properties:
    name:
      type: string
  required:
  - name
  type: object
Initializers:
  additionalProperties: false
  properties:
    pending:
      items:
        $ref: '#/definitions/Initializer'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    result:
      $ref: '#/definitions/Status'
      $schema: http://json-schema.org/draft-04/schema#
  required:
  - pending
  type: object
ListMeta:
  additionalProperties: false
  properties:
    continue:
      type: string
    resourceVersion:
      type: string
    selfLink:
      type: string
  type: object
ManagedNodeGroup:
  additionalProperties: false
  properties:
    ScalingConfig:
      $ref: '#/definitions/ScalingConfig'
      $schema: http://json-schema.org/draft-04/schema#
    amiFamily:
      type: string
    availabilityZones:
      items:
        type: string
      type: array
    iam:
      $ref: '#/definitions/NodeGroupIAM'
    instanceType:
      type: string
    labels:
      patternProperties:
        .*:
          type: string
      type: object
    name:
      type: string
    ssh:
      $ref: '#/definitions/NodeGroupSSH'
    tags:
      patternProperties:
        .*:
          type: string
      type: object
    volumeSize:
      type: integer
  required:
  - name
  - ScalingConfig
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
    bottlerocket:
      $ref: '#/definitions/NodeGroupBottlerocket'
      $schema: http://json-schema.org/draft-04/schema#
    classicLoadBalancerNames:
      items:
        type: string
      type: array
    clusterDNS:
      type: string
    desiredCapacity:
      type: integer
    ebsOptimized:
      type: boolean
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
  - iam
  type: object
NodeGroupBottlerocket:
  additionalProperties: false
  properties:
    enableAdminContainer:
      type: boolean
    settings:
      patternProperties:
        .*:
          additionalProperties: true
          type: object
      type: object
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
    instanceRolePermissionsBoundary:
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
    spotAllocationStrategy:
      type: string
    spotInstancePools:
      type: integer
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
    sourceSecurityGroupIds:
      items:
        type: string
      type: array
  required:
  - allow
  type: object
ObjectMeta:
  additionalProperties: false
  properties:
    annotations:
      patternProperties:
        .*:
          type: string
      type: object
    clusterName:
      type: string
    creationTimestamp:
      $ref: '#/definitions/Time'
      $schema: http://json-schema.org/draft-04/schema#
    deletionGracePeriodSeconds:
      type: integer
    deletionTimestamp:
      $ref: '#/definitions/Time'
    finalizers:
      items:
        type: string
      type: array
    generateName:
      type: string
    generation:
      type: integer
    initializers:
      $ref: '#/definitions/Initializers'
      $schema: http://json-schema.org/draft-04/schema#
    labels:
      patternProperties:
        .*:
          type: string
      type: object
    name:
      type: string
    namespace:
      type: string
    ownerReferences:
      items:
        $ref: '#/definitions/OwnerReference'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    resourceVersion:
      type: string
    selfLink:
      type: string
    uid:
      type: string
  type: object
OwnerReference:
  additionalProperties: false
  properties:
    apiVersion:
      type: string
    blockOwnerDeletion:
      type: boolean
    controller:
      type: boolean
    kind:
      type: string
    name:
      type: string
    uid:
      type: string
  required:
  - apiVersion
  - kind
  - name
  - uid
  type: object
ScalingConfig:
  additionalProperties: false
  properties:
    desiredCapacity:
      type: integer
    maxSize:
      type: integer
    minSize:
      type: integer
  type: object
SecretsEncryption:
  additionalProperties: false
  properties:
    keyARN:
      type: string
  type: object
Status:
  additionalProperties: false
  properties:
    TypeMeta:
      $ref: '#/definitions/TypeMeta'
    code:
      type: integer
    details:
      $ref: '#/definitions/StatusDetails'
      $schema: http://json-schema.org/draft-04/schema#
    message:
      type: string
    metadata:
      $ref: '#/definitions/ListMeta'
      $schema: http://json-schema.org/draft-04/schema#
    reason:
      type: string
    status:
      type: string
  required:
  - TypeMeta
  type: object
StatusCause:
  additionalProperties: false
  properties:
    field:
      type: string
    message:
      type: string
    reason:
      type: string
  type: object
StatusDetails:
  additionalProperties: false
  properties:
    causes:
      items:
        $ref: '#/definitions/StatusCause'
        $schema: http://json-schema.org/draft-04/schema#
      type: array
    group:
      type: string
    kind:
      type: string
    name:
      type: string
    retryAfterSeconds:
      type: integer
    uid:
      type: string
  type: object
Time:
  additionalProperties: false
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
