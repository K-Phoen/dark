# Installing the operator

**Important:** be sure to review any manifest before applying them.

DARK needs a few things to be up and running:

* a [`CRD`](../../config/crd)
* a [`ServiceAccount`](../../config/rbac/service_account.yaml) with [some permissions](../../config/rbac/role.yaml)
* a [`Deployment`](../../config/operator)

## Using kustomize

Kustomize can be used to fit these manifests to your needs.
Write your own `kustomization.yaml` using ours as a "base" and write patches to tweak the configuration.

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - https://github.com/K-Phoen/dark/config/crd
  - https://github.com/K-Phoen/dark/config/rbac
  - https://github.com/K-Phoen/dark/config/operator

namespace: monitoring
```

Generate the manifests:

```shell
kubectl kustomize . -o dark-all.yaml
```

Review their content, and deploy:

```shell
kubectl apply -f dark-all.yaml
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)
