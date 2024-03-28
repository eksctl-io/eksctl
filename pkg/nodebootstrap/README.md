# Nodebootstrap

Until recently, `eksctl` provided its own bootstrapping/userdata logic on disk,
over-writing or ignore those which came with the AMIs. This caused numerous headaches
after changes upstream and we got tired of maintaining these extra pieces.

## Implementation

For Unmanaged nodes there is an interface:

```go
type Bootstrapper interface {
  UserData() (string, error)
}
```

AMI families Ubuntu, AmazonLinux2, Bottlerocket and Windows all fulfil this.

As of `eksctl` version `0.45.0` unmanaged nodes of these families, as well as managed nodes (different interface),
will defer to the native bootstrap script which comes built into the image.

This script is found on disk at `/etc/eks/bootstrap.sh`. `UserData` will provide a wrapper script
which will set custom values and delegate to the official bootstrap script.

### Ubuntu & AmazonLinux2

The bootstrapping "prep" for these 2 are fairly similar. Common setup lives in `userdata.go`.
Individual scripts are prepped in `ubuntu.go` and `al2.go`.

Non-dynamic assets live in `assets/`.

Both bootstrappers add `assets/bootstrap.helper.sh` to the node along with either `assets/bootstrap.ubuntu.sh`
or `assets/bootstrap.al2.sh`.

The call to `UserData` will also dynamically add the following:
- `kubelet-extra.json` - user configuration for kubelet
- `docker-extra.json` - extra config for docker daemon
- `kubelet.env` - env vars for kubelet

The bootstrap wrapper scripts will use `jq` and `sed` to get user and our config into various files,
and then call `/etc/eks/bootstrap.sh`.

For AL2, enabling either SSM or EFA will add `assets/install-ssm.al2.sh` or `assets/efa.al2.sh`.

### AmazonLinux2023

While AL2023 implements the `Bootstrapper` interface, the underlying userdata will be entirely different from other AMI families. Specifically, AL2023 introduces a new node initialization process nodeadm that uses a YAML configuration schema, dropping the use of `/etc/eks/bootstrap.sh` script. For self-managed nodes, and for EKS-managed nodes based on custom AMIs, eksctl will populate userdata in the fashion below:

```
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=//

--//
Content-Type: application/node.eks.aws

apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    apiServerEndpoint: https://XXXX.us-west-2.eks.amazonaws.com
    certificateAuthority: XXXX
    cidr: 10.100.0.0/16
    name: my-cluster
  kubelet:
    config:
      clusterDNS:
      - 10.100.0.10
    flags:
    - --node-labels=alpha.eksctl.io/cluster-name=my-cluster,alpha.eksctl.io/nodegroup-name=my-nodegroup
    - --register-with-taints=special=true:NoSchedule (only for EKS-managed nodes)

--//--

```

For EKS-managed nodes based on native AMIs, the userdata above is fulfilled automatically by the AWS SSM agent. 

## Troubleshooting

### Ubuntu

```sh
sudo snap logs kubelet-eks [-n=all/20]
systemctl status docker.service
```

Files:
```sh
/etc/eks/bootstrap.sh
/var/lib/cloud/scripts/eksctl/bootstrap.ubuntu.sh
/etc/kubernetes/kubelet/kubelet-config.json
/etc/docker/daemon.json
```

### AmazonLinux2

Status:
```sh
systemctl status kubelet
systemctl status docker
```

Logs:
```sh
journalctl -u kubelet.service
```

Files:
```sh
/etc/eks/bootstrap.sh
/var/lib/cloud/scripts/eksctl/bootstrap.al2.sh
/etc/kubernetes/kubelet/kubelet-config.json
/etc/docker/daemon.json
/var/lib/cloud/scripts/eksctl/efa.al2.sh
/var/lib/cloud/scripts/eksctl/install-ssm.sh
```
