# Declaring a Tempo data source

A Tempo data source allows the integration of [Tempo](https://grafana.com/oss/tempo/) into Grafana.

## Example usage

The following example will create a `my-tempo` data source in Grafana:

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: Datasource
metadata:
  name: my-tempo
spec:
  tempo:
    url: "http://tempo-server"
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
  tempo:
    # Makes this data source the default one.
    default: false # Optional. Default value: false.

    # ---- #
    # HTTP #
    # ---- #

    # URL of the Tempo server.
    # Required.
    url: "http://tempo-server:3100"

    # HTTP request timeout.
    # Optional. Default: ""
    timeout: "60s"

    # List of cookies (by name) that should be forwarded to Grafana.
    # Optional. Default: []
    forward_cookies: []

    # ---- #
    # Auth #
    # ---- #
    
    # Forward the user's upstream OAuth identity to Tempo (their access token gets passed along)
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

    # Enable basic authentication to the Tempo server.
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

    # ---------- #
    # Node Graph #
    # ---------- #

    # Enables the Node Graph visualization in the trace viewer.
    # Optional. Default: false
    node_graph: false

    # ------------- #
    # Trace to logs #
    # ------------- #

    # Trace to logs lets you navigate from a trace span to the selected data source's log.
    # Optional. Default: none
    trace_to_logs:
      - tags: [] # Tags that will be used in the Loki query. Optional. Default: [cluster, hostname, pod, namespace].

        # The data source the trace is going to navigate to.
        # Optional. Default: none
        datasource:
          # Data source UID.
          # Optional. Default: ""
          uid: ""

          # Data source name.
          name: "some-data-source"

        # Shifts the start time of the span.
        # Optional. Default: none
        span_start_shift: "0h"

        # Shifts the end time of the span.
        # Optional. Default: none
        span_end_shift: "0h"

        # Filters logs by trace ID.
        # Optional. Default: false
        filter_by_trace: false

        # Filters logs by span ID.
        # Optional. Default: false
        filter_by_span: false
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)