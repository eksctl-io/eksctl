# Nodegroups

## Creating nodegroups
You can add one or more nodegroups in addition to the initial nodegroup created along with the cluster.

To create an additional nodegroup, use:

```
eksctl create nodegroup --cluster=<clusterName> [--name=<nodegroupName>]
```
 
???+ note
    `--version` flag is not supported for managed nodegroups. It always inherits the version from control plane.

    By default, new unmanaged nodegroups inherit the version from the control plane (`--version=auto`), but you can specify a different
    version e.g. `--version=1.10`, you can also use `--version=latest` to force use of whichever is the latest version.

Additionally, you can use the same config file used for `eksctl create cluster`:

```
eksctl create nodegroup --config-file=<path>
```

### Creating a nodegroup from a config file

Nodegroups can also be created through a cluster definition or config file. Given the following example config file
and an existing cluster called `dev-cluster`:

```yaml
# dev-cluster.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev-cluster
  region: eu-north-1

managedNodeGroups:
  - name: ng-1-workers
    labels: { role: workers }
    instanceType: m5.xlarge
    desiredCapacity: 10
    volumeSize: 80
    privateNetworking: true
  - name: ng-2-builders
    labels: { role: builders }
    instanceType: m5.2xlarge
    desiredCapacity: 2
    volumeSize: 100
    privateNetworking: true
```

The nodegroups `ng-1-workers` and `ng-2-builders` can be created with this command:

```bash
eksctl create nodegroup --config-file=dev-cluster.yaml
```

#### Load Balancing
If you have already prepared for attaching existing classic load balancers or/and target groups to the nodegroups,
you can specify these in the config file. The classic load balancers or/and target groups are automatically associated with the ASG when creating nodegroups. This is only supported for self-managed nodegroups defined via the `nodeGroups` field.

```yaml
# dev-cluster-with-lb.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev-cluster
  region: eu-north-1

nodeGroups:
  - name: ng-1-web
    labels: { role: web }
    instanceType: m5.xlarge
    desiredCapacity: 10
    privateNetworking: true
    classicLoadBalancerNames:
      - dev-clb-1
      - dev-clb-2
    asgMetricsCollection:
      - granularity: 1Minute
        metrics:
          - GroupMinSize
          - GroupMaxSize
          - GroupDesiredCapacity
          - GroupInServiceInstances
          - GroupPendingInstances
          - GroupStandbyInstances
          - GroupTerminatingInstances
          - GroupTotalInstances
  - name: ng-2-api
    labels: { role: api }
    instanceType: m5.2xlarge
    desiredCapacity: 2
    privateNetworking: true
    targetGroupARNs:
      - arn:aws:elasticloadbalancing:eu-north-1:01234567890:targetgroup/dev-target-group-1/abcdef0123456789
```

## Nodegroup selection in config files

To perform a `create` or `delete` operation on only a subset of the nodegroups specified in a config file, there are two
CLI flags that accept a list of globs, `--include=<glob,glob,...>` and `--exclude=<glob,glob,...>`, e.g.: 

```
eksctl create nodegroup --config-file=<path> --include='ng-prod-*-??' --exclude='ng-test-1-ml-a,ng-test-2-?'
```

Using the example config file above, one can create all the workers nodegroup except the workers one with the following
command:

```bash
eksctl create nodegroup --config-file=dev-cluster.yaml --exclude=ng-1-workers
```

Or one could delete the builders nodegroup with:

```bash
eksctl delete nodegroup --config-file=dev-cluster.yaml --include=ng-2-builders --approve
```

In this case, we also need to supply the `--approve` command to actually delete the nodegroup.


### Include and exclude rules

- if no `--include` or `--exclude` is specified everything is included
- if only `--include` is specified, only nodegroups that match those globs will be included
- if only `--exclude` is specified, all nodegroups that do not match those globs are included
- if both are specified then `--exclude` rules take precedence over `--include` (i.e. nodegroups that match rules in
both groups will be excluded)

## Listing nodegroups

To list the details about a nodegroup or all of the nodegroups, use:

```bash
eksctl get nodegroup --cluster=<clusterName> [--name=<nodegroupName>]
```

To list one or more nodegroups in YAML or JSON format, which outputs more info than the default log table, use:
```bash
# YAML format
eksctl get nodegroup --cluster=<clusterName> [--name=<nodegroupName>] --output=yaml

# JSON format
eksctl get nodegroup --cluster=<clusterName> [--name=<nodegroupName>] --output=json
```

## Nodegroup immutability

