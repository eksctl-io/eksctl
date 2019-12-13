---
title: "EKS Fargate Support"
weight: 160
url: usage/fargate
---

## Fargate

[AWS Fargate][fargate] is a managed compute engine for Amazon ECS that can run containers. In Fargate you don't need to
manage servers or clusters.

[Amazon EKS can now launch pods onto AWS Fargate][eks-fargate]. This removes the need to worry about how you provision or manage
infrastructure for pods and makes it easier to build and run performant, highly-available Kubernetes applications on AWS.

### Creating a cluster with Fargate support

You can add a cluster with Fargate support with:

```console
$ eksctl create cluster --fargate
[ℹ]  eksctl version 0.11.0
[ℹ]  using region ap-northeast-1
[ℹ]  setting availability zones to [ap-northeast-1a ap-northeast-1d ap-northeast-1c]
[ℹ]  subnets for ap-northeast-1a - public:192.168.0.0/19 private:192.168.96.0/19
[ℹ]  subnets for ap-northeast-1d - public:192.168.32.0/19 private:192.168.128.0/19
[ℹ]  subnets for ap-northeast-1c - public:192.168.64.0/19 private:192.168.160.0/19
[ℹ]  nodegroup "ng-dba9d731" will use "ami-02e124a380df41614" [AmazonLinux2/1.14]
[ℹ]  using Kubernetes version 1.14
[ℹ]  creating EKS cluster "ridiculous-painting-1574859263" in "ap-northeast-1" region
[ℹ]  will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup
[ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=ap-northeast-1 --cluster=ridiculous-painting-1574859263'
[ℹ]  CloudWatch logging will not be enabled for cluster "ridiculous-painting-1574859263" in "ap-northeast-1"
[ℹ]  you can enable it with 'eksctl utils update-cluster-logging --region=ap-northeast-1 --cluster=ridiculous-painting-1574859263'
[ℹ]  Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster "ridiculous-painting-1574859263" in "ap-northeast-1"
[ℹ]  2 sequential tasks: { create cluster control plane "ridiculous-painting-1574859263", create nodegroup "ng-dba9d731" }
[ℹ]  building cluster stack "eksctl-ridiculous-painting-1574859263-cluster"
[ℹ]  deploying stack "eksctl-ridiculous-painting-1574859263-cluster"
[ℹ]  building nodegroup stack "eksctl-ridiculous-painting-1574859263-nodegroup-ng-dba9d731"
[ℹ]  --nodes-min=2 was set automatically for nodegroup ng-dba9d731
[ℹ]  --nodes-max=2 was set automatically for nodegroup ng-dba9d731
[ℹ]  deploying stack "eksctl-ridiculous-painting-1574859263-nodegroup-ng-dba9d731"
[✔]  all EKS cluster resources for "ridiculous-painting-1574859263" have been created
[✔]  saved kubeconfig as "/Users/marc/.kube/config"
[ℹ]  adding identity "arn:aws:iam::123456789012:role/eksctl-ridiculous-painting-157485-NodeInstanceRole-104DXUJOFDPO5" to auth ConfigMap
[ℹ]  nodegroup "ng-dba9d731" has 0 node(s)
[ℹ]  waiting for at least 2 node(s) to become ready in "ng-dba9d731"
[ℹ]  nodegroup "ng-dba9d731" has 2 node(s)
[ℹ]  node "ip-192-168-27-156.ap-northeast-1.compute.internal" is ready
[ℹ]  node "ip-192-168-95-177.ap-northeast-1.compute.internal" is ready
[ℹ]  creating Fargate profile "default" on EKS cluster "ridiculous-painting-1574859263"
[ℹ]  created Fargate profile "default" on EKS cluster "ridiculous-painting-1574859263"
[ℹ]  kubectl command should work with "/Users/marc/.kube/config", try 'kubectl get nodes'
[✔]  EKS cluster "ridiculous-painting-1574859263" in "ap-northeast-1" region is ready
```

This command will have created a cluster and a Fargate profile. This profile contains certain information needed by AWS to instantiate
pods in Fargate. These are:

- pod execution role to define the permissions required to run the pod and the
  networking location (subnet) to run the pod. This allows the same networking
  and security permissions to be applied to multiple Fargate pods and makes it
  easier to migrate existing pods on a cluster to Fargate.
- Selector to define which pods should run on Fargate. This is composed by a
  `namespace` and `labels`.

