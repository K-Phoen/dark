# Creating API keys

The `APIKey` manifest allows the creation of [Grafana API keys](https://grafana.com/docs/grafana/latest/administration/api-keys/create-api-key/),
that will be exposed as Kubernetes secrets.

Consider the following `APIKey`:

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: APIKey
metadata:
  name: my-service-key
  namespace: apps
spec:
  role: viewer # or: 'admin', 'editor' 
```

Check the result with:

```sh
kubectl get apikeys
```

Once successfully applied, a `my-service-key` Kubernetes secret will be created in the `apps` namespace.
The actual API key will be stored in it, under the `token` key:

```sh
kubectl get secrets my-service-editor-key -n apps --template="{{ .data.token | base64decode }}"
```

### Roles

Valid roles are:

* `admin`
* `editor`
* `viewer`

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)