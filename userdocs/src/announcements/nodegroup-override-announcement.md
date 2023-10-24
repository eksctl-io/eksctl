# Nodegroup Bootstrap Override For Custom AMIs

This change was announced in the issue [Breaking: overrideBootstrapCommand soon...](https://github.com/eksctl-io/eksctl/issues/3563).
Now, it has come to pass in [this](https://github.com/eksctl-io/eksctl/pull/4968) PR. Please read the attached issue carefully about
why we decided to move away from supporting custom AMIs without bootstrap scripts or with partial bootstrap scripts.

We still provide a helper! Migrating hopefully is not that painful. `eksctl` still provides a script, which when sourced,
will export a couple of helpful environment properties and settings. This script is located [here](https://github.com/eksctl-io/eksctl/blob/70a289d62e3c82e6177930cf2469c2572c82e104/pkg/nodebootstrap/assets/scripts/bootstrap.helper.sh).

The following environment properties will be at your disposal:

```bash
API_SERVER_URL
B64_CLUSTER_CA
INSTANCE_ID
INSTANCE_LIFECYCLE
CLUSTER_DNS
NODE_TAINTS
MAX_PODS
NODE_LABELS
CLUSTER_NAME
CONTAINER_RUNTIME # default is docker
KUBELET_EXTRA_ARGS # for details, look at the script
```

The minimum that needs to be used when overriding so `eksctl` doesn't fail, is labels! `eksctl` relies on a specific set of
labels to be on the node, so it can find them. When defining the override, please provide this **bare minimum** override
command:

```yaml
    overrideBootstrapCommand: |
      #!/bin/bash

      source /var/lib/cloud/scripts/eksctl/bootstrap.helper.sh

      # Note "--node-labels=${NODE_LABELS}" needs the above helper sourced to work, otherwise will have to be defined manually.
      /etc/eks/bootstrap.sh ${CLUSTER_NAME} --container-runtime containerd --kubelet-extra-args "--node-labels=${NODE_LABELS}"
```

For nodegroups that have no outbound internet access, you'll need to supply `--apiserver-endpoint` and `--b64-cluster-ca`
to the bootstrap script as follows:

```yaml
    overrideBootstrapCommand: |
      #!/bin/bash

      source /var/lib/cloud/scripts/eksctl/bootstrap.helper.sh

      # Note "--node-labels=${NODE_LABELS}" needs the above helper sourced to work, otherwise will have to be defined manually.
      /etc/eks/bootstrap.sh ${CLUSTER_NAME} --container-runtime containerd --kubelet-extra-args "--node-labels=${NODE_LABELS}" \
        --apiserver-endpoint ${API_SERVER_URL} --b64-cluster-ca ${B64_CLUSTER_CA}
```

Note the _`--node-labels`_ setting. If this is not defined, the node will join the cluster, but `eksctl` will ultimately
time out on the last step when it's waiting for the nodes to be `Ready`. It's doing a Kubernetes lookup for nodes that
have the label `alpha.eksctl.io/nodegroup-name=<cluster-name>`. This is only true for unmanaged nodegroups. For managed
it's using a different label.

If, at all, it's possible to switch to managed nodegroups to avoid this overhead, the time has come now to do that. Makes
all the overriding a lot easier.