When the profile is not specified but support for Fargate is enabled with `--fargate` a default Fargate profile is
created. This profile targets the `default` and the `kube-system` namespaces so pods in those namespaces will run on
Fargate.

The Fargate profile that was created can be checked with the following command:

```console
$ eksctl get fargateprofile --cluster ridiculous-painting-1574859263 -o yaml
- name: default
  podExecutionRoleARN: arn:aws:iam::123456789012:role/eksctl-ridiculous-painting-1574859263-ServiceRole-EIFQOH0S1GE7
  selectors:
  - namespace: default
  - namespace: kube-system
  subnets:
  - subnet-0b3a5522f3b48a742
  - subnet-0c35f1497067363f3
  - subnet-0a29aa00b25082021
```

To learn more about selectors see [Designing Fargate profiles](#designing-fargate-profiles).

### Creating a cluster with Fargate support using a config file

The following config file declares an EKS cluster with both a nodegroup composed of one EC2 `m5.large` instance and two
Fargate profiles. All pods defined in the `default` and `kube-system` namespaces will run on Fargate. All pods in the
`dev` namespace that also have the label `dev=passed` will also run on Fargate. Any other pods will be scheduled on the
node in `ng-1`.

```yaml
# An example of ClusterConfig with a normal nodegroup and a Fargate profile.
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: fargate-cluster
  region: ap-northeast-1

nodeGroups:
  - name: ng-1
    instanceType: m5.large
    desiredCapacity: 1

fargateProfiles:
  - name: default
    selectors:
      # All workloads in the "default" Kubernetes namespace will be
      # scheduled onto Fargate:
      - namespace: default
      # All workloads in the "kube-system" Kubernetes namespace will be
      # scheduled onto Fargate:
      - namespace: kube-system
  - name: dev
    selectors:
      # All workloads in the "dev" Kubernetes namespace matching the following
      # label selectors will be scheduled onto Fargate:
      - namespace: dev
        labels:
          env: dev
          checks: passed
```

```console
$ eksctl create cluster -f cluster-fargate.yaml
[ℹ]  eksctl version 0.11.0
[ℹ]  using region ap-northeast-1
[ℹ]  setting availability zones to [ap-northeast-1c ap-northeast-1a ap-northeast-1d]
[ℹ]  subnets for ap-northeast-1c - public:192.168.0.0/19 private:192.168.96.0/19
[ℹ]  subnets for ap-northeast-1a - public:192.168.32.0/19 private:192.168.128.0/19
[ℹ]  subnets for ap-northeast-1d - public:192.168.64.0/19 private:192.168.160.0/19
[ℹ]  nodegroup "ng-1" will use "ami-02e124a380df41614" [AmazonLinux2/1.14]
[ℹ]  using Kubernetes version 1.14
[ℹ]  creating EKS cluster "fargate-cluster" in "ap-northeast-1" region with Fargate profile and un-managed nodes
[ℹ]  1 nodegroup (ng-1) was included (based on the include/exclude rules)
[ℹ]  will create a CloudFormation stack for cluster itself and 1 nodegroup stack(s)
[ℹ]  will create a CloudFormation stack for cluster itself and 0 managed nodegroup stack(s)
[ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=ap-northeast-1 --cluster=fargate-cluster'
[ℹ]  CloudWatch logging will not be enabled for cluster "fargate-cluster" in "ap-northeast-1"
[ℹ]  you can enable it with 'eksctl utils update-cluster-logging --region=ap-northeast-1 --cluster=fargate-cluster'
[ℹ]  Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster "fargate-cluster" in "ap-northeast-1"
[ℹ]  2 sequential tasks: { create cluster control plane "fargate-cluster", create nodegroup "ng-1" }
[ℹ]  building cluster stack "eksctl-fargate-cluster-cluster"
[ℹ]  deploying stack "eksctl-fargate-cluster-cluster"
[ℹ]  building nodegroup stack "eksctl-fargate-cluster-nodegroup-ng-1"
[ℹ]  --nodes-min=1 was set automatically for nodegroup ng-1
[ℹ]  --nodes-max=1 was set automatically for nodegroup ng-1
[ℹ]  deploying stack "eksctl-fargate-cluster-nodegroup-ng-1"
[✔]  all EKS cluster resources for "fargate-cluster" have been created
[✔]  saved kubeconfig as "/home/user1/.kube/config"
[ℹ]  adding identity "arn:aws:iam::123456789012:role/eksctl-fargate-cluster-nod-NodeInstanceRole-42Q80B2Z147I" to auth ConfigMap
[ℹ]  nodegroup "ng-1" has 0 node(s)
[ℹ]  waiting for at least 1 node(s) to become ready in "ng-1"
[ℹ]  nodegroup "ng-1" has 1 node(s)
[ℹ]  node "ip-192-168-71-83.ap-northeast-1.compute.internal" is ready
[ℹ]  creating Fargate profile "default" on EKS cluster "fargate-cluster"
[ℹ]  created Fargate profile "default" on EKS cluster "fargate-cluster"
[ℹ]  creating Fargate profile "dev" on EKS cluster "fargate-cluster"
[ℹ]  created Fargate profile "dev" on EKS cluster "fargate-cluster"
[ℹ]  "coredns" is now schedulable onto Fargate
[ℹ]  "coredns" is now scheduled onto Fargate
[ℹ]  "coredns" is now scheduled onto Fargate
[ℹ]  "coredns" pods are now scheduled onto Fargate
[ℹ]  kubectl command should work with "/home/user1/.kube/config", try 'kubectl get nodes'
[✔]  EKS cluster "fargate-cluster" in "ap-northeast-1" region is ready
```

### Designing Fargate profiles

Each selector entry has up to two components, namespace and a list of key-value
pairs. Only the namespace component is required to create a selector entry. All rules
(namespaces, key value pairs) must apply to a pod to match a selector entry. A pod
only needs to match one selector entry to run on the profile.
Any pod that matches all the conditions in a selector field would be scheduled to be run on
Fargate. Any pods not matching either the whitelisted Namespaces but where the
user manually set the scheduler: fargate-scheduler filed would be stuck in a Pending
state, as they were not authorized to run on Fargate.

Profiles must meet the following requirements:

- One selector is mandatory per profile
- Each selector must include a namespace; labels are optional

#### Example: scheduling workload in Fargate

To schedule pods on Fargate for the example mentioned above, one could, for example, create a namespace called `dev` and
deploy the workload there:

```console
$ kubectl create namespace dev
namespace/dev created

$ kubectl run nginx --image=nginx --restart=Never --namespace dev
pod/nginx created

$ kubectl get pods --all-namespaces --output wide
NAMESPACE     NAME                       READY   STATUS    AGE   IP                NODE
dev           nginx                      1/1     Running   75s   192.168.183.140   fargate-ip-192-168-183-140.ap-northeast-1.compute.internal
kube-system   aws-node-44qst             1/1     Running   21m   192.168.70.246    ip-192-168-70-246.ap-northeast-1.compute.internal
kube-system   aws-node-4vr66             1/1     Running   21m   192.168.23.122    ip-192-168-23-122.ap-northeast-1.compute.internal
kube-system   coredns-699bb99bf8-84x74   1/1     Running   26m   192.168.2.95      ip-192-168-23-122.ap-northeast-1.compute.internal
kube-system   coredns-699bb99bf8-f6x6n   1/1     Running   26m   192.168.90.73     ip-192-168-70-246.ap-northeast-1.compute.internal
kube-system   kube-proxy-brxhg           1/1     Running   21m   192.168.23.122    ip-192-168-23-122.ap-northeast-1.compute.internal
kube-system   kube-proxy-zd7s8           1/1     Running   21m   192.168.70.246    ip-192-168-70-246.ap-northeast-1.compute.internal
```

From the output of the last `kubectl get pods` command we can see that the `nginx` pod is deployed in a node called
`fargate-ip-192-168-183-140.ap-northeast-1.compute.internal`.

### Managing Fargate profiles

To deploy Kubernetes workloads on Fargate, EKS needs a Fargate profile. When creating a cluster like in the examples
above, `eksctl` takes care of this by creating a default profile. Given an already existing cluster, it's also possible to
create a Fargate profile with the `eksctl create fargateprofile` command:

> NOTE: This operation is only supported on clusters that run on the EKS platform version `eks.5` or higher.
>
> NOTE: If the existing was created with a version of `eksctl` prior to 0.11.0, you will  need to run `eksctl update
> cluster` before creating the Fargate profile.

```console
$ eksctl create fargateprofile --namespace dev --cluster fargate-example-cluster
[ℹ]  creating Fargate profile "fp-9bfc77ad" on EKS cluster "fargate-example-cluster"
[ℹ]  created Fargate profile "fp-9bfc77ad" on EKS cluster "fargate-example-cluster"
```

You can also specify the name of the Fargate profile to be created. This name must not start with the prefix `eks-`.

```console
$ eksctl create fargateprofile --name eks-dev --namespace dev --cluster fargate-example-cluster --name fp-development
[ℹ]  created Fargate profile "fp-development" on EKS cluster "fargate-example-cluster"
```

Using this command with CLI flags eksctl can only create a single Fargate profile with a simple selector. For more
complex selectors, for example with more namespaces, eksctl supports using a config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: fargate-example-cluster
  region: ap-northeast-1

fargateProfiles:
  - name: default
    selectors:
      # All workloads in the "default" Kubernetes namespace will be
      # scheduled onto Fargate:
      - namespace: default
      # All workloads in the "kube-system" Kubernetes namespace will be
      # scheduled onto Fargate:
      - namespace: kube-system
  - name: dev
    selectors:
      # All workloads in the "dev" Kubernetes namespace matching the following
      # label selectors will be scheduled onto Fargate:
      - namespace: dev
        labels:
          env: dev
          checks: passed

```

```console
$ eksctl create fargateprofile -f fargate-example-cluster.yaml
[ℹ]  creating Fargate profile "default" on EKS cluster "fargate-example-cluster"
[ℹ]  created Fargate profile "default" on EKS cluster "fargate-example-cluster"
[ℹ]  creating Fargate profile "dev" on EKS cluster "fargate-example-cluster"
[ℹ]  created Fargate profile "dev" on EKS cluster "fargate-example-cluster"
[ℹ]  "coredns" is now scheduled onto Fargate
[ℹ]  "coredns" pods are now scheduled onto Fargate
```

To see existing Fargate profiles in a cluster:

```console
$ eksctl get fargateprofile --cluster fargate-example-cluster
NAME         POD_EXECUTION_ROLE_ARN                                                                   SUBNETS                                                                     SELECTOR_NAMESPACE  SELECTOR_LABELS
fp-9bfc77ad  arn:aws:iam::123456789012:role/eksctl-fargate-example-cluster-ServiceRole-1T5F78E5FSH79  subnet-00adf1d8c99f83381,subnet-04affb163ffab17d4,subnet-035b34379d5ef5473  dev                 <none>
```

And to see them in `yaml` format:

```console
$ eksctl get fargateprofile --cluster fargate-example-cluster -o yaml
- name: fp-9bfc77ad
  podExecutionRoleARN: arn:aws:iam::123456789012:role/eksctl-fargate-example-cluster-ServiceRole-1T5F78E5FSH79
  selectors:
  - namespace: dev
  subnets:
  - subnet-00adf1d8c99f83381
  - subnet-04affb163ffab17d4
  - subnet-035b34379d5ef5473
```

Or in `json` format:

```console
$ eksctl get fargateprofile --cluster fargate-example-cluster -o json
[
    {
        "name": "fp-9bfc77ad",
        "podExecutionRoleARN": "arn:aws:iam::123456789012:role/eksctl-fargate-example-cluster-ServiceRole-1T5F78E5FSH79",
        "selectors": [
            {
                "namespace": "dev"
            }
        ],
        "subnets": [
            "subnet-00adf1d8c99f83381",
            "subnet-04affb163ffab17d4",
            "subnet-035b34379d5ef5473"
        ]
    }
]
```

Fargate profiles are immutable by design. To change something, create a new Fargate profile with the desired changes and
delete the old one with the `eksctl delete fargateprofile` command like in the following example:

```console
$ eksctl delete fargateprofile --cluster fargate-example-cluster --name fp-9bfc77ad --wait
2019-11-27T19:04:26+09:00 [ℹ]  deleting Fargate profile "fp-9bfc77ad"
  ClusterName: "fargate-example-cluster",
  FargateProfileName: "fp-9bfc77ad"
}
```

Note that the profile deletion is a process that can take up to a few minutes. When the `--wait` flag is not specified,
`eksctl` optimistically expects the profile to be deleted and returns as soon as the aws request has been sent. To make
`eksctl` wait until it has been successfully deleted use `--wait` like in the example above.

### Further reading

- [Fargate][fargate]
- [Fargate from EKS][eks-fargate]

[fargate]: https://aws.amazon.com/fargate/
[eks-fargate]: https://docs.aws.amazon.com/eks/latest/userguide/fargate.html
