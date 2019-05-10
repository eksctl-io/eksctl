# Design Proposal #001: Add-ons

## What's an add-on?

An add-on extends functionality of a Kubernetes cluster, it may consist of a workload and/or configuration within the given cluster or the cloud provider. The workload, if present, may classify as an operator or a controller, however that is not necessary.

The difference between add-ons and anything a user decided to run themselves, is that from a user's point of view, an add-on is something they don't need to maintain directly. A cluster provider or bootstrap tool (i.e. eksctl), is expected to ensure ease of use, by providing a minimum viable solution (that configuration to ensure end-to-end integration and compatibility), as well as catering for reconfiguration and upgrades.

## Design Requirements and Goals

An add-on must expose all of the following attributes:

- name (fixed) *[TBD]*
- version tag or revision (has default value that the user can override)
- source URL (has default value that user can override)
- set of customisations (no default value, user defined)

Internally, an add-on may consist of any and/or all of the following elements:

- Kubernetes workload definitions
- Kubernetes configuration objects - PersistentVolumes, ConfigMaps, ...etc
- cloud provider resource definitions (e.g. CloudFormation sub-resources)
- dependencies (supported Kubernetes versions and corresponding cloud providers, list of references to other add-ons)
- incompatibilities (supported Kubernetes versions and corresponding cloud providers, list of reference to other add-ons)

In other words, it ought to be possible to define an add-on that only includes:

- dependencies on other add-ons
- cloud provider resource definitions
- Kubernetes configuration objects
- Kubernetes workload definitions

The implementation must also ensure:

- fully deterministic behaviour for a given version of eksctl, thereby it should have
    - strict versioning of built-in add-ons by default
    - any externally-sourced config has to be pinned down or vendored *[TBD]*
    - there should be an option to use the latest version or alternative source URL
