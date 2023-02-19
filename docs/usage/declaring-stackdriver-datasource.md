# Declaring a Stackdriver (Google Cloud Monitoring) data source

A Stackdriver data source allows the integration of the Google Cloud Monitoring into Grafana.

## Example usage

The following example will create a `my-stackdriver` data source in Grafana:

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: Datasource
metadata:
  name: my-stackdriver
spec:
  # Uses the default GCE service account for authentication.
  stackdriver: {}
```

Check the result with:

```sh
kubectl get datasources
```

## Reference

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: Datasource
metadata:
  name: datasource-name
spec:
  # Will use the default GCE service account for authentication by default
  stackdriver:
    # Makes this data source the default one.
    default: false # Optional. Default value: false.

    # Service account key to use.
    # Optional. Default: none (the default GCE service account is used)
    jwt_authentication:
      # Service account key to use, as plain text. This is not recommended.
      # Optional. Default: ''
      value: ''

      # Reference to a secret containing the service account key to use.
      # Optional. Default: none
      valueFrom:
        secretKeyRef:
          name: 'secret-name' # name of the secret
          key: 'sa.json' # Key within the secret
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)
