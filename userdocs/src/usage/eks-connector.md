# Registering non-EKS clusters with EKS Connector
The EKS Console provides a single pane of glass to manage all your Kubernetes clusters, including those hosted on
other cloud providers, via [EKS Connector][eks-user-guide]. This process requires registering the cluster with EKS and running the
EKS Connector agent on the external Kubernetes cluster.

`eksctl` simplifies registering non-EKS clusters by creating the required AWS resources and generating Kubernetes manifests
for EKS Connector to apply to the external cluster.


## Register Cluster
To register or connect a non-EKS Kubernetes cluster, run

```shell
$ eksctl register cluster --name <name> --provider <provider>
2021-08-19 13:47:26 [ℹ]  creating IAM role "eksctl-20210819194112186040"
2021-08-19 13:47:26 [ℹ]  registered cluster "<name>" successfully
2021-08-19 13:47:26 [ℹ]  wrote file eks-connector.yaml to <current directory>
2021-08-19 13:47:26 [ℹ]  wrote file eks-connector-clusterrole.yaml to <current directory>
2021-08-19 13:47:26 [ℹ]  wrote file eks-connector-console-dashboard-full-access-group.yaml to <current directory>
2021-08-19 13:47:26 [!]  note: "eks-connector-clusterrole.yaml" and "eks-connector-console-dashboard-full-access-group.yaml" give full EKS Console access to IAM identity "<aws-arn>", edit if required; read https://eksctl.io/usage/eks-connector for more info
2021-08-19 13:47:26 [ℹ]  run `kubectl apply -f eks-connector.yaml,eks-connector-clusterrole.yaml,eks-connector-console-dashboard-full-access-group.yaml` before <expiry> to connect the cluster

```

This command will register the cluster and write three files that contain the Kubernetes manifests for EKS Connector
that must be applied to the external cluster before the registration expires.

!!!note
`eks-connector-clusterrole.yaml` and `eks-connector-console-dashboard-full-access-clusterrole.yaml` give `get` and `list` permissions for Kubernetes resources
in all namespaces to the calling IAM identity and must be edited accordingly if required before applying them to the cluster.
To configure more restricted access, see [Granting access to a user to view a cluster](https://docs.aws.amazon.com/eks/latest/userguide/connector-grant-access.html).

To provide an existing IAM role to use for EKS Connector, pass it via `--role-arn` as in:

```shell
$ eksctl register cluster --name <name> --provider <provider> --role-arn=<role-arn>
```


If the cluster already exists, eksctl will return an error.


## Deregister cluster

To deregister or disconnect a registered cluster, run

```shell
$ eksctl deregister cluster --name <name>
2021-08-19 16:04:09 [ℹ]  unregistered cluster "<name>" successfully
2021-08-19 16:04:09 [ℹ]  run `kubectl delete namespace eks-connector` and `kubectl delete -f eks-connector-binding.yaml` on your cluster to remove EKS Connector resources
```

This command will deregister the external cluster and remove its associated AWS resources, but you are required to remove the
EKS connector Kubernetes resources from the cluster.


## Further information

- [EKS Connector][eks-user-guide]

[eks-user-guide]: https://docs.aws.amazon.com/eks/latest/userguide/eks-connector.html