By design, nodegroups are immutable. This means that if you need to change something (other than scaling) like the
AMI or the instance type of a nodegroup, you would need to create a new nodegroup with the desired changes, move the
load and delete the old one. See the [Deleting and draining nodegroups](#deleting-and-draining-nodegroups) section.

## Scaling nodegroups

Nodegroup scaling is a process that can take up to a few minutes. When the `--wait` flag is not specified,
`eksctl` optimistically expects the nodegroup to be scaled and returns as soon as the AWS API request has been sent. To make
`eksctl` wait until the nodes are available, add a `--wait` flag like the example below.

???+ note
    Scaling a nodegroup down/in (i.e. reducing the number of nodes) may result in errors as we rely purely on changes to the ASG. This means that the node(s) being removed/terminated aren't explicitly drained. This may be an area for improvement in the future.

Scaling a managed nodegroup is achieved by directly calling the EKS API that updates a managed node group configuration. 

### Scaling a single nodegroup

A nodegroup can be scaled by using the `eksctl scale nodegroup` command:

```
eksctl scale nodegroup --cluster=<clusterName> --nodes=<desiredCount> --name=<nodegroupName> [ --nodes-min=<minSize> ] [ --nodes-max=<maxSize> ] --wait
```

For example, to scale nodegroup `ng-a345f4e1` in `cluster-1` to 5 nodes, run:

```
eksctl scale nodegroup --cluster=cluster-1 --nodes=5 ng-a345f4e1
```

A nodegroup can also be scaled by using a config file passed to `--config-file` and specifying the name of the nodegroup that should be scaled with `--name`. Eksctl will search the config file and discover that nodegroup as well as its scaling configuration values.

If the desired number of nodes is `NOT` within the range of current minimum and current maximum number nodes, one specific error will be shown.
These values can also be passed with flags `--nodes-min` and `--nodes-max` respectively.

### Scaling multiple nodegroups

Eksctl can discover and scale all the nodegroups found in a config file that is passed with `--config-file`. 

Similarly to scaling a single nodegroup, the same set of validations apply to each nodegroup. For example, the desired number of nodes must be within the range of the minimum and maximum number of nodes.

## Deleting and draining nodegroups

To delete a nodegroup, run:

```
eksctl delete nodegroup --cluster=<clusterName> --name=<nodegroupName>
```

[Include and exclude rules](#include-and-exclude-rules) can also be used with this command.


???+ note
    This will drain all pods from that nodegroup before the instances are deleted.

To skip eviction rules during the drain process, run:

```
eksctl delete nodegroup --cluster=<clusterName> --name=<nodegroupName> --disable-eviction
```

All nodes are cordoned and all pods are evicted from a nodegroup on deletion,
but if you need to drain a nodegroup without deleting it, run:

```
eksctl drain nodegroup --cluster=<clusterName> --name=<nodegroupName>
```

To uncordon a nodegroup, run:

```
eksctl drain nodegroup --cluster=<clusterName> --name=<nodegroupName> --undo
```

To ignore eviction rules such as PodDisruptionBudget settings, run:

```
eksctl drain nodegroup --cluster=<clusterName> --name=<nodegroupName> --disable-eviction
```

To speed up the drain process you can specify `--parallel <value>` for the number of nodes to drain in parallel.

## Other features
You can also enable SSH, ASG access and other features for a nodegroup, e.g.:

```
eksctl create nodegroup --cluster=cluster-1 --node-labels="autoscaling=enabled,purpose=ci-worker" --asg-access --full-ecr-access --ssh-access
```

### Update labels

There are no specific commands in `eksctl` to update the labels of a nodegroup, but it can easily be achieved using
`kubectl`, e.g.:

```bash
kubectl label nodes -l alpha.eksctl.io/nodegroup-name=ng-1 new-label=foo
```

### SSH Access
You can enable SSH access for nodegroups by configuring one of `publicKey`, `publicKeyName` and `publicKeyPath` in your
nodegroup configuration. Alternatively you can use [AWS Systems Manager (SSM)](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-sessions-start.html#sessions-start-cli) to SSH onto nodes, by configuring the nodegroup with `enableSsm`:


```yaml
managedNodeGroups:
  - name: ng-1
    instanceType: m5.large
    desiredCapacity: 1
    ssh: # import public key from file
      publicKeyPath: ~/.ssh/id_rsa_tests.pub
  - name: ng-2
    instanceType: m5.large
    desiredCapacity: 1
    ssh: # use existing EC2 key
      publicKeyName: ec2_dev_key
  - name: ng-3
    instanceType: m5.large
    desiredCapacity: 1
    ssh: # import inline public key
      publicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDqZEdzvHnK/GVP8nLngRHu/GDi/3PeES7+Bx6l3koXn/Oi/UmM9/jcW5XGziZ/oe1cPJ777eZV7muEvXg5ZMQBrYxUtYCdvd8Rt6DIoSqDLsIPqbuuNlQoBHq/PU2IjpWnp/wrJQXMk94IIrGjY8QHfCnpuMENCucVaifgAhwyeyuO5KiqUmD8E0RmcsotHKBV9X8H5eqLXd8zMQaPl+Ub7j5PG+9KftQu0F/QhdFvpSLsHaxvBzA5nhIltjkaFcwGQnD1rpCM3+UnQE7Izoa5Yt1xoUWRwnF+L2TKovW7+bYQ1kxsuuiX149jXTCJDVjkYCqi7HkrXYqcC1sbsror someuser@hostname"
  - name: ng-4
    instanceType: m5.large
    desiredCapacity: 1
    ssh: # enable SSH using SSM
      enableSsm: true
```