- that a user is able to:
    - remove any add-on (including it's dependencies) at any time
    - customise an add-on in an arbitrary fashion
    - upgrade an add-on
    - install an add-on at any time
- an add-on may consist exclusively of
    - a reference to any Helm chart (or a copy of one)
    - a reference to a Ksonnet application definition

## Review of Prior Art

This section aims to briefly review some of the existing add-on management implementations.

### kube-addon-manager

***Check where does the code live, describe some of the functionality based on what the does***

This is a long-lived component, which started as a shell script (part of `kube-up.sh`). It's been used by a few
cluster providers (including GKE), yet it's behaviour is very opaque to most users. It's mostly a very simple
component, yet it's not designed to be consumed by a user and is primarily aimed to solve private needs of a
cluster provider.

### kubernetes/kops



### kubernetes/minikube

Add-ons are defined in `deploy/addons` and referenced in `pkg/minikube/assets/addons.go`. All of the manifests
are flat YAML files. The `kubernetes.io/minikube-addons` and `kubernetes.io/minikube-addons-endpoint` labels are
commonly used. Custom add-ons can be provided via `~/.minikube/addons`.

It uses `kube-addon-manager`. The manifests are compiled-in, and extract to `/etc/kubernetes/addons` at runtime,
which is mounted into `kube-addon-manager` pod.

Many of the add-ons are using `RepicationControllers`, and it's (probably) `kube-addon-manager` job to do rolling
update in old-fashioned way with `kubectl rolling-update`.

Commands:

- `minikube addons list`
- `minikube addons enable|disable <name>`
- `minikube addons open <name>` forwards service endpoint and opens it in a local browser

Built-in add-ons include kube-dns, coredns, dashboard, registry, heapster and freshpod (to name a few).

With this design in mind, upgrades are either unsupported, or implicit on cluster upgrade or restart, or on CLI
upgrade. In any case, `kube-addon-manager` handles upgrades transparently, however user has no control at all.
It's also not clear if an add-on can be customised by the user without having to taken on maintaining the add-on
a whole.

### kubernetes-sigs/cluster-api

In the Cluster API project, the notion of add-on is very simple – it's anything that can be used with `kubectl apply`,
i.e. one or more manifest bundled in a file or a directory. Add-ons are applied once API server is up, there isn't
anything else to it. Definition of cluster readiness.

### kubeadm

***This needs to be reviewed, it was written from memory***

Add-ons in kubeadm are much more tightly coupled than elsewhere, and the only add-ons are the DNS component (orignally
it was kube-dns, and more recently CoreDNS was adopted; both choices are still offered), and kube-proxy.
The implementation varies between add-on, one is defined in terms of `k8s.io/api` structs directly, and the other is
defined as an in-line string that is used with `text/template`. Originally, both add-ons used structs, but that has
changed at one point, when maintainers concluded that DNS component was easier to maintain by copying YAML from another
(canonical) directory (in the same `k8s.io/kubernetes` repository).
The add-ons are upgraded during cluster upgrades. There is a way for user to create a cluster without the two add-ons,
if they wished to do so.

### GoogleCloudPlatform/k8s-cluster-bundle (TBC)

***More technical information is needed***

This project aims at defining an API for defining custom collections of manifests.

### Conclusion

Most of the implementations described above have custom add-on management mechanisms. The goal of this proposal is
to provide a design that would work for eksctl short-term (MVP), but also suggest a direction for an implementation
that could be applied more widely and adopted by other projects in the future.
While majority of add-on manifests are based on a canonical source, the way copies are updates is different for each
of the implementations, and it's hard for user to find out whether any of such copies are up-to-date date or not.
Some of the implementations set set custom labels, but it's not clear what customisations had been applied to each
copy without deeper case-by-case analysis.
These findings suggest that there is a need for a common add-on management solution that integrates well with various
community projects, including Cluster API, as well as any configuration/package management solutions, e.g. Kustomize,
Helm, Ksonnet and others.

## Phases of Development

1. discuss and specify add-on definition objects and repository layout
2. discuss and specify how different config management tools would integrate (Helm, K/Jsonnet, kubegen, kustomize)
3. implement ability to
    - install a built-in add-on (must ensure upgrades are catered for)
    - install an external add-on
    - reconfigure customer an installed add-on
    - upgrade any installed add-on

## Design Recommendations

- *[TBD]*

- built-in add-ons:
    - only one version of each add-on is expected to be built-in, alternative versions have to be sourced externally
    - only built-in add-ons may include Go code to create AWS resources (this may change in the future, yet significantly simplifies initial design and ensures security)

## API

### Add-on instance examples

Install a built-in `flux` add-on:

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonInstance
metadata:
    name: flux
params:
    newCodeCommitRepo: true
    witHelmOperator: true
```

Install a built-in `helm` addon:

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonInstance
metadata:
    name: helm
```

It will be possible to list built-in add-ons with `eksctl addons list`, and install them as simply as `eksctl create cluster --addons=<name>` or `eksctl addons install <name>`.
It should be possible to also specify parameters using a flag, the syntax is _TBD_.

Install a external example add-on:

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonInstance
metadata:
    name: kustomize-example
spec:
    source:
        repo: https://github.com/weaveworks/eksctl-addons
        branch: $HEAD
        path: examples/kustomize
    kustomization:
        # by default the base is the add-on module
        # TBD: allow user to append bases
        patchesStrategicMerge:
        - patch.yaml
        namePrefix: example-
```

### Add-on module examples

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonModule
metadata:
    name: flux
spec:
    source:
        repo: https://github.com/weaveworks/eksctl
        branch: $HEAD
        path: addons/flux
    helmTemplate:
        chartDir: addons/flux/chart
```

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonModule
metadata:
    name: ksonnet-example
spec:
    source:
        repo: https://github.com/weaveworks/eksctl-addons
        branch: $HEAD
        path: examples/ksonnet
    ksonnet: {}
```

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonModule
metadata:
    name: kustomize-example
spec:
    source:
        repo: https://github.com/weaveworks/eksctl-addons
        branch: $HEAD
        path: examples/kustomize
    kustomize:
        resources:
        - deployment.yaml
        - service.yaml
        configMapGenerator:
        - name: configmap
          files:
          - config.env
```

TBD: it may be possible to chain different tools, e.g. render a Helm chart (via `helm template`), and pass it to `kustomize`.

### Parameters

A way to define parameters in a portable way will be required for optimal user (and/or add-on author) experience.

At the add-on management layer, the following key aspects of parameters are:

- parameters are a map of strings to a value of a primitive type
- initial set of primitive types - `bool`, `number` and `string`
- non-primive types (`object` and `array`) are initially out of scope (yet could probably be handled at the underlying provider level)
- a parameter may be either optional with a default value or non-optional without a default value

#### How would this map to different providers?

Helm values are essentially `map[string]interface{}`, which implies there is no type information or required params, and default values are defined inside a chart. The management layer could take care of the concerns above and pass an appropriate set of value.

Ksonnet – TBD

With kustomize, there are no params, but they can be introduced with auto-generate ConfigMap and `vars` stanza. So potentially it should be doable, given the underlying add-on base layer was written with parameters in mind, otherwise kustomizations would be required as another layer where params are referenced.

#### Alternatives

- only backend-specific params, i.e.
  - an add-on based on Helm chart would use Helm values format
  - a ksonnet-based add-on will use its own
  - a kustomize-based add-on won't have any

In Helm, we would be ...TBD
