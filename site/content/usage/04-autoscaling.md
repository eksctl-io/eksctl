---
title: "Auto Scaling"
weight: 40
url: usage/autoscaling
---

## Enable Auto Scaling

You can create a cluster (or nodegroup in an existing cluster) with IAM role that will allow use of [cluster autoscaler][]:

```
eksctl create cluster --asg-access
```

Once cluster is running, you will need to install [cluster autoscaler][] itself. This flag also sets `k8s.io/cluster-autoscaler/enabled`
and `k8s.io/cluster-autoscaler/<clusterName>` tags, so nodegroup discovery should work.

### Scaling up from 0

If you'd like to be able to scale your node group up from 0 and you have
labels and/or taints defined on your nodegroups you'll need corresponding
tags on your ASGs.  You can do this with the `tags` key on your node group
definitions.  For example, given a node group with the following labels and
taints:

```yaml
nodeGroups:
  - name: ng1-public
    ...
    labels:
      my-cool-label: pizza
    taints:
      feaster: "true:NoSchedule"
```

You would need to add the following ASG tags:

```yaml
nodeGroups:
  - name: ng1-public
    ...
    labels:
      my-cool-label: pizza
    taints:
      feaster: "true:NoSchedule"
    tags:
      k8s.io/cluster-autoscaler/node-template/label/my-cool-label: pizza
      k8s.io/cluster-autoscaler/node-template/taint/feaster: "true:NoSchedule"
```

You can read more about this
[here](https://github.com/weaveworks/eksctl/issues/1066) and
[here](https://github.com/kubernetes/autoscaler/issues/2418).

[cluster autoscaler]: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md

### Zone-aware Auto Scaling

If your workloads are zone-specific you'll need to create separate nodegroups for each zone. This is because the `cluster-autoscaler` assumes that all nodes in a group are exactly equivalent. So, for example, if a scale-up event is triggered by a pod which needs a zone-specific PVC (e.g. an EBS volume), the new node might get scheduled in the wrong AZ and the pod will fail to start.

You won't need a separate nodegroup for each AZ if your environment meets the following criteria:

- No zone-specific storage requirements.
- No required podAffinity with topology other than host.
- No required nodeAffinity on zone label.
- No nodeSelector on a zone label.

(Read more [here](https://github.com/kubernetes/autoscaler/pull/1802#issuecomment-474295002) and [here](https://github.com/weaveworks/eksctl/pull/647#issuecomment-474698054).)

If you meet all of the above requirements (and possibly others) then you should be safe with a single nodegroup which spans multiple AZs. Otherwise you'll want to create separate, single-AZ nodegroups:

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
