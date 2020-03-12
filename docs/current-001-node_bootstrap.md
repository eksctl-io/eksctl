# Current Design #001: Node Bootstrap

> NOTE: the purpose of this document is to summarize how node bootstrap process is currently implemented in `eksctl`,
> so it can be referred to in any discussions.

## Guiding Principles

To enable reliable node management at large-scale, the bootstrap process of a node should bear the following attributes:

- stateless
- immutable
- deterministic
- simple

It **should avoid**:
- non-essential runtime parameters
- mutation of configuration during execution
- generating complex configuration files
- download of configuration files
- installation of software
- making unnecessary look-up calls to remote APIs
- only call out to commands when absolutely essential

Current implementation is build for official AL2 EKS AMIs, that includes dependencies such as AWS CLI and jq, however
it's ought to be possible for one to create a stripped-down AMI and re-use the bootstrap script.

Aside from all of the above concerns, it's crucial to bear in mind that the dependencies of bootstrap process should
not pose requirements that affect running of the node or the entire cluster past the bootstrap stage. For example, relying
on particular AWS APIs implies instance role has to be adjusted to allow for access to such APIs, which could otherwise
be avoided, or any software required solely for bootstrap function has to be tracked for security vulnerabilities.

## Known Anti-patterns

### Use of Dependencies

If bootstrap process requires use of a dependency, such as AWS CLI and jq, it means this dependency has to be installed
on the AMI, and thereby it affects running of the node beyond the bootstrap stage. If a vulnerability is found, e.g. in
one of the Python packages that AWS CLI uses, the node will have to be upgraded.

### Use of EC2 tags for parameter discovery

This requires multiple remote calls, and parsing JSON. It can be done with AWS CLI, yet it's a Python dependency
that is not essential otherwise. Furthermore, discovering tags requires listing tags of all instances and searching
the list. It's not as simple as asking "what are the tags of this instance?", it's "what is the ID of this instance?",
followed by "what are the tags of all instances? which of all instances is this instance? what are the tags of this
instance?".

Aside from this, this approach requires instance IAM role to have access to EC2 APIs, which is considered non-essential
for any other purpose.

## Current Design

