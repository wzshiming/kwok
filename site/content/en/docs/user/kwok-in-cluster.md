# Deploy `kwok` in a Cluster

{{< hint "info" >}}

This document walks you through how to deploy `kwok` in a Kubernetes cluster.

{{< /hint >}}

## Preparations

``` bash
# Temporary directory
KWOK_WORK_DIR=$(mktemp -d)
```

## Render kustomization yaml

Firstly, generate a kustomization template yaml to the previously-defined temporary directory.

``` bash
cat <<EOF > "${KWOK_WORK_DIR}/kustomization.yaml"
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
images:
  - name: registry.k8s.io/kwok/kwok
    newTag: "v{{< variable "version" >}}"
resources:
  - "https://github.com/{{< variable "repo" >}}/kustomize/kwok?ref=v{{< variable "version" >}}"
EOF
```

Next, render it with the prepared variables.

``` bash
kubectl kustomize "${KWOK_WORK_DIR}" > "${KWOK_WORK_DIR}/kwok.yaml"
```

## `kwok` deployment

Finally, we're able to deploy `kwok`:

``` bash
kubectl apply -f "${KWOK_WORK_DIR}/kwok.yaml"
```

## Next steps

Now, you can use `kwok` to [manage nodes and pods] in the Kubernetes cluster.

[manage nodes and pods]: {{< relref "/docs/user/kwok-manage-nodes-and-pods" >}}
