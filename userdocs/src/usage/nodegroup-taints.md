# Taints

To apply [taints](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) to a specific nodegroup use the `taints` config section like this:

```yaml
    taints:
      - key: your.domain.com/db
        value: "true"
        effect: NoSchedule
      - key: your.domain.com/production
        value: "true"
        effect: NoExecute
```

A full example can be found [here](https://github.com/eksctl-io/eksctl/blob/main/examples/34-taints.yaml).