Summary:
- all configuration files along with [the bootstrap script](https://github.com/weaveworks/eksctl/blob/70041a226bb8ef5c51a229d587235551a2410eda/pkg/nodebootstrap/assets/bootstrap.al2.sh) are passed to the node as cloud-init config
- all configuration files are written to `/etc/eksctl` directory
- the only remote API that's used is the EC2 metadata service, required to obtain instance ID, type and IP address
- the behavior is fully determined by:
    - `ClusterConfig` set in configuration file
    - version of `eksctl`
- single mode of execution

More specifically, as part of CloudFormation template, `eksctl` will define full content of the following files:

- `/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf` - systemd drop-in for kubelet
- `/etc/eksctl/kubelet.yaml` - fully-formed kubelet configuration file (NOTE: currently it doesn't support all the options that flags support)
- `/etc/eksctl/kubeconfig.yaml` - client credentials for kubelet
- `/etc/eksctl/ca.crt`
- `/etc/eksctl/metadata.env` - known metadata for all nodes in a nodegroup, used by the bootstrap process and for kubelet flags that are not otherwise settable via `/etc/eksctl/kubelet.yaml`
        - `AWS_DEFAULT_REGION`
        - `AWS_EKS_CLUSTER_NAME`
        - `AWS_EKS_ENDPOINT`
        - `AWS_EKS_ECR_ACCOUNT`
- `/etc/eksctl/kubelet.env` - kubelet-specific metadata, used for flags that are not otherwise settable via `/etc/eksctl/kubelet.yaml`
        - `NODE_LABELS`
        - `NODE_TAINTS`
        - `MAX_PODS` (optional when user chose to override it)
- `/etc/eksctl/max_pods.map` - used by the bootstrap script (NOTE: based on [`eni-max-pods.txt`](https://raw.github.com/awslabs/amazon-eks-ami/master/files/eni-max-pods.txt), but format differs to ease parsing)

## Source Code

### [Bootstrap Script for Amazon Linux 2](https://github.com/weaveworks/eksctl/blob/70041a226bb8ef5c51a229d587235551a2410eda/pkg/nodebootstrap/assets/bootstrap.al2.sh)

```shell
#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

function get_max_pods() {
  while read instance_type pods; do
    if  [[ "${instance_type}" == "${1}" ]] && [[ "${pods}" =~ ^[0-9]+$ ]] ; then
      echo ${pods}
      return
    fi
  done < /etc/eksctl/max_pods.map
}

NODE_IP="$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)"
INSTANCE_ID="$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)"
INSTANCE_TYPE="$(curl --silent http://169.254.169.254/latest/meta-data/instance-type)"

source /etc/eksctl/kubelet.env # this can override MAX_PODS

cat > /etc/eksctl/kubelet.local.env <<EOF
NODE_IP=${NODE_IP}
INSTANCE_ID=${INSTANCE_ID}
INSTANCE_TYPE=${INSTANCE_TYPE}
MAX_PODS=${MAX_PODS:-$(get_max_pods "${INSTANCE_TYPE}")}
EOF

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
```

Dependencies:
- systemd
- cat
- bash
- curl

### [Bootstrap Script for Ubuntu](https://github.com/weaveworks/eksctl/blob/70041a226bb8ef5c51a229d587235551a2410eda/pkg/nodebootstrap/assets/bootstrap.ubuntu.sh)

Same as above, but Ubuntu-specific commands are used instead of `systemctl` and the drop-in unit:
```shell
### 28 lines omitted

snap alias kubelet-eks.kubelet kubelet
snap alias kubectl-eks.kubectl kubectl
snap stop kubelet-eks
systemctl reset-failed

(
  # TODO: these should be looked at every time kubelet starts up,
  # which is what we do in AL2 (which is based on plain systemd,
  # and meant to be portable to most systemd distros), but it's
  # not clear how to load these from kubelet snap without having
  # to customise the snap itself
  source /etc/eksctl/kubelet.local.env
  source /etc/eksctl/kubelet.env
  source /etc/eksctl/metadata.env

  flags=(
    "node-ip=${NODE_IP}"
    "max-pods=${MAX_PODS}"
    "node-labels=${NODE_LABELS},alpha.eksctl.io/instance-id=${INSTANCE_ID}"
    "pod-infra-container-image=${AWS_EKS_ECR_ACCOUNT}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/eks/pause-amd64:3.1"
    "cloud-provider=aws"
    "cni-bin-dir=/opt/cni/bin"
    "cni-conf-dir=/etc/cni/net.d"
    "container-runtime=docker"
    "network-plugin=cni"
    "register-node=true"
    "register-with-taints=${NODE_TAINTS}"
    "kubeconfig=/etc/eksctl/kubeconfig.yaml"
    "config=/etc/eksctl/kubelet.yaml"
  )

  snap set kubelet-eks "${flags[@]}"
)

snap start kubelet-eks
```

### [Systemd Drop-in Unit for Amazon Linux](https://github.com/weaveworks/eksctl/blob/70041a226bb8ef5c51a229d587235551a2410eda/pkg/nodebootstrap/assets/10-eksclt.al2.conf)

```ini
# eksctl-specific systemd drop-in unit for kubelet, for Amazon Linux 2 (AL2)

[Service]
# Local metadata parameters: REGION, AWS_DEFAULT_REGION
EnvironmentFile=/etc/eksctl/metadata.env
# Global and static parameters: CLUSTER_DNS, NODE_LABELS, NODE_TAINTS
EnvironmentFile=/etc/eksctl/kubelet.env
# Local non-static parameters: NODE_IP, INSTANCE_ID
EnvironmentFile=/etc/eksctl/kubelet.local.env

ExecStart=/usr/bin/kubelet \
  --node-ip=${NODE_IP} \
  --node-labels=${NODE_LABELS},alpha.eksctl.io/instance-id=${INSTANCE_ID} \
  --max-pods=${MAX_PODS} \
  --register-node=true --register-with-taints=${NODE_TAINTS} \
  --cloud-provider=aws \
  --container-runtime=docker \
  --network-plugin=cni \
  --cni-bin-dir=/opt/cni/bin \
  --cni-conf-dir=/etc/cni/net.d \
  --pod-infra-container-image=${AWS_EKS_ECR_ACCOUNT}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/eks/pause-amd64:3.1 \
  --kubeconfig=/etc/eksctl/kubeconfig.yaml \
  --config=/etc/eksctl/kubelet.yaml
```

### [`kubelet.yaml`](https://github.com/weaveworks/eksctl/blob/70041a226bb8ef5c51a229d587235551a2410eda/pkg/nodebootstrap/assets/kubelet.yaml)

```YAML
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1

address: 0.0.0.0
clusterDomain: cluster.local

authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 2m0s
    enabled: true
  x509:
    clientCAFile: /etc/eksctl/ca.crt

authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 5m0s
    cacheUnauthorizedTTL: 30s

serverTLSBootstrap: true

cgroupDriver: cgroupfs

featureGates:
  RotateKubeletServerCertificate: true
```

## Comparison to EKS Amazon Linux 2 script

The [EKS Amazon Linux 2 script](https://github.com/awslabs/amazon-eks-ami/blob/b85ef2fce5f46eddafa2e8881d757b31a9feed81/files/bootstrap.sh):
- has 228 lines of bash code
- has 9 flags representing different modes of execution
- uses `sed` to mutate a YAML file (`/var/lib/kubelet/kubeconfig`)
- uses `jq` to generate JSON files (`/etc/docker/daemon.json` and `/etc/kubernetes/kubelet/kubelet-config.json`)
- calls `aws eks wait cluster-active` and `aws eks describe-cluster`
- implements API call retry logic
- uses `awk` and `grep` to parse structured data obtained from EC2 metadata API

Dependencies:
- systemd
- cat
- bash
- curl
- awk
- sed
- jq
- aws

## Comparison to Ubuntu script

No source for the script used by Ubuntu is publicly available. It can be obtained from the AMI, but it is not
clear if copyright may permit us from sharing the script here.
It is known to be very different from the AL2 script, but similar in complexity and requires access to EC2 APIs
for tag-based discovery (discussed above).
