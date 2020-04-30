# Creating and managing clusters

## Creating a cluster

Create a simple cluster with the following command:

```sh
eksctl create cluster
```

That will create an EKS cluster in your default region (as specified by your AWS CLI configuration) with one
nodegroup containing 2 m5.large nodes.

After the cluster has been created, the appropriate kubernetes configuration will be added to your kubeconfig file.
This is, the file that you have configured in the environment variable `KUBECONFIG` or `~/.kube/config` by default.
The path to the kubeconfig file can be overridden using the `--kubeconfig` flag.

Other flags that can change how the kubeconfig file is written:

| flag                     | type   | use                                                                                                             | default value                 |
|--------------------------|--------|-----------------------------------------------------------------------------------------------------------------|-------------------------------|
| --kubeconfig             | string | path to write kubeconfig (incompatible with –auto-kubeconfig)                                                   | $KUBECONFIG or ~/.kube/config |
| --set-kubeconfig-context | bool   | if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten | true                          |
| --auto-kubeconfig        | bool   | save kubeconfig file by cluster name                                                                            | true                          |
| --write-kubeconfig       | bool   | toggle writing of kubeconfig                                                                                    | true                          |

## Using Config Files

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

Next, run this command:

```
eksctl create cluster -f cluster.yaml
```

This will create a cluster as described.

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

To delete this cluster, run:

```
eksctl delete cluster -f cluster.yaml
```

!!!note

    Without the `--wait` flag, this will only issue a delete operation to the cluster's CloudFormation stack and won't wait for its deletion.

    In some cases, AWS resources using the cluster or its VPC may cause cluster deletion to fail. To ensure any deletion errors are propagated in `eksctl delete cluster`, the `--wait` flag must be used.
    If your delete fails or you forget the wait flag, you may have to go to the CloudFormation GUI and delete the eks stacks from there.

See [`examples/`](https://github.com/weaveworks/eksctl/tree/master/examples) directory for more sample config files.
