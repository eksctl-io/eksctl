# GPU Support

Eksctl supports selecting GPU instance types for nodegroups. Simply supply a
compatible instance type to the create command, or via the config file.

```
eksctl create cluster --node-type=p2.xlarge
```

!!! note
    It is no longer necessary to subscribe to the marketplace AMI for GPU support on EKS.

The AMI resolvers (`auto` and `auto-ssm`) will see that you want to use a
GPU instance type and they will select the correct EKS optimized accelerated AMI.

Eksctl will detect that an AMI with a GPU-enabled instance type has been selected and
will install the [NVIDIA Kubernetes device plugin](https://github.com/NVIDIA/k8s-device-plugin) automatically.

!!! note
    Windows and Ubuntu AMIs do not ship with GPU drivers installed, hence running GPU-accelerated workloads will not work out of the box.

To disable the automatic plugin installation, and manually install a specific version,
use `--install-nvidia-plugin=false` with the create command. For example:

```
eksctl create cluster --node-type=p2.xlarge --install-nvidia-plugin=false

kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/<VERSION>/nvidia-device-plugin.yml
```

The installation of the [NVIDIA Kubernetes device plugin](https://github.com/NVIDIA/k8s-device-plugin) will be skipped if the cluster only includes Bottlerocket nodegroups, since Bottlerocket already handles the execution of the device plugin.
If you use different AMI families in your cluster's configurations, you may need to use taints and tolerations to keep the device plugin from running on Bottlerocket nodes.
