# Declaring a CloudWatch data source

A CloudWatch data source allows the integration of Amazon CloudWatch service into Grafana.

## Example usage

The following example will create a `my-cloudwatch` data source in Grafana:

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: Datasource
metadata:
  name: my-cloudwatch
spec:
  # Uses AWS SDK auth strategy by default
  cloudwatch: {}
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
  # Uses AWS SDK auth strategy by default
  cloudwatch:
    # Makes this data source the default one.
    default: false # Optional. Default value: false.

    # Custom endpoint for the CloudWatch service.
    endpoint: '' # Optional. Default value: ''.

    # Default region to use.
    default_region: '' # Optional. Default value: ''.

    # ARN of a role to assume.
    assume_role_arn: '' # Optional. Default value: ''.

    # External identifier of a role to assume in another account.
    external_id: '' # Optional. Default value: ''.

    # List of namespaces for custom metrics.
    custom_metrics_namespaces: [] # Optional. Default value: [].

    # Authentication mode to use.
    # Optional. Default: none (the default AWS SDK strategy is used)
    auth:
      # Access + secret keys
      keys:
        # Access key.
        access: ''

        # Secret key.
        secret:
          # Secret key to use, as plain text. This is not recommended.
          # Optional. Default: ''
          value: ''

          # Reference to a secret containing the secret key.
          # Optional. Default: none
          valueFrom:
            secretKeyRef:
              name: 'secret-name' # name of the secret
              key: 'cloudwatch_secret_key' # Key within the secret
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)
