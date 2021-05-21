# THIS DESIGN DEPRECATED AS OF 0.47.0

See [pkg/nodebootstrap/README.md](/pkg/nodebootstrap/README.md) for current implementation details.

--------------

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

- [`/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf`](#systemd-unit-for-amazon-linux) - systemd drop-in for kubelet
- [`/etc/eksctl/kubelet.yaml`](#kubelet-configuration) - fully-formed kubelet configuration file (NOTE: currently it doesn't support all the options that flags support)
- [`/etc/eksctl/kubeconfig.yaml`](#kubeconfig) - client credentials for kubelet
- [`/etc/eksctl/ca.crt`](#ca-file)
- [`/etc/eksctl/metadata.env`](#metadata-environment-variables) - known metadata for all nodes in a nodegroup, used by the bootstrap process and for kubelet flags that are not otherwise settable via `/etc/eksctl/kubelet.yaml`
- [`/etc/eksctl/kubelet.env`](#kubelet-environment-variables) - kubelet-specific metadata, used for flags that are not otherwise settable via `/etc/eksctl/kubelet.yaml`
- [`/etc/eksctl/kubelet.local.env`](#kubelet-local-environment-variables) - kubelet-specific metadata that is local to the specific node
- [`/etc/eksctl/max_pods.map`](#max-pods-file) - used by the bootstrap script (NOTE: based on [`eni-max-pods.txt`](https://raw.github.com/awslabs/amazon-eks-ami/master/files/eni-max-pods.txt), but format differs to ease parsing)

## Source Code

The bootstrapping of un-managed nodegroups is done through CloudFormation with the files that eksctl sends as `userData`
(the user data is a property of the launch template used to spawn the nodes in EC2).

The main files are the [bootstrap script](#bootstrap-script-for-amazon-linux-2), the [systemd unit](#systemd-unit-for-amazon-linux), the
[kubelet.yaml](#kubelet-configuration) and the [kubeconfig.yalm](#kubeconfig).


### [Bootstrap Script for Amazon Linux 2](https://github.com/weaveworks/eksctl/blob/master/pkg/nodebootstrap/assets/bootstrap.al2.sh)

This is sent as one of the `runScript`s  in the `userData` and it is used to bootstrap the node. It can also be
[overwritten](/usage/schema/#nodeGroups-overrideBootstrapCommand) by the user or appended with
[extra shell commands](/usage/schema/#nodeGroups-preBootstrapCommands).

Dependencies:
- systemd
- cat
- bash
- curl

<pre>
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
  done < <a href="#max-pods-file">/etc/eksctl/max_pods.map</a>
}

# Use IMDSv2 to get metadata
TOKEN="$(curl --silent -X PUT -H "X-aws-ec2-metadata-token-ttl-seconds: 600" http://169.254.169.254/latest/api/token)"
function get_metadata() {
  curl --silent -H "X-aws-ec2-metadata-token: $TOKEN" "http://169.254.169.254/latest/meta-data/$1"
}

NODE_IP="$(get_metadata local-ipv4)"
INSTANCE_ID="$(get_metadata instance-id)"
INSTANCE_TYPE="$(get_metadata instance-type)"
AWS_SERVICES_DOMAIN="$(get_metadata services/domain)"


source <a href="#kubelet-environment-variables">/etc/eksctl/kubelet.env</a> # this can override MAX_PODS

cat &gt; <a href="#kubelet-local-environment-variables">/etc/eksctl/kubelet.local.env</a>  &lt;&lt;EOF
NODE_IP=${NODE_IP}
INSTANCE_ID=${INSTANCE_ID}
INSTANCE_TYPE=${INSTANCE_TYPE}
AWS_SERVICES_DOMAIN=${AWS_SERVICES_DOMAIN}
MAX_PODS=${MAX_PODS:-$(get_max_pods "${INSTANCE_TYPE}")}
EOF

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet

</pre>


### [Bootstrap Script for Ubuntu](https://github.com/weaveworks/eksctl/blob/master/pkg/nodebootstrap/assets/bootstrap.ubuntu.sh)

Same as above, but Ubuntu-specific commands are used instead of `systemctl` and the drop-in unit.

### [Systemd Unit for Amazon Linux](https://github.com/weaveworks/eksctl/blob/70041a226bb8ef5c51a229d587235551a2410eda/pkg/nodebootstrap/assets/10-eksclt.al2.conf)

This is the systemd unit file used to start the kubelet service. It makes use of other configuration files sent by
eksctl and it is stored in `etc/systemd/system/kubelet.service.d/10-eksctl.al2.conf`.

<pre>
# eksctl-specific systemd drop-in unit for kubelet, for Amazon Linux 2 (AL2)

[Service]
EnvironmentFile=<a href="#metadata-environment-variables">/etc/eksctl/metadata.env</a>
EnvironmentFile=<a href="#kubelet-environment-variables">/etc/eksctl/kubelet.env</a>
EnvironmentFile=<a href="#kubelet-local-environment-variables">/etc/eksctl/kubelet.local.env</a>

ExecStart=
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
  --pod-infra-container-image=${AWS_EKS_ECR_ACCOUNT}.dkr.ecr.${AWS_DEFAULT_REGION}.${AWS_SERVICES_DOMAIN}/eks/pause:3.1-eksbuild.1 \
  --kubeconfig=<a href="#kubeconfig">/etc/eksctl/kubeconfig.yaml</a> \
  --config=<a href="#kubelet-configuration">/etc/eksctl/kubelet.yaml</a>

</pre>


### [Kubelet configuration](https://github.com/weaveworks/eksctl/blob/70041a226bb8ef5c51a229d587235551a2410eda/pkg/nodebootstrap/assets/kubelet.yaml)

The configuration for the kubelet is stored in `/etc/eksctl/kubelet.yaml`.

<pre>
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1

clusterDNS:
- 10.100.0.10

address: 0.0.0.0
clusterDomain: cluster.local
serverTLSBootstrap: true
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 2m0s
    enabled: true
  x509:
    clientCAFile: <a href="#" title="">/etc/eksctl/ca.crt</a>

authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 5m0s
    cacheUnauthorizedTTL: 30s

cgroupDriver: cgroupfs
kubeReserved:
  cpu: 60m
  ephemeral-storage: 1Gi
  memory: 343Mi

featureGates:
  RotateKubeletServerCertificate: true
</pre>

### Kubeconfig

The kubeconfig is stored in `/etc/eksctl/kubeconfig.yaml`.

<pre>
kind: Config
apiVersion: v1

clusters:
- cluster:
    certificate-authority: <a href="#" title="">/etc/eksctl/ca.crt</a>
    server: https://0A3....gr7.us-west-2.eks.amazonaws.com
  name: my-cluster-1.us-west-2.eksctl.io
contexts:
- context:
    cluster: my-cluster-1.us-west-2.eksctl.io
    user: kubelet@my-cluster-1.us-west-2.eksctl.io
  name: kubelet@my-cluster-1.us-west-2.eksctl.io
current-context: kubelet@my-cluster-1.us-west-2.eksctl.io

preferences: {}

users:
- name: kubelet@my-cluster-1.us-west-2.eksctl.io
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - eks
      - get-token
      - --cluster-name
      - my-cluster-1
      - --region
      - us-west-2
      command: aws
      env:
      - name: AWS_STS_REGIONAL_ENDPOINTS
        value: regional
</pre>

### Kubelet local environment variables

Stored in `/etc/eksctl/kubelet.local.env`, this file contains dynamic data and is generated by eksctl when generating
the CloudFormation template.

```
NODE_IP=192.168.72.51
INSTANCE_ID=i-0b704274a75a321a7
INSTANCE_TYPE=m6g.medium
AWS_SERVICES_DOMAIN=amazonaws.com
MAX_PODS=8
```

### Metadata environment variables

`/etc/eksctl/metadata.env`

```
AWS_DEFAULT_REGION=us-west-2
AWS_EKS_CLUSTER_NAME=my-cluster-1
AWS_EKS_ENDPOINT=https://0A123.....gr7.us-west-2.eks.amazonaws.com
AWS_EKS_ECR_ACCOUNT=602401143452

```

### Kubelet environment variables

`/etc/eksctl/kubelet.env`

```
NODE_LABELS=alpha.eksctl.io/cluster-name=my-cluster-1,alpha.eksctl.io/nodegroup-name=ng-1
NODE_TAINTS=
```

### CA File

`etc/eksctl/ca.crt`


```
-----BEGIN CERTIFICATE-----
1231231231231231231231231231231231231231231231231231231231231233
ABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCA
........
-----END CERTIFICATE-----
```


### Max pods file

`/etc/eksctl/max_pods.map`

```
t2.micro 4
x1e.32xlarge 234
c1.xlarge 58
g4dn.2xlarge 29
r5.24xlarge 737
r5a.12xlarge 234
r5ad.xlarge 58
t3.2xlarge 58
c1.medium 12
m5d.16xlarge 737
...
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
