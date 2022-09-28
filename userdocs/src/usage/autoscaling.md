# Auto Scaling

## Enable Auto Scaling

You can create a cluster (or nodegroup in an existing cluster) with IAM role that will allow use of [cluster autoscaler][]:

```console
eksctl create cluster --asg-access
```

This flag also sets `k8s.io/cluster-autoscaler/enabled`
and `k8s.io/cluster-autoscaler/<clusterName>` tags, so nodegroup discovery should work.

Once the cluster is running, you will need to install [Cluster Autoscaler][] itself. 

You should also add the following to your managed or unmanaged nodegroup definition(s) to add the tags required for the Cluster Autoscaler to scale the nodegroup:
```yaml
nodeGroups:
  - name: ng1-public
    iam:
      withAddonPolicies:
        autoScaler: true
```

### Scaling up from 0

If you would like to be able to scale your node group up from 0 and you have
labels and/or taints defined on your nodegroups, you will need to propagate these as
tags on your Auto Scaling Groups (ASGs). 

One way to do this is by setting the ASG tags in the `tags` field of your nodegroup
definitions. For example, given a nodegroup with the following labels and
taints:

```yaml
nodeGroups:
  - name: ng1-public
    ...
    labels:
      my-cool-label: pizza
    taints:
      key: feaster
      value: "true"
      effect: NoSchedule
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

For both managed and unmanaged nodegroups, this can be done automatically by setting `propagateASGTags` to `true`, which will add the labels and taints as tags to the Auto Scaling group:

```yaml
nodeGroups:
  - name: ng1-public
    ...
    labels:
      my-cool-label: pizza
    taints:
      feaster: "true:NoSchedule"
    propagateASGTags: true
```

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
