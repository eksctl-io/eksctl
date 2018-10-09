# Desing Proposal #001: Add-ons

## What's an add-ons?

An add-on extends functionality of a Kubernetes cluster, it may consist of a workload and/or configuration within the give cluster or the cloud provider. The workload, if present, would may classify as an operator or a controller, however that is not neccessary.

The difference between add-ons and anything a user decided to run themselves, is that from user's point of view, add-on is something they don't need to maintain directly. A cluster provider or bootstrap tool (i.e. eksctl), is expected to ensure ease of use, by providing minumum viable (that configuration to ensure end-to-end integration and compatibility), as well as catering for reconfiguration and upgrades.

## Design Requirements and Goals

Add-on must expose all of the following attributes:

- name (fixed) *[TBD]*
- version tag or revision (has default value that user can override)
- source URL (has default value that user can override)
- set of customisations (no default value, user defined)

Internally, an add-on may consist of any and/or all of the following elements:

- Kubernetes workload definitions
- Kubernetes configuration objects - PersistentVolumes, ConfigMaps, ...etc
- cloud provider resource definitons (e.g. CloudFormation sub-resources)
- dependencies (supported Kubernetes versions and corresponding cloud providers, list of references to other add-ons)
- incompatibilities (supported Kubernetes versions and corresponding cloud providers, list of reference to other add-ons)

In other words, it's ought to be possible to define an add-on that only includes:

- dependencies on other add-ons
- cloud provider resource definitions
- Kubernetes configuration objects
- Kubernetes workload definitions

The implementation must also ensure:

- fully deterministic behaviour for a given version of eksctl, thereby it should have
    - strict versioning of built-in add-ons by default
    - any externally-sourced config has to be pinned down or vendored *[TBD]*
    - there should be an option to use latest version or altenative source URL
- user should be able:
    - remove any add-on (including it's dependencies) at any time
    - customise an add-on in an arbitrary fashion
    - upgrade an add-on
    - install add-on at any time
- an add-on may consists exclusively of
    - a reference to any Helm chart (or a copy of one)
    - a reference to a Ksonnet application definition

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

Install a buil-in `flux` addon:

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonInstance
metadata:
    name: flux
params:
    newCodeCommitRepo: true
    witHelmOperator: true
```

Install a buil-in `helm` addon:

```YAML
apiVersion: eksctl.io/v1alpha1
kind: AddonInstance
metadata:
    name: helm
```

It will be possible to list built-in add-ons with `eksctl addons list`, and install them as simply as `eksctl create cluster --addons=<name>` or `eksctl addons install <name>`.
It should be possible to also specify parameters using a flag, the syntax is _TBD_.

Install a external example addon:

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