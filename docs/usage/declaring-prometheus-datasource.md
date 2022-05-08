# Declaring a Prometheus data source

A Prometheus data source allows the integration of Prometheus into Grafana.

## Example usage

The following example will create a `my-prometheus` data source in Grafana:

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: Datasource
metadata:
  name: my-prometheus
spec:
  prometheus:
    url: "http://dark-prometheus-server:9090"

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
  prometheus:
    # Makes this data source the default one.
    default: false # Optional. Default value: false.

    # ---- #
    # HTTP #
    # ---- #

    # URL of the Prometheus server.
    # Required.
    url: "http://prometheus-server:9090"

    # Defined how Prometheus is accessed.
    # "proxy" means Grafana will call Prometheus and proxy the results to the user,
    # "direct" means the user directly accesses Prometheus.
    # Optional. Default: "proxy"
    access_mode: "proxy"

    # List of cookies (by name) that should be forwarded to Grafana.
    # Optional. Default: []
    forward_cookies: []

    # ---- #
    # Auth #
    # ---- #
    
    # Forward the user's upstream OAuth identity to Prometheus (their access token gets passed along)
    # Optional. Default: false
    forward_oauth: false

    # Whether credentials such as cookies or headers should be sent with cross-site requests.
    # Optional. Default: false
    forward_credentials: false

    # Disables SSL certificates verification.
    # Optional. Default: false
    skip_tls_verify: false

    # Needed to verify self-signed certificates.
    # Optional. Default: none
    ca_certificate:
      # CA certificate, as plain text. This is not recommended.
      # Optional. Default: ''
      value: ''

      # Reference to a secret containing the CA certificate.
      # Optional. Default: none
      valueFrom:
        secretKeyRef:
          name: 'secret-name' # name of the secret
          key: 'certificate' # Key within the secret

    # Enable basic authentication to the Prometheus server.
    # Optional. Default: none
    basic_auth:
      username:
        # Username, as plain text.
        # Optional. Default: ''
        value: ''
    
        # Reference to a secret containing the username.
        # Optional. Default: none
        valueFrom:
          secretKeyRef:
            name: 'secret-name' # name of the secret
            key: 'username' # Key within the secret
          
      password:
        # Password, as plain text. This is not recommended.
        # Optional. Default: ''
        value: ''

        # Reference to a secret containing the password.
        # Optional. Default: none
        valueFrom:
          secretKeyRef:
            name: 'secret-name' # name of the secret
            key: 'password' # Key within the secret

    # -------- #
    # Alerting #
    # -------- #

    # Set this to the typical scrape and evaluation interval configured in Prometheus.
    # Optional. Default: "15s"
    scrape_interval: "15s"

    # Prometheus query timeout.
    # Optional. Default: "60s"
    query_timeout: "60s"

    # HTTP method used to query Prometheus. POST is recommended since it allows for bigger queries.
    # Optional. Default: "POST"
    http_method: "POST" # Or "GET"

    # --------- #
    # Exemplars #
    # --------- #

    # Exemplars configuration.
    # Optional. Default: []
    exemplars:
      - label_name: "traceID" # Name of the field in the labels object that should be used to get the trace ID.
        
        # The URL of the trace backend the user would go to see its traces.
        # For external links only.
        # Optional. Default: ""
        url: "https://example.com/${__value.raw}"
        
        # The data source the exemplar is going to navigate to.
        # For internal links only.
        # Only one of the `uid` or `name` field may be specified.
        # Optional. Default: none
        datasource:
          # Data source UID.
          # Optional. Default: ""
          uid: ""

          # Data source name.
          name: "some-jaeger-source"
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)