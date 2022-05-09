# Declaring a Loki data source

A Loki data source allows the integration of [Loki](https://grafana.com/oss/loki/) into Grafana.

## Example usage

The following example will create a `my-loki` data source in Grafana:

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: Datasource
metadata:
  name: my-loki
spec:
  loki:
    url: "http://loki-server:9090"

    derived_fields:
      - name: TraceID
        regex: '(?:traceID|trace_id)=(\w+)'
        url: '${__value.raw}'
        datasource:
          name: my-tempo
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
  loki:
    # Makes this data source the default one.
    default: false # Optional. Default value: false.

    # ---- #
    # HTTP #
    # ---- #

    # URL of the Loki server.
    # Required.
    url: "http://loki-server:3100"

    # HTTP request timeout.
    # Optional. Default: ""
    timeout: "60s"

    # List of cookies (by name) that should be forwarded to Grafana.
    # Optional. Default: []
    forward_cookies: []

    # ---- #
    # Auth #
    # ---- #
    
    # Forward the user's upstream OAuth identity to Loki (their access token gets passed along)
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

    # Enable basic authentication to the Loki server.
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

    # Loki queries must contain a limit of the maximum number of lines returned.
    # Optional. Default: 1000
    maximum_lines: 1000

    # Derived fields can be used to extract new fields from a log message and create a link from its value.
    # Optional. Default: none
    derived_fields:
      - name: "field" # Required. Name of the field.

        # Used to parse and capture some part of the log message. You can use the captured groups in the template.
        # Required.
        regex: ""

        # Required.
        url: "https://example.com/${__value.raw}"

        # Used to override the button label when this derived field is found in a log.
        # Optional. Default: none
        url_label:
        
        # For internal links only.
        # Optional. Default: none
        datasource:
          # Data source UID.
          # Optional. Default: ""
          uid: ""

          # Data source name.
          name: "some-data-source"
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)