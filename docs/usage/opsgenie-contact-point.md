# Defining `opsgenie` contact point types

## Example usage

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: AlertManager
metadata:
  name: alertmanager-example
spec:
  contact_points:
    - name: Team A

      # Contact Team A via Opsgenie
      contacts:
        - opsgenie:
            api_url: https://api.eu.opsgenie.com/v2/alerts
            api_key: { value: "shhhh, it's a secret" }
            auto_close: true
            override_priority: true

  routing:
    # ... omitted
```

## Reference

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: AlertManager
metadata:
  name: alertmanager-example
spec:
  # Alerts not matched by any of the routing rules will be sent to this contact point
  default_contact_point: 'Team A'

  # List of known contact points
  contact_points:
    - name: Team A
      contacts:
        - opsgenie:
            # Alert API Url.
            # Required.
            api_url: https://api.eu.opsgenie.com/v2/alerts

            # API key to use.
            # Required.
            api_key:
              # API key, as plain text. This is not recommended.
              # Optional. Default: ''
              value: ''

              # Reference to a secret containing the API key.
              # Optional. Default: none
              valueFrom:
                secretKeyRef:
                  name: 'secret-name' # name of the secret
                  key: 'key' # Key within the secret

            # Automatically close alerts in Opsgenie when they are closed in Grafana.
            # Optional. Default: false
            auto_close: false

            # Allow the alert priority to be set using the og_priority annotation.
            # Optional. Default: false
            override_priority: false
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)
