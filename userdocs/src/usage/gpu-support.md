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

To disable the automatic plugin installation, and manually install a specific version,
use `--install-nvidia-plugin=false` with the create command. For example:

```
eksctl create cluster --node-type=p2.xlarge --install-nvidia-plugin=false

kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/<VERSION>/nvidia-device-plugin.yml
```
